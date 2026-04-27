package kick

import (
	"context"
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
	user_creator "github.com/twirapp/twir/apps/eventsub/internal/services/user-creator"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/generic"
	kickbus "github.com/twirapp/twir/libs/bus-core/kick"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/redis_keys"
	channelsinfohistory "github.com/twirapp/twir/libs/repositories/channels_info_history"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	streams "github.com/twirapp/twir/libs/repositories/streams"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	usersstatsmodel "github.com/twirapp/twir/libs/repositories/users_stats/model"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
)

const (
	idempotencyProcessingTTL    = 30 * time.Second
	idempotencyTTL              = 10 * time.Minute
	idempotencyKeyPrefix        = "kick:event:"
	idempotencyStatusProcessing = "processing"
	idempotencyStatusProcessed  = "processed"
)

type kickUser struct {
	UserID         int           `json:"user_id"`
	Username       string        `json:"username"`
	ChannelSlug    string        `json:"channel_slug"`
	IsVerified     bool          `json:"is_verified"`
	ProfilePicture string        `json:"profile_picture"`
	Identity       *kickIdentity `json:"identity,omitempty"`
}

type kickIdentity struct {
	UsernameColor string      `json:"username_color"`
	Badges        []kickBadge `json:"badges,omitempty"`
}

type kickBadge struct {
	Text  string `json:"text"`
	Type  string `json:"type"`
	Count int    `json:"count,omitempty"`
}

type kickChatMessagePayload struct {
	MessageID   string      `json:"message_id"`
	Broadcaster kickUser    `json:"broadcaster"`
	Sender      kickUser    `json:"sender"`
	Content     string      `json:"content"`
	Emotes      []kickEmote `json:"emotes,omitempty"`
	CreatedAt   string      `json:"created_at"`
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

type kickLivestreamMetadataPayload struct {
	Broadcaster kickUser `json:"broadcaster"`
	Metadata    struct {
		Title            string `json:"title"`
		Language         string `json:"language"`
		HasMatureContent bool   `json:"has_mature_content"`
		Category         struct {
			ID        int    `json:"id"`
			Name      string `json:"name"`
			Thumbnail string `json:"thumbnail"`
		} `json:"category"`
	} `json:"metadata"`
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
	kickBotsRepo          kickbotsrepository.Repository
	channelsInfoHistoryRepo channelsinfohistory.Repository
	streamsRepo           streamsrepository.Repository
	userCreatorService    *user_creator.UserCreatorService
	prefixCache           *generic_cacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix]
}

type HandlersOpts struct {
	fx.In

	Logger             *slog.Logger
	Redis              *redis.Client
	Bus                *bus_core.Bus
	ChannelsRepo       channelsrepository.Repository
	UsersRepo          usersrepository.Repository
	KickBotsRepo       kickbotsrepository.Repository
	ChannelsInfoHistoryRepo channelsinfohistory.Repository
	StreamsRepo        streamsrepository.Repository
	UserCreatorService *user_creator.UserCreatorService
	PrefixCache        *generic_cacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix]
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
		kickBotsRepo:          opts.KickBotsRepo,
		channelsInfoHistoryRepo: opts.ChannelsInfoHistoryRepo,
		streamsRepo:           opts.StreamsRepo,
		userCreatorService:    opts.UserCreatorService,
		prefixCache:           opts.PrefixCache,
	}
}

