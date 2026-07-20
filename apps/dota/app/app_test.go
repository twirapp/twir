package app

import (
	"testing"

	"go.uber.org/fx"
)

func TestFxGraphValidates(t *testing.T) {
	if err := fx.ValidateApp(App); err != nil {
		t.Fatal(err)
	}
}
