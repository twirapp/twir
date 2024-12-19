package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func RolesToGql(m entity.ChannelRole) gqlmodel.Role {
	permissions := make([]gqlmodel.ChannelRolePermissionEnum, len(m.Permissions))
	for i, permission := range m.Permissions {
		permissions[i] = gqlmodel.ChannelRolePermissionEnum(permission)
	}

	return gqlmodel.Role{
		ID:          m.ID,
		ChannelID:   m.ChannelID,
		Name:        m.Name,
		Type:        gqlmodel.RoleTypeEnum(m.Type.String()),
		Permissions: permissions,
		Settings: &gqlmodel.RoleSettings{
			RequiredWatchTime:         int(m.RequiredWatchTime),
			RequiredMessages:          int(m.RequiredMessages),
			RequiredUserChannelPoints: int(m.RequiredUsedChannelPoints),
		},
	}
}
