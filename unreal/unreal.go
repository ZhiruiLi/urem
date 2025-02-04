package unreal

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/iancoleman/orderedmap"

	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/pwsh"
)

type uprojectFile struct {
	EngineAssociation string
	Category          string
	Description       string
	Modules           []moduleMetaInfo
	FileVersion       int
}

type moduleMetaInfo struct {
	Name         string
	Type         string
	LoadingPhase string
}

// GetEngineVersion 获取工程所用的引擎版本。
func GetEngineVersion(projectFilePath string) (string, error) {
	content, err := os.ReadFile(projectFilePath)
	if err != nil {
		return "", fmt.Errorf("open project file: %w", err)
	}

	var file uprojectFile
	err = json.Unmarshal(content, &file)
	if err != nil {
		return "", fmt.Errorf("unmarshal project file: %w", err)
	}

	if len(file.EngineAssociation) == 0 {
		return "", fmt.Errorf("empty EngineAssociation")
	}

	return file.EngineAssociation, nil
}

// EngineInfo 用于存放 UE 引擎的信息。
type EngineInfo struct {
	Version     string
	InstallPath string
}

func trimPsOutput(s string) string {
	return strings.Trim(strings.TrimSpace(s), "\"")
}

// FindAllEngineInfos 查找所有已安装的引擎信息。
func FindAllEngineInfos() ([]*EngineInfo, error) {
	sh := pwsh.New()

	stdOut, stdErr, err := sh.Execute(
		`(Get-ItemProperty "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\EpicGames\Unreal Engine\*") | ` +
			`%{ Write-Output $_.PSChildName $_.InstalledDirectory }`)

	if stdErr != "" {
		core.LogE("%s", stdErr)
	}

	if err != nil {
		return nil, err
	}

	sp := strings.Split(stdOut, "\n")
	if len(sp) <= 1 {
		return nil, fmt.Errorf("Unreal engine not found")
	}

	var infos []*EngineInfo
	for i := 1; i < len(sp); i += 2 {
		ver := trimPsOutput(sp[i-1])
		path := trimPsOutput(sp[i])
		infos = append(infos, &EngineInfo{
			Version:     ver,
			InstallPath: path,
		})
	}

	return infos, nil
}

// FindEngineInfo 获取特定版本的引擎的路径。如果不指定 version，则返回找到的第一个版本的信息。
func FindEngineInfo(version string) (*EngineInfo, error) {
	infos, err := FindAllEngineInfos()
	if err != nil {
		return nil, err
	}

	if len(version) == 0 && len(infos) != 0 {
		return infos[0], nil
	}

	for _, info := range infos {
		if info.Version == version {
			return info, nil
		}
	}

	return nil, fmt.Errorf("engine with version '%s' no found", version)
}

// ExecuteUbt 执行 Unreal Build Tool 的命令。
func ExecuteUbt(engineDir string, args string) error {
	sh := pwsh.New()

	binPath := filepath.Join(engineDir, "Engine", "Binaries", "DotNET", "UnrealBuildTool", "UnrealBuildTool.exe")
	pwCmd := fmt.Sprintf("& \"%s\" %s", binPath, args)
	stdOut, stdErr, err := sh.Execute(pwCmd)
	if stdOut != "" {
		core.LogD("%s", stdOut)
	}
	if stdErr != "" {
		core.LogE("%s", stdErr)
	}

	return err
}

// generateClangdFlagsFile 生成 clangd_args 文件，用于指定 clangd 的额外参数。
// ref: https://github.com/natsu-anon/ue-assist/
func generateClangdFlagsFile(projectDir string) (string, error) {
	bs, err := core.Global.EmbedFs.ReadFile("resources/compile/clangd_args.tmpl")
	if err != nil {
		return "", fmt.Errorf("load clangd_args file template: %w", err)
	}

	outFile := filepath.Join(projectDir, ".vscode", "clangd_args")
	core.LogD("write file to %s", outFile)
	return outFile, os.WriteFile(outFile, bs, 0644)
}

// ExecuteUbtGenProject 执行 Unreal Build Tool 的工程构建命令。
// 目前的实现参考了 ue-assist 项目，通过生成 vscode 的配置文件来产出 clangd 使用的 DB 文件。
// ref: https://github.com/natsu-anon/ue-assist/
func ExecuteUbtGenProject(engineDir string, projectFilePath string) error {
	projectFileName := filepath.Base(projectFilePath)
	projectName := strings.TrimSuffix(projectFileName, filepath.Ext(projectFileName))
	core.LogD("detect project name %s", projectName)

	args := fmt.Sprintf("-projectfiles -vscode -game -engine -dotnet -progress -noIntelliSense \"%s\"", projectFilePath)
	if err := ExecuteUbt(engineDir, args); err != nil {
		return fmt.Errorf("execute UBT: %w", err)
	}

	core.LogD("execute UBT %s success", projectFilePath)

	projectDir := filepath.Dir(projectFilePath)
	clangdFile, err := generateClangdFlagsFile(projectDir)
	if err != nil {
		return fmt.Errorf("generate clangd_args file: %w", err)
	}

	core.LogD("generate clangd_args file %s success", clangdFile)

	srcDbFileName := fmt.Sprintf("compileCommands_%s.json", projectName)
	srcDbFilePath := filepath.Join(projectDir, ".vscode", srcDbFileName)
	srcDbDataRaw, err := os.ReadFile(srcDbFilePath)

	// var dbDataArray []map[string]interface{}
	var dbDataArray []orderedmap.OrderedMap
	if err := json.Unmarshal(srcDbDataRaw, &dbDataArray); err != nil {
		return fmt.Errorf("unmarshal src clang database: %w", err)
	}

	clangdExtraArgs := fmt.Sprintf("@%s", clangdFile)
	for _, elem := range dbDataArray {
		args, ok := elem.Get("arguments")
		if !ok {
			continue
		}

		args = append(args.([]interface{}), clangdExtraArgs)
		elem.Set("arguments", args)
	}

	dstDbFilePath := filepath.Join(projectDir, "compile_commands.json")
	dstDbDataRaw, err := json.MarshalIndent(dbDataArray, "", "\t")
	if err != nil {
		return fmt.Errorf("marshal dst clang database: %w", err)
	}

	if err := os.WriteFile(dstDbFilePath, dstDbDataRaw, 0644); err != nil {
		return fmt.Errorf("write dst clang database to %s: %w", dstDbFilePath, err)
	}

	core.LogD("generate clang database from %s to %s success", srcDbFilePath, dstDbFilePath)
	return nil
}
