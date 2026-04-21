package kick

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/generic"
	kickbus "github.com/twirapp/twir/libs/bus-core/kick"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	"go.uber.org/fx"
)

const (
	idempotencyProcessingTTL = 30 * time.Second
	idempotencyTTL           = 10 * time.Minute
	idempotencyKeyPrefix     = "kick:event:"
	idempotencyStatusProcessing = "processing"
	idempotencyStatusProcessed  = "processed"
)

type kickUser struct {
	UserID       int            `json:"user_id"`
	Username     string         `json:"username"`
	ChannelSlug  string         `json:"channel_slug"`
	IsVerified   bool           `json:"is_verified"`
	ProfilePicture string       `json:"profile_picture"`
	Identity     *kickIdentity  `json:"identity,omitempty"`
}

type kickIdentity struct {
	UsernameColor string     `json:"username_color"`
	Badges        []kickBadge `json:"badges,omitempty"`
}

type kickBadge struct {
	Text  string `json:"text"`
	Type  string `json:"type"`
	Count int    `json:"count,omitempty"`
}

type kickChatMessagePayload struct {
	MessageID string    `json:"message_id"`
	Broadcaster kickUser `json:"broadcaster"`
	Sender    kickUser   `json:"sender"`
	Content   string     `json:"content"`
	Emotes    []kickEmote `json:"emotes,omitempty"`
	CreatedAt string      `json:"created_at"`
}

type kickEmote struct {
	EmoteID   string `json:"emote_id"`
	Positions []struct {
		S int `json:"s"`
		E int `json:"e"`
	} `json:"positions"`
}

type kickFollowPayload struct {
	Broadcaster kickUser `json:"broadcaster"`
	Follower    kickUser `json:"follower"`
}

type kickLivestreamStatusPayload struct {
	Broadcaster kickUser `json:"broadcaster"`
	IsLive      bool     `json:"is_live"`
	Title       string   `json:"title"`
	StartedAt   string   `json:"started_at"`
	EndedAt     *string  `json:"ended_at,omitempty"`
}

type Handlers struct {
	logger                *slog.Logger
	redis                 *redis.Client
	chatMessagesGeneric   bus_core.Queue[generic.ChatMessage, struct{}]
	processGenericMessage bus_core.Queue[generic.ChatMessage, struct{}]
	eventsFollow          bus_core.Queue[events.FollowMessage, struct{}]
	streamOnline          bus_core.Queue[kickbus.KickStreamOnline, struct{}]
	streamOffline         bus_core.Queue[kickbus.KickStreamOffline, struct{}]
	channelsRepo          channelsrepository.Repository
	usersRepo             usersrepository.Repository
}

type HandlersOpts struct {
	fx.In

	Logger       *slog.Logger
	Redis        *redis.Client
	Bus          *bus_core.Bus
	ChannelsRepo channelsrepository.Repository
	UsersRepo    usersrepository.Repository
}

func NewHandlers(opts HandlersOpts) *Handlers {
	return &Handlers{
		logger:                opts.Logger,
		redis:                 opts.Redis,
		chatMessagesGeneric:   opts.Bus.ChatMessagesGeneric,
		processGenericMessage: opts.Bus.Parser.ProcessGenericMessage,
		eventsFollow:          opts.Bus.Events.Follow,
		streamOnline:          opts.Bus.KickStreamOnline,
		streamOffline:         opts.Bus.KickStreamOffline,
		channelsRepo:          opts.ChannelsRepo,
		usersRepo:             opts.UsersRepo,
	}
}

