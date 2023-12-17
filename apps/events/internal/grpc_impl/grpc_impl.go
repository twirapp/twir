package grpc_impl

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/events/internal/chat_alerts"
	"github.com/satont/twir/apps/events/internal/shared"
	"github.com/satont/twir/apps/events/internal/workflows"
	cfg "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/constants"
	api_events "github.com/satont/twir/libs/grpc/generated/api/events"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/events"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/utils"
	"go.uber.org/fx"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Logger logger.Logger
	Cfg    cfg.Config
	Db     *gorm.DB
	Redis  *redis.Client

	BotsGrpc       bots.BotsClient
	TokensGrpc     tokens.TokensClient
	WebsocketsGrpc websockets.WebsocketClient

	ChatAlerts     *chat_alerts.ChatAlerts
	EventsWorkflow *workflows.EventWorkflow
}

func New(opts Opts) error {
	impl := &EventsGrpcImplementation{
		db:             opts.Db,
		redis:          opts.Redis,
		logger:         opts.Logger,
		cfg:            opts.Cfg,
		botsGrpc:       opts.BotsGrpc,
		tokensGrpc:     opts.TokensGrpc,
		websocketsGrpc: opts.WebsocketsGrpc,
		chatAlerts:     opts.ChatAlerts,
		eventsWorkflow: opts.EventsWorkflow,
	}

	lis, err := net.Listen("tcp", fmt.Sprintf("0.0.0.0:%d", constants.EVENTS_SERVER_PORT))
	if err != nil {
		return err
	}
	grpcServer := grpc.NewServer()
	events.RegisterEventsServer(grpcServer, impl)

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				go grpcServer.Serve(lis)
				opts.Logger.Info("Grpc server started", slog.Int("port", constants.EVENTS_SERVER_PORT))
				return nil
			},
			OnStop: func(ctx context.Context) error {
				grpcServer.GracefulStop()
				return nil
			},
		},
	)

	return nil
}

type EventsGrpcImplementation struct {
	events.UnimplementedEventsServer

	db     *gorm.DB
	redis  *redis.Client
	logger logger.Logger
	cfg    cfg.Config

	botsGrpc       bots.BotsClient
	tokensGrpc     tokens.TokensClient
	websocketsGrpc websockets.WebsocketClient

	chatAlerts     *chat_alerts.ChatAlerts
	eventsWorkflow *workflows.EventWorkflow
}

