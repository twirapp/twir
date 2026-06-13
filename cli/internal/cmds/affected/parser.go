package affected

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type appDeps struct {
	Name    string
	Path    string
	Runtime string // "go" or "js"
	Deps    []string
}

func parseGoModDeps(appDir string) ([]string, error) {
	gomodPath := filepath.Join(appDir, "go.mod")
	f, err := os.Open(gomodPath)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", gomodPath, err)
	}
	defer f.Close()

	var deps []string
	seen := map[string]bool{}
	inRequire := false

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "require (" {
			inRequire = true
			continue
		}
		if inRequire && line == ")" {
			inRequire = false
			continue
		}

		if inRequire || strings.HasPrefix(line, "require ") {
			entry := strings.TrimPrefix(line, "require ")
			entry = strings.TrimSpace(entry)
			if entry == "" || strings.HasPrefix(entry, "//") {
				continue
			}

			parts := strings.Fields(entry)
			if len(parts) < 1 {
				continue
			}

			modPath := parts[0]
			const prefix = "github.com/twirapp/twir/libs/"
			if strings.HasPrefix(modPath, prefix) {
				libName := strings.TrimPrefix(modPath, prefix)
				libName = strings.Split(libName, "/")[0]
				if !seen[libName] {
					seen[libName] = true
					deps = append(deps, libName)
				}
			}
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("scan %s: %w", gomodPath, err)
	}

	return deps, nil
}

type packageJSON struct {
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

func parsePackageJSONDeps(appDir string) ([]string, error) {
	pkgPath := filepath.Join(appDir, "package.json")
	f, err := os.Open(pkgPath)
	if err != nil {
		return nil, fmt.Errorf("open %s: %w", pkgPath, err)
	}
	defer f.Close()

	var pkg packageJSON
	if err := json.NewDecoder(f).Decode(&pkg); err != nil {
		return nil, fmt.Errorf("decode %s: %w", pkgPath, err)
	}

	var deps []string
	seen := map[string]bool{}

	for name, version := range pkg.Dependencies {
		if version == "workspace:*" && strings.HasPrefix(name, "@twir/") {
			libName := strings.TrimPrefix(name, "@twir/")
			if !seen[libName] {
				seen[libName] = true
				deps = append(deps, libName)
			}
		}
	}

	for name, version := range pkg.DevDependencies {
		if version == "workspace:*" && strings.HasPrefix(name, "@twir/") {
			libName := strings.TrimPrefix(name, "@twir/")
			if !seen[libName] {
				seen[libName] = true
				deps = append(deps, libName)
			}
		}
	}

	return deps, nil
}

func discoverApps(rootDir string) ([]appDeps, error) {
	entries, err := os.ReadDir(filepath.Join(rootDir, "apps"))
	if err != nil {
		return nil, fmt.Errorf("read apps dir: %w", err)
	}

	var apps []appDeps

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		appDir := filepath.Join(rootDir, "apps", entry.Name())

		gomodPath := filepath.Join(appDir, "go.mod")
		pkgPath := filepath.Join(appDir, "package.json")

		if _, err := os.Stat(gomodPath); err == nil {
			deps, err := parseGoModDeps(appDir)
			if err != nil {
				return nil, fmt.Errorf("parse go deps for %s: %w", entry.Name(), err)
			}
			apps = append(apps, appDeps{
				Name:    entry.Name(),
				Path:    filepath.Join("apps", entry.Name()),
				Runtime: "go",
				Deps:    deps,
			})
		} else if _, err := os.Stat(pkgPath); err == nil {
			deps, err := parsePackageJSONDeps(appDir)
			if err != nil {
				return nil, fmt.Errorf("parse js deps for %s: %w", entry.Name(), err)
			}
			apps = append(apps, appDeps{
				Name:    entry.Name(),
				Path:    filepath.Join("apps", entry.Name()),
				Runtime: "js",
				Deps:    deps,
			})
		}
	}

	frontendDir := filepath.Join(rootDir, "frontend")
	if frontendEntries, err := os.ReadDir(frontendDir); err == nil {
		for _, entry := range frontendEntries {
			if !entry.IsDir() {
				continue
			}

			appDir := filepath.Join(frontendDir, entry.Name())
			pkgPath := filepath.Join(appDir, "package.json")

			if _, err := os.Stat(pkgPath); err == nil {
				deps, err := parsePackageJSONDeps(appDir)
				if err != nil {
					return nil, fmt.Errorf("parse js deps for frontend/%s: %w", entry.Name(), err)
				}
				apps = append(apps, appDeps{
					Name:    entry.Name(),
					Path:    filepath.Join("frontend", entry.Name()),
					Runtime: "js",
					Deps:    deps,
				})
			}
		}
	}

	webDir := filepath.Join(rootDir, "web")
	if _, err := os.Stat(filepath.Join(webDir, "package.json")); err == nil {
		deps, err := parsePackageJSONDeps(webDir)
		if err != nil {
			return nil, fmt.Errorf("parse js deps for web: %w", err)
		}
		apps = append(apps, appDeps{
			Name:    "web",
			Path:    "web",
			Runtime: "js",
			Deps:    deps,
		})
	}

	return apps, nil
}

func discoverLibDeps(rootDir string) (map[string][]string, error) {
	libDeps := map[string][]string{}

	libsDir := filepath.Join(rootDir, "libs")
	entries, err := os.ReadDir(libsDir)
	if err != nil {
		return nil, fmt.Errorf("read libs dir: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		libDir := filepath.Join(libsDir, entry.Name())
		gomodPath := filepath.Join(libDir, "go.mod")

		if _, err := os.Stat(gomodPath); err == nil {
			deps, err := parseGoModDeps(libDir)
			if err != nil {
				return nil, fmt.Errorf("parse go deps for lib %s: %w", entry.Name(), err)
			}
			if len(deps) > 0 {
				key := "go:" + entry.Name()
				libDeps[key] = deps
			}
		}

		pkgPath := filepath.Join(libDir, "package.json")
		if _, err := os.Stat(pkgPath); err == nil {
			deps, err := parsePackageJSONDeps(libDir)
			if err != nil {
				return nil, fmt.Errorf("parse js deps for lib %s: %w", entry.Name(), err)
			}
			if len(deps) > 0 {
				key := "js:" + entry.Name()
				libDeps[key] = deps
			}
		}
	}

	return libDeps, nil
}
