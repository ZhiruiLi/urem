package newcmd

import (
	"fmt"
)

// Cmd 是 new 子命令的集合。
type Cmd struct {
	NewModCommand       *NewModCmd       `arg:"subcommand:mod"`
	NewFormatCommand    *NewFormatCmd    `arg:"subcommand:fmt"`
	NewIgnoreCommand    *NewIgnoreCmd    `arg:"subcommand:ig"`
	NewAttributeCommand *NewAttributeCmd `arg:"subcommand:attr"`
}

// Run 实现了 subCmd 的接口。
func (cmd *Cmd) Run() error {
	if cmd.NewModCommand != nil {
		return cmd.NewModCommand.Run()
	} else if cmd.NewFormatCommand != nil {
		return cmd.NewFormatCommand.Run()
	} else if cmd.NewIgnoreCommand != nil {
		return cmd.NewIgnoreCommand.Run()
	} else if cmd.NewAttributeCommand != nil {
		return cmd.NewAttributeCommand.Run()
	}

	return fmt.Errorf("missing target: mod/fmt/ig")
}
