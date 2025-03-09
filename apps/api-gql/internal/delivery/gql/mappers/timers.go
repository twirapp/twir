package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/integrations/streamelements"
)

func TimerEntityToGql(m entity.Timer) gqlmodel.Timer {
	responses := make([]gqlmodel.TimerResponse, 0, len(m.Responses))
	for _, r := range m.Responses {
		responses = append(
			responses,
			gqlmodel.TimerResponse{
				ID:         r.ID,
				Text:       r.Text,
				IsAnnounce: r.IsAnnounce,
			},
		)
	}

	return gqlmodel.Timer{
		ID:              m.ID,
		Name:            m.Name,
		Enabled:         m.Enabled,
		TimeInterval:    m.TimeInterval,
		MessageInterval: m.MessageInterval,
		Responses:       responses,
	}
}

func StreamElementsTimerToGql(m streamelements.Timer) gqlmodel.StreamElementsTimer {
	return gqlmodel.StreamElementsTimer{
		ID:        m.Id,
		Name:      m.Name,
		Enabled:   m.Enabled,
		ChatLines: m.ChatLines,
		Message:   m.Message,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}
