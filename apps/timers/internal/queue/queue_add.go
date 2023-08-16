package queue

import (
	"log/slog"
	"time"
)

func (c *Queue) Add(id string) error {
	if err := c.Remove(id); err != nil {
		c.logger.Debug("cannot delete timer", slog.String("id", id), slog.Any("error", err))
	}

	t, err := c.timersRepository.GetById(id)
	if err != nil {
		return err
	}

	var ticker *time.Ticker
	if c.config.AppEnv == "production" {
		ticker = time.NewTicker(time.Duration(t.Interval) * time.Minute)
	} else {
		ticker = time.NewTicker(time.Duration(t.Interval) * time.Second)
	}

	newTimer := timer{
		Timer:        t,
		LastResponse: 0,
		ticker:       ticker,
		doneChann:    make(chan bool),
	}

	c.timers[t.ID] = &newTimer

	go func() {
		for {
			select {
			case <-newTimer.ticker.C:
				c.handle(&newTimer)
			case <-newTimer.doneChann:
				ticker.Stop()
				return
			}
		}
	}()

	c.logger.Info(
		"Added timer",
		slog.String("id", t.ID),
		slog.String("channelId", t.ChannelID),
		slog.Int("interval", t.Interval),
		slog.Int("responses", len(t.Responses)),
	)

	return nil
}
