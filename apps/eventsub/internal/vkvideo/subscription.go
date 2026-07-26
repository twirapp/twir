package vkvideo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
)

func (t *Transport) startBinding(
	binding channelplatformentity.ChannelPlatform,
	lease *Lease,
	owned *ownedConnection,
) error {
	startupCtx, cancel := context.WithTimeout(lease.Context(), leaseExpiry)
	defer cancel()

	channel, err := t.tokens.DiscoverChatChannel(startupCtx, binding.UserID)
	if err != nil {
		return fmt.Errorf("discover VK Video chat channel: %w", err)
	}

	connection, err := t.newConnection(RealtimeClientConfig{
		Channel:       channel,
		QueueCapacity: 128,
		BindingID:     binding.ID,
		Logger:        t.logger,
		Tokens: TokenCallbacks{
			Context: lease.Context(),
			Connection: func(ctx context.Context) (string, error) {
				return t.tokens.ConnectionToken(ctx, binding.UserID)
			},
			Subscription: func(ctx context.Context, channel string) (string, error) {
				return t.tokens.SubscriptionToken(ctx, binding.UserID, channel)
			},
		},
	})
	if err != nil {
		return fmt.Errorf("create VK Video realtime connection: %w", err)
	}
	owned.Set(connection)

	if err := connection.Connect(startupCtx); err != nil {
		return fmt.Errorf("connect VK Video realtime connection: %w", err)
	}
	if err := lease.Context().Err(); err != nil {
		return fmt.Errorf("VK Video chat lease ended before connection started: %w", err)
	}

	active := &activeBinding{lease: lease, connection: owned}
	t.mu.Lock()
	t.bindings[binding.ID] = active
	t.mu.Unlock()

	go t.consume(lease.Context(), binding, active)
	return nil
}

func (t *Transport) consume(
	ctx context.Context,
	binding channelplatformentity.ChannelPlatform,
	active *activeBinding,
) {
	defer t.removeBinding(binding.ID, active)

	for {
		publication, err := active.connection.Receive(ctx)
		if err != nil {
			if ctx.Err() == nil && t.logger != nil {
				t.logger.WarnContext(ctx, "VK Video chat receive failed", slog.String("binding_id", binding.ID.String()), slog.Any("error", err))
			}
			return
		}
		if err := t.handlePublication(ctx, binding, publication); err != nil && t.logger != nil {
			t.logger.WarnContext(ctx, "VK Video chat publication ignored", slog.String("binding_id", binding.ID.String()), slog.Any("error", err))
		}
	}
}

func (t *Transport) removeBinding(bindingID uuid.UUID, active *activeBinding) {
	t.mu.Lock()
	defer t.mu.Unlock()
	if t.bindings[bindingID] == active {
		delete(t.bindings, bindingID)
	}
}
