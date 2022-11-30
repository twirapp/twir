package grpc_impl

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/go-redis/redis/v9"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/commands"
	"github.com/satont/tsuwari/apps/parser/internal/permissions"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	"github.com/satont/tsuwari/apps/parser/internal/variables"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"google.golang.org/protobuf/types/known/emptypb"

	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
)

type GrpcImplOpts struct {
	Redis     *redis.Client
	Variables *variables.Variables
	Commands  *commands.Commands
}

type parserGrpcServer struct {
	parser.UnimplementedParserServer

	redis     *redis.Client
	variables *variables.Variables
	commands  *commands.Commands
}

func NewServer(opts *GrpcImplOpts) *parserGrpcServer {
	return &parserGrpcServer{
		redis:     opts.Redis,
		variables: opts.Variables,
		commands:  opts.Commands,
	}
}

func (c *parserGrpcServer) ProcessCommand(
	ctx context.Context,
	data *parser.ProcessCommandRequest,
) (*parser.ProcessCommandResponse, error) {
	if !strings.HasPrefix(data.Message.Text, "!") {
		return nil, nil
	}
	data.Message.Text = data.Message.Text[1:]

	cmds, err := c.commands.GetChannelCommands(data.Channel.Id)
	if err != nil {
		return nil, errors.New("command not found")
	}

	cmd := c.commands.FindByMessage(data.Message.Text, cmds)

	if cmd.Cmd == nil {
		return nil, errors.New("command not found")
	}

	if cmd.Cmd.Cooldown.Valid && cmd.Cmd.CooldownType == cooldownGlobal &&
		cmd.Cmd.Cooldown.Int64 > 0 &&
		c.shouldCheckCooldown(data.Sender.Badges) {
		key := fmt.Sprintf("commands:%s:cooldowns:global", cmd.Cmd.ID)
		_, rErr := c.redis.Get(context.TODO(), key).
			Result()

		if rErr == redis.Nil {
			c.redis.Set(context.TODO(), key, "", time.Duration(cmd.Cmd.Cooldown.Int64)*time.Second)
		} else {
			fmt.Println(err)
			return nil, errors.New("error while setting redis cooldown for command")
		}
	}

	if cmd.Cmd.Cooldown.Valid && cmd.Cmd.CooldownType == cooldownPerUser &&
		cmd.Cmd.Cooldown.Int64 > 0 &&
		c.shouldCheckCooldown(data.Sender.Badges) {
		key := fmt.Sprintf("commands:%s:cooldowns:user:%s", cmd.Cmd.ID, data.Sender.Id)
		_, rErr := c.redis.Get(context.TODO(), key).
			Result()

		if rErr == redis.Nil {
			c.redis.Set(context.TODO(), key, "", time.Duration(cmd.Cmd.Cooldown.Int64)*time.Second)
		} else {
			fmt.Println(err)
			return nil, errors.New("error while setting redis cooldown for command")
		}
	}

	hasPerm := permissions.UserHasPermissionToCommand(data.Sender.Badges, cmd.Cmd.Permission)

	if !hasPerm {
		return nil, errors.New("have no permissions")
	}

	result := c.commands.ParseCommandResponses(cmd, data)

	return result, nil
}

func (c *parserGrpcServer) ParseTextResponse(
	ctx context.Context,
	data *parser.ParseTextRequestData,
) (*parser.ParseTextResponseData, error) {
	isCommand := lo.IfF(data.ParseVariables != nil, func() bool {
		return *data.ParseVariables
	}).ElseF(func() bool { return false })

	cacheService := variables_cache.New(variables_cache.VariablesCacheOpts{
		Text:        &data.Message.Text,
		SenderId:    data.Sender.Id,
		ChannelName: data.Channel.Name,
		ChannelId:   data.Channel.Id,
		SenderName:  &data.Sender.DisplayName,
		Redis:       c.redis,
		Regexp:      nil,
		Twitch:      c.commands.Twitch,
		DB:          c.commands.Db,
		Nats:        c.commands.Nats,
		IsCommand:   isCommand,
	})

	res := c.variables.ParseInput(cacheService, data.Message.Text)

	return &parser.ParseTextResponseData{
		Responses: []string{res},
	}, nil
}

func (c *parserGrpcServer) GetDefaultCommands(
	ctx context.Context,
	data *emptypb.Empty,
) (*parser.GetDefaultCommandsResponse, error) {
	list := make([]*parser.GetDefaultCommandsResponse_DefaultCommand, len(c.commands.DefaultCommands))

	for i, v := range c.commands.DefaultCommands {
		cmd := &parser.GetDefaultCommandsResponse_DefaultCommand{
			Name:               v.Name,
			Description:        *v.Description,
			Visible:            v.Visible,
			Permission:         v.Permission,
			Module:             *v.Module,
			IsReply:            v.IsReply,
			KeepResponsesOrder: v.KeepResponsesOrder,
		}

		list[i] = cmd
	}

	return &parser.GetDefaultCommandsResponse{
		List: list,
	}, nil
}

func (c *parserGrpcServer) GetDefaultVariables(
	ctx context.Context,
	data *emptypb.Empty,
) (*parser.GetVariablesResponse, error) {
	filteredVars := lo.Filter(c.variables.Store, func(i types.Variable, _i int) bool {
		if i.Visible != nil {
			return *i.Visible
		}
		return true
	})

	vars := lo.Map(
		filteredVars,
		func(v types.Variable, _ int) *parser.GetVariablesResponse_Variable {
			desc := v.Name
			if v.Description != nil {
				desc = *v.Description
			}
			example := v.Name
			if v.Example != nil {
				example = *v.Example
			}
			return &parser.GetVariablesResponse_Variable{
				Name:        v.Name,
				Example:     example,
				Description: desc,
			}
		},
	)

	return &parser.GetVariablesResponse{
		List: vars,
	}, nil
}
