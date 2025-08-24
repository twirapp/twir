package channelsintegrationsvalorant

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/channels_integrations_valorant/model"
)

type Repository interface {
	GetByChannelID(ctx context.Context, channelID string) (model.ChannelIntegrationValorant, error)
}
