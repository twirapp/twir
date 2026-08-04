package goapp

import (
	"os"
	"path/filepath"
	"testing"
)

func TestWireTargets(t *testing.T) {
	repositoryRoot := t.TempDir()
	appPath := filepath.Join(repositoryRoot, "apps", "events")
	baseInjector := filepath.Join(repositoryRoot, "libs", "baseapp", "wire.go")
	appInjector := filepath.Join(appPath, "cmd", "wire.go")

	for _, injector := range []string{baseInjector, appInjector} {
		if err := os.MkdirAll(filepath.Dir(injector), 0o755); err != nil {
			t.Fatalf("create injector directory: %v", err)
		}
		if err := os.WriteFile(injector, []byte("package placeholder\n"), 0o644); err != nil {
			t.Fatalf("create injector: %v", err)
		}
	}

	targets, err := wireTargets(repositoryRoot, appPath)
	if err != nil {
		t.Fatalf("get Wire targets: %v", err)
	}

	want := []string{
		filepath.Join("..", "libs", "baseapp"),
		filepath.Join("..", "apps", "events", "cmd"),
	}
	if len(targets) != len(want) {
		t.Fatalf("got %d targets, want %d: %v", len(targets), len(want), targets)
	}
	for index := range want {
		if targets[index] != want[index] {
			t.Fatalf("target %d = %q, want %q", index, targets[index], want[index])
		}
	}
}

func TestWireTargetsSkipsMissingInjectors(t *testing.T) {
	repositoryRoot := t.TempDir()
	appPath := filepath.Join(repositoryRoot, "apps", "parser")

	targets, err := wireTargets(repositoryRoot, appPath)
	if err != nil {
		t.Fatalf("get Wire targets: %v", err)
	}
	if len(targets) != 0 {
		t.Fatalf("got targets for missing injectors: %v", targets)
	}
}
