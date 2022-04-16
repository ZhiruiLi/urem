package newmod

import (
	"bytes"
	"embed"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/zhiruili/usak/core"
	"github.com/zhiruili/usak/osutil"
	"github.com/zhiruili/usak/regensln"
)

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
			"Type": "Runtime",
			"LoadingPhase": "Default"
		}{{if .HasOtherModules}},
		{{else}}
	{{end}}{{.FormatSuffix}}`

type projectJsonFormatContext struct {
	ModuleName      string
	HasOtherModules bool
	FormatPrefix    string
	FormatSuffix    string
}

func formatProjectJsonText(orignalJson string, moduleName string) string {
	ctx := projectJsonFormatContext{
		ModuleName: moduleName,
	}

	moduleTagIdx := strings.Index(orignalJson, "\"Modules\"")
	if moduleTagIdx < 0 {
		closeParansIdx := strings.LastIndex(orignalJson, "}")
		if closeParansIdx < 0 {
			return orignalJson
		}
		ctx.FormatPrefix = strings.TrimRight(orignalJson[:closeParansIdx], "\t \n") + ",\n\t\"Modules\": ["
		ctx.FormatSuffix = "]\n" + orignalJson[closeParansIdx:]
		ctx.HasOtherModules = false
	} else {
		afterModelTag := orignalJson[moduleTagIdx:]
		modelOpenBracketIdx := strings.Index(afterModelTag, "[")
		if modelOpenBracketIdx < 0 {
			return orignalJson
		}
		afterOpenBracketIdx := moduleTagIdx + modelOpenBracketIdx + 1
		ctx.FormatPrefix = orignalJson[:afterOpenBracketIdx]
		ctx.FormatSuffix = strings.TrimLeft(orignalJson[afterOpenBracketIdx:], "\t \n")
		endBracketIdx := strings.Index(ctx.FormatSuffix, "]")
		ctx.HasOtherModules = strings.TrimSpace(ctx.FormatSuffix[:endBracketIdx]) != ""
	}

	tmplEngine := template.Must(template.New("ProjectJSON").Parse(projectJsonTmpl))
	out := new(bytes.Buffer)
	tmplEngine.Execute(out, &ctx)
	return out.String()
}

type Cmd struct {
	ModuleName string `arg:"positional,required"`
	Output     string `arg:"-o,--out" help:"module file output dir" default:"."`
}

func (cmd *Cmd) getModulePath() string {
	return filepath.Join(cmd.Output, cmd.ModuleName)
}

func (cmd *Cmd) generateFile(info *genFileInfo, modulePath string, fs *embed.FS) (string, error) {
	fileContentTmpl, err := fs.ReadFile(info.resourcePath)
	if err != nil {
		return "", fmt.Errorf("load resouce %s: %w", info.resourcePath, err)
	}

	fileContentTmplEngine := template.Must(template.New("File " + info.name).Parse(string(fileContentTmpl)))
	fileContent := new(bytes.Buffer)
	if err := fileContentTmplEngine.Execute(fileContent, cmd); err != nil {
		core.LogD("resource file %s content:\n%s\n", info.resourcePath, fileContentTmpl)
		return "", fmt.Errorf("format resource content, %w", err)
	}

	filePathTmplEngine := template.Must(template.New("Path " + info.name).Parse(info.targetPath))
	filePathBs := new(bytes.Buffer)
	if err := filePathTmplEngine.Execute(filePathBs, cmd); err != nil {
		core.LogD("resource file %s target path:\n%s\n", info.resourcePath, info.targetPath)
		return "", fmt.Errorf("format target path: %w", err)
	}

	filePath := filepath.Join(modulePath, filePathBs.String())
	fileDir := filepath.Dir(filePath)
	if err := os.MkdirAll(fileDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("create dir %s for file %s", fileDir, filePath)
	}

	if err := ioutil.WriteFile(filePath, fileContent.Bytes(), 0644); err != nil {
		return "", fmt.Errorf("write file %s", filePath)
	}

	return filePath, nil
}

func (cmd *Cmd) updateProjectJson(modulePath string) error {
	filePath, err := osutil.FindFileBottomUp(modulePath, ".uproject", ".uplugin")
	if err != nil {
		return fmt.Errorf("find .uproject or .uplugin file: %w", err)
	}

	if filePath == "" {
		return fmt.Errorf(".uproject or .uplugin file no found")
	}

	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("read file %s", filePath)
	}

	updated := formatProjectJsonText(string(content), cmd.ModuleName)
	if err := ioutil.WriteFile(filePath, []byte(updated), 0644); err != nil {
		return fmt.Errorf("write file %s", filePath)
	}

	return nil
}

func (cmd *Cmd) refreshSln(modulePath string) error {
	filePath, err := osutil.FindFileBottomUp(modulePath, ".uproject")
	if err != nil {
		return fmt.Errorf("find project file: %w", err)
	}
	if filePath == "" {
		return fmt.Errorf(".uproject file no found")
	}
	core.LogD("project file path: %s", filePath)

	rCmd := regensln.Cmd{
		ProjectFile: filePath,
	}
	return rCmd.Run()
}

func (cmd *Cmd) Run() (err error) {
	modulePath := cmd.getModulePath()

	if err = osutil.MkDirIfNotExisted(modulePath); err != nil {
		err = fmt.Errorf("make module dir: %w", err)
		return err
	}

	for _, info := range genFileInfos {
		genFilePath := ""
		if genFilePath, err = cmd.generateFile(info, modulePath, &core.Global.EmbedFs); err != nil {
			err = fmt.Errorf("generate file %s: %w", info.name, err)
			goto ERREND
		}
		core.LogI("generate %s file at %s", info.name, genFilePath)
	}

	if err = cmd.updateProjectJson(modulePath); err != nil {
		err = fmt.Errorf("update project or plugin JSON file: %w", err)
		goto ERREND
	}

	if err := cmd.refreshSln(modulePath); err != nil {
		core.LogE("refresh solution files: %s", err.Error())
	}

	return nil

ERREND:
	if core.Global.Debug {
		return err
	}
	if yes, _ := osutil.IsDir(modulePath); yes {
		os.RemoveAll(modulePath)
	}
	return err
}
