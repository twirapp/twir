package kick

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/google/uuid"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/channels"
	usersrepository "github.com/twirapp/twir/libs/repositories/users"
	"go.uber.org/fx"
)

type SubscriptionLister interface {
	ListSubscriptions(ctx context.Context, broadcasterUserID int) ([]SubscriptionInfo, error)
	SubscribeAll(ctx context.Context, kickChannelID uuid.UUID) error
}

type ResubscribeJob struct {
	subManager   SubscriptionLister
	channelsRepo channels.Repository
	usersRepo    usersrepository.Repository
	logger       *slog.Logger
	config       cfg.Config
	interval     time.Duration
}

type ResubscribeJobOpts struct {
	fx.In

	Lc fx.Lifecycle

	SubManager   *SubscriptionManager
	ChannelsRepo channels.Repository
	UsersRepo    usersrepository.Repository
	Logger       *slog.Logger
	Config       cfg.Config
}

func NewResubscribeJob(opts ResubscribeJobOpts) *ResubscribeJob {
	j := &ResubscribeJob{
		subManager:   opts.SubManager,
		channelsRepo: opts.ChannelsRepo,
		usersRepo:    opts.UsersRepo,
		logger:       opts.Logger,
		config:       opts.Config,
		interval:     23 * time.Hour,
	}

	stopCh := make(chan struct{})
	opts.Lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go j.Start(stopCh)
			return nil
		},
		OnStop: func(_ context.Context) error {
			close(stopCh)
			return nil
		},
	})

	return j
}

func (j *ResubscribeJob) Start(stopCh <-chan struct{}) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-stopCh:
			return
		case <-ticker.C:
			j.run(context.Background())
		}
	}
}

func (j *ResubscribeJob) run(ctx context.Context) {
	hasKick := true
	kickChannels, err := j.channelsRepo.GetMany(ctx, channels.GetManyInput{
		HasKickUserID: &hasKick,
	})
	if err != nil {
		j.logger.ErrorContext(ctx, "resubscribe job: failed to list kick channels", logger.Error(err))
		return
	}

	for _, ch := range kickChannels {
		if ch.KickUserID == nil || !ch.IsEnabled {
			continue
		}

		kickChannelID := ch.KickUserID.String()

		user, err := j.usersRepo.GetByID(ctx, *ch.KickUserID)
		if err != nil {
			j.logger.ErrorContext(
				ctx,
				"resubscribe job: failed to get user for kick channel",
				slog.String("channel_id", ch.ID.String()),
				slog.String("kick_user_id", ch.KickUserID.String()),
				logger.Error(err),
			)
			continue
		}

		broadcasterUserID, err := strconv.Atoi(user.PlatformID)
		if err != nil {
			j.logger.ErrorContext(
				ctx,
				"resubscribe job: failed to parse kick platform user ID",
				slog.String("channel_id", ch.ID.String()),
				slog.String("kick_user_id", ch.KickUserID.String()),
				logger.Error(err),
			)
			continue
		}

		subs, err := j.subManager.ListSubscriptions(ctx, broadcasterUserID)
		if err != nil {
			j.logger.ErrorContext(
				ctx,
				"resubscribe job: failed to list kick subscriptions",
				slog.String("channel_id", ch.ID.String()),
				slog.String("kick_user_id", kickChannelID),
				logger.Error(err),
			)
			continue
		}

		if allEventTypesPresent(subs) {
			continue
		}

		if err := j.subManager.SubscribeAll(ctx, *ch.KickUserID); err != nil {
			j.logger.ErrorContext(
				ctx,
				"resubscribe job: failed to re-subscribe kick eventsub",
				slog.String("channel_id", ch.ID.String()),
				slog.String("kick_user_id", kickChannelID),
				logger.Error(err),
			)
			continue
		}

		j.logger.InfoContext(
			ctx,
			"re-subscribed kick eventsub for channel",
			slog.String("channel_id", ch.ID.String()),
			slog.String("kick_user_id", kickChannelID),
		)
	}
}

func allEventTypesPresent(subs []SubscriptionInfo) bool {
	present := make(map[string]struct{}, len(subs))
	for _, s := range subs {
		present[s.Event] = struct{}{}
	}
	for _, et := range EventTypes {
		if _, ok := present[et]; !ok {
			return false
		}
	}
	return true
}
