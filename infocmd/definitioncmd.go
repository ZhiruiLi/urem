package infocmd

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/unreal"
)

// InfoDefine 是用于查找 UE 类定义位置的子命令。
type InfoDefine struct {
	Version     string   `arg:"-v,--version" help:"UE version"`
	FullPath    bool     `arg:"-f,--full" help:"print file full path"`
	LineNo      bool     `arg:"-l,--line" help:"print definition line no in file"`
	WithDetails bool     `arg:"-d,--detail" help:"print definition with details"`
	ClassNames  []string `arg:"positional,required" help:"target class names"`
}

func getIncSourceSearchPaths(engineDir string) []string {
	srcDir := filepath.Join(engineDir, "Engine", "Source")
	return []string{
		filepath.Join(srcDir, "Runtime"),
		filepath.Join(srcDir, "Editor"),
	}
}

func fmtIncGrepResult(r *grepResult, full bool, line bool) string {
	prefix := r.Matched[1] + " "
	suffix := ""
	if line {
		suffix = fmt.Sprintf(":%d", r.LineNo)
	}

	if full {
		return prefix + r.FileName + suffix
	}

	base := filepath.Base(r.FileName)
	dir := filepath.Base(filepath.Dir(r.FileName))
	if dir == "Public" || dir == "Private" || dir == "" {
		return prefix + base + suffix
	}

	return prefix + dir + "/" + base + suffix
}

var delimeter = color.BlueString(strings.Repeat("─", 120))

func printDetailInfo(r *grepResult) {
	fmt.Println(delimeter)
	fmt.Printf("%s:%s\n", color.GreenString(r.FileName), color.GreenString(strconv.Itoa(r.LineNo)))
	for _, head := range r.HeadLines {
		fmt.Println(head)
	}

	className := r.Matched[1]
	fmtClassName := color.New(color.Underline).Add(color.FgGreen).Sprintf("%s", className)
	fmt.Println(strings.ReplaceAll(r.LineText, className, fmtClassName))
	fmt.Println(delimeter)
}

// Run 执行 UE 类定义位置查找逻辑。
func (cmd *InfoDefine) Run() error {
	info, err := unreal.FindEngineInfo(cmd.Version)
	if err != nil {
		return fmt.Errorf("find Unreal engine info: %w", err)
	}

	if len(cmd.ClassNames) == 0 {
		return fmt.Errorf("empty target class names")
	}

	var patterns []*grepPattern
	for _, namePattern := range cmd.ClassNames {
		fixPattern := strings.ReplaceAll(namePattern, `.`, `[^\s]`)
		expr := `_API\s+(` + fixPattern + `)[\s:{]`
		reg, err := regexp.Compile(expr)
		if err != nil {
			return fmt.Errorf("illegal regex expr %s: %w", expr, err)
		}

		patterns = append(patterns, &grepPattern{namePattern, expr, reg})
	}

	searchDirs := getIncSourceSearchPaths(info.InstallPath)
	results := grepManyDir(patterns, searchDirs)
	for i, result := range results {
		if result.Error != nil {
			core.LogE("fail to search in %s: %s", result.FileName, result.Error.Error())
		} else {
			if cmd.WithDetails && i > 0 {
				fmt.Println("")
			}

			core.LogD("%s match %s in %s:%d", result.Pattern, result.Matched[0], result.FileName, result.LineNo)
			fmt.Println(fmtIncGrepResult(result, cmd.FullPath, cmd.LineNo))
			if cmd.WithDetails {
				printDetailInfo(result)
			}
		}
	}

	return nil
}
