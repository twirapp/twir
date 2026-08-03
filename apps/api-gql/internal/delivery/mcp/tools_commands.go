package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	modelsdk "github.com/modelcontextprotocol/go-sdk/mcp"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands"
	commandsrelations "github.com/twirapp/twir/apps/api-gql/internal/services/commands_with_groups_and_responses"
	"github.com/twirapp/twir/libs/entities/platform"
)

type commandResponseInput struct {
	Text                string              `json:"text" jsonschema:"command response text"`
	Order               int                 `json:"order,omitempty" jsonschema:"response order, starting at zero"`
	TwitchCategoriesIDS []string            `json:"twitchCategoriesIds,omitempty" jsonschema:"Twitch category IDs where this response is available"`
	OnlineOnly          bool                `json:"onlineOnly,omitempty"`
	OfflineOnly         bool                `json:"offlineOnly,omitempty"`
	Platforms           []platform.Platform `json:"platforms,omitempty" jsonschema:"platforms where this response is available: twitch, kick, vk_video_live, youtube; empty means all"`
}

type commandRoleCooldownInput struct {
	RoleID   string `json:"roleId" jsonschema:"channel role UUID returned by list_command_roles"`
	Cooldown int    `json:"cooldown" jsonschema:"role-specific cooldown in seconds"`
}

type createCommandInput struct {
	Name                      string                     `json:"name" jsonschema:"command name without the prefix"`
	Responses                 []commandResponseInput     `json:"responses" jsonschema:"one or more command responses"`
	Aliases                   []string                   `json:"aliases,omitempty"`
	Description               string                     `json:"description,omitempty"`
	Cooldown                  int                        `json:"cooldown,omitempty" jsonschema:"cooldown in seconds"`
	CooldownType              string                     `json:"cooldownType,omitempty" jsonschema:"GLOBAL or PER_USER"`
	Enabled                   bool                       `json:"enabled"`
	Visible                   bool                       `json:"visible"`
	IsReply                   bool                       `json:"isReply,omitempty"`
	KeepResponsesOrder        bool                       `json:"keepResponsesOrder,omitempty"`
	DeniedUsersIDS            []string                   `json:"deniedUsersIds,omitempty" jsonschema:"internal Twir user UUIDs denied access"`
	AllowedUsersIDS           []string                   `json:"allowedUsersIds,omitempty" jsonschema:"internal Twir user UUIDs explicitly allowed access"`
	RolesIDS                  []string                   `json:"rolesIds,omitempty" jsonschema:"channel role UUIDs allowed access; use list_command_roles to discover them"`
	OnlineOnly                bool                       `json:"onlineOnly,omitempty"`
	OfflineOnly               bool                       `json:"offlineOnly,omitempty"`
	EnabledCategories         []string                   `json:"enabledCategories,omitempty" jsonschema:"Twitch category IDs where the command is enabled"`
	RequiredWatchTime         int                        `json:"requiredWatchTime,omitempty" jsonschema:"required watch time in minutes"`
	RequiredMessages          int                        `json:"requiredMessages,omitempty" jsonschema:"required channel message count"`
	RequiredUsedChannelPoints int                        `json:"requiredUsedChannelPoints,omitempty" jsonschema:"required channel points spent"`
	GroupID                   *string                    `json:"groupId,omitempty" jsonschema:"command group UUID"`
	ExpiresAt                 *int                       `json:"expiresAt,omitempty" jsonschema:"expiration time as Unix milliseconds"`
	ExpiresType               *string                    `json:"expiresType,omitempty" jsonschema:"expiration action: DISABLE or DELETE"`
	RoleCooldowns             []commandRoleCooldownInput `json:"roleCooldowns,omitempty"`
	Platforms                 []platform.Platform        `json:"platforms,omitempty" jsonschema:"platforms where the command is available: twitch, kick, vk_video_live, youtube; empty means all"`
}

