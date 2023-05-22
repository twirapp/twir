package grpc_impl

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser-new/internal/cacher"
	"github.com/satont/tsuwari/apps/parser-new/internal/commands"
	"github.com/satont/tsuwari/apps/parser-new/internal/types"
	"github.com/satont/tsuwari/apps/parser-new/internal/types/services"
	"github.com/satont/tsuwari/apps/parser-new/internal/variables"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/events"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	commandsCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "parser_commands_processed",
		Help: "The total number of processed commands",
	})
	textParseCounter = promauto.NewCounter(prometheus.CounterOpts{
		Name: "parser_text_processed",
		Help: "The total number of processed commands",
	})
)

type parserGrpcServer struct {
	parser.UnimplementedParserServer

	services  *services.Services
	commands  *commands.Commands
	variables *variables.Variables
}

func NewServer(
	services *services.Services,
	commands *commands.Commands,
	variables *variables.Variables,
) *parserGrpcServer {
	return &parserGrpcServer{
		services:  services,
		commands:  commands,
		variables: variables,
	}
}

func (c *parserGrpcServer) ProcessCommand(
	ctx context.Context,
	data *parser.ProcessCommandRequest,
) (*parser.ProcessCommandResponse, error) {
	defer commandsCounter.Inc()

	if !strings.HasPrefix(data.Message.Text, "!") {
		return nil, nil
	}
	data.Message.Text = data.Message.Text[1:]

	cmds, err := c.commands.GetChannelCommands(ctx, data.Channel.Id)
	if err != nil {
		return nil, errors.New("command not found")
	}

	cmd := c.commands.FindChannelCommandInInput(data.Message.Text, cmds)

	if cmd.Cmd == nil {
		return nil, errors.New("command not found")
	}

	if cmd.Cmd.OnlineOnly {
		stream := &model.ChannelsStreams{}
		err = c.services.Gorm.
			WithContext(ctx).
			Where(`"userId" = ?`, data.Channel.Id).
			Find(stream).Error
		if err != nil {
			c.services.Logger.Sugar().Error(err)
			return nil, err
		}
		if stream == nil || stream.ID == "" {
			return &parser.ProcessCommandResponse{}, errors.New("stream is offline")
		}
	}

	if cmd.Cmd.Cooldown.Valid && cmd.Cmd.CooldownType == cooldownGlobal &&
		cmd.Cmd.Cooldown.Int64 > 0 &&
		c.shouldCheckCooldown(data.Sender.Badges) {
		key := fmt.Sprintf("commands:%s:cooldowns:global", cmd.Cmd.ID)
		rErr := c.services.Redis.Get(context.TODO(), key).Err()

		if rErr == redis.Nil {
			c.services.Redis.Set(context.TODO(), key, "", time.Duration(cmd.Cmd.Cooldown.Int64)*time.Second)
		} else if rErr != nil {
			c.services.Logger.Sugar().Error(rErr)
			return nil, errors.New("error while setting redis cooldown for command")
		}
	}

	if cmd.Cmd.Cooldown.Valid && cmd.Cmd.CooldownType == cooldownPerUser &&
		cmd.Cmd.Cooldown.Int64 > 0 &&
		c.shouldCheckCooldown(data.Sender.Badges) {
		key := fmt.Sprintf("commands:%s:cooldowns:user:%s", cmd.Cmd.ID, data.Sender.Id)
		rErr := c.services.Redis.Get(context.TODO(), key).Err()

		if rErr == redis.Nil {
			c.services.Redis.Set(context.TODO(), key, "", time.Duration(cmd.Cmd.Cooldown.Int64)*time.Second)
		} else if rErr != nil {
			zap.S().Error(rErr)
			return nil, errors.New("error while setting redis cooldown for command")
		}
	}

	hasPerm := c.isUserHasPermissionToCommand(
		ctx,
		data.Sender.Id,
		data.Channel.Id,
		data.Sender.Badges,
		cmd.Cmd,
	)

	if !hasPerm {
		return nil, errors.New("have no permissions")
	}

	result := c.commands.ParseCommandResponses(ctx, cmd, data)

	defer c.services.GrpcClients.Events.CommandUsed(context.Background(), &events.CommandUsedMessage{
		BaseInfo:        &events.BaseInfo{ChannelId: data.Channel.Id},
		CommandId:       cmd.Cmd.ID,
		CommandName:     cmd.Cmd.Name,
		CommandInput:    strings.TrimSpace(data.Message.Text[len(cmd.FoundBy):]),
		UserName:        data.Sender.Name,
		UserDisplayName: data.Sender.DisplayName,
		UserId:          data.Sender.Id,
	})

	return result, nil
}