// func (c *EventsGrpcImplementation) Follow(
// 	ctx context.Context,
// 	msg *events.FollowMessage,
// ) (*emptypb.Empty, error) {
// 	go c.processEvent(
// 		msg.BaseInfo.ChannelId,
// 		internal.Data{
// 			UserName:        msg.UserName,
// 			UserDisplayName: msg.UserDisplayName,
// 			UserID:          msg.UserId,
// 		},
// 		model.EventTypeFollow,
// 	)
// 	go c.chatAlerts.ProcessEvent(ctx, msg.BaseInfo.ChannelId, api_events.TwirEventType_FOLLOW, msg)
// 	go c.websocketsGrpc.TriggerKappagenByEvent(
// 		ctx, &websockets.TriggerKappagenByEventRequest{
// 			ChannelId: msg.BaseInfo.ChannelId,
// 			Event:     api_events.TwirEventType_FOLLOW,
// 		},
// 	)
//
// 	return &emptypb.Empty{}, nil
// }
//
// func (c *EventsGrpcImplementation) Subscribe(
// 	ctx context.Context,
// 	msg *events.SubscribeMessage,
// ) (*emptypb.Empty, error) {
// 	go c.processEvent(
// 		msg.BaseInfo.ChannelId,
// 		internal.Data{
// 			UserName:        msg.UserName,
// 			UserDisplayName: msg.UserDisplayName,
// 			SubLevel:        msg.Level,
// 			UserID:          msg.UserId,
// 		},
// 		model.EventTypeSubscribe,
// 	)
// 	go c.chatAlerts.ProcessEvent(
// 		ctx,
// 		msg.BaseInfo.ChannelId,
// 		api_events.TwirEventType_SUBSCRIBE,
// 		chat_alerts.SubscribMessage{
// 			UserName:  msg.UserName,
// 			Months:    0,
// 			ChannelId: msg.BaseInfo.ChannelId,
// 		},
// 	)
// 	go c.websocketsGrpc.TriggerKappagenByEvent(
// 		ctx, &websockets.TriggerKappagenByEventRequest{
// 			ChannelId: msg.BaseInfo.ChannelId,
// 			Event:     api_events.TwirEventType_SUBSCRIBE,
// 		},
// 	)
//
// 	return &emptypb.Empty{}, nil
// }
//
// func (c *EventsGrpcImplementation) ReSubscribe(
// 	ctx context.Context,
// 	msg *events.ReSubscribeMessage,
// ) (*emptypb.Empty, error) {
// 	go c.processEvent(
// 		msg.BaseInfo.ChannelId,
// 		internal.Data{
// 			UserName:        msg.UserName,
// 			UserDisplayName: msg.UserDisplayName,
// 			SubLevel:        msg.Level,
// 			ResubMessage:    msg.Message,
// 			ResubMonths:     msg.Months,
// 			ResubStreak:     msg.Streak,
// 			UserID:          msg.UserId,
// 		},
// 		model.EventTypeResubscribe,
// 	)
// 	go c.chatAlerts.ProcessEvent(ctx, msg.BaseInfo.ChannelId, api_events.TwirEventType_RESUBSCRIBE, msg)
// 	go c.websocketsGrpc.TriggerKappagenByEvent(
// 		ctx, &websockets.TriggerKappagenByEventRequest{
// 			ChannelId: msg.BaseInfo.ChannelId,
// 			Event:     api_events.TwirEventType_RESUBSCRIBE,
// 		},
// 	)
//
// 	return &emptypb.Empty{}, nil
// }
//
// func (c *EventsGrpcImplementation) RedemptionCreated(
// 	ctx context.Context,
// 	msg *events.RedemptionCreatedMessage,
// ) (*emptypb.Empty, error) {
// 	go c.processEvent(
// 		msg.BaseInfo.ChannelId,
// 		internal.Data{
// 			UserName:        msg.UserName,
// 			UserDisplayName: msg.UserDisplayName,
// 			RewardCost:      msg.RewardCost,
// 			RewardInput:     msg.Input,
// 			RewardName:      msg.RewardName,
// 			RewardID:        msg.Id,
// 			UserID:          msg.UserId,
// 		},
// 		model.EventTypeRedemptionCreated,
// 	)
// 	go c.chatAlerts.ProcessEvent(
// 		ctx,
// 		msg.BaseInfo.ChannelId,
// 		api_events.TwirEventType_REDEMPTION_CREATED,
// 		msg,
// 	)
// 	go c.websocketsGrpc.TriggerKappagenByEvent(
// 		ctx, &websockets.TriggerKappagenByEventRequest{
// 			ChannelId: msg.BaseInfo.ChannelId,
// 			Event:     api_events.TwirEventType_REDEMPTION_CREATED,
// 		},
// 	)
//
// 	return &emptypb.Empty{}, nil
// }

