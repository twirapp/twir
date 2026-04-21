package kick

import (
	"context"
	"errors"
	"log/slog"
	"time"

	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/crypto"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/channels"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	"go.uber.org/fx"
)

type SubscriptionLister interface {
	ListSubscriptions(ctx context.Context, broadcasterToken string) ([]SubscriptionInfo, error)
	SubscribeAll(ctx context.Context, kickChannelID string, broadcasterToken string) error
}

type ResubscribeJob struct {
	subManager   SubscriptionLister
	channelsRepo channels.Repository
	kickBotsRepo kickbotsrepository.Repository
	logger       *slog.Logger
	config       cfg.Config
	interval     time.Duration
}

type ResubscribeJobOpts struct {
	fx.In

	Lc fx.Lifecycle

	SubManager   *SubscriptionManager
	ChannelsRepo channels.Repository
	KickBotsRepo kickbotsrepository.Repository
	Logger       *slog.Logger
	Config       cfg.Config
}

func NewResubscribeJob(opts ResubscribeJobOpts) *ResubscribeJob {
	j := &ResubscribeJob{
		subManager:   opts.SubManager,
		channelsRepo: opts.ChannelsRepo,
		kickBotsRepo: opts.KickBotsRepo,
		logger:       opts.Logger,
		config:       opts.Config,
		interval:     23 * time.Hour,
	}

	opts.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go j.Start(ctx)
			return nil
		},
		OnStop: func(_ context.Context) error {
			return nil
		},
	})

	return j
}

func (j *ResubscribeJob) Start(ctx context.Context) {
	ticker := time.NewTicker(j.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			j.run(ctx)
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
		if ch.KickUserID == nil {
			continue
		}

		if ch.KickBotID == nil {
			j.logger.InfoContext(
				ctx,
				"resubscribe job: channel has no kick bot assigned, skipping",
				slog.String("channel_id", ch.ID.String()),
				slog.String("kick_user_id", ch.KickUserID.String()),
			)
			continue
		}

		kickBot, err := j.kickBotsRepo.GetByID(ctx, *ch.KickBotID)
		if err != nil {
			if errors.Is(err, kickbotsrepository.ErrNotFound) {
				j.logger.InfoContext(
					ctx,
					"resubscribe job: kick bot not found, skipping",
					slog.String("channel_id", ch.ID.String()),
					slog.String("kick_bot_id", ch.KickBotID.String()),
				)
			} else {
				j.logger.ErrorContext(
					ctx,
					"resubscribe job: failed to get kick bot",
					slog.String("channel_id", ch.ID.String()),
					slog.String("kick_bot_id", ch.KickBotID.String()),
					logger.Error(err),
				)
			}
			continue
		}

		accessToken, err := crypto.Decrypt(kickBot.AccessToken, j.config.TokensCipherKey)
		if err != nil {
			j.logger.ErrorContext(
				ctx,
				"resubscribe job: failed to decrypt kick access token",
				slog.String("channel_id", ch.ID.String()),
				slog.String("kick_bot_id", ch.KickBotID.String()),
				logger.Error(err),
			)
			continue
		}

		kickChannelID := ch.KickUserID.String()

		subs, err := j.subManager.ListSubscriptions(ctx, accessToken)
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

		if err := j.subManager.SubscribeAll(ctx, kickChannelID, accessToken); err != nil {
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
		present[s.Type] = struct{}{}
	}
	for _, et := range EventTypes {
		if _, ok := present[et]; !ok {
			return false
		}
	}
	return true
}
