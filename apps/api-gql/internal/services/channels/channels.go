package channels

import (
	"context"
	"errors"
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/repositories/channels/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	ChannelsRepository channels.Repository
}

func New(opts Opts) *Service {
	return &Service{
		channelsRepository: opts.ChannelsRepository,
	}
}

type Service struct {
	channelsRepository channels.Repository
}

var ErrNotFound = fmt.Errorf("channel not found")

func (c *Service) mapToEntity(m model.Channel) entity.Channel {
	return entity.Channel{
		ID:             m.ID,
		IsEnabled:      m.IsEnabled,
		IsTwitchBanned: m.IsTwitchBanned,
		IsBotMod:       m.IsBotMod,
		BotID:          m.BotID,
	}
}

func (c *Service) GetByID(ctx context.Context, channelID string) (entity.Channel, error) {
	channel, err := c.channelsRepository.GetByID(ctx, channelID)
	if err != nil {
		if errors.Is(err, channels.ErrNotFound) {
			return entity.ChannelNil, ErrNotFound
		}

		return entity.ChannelNil, err
	}
	
	return c.mapToEntity(channel), nil
}