func (c *parserGrpcServer) ParseTextResponse(
	ctx context.Context,
	data *parser.ParseTextRequestData,
) (*parser.ParseTextResponseData, error) {
	defer textParseCounter.Inc()

	isCommand := lo.IfF(data.ParseVariables != nil, func() bool {
		return *data.ParseVariables
	}).ElseF(func() bool { return false })

	parseCtxChannel := &types.ParseContextChannel{
		ID:   data.Channel.Id,
		Name: data.Channel.Name,
	}
	parseCtxSender := &types.ParseContextSender{
		ID:          data.Sender.Id,
		Name:        data.Sender.Name,
		DisplayName: data.Sender.DisplayName,
		Badges:      data.Sender.Badges,
	}
	parseCtx := &types.ParseContext{
		Channel:   parseCtxChannel,
		Sender:    parseCtxSender,
		Text:      &data.Message.Text,
		Services:  c.services,
		IsCommand: isCommand,
		Cacher: cacher.NewCacher(&cacher.CacherOpts{
			Services:        c.services,
			ParseCtxChannel: parseCtxChannel,
			ParseCtxSender:  parseCtxSender,
			ParseCtxText:    &data.Message.Text,
		}),
	}

	res := c.variables.ParseVariablesInText(ctx, parseCtx, data.Message.Text)

	return &parser.ParseTextResponseData{
		Responses: []string{res},
	}, nil
}

func (c *parserGrpcServer) GetDefaultCommands(
	_ context.Context,
	_ *emptypb.Empty,
) (*parser.GetDefaultCommandsResponse, error) {
	list := make([]*parser.GetDefaultCommandsResponse_DefaultCommand, len(c.commands.DefaultCommands))

	for _, v := range c.commands.DefaultCommands {
		cmd := &parser.GetDefaultCommandsResponse_DefaultCommand{
			Name:               v.ChannelsCommands.Name,
			Description:        v.ChannelsCommands.Description.String,
			Visible:            v.ChannelsCommands.Visible,
			RolesNames:         v.ChannelsCommands.RolesIDS,
			Module:             v.ChannelsCommands.Module,
			IsReply:            v.ChannelsCommands.IsReply,
			KeepResponsesOrder: v.ChannelsCommands.KeepResponsesOrder,
			Aliases:            v.ChannelsCommands.Aliases,
		}

		list = append(list, cmd)
	}

	return &parser.GetDefaultCommandsResponse{
		List: list,
	}, nil
}

func (c *parserGrpcServer) GetDefaultVariables(
	ctx context.Context,
	data *emptypb.Empty,
) (*parser.GetVariablesResponse, error) {
	vars := lo.FilterMap(
		lo.Values(c.variables.Store),
		func(v *types.Variable, _i int) (*parser.GetVariablesResponse_Variable, bool) {
			if v.Visible == nil || !*v.Visible {
				return nil, false
			}

			desc := v.Name
			if v.Description != nil {
				desc = *v.Description
			}
			example := v.Name
			if v.Example != nil {
				example = *v.Example
			}

			visible := true
			if v.Visible != nil {
				visible = *v.Visible
			}

			return &parser.GetVariablesResponse_Variable{
				Name:        v.Name,
				Example:     example,
				Description: desc,
				Visible:     visible,
			}, true
		})

	return &parser.GetVariablesResponse{
		List: vars,
	}, nil
}
