package plugins

type FileTemplate struct {
	PackageName  string
	FileFullName string
	Tpl          string
}

func NewFileTemplate(pkgName, fullName string) *FileTemplate {
	return &FileTemplate{
		PackageName:  pkgName,
		FileFullName: fullName,
	}
}

func (t *FileTemplate) WithBlock(tpl string) *FileTemplate {
	t.Tpl += tpl
	return t
}
