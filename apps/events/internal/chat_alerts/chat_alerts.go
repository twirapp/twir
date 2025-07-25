package chat_alerts

import (
	"context"
	"errors"
	"log/slog"

	"github.com/redis/go-redis/v9"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/api/messages/events"
	buscore "github.com/twirapp/twir/libs/bus-core"
	busevents "github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	chatalertscache "github.com/twirapp/twir/libs/cache/chatalerts"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type ChatAlerts struct {
	db              *gorm.DB
	redis           *redis.Client
	logger          logger.Logger
	cfg             cfg.Config
	websocketsGrpc  websockets.WebsocketClient
	bus             *buscore.Bus
	chatAlertsCache *generic_cacher.GenericCacher[chatalertscache.ChatAlert]
}

type Opts struct {
	fx.In

	DB              *gorm.DB
	Redis           *redis.Client
	Logger          logger.Logger
	Cfg             cfg.Config
	WebsocketsGrpc  websockets.WebsocketClient
	Bus             *buscore.Bus
	ChatAlertsCache *generic_cacher.GenericCacher[chatalertscache.ChatAlert]
}

func New(opts Opts) (*ChatAlerts, error) {
	return &ChatAlerts{
		db:              opts.DB,
		redis:           opts.Redis,
		logger:          opts.Logger,
		cfg:             opts.Cfg,
		bus:             opts.Bus,
		websocketsGrpc:  opts.WebsocketsGrpc,
		chatAlertsCache: opts.ChatAlertsCache,
	}, nil
}

func (c *ChatAlerts) ProcessEvent(
	ctx context.Context,
	channelId string,
	eventType events.TwirEventType,
	data any,
) {
	entity, err := c.chatAlertsCache.Get(ctx, channelId)
	if err != nil {
		if errors.Is(err, chatalertscache.ErrChatAlertNotFound) {
			return
		}

		c.logger.Error("cannot get chat alerts", slog.Any("err", err))
		return
	}

	parsedSettings := entity.ParsedSettings

	eventCooldown, err := c.isOnCooldown(
		ctx,
		channelId,
		eventType.String(),
	)
	if err != nil {
		c.logger.Error("cannot get channel event cooldown", slog.Any("err", err))
		return
	}
	if eventCooldown {
		return
	}

	var cooldown int
	var processErr error

	switch eventType {
	case events.TwirEventType_FOLLOW:
		casted, ok := data.(busevents.FollowMessage)
		if ok {
			cooldown = parsedSettings.Followers.Cooldown
			processErr = c.follow(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_USER_BANNED:
		casted, ok := data.(busevents.ChannelBanMessage)
		if ok {
			cooldown = parsedSettings.Ban.Cooldown
			processErr = c.ban(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_CHAT_CLEAR:
		casted, ok := data.(busevents.ChatClearMessage)
		if ok {
			cooldown = parsedSettings.ChatCleared.Cooldown
			processErr = c.chatCleared(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_DONATE:
		casted, ok := data.(busevents.DonateMessage)
		if ok {
			cooldown = parsedSettings.Donations.Cooldown
			processErr = c.donation(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_RAIDED:
		casted, ok := data.(busevents.RaidedMessage)
		if ok {
			cooldown = parsedSettings.Raids.Cooldown
			processErr = c.raid(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_REDEMPTION_CREATED:
		casted, ok := data.(busevents.RedemptionCreatedMessage)
		if ok {
			cooldown = parsedSettings.Redemptions.Cooldown
			processErr = c.redemption(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_STREAM_OFFLINE:
		casted, ok := data.(twitch.StreamOfflineMessage)
		if ok {
			cooldown = parsedSettings.StreamOffline.Cooldown
			processErr = c.streamOffline(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_STREAM_ONLINE:
		casted, ok := data.(twitch.StreamOnlineMessage)
		if ok {
			cooldown = parsedSettings.StreamOnline.Cooldown
			processErr = c.streamOnline(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_SUBSCRIBE:
		casted, ok := data.(SubscribeMessage)
		if ok {
			cooldown = parsedSettings.Subscribers.Cooldown
			processErr = c.subscribe(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_FIRST_USER_MESSAGE:
		casted, ok := data.(busevents.FirstUserMessageMessage)
		if ok {
			cooldown = parsedSettings.FirstUserMessage.Cooldown
			processErr = c.firstUserMessage(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_CHANNEL_UNBAN_REQUEST_CREATED:
		casted, ok := data.(busevents.ChannelUnbanRequestCreateMessage)
		if ok {
			cooldown = parsedSettings.Ban.Cooldown
			processErr = c.unbanRequestCreate(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_CHANNEL_UNBAN_REQUEST_RESOLVED:
		casted, ok := data.(busevents.ChannelUnbanRequestResolveMessage)
		if ok {
			cooldown = parsedSettings.Ban.Cooldown
			processErr = c.unbanRequestResolved(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_CHANNEL_MESSAGE_DELETE:
		casted, ok := data.(busevents.ChannelMessageDeleteMessage)
		if ok {
			cooldown = parsedSettings.MessageDelete.Cooldown
			processErr = c.messageDelete(ctx, parsedSettings, casted)
		}
	default:
		c.logger.Warn("unknown event", slog.Any("eventType", eventType))
	}

	if processErr != nil {
		c.logger.Error("cannot process event", slog.Any("err", processErr))
		return
	}

	if cooldown != 0 {
		err = c.SetCooldown(ctx, channelId, eventType.String(), cooldown)
		if err != nil {
			c.logger.Error("cannot set cooldown", slog.Any("err", err))
		}
	}

	return
}
