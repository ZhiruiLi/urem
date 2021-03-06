package gencmd

import (
	"fmt"
	"strings"

	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/osutil"
	"github.com/zhiruili/urem/pwsh"
)

// GenVsCmd 是 gen 子命令中负责生成 VS 工程的子命令。
type GenVsCmd struct {
	ProjectFile string `arg:"positional,required"`
}

func refreshSln(projectFilePath string) error {
	sh := pwsh.New()
	stdOut, stdErr, err := sh.Execute(
		`(Get-ItemProperty 'Registry::HKEY_CLASSES_ROOT\Unreal.ProjectFile\DefaultIcon').'(default)'`)
	binPath := strings.Trim(strings.TrimSpace(stdOut), "\"")
	if stdErr != "" {
		core.LogE("%s", stdErr)
	}
	if err != nil {
		return fmt.Errorf("find Unreal gen project bin: %w", err)
	}
	core.LogD("Unreal gen project bin path %s", binPath)

	pwCmd := fmt.Sprintf("& \"%s\" /projectfiles \"%s\"", binPath, projectFilePath)
	core.LogD("command: %s", pwCmd)

	stdOut, stdErr, err = sh.Execute(pwCmd)
	if stdOut != "" {
		core.LogD("%s", stdOut)
	}
	if stdErr != "" {
		core.LogE("%s", stdErr)
	}

	return err
}

// Run 执行生成操作。
func (cmd *GenVsCmd) Run() error {
	return osutil.DoInProjectRoot(cmd.ProjectFile, refreshSln)
}
