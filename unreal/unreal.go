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

func FindEngineDir(version string) (string, error) {
	sh := pwsh.New()

	stdOut, stdErr, err := sh.Execute(
		`(Get-ItemProperty "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\EpicGames\Unreal Engine\` +
			version +
			`" -Name "InstalledDirectory")."InstalledDirectory"`)
	engineDir := strings.Trim(strings.TrimSpace(stdOut), "\"")
	if stdErr != "" {
		core.LogE("%s", stdErr)
	}
	if err != nil {
		return "", err
	}
	return engineDir, nil
}

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

func ExecuteUbtGenProject(engineDir string, mode string, projectFilePath string) error {
	projectFileName := filepath.Base(projectFilePath)
	projectName := strings.TrimSuffix(projectFileName, filepath.Ext(projectFileName))
	core.LogD("detect project name %s", projectName)
	args := fmt.Sprintf("-mode=%s -project=\"%s\" -engine \"%s\" Development Win64", mode, projectFilePath, projectName)
	return ExecuteUbt(engineDir, args)
}