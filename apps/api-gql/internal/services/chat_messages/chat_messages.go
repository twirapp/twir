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
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/generic"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/chat_messages"
	"github.com/twirapp/twir/libs/repositories/chat_messages/model"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"github.com/twirapp/twir/libs/wsrouter"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	ChatMessagesRepository chat_messages.Repository
	ChannelService         *channelservice.ChannelService
	TwirBus                *buscore.Bus
	WsRouter               wsrouter.WsRouter
	Logger                 *slog.Logger
}

const (
	chatMessagesSubscriptionKey          = "api.chatMessages"
	chatMessagesSubscriptionKeyAll       = chatMessagesSubscriptionKey + ".All"
	chatOverlayModerationSubscriptionKey = "api.chatOverlayModeration"
)

func chatMessagesSubscriptionKeyCreate(platform string, channelId string) string {
	return chatMessagesSubscriptionKey + "." + platform + "." + channelId
}

func chatOverlayModerationSubscriptionKeyCreate(platform string, channelId string) string {
	return chatOverlayModerationSubscriptionKey + "." + platform + "." + channelId
}

func New(opts Opts) *Service {
	s := &Service{
		chatMessagesRepository: opts.ChatMessagesRepository,
		channelService:         opts.ChannelService,
		wsRouter:               opts.WsRouter,
		logger:                 opts.Logger,
		chanSubs:               make(map[string]struct{}),
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := opts.TwirBus.ChatMessages.Subscribe(s.handleBusEvent); err != nil {
					return err
				}
				if err := opts.TwirBus.Events.ChannelBan.Subscribe(s.handleChannelBanEvent); err != nil {
					return err
				}
				if err := opts.TwirBus.Events.ChannelMessageDelete.Subscribe(
					s.handleChannelMessageDeleteEvent,
				); err != nil {
					return err
				}
				return opts.TwirBus.Events.ChatClear.Subscribe(s.handleChatClearEvent)
			},
			OnStop: func(ctx context.Context) error {
				opts.TwirBus.ChatMessages.Unsubscribe()
				opts.TwirBus.Events.ChannelBan.Unsubscribe()
				opts.TwirBus.Events.ChannelMessageDelete.Unsubscribe()
				opts.TwirBus.Events.ChatClear.Unsubscribe()
				return nil
			},
		},
	)

	return s
}

type Service struct {
	chatMessagesRepository chat_messages.Repository
	channelService         chatMessagesChannelLookup

	wsRouter   wsrouter.WsRouter
	logger     *slog.Logger
	chanSubs   map[string]struct{}
	chanSubsMu sync.RWMutex
}

