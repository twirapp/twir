package chat_messages

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/chat_messages"
)

type GetManyInput struct {
	Page    int
	PerPage int

	ChannelID    *string
	UserNameLike *string
	TextLike     *string
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
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get chat messages: %w", err)
	}

	convertedMessages := make([]entity.ChatMessage, 0, len(messages))
	for _, m := range messages {
		convertedMessages = append(convertedMessages, c.modelToGql(m))
	}

	return convertedMessages, nil
}
