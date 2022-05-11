package core

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"strings"
)

// Args 是应用的基础命令行选项。
type Args struct {
	Debug   bool `arg:"--debug" help:"run with debug mode"`
	Quite   bool `arg:"-q,--quite" help:"don't ask for user input"`
	Verbose bool `arg:"-v,--verbose" help:"increase log verbosity level"`
}

type globalData struct {
	Args
	EmbedFs embed.FS
}

// Global 存放全局需要的一些数据。
var Global globalData

// StrContains 检查字符串是否在给定数组中。
func StrContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

// StrSliceMap 映射一个字符串数组。
func StrSliceMap(ss []string, mapper func(string) string) []string {
	if len(ss) == 0 {
		return nil
	}

	newss := make([]string, len(ss))
	if mapper == nil {
		copy(newss, ss)
	} else {
		for i, s := range ss {
			newss[i] = mapper(s)
		}
	}
	return newss
}

// LogD 打印 debug 级别的日志。
func LogD(f string, a ...interface{}) {
	if !Global.Verbose {
		return
	}

	fmt.Fprintf(os.Stdout, f+"\n", a...)
}

// LogI 打印 info 级别的日志。
func LogI(f string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, f+"\n", a...)
}

// LogE 打印 error 级别的日志。
func LogE(f string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, f+"\n", a...)
}

// IllegalArg 是用于参数非法时返回的错误类型。
type IllegalArg struct {
	ArgName string
	Message string
}

// Error 实现了 error interface。
func (err *IllegalArg) Error() string {
	return fmt.Sprintf("illegal %s: %s", err.ArgName, err.Message)
}

// IllegalArgErrorf 创建一个 IllegalArg 类型的 error 对象。
func IllegalArgErrorf(argName string, messageF string, a ...interface{}) error {
	return &IllegalArg{
		ArgName: argName,
		Message: fmt.Sprintf(messageF, a...),
	}
}

// GetUserInput 获取用户输入并返回。
func GetUserInput(hint string, availableInputs ...string) string {
	if Global.Quite {
		return availableInputs[0]
	}

	fullHint := fmt.Sprintf("%s (%s) ", hint, strings.Join(availableInputs, "/"))
	fmt.Print(fullHint)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		txt := input.Text()
		if StrContains(availableInputs, txt) {
			return txt
		}
		fmt.Print("illegal input, " + fullHint)
	}

	return ""
}

// GetUserBoolInput 获取用户的 yes or no 输入并返回 bool。
func GetUserBoolInput(hint string) bool {
	input := GetUserInput(hint, "y", "n")
	return input == "y"
}
