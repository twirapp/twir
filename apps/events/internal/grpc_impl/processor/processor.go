package processor

import (
	"encoding/json"
	"fmt"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/events/internal"
	"github.com/valyala/fasttemplate"
	"io"
	"strings"
)

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

			return w.Write([]byte(fmt.Sprintf("%v", v)))
		} else {
			val, ok := m[tag].(string)
			if !ok {
				// not a string
				//return 0, fmt.Errorf("key '%s' is not a string", tag)
				return w.Write([]byte(""))
			}

			return w.Write([]byte(val))
		}
	})

	return s, nil
}

type Processor struct {
	services          *internal.Services
	streamerApiClient *helix.Client

	data      *internal.Data
	channelId string
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
	}
}
