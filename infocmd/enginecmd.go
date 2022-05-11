package infocmd

import (
	"fmt"

	"github.com/zhiruili/urem/osutil"
	"github.com/zhiruili/urem/unreal"
)

// InfoEngineCmd 是用于查找工程使用的 UE 相关信息的子命令。
type InfoEngineCmd struct {
	ProjectFile string `arg:"positional" help:"print engine info used by the given project, if not set, print all"`
}

func printEngineInfo(projectFilePath string) error {
	ver, err := unreal.GetEngineVersion(projectFilePath)
	if err != nil {
		return fmt.Errorf("get Unreal engine version: %s", err.Error())
	}

	fmt.Printf("Unreal Version: %s\n", ver)

	info, err := unreal.FindEngineInfo(ver)
	if err != nil {
		return fmt.Errorf("find Unreal engine info: %w", err)
	}

	fmt.Printf("Install Path: %s\n", info.InstallPath)
	return nil
}

func printAllEngineInfos() error {
	infos, err := unreal.FindAllEngineInfos()
	if err != nil {
		return fmt.Errorf("find Unreal engine info: %w", err)
	}

	for i, info := range infos {
		if i != 0 {
			fmt.Println("")
		}

		fmt.Printf("Unreal Version: %s\n", info.Version)
		fmt.Printf("Install Path: %s\n", info.InstallPath)
	}

	return nil
}

// Run 执行 UE 相关信息查找逻辑。
func (cmd *InfoEngineCmd) Run() error {
	if len(cmd.ProjectFile) == 0 {
		return printAllEngineInfos()
	}

	return osutil.DoInProjectRoot(cmd.ProjectFile, printEngineInfo)
}