func (c *EventsGrpcImplementation) CommandUsed(
	ctx context.Context,
	msg *events.CommandUsedMessage,
) (
	*emptypb.Empty, error,
) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeCommandUsed,
				shared.EvenData{
					ChannelID:       msg.BaseInfo.ChannelId,
					UserName:        msg.UserName,
					UserDisplayName: msg.UserDisplayName,
					CommandName:     msg.CommandName,
					CommandID:       msg.CommandId,
					CommandInput:    msg.CommandInput,
					UserID:          msg.UserId,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	if msg.DefaultCommandName != "kappagen" {
		wg.Go(
			func() {
				_, err := c.websocketsGrpc.TriggerKappagenByEvent(
					ctx,
					&websockets.TriggerKappagenByEventRequest{
						ChannelId: msg.BaseInfo.ChannelId,
						Event:     api_events.TwirEventType_COMMAND_USED,
					},
				)
				if err != nil {
					c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
				}
			},
		)
	}

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) FirstUserMessage(
	ctx context.Context,
	msg *events.FirstUserMessageMessage,
) (*emptypb.Empty, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeFirstUserMessage,
				shared.EvenData{
					ChannelID:       msg.BaseInfo.ChannelId,
					UserName:        msg.UserName,
					UserDisplayName: msg.UserDisplayName,
					UserID:          msg.UserId,
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
				msg.BaseInfo.ChannelId,
				api_events.TwirEventType_FIRST_USER_MESSAGE,
				msg,
			)
		},
	)

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx,
				&websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_FIRST_USER_MESSAGE,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) Raided(
	ctx context.Context,
	msg *events.RaidedMessage,
) (*emptypb.Empty, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeRaided,
				shared.EvenData{
					ChannelID:       msg.BaseInfo.ChannelId,
					UserName:        msg.UserName,
					UserDisplayName: msg.UserDisplayName,
					RaidViewers:     msg.Viewers,
					UserID:          msg.UserId,
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
				msg.BaseInfo.ChannelId,
				api_events.TwirEventType_RAIDED,
				msg,
			)
		},
	)

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx,
				&websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_RAIDED,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) TitleOrCategoryChanged(
	ctx context.Context, msg *events.TitleOrCategoryChangedMessage,
) (*emptypb.Empty, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeTitleOrCategoryChanged,
				shared.EvenData{
					ChannelID:         msg.BaseInfo.ChannelId,
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

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_TITLE_OR_CATEGORY_CHANGED,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) StreamOnline(
	ctx context.Context,
	msg *events.StreamOnlineMessage,
) (
	*emptypb.Empty, error,
) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeStreamOnline,
				shared.EvenData{
					ChannelID:      msg.BaseInfo.ChannelId,
					StreamTitle:    msg.Title,
					StreamCategory: msg.Category,
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
				msg.BaseInfo.ChannelId,
				api_events.TwirEventType_STREAM_ONLINE,
				msg,
			)
		},
	)

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_STREAM_ONLINE,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) StreamOffline(
	ctx context.Context,
	msg *events.StreamOfflineMessage,
) (
	*emptypb.Empty, error,
) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeStreamOffline,
				shared.EvenData{
					ChannelID: msg.BaseInfo.ChannelId,
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
				msg.BaseInfo.ChannelId,
				api_events.TwirEventType_STREAM_OFFLINE,
				msg,
			)
		},
	)

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_STREAM_OFFLINE,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) SubGift(
	ctx context.Context,
	msg *events.SubGiftMessage,
) (*emptypb.Empty, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeSubGift,
				shared.EvenData{
					ChannelID:             msg.BaseInfo.ChannelId,
					TargetUserName:        msg.TargetUserName,
					TargetUserDisplayName: msg.TargetDisplayName,
					SubLevel:              msg.Level,
					UserID:                msg.SenderUserId,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_SUB_GIFT,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) ChatClear(
	ctx context.Context,
	msg *events.ChatClearMessage,
) (*emptypb.Empty, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeOnChatClear,
				shared.EvenData{
					ChannelID: msg.BaseInfo.ChannelId,
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
				msg.BaseInfo.ChannelId,
				api_events.TwirEventType_CHAT_CLEAR,
				msg,
			)
		},
	)

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_CHAT_CLEAR,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) Donate(
	ctx context.Context,
	msg *events.DonateMessage,
) (*emptypb.Empty, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeDonate,
				shared.EvenData{
					ChannelID:      msg.BaseInfo.ChannelId,
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
				msg.BaseInfo.ChannelId,
				api_events.TwirEventType_DONATE,
				msg,
			)
		},
	)

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_DONATE,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) KeywordMatched(
	ctx context.Context,
	msg *events.KeywordMatchedMessage,
) (*emptypb.Empty, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeKeywordMatched,
				shared.EvenData{
					ChannelID:       msg.BaseInfo.ChannelId,
					UserName:        msg.UserName,
					UserDisplayName: msg.UserDisplayName,
					KeywordName:     msg.KeywordName,
					KeywordResponse: msg.KeywordResponse,
					KeywordID:       msg.KeywordId,
					UserID:          msg.UserId,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_KEYWORD_USED,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) GreetingSended(
	ctx context.Context,
	msg *events.GreetingSendedMessage,
) (
	*emptypb.Empty, error,
) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypeGreetingSended,
				shared.EvenData{
					ChannelID:       msg.BaseInfo.ChannelId,
					UserName:        msg.UserName,
					UserDisplayName: msg.UserDisplayName,
					UserID:          msg.UserId,
					GreetingText:    msg.GreetingText,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_GREETING_SENDED,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) PollBegin(
	ctx context.Context,
	msg *events.PollBeginMessage,
) (*emptypb.Empty, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypePollBegin,
				shared.EvenData{
					ChannelID:       msg.BaseInfo.ChannelId,
					UserName:        msg.UserName,
					UserDisplayName: msg.UserDisplayName,
					PollTitle:       msg.Info.Title,
					PollOptionsNames: strings.Join(
						lo.Map(
							msg.Info.Choices, func(item *events.PollInfo_Choice, _ int) string {
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

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_POLL_STARTED,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) PollProgress(
	ctx context.Context,
	msg *events.PollProgressMessage,
) (
	*emptypb.Empty, error,
) {
	totalVotes := lo.Reduce(
		msg.Info.Choices, func(acc int, item *events.PollInfo_Choice, _ int) int {
			return acc + int(item.Votes)
		}, 0,
	)

	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypePollProgress,
				shared.EvenData{
					ChannelID:       msg.BaseInfo.ChannelId,
					UserDisplayName: msg.UserDisplayName,
					PollTitle:       msg.Info.Title,
					PollOptionsNames: strings.Join(
						lo.Map(
							msg.Info.Choices, func(item *events.PollInfo_Choice, _ int) string {
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

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_POLL_VOTED,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) PollEnd(
	ctx context.Context,
	msg *events.PollEndMessage,
) (*emptypb.Empty, error) {
	totalVotes := lo.Reduce(
		msg.Info.Choices, func(acc int, item *events.PollInfo_Choice, _ int) int {
			return acc + int(item.Votes)
		}, 0,
	)

	// find most total votes in choices
	winner := lo.MaxBy(
		msg.Info.Choices, func(a *events.PollInfo_Choice, b *events.PollInfo_Choice) bool {
			return a.Votes > b.Votes
		},
	)

	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypePollEnd,
				shared.EvenData{
					ChannelID:       msg.BaseInfo.ChannelId,
					UserDisplayName: msg.UserDisplayName,
					PollTitle:       msg.Info.Title,
					PollOptionsNames: strings.Join(
						lo.Map(
							msg.Info.Choices, func(item *events.PollInfo_Choice, _ int) string {
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

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_POLL_ENDED,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) PredictionBegin(
	ctx context.Context, msg *events.PredictionBeginMessage,
) (*emptypb.Empty, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypePredictionBegin,
				shared.EvenData{
					ChannelID:       msg.BaseInfo.ChannelId,
					UserName:        msg.UserName,
					UserDisplayName: msg.UserDisplayName,
					PredictionTitle: msg.Info.Title,
					PredictionOptionsNames: strings.Join(
						lo.Map(
							msg.Info.Outcomes, func(item *events.PredictionInfo_OutCome, _ int) string {
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

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_PREDICTION_STARTED,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) PredictionProgress(
	ctx context.Context, msg *events.PredictionProgressMessage,
) (*emptypb.Empty, error) {
	totalPoints := lo.Reduce(
		msg.Info.Outcomes, func(acc int, item *events.PredictionInfo_OutCome, _ int) int {
			return acc + int(item.ChannelPoints)
		}, 0,
	)

	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypePredictionProgress,
				shared.EvenData{
					ChannelID:       msg.BaseInfo.ChannelId,
					UserDisplayName: msg.UserDisplayName,
					PredictionTitle: msg.Info.Title,
					PredictionOptionsNames: strings.Join(
						lo.Map(
							msg.Info.Outcomes, func(item *events.PredictionInfo_OutCome, _ int) string {
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

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_PREDICTION_VOTED,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) PredictionLock(
	ctx context.Context,
	msg *events.PredictionLockMessage,
) (
	*emptypb.Empty, error,
) {
	totalPoints := lo.Reduce(
		msg.Info.Outcomes, func(acc int, item *events.PredictionInfo_OutCome, _ int) int {
			return acc + int(item.ChannelPoints)
		}, 0,
	)

	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypePredictionLock,
				shared.EvenData{
					ChannelID:       msg.BaseInfo.ChannelId,
					UserDisplayName: msg.UserDisplayName,
					PredictionTitle: msg.Info.Title,
					PredictionOptionsNames: strings.Join(
						lo.Map(
							msg.Info.Outcomes, func(item *events.PredictionInfo_OutCome, _ int) string {
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

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_PREDICTION_LOCKED,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) PredictionEnd(
	ctx context.Context,
	msg *events.PredictionEndMessage,
) (
	*emptypb.Empty, error,
) {
	totalPoints := lo.Reduce(
		msg.Info.Outcomes, func(acc int, item *events.PredictionInfo_OutCome, _ int) int {
			return acc + int(item.ChannelPoints)
		}, 0,
	)

	winner, _ := lo.Find(
		msg.Info.Outcomes, func(item *events.PredictionInfo_OutCome) bool {
			return item.Id == msg.WinningOutcomeId
		},
	)

	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventTypePredictionEnd,
				shared.EvenData{
					ChannelID:       msg.BaseInfo.ChannelId,
					UserDisplayName: msg.UserDisplayName,
					PredictionTitle: msg.Info.Title,
					PredictionOptionsNames: strings.Join(
						lo.Map(
							msg.Info.Outcomes, func(item *events.PredictionInfo_OutCome, _ int) string {
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

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_PREDICTION_ENDED,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) StreamFirstUserJoin(
	ctx context.Context, msg *events.StreamFirstUserJoinMessage,
) (*emptypb.Empty, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventStreamFirstUserJoin,
				shared.EvenData{
					ChannelID: msg.BaseInfo.ChannelId,
					UserName:  msg.UserName,
				},
			)
			if err != nil {
				c.logger.Error("Error execute workflow", slog.Any("err", err))
			}
		},
	)

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_FIRST_USER_MESSAGE,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) ChannelBan(
	ctx context.Context,
	msg *events.ChannelBanMessage,
) (*emptypb.Empty, error) {
	wg := utils.NewGoroutinesGroup()

	wg.Go(
		func() {
			err := c.eventsWorkflow.Execute(
				ctx,
				model.EventChannelBan,
				shared.EvenData{
					ChannelID:            msg.BaseInfo.ChannelId,
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
				msg.BaseInfo.ChannelId,
				api_events.TwirEventType_USER_BANNED,
				msg,
			)
		},
	)

	wg.Go(
		func() {
			_, err := c.websocketsGrpc.TriggerKappagenByEvent(
				ctx, &websockets.TriggerKappagenByEventRequest{
					ChannelId: msg.BaseInfo.ChannelId,
					Event:     api_events.TwirEventType_USER_BANNED,
				},
			)
			if err != nil {
				c.logger.Error("Error trigger kappagen by event", slog.Any("err", err))
			}
		},
	)

	wg.Wait()

	return &emptypb.Empty{}, nil
}
