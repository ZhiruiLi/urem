package gencmd

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/osutil"
	"github.com/zhiruili/urem/pwsh"
)

// GenClangCmd 是 gen 子命令中负责生成 clang database 的子命令。
type GenClangCmd struct {
	Version     string `arg:"-v,--version" help:"UE version" default:"4.26"`
	ProjectFile string `arg:"positional,required"`
}

func refreshClang(version string, projectFilePath string) error {
	sh := pwsh.New()
	stdOut, stdErr, err := sh.Execute(`(Get-ItemProperty "Registry::HKEY_LOCAL_MACHINE\SOFTWARE\EpicGames\Unreal Engine\` + version + `" -Name "InstalledDirectory")."InstalledDirectory"`)
	engineDir := strings.Trim(strings.TrimSpace(stdOut), "\"")
	if stdErr != "" {
		core.LogE("%s", stdErr)
	}
	if err != nil {
		return fmt.Errorf("find Unreal engine path: %w", err)
	}

	binPath := filepath.Join(engineDir, "Engine", "Binaries", "DotNET", "UnrealBuildTool.exe")
	core.LogD("UnrealBuildTool file path %s", binPath)

	projectFileName := filepath.Base(projectFilePath)
	projectName := strings.TrimSuffix(projectFileName, filepath.Ext(projectFileName))
	core.LogD("detect project name %s", projectName)

	pwCmd := fmt.Sprintf("& \"%s\" -mode=GenerateClangDatabase -project=\"%s\" -engine \"%s\" Development Win64", binPath, projectFilePath, projectName)
	core.LogD("command: %s", pwCmd)

	stdOut, stdErr, err = sh.Execute(pwCmd)
	if stdOut != "" {
		core.LogD("%s", stdOut)
	}
	if stdErr != "" {
		core.LogE("%s", stdErr)
	}
	if err != nil {
		return fmt.Errorf("generate clang database: %w", err)
	}

	projectDir := filepath.Dir(projectFilePath)
	srcDbFile := filepath.Join(engineDir, "compile_commands.json")
	dstDbFile := filepath.Join(projectDir, "compile_commands.json")
	if err := osutil.CopyFile(srcDbFile, dstDbFile); err != nil {
		return fmt.Errorf("copy clang database from %s to %s: %w", srcDbFile, dstDbFile, err)
	}

	return nil
}

// Run 执行生成操作。
func (cmd *GenClangCmd) Run() error {
	return osutil.DoInProjectRoot(cmd.ProjectFile, func(projPath string) error {
		return refreshClang(cmd.Version, projPath)
	})
}
