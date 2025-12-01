package channelsintegrationslastfm

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/channels_integrations_lastfm/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.ChannelIntegrationLastfm, error)
	Create(ctx context.Context, input CreateInput) (model.ChannelIntegrationLastfm, error)
	Update(ctx context.Context, id int, input UpdateInput) error
	Delete(ctx context.Context, id int) error
}

type CreateInput struct {
	ChannelID  string
	Enabled    bool
	SessionKey *string
	UserName   *string
	Avatar     *string
}

type UpdateInput struct {
	Enabled    *bool
	SessionKey *string
	UserName   *string
	Avatar     *string
}
