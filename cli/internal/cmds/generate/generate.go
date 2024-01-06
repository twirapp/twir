package generate

import (
	"github.com/twirapp/twir/cli/internal/cmds/generate/dockerfile"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:        "generate",
	Usage:       "some generators",
	Aliases:     []string{"gen"},
	Subcommands: []*cli.Command{dockerfile.Dockerfile},
}
