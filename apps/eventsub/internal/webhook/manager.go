package webhook

import (
	"context"
	"fmt"
	"log/slog"

	httpserver "github.com/twirapp/twir/apps/eventsub/internal/http"
	"github.com/twirapp/twir/apps/eventsub/internal/kick"
	"github.com/twirapp/twir/apps/eventsub/internal/tunnel"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/crypto"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/channels"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	"go.uber.org/fx"
)

type Platform interface {
	Name() string
	SubscribeAll(ctx context.Context, channelID string, token string) error
	UnsubscribeAll(ctx context.Context, channelID string) error
	SetCallbackBaseURL(baseURL string)
}

type Manager struct {
	config       cfg.Config
	logger       *slog.Logger
	tunnelMgr    tunnel.Manager
	kickSubMgr   *kick.SubscriptionManager
	channelsRepo channels.Repository
	kickBotsRepo kickbotsrepository.Repository
	platforms    []Platform
}

type Opts struct {
	fx.In

	Lc fx.Lifecycle

	Config       cfg.Config
	Logger       *slog.Logger
	KickSubMgr   *kick.SubscriptionManager
	ChannelsRepo channels.Repository
	KickBotsRepo kickbotsrepository.Repository
	Server       *httpserver.Server
}

func NewManager(opts Opts) *Manager {
	m := &Manager{
		config:       opts.Config,
		logger:       opts.Logger,
		kickSubMgr:   opts.KickSubMgr,
		channelsRepo: opts.ChannelsRepo,
		kickBotsRepo: opts.KickBotsRepo,
	}

	if opts.Config.IsDevelopment() && opts.Config.TunnelEnabled {
		m.tunnelMgr = tunnel.NewPinggyManager(opts.Config, opts.Logger)
	}

	m.platforms = []Platform{m.kickSubMgr}

	opts.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return m.start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			if m.tunnelMgr != nil {
				return m.tunnelMgr.Stop(ctx)
			}
			return nil
		},
	})

	return m
}

func (m *Manager) start(ctx context.Context) error {
	callbackBaseURL := m.config.SiteBaseUrl

	if m.config.IsDevelopment() && m.tunnelMgr != nil {
		publicURL, err := m.tunnelMgr.Start(ctx)
		if err != nil {
			return fmt.Errorf("webhook manager: start tunnel: %w", err)
		}
		callbackBaseURL = publicURL
	}

	for _, p := range m.platforms {
		p.SetCallbackBaseURL(callbackBaseURL)
	}

	if m.config.IsDevelopment() {
		if err := m.unsubscribeAllPlatforms(ctx); err != nil {
			m.logger.ErrorContext(ctx, "webhook manager: unsubscribe all failed", logger.Error(err))
		}
	}

	if err := m.subscribeAllPlatforms(ctx); err != nil {
		m.logger.ErrorContext(ctx, "webhook manager: subscribe all failed", logger.Error(err))
	}

	return nil
}

func (m *Manager) unsubscribeAllPlatforms(ctx context.Context) error {
	hasKick := true
	kickChannels, err := m.channelsRepo.GetMany(ctx, channels.GetManyInput{
		HasKickUserID: &hasKick,
	})
	if err != nil {
		return fmt.Errorf("list kick channels: %w", err)
	}

	for _, ch := range kickChannels {
		if ch.KickUserID == nil {
			continue
		}

		if err := m.kickSubMgr.UnsubscribeAll(ctx, ch.KickUserID.String()); err != nil {
			m.logger.WarnContext(
				ctx,
				"webhook manager: failed to unsubscribe kick",
				slog.String("kick_user_id", ch.KickUserID.String()),
				logger.Error(err),
			)
		}
	}

	return nil
}

func (m *Manager) subscribeAllPlatforms(ctx context.Context) error {
	hasKick := true
	kickChannels, err := m.channelsRepo.GetMany(ctx, channels.GetManyInput{
		HasKickUserID: &hasKick,
	})
	if err != nil {
		return fmt.Errorf("list kick channels: %w", err)
	}

	for _, ch := range kickChannels {
		if ch.KickUserID == nil {
			continue
		}

		kickBot, err := m.kickBotsRepo.GetByKickUserID(ctx, *ch.KickUserID)
		if err != nil {
			m.logger.WarnContext(
				ctx,
				"webhook manager: failed to get kick bot for subscribe",
				slog.String("kick_user_id", ch.KickUserID.String()),
				logger.Error(err),
			)
			continue
		}

		accessToken, err := crypto.Decrypt(kickBot.AccessToken, m.config.TokensCipherKey)
		if err != nil {
			m.logger.WarnContext(
				ctx,
				"webhook manager: failed to decrypt token for kick subscribe",
				slog.String("kick_user_id", ch.KickUserID.String()),
				logger.Error(err),
			)
			continue
		}

		if err := m.kickSubMgr.SubscribeAll(ctx, ch.KickUserID.String(), accessToken); err != nil {
			m.logger.ErrorContext(
				ctx,
				"webhook manager: failed to subscribe kick",
				slog.String("kick_user_id", ch.KickUserID.String()),
				logger.Error(err),
			)
			continue
		}

		m.logger.InfoContext(
			ctx,
			"webhook manager: subscribed kick eventsub",
			slog.String("kick_user_id", ch.KickUserID.String()),
		)
	}

	return nil
}
