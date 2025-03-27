package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func ScheduledVipToGql(m entity.ScheduledVip) gqlmodel.ScheduledVip {
	return gqlmodel.ScheduledVip{
		ID:        m.ID.String(),
		UserID:    m.UserID,
		ChannelID: m.ChannelID,
		CreatedAt: m.CreatedAt,
		RemoveAt:  m.RemoveAt,
	}
}
