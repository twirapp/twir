package chat_client

func (c *ChatClient) Leave(channel string) {
	clientReader, ok := c.channelsToReaders[channel]
	if ok {
		clientReader.size--
		if clientReader.size == 0 {
			clientReader.disconnectChann <- struct{}{}
			delete(c.channelsToReaders, channel)
		}
	}

	for _, r := range c.Readers {
		r.Depart(channel)
	}

	c.Writer.Depart(channel)

	delete(c.RateLimiters.Channels.Items, channel)
}

func (c *ChatClient) readerJoin(reader *BotClientIrc, channel string) {
	reader.Join(channel)
	c.channelsToReaders[channel] = reader
	reader.size++
}

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
