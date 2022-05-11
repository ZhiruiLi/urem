package infocmd

import (
	"fmt"
	"path/filepath"
	"regexp"

	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/unreal"
)

// DefinitionCmd 是用于查找 UE 类定义位置的子命令。
type DefinitionCmd struct {
	Version    string   `arg:"-v,--version" help:"UE version"`
	ClassNames []string `arg:"positional,required" help:"target class names"`
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