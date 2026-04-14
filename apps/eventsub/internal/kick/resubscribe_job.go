package kick

import (
	"context"
	"log/slog"
	"time"

	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/crypto"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	user_platform_accounts "github.com/twirapp/twir/libs/repositories/user_platform_accounts"
	"go.uber.org/fx"
)

type SubscriptionLister interface {
	ListSubscriptions(ctx context.Context, broadcasterToken string) ([]SubscriptionInfo, error)
	SubscribeAll(ctx context.Context, kickChannelID string, broadcasterToken string) error
}

type ResubscribeJob struct {
	subManager               SubscriptionLister
	userPlatformAccountsRepo user_platform_accounts.Repository
	logger                   *slog.Logger
	config                   cfg.Config
	interval                 time.Duration
}

type ResubscribeJobOpts struct {
	fx.In

	Lc fx.Lifecycle

	SubManager               *SubscriptionManager
	UserPlatformAccountsRepo user_platform_accounts.Repository
	Logger                   *slog.Logger
	Config                   cfg.Config
}

func NewResubscribeJob(opts ResubscribeJobOpts) *ResubscribeJob {
	j := &ResubscribeJob{
		subManager:               opts.SubManager,
		userPlatformAccountsRepo: opts.UserPlatformAccountsRepo,
		logger:                   opts.Logger,
		config:                   opts.Config,
		interval:                 23 * time.Hour,
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
	accounts, err := j.userPlatformAccountsRepo.GetAllByPlatform(ctx, platform.PlatformKick)
	if err != nil {
		j.logger.ErrorContext(ctx, "resubscribe job: failed to list kick platform accounts", slog.Any("error", err))
		return
	}

	for _, account := range accounts {
		accessToken, err := crypto.Decrypt(account.AccessToken, j.config.TokensCipherKey)
		if err != nil {
			j.logger.ErrorContext(
				ctx,
				"resubscribe job: failed to decrypt kick access token",
				slog.String("kick_channel_id", account.PlatformUserID),
				logger.Error(err),
			)
			continue
		}

		subs, err := j.subManager.ListSubscriptions(ctx, accessToken)
		if err != nil {
			j.logger.ErrorContext(
				ctx,
				"resubscribe job: failed to list kick subscriptions",
				slog.String("kick_channel_id", account.PlatformUserID),
				slog.Any("error", err),
			)
			continue
		}

		if allEventTypesPresent(subs) {
			continue
		}

		if err := j.subManager.SubscribeAll(ctx, account.PlatformUserID, accessToken); err != nil {
			j.logger.ErrorContext(
				ctx,
				"resubscribe job: failed to re-subscribe kick eventsub",
				slog.String("kick_channel_id", account.PlatformUserID),
				slog.Any("error", err),
			)
			continue
		}

		j.logger.InfoContext(
			ctx,
			"re-subscribed kick eventsub for channel",
			slog.String("kick_channel_id", account.PlatformUserID),
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
