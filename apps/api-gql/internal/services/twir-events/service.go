package twir_events

import (
	"context"

	"github.com/goccy/go-json"
	"github.com/twirapp/twir/libs/wsrouter"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/events"
	bustwitch "github.com/twirapp/twir/libs/bus-core/twitch"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	WsRouter wsrouter.WsRouter
	TwirBus  *buscore.Bus
}

type Service struct {
	wsRouter wsrouter.WsRouter
	twirBus  *buscore.Bus
}

func New(opts Opts) *Service {
	s := &Service{
		wsRouter: opts.WsRouter,
		twirBus:  opts.TwirBus,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				if err := s.twirBus.Events.Follow.SubscribeGroup("api", s.follow); err != nil {
					return err
				}

				if err := s.twirBus.Events.Subscribe.SubscribeGroup(
					"api",
					s.subscribe,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.ReSubscribe.SubscribeGroup(
					"api",
					s.reSubscribe,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.RedemptionCreated.SubscribeGroup(
					"api",
					s.redemptionCreated,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.SubGift.SubscribeGroup("api", s.subGift); err != nil {
					return err
				}

				if err := s.twirBus.Events.CommandUsed.SubscribeGroup(
					"api",
					s.commandUsed,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.FirstUserMessage.SubscribeGroup(
					"api",
					s.firstUserMessage,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.Raided.SubscribeGroup("api", s.raided); err != nil {
					return err
				}

				if err := s.twirBus.Events.TitleOrCategoryChanged.SubscribeGroup(
					"api",
					s.titleOrCategoryChanged,
				); err != nil {
					return err
				}

				if err := s.twirBus.Channel.StreamOnline.SubscribeGroup(
					"api",
					s.streamOnline,
				); err != nil {
					return err
				}

				if err := s.twirBus.Channel.StreamOffline.SubscribeGroup(
					"api",
					s.streamOffline,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.ChatClear.SubscribeGroup(
					"api",
					s.chatClear,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.Donate.SubscribeGroup("api", s.donate); err != nil {
					return err
				}

				if err := s.twirBus.Events.KeywordMatched.SubscribeGroup(
					"api",
					s.keywordMatched,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.GreetingSended.SubscribeGroup(
					"api",
					s.greetingSended,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.PollBegin.SubscribeGroup(
					"api",
					s.pollBegin,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.PollProgress.SubscribeGroup(
					"api",
					s.pollProgress,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.PollEnd.SubscribeGroup("api", s.pollEnd); err != nil {
					return err
				}

				if err := s.twirBus.Events.PredictionBegin.SubscribeGroup(
					"api",
					s.predictionBegin,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.PredictionProgress.SubscribeGroup(
					"api",
					s.predictionProgress,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.PredictionLock.SubscribeGroup(
					"api",
					s.predictionLock,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.PredictionEnd.SubscribeGroup(
					"api",
					s.predictionEnd,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.StreamFirstUserJoin.SubscribeGroup(
					"api",
					s.streamFirstUserJoin,
				); err != nil {
					return err
				}
				if err := s.twirBus.Events.ChannelBan.SubscribeGroup(
					"api",
					s.channelBan,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.ChannelUnbanRequestCreate.SubscribeGroup(
					"api",
					s.channelUnbanRequestCreate,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.ChannelUnbanRequestResolve.SubscribeGroup(
					"api",
					s.channelUnbanRequestResolve,
				); err != nil {
					return err
				}

				if err := s.twirBus.Events.ChannelMessageDelete.SubscribeGroup(
					"api",
					s.channelMessageDelete,
				); err != nil {
					return err
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				s.twirBus.Events.Follow.Unsubscribe()
				s.twirBus.Events.Subscribe.Unsubscribe()
				s.twirBus.Events.ReSubscribe.Unsubscribe()
				s.twirBus.Events.RedemptionCreated.Unsubscribe()
				s.twirBus.Events.SubGift.Unsubscribe()
				s.twirBus.Events.FirstUserMessage.Unsubscribe()
				s.twirBus.Events.Raided.Unsubscribe()
				s.twirBus.Channel.StreamOnline.Unsubscribe()
				s.twirBus.Channel.StreamOffline.Unsubscribe()
				s.twirBus.Events.ChatClear.Unsubscribe()
				s.twirBus.Events.Donate.Unsubscribe()
				s.twirBus.Events.GreetingSended.Unsubscribe()
				s.twirBus.Events.PollBegin.Unsubscribe()
				s.twirBus.Events.PollProgress.Unsubscribe()
				s.twirBus.Events.PollEnd.Unsubscribe()
				s.twirBus.Events.PredictionProgress.Unsubscribe()
				s.twirBus.Events.PredictionLock.Unsubscribe()
				s.twirBus.Events.PredictionEnd.Unsubscribe()
				s.twirBus.Events.StreamFirstUserJoin.Unsubscribe()
				s.twirBus.Events.ChannelBan.Unsubscribe()
				s.twirBus.Events.ChannelUnbanRequestCreate.Unsubscribe()
				s.twirBus.Events.ChannelUnbanRequestResolve.Unsubscribe()
				s.twirBus.Events.ChannelMessageDelete.Unsubscribe()

				return nil
			},
		},
	)

	return s
}

func CreateSubscribeKey(channelID string) string {
	return "api.twirEvents." + channelID
}

type Message struct {
	EventName string
	Data      []byte
}

func createMessage(eventName string, data any) Message {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	return Message{
		EventName: eventName,
		Data:      dataBytes,
	}
}

func (s *Service) follow(ctx context.Context, msg events.FollowMessage) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(
		CreateSubscribeKey(msg.BaseInfo.ChannelID),
		createMessage(events.FollowSubject, msg),
	)
}

func (s *Service) subscribe(ctx context.Context, msg events.SubscribeMessage) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(
		CreateSubscribeKey(msg.BaseInfo.ChannelID),
		createMessage(events.SubscribeSubject, msg),
	)
}

func (s *Service) reSubscribe(ctx context.Context, msg events.ReSubscribeMessage) (
	struct{},
	error,
) {
	return struct{}{}, s.wsRouter.Publish(
		CreateSubscribeKey(msg.BaseInfo.ChannelID),
		createMessage(events.ReSubscribeSubject, msg),
	)
}

func (s *Service) redemptionCreated(
	ctx context.Context,
	msg events.RedemptionCreatedMessage,
) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(
		CreateSubscribeKey(msg.BaseInfo.ChannelID),
		createMessage(events.RedemptionCreatedSubject, msg),
	)
}

func (s *Service) subGift(ctx context.Context, msg events.SubGiftMessage) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(
		CreateSubscribeKey(msg.BaseInfo.ChannelID),
		createMessage(events.SubGiftSubject, msg),
	)
}

func (s *Service) commandUsed(ctx context.Context, msg events.CommandUsedMessage) (
	struct{},
	error,
) {
	return struct{}{}, s.wsRouter.Publish(
		CreateSubscribeKey(msg.BaseInfo.ChannelID),
		createMessage(events.CommandUsedSubject, msg),
	)
}

func (s *Service) firstUserMessage(
	ctx context.Context,
	msg events.FirstUserMessageMessage,
) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(
		CreateSubscribeKey(msg.BaseInfo.ChannelID),
		createMessage(events.FirstUserMessageSubject, msg),
	)
}

func (s *Service) raided(ctx context.Context, msg events.RaidedMessage) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(
		CreateSubscribeKey(msg.BaseInfo.ChannelID),
		createMessage(events.RaidedSubject, msg),
	)
}

func (s *Service) titleOrCategoryChanged(
	ctx context.Context,
	msg events.TitleOrCategoryChangedMessage,
) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}

func (s *Service) streamOnline(ctx context.Context, msg bustwitch.StreamOnlineMessage) (
	struct{},
	error,
) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.ChannelID), msg)
}

func (s *Service) streamOffline(ctx context.Context, msg bustwitch.StreamOfflineMessage) (
	struct{},
	error,
) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.ChannelID), msg)
}

func (s *Service) chatClear(ctx context.Context, msg events.ChatClearMessage) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}

func (s *Service) donate(ctx context.Context, msg events.DonateMessage) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}

func (s *Service) keywordMatched(ctx context.Context, msg events.KeywordMatchedMessage) (
	struct{},
	error,
) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}

