package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func ChatMessageToGQL(input entity.ChatMessage) gqlmodel.ChatMessage {
	return gqlmodel.ChatMessage{
		ID:              input.ID,
		ChannelID:       input.ChannelID,
		UserID:          input.UserID,
		UserName:        input.UserName,
		UserDisplayName: input.UserDisplayName,
		UserColor:       input.UserColor,
		Text:            input.Text,
		CreatedAt:       input.CreatedAt,
		UpdatedAt:       input.UpdatedAt,
	}
}
