package commands_bus

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/parser/internal/cacher"
	"github.com/twirapp/twir/apps/parser/internal/commands"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/types/services"
	"github.com/twirapp/twir/apps/parser/internal/variables"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	generic "github.com/twirapp/twir/libs/bus-core/generic"
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
	redis             *redis.Client
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
		redis:             s.Redis,
	}

	return b
}

func (c *CommandsBus) dedupMessage(ctx context.Context, messageID string) (bool, error) {
	if messageID == "" {
		return false, nil
	}
	key := fmt.Sprintf("parser:dedup:%s", messageID)
	set, err := c.redis.SetNX(ctx, key, "1", 60*time.Second).Result()
	if err != nil {
		return false, err
	}
	return !set, nil
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
			res, err := c.commandService.ProcessChatMessage(ctx, data, "twitch")
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
					Platform:      "twitch",
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

	// TODO(Phase-2): remove ProcessMessageAsCommand subscription once all consumers migrated off it
	c.bus.Parser.ProcessMessageAsCommand.SubscribeGroup(
		"parser",
		func(
			ctx context.Context,
			data twitch.TwitchChatMessage,
		) (struct{}, error) {
			isDup, err := c.dedupMessage(ctx, data.MessageId)
			if err != nil {
				zap.S().Error(err)
			} else if isDup {
				return struct{}{}, nil
			}

			res, err := c.commandService.ProcessChatMessage(ctx, data, "twitch")
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

	c.bus.Parser.ProcessGenericMessage.SubscribeGroup(
		"parser",
		func(
			ctx context.Context,
			msg generic.ChatMessage,
		) (struct{}, error) {
			isDup, err := c.dedupMessage(ctx, msg.MessageID)
			if err != nil {
				zap.S().Error(err)
			} else if isDup {
				return struct{}{}, nil
			}

			badges := make([]twitch.ChatMessageBadge, 0, len(msg.Badges))
			for _, b := range msg.Badges {
				badges = append(badges, twitch.ChatMessageBadge{
					SetId: b.SetID,
					Info:  b.Text,
				})
			}

			twitchMsg := twitch.TwitchChatMessage{
				BroadcasterUserId:    msg.ChannelID,
				BroadcasterUserLogin: msg.PlatformChannelID,
				BroadcasterUserName:  msg.PlatformChannelID,
				ChatterUserId:        msg.UserID,
				ChatterUserLogin:     msg.SenderLogin,
				ChatterUserName:      msg.SenderDisplayName,
				MessageId:            msg.MessageID,
				Badges:               badges,
				Message: &twitch.ChatMessageMessage{
					Text: msg.Text,
				},
			}

			res, err := c.commandService.ProcessChatMessage(ctx, twitchMsg, msg.Platform)
			if err != nil {
				zap.S().Error(err)
				return struct{}{}, err
			}
			if res == nil {
				return struct{}{}, nil
			}

			var replyTo string
			if res.IsReply {
				replyTo = msg.MessageID
			}

			channelName := msg.PlatformChannelID
			for _, r := range res.Responses {
				params := bots.SendMessageRequest{
					ChannelName:       &channelName,
					ChannelId:         msg.ChannelID,
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
	c.bus.Parser.ProcessGenericMessage.Unsubscribe()
	c.bus.Parser.GetDefaultCommands.Unsubscribe()
}
