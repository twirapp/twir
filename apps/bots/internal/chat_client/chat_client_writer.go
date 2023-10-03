package chat_client

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/satont/twir/libs/grpc/generated/tokens"
)

func (c *ChatClient) CreateWriter() {
	client := irc.NewClient(c.TwitchUser.Login, "")
	client.Capabilities = capabilities

	c.Writer = &BotClientIrc{
		size:            0,
		disconnectChann: make(chan struct{}),
		Client:          client,
	}

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
					c.onConnect("Writer")
				},
			)
			client.OnSelfJoinMessage(
				func(m irc.UserJoinMessage) {
					c.onSelfJoin(m, "Writer")
				},
			)

			channels, err := c.getChannels()
			if err != nil {
				c.services.Logger.Error("cannot get channels", slog.Any("err", err))
				return
			}

			client.Join(channels...)

			connectResultCh := make(chan error)
			go func() {
				// Perform the connection attempt in a goroutine.
				err := client.Connect()
				connectResultCh <- err
			}()

		connLoop:
			for {
				select {
				case <-c.Writer.disconnectChann:
					// Signal received, initiate disconnect and break the loop.
					client.Disconnect()
					c.services.Logger.Info("disconnected", slog.Any("err", err))
					break mainLoop
				case err := <-connectResultCh:
					// Handle the result of the connection attempt.
					if err != nil {
						c.services.Logger.Error("disconnected", slog.Any("err", err))
					}
					break connLoop
				}
			}
		}
	}()
}
