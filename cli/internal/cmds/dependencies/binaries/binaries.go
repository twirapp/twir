package binaries

import (
	"os"
	"path/filepath"
)

func isBinaryInstalled(binaryName string) bool {
	wd, err := os.Getwd()
	if err != nil {
		return false
	}

	if _, err := os.Stat(filepath.Join(wd, ".bin", binaryName)); os.IsNotExist(err) {
		return false
	}

	return true
}
