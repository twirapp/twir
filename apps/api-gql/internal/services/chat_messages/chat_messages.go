package chat_messages

import (
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/repositories/chat_messages"
	"github.com/twirapp/twir/libs/repositories/chat_messages/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	ChatMessagesRepository chat_messages.Repository
}

func New(opts Opts) *Service {
	return &Service{
		chatMessagesRepository: opts.ChatMessagesRepository,
	}
}

type Service struct {
	chatMessagesRepository chat_messages.Repository
}

func (c *Service) modelToGql(m model.ChatMessage) entity.ChatMessage {
	return entity.ChatMessage{
		ID:              m.ID,
		ChannelID:       m.ChannelID,
		UserID:          m.UserID,
		UserName:        m.UserName,
		UserDisplayName: m.UserDisplayName,
		UserColor:       m.UserColor,
		Text:            m.Text,
		CreatedAt:       m.CreatedAt,
		UpdatedAt:       m.UpdatedAt,
	}
}
