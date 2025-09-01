package chat_alerts

import (
	"context"
	"errors"
	"log/slog"

	"github.com/redis/go-redis/v9"
	buscore "github.com/twirapp/twir/libs/bus-core"
	busevents "github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	chatalertscache "github.com/twirapp/twir/libs/cache/chatalerts"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/logger"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/events/model"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type ChatAlerts struct {
	db                    *gorm.DB
	redis                 *redis.Client
	logger                logger.Logger
	cfg                   cfg.Config
	websocketsGrpc        websockets.WebsocketClient
	bus                   *buscore.Bus
	chatAlertsCache       *generic_cacher.GenericCacher[chatalertscache.ChatAlert]
	channelEventListsRepo channelseventslist.Repository
}

type Opts struct {
	fx.In

	DB                    *gorm.DB
	Redis                 *redis.Client
	Logger                logger.Logger
	Cfg                   cfg.Config
	WebsocketsGrpc        websockets.WebsocketClient
	Bus                   *buscore.Bus
	ChatAlertsCache       *generic_cacher.GenericCacher[chatalertscache.ChatAlert]
	ChannelEventListsRepo channelseventslist.Repository
}

func New(opts Opts) (*ChatAlerts, error) {
	return &ChatAlerts{
		db:                    opts.DB,
		redis:                 opts.Redis,
		logger:                opts.Logger,
		cfg:                   opts.Cfg,
		bus:                   opts.Bus,
		websocketsGrpc:        opts.WebsocketsGrpc,
		chatAlertsCache:       opts.ChatAlertsCache,
		channelEventListsRepo: opts.ChannelEventListsRepo,
	}, nil
}

func (c *ChatAlerts) ProcessEvent(
	ctx context.Context,
	channelId string,
	eventType model.EventType,
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
	case model.EventTypeFollow:
		casted, ok := data.(busevents.FollowMessage)
		if ok {
			cooldown = parsedSettings.Followers.Cooldown
			processErr = c.follow(ctx, parsedSettings, casted)
		}
	case model.EventTypeChannelBan:
		casted, ok := data.(busevents.ChannelBanMessage)
		if ok {
			cooldown = parsedSettings.Ban.Cooldown
			processErr = c.ban(ctx, parsedSettings, casted)
		}
	case model.EventTypeOnChatClear:
		casted, ok := data.(busevents.ChatClearMessage)
		if ok {
			cooldown = parsedSettings.ChatCleared.Cooldown
			processErr = c.chatCleared(ctx, parsedSettings, casted)
		}
	case model.EventTypeDonate:
		casted, ok := data.(busevents.DonateMessage)
		if ok {
			cooldown = parsedSettings.Donations.Cooldown
			processErr = c.donation(ctx, parsedSettings, casted)
		}
	case model.EventTypeRaided:
		casted, ok := data.(busevents.RaidedMessage)
		if ok {
			cooldown = parsedSettings.Raids.Cooldown
			processErr = c.raid(ctx, parsedSettings, casted)
		}
	case model.EventTypeRedemptionCreated:
		casted, ok := data.(busevents.RedemptionCreatedMessage)
		if ok {
			cooldown = parsedSettings.Redemptions.Cooldown
			processErr = c.redemption(ctx, parsedSettings, casted)
		}
	case model.EventTypeStreamOffline:
		casted, ok := data.(twitch.StreamOfflineMessage)
		if ok {
			cooldown = parsedSettings.StreamOffline.Cooldown
			processErr = c.streamOffline(ctx, parsedSettings, casted)
		}
	case model.EventTypeStreamOnline:
		casted, ok := data.(twitch.StreamOnlineMessage)
		if ok {
			cooldown = parsedSettings.StreamOnline.Cooldown
			processErr = c.streamOnline(ctx, parsedSettings, casted)
		}
	case model.EventTypeSubscribe:
		casted, ok := data.(SubscribeMessage)
		if ok {
			cooldown = parsedSettings.Subscribers.Cooldown
			processErr = c.subscribe(ctx, parsedSettings, casted)
		}
	case model.EventTypeFirstUserMessage:
		casted, ok := data.(busevents.FirstUserMessageMessage)
		if ok {
			cooldown = parsedSettings.FirstUserMessage.Cooldown
			processErr = c.firstUserMessage(ctx, parsedSettings, casted)
		}
	case model.EventTypeChannelUnbanRequestCreate:
		casted, ok := data.(busevents.ChannelUnbanRequestCreateMessage)
		if ok {
			cooldown = parsedSettings.Ban.Cooldown
			processErr = c.unbanRequestCreate(ctx, parsedSettings, casted)
		}
	case model.EventTypeChannelUnbanRequestResolve:
		casted, ok := data.(busevents.ChannelUnbanRequestResolveMessage)
		if ok {
			cooldown = parsedSettings.Ban.Cooldown
			processErr = c.unbanRequestResolved(ctx, parsedSettings, casted)
		}
	case model.EventTypeChannelMessageDelete:
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
