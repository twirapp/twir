package i18n

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/goccy/go-yaml"
)

// LocalesStore example: {"en": {"commands": {"followage": {"title": "qwe"}}}}
type LocalesStore map[string]map[string]map[string]map[string]string

// RawLocalesStore preserves the nested structure for key generation
type RawLocalesStore map[string]map[string]map[string]interface{}

// flattenMap recursively flattens nested maps using dot notation
func flattenMap(data interface{}, prefix string) map[string]string {
	result := make(map[string]string)

	switch v := data.(type) {
	case map[string]interface{}:
		for key, value := range v {
			newKey := key
			if prefix != "" {
				newKey = prefix + "." + key
			}

			if nestedMap, ok := value.(map[string]interface{}); ok {
				// Recursively flatten nested maps
				for k, v := range flattenMap(nestedMap, newKey) {
					result[k] = v
				}
			} else {
				// Convert value to string
				result[newKey] = fmt.Sprintf("%v", value)
			}
		}
	case map[interface{}]interface{}:
		// Handle YAML's default unmarshaling format
		for key, value := range v {
			keyStr := fmt.Sprintf("%v", key)
			newKey := keyStr
			if prefix != "" {
				newKey = prefix + "." + keyStr
			}

			if nestedMap, ok := value.(map[interface{}]interface{}); ok {
				// Recursively flatten nested maps
				for k, v := range flattenMap(nestedMap, newKey) {
					result[k] = v
				}
			} else if nestedStringMap, ok := value.(map[string]interface{}); ok {
				// Recursively flatten nested maps
				for k, v := range flattenMap(nestedStringMap, newKey) {
					result[k] = v
				}
			} else {
				// Convert value to string
				result[newKey] = fmt.Sprintf("%v", value)
			}
		}
	default:
		// For primitive values
		if prefix != "" {
			result[prefix] = fmt.Sprintf("%v", v)
		}
	}

	return result
}

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

				var rawTranslations interface{}
				if err := yaml.Unmarshal(fileContent, &rawTranslations); err != nil {
					return store, fmt.Errorf("error unmarshaling YAML file %s: %w", filePath, err)
				}

				// Flatten the nested structure
				translations := flattenMap(rawTranslations, "")

				if _, ok := store[locale]; !ok {
					store[locale] = map[string]map[string]map[string]string{}
				}

				if _, ok := store[locale][subDirName]; !ok {
					store[locale][subDirName] = map[string]map[string]string{}
				}

				key := file.Name()[0 : len(file.Name())-len(filepath.Ext(file.Name()))]
				store[locale][subDirName][key] = translations
			}
		}
	}

	return store, nil
}

// LoadRawStore loads YAML files preserving their nested structure for key generation
func LoadRawStore(dir string) (RawLocalesStore, error) {
	store := RawLocalesStore{}

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

				var rawTranslations interface{}
				if err := yaml.Unmarshal(fileContent, &rawTranslations); err != nil {
					return store, fmt.Errorf("error unmarshaling YAML file %s: %w", filePath, err)
				}

				if _, ok := store[locale]; !ok {
					store[locale] = map[string]map[string]interface{}{}
				}

				if _, ok := store[locale][subDirName]; !ok {
					store[locale][subDirName] = map[string]interface{}{}
				}

				key := file.Name()[0 : len(file.Name())-len(filepath.Ext(file.Name()))]
				store[locale][subDirName][key] = rawTranslations
			}
		}
	}

	return store, nil
}
