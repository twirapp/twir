package commands

import (
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/commands/model"
)

func (c *Service) modelToEntity(m model.Command) entity.Command {
	var expiresType *entity.CommandExpireType
	if m.ExpiresType != nil {
		expiresType = lo.ToPtr(entity.CommandExpireType(*m.ExpiresType))
	}

	return entity.Command{
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
		CooldownRolesIDs:          m.CooldownRolesIDs,
		EnabledCategories:         m.EnabledCategories,
		RequiredWatchTime:         m.RequiredWatchTime,
		RequiredMessages:          m.RequiredMessages,
		RequiredUsedChannelPoints: m.RequiredUsedChannelPoints,
		GroupID:                   m.GroupID,
		ExpiresAt:                 m.ExpiresAt,
		ExpiresType:               expiresType,
	}
}
