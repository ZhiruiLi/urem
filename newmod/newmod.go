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
	"unicode"

	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/osutil"
	"github.com/zhiruili/urem/regensln"
)

type Cmd struct {
	Copyright    string `arg:"-c,--copyright" help:"copyright owner"`
	LoadingPhase string `arg:"-l,--loading-phase" help:"module loading phase" default:"Default"`
	ModuleName   string `arg:"positional,required" help:"name of the new module"`
	OutputPath   string `arg:"positional,required" help:"module file output dir"`
}

// https://docs.unrealengine.com/4.26/en-US/API/Runtime/Projects/ELoadingPhase__Type/
var legalLoadingPhases = []string{
	"EarliestPossible",
	"PostConfigInit",
	"PostSplashScreen",
	"PreEarlyLoadingScreen",
	"PreLoadingScreen",
	"PreDefault",
	"Default",
	"PostDefault",
	"PostEngineInit",
}

func (cmd *Cmd) getModulePath() string {
	return filepath.Join(cmd.OutputPath, cmd.ModuleName)
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

func formatProjectJsonText(orignalJson string, moduleName string, loadingPhase string) string {
	ctx := projectJsonFormatContext{
		ModuleName:   moduleName,
		LoadingPhase: loadingPhase,
	}

	moduleTagIdx := strings.Index(orignalJson, "\"Modules\"")
	if moduleTagIdx < 0 {
		closeParansIdx := strings.LastIndex(orignalJson, "}")
		if closeParansIdx < 0 {
			return orignalJson
		}
		ctx.FormatPrefix = strings.TrimRight(orignalJson[:closeParansIdx], " \t\r\n") + ",\n\t\"Modules\": ["
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
		ctx.FormatSuffix = strings.TrimLeft(orignalJson[afterOpenBracketIdx:], " \t\r\n")
		endBracketIdx := strings.Index(ctx.FormatSuffix, "]")
		ctx.HasOtherModules = strings.TrimSpace(ctx.FormatSuffix[:endBracketIdx]) != ""
	}

	tmplEngine := template.Must(template.New("ProjectJSON").Parse(projectJsonTmpl))
	out := new(bytes.Buffer)
	tmplEngine.Execute(out, &ctx)
	return out.String()
}

func (cmd *Cmd) updateProjectJson(modulePath string) error {
	filePath, err := osutil.FindFileBottomUp(modulePath, "*.uproject", "*.uplugin")
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

	updated := formatProjectJsonText(string(content), cmd.ModuleName, cmd.LoadingPhase)
	if err := ioutil.WriteFile(filePath, []byte(updated), 0644); err != nil {
		return fmt.Errorf("write file %s", filePath)
	}

	return nil
}

func (cmd *Cmd) refreshSln(modulePath string) error {
	filePath, err := osutil.FindFileBottomUp(modulePath, "*.uproject")
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

func checkModuleName(name string) error {
	runes := []rune(name)
	if len(runes) == 0 {
		return core.IllegalArgErrorf("ModuleName", "empty string")
	}

	if !unicode.IsUpper(runes[0]) {
		cont := core.GetUserBoolInput("Unconventional module name, should start with upper case, continue?")
		if !cont {
			return fmt.Errorf("user cancel")
		}
	}

	return nil
}

func checkOutputPath(outPath string) error {
	baseName := filepath.Base(outPath)
	if baseName != "Source" {
		cont := core.GetUserBoolInput("Unconventional output path, should under Source dir, continue?")
		if !cont {
			return fmt.Errorf("user cancel")
		}
	}

	return nil
}

func checkLoadingPhase(phase string) error {
	if !core.StrContains(legalLoadingPhases, phase) {
		return core.IllegalArgErrorf("LoadingPhase", "illegal enum value, must be oneof: %s", strings.Join(legalLoadingPhases, ", "))
	}

	return nil
}

func (cmd *Cmd) CheckArgs() error {
	if core.Global.Quite {
		return nil
	}

	if absPath, err := filepath.Abs(cmd.OutputPath); err != nil {
		return core.IllegalArgErrorf("OutputPath", "illegal path")
	} else {
		cmd.OutputPath = absPath
	}

	if err := checkModuleName(cmd.ModuleName); err != nil {
		return err
	}

	if err := checkOutputPath(cmd.OutputPath); err != nil {
		return err
	}

	if err := checkLoadingPhase(cmd.LoadingPhase); err != nil {
		return err
	}

	return nil
}

func (cmd *Cmd) Run() (err error) {
	if err = cmd.CheckArgs(); err != nil {
		return err
	}

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
