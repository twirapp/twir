package mappers

import (
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

var commandsExpiresAtMap = map[model.ChannelCommandExpiresType]gqlmodel.CommandExpiresType{
	model.ChannelCommandExpiresTypeDelete:  gqlmodel.CommandExpiresTypeDelete,
	model.ChannelCommandExpiresTypeDisable: gqlmodel.CommandExpiresTypeDisable,
}

func CommandsExpiresAtDbToGql(in model.ChannelCommandExpiresType) gqlmodel.CommandExpiresType {
	return commandsExpiresAtMap[in]
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
	m := gqlmodel.Command{
		ID:                        e.ID.String(),
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
		RolesIds:                  e.RolesIDS,
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
		ID:                  e.ID.String(),
		CommandID:           e.CommandID.String(),
		Text:                "", // will be set later
		TwitchCategoriesIds: e.TwitchCategoryIDs,
	}

	if e.Text != nil {
		m.Text = *e.Text
	}

	return m
}
