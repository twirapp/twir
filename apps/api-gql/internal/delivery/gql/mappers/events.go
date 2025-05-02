package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func MapEventToGQL(event entity.Event) gqlmodel.Event {
	operations := make([]gqlmodel.EventOperation, 0, len(event.Operations))
	for _, op := range event.Operations {
		filters := make([]gqlmodel.EventOperationFilter, 0, len(op.Filters))
		for _, f := range op.Filters {
			filters = append(
				filters, gqlmodel.EventOperationFilter{
					ID:    f.ID,
					Type:  f.Type,
					Left:  f.Left,
					Right: f.Right,
				},
			)
		}

		operations = append(
			operations, gqlmodel.EventOperation{
				ID:             op.ID,
				Type:           op.Type,
				Input:          op.Input,
				Delay:          op.Delay,
				Repeat:         op.Repeat,
				UseAnnounce:    op.UseAnnounce,
				TimeoutTime:    op.TimeoutTime,
				TimeoutMessage: op.TimeoutMessage,
				Target:         op.Target,
				Enabled:        op.Enabled,
				Filters:        filters,
			},
		)
	}

	return gqlmodel.Event{
		ID:          event.ID,
		ChannelID:   event.ChannelID,
		Type:        event.Type,
		RewardID:    event.RewardID,
		CommandID:   event.CommandID,
		KeywordID:   event.KeywordID,
		Description: event.Description,
		Enabled:     event.Enabled,
		OnlineOnly:  event.OnlineOnly,
		Operations:  operations,
	}
}
