package unreal

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/pwsh"
)

type ProjectInfo struct {
	ProjectFilePath string
}

func (pi *ProjectInfo) ProjectFileName() string {
	return filepath.Base(pi.ProjectFilePath)
}

func (pi *ProjectInfo) ProjectName() string {
	projectFileName := pi.ProjectFileName()
	return strings.TrimSuffix(projectFileName, filepath.Ext(projectFileName))
}

func (pi *ProjectInfo) ProjectDir() string {
	return filepath.Dir(pi.ProjectFilePath)
}

func (pi *ProjectInfo) ProjectVscodeDir() string {
	return filepath.Join(pi.ProjectDir(), ".vscode")
}

func (pi *ProjectInfo) ProjectClangDbName() string {
	return fmt.Sprintf("compileCommands_%s.json", pi.ProjectName())
}

func (pi *ProjectInfo) ProjectClangDbPath() string {
	return filepath.Join(pi.ProjectVscodeDir(), pi.ProjectClangDbName())
}

// GetEngineVersion 获取工程所用的引擎版本。
func (pi *ProjectInfo) GetEngineVersion() (string, error) {
	content, err := os.ReadFile(pi.ProjectFilePath)
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

// EngineInfo 用于存放 UE 引擎的信息。
type EngineInfo struct {
	Version     string
	InstallPath string
}

func trimPsOutput(s string) string {
	return strings.Trim(strings.TrimSpace(s), "\"")
}

// FindAllEngineInfos 查找所有已安装的引擎信息。
// 需要根据不同平台参考 UE 自己的实现，这里先只处理 Windows 的情况，参考 UE 的实现：
// FDesktopPlatformWindows::EnumerateEngineInstallations
func FindAllEngineInfos() ([]*EngineInfo, error) {
	sh := pwsh.New()
	var stdOut, stdErr string
	var err error
	var cmdResultList []string

	stdOut, stdErr, err = sh.Execute(
		`(Get-ItemProperty "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\EpicGames\Unreal Engine\*") | ` +
			`%{ Write-Output $_.PSChildName $_.InstalledDirectory }`)

	if stdErr != "" {
		core.LogE("%s", stdErr)
	}

	if err != nil {
		return nil, err
	}

	stdOut = strings.TrimSpace(stdOut)
	core.LogD("find engine info 1 stdOut:\n%s", stdOut)
	cmdResultList = append(cmdResultList, strings.Split(stdOut, "\n")...)

	stdOut, stdErr, err = sh.Execute(
		`(Get-ItemProperty "Registry::HKEY_CURRENT_USER\SOFTWARE\Epic Games\Unreal Engine\Builds").PSObject.Properties | ` +
			`Where-Object { $_.Name -match '^\{.*\}$' } | %{ Write-Output $_.Name $_.Value }`)

	if stdErr != "" {
		core.LogE("%s", stdErr)
	}

	if err != nil {
		return nil, err
	}

	stdOut = strings.TrimSpace(stdOut)
	core.LogD("find engine info 2 stdOut:\n%s", stdOut)
	cmdResultList = append(cmdResultList, strings.Split(stdOut, "\n")...)

	if len(cmdResultList) <= 1 {
		return nil, fmt.Errorf("unreal engine not found")
	}

	var infos []*EngineInfo
	for i := 1; i < len(cmdResultList); i += 2 {
		ver := trimPsOutput(cmdResultList[i-1])
		path := trimPsOutput(cmdResultList[i])
		if len(ver) == 0 || len(path) == 0 {
			continue
		}

		core.LogD("find engine info %d: ver %s path %s", i/2, ver, path)

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

	if core.Global.Verbose {
		var buf bytes.Buffer
		for i, info := range infos {
			if i > 0 {
				buf.WriteString("\n")
			}

			buf.WriteString(info.Version)
			buf.WriteString(": ")
			buf.WriteString(info.InstallPath)
		}

		core.LogD("installed engines:\n%s", buf.String())
	}

	return nil, fmt.Errorf("engine with version '%s' no found", version)
}

// ExecuteUbt 执行 Unreal Build Tool 的命令。
func ExecuteUbt(engineDir string, args string) error {
	sh := pwsh.New()

	// 参考 UE 的实现：
	// FDesktopPlatformWindows::RunUnrealBuildTool
	// 这个函数中使用的路径是：Engine/Build/BatchFiles/Build.bat
	binPath := filepath.Join(engineDir, "Engine", "Build", "BatchFiles", "Build.bat")
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

// GenerateClangdFlagsFile 生成 clangd_args 文件，用于指定 clangd 的额外参数。
// ref: https://github.com/natsu-anon/ue-assist/
func GenerateClangdFlagsFile(projectDir string) (string, error) {
	bs, err := core.Global.EmbedFs.ReadFile("resources/compile/clangd_args.tmpl")
	if err != nil {
		return "", fmt.Errorf("load clangd_args file template: %w", err)
	}

	outFile := filepath.Join(projectDir, ".vscode", "clangd_args")
	core.LogD("write file to %s", outFile)
	return outFile, os.WriteFile(outFile, bs, 0644)
}

// ExecuteUbtGenProject 执行 Unreal Build Tool 的工程构建命令。
// 目前的实现参考了 ue-assist 项目，通过生成 vscode 的配置文件来产出 clangd 使用的 DB 文件，参考：
// https://github.com/natsu-anon/ue-assist/
// 不过内容需要根据 UE 的实现进行调整，参考：
// FDesktopPlatformBase::GenerateProjectFiles
func ExecuteUbtGenProject(engineDir string, projectInfo *ProjectInfo) error {
	projectName := projectInfo.ProjectName()
	core.LogD("detect project name %s", projectName)

	args := fmt.Sprintf("-projectfiles -project='%s' -game -engine -progress", projectInfo.ProjectFilePath)
	if err := ExecuteUbt(engineDir, args); err != nil {
		return fmt.Errorf("execute UBT: %w", err)
	}

	core.LogD("execute UBT %s success", projectInfo.ProjectFilePath)
	return nil
}
