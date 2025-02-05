package gencmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/iancoleman/orderedmap"
	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/osutil"
	"github.com/zhiruili/urem/unreal"
)

// GenClangCmd 是 gen 子命令中负责生成 clang database 的子命令。
type GenClangCmd struct {
	Fast        bool   `arg:"-f,--fast"`
	ProjectFile string `arg:"positional,required"`
}

func (cmd *GenClangCmd) refreshClang(projectInfo *unreal.ProjectInfo) error {
	if !cmd.Fast {
		ver, err := projectInfo.GetEngineVersion()
		if err != nil {
			return fmt.Errorf("get Unreal engine version: %w", err)
		}

		core.LogD("get Unreal engine version: %s", ver)

		info, err := unreal.FindEngineInfo(ver)
		if err != nil {
			return fmt.Errorf("find Unreal engine info: %w", err)
		}

		core.LogD("find Unreal engine info: %s %s", info.Version, info.InstallPath)

		err = unreal.ExecuteUbtGenProject(info.InstallPath, projectInfo)
		if err != nil {
			return fmt.Errorf("generate clang database: %w", err)
		}

		core.LogD("generate clang database success")
	}

	clangdFile, err := unreal.GenerateClangdFlagsFile(projectInfo.ProjectDir())
	if err != nil {
		return fmt.Errorf("generate clangd_args file: %w", err)
	}

	core.LogD("generate clangd_args file %s success", clangdFile)

	srcDbFilePath := projectInfo.ProjectClangDbPath()
	srcDbDataRaw, err := os.ReadFile(srcDbFilePath)

	// 使用 orderedmap.OrderedMap 以保持字段原有的顺序
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

	dstDbFilePath := filepath.Join(projectInfo.ProjectDir(), "compile_commands.json")
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

// Run 执行生成操作。
func (cmd *GenClangCmd) Run() error {
	return osutil.DoInProjectRoot(cmd.ProjectFile, func(projPath string) error {
		return cmd.refreshClang(&unreal.ProjectInfo{ProjectFilePath: projPath})
	})
}
