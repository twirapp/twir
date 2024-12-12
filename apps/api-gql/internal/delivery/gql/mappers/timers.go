package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func TimerEntityToGql(m entity.Timer) gqlmodel.Timer {
	responses := make([]gqlmodel.TimerResponse, 0, len(m.Responses))
	for _, r := range m.Responses {
		responses = append(
			responses,
			gqlmodel.TimerResponse{
				ID:         r.ID.String(),
				Text:       r.Text,
				IsAnnounce: r.IsAnnounce,
			},
		)
	}

	return gqlmodel.Timer{
		ID:              m.ID.String(),
		Name:            m.Name,
		Enabled:         m.Enabled,
		TimeInterval:    m.TimeInterval,
		MessageInterval: m.MessageInterval,
		Responses:       responses,
	}
}
