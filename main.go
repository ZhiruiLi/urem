package main

import (
	"embed"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/genclang"
	"github.com/zhiruili/urem/genvs"
	"github.com/zhiruili/urem/newformat"
	"github.com/zhiruili/urem/newignore"
	"github.com/zhiruili/urem/newmod"
)

type SubCmd interface {
	Run() error
}

type dummyCmd struct{}

func (*dummyCmd) Run() error {
	return nil
}

// verify interfaces
var (
	_ SubCmd = (*newmod.Cmd)(nil)
	_ SubCmd = (*newignore.Cmd)(nil)
	_ SubCmd = (*newformat.Cmd)(nil)
	_ SubCmd = (*genvs.Cmd)(nil)
	_ SubCmd = (*genclang.Cmd)(nil)
	_ SubCmd = (*dummyCmd)(nil)
)

var args struct {
	core.Args
	NewMod    *newmod.Cmd    `arg:"subcommand:newmod"`
	NewIgnore *newignore.Cmd `arg:"subcommand:newignore"`
	NewFormat *newformat.Cmd `arg:"subcommand:newformat"`
	GenVS     *genvs.Cmd     `arg:"subcommand:genvs"`
	GenClang  *genclang.Cmd  `arg:"subcommand:genclang"`
}

//go:embed resources/*/*.tmpl
var embedFs embed.FS

func main() {
	dummy := &dummyCmd{}
	if err := dummy.Run(); err != nil {
		panic(err)
	}

	p := arg.MustParse(&args)

	core.Global.Args = args.Args
	core.Global.EmbedFs = embedFs
	if p.Subcommand() == nil {
		p.Fail("missing subcommand")
	} else if cmd, ok := p.Subcommand().(SubCmd); !ok {
		p.Fail("illegal subcommand")
	} else if err := cmd.Run(); err != nil {
		core.LogE("error: %s", err.Error())
		os.Exit(-1)
	}
}
