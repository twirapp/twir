package kick

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/redis/go-redis/v9"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/generic"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	user_platform_accounts "github.com/twirapp/twir/libs/repositories/user_platform_accounts"
	"go.uber.org/fx"
)

const (
	idempotencyTTL       = 10 * time.Minute
	idempotencyKeyPrefix = "kick:event:"
)

type kickChatMessagePayload struct {
	MessageID            string      `json:"message_id"`
	BroadcasterUserID    string      `json:"broadcaster_user_id"`
	BroadcasterUserLogin string      `json:"broadcaster_user_login"`
	SenderUserID         string      `json:"sender_user_id"`
	SenderUserLogin      string      `json:"sender_user_login"`
	SenderDisplayName    string      `json:"sender_display_name"`
	Content              string      `json:"content"`
	Color                string      `json:"color"`
	Badges               []kickBadge `json:"badges,omitempty"`
}

type kickBadge struct {
	SetID string `json:"set_id"`
	Text  string `json:"text"`
}

type kickFollowPayload struct {
	BroadcasterUserID string `json:"broadcaster_user_id"`
	FollowerUserID    string `json:"follower_user_id"`
	FollowerUserLogin string `json:"follower_user_login"`
}

type kickStreamOnlinePayload struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
}

type kickStreamOfflinePayload struct {
	BroadcasterUserID    string `json:"broadcaster_user_id"`
	BroadcasterUserLogin string `json:"broadcaster_user_login"`
}

type Handlers struct {
	logger                   *slog.Logger
	redis                    *redis.Client
	chatMessagesGeneric      bus_core.Queue[generic.ChatMessage, struct{}]
	processGenericMessage    bus_core.Queue[generic.ChatMessage, struct{}]
	eventsFollow             bus_core.Queue[events.FollowMessage, struct{}]
	channelsRepo             channelsrepository.Repository
	userPlatformAccountsRepo user_platform_accounts.Repository
}

type HandlersOpts struct {
	fx.In

	Logger                   *slog.Logger
	Redis                    *redis.Client
	Bus                      *bus_core.Bus
	ChannelsRepo             channelsrepository.Repository
	UserPlatformAccountsRepo user_platform_accounts.Repository
}

func NewHandlers(opts HandlersOpts) *Handlers {
	return &Handlers{
		logger:                   opts.Logger,
		redis:                    opts.Redis,
		chatMessagesGeneric:      opts.Bus.ChatMessagesGeneric,
		processGenericMessage:    opts.Bus.Parser.ProcessGenericMessage,
		eventsFollow:             opts.Bus.Events.Follow,
		channelsRepo:             opts.ChannelsRepo,
		userPlatformAccountsRepo: opts.UserPlatformAccountsRepo,
	}
}

func (h *Handlers) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	messageID := KickMessageIDFromContext(ctx)
	eventType := KickEventTypeFromContext(ctx)

	idempotencyKey := idempotencyKeyPrefix + messageID
	ok, err := h.redis.SetNX(ctx, idempotencyKey, "1", idempotencyTTL).Result()
	if err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to check idempotency key",
			slog.String("message_id", messageID),
			slog.String("event_type", eventType),
			logger.Error(err),
		)
	} else if !ok {
		h.logger.InfoContext(ctx, "kick: duplicate event, skipping",
			slog.String("message_id", messageID),
			slog.String("event_type", eventType),
		)
		w.WriteHeader(http.StatusOK)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to read request body", logger.Error(err))
		w.WriteHeader(http.StatusOK)
		return
	}

	switch eventType {
	case "chat.message.sent":
		h.handleChatMessage(w, r, body)
	case "channel.follow":
		h.handleChannelFollow(w, r, body)
	case "stream.online":
		h.handleStreamOnline(w, r, body)
	case "stream.offline":
		h.handleStreamOffline(w, r, body)
	default:
		h.logger.InfoContext(ctx, "kick: unknown event type, ignoring",
			slog.String("event_type", eventType),
		)
		w.WriteHeader(http.StatusOK)
	}
}

