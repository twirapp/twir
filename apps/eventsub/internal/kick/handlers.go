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
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	user_creator "github.com/twirapp/twir/apps/eventsub/internal/services/user-creator"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/generic"
	kickbus "github.com/twirapp/twir/libs/bus-core/kick"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/redis_keys"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	channelscommandsprefixrepository "github.com/twirapp/twir/libs/repositories/channels_commands_prefix"
	channelscommandsprefixmodel "github.com/twirapp/twir/libs/repositories/channels_commands_prefix/model"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	channelseventslistmodel "github.com/twirapp/twir/libs/repositories/channels_events_list/model"
	channelsinfohistory "github.com/twirapp/twir/libs/repositories/channels_info_history"
	channelsredemptionshistory "github.com/twirapp/twir/libs/repositories/channels_redemptions_history"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	streams "github.com/twirapp/twir/libs/repositories/streams"
	streamsrepository "github.com/twirapp/twir/libs/repositories/streams"
	streamsmodel "github.com/twirapp/twir/libs/repositories/streams/model"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	usersstatsmodel "github.com/twirapp/twir/libs/repositories/users_stats/model"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
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

type kickSubscriptionPayload struct {
	Broadcaster kickUser `json:"broadcaster"`
	Subscriber  kickUser `json:"subscriber"`
	Duration    int      `json:"duration"`
	CreatedAt   string   `json:"created_at"`
	ExpiresAt   string   `json:"expires_at"`
}

type kickSubscriptionGiftsPayload struct {
	Broadcaster kickUser   `json:"broadcaster"`
	Gifter      kickUser   `json:"gifter"`
	Giftees     []kickUser `json:"giftees"`
	CreatedAt   string     `json:"created_at"`
	ExpiresAt   string     `json:"expires_at"`
}

type kickRewardRedemptionPayload struct {
	ID         string `json:"id"`
	UserInput  string `json:"user_input"`
	Status     string `json:"status"`
	RedeemedAt string `json:"redeemed_at"`
	Reward     struct {
		ID          string `json:"id"`
		Title       string `json:"title"`
		Cost        int    `json:"cost"`
		Description string `json:"description"`
	} `json:"reward"`
	Redeemer struct {
		UserID      int    `json:"user_id"`
		Username    string `json:"username"`
		ChannelSlug string `json:"channel_slug"`
	} `json:"redeemer"`
	Broadcaster struct {
		UserID      int    `json:"user_id"`
		Username    string `json:"username"`
		ChannelSlug string `json:"channel_slug"`
	} `json:"broadcaster"`
}

type kickModerationBannedPayload struct {
	Broadcaster kickUser `json:"broadcaster"`
	Moderator   kickUser `json:"moderator"`
	BannedUser  kickUser `json:"banned_user"`
	Metadata    struct {
		Reason    string  `json:"reason"`
		CreatedAt string  `json:"created_at"`
		ExpiresAt *string `json:"expires_at"`
	} `json:"metadata"`
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
	logger                  *slog.Logger
	redis                   *redis.Client
	chatMessagesGeneric     bus_core.Queue[generic.ChatMessage, struct{}]
	processGenericMessage   bus_core.Queue[generic.ChatMessage, struct{}]
	eventsFollow            bus_core.Queue[events.FollowMessage, struct{}]
	eventsSubscribe         bus_core.Queue[events.SubscribeMessage, struct{}]
	eventsReSubscribe       bus_core.Queue[events.ReSubscribeMessage, struct{}]
	eventsSubGift           bus_core.Queue[events.SubGiftMessage, struct{}]
	eventsRedemptionCreated bus_core.Queue[events.RedemptionCreatedMessage, struct{}]
	eventsChannelBan        bus_core.Queue[events.ChannelBanMessage, struct{}]
	streamOnline            bus_core.Queue[kickbus.KickStreamOnline, struct{}]
	streamOffline           bus_core.Queue[kickbus.KickStreamOffline, struct{}]
	channelsRepo            channelsrepository.Repository
	usersRepo               usersrepository.Repository
	kickBotsRepo            kickbotsrepository.Repository
	eventsListRepo          channelseventslist.Repository
	channelsInfoHistoryRepo channelsinfohistory.Repository
	redemptionsHistoryRepo  channelsredemptionshistory.Repository
	streamsRepo             streamsrepository.Repository
	userCreatorService      *user_creator.UserCreatorService
	prefixCache             *generic_cacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix]
}

