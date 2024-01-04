package shell

import (
	"runtime"
)

func GetShell() string {
	if runtime.GOOS == "windows" {
		return "cmd"
	} else {
		return "sh"
	}
}

func GetShellOption() string {
	if runtime.GOOS == "windows" {
		return "/C"
	} else {
		return "-c"
	}
}
