package chat_messages

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/repositories/chat_messages"
	"github.com/twirapp/twir/libs/repositories/chat_messages/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	ChatMessagesRepository chat_messages.Repository
	TwirBus                *buscore.Bus
}

const allMessagesSubKey = "ALL"

func New(opts Opts) *Service {
	s := &Service{
		chatMessagesRepository: opts.ChatMessagesRepository,
		subs:                   make(map[string]chan entity.ChatMessage),
	}

	s.subs[allMessagesSubKey] = make(chan entity.ChatMessage)

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return opts.TwirBus.ChatMessages.Subscribe(s.handleBusEvent)
			},
			OnStop: func(ctx context.Context) error {
				opts.TwirBus.ChatMessages.Unsubscribe()
				return nil
			},
		},
	)

	return s
}

type Service struct {
	chatMessagesRepository chat_messages.Repository

	subs map[string]chan entity.ChatMessage
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

func (c *Service) handleBusEvent(_ context.Context, data twitch.TwitchChatMessage) struct{} {
	textBuilder := strings.Builder{}
	for _, fragment := range data.Message.Fragments {
		textBuilder.WriteString(fragment.Text)
	}
	msg := entity.ChatMessage{
		ID:              uuid.New(),
		ChannelID:       data.BroadcasterUserId,
		ChannelLogin:    data.BroadcasterUserLogin,
		ChannelName:     data.BroadcasterUserName,
		UserID:          data.ChatterUserId,
		UserName:        data.ChatterUserLogin,
		UserDisplayName: data.ChatterUserName,
		UserColor:       data.Color,
		Text:            textBuilder.String(),
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	if ch, ok := c.subs[data.BroadcasterUserId]; ok {
		ch <- msg
	}

	if ch, ok := c.subs[allMessagesSubKey]; ok {
		ch <- msg
	}

	return struct{}{}
}

func (c *Service) SubscribeToNewMessagesByChannelID(channelID string) <-chan entity.ChatMessage {
	if _, ok := c.subs[channelID]; !ok {
		c.subs[channelID] = make(chan entity.ChatMessage)
	}

	return c.subs[channelID]
}

func (c *Service) SubscribeToNewMessages() <-chan entity.ChatMessage {
	return c.subs[allMessagesSubKey]
}

func (c *Service) UnsubscribeFromNewMessages(channelID string) {
	if ch, ok := c.subs[channelID]; ok {
		close(ch)
		delete(c.subs, channelID)
	}
}