type HandlersOpts struct {
	fx.In

	Logger                  *slog.Logger
	Redis                   *redis.Client
	Bus                     *bus_core.Bus
	ChannelsRepo            channelsrepository.Repository
	UsersRepo               usersrepository.Repository
	KickBotsRepo            kickbotsrepository.Repository
	EventsListRepo          channelseventslist.Repository
	ChannelsInfoHistoryRepo channelsinfohistory.Repository
	RedemptionsHistoryRepo  channelsredemptionshistory.Repository
	StreamsRepo             streamsrepository.Repository
	UserCreatorService      *user_creator.UserCreatorService
	PrefixCache             *generic_cacher.GenericCacher[channelscommandsprefixmodel.ChannelsCommandsPrefix]
}

func NewHandlers(opts HandlersOpts) *Handlers {
	return &Handlers{
		logger:                  opts.Logger,
		redis:                   opts.Redis,
		chatMessagesGeneric:     opts.Bus.ChatMessagesGeneric,
		processGenericMessage:   opts.Bus.Parser.ProcessGenericMessage,
		eventsFollow:            opts.Bus.Events.Follow,
		eventsSubscribe:         opts.Bus.Events.Subscribe,
		eventsReSubscribe:       opts.Bus.Events.ReSubscribe,
		eventsSubGift:           opts.Bus.Events.SubGift,
		eventsRedemptionCreated: opts.Bus.Events.RedemptionCreated,
		eventsChannelBan:        opts.Bus.Events.ChannelBan,
		streamOnline:            opts.Bus.KickStreamOnline,
		streamOffline:           opts.Bus.KickStreamOffline,
		channelsRepo:            opts.ChannelsRepo,
		usersRepo:               opts.UsersRepo,
		kickBotsRepo:            opts.KickBotsRepo,
		eventsListRepo:          opts.EventsListRepo,
		channelsInfoHistoryRepo: opts.ChannelsInfoHistoryRepo,
		redemptionsHistoryRepo:  opts.RedemptionsHistoryRepo,
		streamsRepo:             opts.StreamsRepo,
		userCreatorService:      opts.UserCreatorService,
		prefixCache:             opts.PrefixCache,
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
	case "channel.subscription.new":
		eventLogAttrs, err = h.handleSubscriptionNew(r, body)
	case "channel.subscription.renewal":
		eventLogAttrs, err = h.handleSubscriptionRenewal(r, body)
	case "channel.subscription.gifts":
		eventLogAttrs, err = h.handleSubscriptionGifts(r, body)
	case "channel.reward.redemption.updated":
		eventLogAttrs, err = h.handleRewardRedemptionUpdated(r, body)
	case "livestream.status.updated":
		eventLogAttrs, err = h.handleLivestreamStatus(r, body)
	case "livestream.metadata.updated":
		eventLogAttrs, err = h.handleLivestreamMetadata(r, body)
	case "moderation.banned":
		eventLogAttrs, err = h.handleModerationBanned(r, body)
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
	channelUUID, _, err := h.resolveIDs(r, broadcasterUserID)
	if err != nil {
		return nil, fmt.Errorf("resolve ids for livestream.metadata.updated broadcaster_user_id=%s: %w", broadcasterUserID, err)
	}
	channelID := channelUUID.String()

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
			Platform:  platform.PlatformKick,
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
	channelUUID, _, err := h.resolveIDs(r, broadcasterUserID)
	if err != nil {
		return nil, fmt.Errorf("resolve ids for chat message broadcaster_user_id=%s: %w", broadcasterUserID, err)
	}
	channelID := channelUUID.String()

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
		channel        channelsmodel.Channel
		stream         *streamsmodel.Stream
		commandsPrefix string
	)

	var errwg errgroup.Group

	errwg.Go(func() error {
		var err error
		channel, err = h.channelsRepo.GetByID(ctx, channelUUID)
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
		slog.String("sender_twir_user_id", senderUser.ID.String()),
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
				UserID:            senderUser.ID.String(),
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
		UserID:            senderUser.ID.String(),
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
				ID:                senderUser.ID.String(),
				TokenID:           senderUser.TokenID.Ptr(),
				IsBotAdmin:        senderUser.IsBotAdmin,
				ApiKey:            senderUser.ApiKey,
				IsBanned:          senderUser.IsBanned,
				HideOnLandingPage: senderUser.HideOnLandingPage,
				CreatedAt:         senderUser.CreatedAt,
			},
			DbUserChannelStat: &generic.DbUserChannelStat{
				ID:                senderStats.ID,
				UserID:            senderStats.UserID.String(),
				ChannelID:         senderStats.ChannelID.String(),
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

	return bot.KickUserID.String() == senderUser.ID.String()
}

func (h *Handlers) handleChannelFollow(r *http.Request, body []byte) ([]slog.Attr, error) {
	ctx := r.Context()

	var payload kickFollowPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal follow payload: %w", err)
	}

	broadcasterUserID := strconv.Itoa(payload.Broadcaster.UserID)
	channelUUID, _, err := h.resolveIDs(r, broadcasterUserID)
	if err != nil {
		return nil, fmt.Errorf("resolve ids for follow broadcaster_user_id=%s: %w", broadcasterUserID, err)
	}
	channelID := channelUUID.String()

	if err := h.eventsFollow.Publish(
		ctx,
		events.FollowMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   channelID,
				ChannelName: kickChannelName(payload.Broadcaster),
				Platform:    platform.PlatformKick,
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

func normalizeKickRedemptionStatus(status string) string {
	return strings.ToLower(strings.TrimSpace(status))
}

func isKickRedemptionPending(status string) bool {
	switch normalizeKickRedemptionStatus(status) {
	case "pending", "new":
		return true
	default:
		return false
	}
}

func kickRewardHistoryUUID(rewardID string) uuid.UUID {
	if parsed, err := uuid.Parse(rewardID); err == nil {
		return parsed
	}

	return uuid.NewSHA1(uuid.NameSpaceOID, []byte("kick-reward:"+rewardID))
}

func kickChannelName(user kickUser) string {
	if user.ChannelSlug != "" {
		return user.ChannelSlug
	}

	return user.Username
}

func (h *Handlers) handleSubscriptionNew(r *http.Request, body []byte) ([]slog.Attr, error) {
	ctx := r.Context()

	var payload kickSubscriptionPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal channel.subscription.new payload: %w", err)
	}

	broadcasterUserID := strconv.Itoa(payload.Broadcaster.UserID)
	channelUUID, _, err := h.resolveIDs(r, broadcasterUserID)
	if err != nil {
		return nil, fmt.Errorf("resolve ids for channel.subscription.new broadcaster_user_id=%s: %w", broadcasterUserID, err)
	}
	channelID := channelUUID.String()

	subscriberUserID := strconv.Itoa(payload.Subscriber.UserID)
	if h.eventsListRepo != nil {
		if err := h.eventsListRepo.Create(ctx, channelseventslist.CreateInput{
			ChannelID: channelID,
			UserID:    &subscriberUserID,
			Platform:  platform.PlatformKick,
			Type:      channelseventslistmodel.ChannelEventListItemTypeSubscribe,
			Data: &channelseventslistmodel.ChannelsEventsListItemData{
				SubUserName:        payload.Subscriber.Username,
				SubUserDisplayName: payload.Subscriber.Username,
			},
		}); err != nil {
			return nil, fmt.Errorf("create kick subscribe event list item: %w", err)
		}
	}

	if err := h.eventsSubscribe.Publish(ctx, events.SubscribeMessage{
		BaseInfo: events.BaseInfo{
			ChannelID:   channelID,
			ChannelName: kickChannelName(payload.Broadcaster),
			Platform:    platform.PlatformKick,
		},
		UserID:          subscriberUserID,
		UserName:        payload.Subscriber.Username,
		UserDisplayName: payload.Subscriber.Username,
		Level:           "",
	}); err != nil {
		return nil, fmt.Errorf("publish kick subscribe event: %w", err)
	}

	return []slog.Attr{
		slog.String("channel_id", channelID),
		slog.String("broadcaster_user_id", broadcasterUserID),
		slog.String("subscriber_user_id", subscriberUserID),
		slog.String("subscriber_username", payload.Subscriber.Username),
		slog.Int("duration", payload.Duration),
	}, nil
}

func (h *Handlers) handleSubscriptionRenewal(r *http.Request, body []byte) ([]slog.Attr, error) {
	ctx := r.Context()

	var payload kickSubscriptionPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal channel.subscription.renewal payload: %w", err)
	}

	broadcasterUserID := strconv.Itoa(payload.Broadcaster.UserID)
	channelUUID, _, err := h.resolveIDs(r, broadcasterUserID)
	if err != nil {
		return nil, fmt.Errorf("resolve ids for channel.subscription.renewal broadcaster_user_id=%s: %w", broadcasterUserID, err)
	}
	channelID := channelUUID.String()

	subscriberUserID := strconv.Itoa(payload.Subscriber.UserID)
	months := max(payload.Duration, 0)
	if h.eventsListRepo != nil {
		if err := h.eventsListRepo.Create(ctx, channelseventslist.CreateInput{
			ChannelID: channelID,
			UserID:    &subscriberUserID,
			Platform:  platform.PlatformKick,
			Type:      channelseventslistmodel.ChannelEventListItemTypeReSubscribe,
			Data: &channelseventslistmodel.ChannelsEventsListItemData{
				ReSubUserName:        payload.Subscriber.Username,
				ReSubUserDisplayName: payload.Subscriber.Username,
				ReSubMonths:          strconv.Itoa(months),
				ReSubStreak:          "0",
			},
		}); err != nil {
			return nil, fmt.Errorf("create kick resubscribe event list item: %w", err)
		}
	}

	if err := h.eventsReSubscribe.Publish(ctx, events.ReSubscribeMessage{
		BaseInfo: events.BaseInfo{
			ChannelID:   channelID,
			ChannelName: kickChannelName(payload.Broadcaster),
			Platform:    platform.PlatformKick,
		},
		UserID:          subscriberUserID,
		UserName:        payload.Subscriber.Username,
		UserDisplayName: payload.Subscriber.Username,
		Months:          int64(months),
		Streak:          0,
		IsPrime:         false,
		Message:         "",
		Level:           "",
	}); err != nil {
		return nil, fmt.Errorf("publish kick resubscribe event: %w", err)
	}

	return []slog.Attr{
		slog.String("channel_id", channelID),
		slog.String("broadcaster_user_id", broadcasterUserID),
		slog.String("subscriber_user_id", subscriberUserID),
		slog.String("subscriber_username", payload.Subscriber.Username),
		slog.Int("duration", payload.Duration),
	}, nil
}

func (h *Handlers) handleSubscriptionGifts(r *http.Request, body []byte) ([]slog.Attr, error) {
	ctx := r.Context()

	var payload kickSubscriptionGiftsPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal channel.subscription.gifts payload: %w", err)
	}

	broadcasterUserID := strconv.Itoa(payload.Broadcaster.UserID)
	channelUUID, _, err := h.resolveIDs(r, broadcasterUserID)
	if err != nil {
		return nil, fmt.Errorf("resolve ids for channel.subscription.gifts broadcaster_user_id=%s: %w", broadcasterUserID, err)
	}
	channelID := channelUUID.String()

	gifterUserID := strconv.Itoa(payload.Gifter.UserID)
	for _, giftee := range payload.Giftees {
		gifteeUserID := strconv.Itoa(giftee.UserID)

		if h.eventsListRepo != nil {
			if err := h.eventsListRepo.Create(ctx, channelseventslist.CreateInput{
				ChannelID: channelID,
				UserID:    &gifterUserID,
				Platform:  platform.PlatformKick,
				Type:      channelseventslistmodel.ChannelEventListItemTypeSubGift,
				Data: &channelseventslistmodel.ChannelsEventsListItemData{
					SubGiftUserName:              payload.Gifter.Username,
					SubGiftUserDisplayName:       payload.Gifter.Username,
					SubGiftTargetUserName:        giftee.Username,
					SubGiftTargetUserDisplayName: giftee.Username,
				},
			}); err != nil {
				return nil, fmt.Errorf("create kick subgift event list item for giftee %s: %w", gifteeUserID, err)
			}
		}

		if err := h.eventsSubGift.Publish(ctx, events.SubGiftMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   channelID,
				ChannelName: kickChannelName(payload.Broadcaster),
				Platform:    platform.PlatformKick,
			},
			SenderUserID:      gifterUserID,
			SenderUserName:    payload.Gifter.Username,
			SenderDisplayName: payload.Gifter.Username,
			TargetUserName:    giftee.Username,
			TargetDisplayName: giftee.Username,
			Level:             "",
		}); err != nil {
			return nil, fmt.Errorf("publish kick subgift event for giftee %s: %w", gifteeUserID, err)
		}
	}

	return []slog.Attr{
		slog.String("channel_id", channelID),
		slog.String("broadcaster_user_id", broadcasterUserID),
		slog.String("gifter_user_id", gifterUserID),
		slog.String("gifter_username", payload.Gifter.Username),
		slog.Int("giftee_count", len(payload.Giftees)),
	}, nil
}

