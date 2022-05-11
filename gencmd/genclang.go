package gencmd

import (
	"fmt"
	"path/filepath"

	"github.com/zhiruili/urem/osutil"
	"github.com/zhiruili/urem/unreal"
)

// GenClangCmd 是 gen 子命令中负责生成 clang database 的子命令。
type GenClangCmd struct {
	ProjectFile string `arg:"positional,required"`
}

func refreshClang(projectFilePath string) error {
	ver, err := unreal.GetEngineVersion(projectFilePath)
	if err != nil {
		return fmt.Errorf("get Unreal engine version: %w", err)
	}

	info, err := unreal.FindEngineInfo(ver)
	if err != nil {
		return fmt.Errorf("find Unreal engine info: %w", err)
	}

	err = unreal.ExecuteUbtGenProject(info.InstallPath, "GenerateClangDatabase", projectFilePath)
	if err != nil {
		return fmt.Errorf("generate clang database: %w", err)
	}

	projectDir := filepath.Dir(projectFilePath)
	srcDbFile := filepath.Join(info.InstallPath, "compile_commands.json")
	dstDbFile := filepath.Join(projectDir, "compile_commands.json")
	if err := osutil.CopyFile(srcDbFile, dstDbFile); err != nil {
		return fmt.Errorf("copy clang database from %s to %s: %w", srcDbFile, dstDbFile, err)
	}

	return nil
}

// Run 执行生成操作。
func (cmd *GenClangCmd) Run() error {
	return osutil.DoInProjectRoot(cmd.ProjectFile, func(projPath string) error {
		return refreshClang(projPath)
	})
}
