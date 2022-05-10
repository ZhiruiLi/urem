package infocmd

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/osutil"
	"github.com/zhiruili/urem/unreal"
)

// Cmd 是 info 子命令的集合。
type Cmd struct {
	EngineCommand    *EngineCmd     `arg:"subcommand:ue" help:"print associated engine info of the prject."`
	DefintionCommand *DefinitionCmd `arg:"subcommand:def" help:"find definition file of a UE class."`
}

// Run 实现了 subCmd 的接口。
func (cmd *Cmd) Run() error {
	if cmd.DefintionCommand != nil {
		return cmd.DefintionCommand.Run()
	} else if cmd.EngineCommand != nil {
		return cmd.EngineCommand.Run()
	}

	return fmt.Errorf("missing type: inc")
}

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

// DefinitionCmd 是用于查找 UE 类定义位置的子命令。
type DefinitionCmd struct {
	Version    string   `arg:"-v,--version" help:"UE version" default:"4.26"`
	ClassNames []string `arg:"positional,required"`
}

func getIncSourceSearchPaths(engineDir string) []string {
	srcDir := filepath.Join(engineDir, "Engine", "Source")
	return []string{
		filepath.Join(srcDir, "Runtime"),
		filepath.Join(srcDir, "Editor"),
	}
}

func fmtIncGrepResult(r *grepResult) string {
	prefix := r.Matched[1] + ": "

	base := filepath.Base(r.FileName)
	dir := filepath.Base(filepath.Dir(r.FileName))
	if dir == "Public" || dir == "Private" || dir == "" {
		return prefix + base
	}

	return prefix + dir + "/" + base
}

// Run 执行 UE 类定义位置查找逻辑。
func (cmd *DefinitionCmd) Run() error {
	engineDir, err := unreal.FindEngineDir(cmd.Version)
	if err != nil {
		return fmt.Errorf("find Unreal engine path: %w", err)
	}

	if len(cmd.ClassNames) == 0 {
		return fmt.Errorf("empty target class names")
	}

	var patterns []*grepPattern
	for _, name := range cmd.ClassNames {
		expr := `_API\s+(` + name + `)\s*[:{]`
		reg, err := regexp.Compile(expr)
		if err != nil {
			return fmt.Errorf("illegal regex expr %s: %w", expr, err)
		}

		patterns = append(patterns, &grepPattern{name, expr, reg})
	}

	searchDirs := getIncSourceSearchPaths(engineDir)
	results := grepManyDir(patterns, searchDirs)
	for _, result := range results {
		if result.Error != nil {
			core.LogE("fail to search in %s: %s", result.FileName, result.Error.Error())
		} else {
			core.LogD("%s match %s in %s:%d", result.Pattern, result.Matched[0], result.FileName, result.LineNo)
			fmt.Println(fmtIncGrepResult(result))
		}
	}

	return nil
}
