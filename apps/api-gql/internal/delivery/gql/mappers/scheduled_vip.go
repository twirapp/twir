package mappers

import (
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	scheduledvipsentity "github.com/twirapp/twir/libs/entities/scheduled_vips"
)

var scheduledVipTypeMap = map[scheduledvipsentity.RemoveType]gqlmodel.ScheduledVipRemoveType{
	scheduledvipsentity.RemoveTypeStreamEnd: gqlmodel.ScheduledVipRemoveTypeStreamEnd,
	scheduledvipsentity.RemoveTypeTime:      gqlmodel.ScheduledVipRemoveTypeTime,
}

func ScheduledVipToGql(e scheduledvipsentity.ScheduledVip) gqlmodel.ScheduledVip {
	m := gqlmodel.ScheduledVip{
		ID:        e.ID.String(),
		UserID:    e.UserID,
		ChannelID: e.ChannelID,
		CreatedAt: e.CreatedAt,
		RemoveAt:  e.RemoveAt,
	}

	if e.RemoveType != nil {
		m.RemoveType = lo.ToPtr(scheduledVipTypeMap[*e.RemoveType])
	}

	return m
}
