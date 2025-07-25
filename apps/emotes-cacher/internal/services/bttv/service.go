package bttv

import (
	"context"
	"fmt"
	"log/slog"
	"strings"

	"github.com/goccy/go-json"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/emote"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/emotes_store"
	"github.com/twirapp/twir/apps/emotes-cacher/internal/socket_client"
	"github.com/twirapp/twir/libs/logger"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
	"github.com/twirapp/twir/libs/repositories/channels/model"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Gorm       *gorm.DB
	Logger     logger.Logger
	EmoteStore *emotes_store.EmotesStore
}

type Service struct {
	gorm               *gorm.DB
	logger             logger.Logger
	registeredChannels map[string]struct{}
	socket             *socket_client.WsConnection
	emotesStore        *emotes_store.EmotesStore
}

func New(opts Opts) error {
	c := Service{
		gorm:               opts.Gorm,
		logger:             opts.Logger,
		registeredChannels: make(map[string]struct{}),
		emotesStore:        opts.EmoteStore,
	}

	ctx, cancel := context.WithCancel(context.Background())

	opts.LC.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				socket, err := socket_client.New(
					ctx,
					socket_client.Opts{
						OnMessage: c.onMessage,
						OnReconnect: func(_ context.Context, _ *socket_client.WsConnection) {
							c.logger.Info("Reconnected to BTTV websocket")
							if err := c.start(ctx); err != nil {
								c.logger.Error("failed to re-init bttv connections", "error", err)
							}
						},
						OnConnect: func(_ context.Context, _ *socket_client.WsConnection) {
							c.logger.Info("Connected to BTTV websocket")
						},
						Url:                "wss://sockets.betterttv.net/ws",
						SubscriptionsLimit: 0,
					},
				)
				if err != nil {
					cancel()
					return fmt.Errorf("failed to create socket client: %w", err)
				}

				c.socket = socket

				go func() {
					if err := c.start(ctx); err != nil {
						cancel()
						c.logger.Error("failed to start bttv service", "error", err)
					}
				}()

				return nil
			},
			OnStop: func(ctx context.Context) error {
				cancel()
				return nil
			},
		},
	)

	return nil
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

func (c *Service) AddChannels(ctx context.Context, channelsIDs ...string) error {
	if len(channelsIDs) == 0 {
		return nil
	}

	for _, channelID := range channelsIDs {
		msg := map[string]any{
			"name": "join_channel",
			"data": map[string]any{
				"name": fmt.Sprintf("twitch:%s", channelID),
			},
		}

		if err := c.socket.Subscribe(ctx, msg); err != nil {
			fmt.Printf("Failed to subscribe: %v %v\n", err, msg)
			continue
		}

		c.registeredChannels[channelID] = struct{}{}
	}

	return nil
}

type baseMessage struct {
	Name string `json:"name"`
}

type addMessage struct {
	Name string `json:"name"`
	Data struct {
		Emote struct {
			Id        string `json:"id"`
			Code      string `json:"code"`
			ImageType string `json:"imageType"`
			Animated  bool   `json:"animated"`
			User      struct {
				Id          string `json:"id"`
				Name        string `json:"name"`
				DisplayName string `json:"displayName"`
				ProviderId  string `json:"providerId"`
			} `json:"user"`
		} `json:"emote"`
		Channel string `json:"channel"`
	} `json:"data"`
}

type updateMessage struct {
	Name string `json:"name"`
	Data struct {
		Emote struct {
			Id   string `json:"id"`
			Code string `json:"code"`
		} `json:"emote"`
		Channel string `json:"channel"`
	} `json:"data"`
}

type deleteMessage struct {
	Name string `json:"name"`
	Data struct {
		EmoteId string `json:"emoteId"`
		Channel string `json:"channel"`
	} `json:"data"`
}

func (c *Service) onMessage(ctx context.Context, client *socket_client.WsConnection, data []byte) {
	var baseMsg baseMessage
	if err := json.Unmarshal(data, &baseMsg); err != nil {
		c.logger.Error(
			"Failed to unmarshal message",
			slog.Any("error", err),
			slog.String("data", string(data)),
		)
		return
	}

	switch baseMsg.Name {
	case "emote_create":
		var msg addMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			c.logger.Error(
				"Failed to unmarshal emote_create message",
				slog.Any("error", err),
				slog.String("data", string(data)),
			)
			return
		}

		channel := strings.Split(msg.Data.Channel, ":")
		if len(channel) < 2 {
			c.logger.Error(
				"Invalid channel format in emote_create message",
				slog.String("channel", msg.Data.Channel),
			)
			return
		}

		c.emotesStore.AddEmotes(
			emotes_store.ChannelID(channel[1]),
			emotes_cacher.ServiceNameBTTV,
			emote.Emote{
				ID:   emote.ID(msg.Data.Emote.Id),
				Name: msg.Data.Emote.Code,
			},
		)
	case "emote_update":
		var msg updateMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			c.logger.Error(
				"Failed to unmarshal emote_update message",
				slog.Any("error", err),
				slog.String("data", string(data)),
			)
			return
		}

		channel := strings.Split(msg.Data.Channel, ":")
		if len(channel) < 2 {
			c.logger.Error(
				"Invalid channel format in emote_update message",
				slog.String("channel", msg.Data.Channel),
			)
			return
		}

		c.emotesStore.Update(
			emotes_store.ChannelID(channel[1]),
			emotes_cacher.ServiceNameBTTV,
			emote.ID(msg.Data.Emote.Id),
			emote.Emote{
				ID:   emote.ID(msg.Data.Emote.Id),
				Name: msg.Data.Emote.Code,
			},
		)
	case "emote_delete":
		var msg deleteMessage
		if err := json.Unmarshal(data, &msg); err != nil {
			c.logger.Error(
				"Failed to unmarshal emote_delete message",
				slog.Any("error", err),
				slog.String("data", string(data)),
			)
			return
		}

		channel := strings.Split(msg.Data.Channel, ":")
		if len(channel) < 2 {
			c.logger.Error(
				"Invalid channel format in emote_delete message",
				slog.String("channel", msg.Data.Channel),
			)
			return
		}

		c.emotesStore.RemoveEmoteById(
			emotes_store.ChannelID(channel[1]),
			emotes_cacher.ServiceNameBTTV,
			emote.ID(msg.Data.EmoteId),
		)
	}
}
