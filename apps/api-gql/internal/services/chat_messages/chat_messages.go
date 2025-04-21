package chat_messages

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/wsrouter"
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
	WsRouter               wsrouter.WsRouter
	Logger                 logger.Logger
}

const (
	chatMessagesSubscriptionKey    = "api.chatMessages"
	chatMessagesSubscriptionKeyAll = chatMessagesSubscriptionKey + ".All"
)

func chatMessagesSubscriptionKeyCreate(channelId string) string {
	return chatMessagesSubscriptionKey + "." + channelId
}

func New(opts Opts) *Service {
	s := &Service{
		chatMessagesRepository: opts.ChatMessagesRepository,
		wsRouter:               opts.WsRouter,
		logger:                 opts.Logger,
	}

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

	wsRouter wsrouter.WsRouter
	logger   logger.Logger
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

	err := c.wsRouter.Publish(chatMessagesSubscriptionKeyCreate(data.BroadcasterUserId), msg)
	if err != nil {
		c.logger.Error(
			"Cannot publish some message to separate broadcaster messages",
			slog.Any("err", err),
		)
	}

	err = c.wsRouter.Publish(chatMessagesSubscriptionKeyAll, msg)
	if err != nil {
		c.logger.Error("Cannot publish some message to all messages", slog.Any("err", err))
	}

	return struct{}{}
}

func (c *Service) SubscribeToNewMessagesByChannelID(
	ctx context.Context,
	channelID string,
) <-chan entity.ChatMessage {
	channel := make(chan entity.ChatMessage)
	go func() {
		sub, err := c.wsRouter.Subscribe(
			[]string{
				chatMessagesSubscriptionKeyCreate(channelID),
			},
		)
		if err != nil {
			panic(err)
		}
		defer func() {
			sub.Unsubscribe()
			close(channel)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case data := <-sub.GetChannel():
				var msg entity.ChatMessage
				if err := json.Unmarshal(data, &msg); err != nil {
					panic(err)
				}

				channel <- msg
			}
		}
	}()

	return channel
}

func (c *Service) SubscribeToNewMessages(ctx context.Context) <-chan entity.ChatMessage {
	channel := make(chan entity.ChatMessage)
	go func() {
		sub, err := c.wsRouter.Subscribe(
			[]string{
				chatMessagesSubscriptionKeyAll,
			},
		)
		if err != nil {
			panic(err)
		}
		defer func() {
			sub.Unsubscribe()
			close(channel)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case data := <-sub.GetChannel():
				var msg entity.ChatMessage
				if err := json.Unmarshal(data, &msg); err != nil {
					panic(err)
				}

				channel <- msg
			}
		}
	}()

	return channel
}
