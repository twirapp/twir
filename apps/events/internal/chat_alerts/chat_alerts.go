package chat_alerts

import (
	"context"
	"errors"
	"log/slog"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	cfg "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/libs/api/messages/events"
	buscore "github.com/twirapp/twir/libs/bus-core"
	events_messages "github.com/twirapp/twir/libs/grpc/events"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type ChatAlerts struct {
	db             *gorm.DB
	redis          *redis.Client
	logger         logger.Logger
	cfg            cfg.Config
	tokensGrpc     tokens.TokensClient
	websocketsGrpc websockets.WebsocketClient
	bus            *buscore.Bus
}

type Opts struct {
	fx.In

	DB             *gorm.DB
	Redis          *redis.Client
	Logger         logger.Logger
	Cfg            cfg.Config
	TokensGrpc     tokens.TokensClient
	WebsocketsGrpc websockets.WebsocketClient
	Bus            *buscore.Bus
}

func New(opts Opts) (*ChatAlerts, error) {
	return &ChatAlerts{
		db:             opts.DB,
		redis:          opts.Redis,
		logger:         opts.Logger,
		cfg:            opts.Cfg,
		bus:            opts.Bus,
		tokensGrpc:     opts.TokensGrpc,
		websocketsGrpc: opts.WebsocketsGrpc,
	}, nil
}

func (c *ChatAlerts) ProcessEvent(
	ctx context.Context,
	channelId string,
	eventType events.TwirEventType,
	data any,
) {
	entity := model.ChannelModulesSettings{}

	if err := c.db.
		WithContext(ctx).
		Where(
			`"channelId" = ? AND "userId" IS NULL AND type = 'chat_alerts'`,
			channelId,
		).First(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return
		}

		c.logger.Error("cannot find channel chat alerts settings", slog.Any("err", err))

		return
	}

	parsedSettings := model.ChatAlertsSettings{}
	if err := json.Unmarshal(entity.Settings, &parsedSettings); err != nil {
		c.logger.Error("failed to unmarshal settings: %w", slog.Any("err", err))
		return
	}

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
		casted, ok := data.(*events_messages.FollowMessage)
		if ok {
			cooldown = parsedSettings.Followers.Cooldown
			processErr = c.follow(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_USER_BANNED:
		casted, ok := data.(*events_messages.ChannelBanMessage)
		if ok {
			cooldown = parsedSettings.Ban.Cooldown
			processErr = c.ban(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_CHAT_CLEAR:
		casted, ok := data.(*events_messages.ChatClearMessage)
		if ok {
			cooldown = parsedSettings.ChatCleared.Cooldown
			processErr = c.chatCleared(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_DONATE:
		casted, ok := data.(*events_messages.DonateMessage)
		if ok {
			cooldown = parsedSettings.Donations.Cooldown
			processErr = c.donation(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_RAIDED:
		casted, ok := data.(*events_messages.RaidedMessage)
		if ok {
			cooldown = parsedSettings.Raids.Cooldown
			processErr = c.raid(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_REDEMPTION_CREATED:
		casted, ok := data.(*events_messages.RedemptionCreatedMessage)
		if ok {
			cooldown = parsedSettings.Redemptions.Cooldown
			processErr = c.redemption(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_STREAM_OFFLINE:
		casted, ok := data.(*events_messages.StreamOfflineMessage)
		if ok {
			cooldown = parsedSettings.StreamOffline.Cooldown
			processErr = c.streamOffline(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_STREAM_ONLINE:
		casted, ok := data.(*events_messages.StreamOnlineMessage)
		if ok {
			cooldown = parsedSettings.StreamOnline.Cooldown
			processErr = c.streamOnline(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_SUBSCRIBE:
		casted, ok := data.(*SubscribMessage)
		if ok {
			cooldown = parsedSettings.Subscribers.Cooldown
			processErr = c.subscribe(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_FIRST_USER_MESSAGE:
		casted, ok := data.(*events_messages.FirstUserMessageMessage)
		if ok {
			cooldown = parsedSettings.FirstUserMessage.Cooldown
			processErr = c.firstUserMessage(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_CHANNEL_UNBAN_REQUEST_CREATED:
		casted, ok := data.(*events_messages.ChannelUnbanRequestCreateMessage)
		if ok {
			cooldown = parsedSettings.Ban.Cooldown
			processErr = c.unbanRequestCreate(ctx, parsedSettings, casted)
		}
	case events.TwirEventType_CHANNEL_UNBAN_REQUEST_RESOLVED:
		casted, ok := data.(*events_messages.ChannelUnbanRequestResolveMessage)
		if ok {
			cooldown = parsedSettings.Ban.Cooldown
			processErr = c.unbanRequestResolved(ctx, parsedSettings, casted)
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
