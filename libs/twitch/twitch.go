package twitch

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/kvizyx/twitchy/helix"
	buscore "github.com/twirapp/twir/libs/bus-core"
	cfg "github.com/twirapp/twir/libs/config"
)

func NewAppClient(config cfg.Config, twirBus *buscore.Bus) (*helix.Client, error) {
	return NewAppClientWithContext(context.Background(), config, twirBus)
}

func NewAppClientWithContext(ctx context.Context, config cfg.Config, _ *buscore.Bus) (*helix.Client, error) {
	runtime, err := getRuntime(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("initialize Twitch app client: %w", err)
	}

	return runtime.appClient, nil
}

func NewUserClient(userID uuid.UUID, config cfg.Config, twirBus *buscore.Bus) (*helix.Client, error) {
	return NewUserClientWithContext(context.Background(), userID, config, twirBus)
}

func NewUserClientWithContext(ctx context.Context, userID uuid.UUID, config cfg.Config, _ *buscore.Bus) (*helix.Client, error) {
	runtime, err := getRuntime(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("initialize Twitch user client: %w", err)
	}
	if err := runtime.registerBroadcaster(ctx, userID.String()); err != nil {
		return nil, fmt.Errorf("register Twitch broadcaster: %w", err)
	}

	client, err := runtime.rootClient.AsUser(userID.String())
	if err != nil {
		return nil, fmt.Errorf("resolve Twitch broadcaster: %w", err)
	}
	return client, nil
}

func NewBotClient(botID string, config cfg.Config, twirBus *buscore.Bus) (*helix.Client, error) {
	return NewBotClientWithContext(context.Background(), botID, config, twirBus)
}

func NewBotClientWithContext(ctx context.Context, botID string, config cfg.Config, _ *buscore.Bus) (*helix.Client, error) {
	runtime, err := getRuntime(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("initialize Twitch bot client: %w", err)
	}
	if err := runtime.registerBot(ctx, botID); err != nil {
		return nil, fmt.Errorf("register Twitch bot: %w", err)
	}

	client, err := runtime.rootClient.AsUser(botID)
	if err != nil {
		return nil, fmt.Errorf("resolve Twitch bot: %w", err)
	}
	return client, nil
}

// NewChannelBotClient resolves a bot credential only for the supplied Twitch
// channel. Callers with a channel binding must use this instead of NewBotClient.
func NewChannelBotClient(botID string, channelID string, config cfg.Config) (*helix.Client, error) {
	return NewChannelBotClientWithContext(context.Background(), botID, channelID, config)
}

func NewChannelBotClientWithContext(ctx context.Context, botID string, channelID string, config cfg.Config) (*helix.Client, error) {
	runtime, err := getRuntime(ctx, config)
	if err != nil {
		return nil, fmt.Errorf("initialize Twitch channel bot client: %w", err)
	}
	if err := runtime.registerChannelBot(ctx, botID, channelID); err != nil {
		return nil, fmt.Errorf("register Twitch channel bot: %w", err)
	}

	client, err := runtime.rootClient.AsIntent(channelBotIntent(channelID))
	if err != nil {
		return nil, fmt.Errorf("resolve Twitch channel bot: %w", err)
	}
	return client, nil
}

func RegisterChannelBot(ctx context.Context, botID string, channelID string, config cfg.Config) error {
	runtime, err := getRuntime(ctx, config)
	if err != nil {
		return fmt.Errorf("initialize Twitch channel bot registry: %w", err)
	}
	if err := runtime.registerChannelBot(ctx, botID, channelID); err != nil {
		return fmt.Errorf("register Twitch channel bot: %w", err)
	}
	return nil
}

func RemoveChannelBot(ctx context.Context, botID string, channelID string, config cfg.Config) error {
	runtime, err := getRuntime(ctx, config)
	if err != nil {
		return fmt.Errorf("initialize Twitch channel bot registry: %w", err)
	}
	if err := runtime.removeChannelBot(ctx, botID, channelID); err != nil {
		return fmt.Errorf("remove Twitch channel bot: %w", err)
	}
	return nil
}

// Close releases the package runtime. It is safe to call repeatedly and makes
// the next constructor create a fresh runtime, which is useful for tests.
func Close() error {
	return closeRuntime()
}
