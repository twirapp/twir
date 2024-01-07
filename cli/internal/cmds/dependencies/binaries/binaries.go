package binaries

import (
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