func (h *Handlers) handleRewardRedemptionUpdated(r *http.Request, body []byte) ([]slog.Attr, error) {
	ctx := r.Context()

	var payload kickRewardRedemptionPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal channel.reward.redemption.updated payload: %w", err)
	}

	status := normalizeKickRedemptionStatus(payload.Status)
	broadcasterUserID := strconv.Itoa(payload.Broadcaster.UserID)
	channelUUID, _, err := h.resolveIDs(r, broadcasterUserID)
	if err != nil {
		return nil, fmt.Errorf("resolve ids for channel.reward.redemption.updated broadcaster_user_id=%s: %w", broadcasterUserID, err)
	}
	channelID := channelUUID.String()

	if !isKickRedemptionPending(payload.Status) {
		return []slog.Attr{
			slog.String("channel_id", channelID),
			slog.String("broadcaster_user_id", broadcasterUserID),
			slog.String("reward_id", payload.Reward.ID),
			slog.String("status", status),
			slog.Bool("ignored", true),
		}, nil
	}

	redeemerUserID := strconv.Itoa(payload.Redeemer.UserID)
	if h.redemptionsHistoryRepo != nil {
		if err := h.redemptionsHistoryRepo.Create(ctx, channelsredemptionshistory.CreateInput{
			ChannelID:    channelID,
			UserID:       redeemerUserID,
			Platform:     platform.PlatformKick,
			RewardID:     kickRewardHistoryUUID(payload.Reward.ID),
			RewardPrompt: lo.If(payload.UserInput != "", &payload.UserInput).Else(nil),
			RewardTitle:  payload.Reward.Title,
			RewardCost:   payload.Reward.Cost,
		}); err != nil {
			return nil, fmt.Errorf("create kick redemption history: %w", err)
		}
	}

	if h.eventsListRepo != nil {
		if err := h.eventsListRepo.Create(ctx, channelseventslist.CreateInput{
			ChannelID: channelID,
			UserID:    &redeemerUserID,
			Platform:  platform.PlatformKick,
			Type:      channelseventslistmodel.ChannelEventListItemTypeRedemptionCreated,
			Data: &channelseventslistmodel.ChannelsEventsListItemData{
				RedemptionInput:           payload.UserInput,
				RedemptionTitle:           payload.Reward.Title,
				RedemptionUserName:        payload.Redeemer.Username,
				RedemptionUserDisplayName: payload.Redeemer.Username,
				RedemptionCost:            strconv.Itoa(payload.Reward.Cost),
			},
		}); err != nil {
			return nil, fmt.Errorf("create kick redemption event list item: %w", err)
		}
	}

	if err := h.eventsRedemptionCreated.Publish(ctx, events.RedemptionCreatedMessage{
		ID: payload.Reward.ID,
		BaseInfo: events.BaseInfo{
			ChannelID:   channelID,
			ChannelName: kickChannelName(kickUser{Username: payload.Broadcaster.Username, ChannelSlug: payload.Broadcaster.ChannelSlug}),
			Platform:    platform.PlatformKick,
		},
		UserID:          redeemerUserID,
		UserName:        payload.Redeemer.Username,
		UserDisplayName: payload.Redeemer.Username,
		RewardName:      payload.Reward.Title,
		RewardCost:      strconv.Itoa(payload.Reward.Cost),
		Input:           lo.If(payload.UserInput != "", &payload.UserInput).Else(nil),
	}); err != nil {
		return nil, fmt.Errorf("publish kick redemption created event: %w", err)
	}

	return []slog.Attr{
		slog.String("channel_id", channelID),
		slog.String("broadcaster_user_id", broadcasterUserID),
		slog.String("redeemer_user_id", redeemerUserID),
		slog.String("redeemer_username", payload.Redeemer.Username),
		slog.String("reward_id", payload.Reward.ID),
		slog.String("reward_title", payload.Reward.Title),
		slog.String("status", status),
	}, nil
}

