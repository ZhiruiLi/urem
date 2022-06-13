package infocmd

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/grep"
	"github.com/zhiruili/urem/unreal"
)

// InfoDefine 是用于查找 UE 类定义位置的子命令。
type InfoDefine struct {
	Version      string   `arg:"-v,--version" help:"UE version"`
	WithDetails  bool     `arg:"-d,--detail" help:"print definition with details"`
	ClassOnly    bool     `arg:"-c,--class" help:"search class or struct only"`
	FunctionOnly bool     `arg:"-f,--function" help:"search function only"`
	TargetNames  []string `arg:"positional,required" help:"target class/function names"`
}

func getIncSourceSearchPaths(engineDir string) []string {
	srcDir := filepath.Join(engineDir, "Engine", "Source")
	return []string{
		filepath.Join(srcDir, "Runtime"),
		filepath.Join(srcDir, "Editor"),
	}
}

func fmtIncGrepResult(item *grep.Item) string {
	prefix := item.Matched[1] + " "
	base := filepath.Base(item.FileName)
	dir := filepath.Base(filepath.Dir(item.FileName))
	if dir == "Public" || dir == "Private" || dir == "" {
		return prefix + base
	}

	return prefix + dir + "/" + base
}

var delimeter = color.BlueString(strings.Repeat("─", 80))

func printDetailInfo(item *grep.Item) {
	fmt.Println(delimeter)
	fmt.Printf("%s:%s\n", color.GreenString(item.FileName), color.GreenString(strconv.Itoa(item.LineNo)))
	for _, head := range item.HeadLines {
		fmt.Println(head)
	}

	className := item.Matched[1]
	fmtClassName := color.New(color.Underline).Add(color.FgGreen).Sprintf("%s", className)
	fmt.Println(strings.ReplaceAll(item.LineText, className, fmtClassName))
	fmt.Println(delimeter)
}

// Run 执行 UE 类定义位置查找逻辑。
func (cmd *InfoDefine) Run() error {
	info, err := unreal.FindEngineInfo(cmd.Version)
	if err != nil {
		return fmt.Errorf("find Unreal engine info: %w", err)
	}

	if len(cmd.TargetNames) == 0 {
		return fmt.Errorf("empty target class names")
	}

	var patterns []*grep.Pattern
	for _, namePattern := range cmd.TargetNames {
		fixPattern := strings.ReplaceAll(namePattern, `.`, `[^\s]`)
		classPatten := `class\s+.*_API\s+(` + fixPattern + `)[\s:{]`
		structPatten := `struct\s+.*_API\s+(` + fixPattern + `)[\s:{]`
		functionPattern := `.*_API\s+.*\s+(` + fixPattern + `)\s*\(.*\)`

		exprs := []string{}

		if cmd.ClassOnly {
			exprs = append(exprs, classPatten, structPatten)
		}

		if cmd.FunctionOnly {
			exprs = append(exprs, functionPattern)
		}

		if !cmd.ClassOnly && !cmd.FunctionOnly {
			exprs = append(exprs, classPatten, structPatten, functionPattern)
		}

		for _, expr := range exprs {
			reg, err := regexp.Compile(expr)
			if err != nil {
				return fmt.Errorf("illegal regex expr %s: %w", expr, err)
			}
			patterns = append(patterns, &grep.Pattern{Name: namePattern, Raw: expr, Regexp: reg})
		}
	}

	searchDirs := getIncSourceSearchPaths(info.InstallPath)
	items := grep.GrepResult(patterns, searchDirs, grep.WithExts(".h", ".hpp"))
	for i, result := range items {
		if result.Error != nil {
			core.LogE("fail to search in %s: %s", result.FileName, result.Error.Error())
		} else {
			if cmd.WithDetails && i > 0 {
				fmt.Println("")
			}

			core.LogD("%s match %s in %s:%d", result.Pattern, result.Matched[0], result.FileName, result.LineNo)
			fmt.Println(fmtIncGrepResult(result))
			if cmd.WithDetails {
				printDetailInfo(result)
			}
		}
	}

	return nil
}
