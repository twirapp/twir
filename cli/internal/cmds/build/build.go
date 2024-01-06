package build

import (
	"fmt"
	"math"
	"os"
	"strings"
	"time"

	"github.com/pterm/pterm"
	"github.com/twirapp/twir/cli/internal/shell"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:  "build",
	Usage: "build application",
	Action: func(c *cli.Context) error {
		return build("turbo run build --filter=!./apps/dota")
	},
	Subcommands: []*cli.Command{
		LibsCmd,
	},
}

var LibsCmd = &cli.Command{
	Name: "libs",
	Action: func(context *cli.Context) error {
		return build("turbo run build --filter='./libs/*'")
	},
}

func rgb(i int) (int, int, int) {
	var f = 0.275

	return int(math.Sin(f*float64(i)+4*math.Pi/3)*127 + 128),
		// int(math.Sin(f*float64(i)+2*math.Pi/3)*127 + 128),
		int(45),
		int(math.Sin(f*float64(i)+0)*127 + 128)
}

func rainbow(text string) string {
	var rainbowStr []string
	for index, value := range text {
		r, g, b := rgb(index)
		str := fmt.Sprintf("\033[1m\033[38;2;%d;%d;%dm%c\033[0m\033[0;1m", r, g, b, value)
		rainbowStr = append(rainbowStr, str)
	}

	return strings.Join(rainbowStr, "")
}

func build(cmd string) error {
	wd, err := os.Getwd()
	if err != nil {
		return err
	}

	pterm.Info.Println("Building twir")
	spinner, _ := pterm.DefaultSpinner.Start("Building...")

	startTime := time.Now()

	err = shell.ExecCommand(
		shell.ExecCommandOpts{
			Command: cmd,
			Pwd:     wd,
		},
	)
	if err != nil {
		spinner.Fail(err)
		return err
	}

	if time.Since(startTime).Milliseconds() < 1000 {
		spinner.Success(rainbow(">>> FULL TWIR TURBO ") + "🤙 🤙 🤙")
	} else {
		spinner.Success("Builded")
	}

	return nil
}
