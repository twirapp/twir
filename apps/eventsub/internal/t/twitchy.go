package t

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/imroc/req/v3"
	"github.com/kr/pretty"
	"github.com/kvizyx/twitchy/eventsub"
	buscore "github.com/twirapp/twir/libs/bus-core"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/repositories/channels/model"
	twitchconduits "github.com/twirapp/twir/libs/repositories/twitch_conduits"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Gorm               *gorm.DB
	TwirBus            *buscore.Bus
	Cfg                cfg.Config
	ConduitsRepository twitchconduits.Repository
}

func New(opts Opts) {
	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				channels := []model.Channel{}
				if err := opts.Gorm.Where(`"isEnabled" = true`).Find(&channels).Error; err != nil {
					return err
				}

				e := eventsub.New()

				socket := e.Websocket()

				go func() {
					socket.OnWelcome(
						func(w eventsub.WebsocketWelcomeMessage) {
							fmt.Println(w.Payload.Session.Id)
							conduitRequest := map[string]any{
								"conduit_id": "7cdbdf24-4e79-4aa7-8b34-768387795d7f",
								"shards": []map[string]any{
									{
										"id": "0",
										"transport": map[string]any{
											"method":     "websocket",
											"session_id": w.Payload.Session.Id,
										},
									},
								},
							}
							conduitBody, err := json.Marshal(&conduitRequest)
							if err != nil {
								panic(err)
							}

							appToken, err := opts.TwirBus.Tokens.RequestAppToken.Request(
								context.Background(),
								struct{}{},
							)
							if err != nil {
								panic(err)
							}

							resp, err := req.R().SetHeader(
								"Client-id",
								opts.Cfg.TwitchClientId,
							).SetBearerAuthToken(appToken.Data.AccessToken).SetBodyJsonBytes(conduitBody).SetContentType("application/json").Patch("https://api.twitch.tv/helix/eventsub/conduits/shards")
							if err != nil {
								panic(err)
							}

							if !resp.IsSuccess() {
								panic(resp.String())
							}

							// for _, c := range channels {
							// 	twitchClient, err := twitch.NewBotClient(c.ID, opts.Cfg, opts.TwirBus)
							// 	if err != nil {
							// 		panic(err)
							// 	}
							//
							// 	requestData := map[string]any{
							// 		"type":    "channel.chat.message",
							// 		"version": "1",
							// 		"condition": map[string]any{
							// 			"broadcaster_user_id": c.ID,
							// 			"user_id":             c.ID,
							// 		},
							// 		"transport": map[string]any{
							// 			"method":     "conduit",
							// 			"conduit_id": conduitRequest["conduit_id"],
							// 		},
							// 	}
							// 	requestBytes, err := json.Marshal(&requestData)
							// 	if err != nil {
							// 		panic(err)
							// 	}
							//
							// 	resp, err := req.R().SetHeader(
							// 		"Client-id",
							// 		opts.Cfg.TwitchClientId,
							// 	).SetBearerAuthToken(appToken.Data.AccessToken).SetBodyJsonBytes(requestBytes).SetContentType("application/json").Post("https://api.twitch.tv/helix/eventsub/subscriptions")
							// 	if err != nil {
							// 		panic(err)
							// 	}
							//
							// 	if !resp.IsSuccess() {
							// 		panic(resp.String())
							// 	}
							//
							// 	fmt.Print(resp.String())
							//
							// 	subs, err := twitchClient.GetSubscriptions(
							// 		&helix.SubscriptionsParams{
							// 			BroadcasterID: c.ID,
							// 		},
							// 	)
							// 	pretty.Println(subs.Data)
							// }
						},
					)

					socket.OnChannelChatMessage(
						func(
							ccme eventsub.ChannelChatMessageEvent,
							wnm eventsub.WebsocketNotificationMetadata,
						) {
							fmt.Println(ccme.Message.Text)
						},
					)

					socket.OnChannelUpdate(
						func(cue eventsub.ChannelUpdateEvent, wnm eventsub.WebsocketNotificationMetadata) {
							pretty.Println(cue)
						},
					)

					fmt.Println(socket.Connect(context.Background()))
				}()

				return nil
			},
		},
	)
}
