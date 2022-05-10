package infocmd

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/osutil"
	"github.com/zhiruili/urem/unreal"
)

type Cmd struct {
	IncludeCommand *IncludeCmd `arg:"subcommand:inc" help:"find include dir of a UE class."`
	EngineCommand  *EngineCmd  `arg:"subcommand:engine" help:"print associated engine info of the prject."`
}

// Run 实现了 subCmd 的接口。
func (cmd *Cmd) Run() error {
	if cmd.IncludeCommand != nil {
		return cmd.IncludeCommand.Run()
	} else if cmd.EngineCommand != nil {
		return cmd.EngineCommand.Run()
	}

	return fmt.Errorf("missing type: inc")
}

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

func (cmd *EngineCmd) Run() error {
	return osutil.DoInProjectRoot(cmd.ProjectFile, printEngineInfo)
}

type IncludeCmd struct {
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

func (cmd *IncludeCmd) Run() error {
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
