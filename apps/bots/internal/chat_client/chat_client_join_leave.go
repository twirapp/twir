package chat_client

func (c *ChatClient) Leave(channel string) {
	for _, r := range c.Readers {
		r.Depart(channel)
	}

	reader, ok := c.channelsToReader[channel]
	if ok {
		reader.size--
		if reader.size == 0 {
			reader.disconnectChann <- struct{}{}
			delete(c.channelsToReader, channel)
		}
	}

	c.Writer.Depart(channel)

	delete(c.RateLimiters.Channels.Items, channel)
	c.services.Logger.Info("Leaving channel: " + channel)
}

func (c *ChatClient) readerJoin(reader *BotClientIrc, channel string) {
	reader.Join(channel)
	c.channelsToReader[channel] = reader
	reader.size++
}

const readerCapacity = 50

func (c *ChatClient) Join(channel string) {
	c.joinMu.Lock()
	defer c.joinMu.Unlock()

	c.Writer.Join(channel)

	var ok bool

	for _, reader := range c.Readers {
		if reader.size < readerCapacity {
			ok = true
			c.readerJoin(reader, channel)
		}
	}

	if ok {
		return
	}

	reader := c.createReader()
	c.readerJoin(reader, channel)
}