type updateCommandInput struct {
	ID                        string                     `json:"id" jsonschema:"command UUID"`
	Name                      *string                    `json:"name,omitempty"`
	Responses                 []commandResponseInput     `json:"responses,omitempty"`
	Aliases                   []string                   `json:"aliases,omitempty"`
	Description               *string                    `json:"description,omitempty"`
	Cooldown                  *int                       `json:"cooldown,omitempty"`
	CooldownType              *string                    `json:"cooldownType,omitempty" jsonschema:"GLOBAL or PER_USER"`
	Enabled                   *bool                      `json:"enabled,omitempty"`
	Visible                   *bool                      `json:"visible,omitempty"`
	IsReply                   *bool                      `json:"isReply,omitempty"`
	KeepResponsesOrder        *bool                      `json:"keepResponsesOrder,omitempty"`
	DeniedUsersIDS            []string                   `json:"deniedUsersIds,omitempty"`
	AllowedUsersIDS           []string                   `json:"allowedUsersIds,omitempty"`
	RolesIDS                  []string                   `json:"rolesIds,omitempty" jsonschema:"channel role UUIDs allowed access; use list_command_roles to discover them"`
	OnlineOnly                *bool                      `json:"onlineOnly,omitempty"`
	OfflineOnly               *bool                      `json:"offlineOnly,omitempty"`
	EnabledCategories         []string                   `json:"enabledCategories,omitempty"`
	RequiredWatchTime         *int                       `json:"requiredWatchTime,omitempty"`
	RequiredMessages          *int                       `json:"requiredMessages,omitempty"`
	RequiredUsedChannelPoints *int                       `json:"requiredUsedChannelPoints,omitempty"`
	GroupID                   *string                    `json:"groupId,omitempty" jsonschema:"command group UUID"`
	ExpiresAt                 *int                       `json:"expiresAt,omitempty" jsonschema:"expiration time as Unix milliseconds"`
	ExpiresType               *string                    `json:"expiresType,omitempty" jsonschema:"expiration action: DISABLE or DELETE"`
	RoleCooldowns             []commandRoleCooldownInput `json:"roleCooldowns,omitempty"`
	Platforms                 []platform.Platform        `json:"platforms,omitempty" jsonschema:"platforms where the command is available: twitch, kick, vk_video_live, youtube; empty means all"`
}

type idInput struct {
	ID string `json:"id" jsonschema:"object UUID"`
}

func (h *Handler) addCommandTools(s *modelsdk.Server, requestScope scope) {
	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_command_roles", Description: "List channel roles and their UUIDs for command access and role cooldown settings."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			items, err := h.deps.Roles.GetManyByChannelID(ctx, requestScope.Channel.ID.String())
			return nil, map[string]any{"roles": items}, err
		})

	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_command_groups", Description: "List command groups for this channel."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			items, err := h.deps.CommandGroups.GetManyByChannelID(ctx, requestScope.Channel.ID.String())
			return nil, map[string]any{"groups": items}, err
		})

	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "list_commands", Description: "List all commands, groups, responses, cooldowns, and permissions for this channel."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, _ struct{}) (*modelsdk.CallToolResult, any, error) {
			items, err := h.deps.CommandsRelations.GetManyByChannelID(ctx, requestScope.Channel.ID.String())
			return nil, map[string]any{"commands": items}, err
		})

	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "get_command", Description: "Get one channel command by UUID."},
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

	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "create_command", Description: "Create a command for this channel with dashboard-equivalent access, availability, response, cooldown, grouping, expiration, and platform settings."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input createCommandInput) (*modelsdk.CallToolResult, any, error) {
			serviceInput, err := createCommandServiceInput(requestScope, input)
			if err != nil {
				return nil, nil, err
			}
			result, err := h.deps.Commands.Create(ctx, serviceInput)
			return nil, result, err
		})

	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "update_command", Description: "Update a channel command by UUID with dashboard-equivalent settings."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input updateCommandInput) (*modelsdk.CallToolResult, any, error) {
			id, err := parseID(input.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid command id: %w", err)
			}
			serviceInput, err := updateCommandServiceInput(requestScope, input)
			if err != nil {
				return nil, nil, err
			}
			result, err := h.deps.CommandsRelations.Update(ctx, id, serviceInput)
			return nil, result, err
		})

	addTool(newToolRegistrar(s, requestScope.AccessScopes), &modelsdk.Tool{Name: "delete_command", Description: "Delete a custom command from this channel."},
		func(ctx context.Context, _ *modelsdk.CallToolRequest, input idInput) (*modelsdk.CallToolResult, any, error) {
			id, err := parseID(input.ID)
			if err != nil {
				return nil, nil, fmt.Errorf("invalid command id: %w", err)
			}
			err = h.deps.Commands.Delete(ctx, commands.DeleteInput{ID: id, ChannelID: requestScope.Channel.ID.String(), ActorID: requestScope.ActorID})
			return nil, map[string]bool{"deleted": err == nil}, err
		})
}

