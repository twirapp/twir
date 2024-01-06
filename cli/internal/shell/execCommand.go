package shell

import (
	"fmt"
	"io"
	"os/exec"
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

	if opts.Pwd != "" {
		cmd.Dir = opts.Pwd
	}

	if opts.Stdout != nil {
		cmd.Stdout = opts.Stdout
	}

	if opts.Stderr != nil {
		cmd.Stderr = opts.Stderr
	}

	return cmd, nil
}

func ExecCommand(opts ExecCommandOpts) error {
	cmd, err := CreateCommand(opts)
	if err != nil {
		return err
	}

	return cmd.Run()
}
