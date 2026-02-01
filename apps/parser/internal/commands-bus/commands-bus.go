package commands_bus

import (
	"context"
	"strings"

	"github.com/twirapp/twir/apps/parser/internal/cacher"
	"github.com/twirapp/twir/apps/parser/internal/commands"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/types/services"
	"github.com/twirapp/twir/apps/parser/internal/variables"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/repositories/streams"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
	"go.uber.org/zap"
)

type CommandsBus struct {
	bus               *buscore.Bus
	services          *services.Services
	commandService    *commands.Commands
	variablesService  *variables.Variables
	streamsRepository streams.Repository
}

func New(
	bus *buscore.Bus,
	s *services.Services,
	commandService *commands.Commands,
	variablesService *variables.Variables,
	streamsRepository streams.Repository,
) *CommandsBus {
	b := &CommandsBus{
		bus:               bus,
		services:          s,
		commandService:    commandService,
		variablesService:  variablesService,
		streamsRepository: streamsRepository,
	}

	return b
}

func (c *CommandsBus) Subscribe() error {
	c.bus.Parser.GetDefaultCommands.SubscribeGroup(
		"parser",
		func(ctx context.Context, data struct{}) (parser.GetDefaultCommandsResponse, error) {
			resp := parser.GetDefaultCommandsResponse{
				List: make([]parser.DefaultCommand, 0, len(c.commandService.DefaultCommands)),
			}

			for _, cmd := range c.commandService.DefaultCommands {
				rolesNames := make([]string, len(cmd.RolesIDS))
				for i, roleID := range cmd.RolesIDS {
					rolesNames[i] = roleID
				}

				resp.List = append(
					resp.List,
					parser.DefaultCommand{
						Name:               cmd.Name,
						Description:        cmd.Description.String,
						Visible:            cmd.Visible,
						RolesNames:         rolesNames,
						Module:             cmd.Module,
						IsReply:            cmd.IsReply,
						KeepResponsesOrder: cmd.KeepResponsesOrder,
						Aliases:            cmd.Aliases,
					},
				)
			}

			return resp, nil
		},
	)

	c.bus.Parser.GetCommandResponse.SubscribeGroup(
		"parser",
		func(ctx context.Context, data twitch.TwitchChatMessage) (parser.CommandParseResponse, error) {
			res, err := c.commandService.ProcessChatMessage(ctx, data)
			if err != nil {
				return parser.CommandParseResponse{}, err
			}
			if res == nil {
				return parser.CommandParseResponse{}, nil
			}

			return *res, nil
		},
	)

	c.bus.Parser.ParseVariablesInText.SubscribeGroup(
		"parser",
		func(
			ctx context.Context,
			data parser.ParseVariablesInTextRequest,
		) (parser.ParseVariablesInTextResponse, error) {
			foundStream, err := c.streamsRepository.GetByChannelID(ctx, data.ChannelID)
			if err != nil {
				zap.S().Error(err)
				return parser.ParseVariablesInTextResponse{}, err
			}

			var stream *streamsmodel.Stream
			if foundStream.ID != "" {
				stream = &foundStream
			}

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
					MessageId:     "",
					Channel:       channel,
					Sender:        sender,
					Emotes:        nil,
					Mentions:      data.Mentions,
					Text:          &data.Text,
					RawText:       data.Text,
					IsCommand:     data.IsCommand,
					IsInCustomVar: data.IsInCustomVar,
					Services:      c.services,
					Cacher: cacher.NewCacher(
						&cacher.CacherOpts{
							Services:        c.services,
							ParseCtxChannel: channel,
							ParseCtxSender:  sender,
							ParseCtxText:    &data.Text,
						},
					),
					ChannelStream: stream,
				},
				data.Text,
			)

			return parser.ParseVariablesInTextResponse{
				Text: strings.Join(parsed, " "),
			}, nil
		},
	)

	c.bus.Parser.ProcessMessageAsCommand.SubscribeGroup(
		"parser",
		func(
			ctx context.Context,
			data twitch.TwitchChatMessage,
		) (struct{}, error) {
			res, err := c.commandService.ProcessChatMessage(ctx, data)
			if err != nil {
				zap.S().Error(err)
				return struct{}{}, err
			}
			if res == nil {
				return struct{}{}, nil
			}

			var replyTo string
			if res.IsReply {
				replyTo = data.MessageId
			}

			for _, r := range res.Responses {
				params := bots.SendMessageRequest{
					ChannelName:       &data.BroadcasterUserLogin,
					ChannelId:         data.BroadcasterUserId,
					Message:           r,
					ReplyTo:           replyTo,
					SkipRateLimits:    false,
					SkipToxicityCheck: res.SkipToxicityCheck,
				}

				if res.KeepOrder {
					if _, err := c.bus.Bots.SendMessage.Request(ctx, params); err != nil {
						zap.S().Error(err)
					}
				} else {
					withoutCancel := context.WithoutCancel(ctx)
					go func() {
						if err := c.bus.Bots.SendMessage.Publish(withoutCancel, params); err != nil {
							zap.S().Error(err)
						}
					}()
				}
			}

			return struct{}{}, nil
		},
	)

	return nil
}

func (c *CommandsBus) Unsubscribe() {
	c.bus.Parser.GetCommandResponse.Unsubscribe()
	c.bus.Parser.ParseVariablesInText.Unsubscribe()
	c.bus.Parser.ProcessMessageAsCommand.Unsubscribe()
	c.bus.Parser.GetDefaultCommands.Unsubscribe()
}
