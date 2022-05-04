package osutil

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/zhiruili/urem/core"
)

// IsDir 检查给定路径是否是一个目录。
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

// MkDirIfNotExisted 检查指定路径是否存在，如果不存在，就创建目录。
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

func matchAnyPattern(s string, patterns ...string) (string, bool) {
	for _, p := range patterns {
		ok, err := filepath.Match(p, s)
		if err != nil {
			core.LogE("illegal file pattern %s, when trying to match file %s", p, s)
			continue
		}
		if ok {
			return p, true
		}
	}

	return "", false
}

func findFirstFileMatchPattern(files []fs.FileInfo, patterns ...string) (string, bool) {
	for _, file := range files {
		if !file.IsDir() {
			_, match := matchAnyPattern(file.Name(), patterns...)
			if match {
				return file.Name(), true
			}
		}
	}
	return "", false
}

// FindFileBottomUp 自底向上寻找第一个符合 pattern 的文件。
func FindFileBottomUp(p string, patterns ...string) (string, error) {
	dir, err := firstDirBottomUp(p)
	if err != nil {
		return "", err
	}

	for {
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			return "", err
		}
		name, ok := findFirstFileMatchPattern(files, patterns...)
		if ok {
			return filepath.Join(dir, name), nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", nil
		}

		dir = parent
	}
}

// DoInProjectRoot 在 UE 的 project root 执行 work 操作。
func DoInProjectRoot(path string, work func(string) error) error {
	filePath, err := FindFileBottomUp(path, "*.uproject")
	if err != nil {
		return fmt.Errorf("find .uproject file: %w", err)
	}

	if filePath == "" {
		return fmt.Errorf(".uproject file no found")
	}

	absFilePath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("illegal project file path: %s", filePath)
	}

	core.LogD("project file path: %s", absFilePath)
	return work(absFilePath)
}

// CopyFile 复制文件到指定路径。
func CopyFile(src string, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
