package mappers

import (
	"github.com/google/uuid"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands"
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

var commandsEntityExpiresAtMap = map[entity.CommandExpireType]gqlmodel.CommandExpiresType{
	entity.CommandExpireTypeDelete:  gqlmodel.CommandExpiresTypeDelete,
	entity.CommandExpireTypeDisable: gqlmodel.CommandExpiresTypeDisable,
}

func CommandEntityTo(e entity.Command) gqlmodel.Command {
	rolesIds := make([]string, len(e.RolesIDS))
	for i, v := range e.RolesIDS {
		rolesIds[i] = v.String()
	}

	m := gqlmodel.Command{
		ID:                        e.ID,
		Name:                      e.Name,
		Description:               "", // will be set later
		Aliases:                   e.Aliases,
		Responses:                 nil, // will be set later
		Cooldown:                  0,   // will be set later
		CooldownType:              e.CooldownType,
		Enabled:                   e.Enabled,
		Visible:                   e.Visible,
		Default:                   e.Default,
		DefaultName:               e.DefaultName,
		Module:                    e.Module,
		IsReply:                   e.IsReply,
		KeepResponsesOrder:        e.KeepResponsesOrder,
		DeniedUsersIds:            e.DeniedUsersIDS,
		AllowedUsersIds:           e.AllowedUsersIDS,
		RolesIds:                  rolesIds,
		OnlineOnly:                e.OnlineOnly,
		CooldownRolesIds:          e.CooldownRolesIDs,
		EnabledCategories:         e.EnabledCategories,
		RequiredWatchTime:         e.RequiredWatchTime,
		RequiredMessages:          e.RequiredMessages,
		RequiredUsedChannelPoints: e.RequiredUsedChannelPoints,
		GroupID:                   nil, // will be set later
		Group:                     nil, // will be set later
		ExpiresAt:                 nil, // will be set later
		ExpiresType:               nil, // will be set later
	}

	if e.Cooldown != nil {
		m.Cooldown = *e.Cooldown
	}

	if e.Description != nil {
		m.Description = *e.Description
	}

	if e.ExpiresAt != nil {
		expires := int(e.ExpiresAt.UnixMilli())
		m.ExpiresAt = &expires
	}

	if e.ExpiresType != nil {
		expiresType := commandsEntityExpiresAtMap[*e.ExpiresType]
		m.ExpiresType = &expiresType
	}

	if e.GroupID != nil {
		id := e.GroupID.String()
		m.GroupID = &id
	}

	return m
}

func CommandGroupTo(e entity.CommandGroup) gqlmodel.CommandGroup {
	return gqlmodel.CommandGroup{
		ID:    e.ID.String(),
		Name:  e.Name,
		Color: e.Color,
	}
}

func CommandResponseTo(e entity.CommandResponse) gqlmodel.CommandResponse {
	m := gqlmodel.CommandResponse{
		ID:                  e.ID,
		CommandID:           e.CommandID.String(),
		Text:                "", // will be set later
		TwitchCategoriesIds: e.TwitchCategoryIDs,
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
		CooldownRolesIDs:          input.CooldownRolesIds,
		EnabledCategories:         input.EnabledCategories,
		RequiredWatchTime:         input.RequiredWatchTime,
		RequiredMessages:          input.RequiredMessages,
		RequiredUsedChannelPoints: input.RequiredUsedChannelPoints,
		GroupID:                   groupId,
		ExpiresAt:                 input.ExpiresAt.Value(),
		ExpiresType:               expiresType,
		Responses:                 responses,
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
