package generate

import (
	"github.com/twirapp/twir/cli/internal/cmds/generate/dockerfile"
	"github.com/twirapp/twir/cli/internal/cmds/generate/locales"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:        "generate",
	Usage:       "some generators",
	Aliases:     []string{"gen", "g"},
	Subcommands: []*cli.Command{dockerfile.Dockerfile, locales.Cmd},
}
