package newignore

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/osutil"
)

type Cmd struct {
	ProjectFile string `arg:"positional,required"`
}

func generateIgnoreFile(projectFilePath string) error {
	bs, err := core.Global.EmbedFs.ReadFile("resources/newignore/gitignore.tmpl")
	if err != nil {
		return fmt.Errorf("load ignore file template: %w", err)
	}

	dir := filepath.Dir(projectFilePath)
	outFile := filepath.Join(dir, ".gitignore")
	core.LogD("write file to %s", outFile)
	return ioutil.WriteFile(outFile, bs, 0644)
}

func (cmd *Cmd) Run() error {
	return osutil.DoInProjectRoot(cmd.ProjectFile, generateIgnoreFile)
}
