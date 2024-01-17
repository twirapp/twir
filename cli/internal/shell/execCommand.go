package shell

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

type ExecCommandOpts struct {
	Command string
	Pwd     string

	Stdout io.Writer
	Stderr io.Writer
}

func CreateCommand(opts ExecCommandOpts) (*exec.Cmd, error) {
	if opts.Command == "" {
		return nil, fmt.Errorf("command not specified")
	}

	cmd := exec.Command(
		GetShell(),
		GetShellOption(),
		opts.Command,
	)

	pathVarKey := "PATH"
	path := os.Getenv(pathVarKey)
	pathDelimiter := ":"
	if runtime.GOOS == "windows" {
		pathVarKey = "Path"
		path = os.Getenv(pathVarKey)
		pathDelimiter = ";"
	}

	projectWd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	path += pathDelimiter + filepath.Join(projectWd, ".bin")

	cmd.Env = append(os.Environ(), pathVarKey+"="+path)

	if opts.Pwd != "" {
		cmd.Dir = opts.Pwd
	}

	cmd.Stdout = opts.Stdout
	cmd.Stderr = opts.Stderr

	return cmd, nil
}

func ExecCommand(opts ExecCommandOpts) error {
	cmd, err := CreateCommand(opts)
	if err != nil {
		return err
	}

	return cmd.Run()
}
