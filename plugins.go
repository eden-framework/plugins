package plugins

type Option struct {
	PackageName string
}

type EntryPointPlugins interface {
	GenerateEntryPoint(opt Option, cwd string) string
}

type FilePlugins interface {
	GenerateFilePoint(opt Option, cwd string) []*FileTemplate
}
