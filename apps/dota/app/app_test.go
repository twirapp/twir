package app

import (
	"testing"

	"github.com/twirapp/twir/apps/dota/internal/buslistener"
	"go.uber.org/fx"
)

func TestFxGraphValidates(t *testing.T) {
	if err := fx.ValidateApp(App); err != nil {
		t.Fatal(err)
	}
}

func TestFxGraphProvidesStatsForBusListener(t *testing.T) {
	if err := fx.ValidateApp(
		App,
		fx.Invoke(func(_ buslistener.StatsProvider) {}),
	); err != nil {
		t.Fatal(err)
	}
}
