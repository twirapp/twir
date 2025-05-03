package main

import (
	"log"
	"os"

	"github.com/twirapp/twir/cli/internal/cmds/build"
	"github.com/twirapp/twir/cli/internal/cmds/dependencies"
	"github.com/twirapp/twir/cli/internal/cmds/dev"
	"github.com/twirapp/twir/cli/internal/cmds/execbin"
	"github.com/twirapp/twir/cli/internal/cmds/generate"
	"github.com/twirapp/twir/cli/internal/cmds/kill"
	"github.com/twirapp/twir/cli/internal/cmds/migrations"
	"github.com/twirapp/twir/cli/internal/cmds/proxy"
	"github.com/urfave/cli/v2"
)

func main() {
	// var rootCmd = &cobra.Command{
	// 	Use:   "cli",
	// 	Short: "Cli for manage twir project",
	// 	PreRunE: func(cmd *cobra.Command, args []string) error {
	// 		return dependencies.InstallDependencies().Execute()
	// 	},
	// 	RunE: func(cmd *cobra.Command, args []string) error {
	// 		return nil
	// 	},
	// 	CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	// }
	//
	// rootCmd.AddCommand(dependencies.InstallDependencies())
	// rootCmd.AddCommand(migrations.Migrations())
	//
	// if err := rootCmd.Execute(); err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }

	app := &cli.App{
		Name:        "go run cmd/main.go",
		Description: "TwirApp cli for helping in manage project",
		Commands: []*cli.Command{
			dependencies.Cmd,
			migrations.Cmd,
			proxy.Cmd,
			generate.Cmd,
			build.Cmd,
			dev.CreateDevCommand(),
			execbin.Cmd,
			kill.Cmd,
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