func (s *Service) greetingSended(ctx context.Context, msg events.GreetingSendedMessage) (
	struct{},
	error,
) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}

func (s *Service) pollBegin(ctx context.Context, msg events.PollBeginMessage) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}

func (s *Service) pollProgress(ctx context.Context, msg events.PollProgressMessage) (
	struct{},
	error,
) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}

func (s *Service) pollEnd(ctx context.Context, msg events.PollEndMessage) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}

func (s *Service) predictionBegin(ctx context.Context, msg events.PredictionBeginMessage) (
	struct{},
	error,
) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}

func (s *Service) predictionProgress(
	ctx context.Context,
	msg events.PredictionProgressMessage,
) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}

func (s *Service) predictionLock(ctx context.Context, msg events.PredictionLockMessage) (
	struct{},
	error,
) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}

func (s *Service) predictionEnd(ctx context.Context, msg events.PredictionEndMessage) (
	struct{},
	error,
) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}

func (s *Service) streamFirstUserJoin(
	ctx context.Context,
	msg events.StreamFirstUserJoinMessage,
) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}

func (s *Service) channelBan(ctx context.Context, msg events.ChannelBanMessage) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}

func (s *Service) channelUnbanRequestCreate(
	ctx context.Context,
	msg events.ChannelUnbanRequestCreateMessage,
) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}

func (s *Service) channelUnbanRequestResolve(
	ctx context.Context,
	msg events.ChannelUnbanRequestResolveMessage,
) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}

func (s *Service) channelMessageDelete(
	ctx context.Context,
	msg events.ChannelMessageDeleteMessage,
) (struct{}, error) {
	return struct{}{}, s.wsRouter.Publish(CreateSubscribeKey(msg.BaseInfo.ChannelID), msg)
}
