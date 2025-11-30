package channelsintegrationsvalorant

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/channels_integrations_valorant/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.ChannelIntegrationValorant, error)
	Create(ctx context.Context, input CreateInput) (model.ChannelIntegrationValorant, error)
	Update(ctx context.Context, id int, input UpdateInput) error
	Delete(ctx context.Context, id int) error
}

type CreateInput struct {
	ChannelID            string
	Enabled              bool
	AccessToken          *string
	RefreshToken         *string
	UserName             *string
	ValorantActiveRegion *string
	ValorantPuuid        *string
}

type UpdateInput struct {
	Enabled              *bool
	AccessToken          *string
	RefreshToken         *string
	UserName             *string
	ValorantActiveRegion *string
	ValorantPuuid        *string
}
