package timers

import (
	"context"

	timersentity "github.com/twirapp/twir/libs/entities/timers"
)

func (c *Service) GetAllByChannelID(ctx context.Context, channelID string) ([]timersentity.Timer, error) {
	timers, err := c.timersRepository.GetAllByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	return timers, nil
}
