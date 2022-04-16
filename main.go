package main

import (
	"embed"
	"os"

	"github.com/alexflint/go-arg"
	"github.com/zhiruili/usak/core"
	"github.com/zhiruili/usak/dummy"
	"github.com/zhiruili/usak/newmod"
)

var args struct {
	core.Args
	NewMod *newmod.Cmd `arg:"subcommand:newmod"`
	Dummy  *dummy.Cmd  `arg:"subcommand:dummy"`
}

//go:embed resources/*/*.tmpl
var embedFs embed.FS

func main() {
	arg.MustParse(&args)
	core.Global.Args = args.Args
	core.Global.EmbedFs = embedFs
	var err error
	if args.NewMod != nil {
		err = args.NewMod.Run()
	} else if args.Dummy != nil {
		err = args.Dummy.Run()
	}
	if err != nil {
		core.LogE("error: %s", err.Error())
		os.Exit(-1)
	}
}
