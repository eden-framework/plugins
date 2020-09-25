package plugins

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path"
	"plugin"
	"strings"
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

	defer fmt.Printf("\nDownload finished.\n")

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

	r, err := zip.OpenReader(zipFile)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var pluginFile string
	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".so") && !f.FileInfo().IsDir() {
			reader, err := f.Open()
			if err != nil {
				return nil, err
			}

			pluginFile = path.Join(l.cwd, f.FileInfo().Name())
			file, err := os.Create(pluginFile)
			if err != nil {
				reader.Close()
				return nil, err
			}
			l.fileList = append(l.fileList, file)

			_, err = io.Copy(file, reader)
			if err != nil {
				reader.Close()
				return nil, err
			}

			break
		}
	}

	if pluginFile == "" {
		err := fmt.Errorf("cannot find .so file in zip file: %s", name)
		return nil, err
	}

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