func createCommandServiceInput(requestScope scope, input createCommandInput) (commands.CreateInput, error) {
	groupID, err := parseOptionalCommandID(input.GroupID, "group id")
	if err != nil {
		return commands.CreateInput{}, err
	}
	if err := validateCommandExpiration(input.ExpiresAt, input.ExpiresType, true); err != nil {
		return commands.CreateInput{}, err
	}
	if err := validateCommandPlatforms(input.Platforms, input.Responses); err != nil {
		return commands.CreateInput{}, err
	}
	if err := validateCommandRoleIDs(input.RolesIDS, input.RoleCooldowns); err != nil {
		return commands.CreateInput{}, err
	}

	responses := make([]commands.CreateInputResponse, 0, len(input.Responses))
	for _, response := range input.Responses {
		text := response.Text
		responses = append(responses, commands.CreateInputResponse{
			Text:              &text,
			Order:             response.Order,
			TwitchCategoryIDs: response.TwitchCategoriesIDS,
			OnlineOnly:        response.OnlineOnly,
			OfflineOnly:       response.OfflineOnly,
			Platforms:         response.Platforms,
		})
	}
	roleCooldowns := make([]commands.CreateInputRoleCooldown, 0, len(input.RoleCooldowns))
	for _, roleCooldown := range input.RoleCooldowns {
		roleCooldowns = append(roleCooldowns, commands.CreateInputRoleCooldown{RoleID: roleCooldown.RoleID, Cooldown: roleCooldown.Cooldown})
	}

	return commands.CreateInput{
		ChannelID:                 requestScope.Channel.ID.String(),
		ActorID:                   requestScope.ActorID,
		Name:                      input.Name,
		Cooldown:                  input.Cooldown,
		CooldownType:              input.CooldownType,
		Enabled:                   input.Enabled,
		Aliases:                   input.Aliases,
		Description:               input.Description,
		Visible:                   input.Visible,
		IsReply:                   input.IsReply,
		KeepResponsesOrder:        input.KeepResponsesOrder,
		DeniedUsersIDS:            input.DeniedUsersIDS,
		AllowedUsersIDS:           input.AllowedUsersIDS,
		RolesIDS:                  input.RolesIDS,
		OnlineOnly:                input.OnlineOnly,
		OfflineOnly:               input.OfflineOnly,
		EnabledCategories:         append([]string{}, input.EnabledCategories...),
		RequiredWatchTime:         input.RequiredWatchTime,
		RequiredMessages:          input.RequiredMessages,
		RequiredUsedChannelPoints: input.RequiredUsedChannelPoints,
		GroupID:                   groupID,
		ExpiresAt:                 input.ExpiresAt,
		ExpiresType:               input.ExpiresType,
		Responses:                 responses,
		RoleCooldowns:             roleCooldowns,
		Platforms:                 input.Platforms,
	}, nil
}

