package pwsh

import (
	"bytes"
	"os/exec"

	"github.com/zhiruili/usak/encodeutil"
)

// PowerShell struct
type PowerShell struct {
	powerShell string
}

// New create new session
func New() *PowerShell {
	ps, _ := exec.LookPath("powershell.exe")
	return &PowerShell{
		powerShell: ps,
	}
}

func (p *PowerShell) Execute(args ...string) (stdOut string, stdErr string, err error) {
	args = append([]string{"-NoProfile", "-NonInteractive"}, args...)
	cmd := exec.Command(p.powerShell, args...)

	var outBuf bytes.Buffer
	var errBuf bytes.Buffer
	cmd.Stdout = &outBuf
	cmd.Stderr = &errBuf

	err = cmd.Run()
	stdOut = encodeutil.ByteToString(outBuf.Bytes(), encodeutil.GB18030)
	stdErr = encodeutil.ByteToString(errBuf.Bytes(), encodeutil.GB18030)
	return
}
