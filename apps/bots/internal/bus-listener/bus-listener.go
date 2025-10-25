package bus_listener

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/twirapp/kv"
	"github.com/twirapp/twir/apps/bots/internal/messagehandler"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/apps/bots/internal/workers"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	cfg "github.com/twirapp/twir/libs/config"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/redis_keys"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Logger logger.Logger

	Tracer trace.Tracer

	Gorm           *gorm.DB
	TwitchActions  *twitchactions.TwitchActions
	MessageHandler *messagehandler.MessageHandler
	Bus            *bus_core.Bus
	Cfg            cfg.Config
	WorkersPool    *workers.Pool
	KV             kv.KV
}

func New(opts Opts) (*BusListener, error) {
	listener := &BusListener{
		gorm:           opts.Gorm,
		logger:         opts.Logger,
		config:         opts.Cfg,
		twitchActions:  opts.TwitchActions,
		messageHandler: opts.MessageHandler,
		tracer:         opts.Tracer,
		bus:            opts.Bus,
		kv:             opts.KV,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				err := listener.bus.Bots.SendMessage.SubscribeGroup(
					"bots",
					func(ctx context.Context, data bots.SendMessageRequest) (struct{}, error) {
						err := opts.WorkersPool.SubmitErr(
							func() error {
								return listener.sendMessage(ctx, data)
							},
						).Wait()

						return struct{}{}, err
					},
				)
				if err != nil {
					return err
				}

				err = listener.bus.Bots.DeleteMessage.SubscribeGroup(
					"bots",
					func(ctx context.Context, data bots.DeleteMessageRequest) (struct{}, error) {
						err := opts.WorkersPool.SubmitErr(
							func() error {
								return listener.deleteMessage(ctx, data)
							},
						).Wait()

						return struct{}{}, err
					},
				)
				if err != nil {
					return err
				}

				err = listener.bus.ChatMessages.SubscribeGroup(
					"bots",
					listener.handleChatMessage,
				)
				if err != nil {
					return err
				}

				err = listener.bus.Bots.BanUser.SubscribeGroup(
					"bots",
					func(ctx context.Context, data bots.BanRequest) (struct{}, error) {
						err := opts.WorkersPool.SubmitErr(
							func() error {
								return listener.banUser(ctx, data)
							},
						).Wait()

						return struct{}{}, err
					},
				)
				if err != nil {
					return err
				}

				err = listener.bus.Bots.BanUsers.SubscribeGroup(
					"bots",
					func(ctx context.Context, data []bots.BanRequest) (struct{}, error) {
						err := opts.WorkersPool.SubmitErr(
							func() error {
								return listener.banUsers(ctx, data)
							},
						).Wait()

						return struct{}{}, err
					},
				)
				if err != nil {
					return err
				}

				err = listener.bus.Bots.ShoutOut.SubscribeGroup(
					"bots",
					func(ctx context.Context, data bots.SentShoutOutRequest) (struct{}, error) {
						err := opts.WorkersPool.SubmitErr(
							func() error {
								return listener.handleShoutOut(ctx, data)
							},
						).Wait()

						return struct{}{}, err
					},
				)
				if err != nil {
					return err
				}

				err = listener.bus.Bots.ModeratorAdd.SubscribeGroup(
					"bots",
					func(ctx context.Context, data bots.ModeratorAddRequest) (struct{}, error) {
						err := opts.WorkersPool.SubmitErr(
							func() error {
								return listener.handleModeratorAdd(ctx, data)
							},
						).Wait()

						return struct{}{}, err
					},
				)
				if err != nil {
					return err
				}

				err = listener.bus.Bots.ModeratorRemove.SubscribeGroup(
					"bots",
					func(ctx context.Context, data bots.ModeratorRemoveRequest) (struct{}, error) {
						err := opts.WorkersPool.SubmitErr(
							func() error {
								return listener.handleModeratorRemove(ctx, data)
							},
						).Wait()

						return struct{}{}, err
					},
				)
				if err != nil {
					return err
				}

				err = listener.bus.Events.ChannelUnban.SubscribeGroup(
					"bots",
					func(ctx context.Context, data events.ChannelUnbanMessage) (struct{}, error) {
						err := opts.WorkersPool.SubmitErr(
							func() error {
								return listener.handleUnban(ctx, data)
							},
						).Wait()

						return struct{}{}, err
					},
				)
				if err != nil {
					return err
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				listener.bus.Bots.SendMessage.Unsubscribe()
				listener.bus.ChatMessages.Unsubscribe()
				listener.bus.Bots.DeleteMessage.Unsubscribe()
				listener.bus.Bots.BanUser.Unsubscribe()
				listener.bus.Bots.ShoutOut.Unsubscribe()
				listener.bus.Bots.ModeratorAdd.Unsubscribe()
				listener.bus.Bots.ModeratorRemove.Unsubscribe()
				listener.bus.Events.ChannelUnban.Unsubscribe()

				return nil
			},
		},
	)

	return listener, nil
}

