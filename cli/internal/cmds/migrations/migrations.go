package migrations

import (
	_ "github.com/twirapp/twir/libs/migrations/postgres"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:    "migrations",
	Aliases: []string{"m"},
	Usage:   "manage migrations",
	Subcommands: []*cli.Command{
		MigrateCmd,
		createCmd,
	},
}

type emptyLogWriter struct{}

func (c *emptyLogWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}
