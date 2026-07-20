package dota

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/dota/model"
)

type matchResultApplyingRepository interface {
	ApplyMatchResultOnce(
		context.Context,
		ApplyMatchResultInput,
	) (model.ChannelDotaSettings, error)
}

var _ matchResultApplyingRepository = (Repository)(nil)

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

func TestValidateMatchResultInput(t *testing.T) {
	channelID := uuid.New()

	tests := []struct {
		name  string
		input ApplyMatchResultInput
		want  string
	}{
		{
			name: "requires a channel ID",
			input: ApplyMatchResultInput{
				MatchID: 1,
			},
			want: "channel ID is required",
		},
		{
			name: "requires a positive match ID",
			input: ApplyMatchResultInput{
				ChannelID: channelID,
				MatchID:   0,
			},
			want: "match ID must be positive",
		},
		{
			name: "rejects a negative match ID",
			input: ApplyMatchResultInput{
				ChannelID: channelID,
				MatchID:   -1,
			},
			want: "match ID must be positive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateMatchResultInput(tt.input)
			if err == nil {
				t.Fatal("ValidateMatchResultInput() error = nil, want an error")
			}
			if err.Error() != tt.want {
				t.Errorf("ValidateMatchResultInput() error = %q, want %q", err, tt.want)
			}
		})
	}
}