func (h *Handlers) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	messageID := KickMessageIDFromContext(ctx)
	eventType := KickEventTypeFromContext(ctx)
	eventVersion := KickEventVersionFromContext(ctx)
	subscriptionID := KickSubscriptionIDFromContext(ctx)
	timestamp := KickMessageTimestampFromContext(ctx)

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

	var eventLogAttrs []slog.Attr

	switch eventType {
	case "chat.message.sent":
		eventLogAttrs, err = h.handleChatMessage(r, body)
	case "channel.followed":
		eventLogAttrs, err = h.handleChannelFollow(r, body)
	case "livestream.status.updated":
		eventLogAttrs, err = h.handleLivestreamStatus(r, body)
	case "livestream.metadata.updated":
		eventLogAttrs, err = h.handleLivestreamMetadata(r, body)
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
		if delErr := h.redis.Del(ctx, idempotencyKey).Err(); delErr != nil {
			h.logger.ErrorContext(ctx, "kick: failed to clean up processing key after error",
				slog.String("message_id", messageID),
				logger.Error(delErr),
			)
		}
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

	logAttrs := []slog.Attr{
		slog.String("message_id", messageID),
		slog.String("event_type", eventType),
		slog.String("method", r.Method),
		slog.String("path", r.URL.Path),
	}
	if eventVersion != "" {
		logAttrs = append(logAttrs, slog.String("event_version", eventVersion))
	}
	if subscriptionID != "" {
		logAttrs = append(logAttrs, slog.String("subscription_id", subscriptionID))
	}
	if timestamp != "" {
		logAttrs = append(logAttrs, slog.String("timestamp", timestamp))
	}
	logAttrs = append(logAttrs, eventLogAttrs...)

	h.logger.LogAttrs(ctx, slog.LevelInfo, "received kick webhook event", logAttrs...)

	w.WriteHeader(http.StatusOK)
}

func (h *Handlers) handleLivestreamMetadata(r *http.Request, body []byte) ([]slog.Attr, error) {
	ctx := r.Context()

	var payload kickLivestreamMetadataPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal livestream.metadata.updated payload: %w", err)
	}

	broadcasterUserID := strconv.Itoa(payload.Broadcaster.UserID)
	channelID, _, err := h.resolveIDs(r, broadcasterUserID)
	if err != nil {
		return nil, fmt.Errorf("resolve ids for livestream.metadata.updated broadcaster_user_id=%s: %w", broadcasterUserID, err)
	}

	if h.streamsRepo != nil {
		categoryID := strconv.Itoa(payload.Metadata.Category.ID)
		thumbnailURL := payload.Metadata.Category.Thumbnail
		if err := h.streamsRepo.Update(ctx, channelID, streams.UpdateInput{
			GameId:       &categoryID,
			GameName:     &payload.Metadata.Category.Name,
			Title:        &payload.Metadata.Title,
			Language:     &payload.Metadata.Language,
			ThumbnailUrl: &thumbnailURL,
			IsMature:     &payload.Metadata.HasMatureContent,
		}); err != nil {
			return nil, fmt.Errorf("update kick current stream metadata: %w", err)
		}
	}
	if h.channelsInfoHistoryRepo != nil && payload.Metadata.Category.Name != "" {
		if err := h.channelsInfoHistoryRepo.Create(ctx, channelsinfohistory.CreateInput{
			ChannelID: channelID,
			Title:     payload.Metadata.Title,
			Category:  payload.Metadata.Category.Name,
		}); err != nil {
			return nil, fmt.Errorf("create kick channel info history: %w", err)
		}
	}

	return []slog.Attr{
		slog.String("channel_id", channelID),
		slog.String("broadcaster_user_id", broadcasterUserID),
		slog.String("broadcaster_username", payload.Broadcaster.Username),
		slog.String("title", payload.Metadata.Title),
		slog.String("language", payload.Metadata.Language),
		slog.String("category", payload.Metadata.Category.Name),
		slog.Int("category_id", payload.Metadata.Category.ID),
	}, nil
}

