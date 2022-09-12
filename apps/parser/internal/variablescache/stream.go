package variables_cache

import (
	"context"
	"encoding/json"
	"tsuwari/parser/internal/variables/stream"

	"github.com/nicklaw5/helix"
)

func (c *VariablesCacheService) GetChannelStream() *stream.HelixStream {
	c.locks.stream.Lock()
	defer c.locks.stream.Unlock()

	if c.cache.Stream != nil {
		return c.cache.Stream
	}

	rCtx := context.TODO()
	rKey := "streams:" + c.ChannelId
	cachedStream, _ := c.Services.Redis.Get(rCtx, rKey).Result()

	if cachedStream != "" {
		json.Unmarshal([]byte(cachedStream), &c.cache.Stream)
		return c.cache.Stream
	}

	streams, err := c.Services.Twitch.Client.GetStreams(&helix.StreamsParams{
		UserIDs: []string{c.ChannelId},
	})

	if err != nil || len(streams.Data.Streams) == 0 {
		return nil
	}

	stream := stream.HelixStream{
		Stream:   streams.Data.Streams[0],
		Messages: 0,
	}

	rData, err := json.Marshal(stream)
	if err == nil {
		c.Services.Redis.Set(rCtx, rKey, rData, 0)
	}

	c.cache.Stream = &stream
	return c.cache.Stream
}
