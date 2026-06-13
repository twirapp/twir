package commands_bus

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
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
	platformentity "github.com/twirapp/twir/libs/entities/platform"
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
		func(ctx context.Context, data generic.ChatMessage) (parser.CommandParseResponse, error) {
			res, err := c.commandService.ProcessChatMessage(ctx, data, platformentity.Platform(data.Platform))
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

			twitchUserID, _ := uuid.Parse(data.ChannelTwitchUserID)
			channel := &types.ParseContextChannel{
				ID:           data.ChannelID,
				Name:         data.ChannelName,
				TwitchUserID: twitchUserID,
				DBChannelID:  data.ChannelDBID,
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
					Platform:      platformentity.PlatformTwitch,
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
			data generic.ChatMessage,
		) (struct{}, error) {
			isDup, err := c.dedupMessage(ctx, data.MessageID)
			if err != nil {
				zap.S().Error(err)
			} else if isDup {
				return struct{}{}, nil
			}

			res, err := c.commandService.ProcessChatMessage(ctx, data, platformentity.Platform(data.Platform))
			if err != nil {
				zap.S().Error(err)
				return struct{}{}, err
			}
			if res == nil {
				return struct{}{}, nil
			}

			var replyTo string
			if res.IsReply {
				replyTo = data.MessageID
			}

			for _, r := range res.Responses {
				internalChannelID := data.EnrichedData.DbChannel.ID
				channelName := data.BroadcasterUserLogin
				if channelName == "" {
					channelName = data.PlatformChannelID
				}
				params := bots.SendMessageRequest{
					ChannelName:       &channelName,
					ChannelId:         data.PlatformChannelID,
					InternalChannelID: &internalChannelID,
					PlatformChannelID: data.PlatformChannelID,
					Platform:          data.Platform,
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
func genericToSendMessageRequest(
	msg generic.ChatMessage,
	channelName *string,
	message string,
	replyTo string,
	skipToxicityCheck bool,
) bots.SendMessageRequest {
	var internalChannelID *uuid.UUID
	if parsedChannelID, err := uuid.Parse(msg.ChannelID); err == nil {
		internalChannelID = &parsedChannelID
	}

	return bots.SendMessageRequest{
		ChannelName:       channelName,
		ChannelId:         msg.PlatformChannelID,
		InternalChannelID: internalChannelID,
		PlatformChannelID: msg.PlatformChannelID,
		Platform:          msg.Platform,
		Message:           message,
		ReplyTo:           replyTo,
		SkipRateLimits:    false,
		SkipToxicityCheck: skipToxicityCheck,
	}
}

func (c *CommandsBus) Unsubscribe() {
	c.bus.Parser.GetCommandResponse.Unsubscribe()
	c.bus.Parser.ParseVariablesInText.Unsubscribe()
	c.bus.Parser.ProcessMessageAsCommand.Unsubscribe()
	c.bus.Parser.GetDefaultCommands.Unsubscribe()
}
