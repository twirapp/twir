package hydrator

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/events/internal/shared"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/valyala/fasttemplate"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Db      *gorm.DB
	TwirBus *bus_core.Bus
}

func New(opts Opts) *Hydrator {
	return &Hydrator{
		db:      opts.Db,
		twirBus: opts.TwirBus,
	}
}

type Hydrator struct {
	db      *gorm.DB
	twirBus *bus_core.Bus
}

func (c *Hydrator) HydrateStringWithData(channelId, channelTwitchUserID, channelDBID, str string, data any) (string, error) {
	template := fasttemplate.New(str, "{", "}")

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	m := make(map[string]any)

	if err = json.Unmarshal(bytes, &m); err != nil {
		return "", err
	}

	s := template.ExecuteFuncString(
		func(w io.Writer, tag string) (int, error) {
			splittedTag := strings.Split(tag, ".")
			if len(splittedTag) > 1 {
				val, ok := m[splittedTag[0]].(map[string]any)
				if !ok {
					// key not found in map
					// return 0, fmt.Errorf("key '%s' is not a map[string]interface{}", splittedTag[0])
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
					// return 0, fmt.Errorf("key '%s' is not found", tag)
					return w.Write([]byte(""))
				}

				return w.Write([]byte(fmt.Sprint(val)))
			}
		},
	)

	var userId, userLogin, userName string
	if m["userId"] != nil {
		userId = m["userId"].(string)
	}

	if m["userName"] != nil {
		userName = m["userName"].(string)
	}

	if m["userDisplayName"] != nil {
		userLogin = m["userDisplayName"].(string)
	}

	parsedChannelID, err := uuid.Parse(channelDBID)
	if err != nil {
		return "", fmt.Errorf("parse channel id: %w", err)
	}

	var platformSource *platform.Platform
	if eventData, ok := data.(shared.EventData); ok && eventData.Platform != "" {
		source := eventData.Platform
		platformSource = &source
	}

	resp, err := c.twirBus.Parser.ParseVariablesInText.Request(
		context.Background(), parser.ParseVariablesInTextRequest{
			ChannelID:      parsedChannelID,
			Text:           s,
			UserID:         userId,
			UserLogin:      userName,
			UserName:       userLogin,
			PlatformSource: platformSource,
		},
	)
	if err != nil {
		return "", err
	}

	if resp.Data.Text != "" {
		s = resp.Data.Text
	}

	return s, nil
}
