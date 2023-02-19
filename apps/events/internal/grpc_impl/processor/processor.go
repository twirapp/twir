package processor

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/events/internal"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/websockets"
	"github.com/valyala/fasttemplate"
	"io"
	"strings"
)

var InternalError = errors.New("internal")

func hydrateStringWithData(str string, data any) (string, error) {
	template := fasttemplate.New(str, "{", "}")

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	m := make(map[string]any)

	if err = json.Unmarshal(bytes, &m); err != nil {
		return "", err
	}

	s := template.ExecuteFuncString(func(w io.Writer, tag string) (int, error) {
		splittedTag := strings.Split(tag, ".")
		if len(splittedTag) > 1 {
			val, ok := m[splittedTag[0]].(map[string]any)
			if !ok {
				// key not found in map
				//return 0, fmt.Errorf("key '%s' is not a map[string]interface{}", splittedTag[0])
				return w.Write([]byte(""))
			}

			v, ok := val[splittedTag[1]]
			if !ok {
				// key not found in map
				// return 0, fmt.Errorf("key '%s' is not found in map", splittedTag[1])
				return w.Write([]byte(""))
			}

			return w.Write([]byte(fmt.Sprint(v)))
		} else {
			val, ok := m[tag]
			if !ok {
				// not a found
				//return 0, fmt.Errorf("key '%s' is not found", tag)
				return w.Write([]byte(""))
			}

			return w.Write([]byte(fmt.Sprint(val)))
		}
	})

	return s, nil
}

type ProcessorCache struct {
	channelModerators []helix.Moderator
	channelVips       []helix.ChannelVips
	dbChannel         *model.Channels
}

type Processor struct {
	services          *internal.Services
	streamerApiClient *helix.Client

	websocketsGrpc websockets.WebsocketClient

	data      *internal.Data
	channelId string

	cache ProcessorCache
}

type Opts struct {
	Services          *internal.Services
	StreamerApiClient *helix.Client
	Data              *internal.Data
	ChannelID         string
}

func NewProcessor(opts Opts) *Processor {
	return &Processor{
		services:          opts.Services,
		streamerApiClient: opts.StreamerApiClient,
		data:              opts.Data,
		channelId:         opts.ChannelID,
		cache:             ProcessorCache{},
	}
}

func (c *Processor) getDbChannel() (*model.Channels, error) {
	if c.cache.dbChannel != nil {
		return c.cache.dbChannel, nil
	}

	channel := &model.Channels{}
	err := c.services.DB.Where(`"id" = ?`, c.channelId).Find(channel).Error

	if err != nil {
		return nil, err
	}

	if channel.ID == "" {
		return nil, errors.New("channel not found")
	}

	c.cache.dbChannel = channel

	return channel, nil
}
