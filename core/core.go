package core

import (
	"embed"
	"fmt"
	"os"
)

type Args struct {
	Debug   bool `arg:"--debug" help:"run with debug mode"`
	Quite   bool `arg:"-q,--quite" help:"disable all logs"`
	Verbose bool `arg:"-v,--verbose" help:"verbosity level"`
}

type GlobalData struct {
	Args
	EmbedFs embed.FS
}

var Global GlobalData

func LogD(f string, a ...interface{}) {
	if !Global.Verbose || Global.Quite {
		return
	}

	fmt.Fprintf(os.Stdout, f+"\n", a...)
}

func LogI(f string, a ...interface{}) {
	if Global.Quite {
		return
	}

	fmt.Fprintf(os.Stdout, f+"\n", a...)
}

func LogE(f string, a ...interface{}) {
	if Global.Quite {
		return
	}

	fmt.Fprintf(os.Stderr, f+"\n", a...)
}
