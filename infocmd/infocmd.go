package infocmd

import (
	"fmt"
)

// Cmd 是 info 子命令的集合。
type Cmd struct {
	EngineCommand *InfoEngineCmd `arg:"subcommand:ue" help:"print associated engine info of the prject."`
	EnumCommand   *InfoEnumCmd   `arg:"subcommand:enum" help:"print available enum value."`
}

// Run 实现了 subCmd 的接口。
func (cmd *Cmd) Run() error {
	if cmd.EngineCommand != nil {
		return cmd.EngineCommand.Run()
	} else if cmd.EnumCommand != nil {
		return cmd.EnumCommand.Run()
	}

	return fmt.Errorf("missing subcommand of info cmd")
}
