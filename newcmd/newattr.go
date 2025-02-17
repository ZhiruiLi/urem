package newcmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/osutil"
)

// NewAttributeCmd 是用于创建新 gitattribute 文件的子命令。
type NewAttributeCmd struct {
	ProjectFile string `arg:"positional,required"`
}

func generateAttributeFile(projectFilePath string) error {
	bs, err := core.Global.EmbedFs.ReadFile("resources/newattribute/gitattribute.tmpl")
	if err != nil {
		return fmt.Errorf("load attribute file template: %w", err)
	}

	dir := filepath.Dir(projectFilePath)
	outFile := filepath.Join(dir, ".gitattributes")
	core.LogD("write file to %s", outFile)
	return os.WriteFile(outFile, bs, 0644)
}

// Run 执行创建命令。
func (cmd *NewAttributeCmd) Run() error {
	return osutil.DoInProjectRoot(cmd.ProjectFile, generateAttributeFile)
}
