package newcmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/osutil"
)

// NewIgnoreCmd 是用于创建新 gitignore 文件的子命令。
type NewIgnoreCmd struct {
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

// Run 执行创建命令。
func (cmd *NewIgnoreCmd) Run() error {
	return osutil.DoInProjectRoot(cmd.ProjectFile, generateIgnoreFile)
}
