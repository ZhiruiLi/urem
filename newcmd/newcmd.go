package newcmd

import (
	"fmt"
)

type Cmd struct {
	NewModCommand    *NewModCmd    `arg:"subcommand:mod"`
	NewFormatCommand *NewFormatCmd `arg:"subcommand:fmt"`
	NewIgnoreCommand *NewIgnoreCmd `arg:"subcommand:ig"`
}

func (cmd *Cmd) Run() error {
	if cmd.NewModCommand != nil {
		return cmd.NewModCommand.Run()
	} else if cmd.NewFormatCommand != nil {
		return cmd.NewFormatCommand.Run()
	} else if cmd.NewIgnoreCommand != nil {
		return cmd.NewIgnoreCommand.Run()
	}

	return fmt.Errorf("missing target: mod/fmt/ig")
}
