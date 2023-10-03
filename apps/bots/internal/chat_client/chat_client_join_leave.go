package chat_client

func (c *ChatClient) Leave(channel string) {
	for _, r := range c.Readers {
		r.Depart(channel)
		r.size--

		if r.size == 0 {
			r.disconnectChann <- struct{}{}
		}
	}

	c.Writer.Depart(channel)

	delete(c.RateLimiters.Channels.Items, channel)
	c.services.Logger.Info("Leave channel: " + channel)
}

func (c *ChatClient) readerJoin(reader *BotClientIrc, channel string) {
	reader.Join(channel)
	reader.size++
}

const readerCapacity = 50

func (c *ChatClient) Join(channel string) {
	c.joinMu.Lock()
	defer c.joinMu.Unlock()

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
