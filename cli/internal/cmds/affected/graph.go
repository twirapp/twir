package affected

import (
	"sort"
	"strings"
)

type depGraph struct {
	// appDeps maps "runtime:appName" -> direct lib deps (just lib names, no runtime prefix)
	appDeps map[string][]string
	// libDeps maps "runtime:libName" -> transitive lib deps (just lib names)
	libDeps map[string][]string
}

func newDepGraph(apps []appDeps, rawLibDeps map[string][]string) *depGraph {
	g := &depGraph{
		appDeps: map[string][]string{},
		libDeps: map[string][]string{},
	}

	for _, app := range apps {
		key := app.Runtime + ":" + app.Name
		g.appDeps[key] = app.Deps
	}

	for key, deps := range rawLibDeps {
		g.libDeps[key] = deps
	}

	return g
}

func (g *depGraph) computeAffected(changedLibs map[string]bool) []string {
	affected := map[string]bool{}

	for key, directDeps := range g.appDeps {
		parts := strings.SplitN(key, ":", 2)
		runtime := parts[0]
		appName := parts[1]

		for _, dep := range directDeps {
			if changedLibs[dep] {
				affected[appName] = true
				break
			}

			if g.transitivelyDependsOn(runtime, dep, changedLibs, map[string]bool{}) {
				affected[appName] = true
				break
			}
		}
	}

	result := make([]string, 0, len(affected))
	for name := range affected {
		result = append(result, name)
	}
	sort.Strings(result)
	return result
}

func (g *depGraph) transitivelyDependsOn(runtime, libName string, targets map[string]bool, visited map[string]bool) bool {
	key := runtime + ":" + libName
	if visited[key] {
		return false
	}
	visited[key] = true

	deps, ok := g.libDeps[key]
	if !ok {
		return false
	}

	for _, dep := range deps {
		if targets[dep] {
			return true
		}
		if g.transitivelyDependsOn(runtime, dep, targets, visited) {
			return true
		}
	}

	return false
}

func collectChangedLibs(changedFiles []string) map[string]bool {
	changedLibs := map[string]bool{}

	for _, file := range changedFiles {
		rest := file

		if idx := strings.Index(rest, "libs/"); idx >= 0 {
			rest = rest[idx+5:]
			if slashIdx := strings.Index(rest, "/"); slashIdx > 0 {
				libName := rest[:slashIdx]
				changedLibs[libName] = true
			}
		}
	}

	return changedLibs
}

func isRootFile(file string) bool {
	base := file
	if idx := strings.LastIndex(file, "/"); idx >= 0 {
		base = file[idx+1:]
	}

	rootFiles := map[string]bool{
		"go.work":      true,
		"go.work.sum":  true,
		"package.json":  true,
		"bun.lock":      true,
		".bun-version":  true,
		"bunfig.toml":   true,
	}

	if rootFiles[base] {
		return true
	}

	if strings.Contains(file, ".github/workflows/dock") {
		return true
	}

	return false
}
