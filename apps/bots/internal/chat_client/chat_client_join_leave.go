package chat_client

import (
	"context"
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
	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(10*time.Second))

l:
	for {
		select {
		case <-ctx.Done():
			break
		default:
			if !reader.Connected {
				time.Sleep(50 * time.Millisecond)
				continue
			}

			reader.Join(channel)
			c.channelsToReader.Add(channel, reader)
			reader.size++
			cancel()
			break l
		}
	}
}

const readerCapacity int8 = 50

func (c *ChatClient) Join(channel string) {
	c.joinMu.Lock()
	defer c.joinMu.Unlock()

	c.Leave(channel)

	c.Writer.Join(channel)

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