func updateCommandServiceInput(requestScope scope, input updateCommandInput) (commandsrelations.UpdateInput, error) {
	groupID, err := parseOptionalCommandID(input.GroupID, "group id")
	if err != nil {
		return commandsrelations.UpdateInput{}, err
	}
	if err := validateCommandExpiration(input.ExpiresAt, input.ExpiresType, false); err != nil {
		return commandsrelations.UpdateInput{}, err
	}
	if err := validateCommandPlatforms(input.Platforms, input.Responses); err != nil {
		return commandsrelations.UpdateInput{}, err
	}
	if err := validateCommandRoleIDs(input.RolesIDS, input.RoleCooldowns); err != nil {
		return commandsrelations.UpdateInput{}, err
	}

	var responses []commandsrelations.UpdateInputResponse
	if input.Responses != nil {
		responses = make([]commandsrelations.UpdateInputResponse, 0, len(input.Responses))
		for _, response := range input.Responses {
			text := response.Text
			responses = append(responses, commandsrelations.UpdateInputResponse{
				Text:              &text,
				Order:             response.Order,
				TwitchCategoryIDs: response.TwitchCategoriesIDS,
				OnlineOnly:        response.OnlineOnly,
				OfflineOnly:       response.OfflineOnly,
				Platforms:         response.Platforms,
			})
		}
	}
	var roleCooldowns []commandsrelations.UpdateInputRoleCooldown
	if input.RoleCooldowns != nil {
		roleCooldowns = make([]commandsrelations.UpdateInputRoleCooldown, 0, len(input.RoleCooldowns))
		for _, roleCooldown := range input.RoleCooldowns {
			roleCooldowns = append(roleCooldowns, commandsrelations.UpdateInputRoleCooldown{RoleID: roleCooldown.RoleID, Cooldown: roleCooldown.Cooldown})
		}
	}
	var expiresAt *time.Time
	if input.ExpiresAt != nil {
		expiresAt = new(time.Time)
		*expiresAt = time.UnixMilli(int64(*input.ExpiresAt))
	}

	return commandsrelations.UpdateInput{
		ChannelID:                 requestScope.Channel.ID.String(),
		ActorID:                   requestScope.ActorID,
		Name:                      input.Name,
		Cooldown:                  input.Cooldown,
		CooldownType:              input.CooldownType,
		Enabled:                   input.Enabled,
		Aliases:                   input.Aliases,
		Description:               input.Description,
		Visible:                   input.Visible,
		IsReply:                   input.IsReply,
		KeepResponsesOrder:        input.KeepResponsesOrder,
		DeniedUsersIDS:            input.DeniedUsersIDS,
		AllowedUsersIDS:           input.AllowedUsersIDS,
		RolesIDS:                  input.RolesIDS,
		OnlineOnly:                input.OnlineOnly,
		OfflineOnly:               input.OfflineOnly,
		EnabledCategories:         input.EnabledCategories,
		RequiredWatchTime:         input.RequiredWatchTime,
		RequiredMessages:          input.RequiredMessages,
		RequiredUsedChannelPoints: input.RequiredUsedChannelPoints,
		GroupID:                   groupID,
		ExpiresAt:                 expiresAt,
		ExpiresType:               input.ExpiresType,
		Responses:                 responses,
		RoleCooldowns:             roleCooldowns,
		Platforms:                 input.Platforms,
	}, nil
}

func parseOptionalCommandID(value *string, field string) (*uuid.UUID, error) {
	if value == nil {
		return nil, nil
	}
	parsed, err := uuid.Parse(*value)
	if err != nil {
		return nil, fmt.Errorf("invalid %s: %w", field, err)
	}
	return &parsed, nil
}

func validateCommandExpiration(expiresAt *int, expiresType *string, requirePair bool) error {
	if requirePair && (expiresAt == nil) != (expiresType == nil) {
		return fmt.Errorf("expiresAt and expiresType must be set together")
	}
	if expiresType != nil && *expiresType != "DISABLE" && *expiresType != "DELETE" {
		return fmt.Errorf("invalid expiresType %q: expected DISABLE or DELETE", *expiresType)
	}
	return nil
}

func validateCommandPlatforms(commandPlatforms []platform.Platform, responses []commandResponseInput) error {
	for _, commandPlatform := range commandPlatforms {
		if !commandPlatform.IsValid() {
			return fmt.Errorf("invalid command platform %q", commandPlatform)
		}
	}
	for _, response := range responses {
		for _, responsePlatform := range response.Platforms {
			if !responsePlatform.IsValid() {
				return fmt.Errorf("invalid response platform %q", responsePlatform)
			}
		}
	}
	return nil
}

func validateCommandRoleIDs(roleIDs []string, roleCooldowns []commandRoleCooldownInput) error {
	for _, roleID := range roleIDs {
		if _, err := uuid.Parse(roleID); err != nil {
			return fmt.Errorf("invalid role id %q: %w", roleID, err)
		}
	}
	for _, roleCooldown := range roleCooldowns {
		if _, err := uuid.Parse(roleCooldown.RoleID); err != nil {
			return fmt.Errorf("invalid role cooldown id %q: %w", roleCooldown.RoleID, err)
		}
	}
	return nil
}
