package dota

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/dota/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID uuid.UUID) (model.ChannelDotaSettings, error)
	GetByGsiToken(ctx context.Context, token string) (model.ChannelDotaSettings, error)
	Create(ctx context.Context, input CreateInput) (model.ChannelDotaSettings, error)
	Update(
		ctx context.Context,
		channelID uuid.UUID,
		input UpdateInput,
	) (model.ChannelDotaSettings, error)
	UpdateMatchResult(
		ctx context.Context,
		channelID uuid.UUID,
		won bool,
		mmrDelta int,
	) (model.ChannelDotaSettings, error)
	ApplyMatchResultOnce(
		ctx context.Context,
		input ApplyMatchResultInput,
	) (model.ChannelDotaSettings, error)
	GetMatchState(ctx context.Context, channelID uuid.UUID) (model.MatchState, error)
	ApplyMatchStateTransition(ctx context.Context, input ApplyMatchStateTransitionInput) (bool, error)
	ClaimPredictionActions(
		ctx context.Context,
		input ClaimPredictionActionsInput,
	) ([]model.ClaimedOutboxAction, error)
	CompletePredictionAction(ctx context.Context, actionID uuid.UUID, lockToken uuid.UUID) error
	RetryPredictionAction(
		ctx context.Context,
		actionID uuid.UUID,
		lockToken uuid.UUID,
		availableAt time.Time,
	) error
	ResetSession(ctx context.Context, channelID uuid.UUID) (model.ChannelDotaSettings, error)
	RegenerateGsiToken(
		ctx context.Context,
		channelID uuid.UUID,
	) (model.ChannelDotaSettings, error)
}

type CreateInput struct {
	ChannelID          uuid.UUID
	Enabled            bool
	SteamAccountID     *string
	Mmr                int
	MmrDelta           int
	PredictionSettings model.PredictionSettings
	ChatEvents         model.ChatEvents
	CommandsSettings   *model.CommandsSettings
}

type ApplyMatchResultInput struct {
	ChannelID uuid.UUID
	MatchID   int64
	Won       bool
	MmrDelta  int
}

type ApplyMatchStateTransitionInput struct {
	ChannelID         uuid.UUID
	ExpectedRevision  int64
	ProviderTimestamp int64
	Snapshot          json.RawMessage
	Actions           []model.OutboxActionInput
}

type ClaimPredictionActionsInput struct {
	Limit int
	Lease time.Duration
}

func ValidateMatchResultInput(input ApplyMatchResultInput) error {
	if input.ChannelID == uuid.Nil {
		return errors.New("channel ID is required")
	}
	if input.MatchID <= 0 {
		return errors.New("match ID must be positive")
	}

	return nil
}

func ValidateOutboxActionInput(input model.OutboxActionInput) error {
	if input.ChannelID == uuid.Nil {
		return errors.New("channel ID is required")
	}
	if input.MatchID <= 0 {
		return errors.New("match ID must be positive")
	}
	if !isOutboxAction(input.Action) {
		return errors.New("action is invalid")
	}
	if input.Sequence <= 0 {
		return errors.New("sequence must be positive")
	}
	if err := validateJSON("payload", input.Payload); err != nil {
		return err
	}

	return nil
}

func ValidateApplyMatchStateTransitionInput(input ApplyMatchStateTransitionInput) error {
	if input.ChannelID == uuid.Nil {
		return errors.New("channel ID is required")
	}
	if input.ExpectedRevision < 0 {
		return errors.New("expected revision must not be negative")
	}
	if input.ProviderTimestamp < 0 {
		return errors.New("provider timestamp must not be negative")
	}
	if err := validateJSON("snapshot", input.Snapshot); err != nil {
		return err
	}

	seenActions := make(map[struct {
		matchID int64
		action  model.OutboxAction
	}]struct{})
	for index, action := range input.Actions {
		if action.ChannelID != input.ChannelID {
			return fmt.Errorf("action %d channel ID must match transition channel ID", index)
		}
		if err := ValidateOutboxActionInput(action); err != nil {
			return fmt.Errorf("action %d: %w", index, err)
		}

		key := struct {
			matchID int64
			action  model.OutboxAction
		}{
			matchID: action.MatchID,
			action:  action.Action,
		}
		if _, exists := seenActions[key]; exists {
			return fmt.Errorf("duplicate action for match %d and action %q", action.MatchID, action.Action)
		}
		seenActions[key] = struct{}{}
	}

	return nil
}

func ValidateClaimPredictionActionsInput(input ClaimPredictionActionsInput) error {
	if input.Limit <= 0 {
		return errors.New("claim limit must be positive")
	}
	if input.Lease <= 0 {
		return errors.New("claim lease must be positive")
	}

	return nil
}

func validateJSON(name string, value json.RawMessage) error {
	if len(bytes.TrimSpace(value)) == 0 {
		return fmt.Errorf("%s is required", name)
	}
	if !json.Valid(value) {
		return fmt.Errorf("%s must be valid JSON", name)
	}

	return nil
}

func isOutboxAction(action model.OutboxAction) bool {
	switch action {
	case model.OutboxActionCreate, model.OutboxActionResolve, model.OutboxActionCancel:
		return true
	default:
		return false
	}
}

func CommandSettingsOrDefault(settings *model.CommandsSettings) model.CommandsSettings {
	if settings != nil {
		return *settings
	}

	return model.CommandsSettings{
		Mmr: true,
		Wl:  true,
		Lg:  true,
		Gm:  true,
		Np:  true,
		Wp:  true,
	}
}

type UpdateInput struct {
	Enabled            bool
	SteamAccountID     *string
	Mmr                int
	MmrDelta           int
	PredictionSettings model.PredictionSettings
	ChatEvents         model.ChatEvents
	CommandsSettings   model.CommandsSettings
}

var (
	ErrNotFound                      = fmt.Errorf("dota settings not found")
	ErrPredictionActionOwnershipLost = errors.New("dota prediction action ownership lost")
)
