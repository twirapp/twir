package binaries

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
)

func isBinaryInstalled(binaryName string) bool {
	wd, err := os.Getwd()
	if err != nil {
		return false
	}

	if runtime.GOOS == "windows" {
		binaryName += ".exe"
	}

	if _, err := os.Stat(filepath.Join(wd, ".bin", binaryName)); os.IsNotExist(err) {
		return false
	}

	return true
}

func CreateDir() error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	if err := os.Mkdir(filepath.Join(wd, ".bin"), 0755); err != nil && !errors.Is(err, os.ErrExist) {
		return err
	}

	return nil
}
