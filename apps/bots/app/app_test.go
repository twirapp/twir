package app

import (
	"testing"

	"go.uber.org/fx"
)

func TestAppHasCompleteDependencyGraph(t *testing.T) {
	if err := fx.ValidateApp(App); err != nil {
		t.Fatalf("validate bots dependency graph: %v", err)
	}
}
