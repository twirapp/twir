package mappers

import (
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands"
	commandwithrelationentity "github.com/twirapp/twir/libs/entities/command_with_relations"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/integrations/streamelements"
)

var commandsExpiresAtMap = map[model.ChannelCommandExpiresType]gqlmodel.CommandExpiresType{
	model.ChannelCommandExpiresTypeDelete:  gqlmodel.CommandExpiresTypeDelete,
	model.ChannelCommandExpiresTypeDisable: gqlmodel.CommandExpiresTypeDisable,
}

func CommandsExpiresAtGqlToDb(in gqlmodel.CommandExpiresType) model.ChannelCommandExpiresType {
	for k, v := range commandsExpiresAtMap {
		if v == in {
			return k
		}
	}

	return model.ChannelCommandExpiresTypeDelete
}

var commandsEntityExpiresAtMap = map[commandwithrelationentity.CommandExpireType]gqlmodel.CommandExpiresType{
	commandwithrelationentity.CommandExpireTypeDelete:  gqlmodel.CommandExpiresTypeDelete,
	commandwithrelationentity.CommandExpireTypeDisable: gqlmodel.CommandExpiresTypeDisable,
}

func CommandEntityTo(e commandwithrelationentity.CommandWithGroupAndResponses) gqlmodel.Command {
	rolesIds := make([]string, len(e.Command.RolesIDS))
	for i, v := range e.Command.RolesIDS {
		rolesIds[i] = v.String()
	}

	m := gqlmodel.Command{
		ID:                        e.Command.ID,
		Name:                      e.Command.Name,
		Description:               "", // will be set later
		Aliases:                   e.Command.Aliases,
		Responses:                 nil, // will be set later
		Cooldown:                  0,   // will be set later
		CooldownType:              e.Command.CooldownType,
		Enabled:                   e.Command.Enabled,
		Visible:                   e.Command.Visible,
		Default:                   e.Command.Default,
		DefaultName:               e.Command.DefaultName,
		Module:                    e.Command.Module,
		IsReply:                   e.Command.IsReply,
		KeepResponsesOrder:        e.Command.KeepResponsesOrder,
		DeniedUsersIds:            e.Command.DeniedUsersIDS,
		AllowedUsersIds:           e.Command.AllowedUsersIDS,
		RolesIds:                  rolesIds,
		OnlineOnly:                e.Command.OnlineOnly,
		OfflineOnly:               e.Command.OfflineOnly,
		EnabledCategories:         e.Command.EnabledCategories,
		RequiredWatchTime:         e.Command.RequiredWatchTime,
		RequiredMessages:          e.Command.RequiredMessages,
		RequiredUsedChannelPoints: e.Command.RequiredUsedChannelPoints,
		GroupID:                   nil, // will be set later
		Group:                     nil, // will be set later
		ExpiresAt:                 nil, // will be set later
		ExpiresType:               nil, // will be set later
		RoleCooldowns:             nil, // will be set later
	}

	if e.Command.Cooldown != nil {
		m.Cooldown = *e.Command.Cooldown
	}

	if e.Command.Description != nil {
		m.Description = *e.Command.Description
	}

	if e.Command.ExpiresAt != nil {
		expires := int(e.Command.ExpiresAt.UnixMilli())
		m.ExpiresAt = &expires
	}

	if e.Command.ExpiresType != nil {
		expiresType := commandsEntityExpiresAtMap[*e.Command.ExpiresType]
		m.ExpiresType = &expiresType
	}

	if e.Command.GroupID != nil {
		id := e.Command.GroupID.String()
		m.GroupID = &id
	}

	rolesCooldowns := make([]gqlmodel.CommandRoleCooldown, 0, len(e.RolesCooldowns))
	for _, rc := range e.RolesCooldowns {
		rolesCooldowns = append(
			rolesCooldowns, gqlmodel.CommandRoleCooldown{
				ID:        rc.ID,
				CommandID: rc.CommandID,
				RoleID:    rc.RoleID,
				Cooldown:  rc.Cooldown,
			},
		)
	}
	m.RoleCooldowns = rolesCooldowns

	return m
}

