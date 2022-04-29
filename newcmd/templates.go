package newcmd

type genFileInfo struct {
	name         string
	resourcePath string
	targetPath   string
}

var genFileInfos = []*genFileInfo{
	{
		"build script",
		"resources/newmod/build.cs.tmpl",
		"{{.ModuleName}}.build.cs",
	},
	{
		"log header",
		"resources/newmod/log.h.tmpl",
		"Private/Log.h",
	},
	{
		"log source",
		"resources/newmod/log.cpp.tmpl",
		"Private/Log.cpp",
	},
	{
		"module header",
		"resources/newmod/module.h.tmpl",
		"Public/{{.ModuleName}}.h",
	},
	{
		"module source",
		"resources/newmod/module.cpp.tmpl",
		"Private/{{.ModuleName}}.cpp",
	},
}

const projectJsonTmpl = `{{.FormatPrefix}}
		{
			"Name": "{{.ModuleName}}",
			"Type": "{{.ModuleType}}",
			"LoadingPhase": "{{.LoadingPhase}}"
		}{{if .HasOtherModules}},
		{{else}}
	{{end}}{{.FormatSuffix}}`

type projectJsonFormatContext struct {
	ModuleName      string
	HasOtherModules bool
	ModuleType      string
	LoadingPhase    string
	FormatPrefix    string
	FormatSuffix    string
}