func (h *Handlers) handleChatMessage(r *http.Request, body []byte) ([]slog.Attr, error) {
	ctx := r.Context()

	h.logger.DebugContext(ctx, "kick: handling chat message",
		slog.Int("body_size", len(body)),
	)

	var payload kickChatMessagePayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal chat message payload: %w", err)
	}

	h.logger.DebugContext(ctx, "kick: parsed chat message payload",
		slog.Int("broadcaster_user_id", payload.Broadcaster.UserID),
		slog.String("broadcaster_username", payload.Broadcaster.Username),
		slog.Int("sender_user_id", payload.Sender.UserID),
		slog.String("sender_username", payload.Sender.Username),
	)

	broadcasterUserID := strconv.Itoa(payload.Broadcaster.UserID)
	channelID, _, err := h.resolveIDs(r, broadcasterUserID)
	if err != nil {
		return nil, fmt.Errorf("resolve ids for chat message broadcaster_user_id=%s: %w", broadcasterUserID, err)
	}

	senderPlatformID := strconv.Itoa(payload.Sender.UserID)

	var color string
	var badges []generic.ChatMessageBadge
	var isBroadcaster, isModerator, isVip, isSubscriber bool
	if payload.Sender.Identity != nil {
		color = payload.Sender.Identity.UsernameColor
		for _, b := range payload.Sender.Identity.Badges {
			badges = append(badges, generic.ChatMessageBadge{
				SetID: b.Type,
				Text:  b.Text,
			})
			switch b.Type {
			case "broadcaster":
				isBroadcaster = true
			case "moderator":
				isModerator = true
			case "vip":
				isVip = true
			case "subscriber":
				isSubscriber = true
			}
		}
	}

	if senderPlatformID == broadcasterUserID {
		isBroadcaster = true
	}

	var (
		channel         channelsmodel.Channel
		stream          *streamsmodel.Stream
		commandsPrefix  string
	)

	var errwg errgroup.Group

		errwg.Go(func() error {
			var err error
			channel, err = h.channelsRepo.GetByID(ctx, uuid.MustParse(channelID))
			if err != nil {
			return fmt.Errorf("get channel by id: %w", err)
		}
		return nil
	})

	errwg.Go(func() error {
		var err error
		stream, err = h.getChannelStream(ctx, channelID)
		if err != nil {
			return fmt.Errorf("get channel stream: %w", err)
		}
		return nil
	})

	errwg.Go(func() error {
		var err error
		commandsPrefix, err = h.getChannelCommandPrefix(ctx, channelID)
		if err != nil {
			return fmt.Errorf("get channel command prefix: %w", err)
		}
		return nil
	})

	if err := errwg.Wait(); err != nil {
		return nil, err
	}

	senderUser, err := h.usersRepo.GetByPlatformID(ctx, platform.PlatformKick, senderPlatformID)
	if err != nil {
		if !errors.Is(err, usersmodel.ErrNotFound) {
			return nil, fmt.Errorf("get sender user by platform id: %w", err)
		}
		createdUser, createErr := h.usersRepo.Create(ctx, usersrepository.CreateInput{
			Platform:    platform.PlatformKick,
			PlatformID:  senderPlatformID,
			Login:       payload.Sender.Username,
			DisplayName: payload.Sender.Username,
		})
		if createErr != nil {
			return nil, fmt.Errorf("create sender user: %w", createErr)
		}
		senderUser = createdUser
	}

	eventAttrs := []slog.Attr{
		slog.String("channel_id", channelID),
		slog.String("broadcaster_user_id", broadcasterUserID),
		slog.String("broadcaster_username", payload.Broadcaster.Username),
		slog.String("sender_user_id", senderPlatformID),
		slog.String("sender_username", payload.Sender.Username),
		slog.String("sender_twir_user_id", senderUser.ID),
		slog.Bool("is_broadcaster", isBroadcaster),
		slog.Bool("is_moderator", isModerator),
		slog.Bool("is_vip", isVip),
		slog.Bool("is_subscriber", isSubscriber),
	}

	if h.shouldIgnoreBotSelfMessage(ctx, channel, senderUser, payload.Sender.Username) {
		return append(eventAttrs, slog.Bool("ignored_self_message", true)), nil
	}

	senderStats := &usersstatsmodel.UserStat{}
	if h.userCreatorService != nil {
		_, senderStats, err = h.userCreatorService.UnsureUser(
			ctx,
			user_creator.CreateUserInput{
				UserID:            senderUser.ID,
				PlatformID:        senderPlatformID,
				Platform:          platform.PlatformKick,
				ChannelID:         &channelID,
				IsBroadcaster:     isBroadcaster,
				IsModerator:       isModerator,
				IsVip:             isVip,
				IsSubscriber:      isSubscriber,
				ShouldUpdateStats: stream != nil && stream.ID != "",
			},
		)
		if err != nil {
			return nil, fmt.Errorf("ensure sender user stats: %w", err)
		}
		if senderStats == nil {
			senderStats = &usersstatsmodel.UserStat{}
		}
	}

	genericMsg := generic.ChatMessage{
		Platform:          string(platform.PlatformKick),
		ChannelID:         channelID,
		UserID:            senderUser.ID,
		PlatformChannelID: broadcasterUserID,
		SenderID:          senderPlatformID,
		SenderLogin:       payload.Sender.Username,
		SenderDisplayName: payload.Sender.Username,
		MessageID:         payload.MessageID,
		Text:              payload.Content,
		Badges:            badges,
		Color:             color,
		EnrichedData: generic.ChatMessageEnrichedData{
			ChannelCommandPrefix: commandsPrefix,
			DbChannel:            channel,
			ChannelStream:        stream,
			DbUser: &generic.DbUser{
				ID:                senderUser.ID,
				TokenID:           senderUser.TokenID.Ptr(),
				IsBotAdmin:        senderUser.IsBotAdmin,
				ApiKey:            senderUser.ApiKey,
				IsBanned:          senderUser.IsBanned,
				HideOnLandingPage: senderUser.HideOnLandingPage,
				CreatedAt:         senderUser.CreatedAt,
			},
			DbUserChannelStat: &generic.DbUserChannelStat{
				ID:                senderStats.ID,
				UserID:            senderStats.UserID,
				ChannelID:         senderStats.ChannelID,
				Messages:          senderStats.Messages,
				Watched:           senderStats.Watched,
				UsedChannelPoints: senderStats.UsedChannelPoints,
				IsMod:             senderStats.IsMod,
				IsVip:             senderStats.IsVip,
				IsSubscriber:      senderStats.IsSubscriber,
				Reputation:        senderStats.Reputation,
				Emotes:            senderStats.Emotes,
				CreatedAt:         senderStats.CreatedAt,
				UpdatedAt:         senderStats.UpdatedAt,
			},
			IsChatterBroadcaster: isBroadcaster,
			IsChatterModerator:   isModerator,
			IsChatterVip:         isVip,
			IsChatterSubscriber:  isSubscriber,
		},
	}

	if err := h.chatMessagesGeneric.Publish(ctx, genericMsg); err != nil {
		return nil, fmt.Errorf("publish chat message to ChatMessagesGeneric: %w", err)
	}

	if err := h.processGenericMessage.Publish(ctx, genericMsg); err != nil {
		return nil, fmt.Errorf("publish chat message to Parser.ProcessGenericMessage: %w", err)
	}

	return append(eventAttrs, slog.Bool("ignored_self_message", false)), nil
}

