package kick

import (
	"context"
	"log/slog"
	"strconv"
	"time"

	"github.com/twirapp/twir/apps/eventsub/internal/channelbinding"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	"github.com/twirapp/twir/libs/repositories/channels"
	"go.uber.org/fx"
)

type SubscriptionLister interface {
	ListSubscriptions(ctx context.Context, broadcasterUserID int) ([]SubscriptionInfo, error)
	Subscribe(ctx context.Context, binding channelplatformsmodel.ChannelPlatform) error
}

type ResubscribeJob struct {
	subManager   SubscriptionLister
	channelsRepo channels.Repository
	logger       *slog.Logger
	config       cfg.Config
	interval     time.Duration
}

type ResubscribeJobOpts struct {
	fx.In

	Lc fx.Lifecycle

	SubManager   *SubscriptionManager
	ChannelsRepo channels.Repository
	Logger       *slog.Logger
	Config       cfg.Config
}

func NewResubscribeJob(opts ResubscribeJobOpts) *ResubscribeJob {
	j := &ResubscribeJob{
		subManager:   opts.SubManager,
		channelsRepo: opts.ChannelsRepo,
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
	kickChannels, err := j.channelsRepo.GetAllByBindingPlatform(ctx, platform.PlatformKick)
	if err != nil {
		j.logger.ErrorContext(ctx, "resubscribe job: failed to list kick channels", logger.Error(err))
		return
	}

	for _, ch := range kickChannels {
		binding, ok := channelbinding.Find(ch, platform.PlatformKick)
		if !ok || !binding.Enabled {
			continue
		}

		kickChannelID := binding.UserID.String()

		broadcasterUserID, err := strconv.Atoi(binding.PlatformChannelID)
		if err != nil {
			j.logger.ErrorContext(
				ctx,
				"resubscribe job: failed to parse kick platform channel ID",
				slog.String("channel_id", ch.ID.String()),
				slog.String("kick_user_id", kickChannelID),
				slog.String("kick_platform_channel_id", binding.PlatformChannelID),
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

		if err := j.subManager.Subscribe(ctx, binding); err != nil {
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
