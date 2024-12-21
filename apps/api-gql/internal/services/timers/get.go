package timers

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func (c *Service) GetAllByChannelID(ctx context.Context, channelID string) ([]entity.Timer, error) {
	timers, err := c.timersRepository.GetAllByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	converted := make([]entity.Timer, 0, len(timers))
	for _, timer := range timers {
		converted = append(converted, c.dbToModel(timer))
	}

	return converted, nil
}
