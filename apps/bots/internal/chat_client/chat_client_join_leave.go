package chat_client

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
	reader.Join(channel)
	c.channelsToReader.Add(channel, reader)
	reader.size++
}

const readerCapacity int8 = 1

func (c *ChatClient) Join(channel string) {
	c.joinMu.Lock()
	defer c.joinMu.Unlock()

	// c.Leave(channel)

	c.Writer.Join(channel)

	var reader *BotClientIrc
	for _, r := range c.Readers {
		if r.size >= readerCapacity {
			continue
		}

		reader = r
		break
	}

	if reader == nil {
		reader = c.createReader()
		c.Readers = append(c.Readers, reader)
	}

	c.readerJoin(reader, channel)
}
