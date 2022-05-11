package main

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"runtime/pprof"

	"github.com/alexflint/go-arg"
	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/gencmd"
	"github.com/zhiruili/urem/infocmd"
	"github.com/zhiruili/urem/newcmd"
)

type subCmd interface {
	Run() error
}

type dummyCmd struct{}

func (*dummyCmd) Run() error {
	return nil
}

// 确保 interfaces 有效。
var (
	_ subCmd = (*newcmd.Cmd)(nil)
	_ subCmd = (*gencmd.Cmd)(nil)
	_ subCmd = (*infocmd.Cmd)(nil)
	_ subCmd = (*dummyCmd)(nil)
)

type args struct {
	core.Args
	NewCommand  *newcmd.Cmd  `arg:"subcommand:new"`
	GenCommand  *gencmd.Cmd  `arg:"subcommand:gen"`
	InfoCommand *infocmd.Cmd `arg:"subcommand:info"`
}

func (args) Version() string {
	return "URem 0.1.1"
}

//go:embed resources/*/*.tmpl
var embedFs embed.FS

func main() {
	dummy := &dummyCmd{}
	if err := dummy.Run(); err != nil {
		panic(err)
	}

	var args args
	p := arg.MustParse(&args)

	if len(args.PProfFile) != 0 {
		prof, err := os.OpenFile(args.PProfFile, os.O_CREATE|os.O_RDWR, 0644)
		if err != nil {
			core.LogE("fail open profiling file: %s", err.Error())
			os.Exit(-2)
		}

		defer prof.Close()
		pprof.StartCPUProfile(bufio.NewWriter(prof))
		defer func() {
			pprof.StopCPUProfile()
			fmt.Printf("done write profiling file, run:\n"+
				"    go tool pprof -http=:9999 %s\n"+
				"to see the result", args.PProfFile)
		}()
	}

	core.Global.Args = args.Args
	core.Global.EmbedFs = embedFs
	if p.Subcommand() == nil {
		p.Fail("missing subcommand")
	} else if cmd, ok := p.Subcommand().(subCmd); !ok {
		p.Fail("illegal subcommand")
	} else if err := cmd.Run(); err != nil {
		core.LogE("error: %s", err.Error())
		os.Exit(-1)
	}
}
