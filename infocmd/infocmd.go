package infocmd

import (
	"fmt"
)

// Cmd 是 info 子命令的集合。
type Cmd struct {
	EngineCommand    *InfoEngineCmd `arg:"subcommand:ue" help:"print associated engine info of the prject."`
	EnumCommand      *InfoEnumCmd   `arg:"subcommand:enum" help:"print available enum value."`
	DefintionCommand *InfoDefine    `arg:"subcommand:def" help:"find definition file of a UE class."`
}

// Run 实现了 subCmd 的接口。
func (cmd *Cmd) Run() error {
	if cmd.EngineCommand != nil {
		return cmd.EngineCommand.Run()
	} else if cmd.EnumCommand != nil {
		return cmd.EnumCommand.Run()
	} else if cmd.DefintionCommand != nil {
		return cmd.DefintionCommand.Run()
	}

	return fmt.Errorf("missing subcommand of info cmd")
}
