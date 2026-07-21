package dota

import (
	"context"
	"encoding/json"
	"testing"
	"time"

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

type matchLifecycleRepository interface {
	GetMatchState(context.Context, uuid.UUID) (model.MatchState, error)
	ApplyMatchStateTransition(context.Context, ApplyMatchStateTransitionInput) (bool, error)
	ClaimPredictionActions(context.Context, ClaimPredictionActionsInput) ([]model.ClaimedOutboxAction, error)
	CompletePredictionAction(context.Context, uuid.UUID, uuid.UUID) error
	RetryPredictionAction(context.Context, uuid.UUID, uuid.UUID, time.Time) error
}

var _ matchLifecycleRepository = (Repository)(nil)

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

func TestValidateOutboxActionInput(t *testing.T) {
	channelID := uuid.New()
	validInput := func() model.OutboxActionInput {
		return model.OutboxActionInput{
			ChannelID: channelID,
			MatchID:   42,
			Action:    model.OutboxActionCreate,
			Sequence:  10,
			Payload:   json.RawMessage(`{"kind":"create"}`),
		}
	}

	tests := []struct {
		name  string
		input model.OutboxActionInput
		want  bool
	}{
		{
			name:  "accepts a valid action",
			input: validInput(),
		},
		{
			name: "requires a channel ID",
			input: func() model.OutboxActionInput {
				input := validInput()
				input.ChannelID = uuid.Nil
				return input
			}(),
			want: true,
		},
		{
			name: "rejects a zero match ID",
			input: func() model.OutboxActionInput {
				input := validInput()
				input.MatchID = 0
				return input
			}(),
			want: true,
		},
		{
			name: "rejects a negative match ID",
			input: func() model.OutboxActionInput {
				input := validInput()
				input.MatchID = -1
				return input
			}(),
			want: true,
		},
		{
			name: "rejects an unsupported action",
			input: func() model.OutboxActionInput {
				input := validInput()
				input.Action = "unknown"
				return input
			}(),
			want: true,
		},
		{
			name: "rejects a zero sequence",
			input: func() model.OutboxActionInput {
				input := validInput()
				input.Sequence = 0
				return input
			}(),
			want: true,
		},
		{
			name: "rejects a negative sequence",
			input: func() model.OutboxActionInput {
				input := validInput()
				input.Sequence = -1
				return input
			}(),
			want: true,
		},
		{
			name: "requires a payload",
			input: func() model.OutboxActionInput {
				input := validInput()
				input.Payload = json.RawMessage(" \t")
				return input
			}(),
			want: true,
		},
		{
			name: "rejects invalid payload JSON",
			input: func() model.OutboxActionInput {
				input := validInput()
				input.Payload = json.RawMessage(`{`)
				return input
			}(),
			want: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateOutboxActionInput(tt.input)
			if got := err != nil; got != tt.want {
				t.Errorf("ValidateOutboxActionInput() error = %v, want error = %t", err, tt.want)
			}
		})
	}
}

func TestValidateApplyMatchStateTransitionInput(t *testing.T) {
	channelID := uuid.New()

	validInput := func() ApplyMatchStateTransitionInput {
		return ApplyMatchStateTransitionInput{
			ChannelID:         channelID,
			ExpectedRevision:  0,
			ProviderTimestamp: 1_700_000_000,
			Snapshot:          json.RawMessage(`{"state":"active"}`),
			Actions: []model.OutboxActionInput{{
				ChannelID: channelID,
				MatchID:   42,
				Action:    model.OutboxActionCreate,
				Sequence:  10,
				Payload:   json.RawMessage(`{"kind":"create"}`),
			}},
		}
	}

	tests := []struct {
		name  string
		input ApplyMatchStateTransitionInput
	}{
		{
			name: "requires a channel ID",
			input: ApplyMatchStateTransitionInput{
				Snapshot: json.RawMessage(`{}`),
			},
		},
		{
			name: "rejects a negative expected revision",
			input: func() ApplyMatchStateTransitionInput {
				input := validInput()
				input.ExpectedRevision = -1
				return input
			}(),
		},
		{
			name: "rejects a negative provider timestamp",
			input: func() ApplyMatchStateTransitionInput {
				input := validInput()
				input.ProviderTimestamp = -1
				return input
			}(),
		},
		{
			name: "requires a snapshot",
			input: func() ApplyMatchStateTransitionInput {
				input := validInput()
				input.Snapshot = nil
				return input
			}(),
		},
		{
			name: "rejects invalid snapshot JSON",
			input: func() ApplyMatchStateTransitionInput {
				input := validInput()
				input.Snapshot = json.RawMessage(`{`)
				return input
			}(),
		},
		{
			name: "rejects an action for another channel",
			input: func() ApplyMatchStateTransitionInput {
				input := validInput()
				input.Actions[0].ChannelID = uuid.New()
				return input
			}(),
		},
		{
			name: "rejects a non-positive match ID",
			input: func() ApplyMatchStateTransitionInput {
				input := validInput()
				input.Actions[0].MatchID = 0
				return input
			}(),
		},
		{
			name: "rejects an unknown action",
			input: func() ApplyMatchStateTransitionInput {
				input := validInput()
				input.Actions[0].Action = "unknown"
				return input
			}(),
		},
		{
			name: "rejects a non-positive action sequence",
			input: func() ApplyMatchStateTransitionInput {
				input := validInput()
				input.Actions[0].Sequence = 0
				return input
			}(),
		},
		{
			name: "requires valid action payload JSON",
			input: func() ApplyMatchStateTransitionInput {
				input := validInput()
				input.Actions[0].Payload = json.RawMessage(`{`)
				return input
			}(),
		},
		{
			name: "rejects duplicate match action pairs",
			input: func() ApplyMatchStateTransitionInput {
				input := validInput()
				input.Actions = append(input.Actions, input.Actions[0])
				return input
			}(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateApplyMatchStateTransitionInput(tt.input); err == nil {
				t.Fatal("ValidateApplyMatchStateTransitionInput() error = nil, want an error")
			}
		})
	}

	if err := ValidateApplyMatchStateTransitionInput(validInput()); err != nil {
		t.Fatalf("ValidateApplyMatchStateTransitionInput() error = %v", err)
	}
}

func TestValidateClaimPredictionActionsInput(t *testing.T) {
	tests := []struct {
		name  string
		input ClaimPredictionActionsInput
	}{
		{
			name: "requires a positive limit",
			input: ClaimPredictionActionsInput{
				Limit: 0,
				Lease: time.Second,
			},
		},
		{
			name: "requires a positive lease",
			input: ClaimPredictionActionsInput{
				Limit: 1,
				Lease: 0,
			},
		},
		{
			name: "rejects a negative lease",
			input: ClaimPredictionActionsInput{
				Limit: 1,
				Lease: -time.Second,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateClaimPredictionActionsInput(tt.input); err == nil {
				t.Fatal("ValidateClaimPredictionActionsInput() error = nil, want an error")
			}
		})
	}

	if err := ValidateClaimPredictionActionsInput(ClaimPredictionActionsInput{
		Limit: 2,
		Lease: time.Minute,
	}); err != nil {
		t.Fatalf("ValidateClaimPredictionActionsInput() error = %v", err)
	}
}
