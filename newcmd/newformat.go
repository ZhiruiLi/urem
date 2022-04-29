package newcmd

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/zhiruili/urem/core"
	"github.com/zhiruili/urem/osutil"
)

type NewFormatCmd struct {
	ProjectFile string `arg:"positional,required"`
}

func generateClangFormatFile(projectFilePath string) error {
	bs, err := core.Global.EmbedFs.ReadFile("resources/newformat/clang-format.tmpl")
	if err != nil {
		return fmt.Errorf("load ignore file template: %w", err)
	}

	dir := filepath.Dir(projectFilePath)
	outFile := filepath.Join(dir, ".clang-format")
	core.LogD("write file to %s", outFile)
	return ioutil.WriteFile(outFile, bs, 0644)
}

func (cmd *NewFormatCmd) Run() error {
	return osutil.DoInProjectRoot(cmd.ProjectFile, generateClangFormatFile)
}
