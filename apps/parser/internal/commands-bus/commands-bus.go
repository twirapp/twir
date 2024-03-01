package commands_bus

import (
	"context"

	"github.com/satont/twir/apps/parser/internal/cacher"
	"github.com/satont/twir/apps/parser/internal/commands"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/internal/types/services"
	"github.com/satont/twir/apps/parser/internal/variables"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"go.uber.org/zap"
)

type CommandsBus struct {
	bus              *buscore.Bus
	services         *services.Services
	commandService   *commands.Commands
	variablesService *variables.Variables
}

func New(
	bus *buscore.Bus,
	s *services.Services,
	commandService *commands.Commands,
	variablesService *variables.Variables,
) *CommandsBus {
	b := &CommandsBus{
		bus:              bus,
		services:         s,
		commandService:   commandService,
		variablesService: variablesService,
	}

	return b
}

func (c *CommandsBus) Subscribe() error {
	c.bus.Parser.GetCommandResponse.SubscribeGroup(
		"parser",
		func(ctx context.Context, data twitch.TwitchChatMessage) parser.CommandParseResponse {
			res, err := c.commandService.ProcessChatMessage(ctx, data)
			if err != nil || res == nil {
				return parser.CommandParseResponse{}
			}

			return *res
		},
	)

	c.bus.Parser.ParseVariablesInText.SubscribeGroup(
		"parser",
		func(
			ctx context.Context,
			data parser.ParseVariablesInTextRequest,
		) parser.ParseVariablesInTextResponse {
			channel := &types.ParseContextChannel{
				ID:   data.ChannelID,
				Name: data.ChannelName,
			}
			sender := &types.ParseContextSender{
				ID:          data.UserID,
				Name:        data.UserLogin,
				DisplayName: data.UserName,
			}
			parsed := c.variablesService.ParseVariablesInText(
				ctx,
				&types.ParseContext{
					MessageId: "",
					Channel:   channel,
					Sender:    sender,
					Emotes:    nil,
					Mentions:  nil,
					Text:      &data.Text,
					RawText:   data.Text,
					IsCommand: false,
					Services:  c.services,
					Cacher: cacher.NewCacher(
						&cacher.CacherOpts{
							Services:        c.services,
							ParseCtxChannel: channel,
							ParseCtxSender:  sender,
							ParseCtxText:    &data.Text,
						},
					),
				},
				data.Text,
			)

			return parser.ParseVariablesInTextResponse{
				Text: parsed,
			}
		},
	)

	c.bus.Parser.ProcessMessageAsCommand.SubscribeGroup(
		"parser",
		func(
			ctx context.Context,
			data twitch.TwitchChatMessage,
		) struct{} {
			res, err := c.commandService.ProcessChatMessage(ctx, data)
			if err != nil {
				zap.S().Error(err)
				return struct{}{}
			}
			if res == nil {
				return struct{}{}
			}

			var replyTo string
			if res.IsReply {
				replyTo = data.MessageId
			}

			for _, r := range res.Responses {
				c.bus.Bots.SendMessage.Publish(
					bots.SendMessageRequest{
						ChannelId:      data.BroadcasterUserId,
						ChannelName:    &data.BroadcasterUserLogin,
						Message:        r,
						SkipRateLimits: false,
						ReplyTo:        replyTo,
					},
				)
			}

			return struct{}{}
		},
	)

	return nil
}

func (c *CommandsBus) Unsubscribe() {
	c.bus.Parser.GetCommandResponse.Unsubscribe()
	c.bus.Parser.ParseVariablesInText.Unsubscribe()
	c.bus.Parser.ProcessMessageAsCommand.Unsubscribe()
}
