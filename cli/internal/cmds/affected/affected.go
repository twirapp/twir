package affected

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
	"github.com/urfave/cli/v2"
)

var Cmd = &cli.Command{
	Name:  "affected",
	Usage: "determine which apps are affected by changed files",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:     "files",
			Usage:    "comma-separated list of changed file paths",
			Required: true,
		},
		&cli.StringFlag{
			Name:  "output",
			Usage: "output format: json or list",
			Value: "json",
		},
	},
	Action: func(c *cli.Context) error {
		filesRaw := c.String("files")
		outputFormat := c.String("output")

		rootDir, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("get working directory: %w", err)
		}

		if filesRaw == "ALL" {
			allApps := allAppNames(rootDir)
			printOutput(allApps, outputFormat)
			return nil
		}

		changedFiles := splitFiles(filesRaw)
		if len(changedFiles) == 0 {
			pterm.Warning.Println("No changed files provided")
			printOutput(nil, outputFormat)
			return nil
		}

		for _, f := range changedFiles {
			if isRootFile(f) {
				pterm.Info.Println("Root config file changed, all apps affected")
				allApps := allAppNames(rootDir)
				printOutput(allApps, outputFormat)
				return nil
			}
		}

		apps, err := discoverApps(rootDir)
		if err != nil {
			return fmt.Errorf("discover apps: %w", err)
		}

		libDeps, err := discoverLibDeps(rootDir)
		if err != nil {
			return fmt.Errorf("discover lib deps: %w", err)
		}

		graph := newDepGraph(apps, libDeps)
		changedLibs := collectChangedLibs(changedFiles)

		affected := graph.computeAffected(changedLibs)

		for _, file := range changedFiles {
			for _, app := range apps {
				if strings.HasPrefix(file, app.Path+"/") {
					if !contains(affected, app.Name) {
						affected = append(affected, app.Name)
					}
				}
			}
		}

		if strings.HasPrefix(filesRaw, "web/") || strings.Contains(filesRaw, ",web/") {
			if !contains(affected, "web") {
				affected = append(affected, "web")
			}
		}

		printOutput(affected, outputFormat)
		return nil
	},
}

func splitFiles(raw string) []string {
	parts := strings.Split(raw, ",")
	var result []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}
	return result
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func allAppNames(rootDir string) []string {
	var names []string

	entries, err := os.ReadDir(fmt.Sprintf("%s/apps", rootDir))
	if err == nil {
		for _, e := range entries {
			if e.IsDir() && e.Name() != "twitch-mock" {
				names = append(names, e.Name())
			}
		}
	}

	entries, err = os.ReadDir(fmt.Sprintf("%s/frontend", rootDir))
	if err == nil {
		for _, e := range entries {
			if e.IsDir() {
				names = append(names, e.Name())
			}
		}
	}

	names = append(names, "web")

	return names
}

func printOutput(apps []string, format string) {
	switch format {
	case "json":
		data, _ := json.Marshal(apps)
		fmt.Println(string(data))
	case "list":
		for _, app := range apps {
			fmt.Println(app)
		}
	default:
		data, _ := json.Marshal(apps)
		fmt.Println(string(data))
	}
}
