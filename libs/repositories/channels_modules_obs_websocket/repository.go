package channels_modules_obs_websocket

import (
	"context"

	"github.com/twirapp/twir/libs/entities/obs"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (obs.ObsWebsocketData, error)
	Create(ctx context.Context, input CreateInput) (obs.ObsWebsocketData, error)
	Update(ctx context.Context, id int, input UpdateInput) error
	UpdateSources(ctx context.Context, channelID string, input UpdateSourcesInput) error
	Delete(ctx context.Context, id int) error
}

type UpdateInput struct {
	ServerPort     *int
	ServerAddress  *string
	ServerPassword *string
}

type CreateInput struct {
	ChannelID      string
	ServerPort     int
	ServerAddress  string
	ServerPassword string
}

type UpdateSourcesInput struct {
	Scenes       []string
	Sources      []string
	AudioSources []string
}