func (h *Handlers) handleChatMessage(w http.ResponseWriter, r *http.Request, body []byte) {
	ctx := r.Context()
	defer w.WriteHeader(http.StatusOK)

	var payload kickChatMessagePayload
	if err := json.Unmarshal(body, &payload); err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to unmarshal chat message payload", logger.Error(err))
		return
	}

	channelID, userID, err := h.resolveIDs(r, payload.BroadcasterUserID)
	if err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to resolve IDs for chat message",
			slog.String("broadcaster_user_id", payload.BroadcasterUserID),
			logger.Error(err),
		)
		return
	}

	badges := make([]generic.ChatMessageBadge, 0, len(payload.Badges))
	for _, b := range payload.Badges {
		badges = append(badges, generic.ChatMessageBadge{
			SetID: b.SetID,
			Text:  b.Text,
		})
	}

	genericMsg := generic.ChatMessage{
		Platform:          string(platform.PlatformKick),
		ChannelID:         channelID,
		UserID:            userID,
		PlatformChannelID: payload.BroadcasterUserID,
		SenderID:          payload.SenderUserID,
		SenderLogin:       payload.SenderUserLogin,
		SenderDisplayName: payload.SenderDisplayName,
		MessageID:         payload.MessageID,
		Text:              payload.Content,
		Badges:            badges,
		Color:             payload.Color,
	}

	if err := h.chatMessagesGeneric.Publish(ctx, genericMsg); err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to publish to ChatMessagesGeneric", logger.Error(err))
	}

	if err := h.processGenericMessage.Publish(ctx, genericMsg); err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to publish to Parser.ProcessGenericMessage", logger.Error(err))
	}
}

func (h *Handlers) handleChannelFollow(w http.ResponseWriter, r *http.Request, body []byte) {
	ctx := r.Context()
	defer w.WriteHeader(http.StatusOK)

	var payload kickFollowPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to unmarshal follow payload", logger.Error(err))
		return
	}

	channelID, _, err := h.resolveIDs(r, payload.BroadcasterUserID)
	if err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to resolve IDs for follow",
			slog.String("broadcaster_user_id", payload.BroadcasterUserID),
			logger.Error(err),
		)
		return
	}

	if err := h.eventsFollow.Publish(
		ctx,
		events.FollowMessage{
			BaseInfo: events.BaseInfo{
				ChannelID: channelID,
			},
			UserID:   payload.FollowerUserID,
			UserName: payload.FollowerUserLogin,
		},
	); err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to publish follow event", logger.Error(err))
	}
}

func (h *Handlers) handleStreamOnline(w http.ResponseWriter, r *http.Request, body []byte) {
	ctx := r.Context()
	defer w.WriteHeader(http.StatusOK)

	var payload kickStreamOnlinePayload
	if err := json.Unmarshal(body, &payload); err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to unmarshal stream.online payload", logger.Error(err))
		return
	}

	channelID, _, err := h.resolveIDs(r, payload.BroadcasterUserID)
	if err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to resolve IDs for stream.online",
			slog.String("broadcaster_user_id", payload.BroadcasterUserID),
			logger.Error(err),
		)
		return
	}

	// TODO: publish to generic stream online/offline bus topic when available
	h.logger.InfoContext(ctx, "kick: stream online",
		slog.String("channel_id", channelID),
		slog.String("broadcaster_user_id", payload.BroadcasterUserID),
	)
}

func (h *Handlers) handleStreamOffline(w http.ResponseWriter, r *http.Request, body []byte) {
	ctx := r.Context()
	defer w.WriteHeader(http.StatusOK)

	var payload kickStreamOfflinePayload
	if err := json.Unmarshal(body, &payload); err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to unmarshal stream.offline payload", logger.Error(err))
		return
	}

	channelID, _, err := h.resolveIDs(r, payload.BroadcasterUserID)
	if err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to resolve IDs for stream.offline",
			slog.String("broadcaster_user_id", payload.BroadcasterUserID),
			logger.Error(err),
		)
		return
	}

	// TODO: publish to generic stream online/offline bus topic when available
	h.logger.InfoContext(ctx, "kick: stream offline",
		slog.String("channel_id", channelID),
		slog.String("broadcaster_user_id", payload.BroadcasterUserID),
	)
}

func (h *Handlers) resolveIDs(r *http.Request, broadcasterUserID string) (string, string, error) {
	ctx := r.Context()

	account, err := h.userPlatformAccountsRepo.GetByPlatformUserID(ctx, platform.PlatformKick, broadcasterUserID)
	if err != nil {
		if errors.Is(err, user_platform_accounts.ErrNotFound) {
			return "", "", fmt.Errorf("no kick platform account for broadcaster_user_id=%s", broadcasterUserID)
		}
		return "", "", fmt.Errorf("get platform account: %w", err)
	}

	channel, err := h.channelsRepo.GetByUserIDAndPlatform(ctx, account.UserID, platform.PlatformKick)
	if err != nil {
		return "", "", fmt.Errorf("get channel by user_id and platform: %w", err)
	}
	if channel.IsNil() {
		return "", "", fmt.Errorf("channel not found for user_id=%s platform=kick", account.UserID)
	}

	return channel.ID, account.UserID.String(), nil
}
