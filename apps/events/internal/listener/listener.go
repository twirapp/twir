package listener

import (
	"context"
	"log/slog"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/events/internal/chat_alerts"
	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/apps/events/internal/song_request"
	"github.com/twirapp/twir/apps/events/internal/workflows"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	cfg "github.com/twirapp/twir/libs/config"
	deprecatedgormmodel "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/logger"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	channelseventslistmodel "github.com/twirapp/twir/libs/repositories/channels_events_list/model"
	"github.com/twirapp/twir/libs/repositories/events/model"
	twitchlib "github.com/twirapp/twir/libs/twitch"
	"github.com/twirapp/twir/libs/utils"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Logger logger.Logger
	Cfg    cfg.Config
	Db     *gorm.DB
	Redis  *redis.Client

	WebsocketsGrpc websockets.WebsocketClient

	ChatAlerts             *chat_alerts.ChatAlerts
	EventsWorkflow         *workflows.EventWorkflow
	SongRequest            *song_request.SongRequest
	TwirBus                *buscore.Bus
	ChannelsEventsListRepo channelseventslist.Repository
}

func New(opts Opts) error {
	impl := &EventsGrpcImplementation{
		db:                     opts.Db,
		redis:                  opts.Redis,
		logger:                 opts.Logger,
		cfg:                    opts.Cfg,
		websocketsGrpc:         opts.WebsocketsGrpc,
		chatAlerts:             opts.ChatAlerts,
		eventsWorkflow:         opts.EventsWorkflow,
		songsRequest:           opts.SongRequest,
		twirBus:                opts.TwirBus,
		channelsEventsListRepo: opts.ChannelsEventsListRepo,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := impl.twirBus.Events.Follow.SubscribeGroup("events", impl.Follow); err != nil {
					return err
				}

				if err := impl.twirBus.Events.Subscribe.SubscribeGroup(
					"events",
					impl.Subscribe,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.ReSubscribe.SubscribeGroup(
					"events",
					impl.ReSubscribe,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.RedemptionCreated.SubscribeGroup(
					"events",
					impl.RedemptionCreated,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.SubGift.SubscribeGroup("events", impl.SubGift); err != nil {
					return err
				}

				if err := impl.twirBus.Events.CommandUsed.SubscribeGroup(
					"events",
					impl.CommandUsed,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.FirstUserMessage.SubscribeGroup(
					"events",
					impl.FirstUserMessage,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.Raided.SubscribeGroup("events", impl.Raided); err != nil {
					return err
				}

				if err := impl.twirBus.Events.TitleOrCategoryChanged.SubscribeGroup(
					"events",
					impl.TitleOrCategoryChanged,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Channel.StreamOnline.SubscribeGroup(
					"events",
					impl.StreamOnline,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Channel.StreamOffline.SubscribeGroup(
					"events",
					impl.StreamOffline,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.ChatClear.SubscribeGroup(
					"events",
					impl.ChatClear,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.Donate.SubscribeGroup("events", impl.Donate); err != nil {
					return err
				}

				if err := impl.twirBus.Events.KeywordMatched.SubscribeGroup(
					"events",
					impl.KeywordMatched,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.GreetingSended.SubscribeGroup(
					"events",
					impl.GreetingSended,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.PollBegin.SubscribeGroup(
					"events",
					impl.PollBegin,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.PollProgress.SubscribeGroup(
					"events",
					impl.PollProgress,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.PollEnd.SubscribeGroup("events", impl.PollEnd); err != nil {
					return err
				}

				if err := impl.twirBus.Events.PredictionBegin.SubscribeGroup(
					"events",
					impl.PredictionBegin,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.PredictionProgress.SubscribeGroup(
					"events",
					impl.PredictionProgress,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.PredictionLock.SubscribeGroup(
					"events",
					impl.PredictionLock,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.PredictionEnd.SubscribeGroup(
					"events",
					impl.PredictionEnd,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.StreamFirstUserJoin.SubscribeGroup(
					"events",
					impl.StreamFirstUserJoin,
				); err != nil {
					return err
				}
				if err := impl.twirBus.Events.ChannelBan.SubscribeGroup(
					"events",
					impl.ChannelBan,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.ChannelUnbanRequestCreate.SubscribeGroup(
					"events",
					impl.ChannelUnbanRequestCreate,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.ChannelUnbanRequestResolve.SubscribeGroup(
					"events",
					impl.ChannelUnbanRequestResolve,
				); err != nil {
					return err
				}

				if err := impl.twirBus.Events.ChannelMessageDelete.SubscribeGroup(
					"events",
					impl.ChannelMessageDelete,
				); err != nil {
					return err
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				impl.twirBus.Events.Follow.Unsubscribe()
				impl.twirBus.Events.Subscribe.Unsubscribe()
				impl.twirBus.Events.ReSubscribe.Unsubscribe()
				impl.twirBus.Events.RedemptionCreated.Unsubscribe()
				impl.twirBus.Events.SubGift.Unsubscribe()
				impl.twirBus.Events.FirstUserMessage.Unsubscribe()
				impl.twirBus.Events.Raided.Unsubscribe()
				impl.twirBus.Channel.StreamOnline.Unsubscribe()
				impl.twirBus.Channel.StreamOffline.Unsubscribe()
				impl.twirBus.Events.ChatClear.Unsubscribe()
				impl.twirBus.Events.Donate.Unsubscribe()
				impl.twirBus.Events.GreetingSended.Unsubscribe()
				impl.twirBus.Events.PollBegin.Unsubscribe()
				impl.twirBus.Events.PollProgress.Unsubscribe()
				impl.twirBus.Events.PollEnd.Unsubscribe()
				impl.twirBus.Events.PredictionProgress.Unsubscribe()
				impl.twirBus.Events.PredictionLock.Unsubscribe()
				impl.twirBus.Events.PredictionEnd.Unsubscribe()
				impl.twirBus.Events.StreamFirstUserJoin.Unsubscribe()
				impl.twirBus.Events.ChannelBan.Unsubscribe()
				impl.twirBus.Events.ChannelUnbanRequestCreate.Unsubscribe()
				impl.twirBus.Events.ChannelUnbanRequestResolve.Unsubscribe()
				impl.twirBus.Events.ChannelMessageDelete.Unsubscribe()

				return nil
			},
		},
	)

	return nil
}

type EventsGrpcImplementation struct {
	db     *gorm.DB
	redis  *redis.Client
	logger logger.Logger
	cfg    cfg.Config

	websocketsGrpc websockets.WebsocketClient

	chatAlerts             *chat_alerts.ChatAlerts
	eventsWorkflow         *workflows.EventWorkflow
	songsRequest           *song_request.SongRequest
	twirBus                *buscore.Bus
	channelsEventsListRepo channelseventslist.Repository
}

func (c *EventsGrpcImplementation) Follow(
	ctx context.Context,
	msg events.FollowMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			var stream *deprecatedgormmodel.ChannelsStreams
			if err := c.db.Where(`"userId" = ?`, msg.BaseInfo.ChannelID).
				Find(&stream).Error; err != nil {
				c.logger.Error("Error get stream", slog.Any("err", err))
				return
			}

			var streamFollowersCount int64
			if stream != nil && stream.ID != "" {
				t := channelseventslistmodel.ChannelEventListItemTypeFollow
				count, err := c.channelsEventsListRepo.CountBy(
					ctx,
					channelseventslist.CountByInput{
						ChannelID:    &msg.BaseInfo.ChannelID,
						CreatedAtGTE: &stream.StartedAt,
						Type:         &t,
					},
				)
				if err != nil {
					c.logger.Error("Error get stream followers count", slog.Any("err", err))
					return
				}

				streamFollowersCount = count
			}

			twitchClient, err := twitchlib.NewUserClientWithContext(
				ctx,
				msg.BaseInfo.ChannelID,
				c.cfg,
				c.twirBus,
			)
			if err != nil {
				c.logger.Error("Error create twitch client", slog.Any("err", err))
				return
			}

			followersReq, err := twitchClient.GetChannelFollows(
				&helix.GetChannelFollowsParams{
					BroadcasterID: msg.BaseInfo.ChannelID,
				},
			)
			if err != nil {
				c.logger.Error("Error get channel followers", slog.Any("err", err))
				return
			}

			err = c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeFollow,
				shared.EventData{
					ChannelID:              msg.BaseInfo.ChannelID,
					UserName:               msg.UserName,
					UserDisplayName:        msg.UserDisplayName,
					UserID:                 msg.UserID,
					ChannelFollowers:       int64(followersReq.Data.Total),
					ChannelStreamFollowers: streamFollowersCount,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Go(
		func() {
			c.chatAlerts.ProcessEvent(
				ctx,
				msg.BaseInfo.ChannelID,
				model.EventTypeFollow,
				msg,
			)
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) Subscribe(
	ctx context.Context,
	msg events.SubscribeMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeSubscribe,
				shared.EventData{
					ChannelID:       msg.BaseInfo.ChannelID,
					UserDisplayName: msg.UserDisplayName,
					SubLevel:        msg.Level,
					UserID:          msg.UserID,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Go(
		func() {
			c.chatAlerts.ProcessEvent(
				ctx,
				msg.BaseInfo.ChannelID,
				model.EventTypeSubscribe,
				chat_alerts.SubscribeMessage{
					UserName:  msg.UserName,
					Months:    0,
					ChannelId: msg.BaseInfo.ChannelID,
				},
			)
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) ReSubscribe(
	ctx context.Context,
	msg events.ReSubscribeMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeResubscribe,
				shared.EventData{
					ChannelID:       msg.BaseInfo.ChannelID,
					UserDisplayName: msg.UserDisplayName,
					SubLevel:        msg.Level,
					ResubMessage:    msg.Message,
					ResubMonths:     msg.Months,
					ResubStreak:     msg.Streak,
					UserID:          msg.UserID,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Go(
		func() {
			c.chatAlerts.ProcessEvent(
				ctx,
				msg.BaseInfo.ChannelID,
				model.EventTypeResubscribe,
				msg,
			)
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) RedemptionCreated(
	ctx context.Context,
	msg events.RedemptionCreatedMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeRedemptionCreated,
				shared.EventData{
					ChannelID:       msg.BaseInfo.ChannelID,
					UserName:        msg.UserName,
					UserDisplayName: msg.UserDisplayName,
					RewardCost:      msg.RewardCost,
					RewardInput:     msg.Input,
					RewardName:      msg.RewardName,
					RewardID:        msg.ID,
					UserID:          msg.UserID,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Go(
		func() {
			c.chatAlerts.ProcessEvent(
				ctx,
				msg.BaseInfo.ChannelID,
				model.EventTypeRedemptionCreated,
				msg,
			)
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) CommandUsed(
	ctx context.Context,
	msg events.CommandUsedMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeCommandUsed,
				shared.EventData{
					ChannelID:       msg.BaseInfo.ChannelID,
					UserName:        msg.UserName,
					UserDisplayName: msg.UserDisplayName,
					CommandName:     msg.CommandName,
					CommandID:       msg.CommandID,
					CommandInput:    msg.CommandInput,
					UserID:          msg.UserID,
					ChatMessageId:   msg.MessageID,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) FirstUserMessage(
	ctx context.Context,
	msg events.FirstUserMessageMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeFirstUserMessage,
				shared.EventData{
					ChannelID:       msg.BaseInfo.ChannelID,
					UserName:        msg.UserName,
					UserDisplayName: msg.UserDisplayName,
					UserID:          msg.UserID,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Go(
		func() {
			c.chatAlerts.ProcessEvent(
				ctx,
				msg.BaseInfo.ChannelID,
				model.EventTypeFirstUserMessage,
				msg,
			)
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) Raided(
	ctx context.Context,
	msg events.RaidedMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeRaided,
				shared.EventData{
					ChannelID:       msg.BaseInfo.ChannelID,
					UserName:        msg.UserName,
					UserDisplayName: msg.UserDisplayName,
					RaidViewers:     msg.Viewers,
					UserID:          msg.UserID,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Go(
		func() {
			c.chatAlerts.ProcessEvent(
				ctx,
				msg.BaseInfo.ChannelID,
				model.EventTypeRaided,
				msg,
			)
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) TitleOrCategoryChanged(
	ctx context.Context,
	msg events.TitleOrCategoryChangedMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeTitleOrCategoryChanged,
				shared.EventData{
					ChannelID:         msg.BaseInfo.ChannelID,
					OldStreamCategory: msg.OldCategory,
					NewStreamCategory: msg.NewCategory,
					OldStreamTitle:    msg.OldTitle,
					NewStreamTitle:    msg.NewTitle,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) StreamOnline(
	ctx context.Context,
	msg twitch.StreamOnlineMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeStreamOnline,
				shared.EventData{
					ChannelID:      msg.ChannelID,
					StreamTitle:    msg.Title,
					StreamCategory: msg.CategoryName,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Go(
		func() {
			c.chatAlerts.ProcessEvent(
				ctx,
				msg.ChannelID,
				model.EventTypeStreamOnline,
				msg,
			)
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) StreamOffline(
	ctx context.Context,
	msg twitch.StreamOfflineMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeStreamOffline,
				shared.EventData{
					ChannelID: msg.ChannelID,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Go(
		func() {
			c.chatAlerts.ProcessEvent(
				ctx,
				msg.ChannelID,
				model.EventTypeStreamOffline,
				msg,
			)
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) SubGift(
	ctx context.Context,
	msg events.SubGiftMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeSubGift,
				shared.EventData{
					ChannelID:             msg.BaseInfo.ChannelID,
					TargetUserName:        msg.TargetUserName,
					TargetUserDisplayName: msg.TargetDisplayName,
					SubLevel:              msg.Level,
					UserID:                msg.SenderUserID,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) ChatClear(
	ctx context.Context,
	msg events.ChatClearMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeOnChatClear,
				shared.EventData{
					ChannelID: msg.BaseInfo.ChannelID,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Go(
		func() {
			c.chatAlerts.ProcessEvent(
				ctx,
				msg.BaseInfo.ChannelID,
				model.EventTypeOnChatClear,
				msg,
			)
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) Donate(
	ctx context.Context,
	msg events.DonateMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeDonate,
				shared.EventData{
					ChannelID:      msg.BaseInfo.ChannelID,
					UserName:       msg.UserName,
					DonateAmount:   msg.Amount,
					DonateCurrency: msg.Currency,
					DonateMessage:  msg.Message,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Go(
		func() {
			c.chatAlerts.ProcessEvent(
				ctx,
				msg.BaseInfo.ChannelID,
				model.EventTypeDonate,
				msg,
			)
		},
	)

	wg.Go(
		func() {
			err := c.songsRequest.ProcessFromDonation(
				ctx, song_request.ProcessFromDonationInput{
					Text:      msg.Message,
					ChannelID: msg.BaseInfo.ChannelID,
				},
			)

			if err != nil {
				c.logger.Error("Error processing donation", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) KeywordMatched(
	ctx context.Context,
	msg events.KeywordMatchedMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeKeywordMatched,
				shared.EventData{
					ChannelID:       msg.BaseInfo.ChannelID,
					UserName:        msg.UserName,
					UserDisplayName: msg.UserDisplayName,
					KeywordName:     msg.KeywordName,
					KeywordResponse: msg.KeywordResponse,
					KeywordID:       msg.KeywordID,
					UserID:          msg.UserID,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) GreetingSended(
	ctx context.Context,
	msg events.GreetingSendedMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeGreetingSended,
				shared.EventData{
					ChannelID:       msg.BaseInfo.ChannelID,
					UserName:        msg.UserName,
					UserDisplayName: msg.UserDisplayName,
					UserID:          msg.UserID,
					GreetingText:    msg.GreetingText,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) PollBegin(
	ctx context.Context,
	msg events.PollBeginMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypePollBegin,
				shared.EventData{
					ChannelID:       msg.BaseInfo.ChannelID,
					UserName:        msg.UserName,
					UserDisplayName: msg.UserDisplayName,
					PollTitle:       msg.Info.Title,
					PollOptionsNames: strings.Join(
						lo.Map(
							msg.Info.Choices, func(item events.PollChoice, _ int) string {
								return item.Title
							},
						), " · ",
					),
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) PollProgress(
	ctx context.Context,
	msg events.PollProgressMessage,
) (struct{}, error) {
	totalVotes := lo.Reduce(
		msg.Info.Choices, func(acc int, item events.PollChoice, _ int) int {
			return acc + int(item.Votes)
		}, 0,
	)

	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypePollProgress,
				shared.EventData{
					ChannelID:       msg.BaseInfo.ChannelID,
					UserDisplayName: msg.UserDisplayName,
					PollTitle:       msg.Info.Title,
					PollOptionsNames: strings.Join(
						lo.Map(
							msg.Info.Choices, func(item events.PollChoice, _ int) string {
								return item.Title
							},
						), " · ",
					),
					PollTotalVotes: totalVotes,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) PollEnd(
	ctx context.Context,
	msg events.PollEndMessage,
) (struct{}, error) {
	totalVotes := lo.Reduce(
		msg.Info.Choices, func(acc int, item events.PollChoice, _ int) int {
			return acc + int(item.Votes)
		}, 0,
	)

	// find most total votes in choices
	winner := lo.MaxBy(
		msg.Info.Choices, func(a events.PollChoice, b events.PollChoice) bool {
			return a.Votes > b.Votes
		},
	)

	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypePollEnd,
				shared.EventData{
					ChannelID:       msg.BaseInfo.ChannelID,
					UserDisplayName: msg.UserDisplayName,
					PollTitle:       msg.Info.Title,
					PollOptionsNames: strings.Join(
						lo.Map(
							msg.Info.Choices, func(item events.PollChoice, _ int) string {
								return item.Title
							},
						), " · ",
					),
					PollWinnerTitle:               winner.Title,
					PollWinnerBitsVotes:           int(winner.BitsVotes),
					PollWinnerChannelsPointsVotes: int(winner.ChannelsPointsVotes),
					PollWinnerTotalVotes:          int(winner.Votes),
					PollTotalVotes:                totalVotes,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) PredictionBegin(
	ctx context.Context, msg events.PredictionBeginMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypePredictionBegin,
				shared.EventData{
					ChannelID:       msg.BaseInfo.ChannelID,
					UserName:        msg.UserName,
					UserDisplayName: msg.UserDisplayName,
					PredictionTitle: msg.Info.Title,
					PredictionOptionsNames: strings.Join(
						lo.Map(
							msg.Info.Outcomes, func(item events.PredictionOutcome, _ int) string {
								return item.Title
							},
						), " · ",
					),
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) PredictionProgress(
	ctx context.Context, msg events.PredictionProgressMessage,
) (struct{}, error) {
	totalPoints := lo.Reduce(
		msg.Info.Outcomes, func(acc int, item events.PredictionOutcome, _ int) int {
			return acc + int(item.ChannelPoints)
		}, 0,
	)

	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypePredictionProgress,
				shared.EventData{
					ChannelID:       msg.BaseInfo.ChannelID,
					UserDisplayName: msg.UserDisplayName,
					PredictionTitle: msg.Info.Title,
					PredictionOptionsNames: strings.Join(
						lo.Map(
							msg.Info.Outcomes, func(item events.PredictionOutcome, _ int) string {
								return item.Title
							},
						), " · ",
					),
					PredictionTotalChannelPoints: totalPoints,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) PredictionLock(
	ctx context.Context,
	msg events.PredictionLockMessage,
) (struct{}, error) {
	totalPoints := lo.Reduce(
		msg.Info.Outcomes, func(acc int, item events.PredictionOutcome, _ int) int {
			return acc + int(item.ChannelPoints)
		}, 0,
	)

	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypePredictionLock,
				shared.EventData{
					ChannelID:       msg.BaseInfo.ChannelID,
					UserDisplayName: msg.UserDisplayName,
					PredictionTitle: msg.Info.Title,
					PredictionOptionsNames: strings.Join(
						lo.Map(
							msg.Info.Outcomes, func(item events.PredictionOutcome, _ int) string {
								return item.Title
							},
						), " · ",
					),
					PredictionTotalChannelPoints: totalPoints,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) PredictionEnd(
	ctx context.Context,
	msg events.PredictionEndMessage,
) (struct{}, error) {
	totalPoints := lo.Reduce(
		msg.Info.Outcomes, func(acc int, item events.PredictionOutcome, _ int) int {
			return acc + int(item.ChannelPoints)
		}, 0,
	)

	winner, _ := lo.Find(
		msg.Info.Outcomes, func(item events.PredictionOutcome) bool {
			return item.ID == msg.WinningOutcomeID
		},
	)

	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypePredictionEnd,
				shared.EventData{
					ChannelID:       msg.BaseInfo.ChannelID,
					UserDisplayName: msg.UserDisplayName,
					PredictionTitle: msg.Info.Title,
					PredictionOptionsNames: strings.Join(
						lo.Map(
							msg.Info.Outcomes, func(item events.PredictionOutcome, _ int) string {
								return item.Title
							},
						), " · ",
					),
					PredictionWinner: &shared.EventDataPredictionOutCome{
						Title:       winner.Title,
						TotalUsers:  int(winner.Users),
						TotalPoints: int(winner.ChannelPoints),
						TopUsers:    predictionMapTopPredictors(winner.TopPredictors),
					},
					PredictionTotalChannelPoints: totalPoints,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) StreamFirstUserJoin(
	ctx context.Context, msg events.StreamFirstUserJoinMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeFirstUserMessage,
				shared.EventData{
					ChannelID: msg.BaseInfo.ChannelID,
					UserName:  msg.UserLogin,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) ChannelBan(
	ctx context.Context,
	msg events.ChannelBanMessage,
) (struct{}, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeChannelBan,
				shared.EventData{
					ChannelID:            msg.BaseInfo.ChannelID,
					UserDisplayName:      msg.UserName,
					ModeratorDisplayName: msg.ModeratorUserName,
					ModeratorName:        msg.ModeratorUserLogin,
					BanReason:            msg.Reason,
					BanEndsInMinutes:     msg.EndsAt,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Go(
		func() {
			c.chatAlerts.ProcessEvent(
				ctx,
				msg.BaseInfo.ChannelID,
				model.EventTypeChannelBan,
				msg,
			)
		},
	)

	wg.Wait()

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) ChannelUnbanRequestCreate(
	ctx context.Context,
	msg events.ChannelUnbanRequestCreateMessage,
) (struct{}, error) {
	c.chatAlerts.ProcessEvent(
		ctx,
		msg.BaseInfo.ChannelID,
		model.EventTypeRedemptionCreated,
		msg,
	)

	err := c.eventsWorkflow.Execute(
		ctx,
		model.EventTypeRedemptionCreated,
		shared.EventData{
			ChannelID:       msg.BaseInfo.ChannelID,
			UserName:        msg.UserLogin,
			UserDisplayName: msg.UserName,
			Message:         msg.Text,
		},
	)
	if err != nil {
		c.logger.Error("Error execute workflow", slog.Any("err", err))
	}

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) ChannelUnbanRequestResolve(
	ctx context.Context,
	msg events.ChannelUnbanRequestResolveMessage,
) (struct{}, error) {
	err := c.eventsWorkflow.Execute(
		ctx,
		model.EventTypeChannelUnbanRequestResolve,
		shared.EventData{
			ChannelID:                          msg.BaseInfo.ChannelID,
			UserName:                           msg.UserLogin,
			UserDisplayName:                    msg.UserName,
			Message:                            msg.Reason,
			ChannelUnbanRequestResolveDeclined: msg.Declined,
			ModeratorName:                      msg.ModeratorUserLogin,
			ModeratorDisplayName:               msg.ModeratorUserName,
		},
	)
	if err != nil {
		c.logger.Error("Error execute workflow", slog.Any("err", err))
	}

	return struct{}{}, nil
}

func (c *EventsGrpcImplementation) ChannelMessageDelete(
	ctx context.Context,
	msg events.ChannelMessageDeleteMessage,
) (struct{}, error) {
	c.chatAlerts.ProcessEvent(
		ctx,
		msg.BaseInfo.ChannelID,
		model.EventTypeChannelMessageDelete,
		msg,
	)

	err := c.eventsWorkflow.Execute(
		ctx,
		model.EventTypeChannelMessageDelete,
		shared.EventData{
			ChannelID:       msg.BaseInfo.ChannelID,
			UserName:        msg.UserLogin,
			UserDisplayName: msg.UserName,
			ChatMessageId:   msg.MessageId,
		},
	)
	if err != nil {
		c.logger.Error("Error execute workflow", slog.Any("err", err))
	}

	return struct{}{}, nil
}
