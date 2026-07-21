package app

import (
	"testing"

	"github.com/twirapp/twir/apps/dota/internal/predictions"
	"go.uber.org/fx"
)

func TestFxGraphValidates(t *testing.T) {
	if err := fx.ValidateApp(App); err != nil {
		t.Fatal(err)
	}
}

func TestFxGraphProvidesLifecycleWorker(t *testing.T) {
	if err := fx.ValidateApp(App, fx.Invoke(func(*predictions.LifecycleWorker) {})); err != nil {
		t.Fatal(err)
	}
}
