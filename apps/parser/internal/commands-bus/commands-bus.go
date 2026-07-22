package commands_bus

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/parser/internal/cacher"
	"github.com/twirapp/twir/apps/parser/internal/channelbinding"
	"github.com/twirapp/twir/apps/parser/internal/commands"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/types/services"
	"github.com/twirapp/twir/apps/parser/internal/variables"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	generic "github.com/twirapp/twir/libs/bus-core/generic"
	"github.com/twirapp/twir/libs/bus-core/parser"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	"github.com/twirapp/twir/libs/repositories/streams"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
	"github.com/twirapp/twir/libs/repositories/userswithstats"
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

func (c *CommandsBus) enrichChatMessage(
	ctx context.Context,
	data generic.ChatMessage,
) (commands.ChatMessageContext, error) {
	messagePlatform := platformentity.Platform(data.Platform)
	if messagePlatform == "" {
		messagePlatform = platformentity.PlatformTwitch
	}

	channelID, err := uuid.Parse(data.ChannelID)
	if err != nil {
		return commands.ChatMessageContext{}, fmt.Errorf("parse message channel id: %w", err)
	}
	channel, err := c.services.ChannelService.GetChannelByID(ctx, channelID)
	if err != nil {
		return commands.ChatMessageContext{}, fmt.Errorf("get message channel: %w", err)
	}

	userID, err := uuid.Parse(data.UserID)
	if err != nil {
		return commands.ChatMessageContext{}, fmt.Errorf("parse message user id: %w", err)
	}
	userWithStats, err := c.services.UsersWithStatsRepository.GetByUserAndChannelID(
		ctx,
		userswithstats.GetByUserAndChannelIDInput{
			UserID:    userID,
			ChannelID: channelID,
		},
	)
	if err != nil {
		return commands.ChatMessageContext{}, fmt.Errorf("get message user and stats: %w", err)
	}

	commandsPrefix := "!"
	prefix, err := c.services.CommandsPrefixCache.Get(ctx, channelID.String())
	if err != nil && !errors.Is(err, channelscommandsprefixrepository.ErrNotFound) {
		return commands.ChatMessageContext{}, fmt.Errorf("get message command prefix: %w", err)
	}
	if err == nil && prefix.Prefix != "" {
		commandsPrefix = prefix.Prefix
	}

	stream, err := c.streamsRepository.GetByChannelID(ctx, channelID, messagePlatform)
	if err != nil {
		return commands.ChatMessageContext{}, fmt.Errorf("get message channel stream: %w", err)
	}
	var channelStream *streamsmodel.Stream
	if stream.ID != "" {
		channelStream = &stream
	}

	return commands.ChatMessageContext{
		ChatMessage:   data,
		Channel:       channel,
		Stream:        channelStream,
		User:          userWithStats.User,
		UserStats:     userWithStats.Stats,
		CommandPrefix: commandsPrefix,
	}, nil
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
			message, err := c.enrichChatMessage(ctx, data)
			if err != nil {
				return parser.CommandParseResponse{}, err
			}

			res, err := c.commandService.ProcessChatMessage(
				ctx,
				message,
				platformentity.Platform(data.Platform),
			)
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
			platformSource := platformentity.PlatformTwitch
			if data.PlatformSource != nil {
				platformSource = *data.PlatformSource
			}

			channelModel, err := c.services.ChannelService.GetChannelByID(ctx, data.ChannelID)
			if err != nil {
				zap.S().Error(err)
				return parser.ParseVariablesInTextResponse{}, err
			}

			channel, ok := channelbinding.NewParseContextChannel(
				channelModel,
				platformSource,
				data.ChannelName,
			)
			if !ok {
				return parser.ParseVariablesInTextResponse{}, fmt.Errorf("channel %s is not connected to %s", data.ChannelID, platformSource)
			}

			foundStream, err := c.streamsRepository.GetByChannelID(ctx, data.ChannelID, platformSource)
			if err != nil {
				zap.S().Error(err)
				return parser.ParseVariablesInTextResponse{}, err
			}

			var stream *streamsmodel.Stream
			if foundStream.ID != "" {
				stream = &foundStream
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
					Platform:      platformSource,
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

			message, err := c.enrichChatMessage(ctx, data)
			if err != nil {
				zap.S().Error(err)
				return struct{}{}, err
			}

			res, err := c.commandService.ProcessChatMessage(
				ctx,
				message,
				platformentity.Platform(data.Platform),
			)
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

			messagePlatform := platformentity.Platform(data.Platform)
			if messagePlatform == "" {
				messagePlatform = platformentity.PlatformTwitch
			}

			for _, r := range res.Responses {
				params := bots.SendMessageRequest{
					ChannelID:         message.Channel.ID,
					Platforms:         []platformentity.Platform{messagePlatform},
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
