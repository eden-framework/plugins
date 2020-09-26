package plugins

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path"
	"plugin"
)

type Loader struct {
	cwd      string
	fileList []*os.File
}

func NewLoader(projectCwd string) *Loader {
	ldr := &Loader{
		cwd: projectCwd,
	}
	return ldr
}

func (l *Loader) Load(name, zipUrl string) (*plugin.Plugin, error) {
	fmt.Printf("Start download plugin: %s\n", name)
	fmt.Printf("plugin zip url: %s\n", zipUrl)

	resp, err := http.Get(zipUrl)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	zipFile := path.Join(l.cwd, fmt.Sprintf("%s.zip", name))
	w, err := os.Create(zipFile)
	if err != nil {
		return nil, err
	}
	l.fileList = append(l.fileList, w)

	counter := NewWriteCounter(uint64(resp.ContentLength))
	_, err = io.Copy(w, io.TeeReader(resp.Body, counter))
	if err != nil {
		return nil, err
	}
	w.Close()
	fmt.Printf("\nDownload finished.\n")

	fmt.Println("Decompressing zip...")
	rootPath, err := l.Decompress(zipFile, l.cwd)
	if err != nil {
		return nil, err
	}

	fmt.Println("Building plugin...")
	cmd := exec.Command("make")
	cmd.Dir = rootPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return nil, err
	}

	fmt.Println("Loading plugin...")
	pluginFile := path.Join(rootPath, "generator.so")
	return plugin.Open(pluginFile)
}

func (l *Loader) Clear() error {
	fmt.Println("Clear loader cache files...")
	defer fmt.Printf("Clear finished.\n")

	for _, f := range l.fileList {
		fmt.Printf("Removing file: %s ...\n", f.Name())
		err := os.Remove(f.Name())
		if err != nil {
			return err
		}
	}

	return nil
}

func (l *Loader) Decompress(zipFile, outputPath string) (rootPath string, err error) {
	reader, err := zip.OpenReader(zipFile)
	if err != nil {
		return
	}
	defer reader.Close()

	for _, file := range reader.File {
		filename := path.Join(outputPath, file.Name)
		if file.FileInfo().Name() == "Makefile" && !file.FileInfo().IsDir() && rootPath == "" {
			rootPath = path.Dir(filename)
		}

		if file.FileInfo().IsDir() {
			err = os.MkdirAll(filename, 0755)
			if err != nil {
				return rootPath, err
			}
			continue
		} else {
			err = os.MkdirAll(path.Dir(filename), 0755)
			if err != nil {
				return rootPath, err
			}
		}

		w, err := os.Create(filename)
		if err != nil {
			return rootPath, err
		}
		rc, err := file.Open()
		if err != nil {
			w.Close()
			return rootPath, err
		}
		_, err = io.Copy(w, rc)
		if err != nil {
			w.Close()
			rc.Close()
			return rootPath, err
		}
		w.Close()
		rc.Close()
	}

	if rootPath == "" {
		err = fmt.Errorf("cannot find Makefile in plugins")
	}

	root, err := os.Open(rootPath)
	if err != nil {
		err = fmt.Errorf("open rootPath %s failed", rootPath)
		return
	}
	l.fileList = append(l.fileList, root)
	return
}
