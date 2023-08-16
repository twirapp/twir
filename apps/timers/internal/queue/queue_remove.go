package queue

import (
	"errors"
	"log/slog"
)

func (c *Queue) Remove(id string) error {
	t, ok := c.timers[id]

	if !ok {
		return errors.New("timer not found in queue")
	}

	// stop ticker
	t.doneChann <- true

	delete(c.timers, id)

	c.logger.Info(
		"Removed timer",
		slog.String("id", t.ID),
		slog.String("channelId", t.ChannelID),
		t.Interval,
	)

	return nil
}
