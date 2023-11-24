package chat_client

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/google/uuid"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/events"
	"github.com/satont/twir/libs/grpc/generated/tokens"
)

func (c *ChatClient) createReader() *BotClientIrc {
	client := irc.NewClient(c.TwitchUser.Login, "")
	client.Capabilities = capabilities

	reader := &BotClientIrc{
		size:            0,
		disconnectChann: make(chan struct{}),
		Client:          client,
	}
	shardId := len(c.Readers) + 1

	go func() {
	mainLoop:
		for {
			time.Sleep(500 * time.Millisecond)

			tokenReqCtx, cancelTokenReqCtx := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancelTokenReqCtx()

			token, err := c.services.TokensGrpc.RequestBotToken(
				tokenReqCtx,
				&tokens.GetBotTokenRequest{
					BotId: c.Model.ID,
				},
			)
			if err != nil {
				c.services.Logger.Error(err.Error())
				continue
			}

			c.services.TwitchClient.SetUserAccessToken(token.AccessToken)

			// joinChannels(opts.DB, opts.Cfg, opts.Logger, &client)
			client.SetIRCToken(fmt.Sprintf("oauth:%s", token.AccessToken))

			client.OnConnect(
				func() {
					reader.Connected = true
					c.onConnect(fmt.Sprintf("Reader #%v", shardId))
				},
			)
			client.OnSelfJoinMessage(
				func(m irc.UserJoinMessage) {
					c.onSelfJoin(m, fmt.Sprint(shardId))
				},
			)
			client.OnUserStateMessage(
				func(message irc.UserStateMessage) {
					if message.User.ID == c.TwitchUser.ID && c.services.Cfg.AppEnv != "development" {
						return
					}
					c.OnUserStateMessage(message)
				},
			)
			client.OnUserNoticeMessage(
				func(message irc.UserNoticeMessage) {
					c.counters.messagesCounter.Inc()

					if message.User.ID == c.TwitchUser.ID && c.services.Cfg.AppEnv != "development" {
						return
					}
					c.onMessage(
						Message{
							ID: message.ID,
							Channel: MessageChannel{
								ID:   message.RoomID,
								Name: message.Channel,
							},
							User: MessageUser{
								ID:          message.User.ID,
								Name:        message.User.Name,
								DisplayName: message.User.DisplayName,
								Badges:      message.User.Badges,
							},
							Message: message.Message,
							Emotes:  message.Emotes,
							Tags:    message.Tags,
						},
					)
				},
			)
			client.OnPrivateMessage(
				func(message irc.PrivateMessage) {
					c.counters.messagesCounter.Inc()

					if message.User.ID == c.TwitchUser.ID && c.services.Cfg.AppEnv != "development" {
						return
					}
					c.onMessage(
						Message{
							ID: message.ID,
							Channel: MessageChannel{
								ID:   message.RoomID,
								Name: message.Channel,
							},
							User: MessageUser{
								ID:          message.User.ID,
								Name:        message.User.Name,
								DisplayName: message.User.DisplayName,
								Badges:      message.User.Badges,
							},
							Message: message.Message,
							Emotes:  message.Emotes,
							Tags:    message.Tags,
						},
					)
				},
			)
			client.OnClearChatMessage(
				func(message irc.ClearChatMessage) {
					c.counters.messagesCounter.Inc()

					if message.TargetUserID != "" {
						return
					}
					c.services.DB.Create(
						&model.ChannelsEventsListItem{
							ID:        uuid.New().String(),
							ChannelID: message.RoomID,
							Type:      model.ChannelEventListItemTypeChatClear,
							CreatedAt: time.Now().UTC(),
							Data:      &model.ChannelsEventsListItemData{},
						},
					)
					c.services.EventsGrpc.ChatClear(
						context.Background(), &events.ChatClearMessage{
							BaseInfo: &events.BaseInfo{ChannelId: message.RoomID},
						},
					)
				},
			)
			client.OnUserNoticeMessage(c.onNotice)
			client.OnUserJoinMessage(c.onUserJoin)

			if err != nil {
				c.services.Logger.Error("cannot get channels", slog.Any("err", err))
				return
			}

			connectResultCh := make(chan error)
			go func() {
				// Perform the connection attempt in a goroutine.
				err := client.Connect()
				connectResultCh <- err
			}()

		connLoop:
			for {
				select {
				case <-reader.disconnectChann:
					// Signal received, initiate disconnect and break the loop.
					client.Disconnect()
					reader.Connected = false
					c.services.Logger.Info("reader disconnected manually", slog.Int("shardId", shardId))
					c.Readers = lo.Filter(
						c.Readers,
						func(item *BotClientIrc, _ int) bool {
							return item != reader
						},
					)
					break mainLoop
				case err := <-connectResultCh:
					// Handle the result of the connection attempt.
					if err != nil {
						reader.Connected = false
						c.services.Logger.Error(
							"reader disconnected by error",
							slog.String("err", err.Error()),
						)
					}
					close(connectResultCh)
					break connLoop
				}
			}
		}
	}()

	return reader
}
