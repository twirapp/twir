package bus_listener

import (
	"context"
	"log/slog"

	"github.com/satont/twir/apps/bots/internal/messagehandler"
	mod_task_queue "github.com/satont/twir/apps/bots/internal/mod-task-queue"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	"github.com/satont/twir/apps/bots/internal/workers"
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

	Logger logger.Logger

	TokensGrpc         tokens.TokensClient
	Tracer             trace.Tracer
	ModTaskDistributor mod_task_queue.TaskDistributor

	Gorm           *gorm.DB
	TwitchActions  *twitchactions.TwitchActions
	MessageHandler *messagehandler.MessageHandler
	Bus            *bus_core.Bus
	Cfg            cfg.Config
	WorkersPool    *workers.Pool
}

func New(opts Opts) (*BusListener, error) {
	listener := &BusListener{
		gorm:               opts.Gorm,
		logger:             opts.Logger,
		config:             opts.Cfg,
		tokensGrpc:         opts.TokensGrpc,
		twitchActions:      opts.TwitchActions,
		messageHandler:     opts.MessageHandler,
		tracer:             opts.Tracer,
		bus:                opts.Bus,
		modTaskDistributor: opts.ModTaskDistributor,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				listener.bus.Bots.SendMessage.SubscribeGroup(
					"bots",
					func(ctx context.Context, data bots.SendMessageRequest) struct{} {
						opts.WorkersPool.Submit(
							func() {
								listener.sendMessage(ctx, data)
							},
						).Wait()

						return struct{}{}
					},
				)
				listener.bus.Bots.DeleteMessage.SubscribeGroup(
					"bots",
					func(ctx context.Context, data bots.DeleteMessageRequest) struct{} {
						opts.WorkersPool.Submit(
							func() {
								listener.deleteMessage(ctx, data)
							},
						).Wait()

						return struct{}{}
					},
				)
				listener.bus.ChatMessages.SubscribeGroup(
					"bots",
					listener.handleChatMessage,
				)
				listener.bus.Bots.BanUser.SubscribeGroup(
					"bots",
					func(ctx context.Context, data bots.BanRequest) struct{} {
						opts.WorkersPool.Submit(
							func() {
								listener.banUser(ctx, data)
							},
						).Wait()

						return struct{}{}
					},
				)
				listener.bus.Bots.BanUsers.SubscribeGroup(
					"bots",
					func(ctx context.Context, data []bots.BanRequest) struct{} {
						opts.WorkersPool.Submit(
							func() {
								listener.banUsers(ctx, data)
							},
						).Wait()

						return struct{}{}
					},
				)
				listener.bus.Bots.ShoutOut.SubscribeGroup(
					"bots",
					func(ctx context.Context, data bots.SentShoutOutRequest) struct{} {
						opts.WorkersPool.Submit(
							func() {
								listener.handleShoutOut(ctx, data)
							},
						).Wait()

						return struct{}{}
					},
				)

				return nil
			},
			OnStop: func(ctx context.Context) error {
				listener.bus.Bots.SendMessage.Unsubscribe()
				listener.bus.ChatMessages.Unsubscribe()
				listener.bus.Bots.DeleteMessage.Unsubscribe()
				listener.bus.Bots.BanUser.Unsubscribe()
				listener.bus.Bots.ShoutOut.Unsubscribe()

				return nil
			},
		},
	)

	return listener, nil
}

type BusListener struct {
	logger             logger.Logger
	tokensGrpc         tokens.TokensClient
	tracer             trace.Tracer
	modTaskDistributor mod_task_queue.TaskDistributor
	gorm               *gorm.DB
	twitchActions      *twitchactions.TwitchActions
	messageHandler     *messagehandler.MessageHandler
	bus                *bus_core.Bus
	config             cfg.Config
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
				e := c.twitchActions.DeleteMessage(
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
	if req.ChannelId == "" {
		return struct{}{}
	}

	err := c.twitchActions.SendMessage(
		ctx,
		twitchactions.SendMessageOpts{
			BroadcasterID:        req.ChannelId,
			SenderID:             "",
			Message:              req.Message,
			ReplyParentMessageID: req.ReplyTo,
			IsAnnounce:           req.IsAnnounce,
			SkipToxicityCheck:    req.SkipToxicityCheck,
			SkipRateLimits:       req.SkipRateLimits,
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

func (c *BusListener) handleShoutOut(ctx context.Context, req bots.SentShoutOutRequest) struct{} {
	err := c.twitchActions.ShoutOut(
		ctx,
		twitchactions.ShoutOutInput{
			BroadcasterID: req.ChannelID,
			TargetID:      req.TargetID,
		},
	)
	if err != nil {
		c.logger.Error("cannot send shoutout", slog.Any("err", err))
	}
	return struct{}{}
}