func (h *Handlers) shouldIgnoreBotSelfMessage(
	ctx context.Context,
	channel channelsmodel.Channel,
	senderUser usersmodel.User,
	senderUsername string,
) bool {
	if h.kickBotsRepo == nil || channel.KickBotID == nil {
		return false
	}

	bot, err := h.kickBotsRepo.GetByID(ctx, *channel.KickBotID)
	if err != nil {
		if errors.Is(err, kickbotsrepository.ErrNotFound) {
			return false
		}

		h.logger.DebugContext(ctx, "kick: failed to resolve assigned bot for self-message guard",
			slog.String("channel_id", channel.ID.String()),
			slog.String("kick_bot_id", channel.KickBotID.String()),
			slog.String("sender_username", senderUsername),
			logger.Error(err),
		)
		return false
	}

	return bot.KickUserID.String() == senderUser.ID
}

func (h *Handlers) handleChannelFollow(r *http.Request, body []byte) ([]slog.Attr, error) {
	ctx := r.Context()

	var payload kickFollowPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal follow payload: %w", err)
	}

	broadcasterUserID := strconv.Itoa(payload.Broadcaster.UserID)
	channelID, _, err := h.resolveIDs(r, broadcasterUserID)
	if err != nil {
		return nil, fmt.Errorf("resolve ids for follow broadcaster_user_id=%s: %w", broadcasterUserID, err)
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
		return nil, fmt.Errorf("publish follow event: %w", err)
	}

	return []slog.Attr{
		slog.String("channel_id", channelID),
		slog.String("broadcaster_user_id", broadcasterUserID),
		slog.String("broadcaster_username", payload.Broadcaster.Username),
		slog.String("follower_user_id", strconv.Itoa(payload.Follower.UserID)),
		slog.String("follower_username", payload.Follower.Username),
	}, nil
}

