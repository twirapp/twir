package chat_messages

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/chat_messages"
)

type GetManyInput struct {
	ChannelID    *string
	UserNameLike *string
	TextLike     *string
	Page         int
	PerPage      int
	UserIDs      []string
}

func (c *Service) GetMany(ctx context.Context, input GetManyInput) ([]entity.ChatMessage, error) {
	messages, err := c.chatMessagesRepository.GetMany(
		ctx,
		chat_messages.GetManyInput{
			Page:         input.Page,
			PerPage:      input.PerPage,
			ChannelID:    input.ChannelID,
			UserNameLike: input.UserNameLike,
			TextLike:     input.TextLike,
			UserIDs:      input.UserIDs,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat messages: %w", err)
	}

	convertedMessages := make([]entity.ChatMessage, len(messages))
	for i, message := range messages {
		convertedMessages[i] = c.modelToGql(message)
	}

	return convertedMessages, nil
}
