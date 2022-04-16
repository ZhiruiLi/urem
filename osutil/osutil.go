package osutil

import (
	"errors"
	"fmt"
	"os"
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
