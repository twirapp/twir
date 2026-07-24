package webhook

import (
	"context"
	"fmt"
	"log/slog"

	httpserver "github.com/twirapp/twir/apps/eventsub/internal/http"
	eventplatforms "github.com/twirapp/twir/apps/eventsub/internal/platforms"
	cfg "github.com/twirapp/twir/libs/config"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	platformsregistry "github.com/twirapp/twir/libs/platforms"
	"github.com/twirapp/twir/libs/repositories/channels"
	"go.uber.org/fx"
)

type Manager struct {
	config       cfg.Config
	logger       *slog.Logger
	channelsRepo channels.Repository
	transports   *platformsregistry.Registry[eventplatforms.EventTransport]
}

type Opts struct {
	fx.In

	Lc fx.Lifecycle

	Config       cfg.Config
	Logger       *slog.Logger
	ChannelsRepo channels.Repository
	Transports   *platformsregistry.Registry[eventplatforms.EventTransport]
	Server       *httpserver.Server
}

func NewManager(opts Opts) *Manager {
	m := &Manager{
		config:       opts.Config,
		logger:       opts.Logger,
		channelsRepo: opts.ChannelsRepo,
		transports:   opts.Transports,
	}

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

func (m *Manager) registeredTransports() []eventplatforms.EventTransport {
	if m.transports == nil {
		return nil
	}

	transports := make([]eventplatforms.EventTransport, 0, len(platformentity.All()))
	for _, platform := range platformentity.All() {
		transport, ok := m.transports.Get(platform)
		if ok {
			transports = append(transports, transport)
		}
	}

	return transports
}

func (m *Manager) bindingsForPlatform(
	ctx context.Context,
	platform platformentity.Platform,
) ([]channelplatformentity.ChannelPlatform, error) {
	channels, err := m.channelsRepo.GetAllByBindingPlatform(ctx, platform)
	if err != nil {
		return nil, fmt.Errorf("list %s channels: %w", platform, err)
	}

	bindings := make([]channelplatformentity.ChannelPlatform, 0, len(channels))
	for _, channel := range channels {
		binding, ok := channel.Binding(platform)
		if !ok {
			continue
		}

		bindings = append(bindings, binding)
	}

	return bindings, nil
}

func (m *Manager) start(ctx context.Context) error {
	callbackBaseURL := m.config.SiteBaseUrl
	if m.config.EventSubCallbackBaseUrl != "" {
		callbackBaseURL = m.config.EventSubCallbackBaseUrl
	}

	transports := m.registeredTransports()
	m.logger.InfoContext(
		ctx,
		"webhook manager: starting",
		slog.String("callback_base_url", callbackBaseURL),
		slog.Bool("is_development", m.config.IsDevelopment()),
		slog.Int("platforms_count", len(transports)),
	)

	for _, transport := range transports {
		transport.SetCallbackBaseURL(callbackBaseURL)
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
	for _, transport := range m.registeredTransports() {
		bindings, err := m.bindingsForPlatform(ctx, transport.Platform())
		if err != nil {
			return err
		}

		if err := eventplatforms.UnsubscribeAll(ctx, m.transports, bindings); err != nil {
			m.logger.WarnContext(
				ctx,
				"webhook manager: failed to unsubscribe transport bindings",
				slog.String("platform", transport.Platform().String()),
				logger.Error(err),
			)
		}
	}

	return nil
}

func (m *Manager) subscribeAllPlatforms(ctx context.Context) error {
	for _, transport := range m.registeredTransports() {
		bindings, err := m.bindingsForPlatform(ctx, transport.Platform())
		if err != nil {
			return err
		}

		m.logger.InfoContext(
			ctx,
			"webhook manager: subscribing to platform events",
			slog.String("platform", transport.Platform().String()),
			slog.Int("bindings_count", len(bindings)),
		)

		if err := eventplatforms.SubscribeAll(ctx, m.transports, bindings); err != nil {
			m.logger.ErrorContext(
				ctx,
				"webhook manager: failed to subscribe transport bindings",
				slog.String("platform", transport.Platform().String()),
				logger.Error(err),
			)
		}

		m.logger.InfoContext(
			ctx,
			"webhook manager: finished subscribing to platform events",
			slog.String("platform", transport.Platform().String()),
			slog.Int("bindings_count", len(bindings)),
		)
	}

	return nil
}
