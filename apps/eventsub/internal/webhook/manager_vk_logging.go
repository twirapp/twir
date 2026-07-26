package webhook

import (
	"context"
	"log/slog"

	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

func (m *Manager) isTransportRegistered(platform platformentity.Platform) bool {
	if m.transports == nil {
		return false
	}

	_, ok := m.transports.Get(platform)
	return ok
}

func countEnabledBindings(bindings []channelplatformentity.ChannelPlatform) int {
	enabled := 0
	for _, binding := range bindings {
		if binding.Enabled {
			enabled++
		}
	}

	return enabled
}

func (m *Manager) logVKStartup(ctx context.Context) {
	m.logger.InfoContext(
		ctx,
		"webhook manager: VK EventSub transport registration state",
		slog.Bool("vk_registered", m.isTransportRegistered(platformentity.PlatformVKVideoLive)),
	)
}

func (m *Manager) logVKSubscribeSummary(ctx context.Context) error {
	bindings, err := m.bindingsForPlatform(ctx, platformentity.PlatformVKVideoLive)
	if err != nil {
		return err
	}

	registered := m.isTransportRegistered(platformentity.PlatformVKVideoLive)
	eligible := len(bindings)
	routed := 0
	skipReason := ""

	if registered {
		routed = countEnabledBindings(bindings)
		if routed < eligible {
			skipReason = "disabled_binding"
		}
	} else if eligible > 0 {
		skipReason = "transport_not_registered"
	}

	attrs := []any{
		slog.String("platform", platformentity.PlatformVKVideoLive.String()),
		slog.Bool("transport_registered", registered),
		slog.Int("eligible_bindings", eligible),
		slog.Int("routed_bindings", routed),
		slog.Int("skipped_bindings", eligible-routed),
	}

	if skipReason != "" {
		m.logger.WarnContext(ctx, "webhook manager: VK EventSub subscribe summary", append(attrs, slog.String("skip_reason", skipReason))...)
		return nil
	}

	m.logger.InfoContext(ctx, "webhook manager: VK EventSub subscribe summary", attrs...)

	return nil
}

func (m *Manager) logVKUnsubscribeSummary(ctx context.Context) error {
	bindings, err := m.bindingsForPlatform(ctx, platformentity.PlatformVKVideoLive)
	if err != nil {
		return err
	}

	registered := m.isTransportRegistered(platformentity.PlatformVKVideoLive)
	eligible := len(bindings)
	routed := eligible
	skipReason := ""

	if !registered && eligible > 0 {
		routed = 0
		skipReason = "transport_not_registered"
	}

	attrs := []any{
		slog.String("platform", platformentity.PlatformVKVideoLive.String()),
		slog.Bool("transport_registered", registered),
		slog.Int("eligible_bindings", eligible),
		slog.Int("routed_bindings", routed),
		slog.Int("skipped_bindings", eligible-routed),
	}

	if skipReason != "" {
		m.logger.WarnContext(ctx, "webhook manager: VK EventSub unsubscribe summary", append(attrs, slog.String("skip_reason", skipReason))...)
		return nil
	}

	m.logger.InfoContext(ctx, "webhook manager: VK EventSub unsubscribe summary", attrs...)

	return nil
}
