package seventv

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"

	"github.com/avast/retry-go/v4"
	"github.com/satont/twir/apps/emotes-cacher/internal/emote"
	"github.com/satont/twir/apps/emotes-cacher/internal/emotes_store"
	dispatchtypes "github.com/satont/twir/apps/emotes-cacher/internal/services/seventv/dispatch_types"
	"github.com/satont/twir/apps/emotes-cacher/internal/services/seventv/messages"
	"github.com/satont/twir/apps/emotes-cacher/internal/services/seventv/operations"
	"github.com/satont/twir/apps/emotes-cacher/internal/socket_client"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	emotes_cacher "github.com/twirapp/twir/libs/bus-core/emotes-cacher"
	"github.com/twirapp/twir/libs/integrations/seventv"
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
		registeredChannelsWithEmoteSetId: make(map[string]string),
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

type ConnData struct {
	SessionId string
}

type socketInstance struct {
	SessionID string
	Instance  *socket_client.WsConnection
}

type Service struct {
	sockets []*socketInstance

	gorm                             *gorm.DB
	sevenTvApiClient                 seventv.Client
	registeredChannelsWithEmoteSetId map[string]string
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

func (c *Service) AddChannels(ctx context.Context, channelsIDs ...string) error {
	if len(channelsIDs) == 0 {
		return nil
	}

	channelsWithEmoteSets := map[string]string{}

	var wg sync.WaitGroup
	var my sync.Mutex

	for _, channel := range channelsIDs {
		if _, ok := c.registeredChannelsWithEmoteSetId[channel]; ok {
			continue
		}

		wg.Add(1)

		go func() {
			defer wg.Done()
			retry.Do(
				func() error {
					profile, err := c.sevenTvApiClient.GetProfileByTwitchId(ctx, channel)
					if err != nil {
						return fmt.Errorf(
							"failed to fetch profile for channel %s: %w",
							channel,
							err,
						)
					}
					if profile == nil || profile.Users.UserByConnection == nil || profile.Users.UserByConnection.Style.ActiveEmoteSet == nil {
						fmt.Println("No active emote set found for channel:", channel)
						return nil
					}

					my.Lock()
					channelsWithEmoteSets[channel] = profile.Users.UserByConnection.Style.ActiveEmoteSet.Id
					my.Unlock()

					return nil
				},
				retry.RetryIf(
					func(err error) bool {
						return err != nil
					},
				),
				retry.Attempts(10),
			)
		}()
	}

	wg.Wait()

	subMessages := make([]map[string]interface{}, 0, len(channelsIDs))
	for _, channelId := range channelsIDs {
		subMessages = append(
			subMessages, map[string]interface{}{
				"op": operations.OutgoingOpSubscribe,
				"d": map[string]any{
					"type": "emote_set.*",
					"condition": map[string]string{
						"ctx":      "channel",
						"id":       channelId,
						"platform": "TWITCH",
					},
				},
			},
		)

		if emoteSetId, ok := channelsWithEmoteSets[channelId]; ok {
			c.registeredChannelsWithEmoteSetId[channelId] = emoteSetId

			subMessages = append(
				subMessages, map[string]interface{}{
					"op": operations.OutgoingOpSubscribe,
					"d": map[string]any{
						"type": dispatchtypes.UpdateEmoteSet,
						"condition": map[string]string{
							"object_id": emoteSetId,
						},
					},
				},
			)
		}
	}

	for _, msg := range subMessages {
		var instance *socketInstance
		for _, ws := range c.sockets {
			if ws.Instance.SubscriptionsCount < ws.Instance.SubscriptionsLimit {
				instance = ws
			}
		}
		if instance == nil {
			newConn, err := socket_client.New(
				ctx,
				socket_client.Opts{
					OnMessage:          c.onMessage,
					OnReconnect:        nil,
					OnConnect:          nil,
					Url:                "wss://events.7tv.io/v3",
					SubscriptionsLimit: 0,
				},
			)
			if err != nil {
				return err
			}
			instance = &socketInstance{
				Instance: newConn,
			}

			c.sockets = append(
				c.sockets,
				instance,
			)
		}

		for {
			if instance.SessionID == "" {
				continue
			}
			if err := instance.Instance.Subscribe(ctx, msg); err != nil {
				fmt.Printf("Failed to subscribe: %v %v\n", err, msg)
				continue
			}

			break
		}
	}

	return nil
}

func (c *Service) stop() error {
	for _, socket := range c.sockets {
		if err := socket.Instance.Close(); err != nil {
			return err
		}
	}
	c.sockets = nil

	return nil
}

func (c *Service) onMessage(
	ctx context.Context,
	client *socket_client.WsConnection,
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
		client.SubscriptionsLimit = int(helloMsg.Data.SubscriptionLimit)
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
			c.logger.Warn(
				"Received non-emote-set update dispatch",
				slog.String("type", baseWithType.Data.Type),
			)
			return
		}

		if err := c.handleEmoteSetUpdate(ctx, baseWithType.Data); err != nil {
			c.logger.Error("Failed to handle emote set update", slog.Any("error", err))
		}
	}
}

func (c *Service) onReconnect(ctx context.Context, client *socket_client.WsConnection) {
	// Re-subscribe to all channels
	for _, socket := range c.sockets {
		if socket.Instance == client {
			c.logger.Info("Reconnecting to 7TV websocket", slog.String("session_id", socket.SessionID))
			resumeMsg := map[string]interface{}{
				"op": 34,
				"d": map[string]string{
					"session_id": socket.SessionID,
				},
			}

			if err := socket.Instance.SendMessage(ctx, resumeMsg); err != nil {
				c.logger.Error("Failed to send resume message", slog.Any("error", err))
			}
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
			slog.String("channel_id", data.Body.Actor.Id),
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
			slog.String("channel_id", data.Body.Actor.Id),
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
			slog.String("channel_id", data.Body.Actor.Id),
		)

		c.emotesStore.RemoveEmoteById(
			emotes_store.ChannelID(channelID),
			emotes_cacher.ServiceNameSevenTV,
			emote.ID(emoteData.Id),
		)
	}

	return nil
}