func (h *Handlers) handleModerationBanned(r *http.Request, body []byte) ([]slog.Attr, error) {
	ctx := r.Context()

	var payload kickModerationBannedPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal moderation.banned payload: %w", err)
	}

	broadcasterUserID := strconv.Itoa(payload.Broadcaster.UserID)
	channelUUID, _, err := h.resolveIDs(r, broadcasterUserID)
	if err != nil {
		return nil, fmt.Errorf("resolve ids for moderation.banned broadcaster_user_id=%s: %w", broadcasterUserID, err)
	}
	channelID := channelUUID.String()

	bannedUserID := strconv.Itoa(payload.BannedUser.UserID)
	moderatorUserID := strconv.Itoa(payload.Moderator.UserID)
	endsAt := "permanent"
	isPermanent := true
	if payload.Metadata.ExpiresAt != nil && *payload.Metadata.ExpiresAt != "" {
		isPermanent = false
		if expiresAt, err := time.Parse(time.RFC3339, *payload.Metadata.ExpiresAt); err == nil {
			minutes := int(time.Until(expiresAt).Round(time.Minute).Minutes())
			if minutes <= 0 {
				minutes = 1
			}
			endsAt = strconv.Itoa(minutes)
		}
	}

	if err := h.eventsChannelBan.Publish(ctx, events.ChannelBanMessage{
		BaseInfo: events.BaseInfo{
			ChannelID:   channelID,
			ChannelName: kickChannelName(payload.Broadcaster),
			Platform:    platform.PlatformKick,
		},
		UserID:               bannedUserID,
		UserName:             payload.BannedUser.Username,
		UserLogin:            payload.BannedUser.Username,
		BroadcasterUserName:  payload.Broadcaster.Username,
		BroadcasterUserLogin: payload.Broadcaster.Username,
		ModeratorUserID:      moderatorUserID,
		ModeratorUserName:    payload.Moderator.Username,
		ModeratorUserLogin:   payload.Moderator.Username,
		Reason:               payload.Metadata.Reason,
		EndsAt:               endsAt,
		IsPermanent:          isPermanent,
	}); err != nil {
		return nil, fmt.Errorf("publish kick channel ban event: %w", err)
	}

	if h.eventsListRepo != nil {
		if err := h.eventsListRepo.Create(ctx, channelseventslist.CreateInput{
			ChannelID: channelID,
			UserID:    &bannedUserID,
			Platform:  platform.PlatformKick,
			Type:      channelseventslistmodel.ChannelEventListItemTypeChannelBan,
			Data: &channelseventslistmodel.ChannelsEventsListItemData{
				BanReason:            payload.Metadata.Reason,
				BanEndsInMinutes:     endsAt,
				BannedUserLogin:      payload.BannedUser.Username,
				BannedUserName:       payload.BannedUser.Username,
				ModeratorDisplayName: payload.Moderator.Username,
				ModeratorName:        payload.Moderator.Username,
			},
		}); err != nil {
			return nil, fmt.Errorf("create kick channel ban event list item: %w", err)
		}
	}

	return []slog.Attr{
		slog.String("channel_id", channelID),
		slog.String("broadcaster_user_id", broadcasterUserID),
		slog.String("moderator_user_id", moderatorUserID),
		slog.String("banned_user_id", bannedUserID),
		slog.String("reason", payload.Metadata.Reason),
		slog.Bool("is_permanent", isPermanent),
	}, nil
}

