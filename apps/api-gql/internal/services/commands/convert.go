package commands

import (
	"github.com/samber/lo"
	commandwithrelationentity "github.com/twirapp/twir/libs/entities/command_with_relations"
	"github.com/twirapp/twir/libs/repositories/commands/model"
)

func (c *Service) modelToEntity(m model.Command) commandwithrelationentity.Command {
	var expiresType *commandwithrelationentity.CommandExpireType
	if m.ExpiresType != nil {
		expiresType = lo.ToPtr(commandwithrelationentity.CommandExpireType(*m.ExpiresType))
	}

	return commandwithrelationentity.Command{
		ID:                        m.ID,
		Name:                      m.Name,
		Cooldown:                  m.Cooldown,
		CooldownType:              m.CooldownType,
		Enabled:                   m.Enabled,
		Aliases:                   m.Aliases,
		Description:               m.Description,
		Visible:                   m.Visible,
		ChannelID:                 m.ChannelID,
		Default:                   m.Default,
		DefaultName:               m.DefaultName,
		Module:                    m.Module,
		IsReply:                   m.IsReply,
		KeepResponsesOrder:        m.KeepResponsesOrder,
		DeniedUsersIDS:            m.DeniedUsersIDS,
		AllowedUsersIDS:           m.AllowedUsersIDS,
		RolesIDS:                  m.RolesIDS,
		OnlineOnly:                m.OnlineOnly,
		OfflineOnly:               m.OfflineOnly,
		EnabledCategories:         m.EnabledCategories,
		RequiredWatchTime:         m.RequiredWatchTime,
		RequiredMessages:          m.RequiredMessages,
		RequiredUsedChannelPoints: m.RequiredUsedChannelPoints,
		GroupID:                   m.GroupID,
		ExpiresAt:                 m.ExpiresAt,
		ExpiresType:               expiresType,
	}
}
