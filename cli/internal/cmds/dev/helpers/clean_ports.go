package helpers

import (
	"errors"
	"fmt"
	"runtime"
	"strings"

	"github.com/pterm/pterm"
	"github.com/twirapp/twir/cli/internal/shell"
	"github.com/twirapp/twir/libs/grpc/constants"
	"github.com/urfave/cli/v2"
)

var CleanPortsCmd = &cli.Command{
	Name:    "clean-ports",
	Usage:   "clean ports needed for twir(maybe dangerous)",
	Aliases: []string{"clean"},
	Flags: []cli.Flag{
		&cli.BoolFlag{
			Name:  "kill-defaults",
			Value: true,
			Usage: "Will kill default twir applications",
		},
		&cli.IntSliceFlag{
			Name:  "additional-ports",
			Usage: "Set additional ports to clean",
		},
		&cli.IntFlag{
			Name:  "linux-signal",
			Value: 9,
			Usage: "Signal to send to the process. Defaults to SIGKILL.",
		},
	},
	Action: func(ctx *cli.Context) error {
		ports := []int{}
		if ctx.Bool("kill-defaults") {
			ports = append(ports, constants.ServerPorts...)
		}
		ports = append(ports, ctx.IntSlice("additional-ports")...)

		signal := ctx.Int("linux-signal")
		if signal <= 1 || signal >= 73 {
			pterm.Fatal.Println("Signal must be a number between 1 and 73")
		}

		pterm.Warning.Println("Cleaning up ports", ports)
		if runtime.GOOS == "windows" {
			pterm.Fatal.Println("Windows is not supported yet")
			return cleanPortsWindows()
		} else if runtime.GOOS == "linux" {
			err := cleanPortsLinux(ctx.Int("signal"), ports...)
			if err != nil {
				return err
			}
			pterm.Success.Println("Ports cleaned")
		}

		return nil
	},
}

func cleanPortsLinux(signal int, ports ...int) error {
	for _, port := range ports {
		exec, err := shell.CreateCommand(shell.ExecCommandOpts{
			Command: fmt.Sprintf("kill -%d $(lsof -t -i:%d)", signal, port),
		})
		if err != nil {
			return err
		}

		err = exec.Run()
		if err != nil {
			if !strings.Contains(err.Error(), "exit status 2") {
				return err
			}
		}
	}

	return nil
}

func cleanPortsWindows() error {
	return errors.New("not implemented")
}
