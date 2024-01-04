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

func ExecCommand(opts ExecCommandOpts) error {
	if opts.Command == "" {
		return fmt.Errorf("command not specified")
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

	return cmd.Run()
}
