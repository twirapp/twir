package seventv

import (
	"context"
	"encoding/json"
	"log/slog"

	"github.com/twirapp/twir/apps/emotes-cacher/internal/emote"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/emotes_store"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/services/seventv/messages"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/services/seventv/operations"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/integrations/seventv"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/channels/model"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Gorm        *gorm.DB
	Config      config.Config
	Logger      logger.Logger
	EmotesStore *emotes_store.EmotesStore
}

func New(opts Opts) error {
	s := Service{
		sockets:                          nil,
		gorm:                             opts.Gorm,
		sevenTvApiClient:                 seventv.NewClient(opts.Config.SevenTvToken),
		logger:                           opts.Logger,
		emotesStore:                      opts.EmotesStore,
		registeredChannelsWithEmoteSetId: make(channelsWithEmotesSetsIds),
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go s.start(context.Background())
				return nil
			},
			OnStop: func(ctx context.Context) error {
				return s.stop()
			},
		},
	)

	return nil
}

type socketInstance struct {
	Instance  *conn
	ShardID   uint8
	SessionID string
}

type Service struct {
	sockets []*socketInstance

	gorm                             *gorm.DB
	sevenTvApiClient                 seventv.Client
	registeredChannelsWithEmoteSetId channelsWithEmotesSetsIds
	logger                           logger.Logger
	emotesStore                      *emotes_store.EmotesStore
}

func (c *Service) start(ctx context.Context) error {
	var channels []string
	if err := c.gorm.
		WithContext(ctx).
		Model(&model.Channel{}).
		Select("id").
		Where(`"isEnabled" = true`).
		Find(&channels).Error; err != nil {
		return err
	}

	return c.AddChannels(ctx, channels...)
}

func (c *Service) stop() error {
	for _, socket := range c.sockets {
		if err := socket.Instance.Stop(); err != nil {
			return err
		}
	}
	c.sockets = nil

	return nil
}

func (c *Service) onMessage(
	ctx context.Context,
	client *conn,
	msg []byte,
) {
	var base messages.BaseMessageWithoutData
	if err := json.Unmarshal(msg, &base); err != nil {
		c.logger.Error("Failed to unmarshal base message", slog.Any("error", err))
		return
	}

	switch base.Operation {
	case operations.IncomingOpHello:
		var helloMsg messages.BaseMessage[messages.HelloMessage]
		if err := json.Unmarshal(msg, &helloMsg); err != nil {
			c.logger.Error("Failed to unmarshal hello message", slog.Any("error", err))
			return
		}

		client.maxCapacity = int(helloMsg.Data.SubscriptionLimit)
		for _, socket := range c.sockets {
			if socket.Instance == client {
				socket.SessionID = helloMsg.Data.SessionID
				break
			}
		}
		c.logger.Info("Connected to 7TV websocket", slog.String("session_id", helloMsg.Data.SessionID))
	case operations.IncomingOpDispatch:
		var baseWithType messages.BaseMessage[messages.Dispatch]
		if err := json.Unmarshal(msg, &baseWithType); err != nil {
			c.logger.Error("Failed to unmarshal dispatch message", slog.Any("error", err))
			return
		}

		if baseWithType.Data.Type != "emote_set.update" {
			return
		}

		if err := c.handleEmoteSetUpdate(ctx, baseWithType.Data); err != nil {
			c.logger.Error("Failed to handle emote set update", slog.Any("error", err))
		}
	}
}

func (c *Service) handleEmoteSetUpdate(_ context.Context, data messages.Dispatch) error {
	for _, emotes := range data.Body.Pushed {
		if emotes.Value == nil || emotes.Value.Data == nil {
			continue
		}

		var channelID string
		for key, value := range c.registeredChannelsWithEmoteSetId {
			if value == data.Body.Id {
				channelID = key
				break
			}
		}

		if channelID == "" {

		}

		emoteData := emotes.Value.Data

		c.logger.Info(
			"Received new emote",
			slog.String("id", emoteData.Id),
			slog.String("name", emoteData.Name),
			slog.String("channel_id", channelID),
			slog.String("emote_set", data.Body.Id),
		)

		c.emotesStore.AddEmotes(
			emotes_store.ChannelID(channelID),
			emotes_cacher.ServiceNameSevenTV,
			emote.Emote{
				ID:   emote.ID(emoteData.Id),
				Name: emoteData.Name,
			},
		)
	}

	for _, emotes := range data.Body.Updated {
		if emotes.Value == nil || emotes.Value.Data == nil || emotes.OldValue == nil || emotes.OldValue.Data == nil {
			continue
		}

		var channelID string
		for key, value := range c.registeredChannelsWithEmoteSetId {
			if value == data.Body.Id {
				channelID = key
				break
			}
		}

		if channelID == "" {
			continue
		}

		emoteData := emotes.Value.Data

		c.logger.Info(
			"Received updated emote",
			slog.String("id", emoteData.Id),
			slog.String("name", emoteData.Name),
			slog.String("channel_id", channelID),
			slog.String("emote_set", data.Body.Id),
		)

		c.emotesStore.Update(
			emotes_store.ChannelID(channelID),
			emotes_cacher.ServiceNameSevenTV,
			emote.ID(emoteData.Id),
			emote.Emote{
				ID:   emote.ID(emoteData.Id),
				Name: emoteData.Name,
			},
		)
	}

	for _, emotes := range data.Body.Pulled {
		if emotes.Value != nil || emotes.OldValue == nil {
			continue
		}

		var channelID string
		for key, value := range c.registeredChannelsWithEmoteSetId {
			if value == data.Body.Id {
				channelID = key
				break
			}
		}

		if channelID == "" {
			continue
		}

		emoteData := emotes.OldValue.Data

		c.logger.Info(
			"Received deleted emote",
			slog.String("id", emoteData.Id),
			slog.String("name", emoteData.Name),
			slog.String("channel_id", channelID),
			slog.String("emote_set", data.Body.Id),
		)

		c.emotesStore.RemoveEmoteById(
			emotes_store.ChannelID(channelID),
			emotes_cacher.ServiceNameSevenTV,
			emote.ID(emoteData.Id),
		)
	}

	return nil
}
