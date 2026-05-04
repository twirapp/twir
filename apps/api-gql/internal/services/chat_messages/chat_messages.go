package chat_messages

import (
	"context"
	"encoding/json"
	"log/slog"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/generic"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/chat_messages"
	"github.com/twirapp/twir/libs/repositories/chat_messages/model"
	"github.com/twirapp/twir/libs/wsrouter"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	ChatMessagesRepository chat_messages.Repository
	TwirBus                *buscore.Bus
	WsRouter               wsrouter.WsRouter
	Logger                 *slog.Logger
}

const (
	chatMessagesSubscriptionKey    = "api.chatMessages"
	chatMessagesSubscriptionKeyAll = chatMessagesSubscriptionKey + ".All"
)

func chatMessagesSubscriptionKeyCreate(platform string, channelId string) string {
	return chatMessagesSubscriptionKey + "." + platform + "." + channelId
}

func New(opts Opts) *Service {
	s := &Service{
		chatMessagesRepository: opts.ChatMessagesRepository,
		wsRouter:               opts.WsRouter,
		logger:                 opts.Logger,
		chanSubs:               make(map[string]struct{}),
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

	wsRouter   wsrouter.WsRouter
	logger     *slog.Logger
	chanSubs   map[string]struct{}
	chanSubsMu sync.RWMutex
}

func (c *Service) modelToGql(m model.ChatMessage) entity.ChatMessage {
	return entity.ChatMessage{
		ID:              m.ID,
		Platform:        m.Platform,
		ChannelID:       m.ChannelID,
		UserID:          m.UserID,
		UserName:        m.UserName,
		UserDisplayName: m.UserDisplayName,
		UserColor:       m.UserColor,
		Text:            m.Text,
		CreatedAt:       m.CreatedAt,
	}
}

func (c *Service) handleBusEvent(_ context.Context, data generic.ChatMessage) (
	struct{},
	error,
) {
	textBuilder := strings.Builder{}
	if data.Message != nil {
		for _, fragment := range data.Message.Fragments {
			textBuilder.WriteString(fragment.Text)
		}
	}
	if textBuilder.Len() == 0 {
		if data.Message != nil {
			textBuilder.WriteString(data.Message.Text)
		} else {
			textBuilder.WriteString(data.Text)
		}
	}
	msg := entity.ChatMessage{
		ID:              uuid.New(),
		Platform:        data.Platform,
		ChannelID:       data.PlatformChannelID,
		ChannelLogin:    data.BroadcasterUserLogin,
		ChannelName:     data.BroadcasterUserName,
		UserID:          data.ChatterUserId,
		UserName:        data.ChatterUserLogin,
		UserDisplayName: data.ChatterUserName,
		UserColor:       data.Color,
		Text:            textBuilder.String(),
		CreatedAt:       time.Now(),
	}

	channelSubKey := chatMessagesSubscriptionKeyCreate(data.Platform, data.PlatformChannelID)
	c.chanSubsMu.RLock()
	if _, ok := c.chanSubs[channelSubKey]; ok {
		err := c.wsRouter.Publish(channelSubKey, msg)
		if err != nil {
			c.logger.Error(
				"cannot publish some message to separate broadcaster messages",
				logger.Error(err),
			)
		}
	}
	c.chanSubsMu.RUnlock()

	err := c.wsRouter.Publish(chatMessagesSubscriptionKeyAll, msg)
	if err != nil {
		c.logger.Error("cannot publish some message to all messages", logger.Error(err))
		return struct{}{}, err
	}

	return struct{}{}, nil
}

func (c *Service) SubscribeToNewMessagesByChannelIDs(
	ctx context.Context,
	channelPairs []chat_messages.PlatformChannelIdentity,
) <-chan entity.ChatMessage {
	channelSubKeys := make([]string, 0, len(channelPairs))
	for _, pair := range channelPairs {
		channelSubKeys = append(channelSubKeys, chatMessagesSubscriptionKeyCreate(pair.Platform, pair.PlatformChannelID))
	}

	c.chanSubsMu.Lock()
	for _, channelSubKey := range channelSubKeys {
		c.chanSubs[channelSubKey] = struct{}{}
	}
	c.chanSubsMu.Unlock()

	channel := make(chan entity.ChatMessage)
	go func() {
		sub, err := c.wsRouter.Subscribe(channelSubKeys)
		if err != nil {
			panic(err)
		}
		defer func() {
			c.chanSubsMu.Lock()
			for _, channelSubKey := range channelSubKeys {
				delete(c.chanSubs, channelSubKey)
			}
			c.chanSubsMu.Unlock()
			sub.Unsubscribe()
			close(channel)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-sub.GetChannel():
				if !ok {
					return
				}

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
			case data, ok := <-sub.GetChannel():
				if !ok {
					return
				}
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
