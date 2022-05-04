package pwsh

import (
	"bytes"
	"os/exec"

	"github.com/zhiruili/urem/encodeutil"
)

// PowerShell 是 PowerShell 对应的结构体。
type PowerShell struct {
	powerShell string
}

// New 创建一个新的 session。
func New() *PowerShell {
	ps, _ := exec.LookPath("powershell.exe")
	return &PowerShell{
		powerShell: ps,
	}
}

// Execute 执行一段 PowerShell 命令。
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
