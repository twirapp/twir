package grpc_impl

import (
	"context"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/commands"
	"github.com/satont/twir/apps/parser/internal/types"
	"github.com/satont/twir/apps/parser/internal/types/services"
	"github.com/satont/twir/apps/parser/internal/variables"
	"github.com/twirapp/twir/libs/grpc/parser"
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
