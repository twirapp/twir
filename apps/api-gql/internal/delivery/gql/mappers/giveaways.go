package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func GiveawayEntityTo(e entity.ChannelGiveaway) gqlmodel.ChannelGiveaway {
	return gqlmodel.ChannelGiveaway{
		ID:              e.ID.String(),
		ChannelID:       e.ChannelID,
		CreatedAt:       e.CreatedAt,
		UpdatedAt:       e.UpdatedAt,
		StartedAt:       e.StartedAt,
		EndedAt:         e.EndedAt,
		StoppedAt:       e.StoppedAt,
		Keyword:         e.Keyword,
		ArchivedAt:      e.ArchivedAt,
		CreatedByUserID: e.CreatedByUserID,
	}
}

func GiveawayParticipantEntityTo(
	e entity.ChannelGiveawayParticipant,
) gqlmodel.ChannelGiveawayParticipants {
	return gqlmodel.ChannelGiveawayParticipants{
		DisplayName: e.DisplayName,
		UserID:      e.UserID,
		IsWinner:    e.IsWinner,
		ID:          e.ID.String(),
		GiveawayID:  e.GiveawayID.String(),
	}
}

func GiveawayWinnerEntityTo(
	e entity.ChannelGiveawayWinner,
) gqlmodel.ChannelGiveawayWinner {
	return gqlmodel.ChannelGiveawayWinner{
		DisplayName: e.DisplayName,
		UserID:      e.UserID,
	}
}
