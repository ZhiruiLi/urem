package dummy

import "github.com/zhiruili/usak/core"

type Cmd struct {
}

func (cmd *Cmd) Run() error {
	core.LogD("hello world")
	return nil
}
