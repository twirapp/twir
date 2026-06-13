package modutil

import (
	"os"
	"path/filepath"
	"strings"

	"golang.org/x/mod/modfile"
)

type LibDependencyResolver struct {
	libToApps map[string][]string
	appToLibs map[string][]string
	workDir   string
}

func NewLibDependencyResolver() (*LibDependencyResolver, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	return &LibDependencyResolver{
		libToApps: make(map[string][]string),
		appToLibs: make(map[string][]string),
		workDir:   wd,
	}, nil
}

func (r *LibDependencyResolver) ResolveApp(appName string) error {
	appPath := filepath.Join(r.workDir, "apps", appName)
	gomodPath := filepath.Join(appPath, "go.mod")

	gomodBytes, err := os.ReadFile(gomodPath)
	if err != nil {
		return err
	}

	parsed, err := modfile.Parse("go.mod", gomodBytes, nil)
	if err != nil {
		return err
	}

	var appLibs []string

	for _, req := range parsed.Require {
		if req == nil || req.Mod.Path == "" {
			continue
		}

		if !strings.HasPrefix(req.Mod.Path, "github.com/twirapp/twir/libs/") {
			continue
		}

		libName := filepath.Base(req.Mod.Path)
		libPath := filepath.Join(r.workDir, "libs", libName)

		if _, err := os.Stat(libPath); os.IsNotExist(err) {
			continue
		}

		for _, rep := range parsed.Replace {
			if rep == nil || rep.Old.Path == "" {
				continue
			}
			if rep.Old.Path == req.Mod.Path {
				if rep.New.Path != "" {
					resolvedPath := rep.New.Path
					if !filepath.IsAbs(resolvedPath) {
						resolvedPath = filepath.Join(appPath, resolvedPath)
					}
					resolvedPath = filepath.Clean(resolvedPath)

					if strings.Contains(resolvedPath, string(filepath.Separator)+"libs"+string(filepath.Separator)) ||
						strings.HasSuffix(resolvedPath, string(filepath.Separator)+"libs") {
						libPath = resolvedPath
					}
				}
				break
			}
		}

		appLibs = append(appLibs, libPath)
		r.libToApps[libPath] = append(r.libToApps[libPath], appName)
	}

	r.appToLibs[appName] = appLibs
	return nil
}

func (r *LibDependencyResolver) GetAppsForLib(libPath string) []string {
	libPath = filepath.Clean(libPath)
	return r.libToApps[libPath]
}

func (r *LibDependencyResolver) GetAppsForFile(filePath string) []string {
	filePath = filepath.Clean(filePath)

	for {
		if apps, ok := r.libToApps[filePath]; ok {
			return apps
		}

		parent := filepath.Dir(filePath)
		if parent == filePath {
			break
		}
		filePath = parent
	}

	return nil
}

func (r *LibDependencyResolver) GetLibsForApp(appName string) []string {
	return r.appToLibs[appName]
}

func (r *LibDependencyResolver) GetAllWatchedPaths() []string {
	paths := make([]string, 0, len(r.libToApps))
	for path := range r.libToApps {
		paths = append(paths, path)
	}
	return paths
}