func (h *Handlers) handleLivestreamStatus(r *http.Request, body []byte) ([]slog.Attr, error) {
	ctx := r.Context()

	var payload kickLivestreamStatusPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal livestream.status.updated payload: %w", err)
	}

	broadcasterUserID := strconv.Itoa(payload.Broadcaster.UserID)
	channelID, _, err := h.resolveIDs(r, broadcasterUserID)
	if err != nil {
		return nil, fmt.Errorf("resolve ids for livestream.status.updated broadcaster_user_id=%s: %w", broadcasterUserID, err)
	}

	if payload.IsLive {
		startedAt := time.Now().UTC()
		if payload.StartedAt != "" {
			if parsedStartedAt, err := time.Parse(time.RFC3339, payload.StartedAt); err == nil {
				startedAt = parsedStartedAt
			}
		}

		if h.streamsRepo != nil {
			if err := h.streamsRepo.Save(ctx, streams.SaveInput{
				ID:        channelID,
				UserId:    channelID,
				UserLogin: payload.Broadcaster.ChannelSlug,
				UserName:  payload.Broadcaster.Username,
				Type:      "live",
				Title:     payload.Title,
				StartedAt: startedAt,
			}); err != nil {
				return nil, fmt.Errorf("save kick current stream: %w", err)
			}
		}

		if err := h.streamOnline.Publish(ctx, kickbus.KickStreamOnline{
			BroadcasterUserID:    broadcasterUserID,
			BroadcasterUserLogin: payload.Broadcaster.Username,
		}); err != nil {
			return nil, fmt.Errorf("publish stream online event: %w", err)
		}
	} else {
		if h.streamsRepo != nil {
			if err := h.streamsRepo.DeleteByChannelID(ctx, channelID); err != nil {
				return nil, fmt.Errorf("delete kick current stream: %w", err)
			}
		}

		if err := h.streamOffline.Publish(ctx, kickbus.KickStreamOffline{
			BroadcasterUserID:    broadcasterUserID,
			BroadcasterUserLogin: payload.Broadcaster.Username,
		}); err != nil {
			return nil, fmt.Errorf("publish stream offline event: %w", err)
		}
	}

	return []slog.Attr{
		slog.String("channel_id", channelID),
		slog.String("broadcaster_user_id", broadcasterUserID),
		slog.String("broadcaster_username", payload.Broadcaster.Username),
		slog.Bool("is_live", payload.IsLive),
		slog.String("title", payload.Title),
	}, nil
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

func (h *Handlers) getChannelCommandPrefix(ctx context.Context, channelId string) (
	string,
	error,
) {
	commandsPrefix := "!"
	if h.prefixCache == nil {
		return commandsPrefix, nil
	}

	fetchedCommandsPrefix, err := h.prefixCache.Get(ctx, channelId)
	if err != nil && !errors.Is(err, channelscommandsprefixrepository.ErrNotFound) {
		return "", err
	}

	if fetchedCommandsPrefix != channelscommandsprefixmodel.Nil {
		commandsPrefix = fetchedCommandsPrefix.Prefix
	} else {
		prefixCtx := context.WithoutCancel(ctx)

		go func() {
			if err := h.prefixCache.SetValue(
				prefixCtx,
				channelId,
				channelscommandsprefixmodel.ChannelsCommandsPrefix{
					ID:        uuid.New(),
					ChannelID: channelId,
					Prefix:    commandsPrefix,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				},
			); err != nil {
				h.logger.Error("cannot set default command prefix", logger.Error(err))
			}
		}()
	}

	return commandsPrefix, nil
}

func (h *Handlers) getChannelStream(
	ctx context.Context,
	channelId string,
) (*streamsmodel.Stream, error) {
	if h.streamsRepo == nil {
		return nil, nil
	}

	cacheKey := redis_keys.StreamByChannelID(channelId)
	cachedBytes, err := h.redis.Get(ctx, cacheKey).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, fmt.Errorf("failed to get stream cache: %w", err)
	}

	if len(cachedBytes) > 0 {
		var stream streamsmodel.Stream
		if err := json.Unmarshal(cachedBytes, &stream); err != nil {
			return nil, err
		}

		return &stream, nil
	}

	stream, err := h.streamsRepo.GetByChannelID(ctx, channelId)
	if err != nil {
		return nil, fmt.Errorf("failed to get stream by channel id: %w", err)
	}

	if stream.ID == "" {
		return nil, nil
	}

	streamBytes, err := json.Marshal(stream)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal stream: %w", err)
	}

	if err := h.redis.Set(ctx, cacheKey, streamBytes, 30*time.Second).Err(); err != nil {
		return nil, fmt.Errorf("failed to set stream cache: %w", err)
	}

	return &stream, nil
}
