package mcp

import (
	"context"
	"fmt"

	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands"
	commandsrelations "github.com/twirapp/twir/apps/api-gql/internal/services/commands_with_groups_and_responses"
)

type commandResponseInput struct {
	Text        string `json:"text" jsonschema:"command response text"`
	Order       int    `json:"order,omitempty" jsonschema:"response order, starting at zero"`
	OnlineOnly  bool   `json:"onlineOnly,omitempty"`
	OfflineOnly bool   `json:"offlineOnly,omitempty"`
}

type createCommandInput struct {
	Name               string                 `json:"name" jsonschema:"command name without the prefix"`
	Responses          []commandResponseInput `json:"responses" jsonschema:"one or more command responses"`
	Aliases            []string               `json:"aliases,omitempty"`
	Description        string                 `json:"description,omitempty"`
	Cooldown           int                    `json:"cooldown,omitempty" jsonschema:"cooldown in seconds"`
	CooldownType       string                 `json:"cooldownType,omitempty" jsonschema:"GLOBAL or PER_USER"`
	Enabled            bool                   `json:"enabled"`
	Visible            bool                   `json:"visible"`
	OnlineOnly         bool                   `json:"onlineOnly,omitempty"`
	KeepResponsesOrder bool                   `json:"keepResponsesOrder,omitempty"`
}

type updateCommandInput struct {
	ID                 string                 `json:"id" jsonschema:"command UUID"`
	Name               *string                `json:"name,omitempty"`
	Responses          []commandResponseInput `json:"responses,omitempty"`
	Aliases            []string               `json:"aliases,omitempty"`
	Description        *string                `json:"description,omitempty"`
	Cooldown           *int                   `json:"cooldown,omitempty"`
	CooldownType       *string                `json:"cooldownType,omitempty"`
	Enabled            *bool                  `json:"enabled,omitempty"`
	Visible            *bool                  `json:"visible,omitempty"`
	OnlineOnly         *bool                  `json:"onlineOnly,omitempty"`
	KeepResponsesOrder *bool                  `json:"keepResponsesOrder,omitempty"`
}

type idInput struct {
	ID string `json:"id" jsonschema:"object UUID"`
}

func (h *Handler) addCommandTools(s *modelsdk.Server, requestScope scope) {
	modelsdk.AddTool(s, &modelsdk.Tool{Name: "list_commands", Description: "List all commands, groups, responses, cooldowns, and permissions for this channel."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			items, err := h.deps.CommandsRelations.GetManyByChannelID(ctx, requestScope.Channel.ID.String())
			return nil, map[string]any{"commands": items}, err
		})

	modelsdk.AddTool(s, &modelsdk.Tool{Name: "get_command", Description: "Get one channel command by UUID."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
			id, err := parseID(input.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid command id: %w", err)
			}
			items, err := h.deps.CommandsRelations.GetManyByChannelID(ctx, requestScope.Channel.ID.String())
			if err != nil {
				return nil, nil, err
			}
			for _, item := range items {
				if item.Command.ID == id {
					return nil, item, nil
				}
			}
			return nil, nil, fmt.Errorf("command not found")
		})

	modelsdk.AddTool(s, &modelsdk.Tool{Name: "create_command", Description: "Create a command for this channel."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input createCommandInput) (*modelsdk.CallToolResult, any, error) {
			responses := make([]commands.CreateInputResponse, 0, len(input.Responses))
			for _, response := range input.Responses {
				responses = append(responses, commands.CreateInputResponse{Text: &response.Text, Order: response.Order, OnlineOnly: response.OnlineOnly, OfflineOnly: response.OfflineOnly})
			}
			result, err := h.deps.Commands.Create(ctx, commands.CreateInput{
				ChannelID: requestScope.Channel.ID.String(), ActorID: requestScope.ActorID,
				Name: input.Name, Responses: responses, Aliases: input.Aliases,
				Description: input.Description, Cooldown: input.Cooldown, CooldownType: input.CooldownType,
				Enabled: input.Enabled, Visible: input.Visible, OnlineOnly: input.OnlineOnly,
				KeepResponsesOrder: input.KeepResponsesOrder,
			})
			return nil, result, err
		})

	modelsdk.AddTool(s, &modelsdk.Tool{Name: "update_command", Description: "Update a channel command by UUID."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input updateCommandInput) (*modelsdk.CallToolResult, any, error) {
			id, err := parseID(input.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid command id: %w", err)
			}
			responses := make([]commandsrelations.UpdateInputResponse, 0, len(input.Responses))
			for _, response := range input.Responses {
				responses = append(responses, commandsrelations.UpdateInputResponse{Text: &response.Text, Order: response.Order, OnlineOnly: response.OnlineOnly, OfflineOnly: response.OfflineOnly})
			}
			result, err := h.deps.CommandsRelations.Update(ctx, id, commandsrelations.UpdateInput{
				ChannelID: requestScope.Channel.ID.String(), ActorID: requestScope.ActorID,
				Name: input.Name, Responses: responses, Aliases: input.Aliases, Description: input.Description,
				Cooldown: input.Cooldown, CooldownType: input.CooldownType, Enabled: input.Enabled,
				Visible: input.Visible, OnlineOnly: input.OnlineOnly, KeepResponsesOrder: input.KeepResponsesOrder,
			})
			return nil, result, err
		})

	modelsdk.AddTool(s, &modelsdk.Tool{Name: "delete_command", Description: "Delete a custom command from this channel."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
			id, err := parseID(input.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid command id: %w", err)
			}
			err = h.deps.Commands.Delete(ctx, commands.DeleteInput{ID: id, ChannelID: requestScope.Channel.ID.String(), ActorID: requestScope.ActorID})
			return nil, map[string]bool{"deleted": err == nil}, err
		})
}
