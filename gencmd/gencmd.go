package gencmd

import "fmt"

// Cmd 是 gen 子命令的集合。
type Cmd struct {
	GenVsCommand    *GenVsCmd    `arg:"subcommand:vs"`
	GenClangCommand *GenClangCmd `arg:"subcommand:clang"`
}

// Run 实现了 subCmd 的接口。
func (cmd *Cmd) Run() error {
	if cmd.GenVsCommand != nil {
		return cmd.GenVsCommand.Run()
	} else if cmd.GenClangCommand != nil {
		return cmd.GenClangCommand.Run()
	}

	return fmt.Errorf("missing target: vs/clang")
}
