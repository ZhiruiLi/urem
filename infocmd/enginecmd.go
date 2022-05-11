package infocmd

import (
	"fmt"

	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/osutil"
	"github.com/zhiruili/urem/unreal"
)

// EngineCmd 是用于查找工程使用的 UE 相关信息的子命令。
type EngineCmd struct {
	ProjectFile string `arg:"positional,required"`
}

func printEngineInfo(projectFilePath string) error {
	ver, err := unreal.GetEngineVersion(projectFilePath)
	if err != nil {
		core.LogE("get Unreal engine version: %s", err.Error())
		return nil
	}

	fmt.Printf("Unreal Version: %s\n", ver)

	dir, err := unreal.FindEngineDir(ver)
	if err != nil {
		return fmt.Errorf("find Unreal engine path: %w", err)
	}

	fmt.Printf("Install Path: %s\n", dir)
	return nil
}

// Run 执行 UE 相关信息查找逻辑。
func (cmd *EngineCmd) Run() error {
	return osutil.DoInProjectRoot(cmd.ProjectFile, printEngineInfo)
}
