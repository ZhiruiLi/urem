package core

import (
	"bufio"
	"embed"
	"fmt"
	"os"
	"strings"
)

type Args struct {
	Debug   bool `arg:"--debug" help:"run with debug mode"`
	Quite   bool `arg:"-q,--quite" help:"don't ask for user input"`
	Verbose bool `arg:"-v,--verbose" help:"increase log verbosity level"`
}

type globalData struct {
	Args
	EmbedFs embed.FS
}

var Global globalData

func StrContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func LogD(f string, a ...interface{}) {
	if !Global.Verbose {
		return
	}

	fmt.Fprintf(os.Stdout, f+"\n", a...)
}

func LogI(f string, a ...interface{}) {
	fmt.Fprintf(os.Stdout, f+"\n", a...)
}

func LogE(f string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, f+"\n", a...)
}

type IllegalArg struct {
	ArgName string
	Message string
}

func (err *IllegalArg) Error() string {
	return fmt.Sprintf("illegal %s: %s", err.ArgName, err.Message)
}

func IllegalArgErrorf(argName string, messageF string, a ...interface{}) error {
	return &IllegalArg{
		ArgName: argName,
		Message: fmt.Sprintf(messageF, a...),
	}
}

func GetUserInput(hint string, availableInputs ...string) string {
	if Global.Quite {
		return availableInputs[0]
	}

	fullHint := fmt.Sprintf("%s (%s) ", hint, strings.Join(availableInputs, "/"))
	fmt.Print(fullHint)
	input := bufio.NewScanner(os.Stdin)
	for input.Scan() {
		txt := input.Text()
		if StrContains(availableInputs, txt) {
			return txt
		}
		fmt.Print("illegal input, " + fullHint)
	}

	return ""
}

func GetUserBoolInput(hint string) bool {
	input := GetUserInput(hint, "y", "n")
	return input == "y"
}
