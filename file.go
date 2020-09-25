package plugins

type FileTemplate struct {
	FileName string
	Tpl      string
}

func NewFileTemplate(name string) *FileTemplate {
	return &FileTemplate{
		FileName: name,
	}
}

func (t *FileTemplate) WithBlock(tpl string) *FileTemplate {
	t.Tpl += tpl
	return t
}
