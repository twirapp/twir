package i18n

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
	"github.com/pterm/pterm"
)

func NewStore(dir string) (LocalesStore, error) {
	store := LocalesStore{}

	locales, err := os.ReadDir(dir)
	if err != nil {
		return store, fmt.Errorf("error reading locales directory: %w", err)
	}

	for _, f := range locales {
		if !f.IsDir() {
			continue
		}

		locale := f.Name()
		localePath := filepath.Join(dir, locale)

		localeSubDirs, err := os.ReadDir(localePath)
		if err != nil {
			return store, fmt.Errorf("error reading locale directory %s: %w", locale, err)
		}

		for _, subDir := range localeSubDirs {
			if !subDir.IsDir() {
				continue
			}

			subDirName := subDir.Name()
			subDirPath := filepath.Join(localePath, subDirName)

			files, err := os.ReadDir(subDirPath)
			if err != nil {
				return store, fmt.Errorf("error reading subdirectory %s: %w", subDirPath, err)
			}

			for _, file := range files {
				if file.IsDir() {
					continue
				}

				filePath := filepath.Join(subDirPath, file.Name())
				fileContent, err := os.ReadFile(filePath)
				if err != nil {
					return store, fmt.Errorf("error reading file %s: %w", filePath, err)
				}

				var translations map[string]string
				if err := yaml.Unmarshal(fileContent, &translations); err != nil {
					return store, fmt.Errorf("error unmarshaling YAML file %s: %w", filePath, err)
				}

				if _, ok := store[locale]; !ok {
					store[locale] = map[string]map[string]map[string]string{}
				}

				if _, ok := store[locale][subDirName]; !ok {
					store[locale][subDirName] = map[string]map[string]string{}
				}

				key := file.Name()[0 : len(file.Name())-len(filepath.Ext(file.Name()))]
				store[locale][subDirName][key] = translations

				pterm.Success.Printfln(
					"Loaded %d translations for %s/%s/%s",
					len(translations),
					locale,
					subDirName,
					key,
				)
			}
		}
	}

	return store, nil
}
