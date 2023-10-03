package chat_client

import (
	"context"
	"log/slog"
	"time"

	"github.com/samber/lo"
)

func (c *ChatClient) Leave(channel string) {
	for _, r := range c.Readers {
		r.Depart(channel)
	}

	reader, ok := c.channelsToReader.Get(channel)
	if ok {
		reader.size--
		if reader.size == 0 {
			reader.disconnectChann <- struct{}{}
			c.channelsToReader.Delete(channel)
		}
	}

	c.Writer.Depart(channel)

	delete(c.RateLimiters.Channels.Items, channel)
}

func (c *ChatClient) readerJoin(reader *BotClientIrc, channel string) {
	c.joinRateLimiter.Wait(context.TODO())

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))
	defer cancel()

l:
	for {
		select {
		case <-ctx.Done():
			break l
		default:
			if !reader.Connected {
				c.services.Logger.Info(
					"Reader is not connected, waiting...",
					slog.String("channel", channel),
					slog.Int("reader", lo.IndexOf(c.Readers, reader)+1),
				)
				time.Sleep(500 * time.Millisecond)
				continue
			}

			reader.Join(channel)
			c.channelsToReader.Add(channel, reader)
			reader.size++
			break l
		}
	}
}

const readerCapacity int8 = 50

func (c *ChatClient) Join(channel string) {
	c.joinMu.Lock()
	defer c.joinMu.Unlock()

	c.Leave(channel)

	reader, ok := lo.Find(
		c.Readers, func(r *BotClientIrc) bool {
			return r.size < readerCapacity
		},
	)

	if !ok {
		reader = c.createReader()
		c.Readers = append(c.Readers, reader)
	}

	c.readerJoin(reader, channel)
}