type chatMessagesChannelLookup interface {
	GetChannelByID(context.Context, uuid.UUID) (channelentity.Channel, error)
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
		MessageID:       data.MessageID,
		MessageType:     data.MessageType,
		SenderID:        data.SenderID,
		AnnounceColor:   data.AnnounceColor,
		Badges:          mapChatMessageBadges(data.Badges),
		Fragments:       mapChatMessageFragments(data.Message),
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
		channelSubKeys = append(
			channelSubKeys,
			chatMessagesSubscriptionKeyCreate(pair.Platform, pair.PlatformChannelID),
		)
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

func mapChatMessageBadges(badges []generic.ChatMessageBadge) []entity.ChatMessageBadge {
	if len(badges) == 0 {
		return nil
	}

	result := make([]entity.ChatMessageBadge, 0, len(badges))
	for _, badge := range badges {
		result = append(
			result, entity.ChatMessageBadge{
				SetID:     badge.SetID,
				VersionID: badge.ID,
				Text:      badge.Text,
			},
		)
	}

	return result
}

func mapChatMessageFragments(message *generic.ChatMessageMessage) []entity.ChatMessageFragment {
	if message == nil {
		return nil
	}

	fragments := make([]entity.ChatMessageFragment, 0, len(message.Fragments))
	for _, fragment := range message.Fragments {
		mapped := entity.ChatMessageFragment{
			Type: "text",
			Text: fragment.Text,
		}

		if fragment.Type == generic.FragmentType_EMOTE && fragment.Emote != nil {
			mapped.Type = "emote"
			mapped.EmoteID = fragment.Emote.ID
			mapped.EmoteURL = fragment.Emote.URL
		}

		fragments = append(fragments, mapped)
	}

	return fragments
}

func (c *Service) publishModerationEvent(
	platform string,
	platformChannelID string,
	event entity.ChatOverlayModerationEvent,
) {
	if platformChannelID == "" {
		return
	}

	key := chatOverlayModerationSubscriptionKeyCreate(platform, platformChannelID)

	c.chanSubsMu.RLock()
	_, ok := c.chanSubs[key]
	c.chanSubsMu.RUnlock()
	if !ok {
		return
	}

	if err := c.wsRouter.Publish(key, event); err != nil {
		c.logger.Error(
			"cannot publish chat overlay moderation event",
			logger.Error(err),
		)
	}
}

func (c *Service) handleChannelBanEvent(
	ctx context.Context,
	msg events.ChannelBanMessage,
) (struct{}, error) {
	platform := msg.BaseInfo.Platform
	if platform == "" {
		platform = platformentity.PlatformTwitch
	}

	platformChannelID := msg.BaseInfo.ChannelPlatformID
	if platform == platformentity.PlatformKick {
		channelUUID, err := uuid.Parse(msg.BaseInfo.ChannelPlatformID)
		if err != nil {
			c.logger.Error("cannot parse kick ban channel id", logger.Error(err))
			return struct{}{}, nil
		}

		channel, err := c.channelService.GetChannelByID(ctx, channelUUID)
		if err != nil {
			c.logger.Error("cannot resolve kick channel for ban event", logger.Error(err))
			return struct{}{}, nil
		}
		if channel.IsNil() {
			return struct{}{}, nil
		}

		binding, hasBinding := channel.Binding(platform)
		if !hasBinding {
			return struct{}{}, nil
		}

		platformChannelID = binding.PlatformChannelID
	}

	c.publishModerationEvent(
		platform.String(),
		platformChannelID,
		entity.ChatOverlayModerationEvent{
			Type:      entity.ChatOverlayModerationEventUserBanned,
			Platform:  platform.String(),
			UserLogin: msg.UserLogin,
		},
	)

	return struct{}{}, nil
}

func (c *Service) handleChannelMessageDeleteEvent(
	_ context.Context,
	msg events.ChannelMessageDeleteMessage,
) (struct{}, error) {
	platform := msg.BaseInfo.Platform
	if platform == "" {
		platform = platformentity.PlatformTwitch
	}

	c.publishModerationEvent(
		platform.String(),
		msg.BaseInfo.ChannelPlatformID,
		entity.ChatOverlayModerationEvent{
			Type:      entity.ChatOverlayModerationEventMessageDeleted,
			Platform:  platform.String(),
			MessageID: msg.MessageId,
		},
	)

	return struct{}{}, nil
}

func (c *Service) handleChatClearEvent(
	_ context.Context,
	msg events.ChatClearMessage,
) (struct{}, error) {
	platform := msg.BaseInfo.Platform
	if platform == "" {
		platform = platformentity.PlatformTwitch
	}

	c.publishModerationEvent(
		platform.String(),
		msg.BaseInfo.ChannelPlatformID,
		entity.ChatOverlayModerationEvent{
			Type:     entity.ChatOverlayModerationEventChatCleared,
			Platform: platform.String(),
		},
	)

	return struct{}{}, nil
}

func (c *Service) SubscribeToOverlayModerationEvents(
	ctx context.Context,
	channelPairs []chat_messages.PlatformChannelIdentity,
) <-chan entity.ChatOverlayModerationEvent {
	channelSubKeys := make([]string, 0, len(channelPairs))
	for _, pair := range channelPairs {
		channelSubKeys = append(
			channelSubKeys,
			chatOverlayModerationSubscriptionKeyCreate(pair.Platform, pair.PlatformChannelID),
		)
	}

	c.chanSubsMu.Lock()
	for _, channelSubKey := range channelSubKeys {
		c.chanSubs[channelSubKey] = struct{}{}
	}
	c.chanSubsMu.Unlock()

	channel := make(chan entity.ChatOverlayModerationEvent)
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

				var event entity.ChatOverlayModerationEvent
				if err := json.Unmarshal(data, &event); err != nil {
					panic(err)
				}

				channel <- event
			}
		}
	}()

	return channel
}