func CommandGroupTo(e commandwithrelationentity.CommandGroup) gqlmodel.CommandGroup {
	return gqlmodel.CommandGroup{
		ID:    e.ID.String(),
		Name:  e.Name,
		Color: e.Color,
	}
}

func CommandResponseTo(e commandwithrelationentity.CommandResponse) gqlmodel.CommandResponse {
	m := gqlmodel.CommandResponse{
		ID:                  e.ID,
		CommandID:           e.CommandID.String(),
		Text:                "", // will be set later
		TwitchCategoriesIds: e.TwitchCategoryIDs,
		OnlineOnly:          e.OnlineOnly,
		OfflineOnly:         e.OfflineOnly,
	}

	if e.Text != nil {
		m.Text = *e.Text
	}

	return m
}

func CommandGqlInputToService(
	channelID, actorID string,
	input gqlmodel.CommandsCreateOpts,
) commands.CreateInput {
	responses := make([]commands.CreateInputResponse, len(input.Responses))
	for idx, res := range input.Responses {
		responses[idx] = commands.CreateInputResponse{
			Text:              &res.Text,
			Order:             idx,
			TwitchCategoryIDs: res.TwitchCategoriesIds,
			OnlineOnly:        res.OnlineOnly,
			OfflineOnly:       res.OfflineOnly,
		}
	}

	var groupId *uuid.UUID
	if input.GroupID.IsSet() && input.GroupID.Value() != nil {
		parsedGroupId, err := uuid.Parse(*input.GroupID.Value())
		if err == nil {
			groupId = &parsedGroupId
		}
	}

	var expiresType *string
	if input.ExpiresType.IsSet() && input.ExpiresType.Value() != nil {
		expiresType = lo.ToPtr(input.ExpiresType.Value().String())
	}

	roleCooldowns := make([]commands.CreateInputRoleCooldown, 0, len(input.RoleCooldowns))
	for _, rc := range input.RoleCooldowns {
		roleCooldowns = append(
			roleCooldowns, commands.CreateInputRoleCooldown{
				RoleID:   rc.RoleID.String(),
				Cooldown: rc.Cooldown,
			},
		)
	}

	return commands.CreateInput{
		ChannelID:                 channelID,
		ActorID:                   actorID,
		Name:                      input.Name,
		Cooldown:                  input.Cooldown,
		CooldownType:              input.CooldownType,
		Enabled:                   input.Enabled,
		Aliases:                   input.Aliases,
		Description:               input.Description,
		Visible:                   input.Visible,
		IsReply:                   input.IsReply,
		KeepResponsesOrder:        input.KeepResponsesOrder,
		DeniedUsersIDS:            input.DeniedUsersIds,
		AllowedUsersIDS:           input.AllowedUsersIds,
		RolesIDS:                  input.RolesIds,
		OnlineOnly:                input.OnlineOnly,
		EnabledCategories:         input.EnabledCategories,
		RequiredWatchTime:         input.RequiredWatchTime,
		RequiredMessages:          input.RequiredMessages,
		RequiredUsedChannelPoints: input.RequiredUsedChannelPoints,
		GroupID:                   groupId,
		ExpiresAt:                 input.ExpiresAt.Value(),
		ExpiresType:               expiresType,
		Responses:                 responses,
		RoleCooldowns:             roleCooldowns,
	}
}

func StreamElementsCommandToGql(m streamelements.Command) gqlmodel.StreamElementsCommand {
	return gqlmodel.StreamElementsCommand{
		ID:             m.ID,
		Name:           m.Name,
		Enabled:        m.Enabled,
		Cooldown:       m.Cooldown.Global,
		Aliases:        m.Aliases,
		Response:       m.Response,
		AccessLevel:    m.AccessLevel,
		EnabledOnline:  m.EnabledOnline,
		EnabledOffline: m.EnabledOffline,
		Hidden:         m.Hidden,
		Type:           m.Type,
		CreatedAt:      m.CreatedAt,
		UpdatedAt:      m.UpdatedAt,
	}
}
