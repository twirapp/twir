package grpc_impl

import (
	"context"
	"strings"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/events/internal"
	"github.com/satont/twir/apps/events/internal/grpc_impl/chat_alerts"
	model "github.com/satont/twir/libs/gomodels"
	events2 "github.com/satont/twir/libs/grpc/generated/api/events"
	"github.com/satont/twir/libs/grpc/generated/events"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"google.golang.org/protobuf/types/known/emptypb"
)

type EventsGrpcImplementation struct {
	events.UnimplementedEventsServer

	services *internal.Services
}

func NewEvents(services *internal.Services) *EventsGrpcImplementation {
	return &EventsGrpcImplementation{
		services: services,
	}
}

func (c *EventsGrpcImplementation) Follow(
	ctx context.Context,
	msg *events.FollowMessage,
) (*emptypb.Empty, error) {
	c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			UserID:          msg.UserId,
		},
		model.EventTypeFollow,
	)

	chatAlerts, err := chat_alerts.New(msg.BaseInfo.ChannelId, c.services)
	if err == nil {
		chatAlerts.Follow(ctx, msg)
	}

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_FOLLOW,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) Subscribe(
	ctx context.Context,
	msg *events.SubscribeMessage,
) (*emptypb.Empty, error) {
	c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			SubLevel:        msg.Level,
			UserID:          msg.UserId,
		},
		model.EventTypeSubscribe,
	)

	chatAlerts, err := chat_alerts.New(msg.BaseInfo.ChannelId, c.services)
	if err == nil {
		chatAlerts.Subscribe(ctx, 1, msg.UserName, msg.BaseInfo.ChannelId)
	}

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_SUBSCRIBE,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) ReSubscribe(
	ctx context.Context,
	msg *events.ReSubscribeMessage,
) (*emptypb.Empty, error) {
	c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			SubLevel:        msg.Level,
			ResubMessage:    msg.Message,
			ResubMonths:     msg.Months,
			ResubStreak:     msg.Streak,
			UserID:          msg.UserId,
		},
		model.EventTypeResubscribe,
	)

	chatAlerts, err := chat_alerts.New(msg.BaseInfo.ChannelId, c.services)
	if err == nil {
		chatAlerts.Subscribe(ctx, int(msg.Months), msg.UserName, msg.BaseInfo.ChannelId)
	}

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_RESUBSCRIBE,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) RedemptionCreated(
	ctx context.Context,
	msg *events.RedemptionCreatedMessage,
) (*emptypb.Empty, error) {
	c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			RewardCost:      msg.RewardCost,
			RewardInput:     msg.Input,
			RewardName:      msg.RewardName,
			RewardID:        msg.Id,
			UserID:          msg.UserId,
		},
		model.EventTypeRedemptionCreated,
	)

	chatAlerts, err := chat_alerts.New(msg.BaseInfo.ChannelId, c.services)
	if err == nil {
		chatAlerts.Redemption(ctx, msg)
	}

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_REDEMPTION_CREATED,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) CommandUsed(
	ctx context.Context,
	msg *events.CommandUsedMessage,
) (
	*emptypb.Empty, error,
) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			CommandName:     msg.CommandName,
			CommandID:       msg.CommandId,
			CommandInput:    msg.CommandInput,
			UserID:          msg.UserId,
		},
		model.EventTypeCommandUsed,
	)

	if msg.DefaultCommandName != "kappagen" {
		c.services.WebsocketsGrpc.TriggerKappagenByEvent(
			ctx, &websockets.TriggerKappagenByEventRequest{
				ChannelId: msg.BaseInfo.ChannelId,
				Event:     events2.TwirEventType_COMMAND_USED,
			},
		)
	}

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) FirstUserMessage(
	ctx context.Context,
	msg *events.FirstUserMessageMessage,
) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			UserID:          msg.UserId,
		},
		model.EventTypeFirstUserMessage,
	)

	chatAlerts, err := chat_alerts.New(msg.BaseInfo.ChannelId, c.services)
	if err == nil {
		chatAlerts.FirstUserMessage(ctx, msg)
	}

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_FIRST_USER_MESSAGE,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) Raided(
	ctx context.Context,
	msg *events.RaidedMessage,
) (*emptypb.Empty, error) {
	c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			RaidViewers:     msg.Viewers,
			UserID:          msg.UserId,
		},
		model.EventTypeRaided,
	)

	chatAlerts, err := chat_alerts.New(msg.BaseInfo.ChannelId, c.services)
	if err == nil {
		chatAlerts.Raid(ctx, msg)
	}

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_RAIDED,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) TitleOrCategoryChanged(
	ctx context.Context, msg *events.TitleOrCategoryChangedMessage,
) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			OldStreamCategory: msg.OldCategory,
			NewStreamCategory: msg.NewCategory,
			OldStreamTitle:    msg.OldTitle,
			NewStreamTitle:    msg.NewTitle,
		},
		model.EventTypeTitleOrCategoryChanged,
	)

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_TITLE_OR_CATEGORY_CHANGED,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) StreamOnline(
	ctx context.Context,
	msg *events.StreamOnlineMessage,
) (
	*emptypb.Empty, error,
) {
	c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			StreamTitle:    msg.Title,
			StreamCategory: msg.Category,
		},
		model.EventTypeStreamOnline,
	)

	chatAlerts, err := chat_alerts.New(msg.BaseInfo.ChannelId, c.services)
	if err == nil {
		chatAlerts.StreamOnline(ctx, msg)
	}

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_STREAM_ONLINE,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) StreamOffline(
	ctx context.Context,
	msg *events.StreamOfflineMessage,
) (
	*emptypb.Empty, error,
) {
	c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{},
		model.EventTypeStreamOffline,
	)

	chatAlerts, err := chat_alerts.New(msg.BaseInfo.ChannelId, c.services)
	if err == nil {
		chatAlerts.StreamOffline(ctx, msg)
	}

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_STREAM_OFFLINE,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) SubGift(
	ctx context.Context,
	msg *events.SubGiftMessage,
) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			TargetUserName:        msg.TargetUserName,
			TargetUserDisplayName: msg.TargetDisplayName,
			SubLevel:              msg.Level,
			UserID:                msg.SenderUserId,
		},
		model.EventTypeSubGift,
	)

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_SUB_GIFT,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) ChatClear(
	ctx context.Context,
	msg *events.ChatClearMessage,
) (*emptypb.Empty, error) {
	c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{},
		model.EventTypeOnChatClear,
	)

	chatAlerts, err := chat_alerts.New(msg.BaseInfo.ChannelId, c.services)
	if err == nil {
		chatAlerts.ChatCleared(ctx, msg)
	}

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_CHAT_CLEAR,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) Donate(
	ctx context.Context,
	msg *events.DonateMessage,
) (*emptypb.Empty, error) {
	c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:       msg.UserName,
			DonateAmount:   msg.Amount,
			DonateCurrency: msg.Currency,
			DonateMessage:  msg.Message,
		},
		model.EventTypeDonate,
	)

	chatAlerts, err := chat_alerts.New(msg.BaseInfo.ChannelId, c.services)
	if err == nil {
		chatAlerts.Donation(ctx, msg)
	}

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_DONATE,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) KeywordMatched(
	ctx context.Context,
	msg *events.KeywordMatchedMessage,
) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			KeywordName:     msg.KeywordName,
			KeywordResponse: msg.KeywordResponse,
			KeywordID:       msg.KeywordId,
			UserID:          msg.UserId,
		},
		model.EventTypeKeywordMatched,
	)

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_KEYWORD_USED,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) GreetingSended(
	ctx context.Context,
	msg *events.GreetingSendedMessage,
) (
	*emptypb.Empty, error,
) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:        msg.UserName,
			UserDisplayName: msg.UserDisplayName,
			UserID:          msg.UserId,
			GreetingText:    msg.GreetingText,
		},
		model.EventTypeGreetingSended,
	)

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_GREETING_SENDED,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) PollBegin(
	ctx context.Context,
	msg *events.PollBeginMessage,
) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
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
		model.EventTypePollBegin,
	)

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_POLL_STARTED,
		},
	)

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

	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
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
			PollTotalVotes: totalVotes,
		},
		model.EventTypePollProgress,
	)

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_POLL_VOTED,
		},
	)

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

	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
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
			PollWinnerTitle:               winner.Title,
			PollWinnerBitsVotes:           int(winner.BitsVotes),
			PollWinnerChannelsPointsVotes: int(winner.ChannelsPointsVotes),
			PollWinnerTotalVotes:          int(winner.Votes),
			PollTotalVotes:                totalVotes,
		},
		model.EventTypePollEnd,
	)

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_POLL_ENDED,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) PredictionBegin(
	ctx context.Context, msg *events.PredictionBeginMessage,
) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
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
		model.EventTypePredictionBegin,
	)

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_PREDICTION_STARTED,
		},
	)

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

	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
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
			PredictionTotalChannelPoints: totalPoints,
		},
		model.EventTypePredictionProgress,
	)

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_PREDICTION_VOTED,
		},
	)

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

	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
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
			PredictionTotalChannelPoints: totalPoints,
		},
		model.EventTypePredictionLock,
	)

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_PREDICTION_LOCKED,
		},
	)

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

	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
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
			PredictionWinner: internal.PredictionOutCome{
				Title:       winner.Title,
				TotalUsers:  int(winner.Users),
				TotalPoints: int(winner.ChannelPoints),
				TopUsers:    predictionMapTopPredictors(winner.TopPredictors),
			},
			PredictionTotalChannelPoints: totalPoints,
		},
		model.EventTypePredictionEnd,
	)

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_PREDICTION_ENDED,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) StreamFirstUserJoin(
	ctx context.Context, msg *events.StreamFirstUserJoinMessage,
) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName: msg.UserName,
		},
		model.EventStreamFirstUserJoin,
	)

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_FIRST_USER_MESSAGE,
		},
	)

	return &emptypb.Empty{}, nil
}

func (c *EventsGrpcImplementation) ChannelBan(
	ctx context.Context,
	msg *events.ChannelBanMessage,
) (*emptypb.Empty, error) {
	go c.processEvent(
		msg.BaseInfo.ChannelId,
		internal.Data{
			UserName:             msg.UserLogin,
			UserDisplayName:      msg.UserName,
			ModeratorDisplayName: msg.ModeratorUserName,
			ModeratorName:        msg.ModeratorUserLogin,
			BanReason:            msg.Reason,
			BanEndsInMinutes:     msg.EndsAt,
		},
		model.EventChannelBan,
	)

	chatAlerts, err := chat_alerts.New(msg.BaseInfo.ChannelId, c.services)
	if err == nil {
		chatAlerts.Ban(ctx, msg)
	}

	c.services.WebsocketsGrpc.TriggerKappagenByEvent(
		ctx, &websockets.TriggerKappagenByEventRequest{
			ChannelId: msg.BaseInfo.ChannelId,
			Event:     events2.TwirEventType_USER_BANNED,
		},
	)

	return &emptypb.Empty{}, nil
}