type BusListener struct {
	logger         logger.Logger
	tracer         trace.Tracer
	gorm           *gorm.DB
	twitchActions  *twitchactions.TwitchActions
	messageHandler *messagehandler.MessageHandler
	bus            *bus_core.Bus
	config         cfg.Config
	kv             kv.KV
}

func (c *BusListener) deleteMessage(ctx context.Context, req bots.DeleteMessageRequest) error {
	channel := model.Channels{}
	err := c.gorm.WithContext(ctx).Where("id = ?", req.ChannelId).Find(&channel).Error
	if err != nil {
		c.logger.Error(
			"cannot get channel",
			slog.String("channelId", req.ChannelId),
		)
		return err
	}

	if channel.ID == "" {
		return nil
	}

	wg, wgCtx := errgroup.WithContext(ctx)

	for _, m := range req.MessageIds {
		wg.Go(
			func() error {
				e := c.twitchActions.DeleteMessage(
					wgCtx,
					twitchactions.DeleteMessageOpts{
						BroadcasterID: req.ChannelId,
						ModeratorID:   channel.BotID,
						MessageID:     m,
					},
				)
				if e != nil {
					c.logger.Error("cannot delete message", slog.Any("err", e))
					return e
				}

				return nil
			},
		)
	}

	if err := wg.Wait(); err != nil {
		return err
	}

	return nil
}

func (c *BusListener) sendMessage(ctx context.Context, req bots.SendMessageRequest) error {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(
		attribute.String("channel_id", req.ChannelId),
		attribute.String("message", req.Message),
		attribute.String("reply_to", req.ReplyTo),
	)

	if req.ChannelId == "" {
		return fmt.Errorf("channel id is empty")
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
			AnnounceColor:        req.AnnounceColor,
		},
	)
	if err != nil {
		c.logger.Error("cannot send message", slog.Any("err", err))
		return err
	}
	return nil
}

func (c *BusListener) handleChatMessage(
	ctx context.Context,
	req twitch.TwitchChatMessage,
) (struct{}, error) {
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
		return struct{}{}, err
	}

	return struct{}{}, nil
}

func (c *BusListener) handleShoutOut(ctx context.Context, req bots.SentShoutOutRequest) error {
	err := c.twitchActions.ShoutOut(
		ctx,
		twitchactions.ShoutOutInput{
			BroadcasterID: req.ChannelID,
			TargetID:      req.TargetID,
		},
	)
	if err != nil {
		c.logger.Error("cannot send shoutout", slog.Any("err", err))
		return err
	}
	return nil
}

func (c *BusListener) handleModeratorAdd(
	ctx context.Context,
	req bots.ModeratorAddRequest,
) error {
	return c.twitchActions.AddModerator(ctx, req.ChannelID, req.TargetID)
}

func (c *BusListener) handleModeratorRemove(
	ctx context.Context,
	req bots.ModeratorRemoveRequest,
) error {
	return c.twitchActions.RemoveModerator(ctx, req.ChannelID, req.TargetID)
}

func (c *BusListener) handleUnban(
	ctx context.Context,
	data events.ChannelUnbanMessage,
) error {
	modTaskExists, err := c.kv.Exists(
		ctx,
		redis_keys.CreateDistributedModTaskKey(data.ModeratorUserID, data.UserID),
	)
	if err != nil {
		return fmt.Errorf("cannot check distributed mod task existence: %w", err)
	}

	if modTaskExists {
		err := c.twitchActions.AddModerator(ctx, data.ModeratorUserID, data.UserID)
		if err != nil {
			return fmt.Errorf("cannot add moderator after unban: %w", err)
		}

		c.logger.Info(
			"added moderator after unban",
			slog.String("channel_id", data.ModeratorUserID),
			slog.String("user_id", data.UserID),
		)
	}

	return nil
}
