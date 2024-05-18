package bus_listener

import (
	"context"
	"log/slog"

	"github.com/satont/twir/apps/bots/internal/messagehandler"
	mod_task_queue "github.com/satont/twir/apps/bots/internal/mod-task-queue"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	cfg "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/utils"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/grpc/tokens"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Gorm   *gorm.DB
	Logger logger.Logger
	Cfg    cfg.Config

	TokensGrpc         tokens.TokensClient
	TwitchActions      *twitchactions.TwitchActions
	MessageHandler     *messagehandler.MessageHandler
	Tracer             trace.Tracer
	Bus                *bus_core.Bus
	ModTaskDistributor mod_task_queue.TaskDistributor
}

func New(opts Opts) (*BusListener, error) {
	listener := &BusListener{
		gorm:               opts.Gorm,
		logger:             opts.Logger,
		config:             opts.Cfg,
		tokensGrpc:         opts.TokensGrpc,
		twitchactions:      opts.TwitchActions,
		messageHandler:     opts.MessageHandler,
		tracer:             opts.Tracer,
		bus:                opts.Bus,
		modTaskDistributor: opts.ModTaskDistributor,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				listener.bus.Bots.SendMessage.SubscribeGroup("bots", listener.sendMessage)
				listener.bus.Bots.DeleteMessage.SubscribeGroup("bots", listener.deleteMessage)
				listener.bus.Bots.ProcessMessage.SubscribeGroup("bots", listener.handleChatMessage)
				listener.bus.Bots.BanUser.SubscribeGroup("bots", listener.banUser)

				return nil
			},
			OnStop: func(ctx context.Context) error {
				listener.bus.Bots.SendMessage.Unsubscribe()
				listener.bus.Bots.ProcessMessage.Unsubscribe()
				listener.bus.Bots.DeleteMessage.Unsubscribe()
				listener.bus.Bots.BanUser.Unsubscribe()

				return nil
			},
		},
	)

	return listener, nil
}

type BusListener struct {
	gorm               *gorm.DB
	logger             logger.Logger
	config             cfg.Config
	tokensGrpc         tokens.TokensClient
	twitchactions      *twitchactions.TwitchActions
	messageHandler     *messagehandler.MessageHandler
	tracer             trace.Tracer
	bus                *bus_core.Bus
	modTaskDistributor mod_task_queue.TaskDistributor
}

func (c *BusListener) deleteMessage(ctx context.Context, req bots.DeleteMessageRequest) struct{} {
	channel := model.Channels{}
	err := c.gorm.WithContext(ctx).Where("id = ?", req.ChannelId).Find(&channel).Error
	if err != nil {
		c.logger.Error(
			"cannot get channel",
			slog.String("channelId", req.ChannelId),
		)
		return struct{}{}
	}

	if channel.ID == "" {
		return struct{}{}
	}

	wg := utils.NewGoroutinesGroup()

	for _, m := range req.MessageIds {
		wg.Go(
			func() {
				e := c.twitchactions.DeleteMessage(
					ctx,
					twitchactions.DeleteMessageOpts{
						BroadcasterID: req.ChannelId,
						ModeratorID:   channel.BotID,
						MessageID:     m,
					},
				)
				if e != nil {
					c.logger.Error("cannot delete message", slog.Any("err", e))
				}
			},
		)
	}

	wg.Wait()

	return struct{}{}
}

func (c *BusListener) sendMessage(ctx context.Context, req bots.SendMessageRequest) struct{} {
	channel := model.Channels{}
	err := c.gorm.WithContext(ctx).Where("id = ?", req.ChannelId).Find(&channel).Error
	if err != nil {
		c.logger.Error(
			"cannot get channel",
			slog.String("channelId", req.ChannelId),
		)
		return struct{}{}
	}

	if channel.ID == "" {
		return struct{}{}
	}

	err = c.twitchactions.SendMessage(
		ctx,
		twitchactions.SendMessageOpts{
			BroadcasterID:        req.ChannelId,
			SenderID:             channel.BotID,
			Message:              req.Message,
			ReplyParentMessageID: req.ReplyTo,
			IsAnnounce:           req.IsAnnounce,
		},
	)
	if err != nil {
		c.logger.Error("cannot send message", slog.Any("err", err))
	}
	return struct{}{}
}

func (c *BusListener) handleChatMessage(
	ctx context.Context,
	req twitch.TwitchChatMessage,
) struct{} {
	span := trace.SpanFromContext(ctx)
	// End the span when the operation we are measuring is done.
	defer span.End()
	span.SetAttributes(
		attribute.String("message_id", req.MessageId),
		attribute.String("channel_id", req.BroadcasterUserId),
	)

	err := c.messageHandler.Handle(ctx, req)
	if err != nil {
		c.logger.Error(
			"cannot handle message",
			slog.String("channelId", req.BroadcasterUserId),
			slog.String("channelName", req.BroadcasterUserLogin),
			slog.Any("err", err),
		)
	}

	return struct{}{}
}
