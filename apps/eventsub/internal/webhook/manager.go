package webhook

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	httpserver "github.com/twirapp/twir/apps/eventsub/internal/http"
	"github.com/twirapp/twir/apps/eventsub/internal/kick"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/crypto"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/channels"
	kickbotsrepository "github.com/twirapp/twir/libs/repositories/kick_bots"
	"go.uber.org/fx"
)

type Platform interface {
	Name() string
	SubscribeAll(ctx context.Context, channelID string, token string, botID uuid.UUID, encryptedRefreshToken string) error
	UnsubscribeAll(ctx context.Context, channelID string, token string, botID uuid.UUID, encryptedRefreshToken string) error
	SetCallbackBaseURL(baseURL string)
}

type Manager struct {
	config       cfg.Config
	logger       *slog.Logger
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

	m.platforms = []Platform{m.kickSubMgr}

	opts.Lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			return m.start(ctx)
		},
		OnStop: func(ctx context.Context) error {
			return nil
		},
	})

	return m
}

func (m *Manager) start(ctx context.Context) error {
	callbackBaseURL := m.config.SiteBaseUrl

	m.logger.InfoContext(
		ctx,
		"webhook manager: starting",
		slog.String("callback_base_url", callbackBaseURL),
		slog.Bool("is_development", m.config.IsDevelopment()),
		slog.Int("platforms_count", len(m.platforms)),
	)

	for _, p := range m.platforms {
		p.SetCallbackBaseURL(callbackBaseURL)
	}

	if m.config.IsDevelopment() {
		m.logKickWebhookInstructions(ctx, callbackBaseURL)
		if err := m.unsubscribeAllPlatforms(ctx); err != nil {
			m.logger.ErrorContext(ctx, "webhook manager: unsubscribe all failed", logger.Error(err))
		}
	}

	if err := m.subscribeAllPlatforms(ctx); err != nil {
		m.logger.ErrorContext(ctx, "webhook manager: subscribe all failed", logger.Error(err))
	}

	m.logger.InfoContext(ctx, "webhook manager: started successfully")

	return nil
}

func (m *Manager) logKickWebhookInstructions(ctx context.Context, callbackBaseURL string) {
	webhookURL := callbackBaseURL + "/webhook/kick"
	m.logger.WarnContext(ctx,
		"============================================================\n"+
			"KICK WEBHOOK SETUP REQUIRED IF YOU WANT TO RECEIVE EVENTS\n"+
			"============================================================\n"+
			"1. Go to https://kick.com/settings/developer\n"+
			"2. Open your app settings\n"+
			"3. Enable 'Webhooks' toggle\n"+
			"4. Set Webhook URL to: "+webhookURL+"\n"+
			"5. Save changes\n"+
			"6. Restart eventsub if needed\n"+
			"============================================================",
		slog.String("webhook_url", webhookURL),
	)
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

		if ch.KickBotID == nil {
			m.logger.InfoContext(
				ctx,
				"webhook manager: channel has no kick bot assigned, skipping unsubscribe",
				slog.String("channel_id", ch.ID.String()),
				slog.String("kick_user_id", ch.KickUserID.String()),
			)
			continue
		}

		kickBot, err := m.kickBotsRepo.GetByID(ctx, *ch.KickBotID)
		if err != nil {
			if errors.Is(err, kickbotsrepository.ErrNotFound) {
				m.logger.InfoContext(
					ctx,
					"webhook manager: kick bot not found, skipping unsubscribe",
					slog.String("channel_id", ch.ID.String()),
					slog.String("kick_bot_id", ch.KickBotID.String()),
				)
				continue
			}

			m.logger.WarnContext(
				ctx,
				"webhook manager: failed to get kick bot for unsubscribe",
				slog.String("channel_id", ch.ID.String()),
				slog.String("kick_bot_id", ch.KickBotID.String()),
				logger.Error(err),
			)
			continue
		}

		accessToken, err := crypto.Decrypt(kickBot.AccessToken, m.config.TokensCipherKey)
		if err != nil {
			m.logger.WarnContext(
				ctx,
				"webhook manager: failed to decrypt kick token for unsubscribe",
				slog.String("channel_id", ch.ID.String()),
				slog.String("kick_bot_id", ch.KickBotID.String()),
				logger.Error(err),
			)
			continue
		}

		if err := m.kickSubMgr.UnsubscribeAll(ctx, ch.KickUserID.String(), accessToken, *ch.KickBotID, kickBot.RefreshToken); err != nil {
			m.logger.WarnContext(
				ctx,
				"webhook manager: failed to unsubscribe kick",
				slog.String("channel_id", ch.ID.String()),
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

	m.logger.InfoContext(
		ctx,
		"webhook manager: subscribing to kick events",
		slog.Int("channels_count", len(kickChannels)),
	)

	for _, ch := range kickChannels {
		if ch.KickUserID == nil || !ch.IsEnabled {
			continue
		}

		if ch.KickBotID == nil {
			m.logger.InfoContext(
				ctx,
				"webhook manager: channel has no kick bot assigned, skipping subscribe",
				slog.String("channel_id", ch.ID.String()),
				slog.String("kick_user_id", ch.KickUserID.String()),
			)
			continue
		}

		kickBot, err := m.kickBotsRepo.GetByID(ctx, *ch.KickBotID)
		if err != nil {
			if errors.Is(err, kickbotsrepository.ErrNotFound) {
				m.logger.InfoContext(
					ctx,
					"webhook manager: kick bot not found, skipping subscribe",
					slog.String("channel_id", ch.ID.String()),
					slog.String("kick_bot_id", ch.KickBotID.String()),
				)
			} else {
				m.logger.WarnContext(
					ctx,
					"webhook manager: failed to get kick bot for subscribe",
					slog.String("channel_id", ch.ID.String()),
					slog.String("kick_bot_id", ch.KickBotID.String()),
					logger.Error(err),
				)
			}
			continue
		}

		accessToken, err := crypto.Decrypt(kickBot.AccessToken, m.config.TokensCipherKey)
		if err != nil {
			m.logger.WarnContext(
				ctx,
				"webhook manager: failed to decrypt token for kick subscribe",
				slog.String("channel_id", ch.ID.String()),
				slog.String("kick_bot_id", ch.KickBotID.String()),
				logger.Error(err),
			)
			continue
		}

		if err := m.kickSubMgr.SubscribeAll(ctx, ch.KickUserID.String(), accessToken, *ch.KickBotID, kickBot.RefreshToken); err != nil {
			m.logger.ErrorContext(
				ctx,
				"webhook manager: failed to subscribe kick",
				slog.String("channel_id", ch.ID.String()),
				slog.String("kick_user_id", ch.KickUserID.String()),
				logger.Error(err),
			)
			continue
		}

		m.logger.InfoContext(
			ctx,
			"webhook manager: subscribed kick eventsub",
			slog.String("channel_id", ch.ID.String()),
			slog.String("kick_user_id", ch.KickUserID.String()),
		)
	}

	m.logger.InfoContext(
		ctx,
		"webhook manager: finished subscribing to kick events",
		slog.Int("channels_count", len(kickChannels)),
	)

	return nil
}
