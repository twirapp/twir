package seventv

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/avast/retry-go/v4"
	"github.com/kr/pretty"
	dispatchtypes "github.com/satont/twir/apps/emotes-cacher/internal/services/seventv/dispatch_types"
	"github.com/satont/twir/apps/emotes-cacher/internal/services/seventv/messages"
	"github.com/satont/twir/apps/emotes-cacher/internal/services/seventv/operations"
	config "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/libs/integrations/seventv"
	"github.com/twirapp/twir/libs/repositories/channels/model"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Gorm   *gorm.DB
	Config config.Config
}

func New(opts Opts) error {
	s := Service{
		sockets:          nil,
		gorm:             opts.Gorm,
		sevenTvApiClient: seventv.NewClient(opts.Config.SevenTvToken),
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

type Service struct {
	sockets []*wsConnection

	gorm             *gorm.DB
	sevenTvApiClient seventv.Client
}

func (c *Service) start(ctx context.Context) error {
	var channels []model.Channel
	if err := c.gorm.Select("id").Where(`"isEnabled" = true`).Find(&channels).Error; err != nil {
		return err
	}

	channelsWithEmoteSets := map[string]string{}

	var wg sync.WaitGroup
	var my sync.Mutex

	for _, channel := range channels {
		wg.Add(1)

		go func() {
			defer wg.Done()
			retry.Do(
				func() error {
					profile, err := c.sevenTvApiClient.GetProfileByTwitchId(ctx, channel.ID)
					if err != nil {
						return fmt.Errorf(
							"failed to fetch profile for channel %s: %w",
							channel.ID,
							err,
						)
					}
					if profile == nil || profile.Users.UserByConnection == nil || profile.Users.UserByConnection.Style.ActiveEmoteSet == nil {
						fmt.Println("No active emote set found for channel:", channel.ID)
						return nil
					}

					my.Lock()
					channelsWithEmoteSets[channel.ID] = profile.Users.UserByConnection.Style.ActiveEmoteSet.Id
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

	messages := make([]map[string]interface{}, 0, len(channelsWithEmoteSets))
	for channelId, currentEmoteSetId := range channelsWithEmoteSets {
		messages = append(
			messages, map[string]interface{}{
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

		if currentEmoteSetId != "" {
			messages = append(
				messages, map[string]interface{}{
					"op": operations.OutgoingOpSubscribe,
					"d": map[string]any{
						"type": dispatchtypes.UpdateEmoteSet,
						"condition": map[string]string{
							"object_id": currentEmoteSetId,
						},
					},
				},
			)
		}
	}

	for _, msg := range messages {
		var conn *wsConnection
		for _, ws := range c.sockets {
			if ws.subscriptionsCount < ws.subscriptionsLimit {
				conn = ws
			}
		}
		if conn == nil {
			newConn, err := createConn(ctx, c.onMessage)
			if err != nil {
				return err
			}
			conn = newConn
			c.sockets = append(c.sockets, newConn)
		}

		if err := conn.subscribe(ctx, msg); err != nil {
			fmt.Printf("Failed to subscribe: %v %v\n", err, msg)
		}
	}

	return nil
}

func (c *Service) stop() error {
	for _, socket := range c.sockets {
		if err := socket.Close(); err != nil {
			return err
		}
	}
	c.sockets = nil

	return nil
}

func (c *Service) onMessage(msg []byte) {
	var base messages.BaseMessageWithoutData
	if err := json.Unmarshal(msg, &base); err != nil {
		fmt.Println("Failed to unmarshal base message:", err)
		return
	}

	switch base.Operation {
	case operations.IncomingOpDispatch:
		var baseWithType messages.BaseMessage[messages.Dispatch]
		if err := json.Unmarshal(msg, &baseWithType); err != nil {
			fmt.Println("Failed to unmarshal dispatch message:", err)
			return
		}

		if baseWithType.Data.Type != "emote_set.update" {
			fmt.Println("Received non-emote-set update dispatch:", baseWithType.Data.Type)
			return
		}

		if err := c.handleEmoteSetUpdate(baseWithType.Data); err != nil {
			fmt.Println("Failed to handle emote set update:", err)
		}
	}
}

func (c *Service) handleEmoteSetUpdate(data messages.Dispatch) error {
	pretty.Println(data)
	return nil
}
