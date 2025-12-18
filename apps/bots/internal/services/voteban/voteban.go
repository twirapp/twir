package voteban

import (
	"context"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/twirapp/twir/apps/bots/internal/services/channel"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"github.com/twirapp/twir/libs/logger"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	TwirBus        *buscore.Bus
	Logger         *slog.Logger
	ChannelService *channel.Service
}

func New(opts Opts) *Service {
	s := &Service{
		inProgressVotebans: make(map[voteBanChannelId]*session),
		twirBus:            opts.TwirBus,
		logger:             opts.Logger,
		channelService:     opts.ChannelService,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return s.twirBus.Bots.VotebanRegister.SubscribeGroup("bots", s.tryRegisterVoteban)
			},
			OnStop: func(ctx context.Context) error {
				s.twirBus.Bots.VotebanRegister.Unsubscribe()
				return nil
			},
		},
	)

	return s
}

type voteBanChannelId = string

type Service struct {
	inProgressVotebans map[voteBanChannelId]*session
	mu                 sync.RWMutex
	twirBus            *buscore.Bus
	logger             *slog.Logger
	channelService     *channel.Service
}

func (s *Service) TryRegisterVote(msg twitch.TwitchChatMessage) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	sess, exists := s.inProgressVotebans[msg.BroadcasterUserId]
	if !exists {
		return
	}

	sess.tryRegisterVote(msg)
}

func (s *Service) tryRegisterVoteban(_ context.Context, req bots.VotebanRegisterRequest) (bots.VotebanRegisterResponse, error) {
	s.mu.RLock()
	if _, ok := s.inProgressVotebans[req.Data.ChannelID]; ok {
		s.mu.RUnlock()
		return bots.VotebanRegisterResponse{
			AlreadyInProgress: true,
		}, nil
	}
	s.mu.RUnlock()

	sess := newSession(
		req.Data,
		req.TargerUser.UserId,
		req.TargerUser.UserLogin,
		req.InitiatorIsModerator,
	)

	s.mu.Lock()
	if _, ok := s.inProgressVotebans[req.Data.ChannelID]; ok {
		return bots.VotebanRegisterResponse{
			AlreadyInProgress: true,
		}, nil
	}

	s.inProgressVotebans[req.Data.ChannelID] = sess
	s.mu.Unlock()

	s.logger.Info(
		"voteban started",
		slog.String("channel_id", req.Data.ChannelID),
		slog.Group(
			"user",
			slog.String("id", req.TargerUser.UserId),
			slog.String("name", req.TargerUser.UserLogin),
		),
	)

	go func() {
		result, ok := sess.waitResult()
		if !ok {
			return
		}

		s.logger.Info(
			"voteban finished",
			slog.String("channel_id", result.channelId),
			slog.Group(
				"user",
				slog.String("id", result.targetUserId),
			),
		)

		if err := s.processSessionResult(req.Data.ChannelID, result); err != nil {
			s.logger.Error("failed to process voteban session result", logger.Error(err))
		}
	}()

	return bots.VotebanRegisterResponse{}, nil
}

func (s *Service) processSessionResult(channelId string, result sessionResult) error {
	defer func() {
		s.mu.Lock()
		delete(s.inProgressVotebans, channelId)
		s.mu.Unlock()
	}()

	if !result.haveDecision {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()

	group, _ := errgroup.WithContext(ctx)

	group.Go(
		func() error {
			if err := s.channelService.SendMessage(
				ctx, bots.SendMessageRequest{
					ChannelId:         result.channelId,
					Message:           result.message,
					IsAnnounce:        true,
					SkipRateLimits:    true,
					SkipToxicityCheck: false,
					AnnounceColor:     bots.RandomAnnounceColor(),
				},
			); err != nil {
				return fmt.Errorf("send message: %w", err)
			}

			return nil
		},
	)

	if result.isBan {
		group.Go(
			func() error {
				if err := s.channelService.Ban(
					ctx, bots.BanRequest{
						ChannelID:      result.channelId,
						UserID:         result.targetUserId,
						Reason:         result.message,
						BanTime:        result.banDuration,
						IsModerator:    result.isModerator,
						AddModAfterBan: result.isModerator,
					},
				); err != nil {
					return fmt.Errorf("ban: %w", err)
				}

				return nil
			},
		)
	}

	return group.Wait()
}