func (h *Handlers) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	messageID := KickMessageIDFromContext(ctx)
	eventType := KickEventTypeFromContext(ctx)

	idempotencyKey := idempotencyKeyPrefix + messageID
	ok, err := h.redis.SetNX(ctx, idempotencyKey, idempotencyStatusProcessing, idempotencyProcessingTTL).Result()
	if err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to claim idempotency key",
			slog.String("message_id", messageID),
			slog.String("event_type", eventType),
			logger.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !ok {
		status, err := h.redis.Get(ctx, idempotencyKey).Result()
		if err != nil && !errors.Is(err, redis.Nil) {
			h.logger.ErrorContext(ctx, "kick: failed to read idempotency status",
				slog.String("message_id", messageID),
				slog.String("event_type", eventType),
				logger.Error(err),
			)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		switch status {
		case idempotencyStatusProcessed:
			h.logger.InfoContext(ctx, "kick: duplicate processed event, skipping",
				slog.String("message_id", messageID),
				slog.String("event_type", eventType),
			)
			w.WriteHeader(http.StatusOK)
		case idempotencyStatusProcessing:
			h.logger.InfoContext(ctx, "kick: event already processing, deferring",
				slog.String("message_id", messageID),
				slog.String("event_type", eventType),
			)
			w.WriteHeader(http.StatusAccepted)
		default:
			h.logger.InfoContext(ctx, "kick: idempotency key already exists, skipping",
				slog.String("message_id", messageID),
				slog.String("event_type", eventType),
				slog.String("status", status),
			)
			w.WriteHeader(http.StatusOK)
		}

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to read request body", logger.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	h.logger.InfoContext(ctx, "kick: received webhook event",
		slog.String("event_type", eventType),
		slog.String("message_id", messageID),
		slog.Int("body_size", len(body)),
	)

	switch eventType {
	case "chat.message.sent":
		err = h.handleChatMessage(r, body)
	case "channel.followed":
		err = h.handleChannelFollow(r, body)
	case "livestream.status.updated":
		err = h.handleLivestreamStatus(r, body)
	default:
		h.logger.InfoContext(ctx, "kick: unknown event type, ignoring",
			slog.String("event_type", eventType),
		)
	}

	if err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to process webhook event",
			slog.String("message_id", messageID),
			slog.String("event_type", eventType),
			logger.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	err = h.redis.Set(ctx, idempotencyKey, idempotencyStatusProcessed, idempotencyTTL).Err()
	if err != nil {
		h.logger.ErrorContext(ctx, "kick: failed to mark event as processed",
			slog.String("message_id", messageID),
			slog.String("event_type", eventType),
			logger.Error(err),
		)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) handleChatMessage(r *http.Request, body []byte) error {
	ctx := r.Context()

	h.logger.InfoContext(ctx, "kick: handling chat message",
		slog.Int("body_size", len(body)),
	)

	var payload kickChatMessagePayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("unmarshal chat message payload: %w", err)
	}

	h.logger.InfoContext(ctx, "kick: parsed chat message payload",
		slog.Int("broadcaster_user_id", payload.Broadcaster.UserID),
		slog.String("broadcaster_username", payload.Broadcaster.Username),
		slog.Int("sender_user_id", payload.Sender.UserID),
		slog.String("sender_username", payload.Sender.Username),
	)

	broadcasterUserID := strconv.Itoa(payload.Broadcaster.UserID)
	channelID, userID, err := h.resolveIDs(r, broadcasterUserID)
	if err != nil {
		return fmt.Errorf("resolve ids for chat message broadcaster_user_id=%s: %w", broadcasterUserID, err)
	}

	h.logger.InfoContext(ctx, "kick: resolved IDs for chat message",
		slog.String("channel_id", channelID),
		slog.String("user_id", userID),
	)

	var color string
	var badges []generic.ChatMessageBadge
	if payload.Sender.Identity != nil {
		color = payload.Sender.Identity.UsernameColor
		for _, b := range payload.Sender.Identity.Badges {
			badges = append(badges, generic.ChatMessageBadge{
				SetID: b.Type,
				Text:  b.Text,
			})
		}
	}

	genericMsg := generic.ChatMessage{
		Platform:          string(platform.PlatformKick),
		ChannelID:         channelID,
		UserID:            userID,
		PlatformChannelID: broadcasterUserID,
		SenderID:          strconv.Itoa(payload.Sender.UserID),
		SenderLogin:       payload.Sender.Username,
		SenderDisplayName: payload.Sender.Username,
		MessageID:         payload.MessageID,
		Text:              payload.Content,
		Badges:            badges,
		Color:             color,
	}

	if err := h.chatMessagesGeneric.Publish(ctx, genericMsg); err != nil {
		return fmt.Errorf("publish chat message to ChatMessagesGeneric: %w", err)
	}

	if err := h.processGenericMessage.Publish(ctx, genericMsg); err != nil {
		return fmt.Errorf("publish chat message to Parser.ProcessGenericMessage: %w", err)
	}

	return nil
}

