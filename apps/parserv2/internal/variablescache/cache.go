package variablescache

import (
	"regexp"
	"sync"
	"tsuwari/parser/internal/helix"
	"tsuwari/parser/pkg/helpers"

	"github.com/go-redis/redis/v9"
)

type VariablesCacheServices struct {
	Redis     *redis.Client
	Regexp    regexp.Regexp
	TwitchApi *helix.Client
}

type VariablesCacheContext struct {
	ChannelId  string
	SenderId   string
	SenderName string
	Text       string
}

type VariablesCacheService struct {
	Context  VariablesCacheContext
	Services VariablesCacheServices
	Cache    VariablesCache
}

type VariablesCache struct {
	Stream *helix.Stream
}

func New(text string, senderId string, channelId string, senderName *string, redis *redis.Client, r regexp.Regexp) *VariablesCacheService {
	cache := &VariablesCacheService{
		Context: VariablesCacheContext{
			ChannelId:  channelId,
			SenderId:   senderId,
			SenderName: *senderName,
			Text:       text,
		},
		Services: VariablesCacheServices{
			Redis:  redis,
			Regexp: r,
		},
		Cache: VariablesCache{
			Stream: nil,
		},
	}

	cache.fillCache()

	return cache
}

func (c *VariablesCacheService) fillCache() {
	matches := c.Services.Regexp.FindAllStringSubmatch(c.Context.Text, len(c.Context.Text))
	myMap := map[string]interface{}{
		"streamId": c.setChannelStream,
	}
	requesting := []string{}
	wg := sync.WaitGroup{}

	for _, match := range matches {
		if match[1] == "" {
			continue
		}

		if helpers.Contains(requesting, match[1]) {
			continue
		}

		if val, ok := myMap[match[1]]; ok {
			wg.Add(1)

			go val.(func(wg *sync.WaitGroup))(&wg)
		}
	}

	wg.Wait()
}

func (c *VariablesCacheService) setChannelStream(wg *sync.WaitGroup) {
	defer wg.Done()

	stream, err := c.Services.TwitchApi.GetStreams(&helix.StreamsParams{
		UserIDs: []string{c.Context.ChannelId},
	})

	if err != nil || stream.Data.Streams == nil {
		return
	}

	c.Cache.Stream = &stream.Data.Streams[0]
}
