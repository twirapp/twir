package bus_listener

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/twirapp/kv"
	discordmessagesupdater "github.com/twirapp/twir/apps/bots/internal/discord/messages_updater"
	"github.com/twirapp/twir/apps/bots/internal/messagehandler"
	"github.com/twirapp/twir/apps/bots/internal/services/channel"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/apps/bots/internal/workers"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/redis_keys"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Logger         *slog.Logger
	Tracer         trace.Tracer
	Gorm           *gorm.DB
	TwitchActions  *twitchactions.TwitchActions
	MessageHandler *messagehandler.MessageHandler
	Bus            *buscore.Bus
	Cfg            cfg.Config
	WorkersPool    *workers.Pool
	KV             kv.KV

	ChannelService         *channel.Service
	DiscordMessagesUpdater *discordmessagesupdater.MessagesUpdater
}

func New(opts Opts) (*BusListener, error) {
	listener := &BusListener{
		gorm:                   opts.Gorm,
		logger:                 opts.Logger,
		config:                 opts.Cfg,
		twitchActions:          opts.TwitchActions,
		messageHandler:         opts.MessageHandler,
		tracer:                 opts.Tracer,
		bus:                    opts.Bus,
		kv:                     opts.KV,
		channelService:         opts.ChannelService,
		discordMessagesUpdater: opts.DiscordMessagesUpdater,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				err := listener.bus.Bots.SendMessage.SubscribeGroup(
					"bots",
					func(ctx context.Context, req bots.SendMessageRequest) (struct{}, error) {
						if err := listener.channelService.SendMessage(ctx, req); err != nil {
							return struct{}{}, fmt.Errorf("send message: %w", err)
						}

						return struct{}{}, nil
					},
				)
				if err != nil {
					return err
				}

				err = listener.bus.Bots.DeleteMessage.SubscribeGroup(
					"bots",
					func(ctx context.Context, req bots.DeleteMessageRequest) (struct{}, error) {
						if err := listener.channelService.DeleteMessage(ctx, req); err != nil {
							return struct{}{}, fmt.Errorf("delete message: %w", err)
						}

						return struct{}{}, nil
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
					func(ctx context.Context, req bots.BanRequest) (struct{}, error) {
						if err := listener.channelService.Ban(ctx, req); err != nil {
							return struct{}{}, fmt.Errorf("ban: %w", err)
						}

						return struct{}{}, nil
					},
				)
				if err != nil {
					return err
				}

				err = listener.bus.Bots.BanUsers.SubscribeGroup(
					"bots",
					func(ctx context.Context, reqs []bots.BanRequest) (struct{}, error) {
						if err := listener.channelService.BanMany(ctx, reqs); err != nil {
							return struct{}{}, fmt.Errorf("ban many: %w", err)
						}

						return struct{}{}, nil
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

				err = listener.bus.Channel.StreamOffline.SubscribeGroup(
					"bots",
					func(ctx context.Context, data twitch.StreamOfflineMessage) (struct{}, error) {
						err := opts.WorkersPool.SubmitErr(
							func() error {
								return listener.handleStreamOffline(ctx, data)
							},
						).Wait()

						return struct{}{}, err
					},
				)

				err = listener.bus.Channel.StreamOnline.SubscribeGroup(
					"bots",
					func(ctx context.Context, data twitch.StreamOnlineMessage) (struct{}, error) {
						err := opts.WorkersPool.SubmitErr(
							func() error {
								return listener.handleStreamOnline(ctx, data)
							},
						).Wait()

						return struct{}{}, err
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
				listener.bus.Bots.ModeratorAdd.Unsubscribe()
				listener.bus.Bots.ModeratorRemove.Unsubscribe()
				listener.bus.Events.ChannelUnban.Unsubscribe()
				listener.bus.Channel.StreamOnline.Unsubscribe()
				listener.bus.Channel.StreamOffline.Unsubscribe()

				return nil
			},
		},
	)

	return listener, nil
}

type BusListener struct {
	logger                 *slog.Logger
	tracer                 trace.Tracer
	gorm                   *gorm.DB
	twitchActions          *twitchactions.TwitchActions
	messageHandler         *messagehandler.MessageHandler
	bus                    *buscore.Bus
	config                 cfg.Config
	kv                     kv.KV
	channelService         *channel.Service
	discordMessagesUpdater *discordmessagesupdater.MessagesUpdater
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
			logger.Error(err),
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
		c.logger.Error("cannot send shoutout", logger.Error(err))
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
		defer c.kv.Delete(
			ctx,
			redis_keys.CreateDistributedModTaskKey(data.ModeratorUserID, data.UserID),
		)

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

func (c *BusListener) handleStreamOnline(ctx context.Context, data twitch.StreamOnlineMessage) error {
	if err := c.discordMessagesUpdater.ProcessOnline(ctx, data.ChannelID); err != nil {
		c.logger.Error(
			"cannot handle discord stream online",
			logger.Error(err),
			slog.String("channel_id", data.ChannelID),
		)
	}

	return nil
}

func (c *BusListener) handleStreamOffline(ctx context.Context, data twitch.StreamOfflineMessage) error {
	if err := c.discordMessagesUpdater.ProcessOffline(ctx, data.ChannelID); err != nil {
		c.logger.Error(
			"cannot handle discord stream offline",
			logger.Error(err),
			slog.String("channel_id", data.ChannelID),
		)
	}

	return nil
}
