package unreal

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"

	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/pwsh"
)

type uprojectFile struct {
	FileVersion       int
	EngineAssociation string
	Category          string
	Description       string
	Modules           []moduleMetaInfo
}

type moduleMetaInfo struct {
	Name         string
	Type         string
	LoadingPhase string
}

// GetEngineVersion 获取工程所用的引擎版本。
func GetEngineVersion(projectFilePath string) (string, error) {
	content, err := ioutil.ReadFile(projectFilePath)
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

// FindEngineDir 获取特定版本的引擎的路径。
func FindEngineDir(version string) (string, error) {
	sh := pwsh.New()

	if len(version) == 0 {
		version = "*"
	}

	stdOut, stdErr, err := sh.Execute(
		`(Get-ItemProperty "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\EpicGames\Unreal Engine\` +
			version + `").InstalledDirectory`)

	sp := strings.Split(stdOut, "\n")
	if len(sp) == 0 {
		return "", fmt.Errorf("Unreal engine not found")
	}

	for _, dir := range sp {
		if len(dir) != 0 {
			core.LogD("found engine dir: %s", dir)
		}
	}

	engineDir := strings.Trim(strings.TrimSpace(sp[0]), "\"")
	if stdErr != "" {
		core.LogE("%s", stdErr)
	}

	if err != nil {
		return "", err
	}

	return engineDir, nil
}

// ExecuteUbt 执行 Unreal Build Tool 的命令。
func ExecuteUbt(engineDir string, args string) error {
	sh := pwsh.New()

	binPath := filepath.Join(engineDir, "Engine", "Binaries", "DotNET", "UnrealBuildTool.exe")
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

// ExecuteUbtGenProject 执行 Unreal Build Tool 的工程构建命令。
func ExecuteUbtGenProject(engineDir string, mode string, projectFilePath string) error {
	projectFileName := filepath.Base(projectFilePath)
	projectName := strings.TrimSuffix(projectFileName, filepath.Ext(projectFileName))
	core.LogD("detect project name %s", projectName)
	args := fmt.Sprintf("-mode=%s -project=\"%s\" -engine \"%s\" Development Win64", mode, projectFilePath, projectName)
	return ExecuteUbt(engineDir, args)
}
