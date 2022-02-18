package logger

import (
	"fmt"
	"os"
)

func isNotExist(src string) bool {
	_, err := os.Stat(src)
	return os.IsNotExist(err)
}

func isPermissionDenied(src string) bool {
	_, err := os.Stat(src)
	return os.IsPermission(err)
}

func isNotExistMkDir(src string) error {
	if isNotExist(src) {
		if err := mkDir(src); err != nil {
			return err
		}
	}
	return nil
}

func mkDir(src string) error {
	err := os.MkdirAll(src, os.ModePerm)
	if err != nil {
		return err
	}

	return nil
}

func mustOpen(fileName, dir string) (*os.File, error) {
	if isPermissionDenied(dir) {
		return nil, fmt.Errorf("permission denied dir: %s", dir)
	}

	if err := isNotExistMkDir(dir); err != nil {
		return nil, fmt.Errorf("error during make dir %s, err: %s", dir, err)
	}

	f, err := os.OpenFile(dir+string(os.PathSeparator)+fileName, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, fmt.Errorf("fail to open file, err: %s", err)
	}

	return f, nil
}
