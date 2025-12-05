package channels_modules_obs_websocket

import (
	"context"

	"github.com/twirapp/twir/libs/entities/obs"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (obs.ObsWebsocketData, error)
	Upsert(ctx context.Context, input UpsertInput) (obs.ObsWebsocketData, error)
	Delete(ctx context.Context, id int) error
}

type UpdateInput struct {
	ServerPort     *int
	ServerAddress  *string
	ServerPassword *string
}

type UpsertInput struct {
	ChannelID      string
	ServerPort     *int
	ServerAddress  *string
	ServerPassword *string
	Scenes         *[]string
	Sources        *[]string
	AudioSources   *[]string
}