func (h *Handlers) handleChannelFollow(r *http.Request, body []byte) error {
	ctx := r.Context()

	var payload kickFollowPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("unmarshal follow payload: %w", err)
	}

	broadcasterUserID := strconv.Itoa(payload.Broadcaster.UserID)
	channelID, _, err := h.resolveIDs(r, broadcasterUserID)
	if err != nil {
		return fmt.Errorf("resolve ids for follow broadcaster_user_id=%s: %w", broadcasterUserID, err)
	}

	if err := h.eventsFollow.Publish(
		ctx,
		events.FollowMessage{
			BaseInfo: events.BaseInfo{
				ChannelID: channelID,
			},
			UserID:   strconv.Itoa(payload.Follower.UserID),
			UserName: payload.Follower.Username,
		},
	); err != nil {
		return fmt.Errorf("publish follow event: %w", err)
	}

	return nil
}

func (h *Handlers) handleLivestreamStatus(r *http.Request, body []byte) error {
	ctx := r.Context()

	var payload kickLivestreamStatusPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return fmt.Errorf("unmarshal livestream.status.updated payload: %w", err)
	}

	broadcasterUserID := strconv.Itoa(payload.Broadcaster.UserID)
	channelID, _, err := h.resolveIDs(r, broadcasterUserID)
	if err != nil {
		return fmt.Errorf("resolve ids for livestream.status.updated broadcaster_user_id=%s: %w", broadcasterUserID, err)
	}

	if payload.IsLive {
		h.logger.InfoContext(ctx, "kick: stream online",
			slog.String("channel_id", channelID),
			slog.String("broadcaster_user_id", broadcasterUserID),
			slog.String("title", payload.Title),
		)

		if err := h.streamOnline.Publish(ctx, kickbus.KickStreamOnline{
			BroadcasterUserID:    broadcasterUserID,
			BroadcasterUserLogin: payload.Broadcaster.Username,
		}); err != nil {
			return fmt.Errorf("publish stream online event: %w", err)
		}
	} else {
		h.logger.InfoContext(ctx, "kick: stream offline",
			slog.String("channel_id", channelID),
			slog.String("broadcaster_user_id", broadcasterUserID),
		)

		if err := h.streamOffline.Publish(ctx, kickbus.KickStreamOffline{
			BroadcasterUserID:    broadcasterUserID,
			BroadcasterUserLogin: payload.Broadcaster.Username,
		}); err != nil {
			return fmt.Errorf("publish stream offline event: %w", err)
		}
	}

	return nil
}

func (h *Handlers) resolveIDs(r *http.Request, broadcasterUserID string) (string, string, error) {
	ctx := r.Context()

	user, err := h.usersRepo.GetByPlatformID(ctx, platform.PlatformKick, broadcasterUserID)
	if err != nil {
		if errors.Is(err, usersmodel.ErrNotFound) {
			return "", "", fmt.Errorf("no kick user for broadcaster_user_id=%s", broadcasterUserID)
		}
		return "", "", fmt.Errorf("get user by platform id: %w", err)
	}

	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		return "", "", fmt.Errorf("parse user id as uuid: %w", err)
	}

	channel, err := h.channelsRepo.GetByKickUserID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, channelsrepository.ErrNotFound) {
			return "", "", fmt.Errorf("channel not found for user_id=%s platform=kick", user.ID)
		}
		return "", "", fmt.Errorf("get channel by kick user id: %w", err)
	}

	return channel.ID.String(), user.ID, nil
}
