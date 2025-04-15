package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func ChatMessageToGQL(input entity.ChatMessage) gqlmodel.ChatMessage {
	return gqlmodel.ChatMessage{
		ID:              input.ID,
		ChannelID:       input.ChannelID,
		ChannelName:     input.ChannelName,
		ChannelLogin:    input.ChannelLogin,
		UserID:          input.UserID,
		UserName:        input.UserName,
		UserDisplayName: input.UserDisplayName,
		UserColor:       input.UserColor,
		Text:            input.Text,
		CreatedAt:       input.CreatedAt,
		UpdatedAt:       input.UpdatedAt,
	}
}
