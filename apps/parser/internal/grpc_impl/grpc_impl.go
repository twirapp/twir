package grpc_impl

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/lib/pq"
	"github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/cacher"
	"github.com/satont/twir/apps/parser/internal/commands"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/internal/types/services"
	"github.com/satont/twir/apps/parser/internal/variables"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/parser"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	commandsCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "parser_commands_processed",
			Help: "The total number of processed commands",
		},
	)
	textParseCounter = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "parser_text_processed",
			Help: "The total number of processed commands",
		},
	)
)

type ParserGrpcServer struct {
	parser.UnimplementedParserServer

	services  *services.Services
	commands  *commands.Commands
	variables *variables.Variables
}

func NewServer(
	services *services.Services,
	commands *commands.Commands,
	variables *variables.Variables,
) *ParserGrpcServer {
	return &ParserGrpcServer{
		services:  services,
		commands:  commands,
		variables: variables,
	}
}

func (c *ParserGrpcServer) ProcessCommand(
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
		return nil, status.Errorf(codes.Internal, "cannot get channel commands: %w", err)
	}

	cmd := c.commands.FindChannelCommandInInput(data.Message.Text, cmds)
	if cmd.Cmd == nil {
		return nil, status.Error(codes.NotFound, "command not found")
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

	dbUser, _, userRoles, commandRoles, err := c.prepareCooldownAndPermissionsCheck(
		ctx,
		data.Sender.Id,
		data.Channel.Id,
		data.Sender.Badges,
		cmd.Cmd,
	)
	if err != nil {
		return nil, err
	}

	shouldCheckCooldown := c.shouldCheckCooldown(data.Sender.Badges, cmd.Cmd, userRoles)
	if cmd.Cmd.CooldownType == cooldownGlobal && cmd.Cmd.Cooldown.Int64 > 0 && shouldCheckCooldown {
		key := fmt.Sprintf("commands:%s:cooldowns:global", cmd.Cmd.ID)
		rErr := c.services.Redis.Get(ctx, key).Err()

		if rErr == redis.Nil {
			c.services.Redis.Set(ctx, key, "", time.Duration(cmd.Cmd.Cooldown.Int64)*time.Second)
		} else if rErr != nil {
			c.services.Logger.Sugar().Error(rErr)
			return &parser.ProcessCommandResponse{}, errors.New("error while setting redis cooldown for command")
		} else {
			return &parser.ProcessCommandResponse{}, nil
		}
	}

	if cmd.Cmd.CooldownType == cooldownPerUser && cmd.Cmd.Cooldown.Int64 > 0 && shouldCheckCooldown {
		key := fmt.Sprintf("commands:%s:cooldowns:user:%s", cmd.Cmd.ID, data.Sender.Id)
		rErr := c.services.Redis.Get(ctx, key).Err()

		if rErr == redis.Nil {
			c.services.Redis.Set(ctx, key, "", time.Duration(cmd.Cmd.Cooldown.Int64)*time.Second)
		} else if rErr != nil {
			zap.S().Error(rErr)
			return nil, errors.New("error while setting redis cooldown for command")
		} else {
			return &parser.ProcessCommandResponse{}, nil
		}
	}

	hasPerm := c.isUserHasPermissionToCommand(
		data.Sender.Id,
		data.Channel.Id,
		cmd.Cmd,
		dbUser,
		userRoles,
		commandRoles,
	)

	if !hasPerm {
		return nil, errors.New("have no permissions")
	}

	go func() {
		gCtx := context.Background()

		c.services.GrpcClients.Events.CommandUsed(
			// this should be background, because we don't want to wait for response
			gCtx,
			&events.CommandUsedMessage{
				BaseInfo:           &events.BaseInfo{ChannelId: data.Channel.Id},
				CommandId:          cmd.Cmd.ID,
				CommandName:        cmd.Cmd.Name,
				CommandInput:       strings.TrimSpace(data.Message.Text[len(cmd.FoundBy):]),
				UserName:           data.Sender.Name,
				UserDisplayName:    data.Sender.DisplayName,
				UserId:             data.Sender.Id,
				IsDefault:          cmd.Cmd.Default,
				DefaultCommandName: cmd.Cmd.DefaultName.String,
			},
		)

		alert := model.ChannelAlert{}
		if err := c.services.Gorm.Where(
			"channel_id = ? AND command_ids && ?",
			data.Channel.Id,
			pq.StringArray{cmd.Cmd.ID},
		).Find(&alert).Error; err != nil {
			zap.S().Error(err)
			return
		}

		if alert.ID == "" {
			return
		}
		c.services.GrpcClients.WebSockets.TriggerAlert(
			gCtx,
			&websockets.TriggerAlertRequest{
				ChannelId: data.Channel.Id,
				AlertId:   alert.ID,
			},
		)
	}()

	result := c.commands.ParseCommandResponses(ctx, cmd, data)

	return result, nil
}

func (c *ParserGrpcServer) ParseTextResponse(
	ctx context.Context,
	data *parser.ParseTextRequestData,
) (*parser.ParseTextResponseData, error) {
	defer textParseCounter.Inc()

	isCommand := lo.IfF(
		data.ParseVariables != nil, func() bool {
			return *data.ParseVariables
		},
	).ElseF(func() bool { return false })

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
		Cacher: cacher.NewCacher(
			&cacher.CacherOpts{
				Services:        c.services,
				ParseCtxChannel: parseCtxChannel,
				ParseCtxSender:  parseCtxSender,
				ParseCtxText:    &data.Message.Text,
			},
		),
	}

	res := c.variables.ParseVariablesInText(ctx, parseCtx, data.Message.Text)

	return &parser.ParseTextResponseData{
		Responses: []string{res},
	}, nil
}

func (c *ParserGrpcServer) GetDefaultCommands(
	_ context.Context,
	_ *emptypb.Empty,
) (*parser.GetDefaultCommandsResponse, error) {
	list := make(
		[]*parser.GetDefaultCommandsResponse_DefaultCommand,
		0,
		len(c.commands.DefaultCommands),
	)

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

func (c *ParserGrpcServer) GetDefaultVariables(
	_ context.Context,
	_ *emptypb.Empty,
) (*parser.GetVariablesResponse, error) {
	vars := lo.FilterMap(
		lo.Values(c.variables.Store),
		func(v *types.Variable, _i int) (*parser.GetVariablesResponse_Variable, bool) {
			if v.Visible != nil && !*v.Visible {
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
				Name:                v.Name,
				Example:             example,
				Description:         desc,
				Visible:             visible,
				CanBeUsedInRegistry: v.CanBeUsedInRegistry,
			}, true
		},
	)

	return &parser.GetVariablesResponse{
		List: vars,
	}, nil
}
