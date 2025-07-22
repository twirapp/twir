package emotes_store

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"sync"
	"unsafe"

	"github.com/satont/twir/apps/emotes-cacher/internal/emote"
	"github.com/satont/twir/libs/logger"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
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
		channels: map[ChannelID]map[emotes_cacher.ServiceName]Service{
			GlobalChannelID: {
				emotes_cacher.ServiceNameBTTV:    Service{},
				emotes_cacher.ServiceNameSevenTV: Service{},
				emotes_cacher.ServiceNameFFZ:     Service{},
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

				go func() {
					s.fillGlobal()
				}()

				return nil
			},
		},
	)

	return s
}

type EmotesStore struct {
	channels map[ChannelID]map[emotes_cacher.ServiceName]Service
	mu       sync.RWMutex

	logger logger.Logger
	gorm   *gorm.DB
}

type ChannelID string

type Service map[emote.ID]emote.Emote

func (c *EmotesStore) AddEmotes(
	channelID ChannelID,
	service emotes_cacher.ServiceName,
	emotes ...emote.Emote,
) {
	c.mu.Lock()
	defer func() {
		c.mu.Unlock()
		c.logCount()
	}()

	if _, ok := c.channels[channelID]; !ok {
		c.channels[channelID] = make(map[emotes_cacher.ServiceName]Service)
	}

	if _, ok := c.channels[channelID][service]; !ok {
		c.channels[channelID][service] = Service{}
	}

	for _, emote := range emotes {
		c.channels[channelID][service][emote.ID] = emote
	}
}

func (c *EmotesStore) GetChannelEmotesServices(channelID ChannelID) map[emotes_cacher.ServiceName]Service {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if emotes, ok := c.channels[channelID]; ok {
		return emotes
	}

	return nil
}

func (c *EmotesStore) GetEmotesByService(
	channelID ChannelID,
	service emotes_cacher.ServiceName,
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
	service emotes_cacher.ServiceName,
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
	service emotes_cacher.ServiceName,
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
	mapSize := c.SizeInBytes()

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

func (c *EmotesStore) SizeInBytes() int {
	c.mu.RLock()
	defer c.mu.RUnlock()

	var totalBytes int

	// Level 1: Iterate over channels map
	for channelID, servicesMap := range c.channels {
		// Add size of the ChannelID string key
		totalBytes += len(channelID)

		// Level 2: Iterate over services map
		for serviceName, service := range servicesMap {
			// Add size of the ServiceName string key
			totalBytes += len(string(serviceName)) // Assuming ServiceName is a string-based type

			// Level 3: Iterate over the innermost Service map
			for emoteID, emote := range service {
				// Add size of the emote.ID string key
				totalBytes += len(emoteID)
				// Add size of the emote.Emote struct value
				totalBytes += sizeOfEmote(emote)
			}
		}
	}

	return totalBytes
}

// sizeOfEmote calculates the "deep" size of an emote.Emote struct.
// This is a helper function that needs to be adapted to the actual
// fields of your emote.Emote struct.
func sizeOfEmote(e emote.Emote) int {
	// Start with the "shallow" size of the struct itself.
	size := int(unsafe.Sizeof(e))

	// Use reflection to inspect the struct fields and add the size
	// of any dynamically-sized data (like strings, slices, or pointers).
	val := reflect.ValueOf(e)
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		switch field.Kind() {
		// If a field is a string, add the length of its content.
		case reflect.String:
			size += field.Len()

		// If a field is a slice, you'd add size of underlying array.
		// For example: size += field.Cap() * int(field.Type().Elem().Size())
		case reflect.Slice:
			// This is a simplified example. A full implementation would be more complex.
			size += field.Cap() * int(field.Type().Elem().Size())
		}
	}
	return size
}
