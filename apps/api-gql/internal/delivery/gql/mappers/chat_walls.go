package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func MapChatWallToGql(input entity.ChatWall) gqlmodel.ChatWall {
	return gqlmodel.ChatWall{
		ID:                     input.ID,
		Phrase:                 input.Phrase,
		Enabled:                input.Enabled,
		Action:                 gqlmodel.ChatWallAction(input.Action),
		DurationSeconds:        input.DurationSeconds,
		TimeoutDurationSeconds: input.TimeoutDurationSeconds,
		AffectedMessages:       input.AffectedMessages,
		CreatedAt:              input.CreatedAt,
		UpdatedAt:              input.UpdatedAt,
	}
}

func MapChatWallLogToGql(input entity.ChatWallLog) gqlmodel.ChatWallLog {
	return gqlmodel.ChatWallLog{
		ID:        input.ID,
		UserID:    input.UserID,
		CreatedAt: input.CreatedAt,
		Text:      input.Text,
	}
}
