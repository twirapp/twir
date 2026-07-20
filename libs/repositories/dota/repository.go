package dota

import (
	"context"
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
