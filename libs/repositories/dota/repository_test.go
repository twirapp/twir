package dota

import (
	"testing"

	"github.com/twirapp/twir/libs/repositories/dota/model"
)

func TestCommandSettingsOrDefault(t *testing.T) {
	t.Run("nil enables every Dota command", func(t *testing.T) {
		want := model.CommandsSettings{
			Mmr: true,
			Wl:  true,
			Lg:  true,
			Gm:  true,
			Np:  true,
			Wp:  true,
		}

		if got := CommandSettingsOrDefault(nil); got != want {
			t.Errorf("CommandSettingsOrDefault(nil) = %#v, want %#v", got, want)
		}
	})

	t.Run("explicit all-false settings are preserved", func(t *testing.T) {
		explicitlyDisabled := model.CommandsSettings{}

		if got := CommandSettingsOrDefault(&explicitlyDisabled); got != explicitlyDisabled {
			t.Errorf("CommandSettingsOrDefault(%#v) = %#v, want %#v", explicitlyDisabled, got, explicitlyDisabled)
		}
	})
}