func (h *Handlers) handleLivestreamStatus(r *http.Request, body []byte) ([]slog.Attr, error) {
	ctx := r.Context()

	var payload kickLivestreamStatusPayload
	if err := json.Unmarshal(body, &payload); err != nil {
		return nil, fmt.Errorf("unmarshal livestream.status.updated payload: %w", err)
	}

	broadcasterUserID := strconv.Itoa(payload.Broadcaster.UserID)
	channelUUID, _, err := h.resolveIDs(r, broadcasterUserID)
	if err != nil {
		return nil, fmt.Errorf("resolve ids for livestream.status.updated broadcaster_user_id=%s: %w", broadcasterUserID, err)
	}
	channelID := channelUUID.String()

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

func (h *Handlers) resolveIDs(r *http.Request, broadcasterUserID string) (uuid.UUID, uuid.UUID, error) {
	ctx := r.Context()

	user, err := h.usersRepo.GetByPlatformID(ctx, platform.PlatformKick, broadcasterUserID)
	if err != nil {
		if errors.Is(err, usersmodel.ErrNotFound) {
			return uuid.Nil, uuid.Nil, fmt.Errorf("no kick user for broadcaster_user_id=%s", broadcasterUserID)
		}
		return uuid.Nil, uuid.Nil, fmt.Errorf("get user by platform id: %w", err)
	}

	channel, err := h.channelsRepo.GetByKickUserID(ctx, user.ID)
	if err != nil {
		if errors.Is(err, channelsrepository.ErrNotFound) {
			return uuid.Nil, uuid.Nil, fmt.Errorf("channel not found for user_id=%s platform=kick", user.ID)
		}
		return uuid.Nil, uuid.Nil, fmt.Errorf("get channel by kick user id: %w", err)
	}

	return channel.ID, user.ID, nil
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
