package voteban

import (
	"context"
	"log/slog"
	"sync"
	"time"

	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	TwirBus *buscore.Bus
	Logger  *slog.Logger
}

func New(opts Opts) *Service {
	s := &Service{
		inProgressVotebans: make(map[voteBanChannelId]*votebanSession),
		mu:                 sync.Mutex{},
		twirBus:            opts.TwirBus,
		logger:             opts.Logger,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				s.twirBus.Bots.VotebanRegister.SubscribeGroup("bots", s.tryRegisterVoteban)
				return nil
			},
			OnStop: func(ctx context.Context) error {
				s.twirBus.Bots.VotebanRegister.Unsubscribe()
				return nil
			},
		},
	)

	return s
}

type voteBanChannelId string

type Service struct {
	inProgressVotebans map[voteBanChannelId]*votebanSession
	mu                 sync.Mutex
	twirBus            *buscore.Bus
	logger             *slog.Logger
}

func (c *Service) tryRegisterVoteban(_ context.Context, req bots.VotebanRegisterRequest) (bots.VotebanRegisterResponse, error) {
	if _, ok := c.inProgressVotebans[voteBanChannelId(req.Data.ChannelID)]; ok {
		return bots.VotebanRegisterResponse{
			AlreadyInProgress: true,
		}, nil
	}

	session := createVotebanSession(
		req.Data,
		req.InitiatorUserID,
		req.TargerUser.UserId,
		req.TargerUser.UserLogin,
		req.InitiatorIsModerator,
	)
	c.mu.Lock()
	c.inProgressVotebans[voteBanChannelId(req.Data.ChannelID)] = session
	c.mu.Unlock()

	c.logger.Info(
		"voteban started",
		slog.String("channel_id", req.Data.ChannelID),
		slog.Group(
			"user",
			slog.String("id", req.TargerUser.UserId),
			slog.String("name", req.TargerUser.UserLogin),
		),
	)

	finishChann := session.start()

	go func() {
		for data := range finishChann {
			c.logger.Info(
				"voteban finished",
				slog.String("channel_id", data.channelId),
				slog.Group(
					"user",
					slog.String("id", data.targerUserId),
				),
			)

			c.mu.Lock()
			defer func() {
				delete(c.inProgressVotebans, voteBanChannelId(session.data.ChannelID))
				c.mu.Unlock()
			}()
			if !data.haveDecision {
				continue
			}

			ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
			defer cancel()

			c.twirBus.Bots.SendMessage.Publish(
				ctx, bots.SendMessageRequest{
					ChannelId:         data.channelId,
					Message:           data.message,
					IsAnnounce:        true,
					SkipRateLimits:    true,
					SkipToxicityCheck: false,
					AnnounceColor:     bots.RandomAnnounceColor(),
				},
			)

			if data.isBan {
				c.twirBus.Bots.BanUser.Publish(
					ctx, bots.BanRequest{
						ChannelID:      data.channelId,
						UserID:         data.targerUserId,
						Reason:         data.message,
						BanTime:        data.banDuration,
						IsModerator:    data.isModerator,
						AddModAfterBan: data.isModerator,
					},
				)
			}
		}
	}()

	return bots.VotebanRegisterResponse{}, nil
}

func (c *Service) HandleTwitchMessage(ctx context.Context, msg twitch.TwitchChatMessage) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	for channelId, s := range c.inProgressVotebans {
		if channelId != voteBanChannelId(msg.BroadcasterUserId) {
			continue
		}

		s.tryRegisterVote(msg)
	}

	return nil
}
