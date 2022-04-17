package osutil

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/zhiruili/urem/core"
)

func IsDir(path string) (bool, error) {
	stat, err := os.Stat(path)

	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return false, nil
		}

		return false, err
	}

	return stat.IsDir(), nil
}

func MkDirIfNotExisted(path string) error {
	yes, err := IsDir(path)
	if err != nil {
		return err
	}
	if yes {
		return fmt.Errorf("directory already existed")
	}

	return os.MkdirAll(path, os.ModePerm)
}

func firstDirBottomUp(p string) (string, error) {
	for {
		if yes, err := IsDir(p); err != nil {
			return "", err
		} else if yes {
			break
		} else {
			parent := filepath.Dir(p)
			if parent != p {
				p = parent
			} else {
				return "", nil
			}
		}
	}

	return p, nil
}

func FindFileBottomUp(p string, exts ...string) (string, error) {
	d, err := firstDirBottomUp(p)
	if err != nil {
		return "", err
	}

	files, err := ioutil.ReadDir(d)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if !file.IsDir() {
			fext := filepath.Ext(file.Name())
			if core.StrContains(exts, fext) {
				return filepath.Join(d, file.Name()), nil
			}
		}
	}

	parent := filepath.Dir(d)
	if parent == d {
		return "", nil
	}

	return FindFileBottomUp(parent, exts...)
}
