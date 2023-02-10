package processor

import (
	"encoding/json"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/events/internal"
	"github.com/valyala/fasttemplate"
)

func hydrateStringWithData(str string, data internal.Data) (string, error) {
	template := fasttemplate.New(str, "{", "}")

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	m := make(map[string]any)

	if err = json.Unmarshal(bytes, &m); err != nil {
		return "", err
	}

	s := template.ExecuteString(m)

	return s, nil
}

type Processor struct {
	services          *internal.Services
	streamerApiClient *helix.Client

	data      internal.Data
	channelId string
}

type Opts struct {
	Services          *internal.Services
	StreamerApiClient *helix.Client
	Data              internal.Data
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
