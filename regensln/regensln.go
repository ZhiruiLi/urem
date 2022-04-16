package regensln

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/pwsh"
)

type Cmd struct {
	ProjectFile string `arg:"positional,required"`
}

func refreshSln(projectFilePath string) error {
	sh := pwsh.New()
	stdOut, stdErr, err := sh.Execute(`(Get-ItemProperty 'Registry::HKEY_CLASSES_ROOT\Unreal.ProjectFile\DefaultIcon').'(default)'`)
	binPath := strings.Trim(strings.TrimSpace(stdOut), "\"")
	if stdErr != "" {
		core.LogE("%s", stdErr)
	}
	if err != nil {
		return fmt.Errorf("find Unreal gen project bin: %w", err)
	}
	core.LogD("Unreal gen project bin path %s", binPath)

	pwCmd := fmt.Sprintf("& \"%s\" /projectfiles \"%s\"", binPath, projectFilePath)
	core.LogD("Command: %s", pwCmd)

	stdOut, stdErr, err = sh.Execute(pwCmd)
	if stdOut != "" {
		core.LogD("%s", stdOut)
	}
	if stdErr != "" {
		core.LogE("%s", stdErr)
	}

	return err
}

func (cmd *Cmd) Run() error {
	ext := filepath.Ext(cmd.ProjectFile)
	if ext != ".uproject" {
		return fmt.Errorf("illegal project file, must be *.uproject file, got %s", cmd.ProjectFile)
	}

	absProjectFilePath, err := filepath.Abs(cmd.ProjectFile)
	if err != nil {
		return fmt.Errorf("illegal project file path %s", cmd.ProjectFile)
	}

	return refreshSln(absProjectFilePath)
}
