package emotes_store

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"unsafe"

	"github.com/satont/twir/apps/emotes-cacher/internal/emote"
	"github.com/satont/twir/apps/emotes-cacher/internal/services"
	"github.com/satont/twir/libs/logger"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Gorm   *gorm.DB
	Logger logger.Logger
}

const GlobalChannelID = "global"

func New(opts Opts) *EmotesStore {
	s := &EmotesStore{
		channels: map[ChannelID]map[services.ServiceName]Service{
			GlobalChannelID: {
				services.ServiceBttv:    Service{},
				services.ServiceSevenTV: Service{},
				services.ServiceFFZ:     Service{},
			},
		},
		logger: opts.Logger,
		gorm:   opts.Gorm,
		mu:     sync.RWMutex{},
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go func() {
					s.fillChannels()
				}()

				return nil
			},
		},
	)

	return s
}

type EmotesStore struct {
	channels map[ChannelID]map[services.ServiceName]Service
	mu       sync.RWMutex

	logger logger.Logger
	gorm   *gorm.DB
}

type ChannelID string

type Service map[emote.ID]emote.Emote

func (c *EmotesStore) AddEmotes(
	channelID ChannelID,
	service services.ServiceName,
	emotes ...emote.Emote,
) {
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		c.logCount()
	}()

	if _, ok := c.channels[channelID]; !ok {
		c.channels[channelID] = make(map[services.ServiceName]Service)
	}

	if _, ok := c.channels[channelID][service]; !ok {
		c.channels[channelID][service] = Service{}
	}

	for _, emote := range emotes {
		c.channels[channelID][service][emote.ID] = emote
	}
}

func (c *EmotesStore) GetChannelEmotesServices(channelID ChannelID) map[services.ServiceName]Service {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if emotes, ok := c.channels[channelID]; ok {
		return emotes
	}

	return nil
}

func (c *EmotesStore) GetEmotesByService(
	channelID ChannelID,
	service services.ServiceName,
) []emote.Emote {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if _, ok := c.channels[channelID]; !ok {
		return nil
	}

	if _, ok := c.channels[channelID][service]; !ok {
		return nil
	}

	emotes := c.channels[channelID][service]
	result := make([]emote.Emote, 0, len(emotes))
	for _, emote := range emotes {
		result = append(result, emote)
	}

	return nil
}

func (c *EmotesStore) RemoveEmoteById(
	channelID ChannelID,
	service services.ServiceName,
	emoteId emote.ID,
) {
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		c.logCount()
	}()

	if _, ok := c.channels[channelID]; !ok {
		return
	}

	if _, ok := c.channels[channelID][service]; !ok {
		return
	}

	delete(c.channels[channelID][service], emoteId)
}

func (c *EmotesStore) Update(
	channelID ChannelID,
	service services.ServiceName,
	emoteId emote.ID,
	newEmote emote.Emote,
) {
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		c.logCount()
	}()

	if _, ok := c.channels[channelID]; !ok {
		return
	}

	if _, ok := c.channels[channelID][service]; !ok {
		return
	}

	if _, ok := c.channels[channelID][service][emoteId]; !ok {
		return
	}

	c.channels[channelID][service][emoteId] = newEmote
}

func (c *EmotesStore) logCount() {
	mapSize := unsafe.Sizeof(c.channels)

	c.logger.Info(
		"EmotesStore count changed",
		slog.Int("channels", c.GetChannelsCount()),
		slog.Int("emotes", c.GetEmotesCount()),
		slog.String("map_size_bytes", fmt.Sprint(mapSize)),
		slog.String("map_size_mb", fmt.Sprintf("%.2f MB", float64(mapSize)/1024/1024)),
	)
}

func (c *EmotesStore) GetChannelsCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	return len(c.channels)
}

func (c *EmotesStore) GetEmotesCount() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	count := 0
	for _, s := range c.channels {
		for _, emotes := range s {
			count += len(emotes)
		}
	}

	return count
}
