package gencmd

import "fmt"

type Cmd struct {
	GenVsCommand    *GenVsCmd    `arg:"subcommand:vs"`
	GenClangCommand *GenClangCmd `arg:"subcommand:clang"`
}

func (cmd *Cmd) Run() error {
	if cmd.GenVsCommand != nil {
		return cmd.GenVsCommand.Run()
	} else if cmd.GenClangCommand != nil {
		return cmd.GenClangCommand.Run()
	}

	return fmt.Errorf("missing target: vs/clang")
}
