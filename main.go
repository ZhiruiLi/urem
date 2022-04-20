package main

import (
	"embed"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/newformat"
	"github.com/zhiruili/urem/newignore"
	"github.com/zhiruili/urem/newmod"
	"github.com/zhiruili/urem/regensln"
)

type SubCmd interface {
	Run() error
}

// verify interfaces
var (
	_ SubCmd = (*newmod.Cmd)(nil)
	_ SubCmd = (*newignore.Cmd)(nil)
	_ SubCmd = (*newformat.Cmd)(nil)
	_ SubCmd = (*regensln.Cmd)(nil)
)

var args struct {
	core.Args
	NewMod    *newmod.Cmd    `arg:"subcommand:newmod"`
	NewIgnore *newignore.Cmd `arg:"subcommand:newignore"`
	NewFormat *newformat.Cmd `arg:"subcommand:newformat"`
	RegenSln  *regensln.Cmd  `arg:"subcommand:regen"`
}

//go:embed resources/*/*.tmpl
var embedFs embed.FS

func main() {
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
