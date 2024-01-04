package dockerfile

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/urfave/cli/v2"
)

var Dockerfile = &cli.Command{
	Name:  "dockerfile",
	Usage: "dockerfile package.json and go.mod",
	Action: func(context *cli.Context) error {
		nodePaths := append(findPackageJsonFiles("./libs"), findPackageJsonFiles("./apps")...)
		nodePaths = append(nodePaths, findPackageJsonFiles("./frontend")...)

		goPaths := append(findGoModFiles("./libs"), findGoModFiles("./apps")...)

		var generatedStrings []string
		for _, p := range nodePaths {
			generatedStrings = append(
				generatedStrings,
				fmt.Sprintf("COPY %s/package.json %s/package.json", p, p),
			)
		}
		for _, p := range goPaths {
			generatedStrings = append(generatedStrings, fmt.Sprintf("COPY %s/go.mod %s/go.mod", p, p))
		}

		originalFile, err := os.ReadFile("./base.Dockerfile")
		if err != nil {
			return fmt.Errorf("error reading file: %w", err)
		}

		newFileContent := replaceBetweenComments(
			string(originalFile),
			strings.Join(generatedStrings, "\n"),
		)

		err = os.WriteFile("./base.Dockerfile", []byte(newFileContent), 0644)
		if err != nil {
			return fmt.Errorf("error writing file: %w", err)
		}

		return nil
	},
}

func findPackageJsonFiles(startPath string) []string {
	var result []string

	if _, err := os.Stat(startPath); os.IsNotExist(err) {
		fmt.Println("Directory path does not exist:", startPath)
		return result
	}

	files, err := os.ReadDir(startPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return result
	}

	for _, file := range files {
		fullPath := filepath.Join(startPath, file.Name())

		if !file.IsDir() {
			continue
		}

		if _, err := os.Stat(filepath.Join(fullPath, "package.json")); err == nil {
			result = append(result, fullPath)
		}
	}

	return result
}

func findGoModFiles(startPath string) []string {
	var result []string

	if _, err := os.Stat(startPath); os.IsNotExist(err) {
		fmt.Println("Directory path does not exist:", startPath)
		return result
	}

	files, err := os.ReadDir(startPath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return result
	}

	for _, file := range files {
		fullPath := filepath.Join(startPath, file.Name())

		if !file.IsDir() {
			continue
		}

		if _, err := os.Stat(filepath.Join(fullPath, "go.mod")); err == nil {
			result = append(result, fullPath)
		}
	}

	return result
}

func replaceBetweenComments(input string, newContent string) string {
	startComment := "# START COPYGEN"
	endComment := "# END COPYGEN"

	regex := fmt.Sprintf("%s[\\s\\S]*%s", startComment, endComment)
	replacedContent := regexp.MustCompile(regex).ReplaceAllString(
		input,
		fmt.Sprintf("%s\n%s\n%s", startComment, newContent, endComment),
	)

	return replacedContent
}
