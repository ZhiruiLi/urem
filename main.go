package main

import (
	"embed"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/gencmd"
	"github.com/zhiruili/urem/lscmd"
	"github.com/zhiruili/urem/newcmd"
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
	_ SubCmd = (*newcmd.Cmd)(nil)
	_ SubCmd = (*gencmd.Cmd)(nil)
	_ SubCmd = (*lscmd.Cmd)(nil)
	_ SubCmd = (*dummyCmd)(nil)
)

var args struct {
	core.Args
	NewCommand *newcmd.Cmd `arg:"subcommand:new"`
	GenCommand *gencmd.Cmd `arg:"subcommand:gen"`
	LsCommand  *lscmd.Cmd  `arg:"subcommand:ls"`
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
