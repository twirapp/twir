package dota

import (
	"context"
	"errors"
	"fmt"

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

func ValidateMatchResultInput(input ApplyMatchResultInput) error {
	if input.ChannelID == uuid.Nil {
		return errors.New("channel ID is required")
	}
	if input.MatchID <= 0 {
		return errors.New("match ID must be positive")
	}

	return nil
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

var ErrNotFound = fmt.Errorf("dota settings not found")
