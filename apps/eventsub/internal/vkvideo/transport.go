package vkvideo

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/go-redsync/redsync/v4"
	redsyncgoredis "github.com/go-redsync/redsync/v4/redis/goredis/v9"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	user_creator "github.com/twirapp/twir/apps/eventsub/internal/services/user-creator"
	buscore "github.com/twirapp/twir/libs/bus-core"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/integrations/vk"
	"go.uber.org/fx"
)

const (
	leaseExpiry        = 30 * time.Second
	leaseRenewInterval = 10 * time.Second
)

type Opts struct {
	fx.In

	Logger               *slog.Logger
	Redis                *redis.Client
	Bus                  *buscore.Bus
	UserCreator          *user_creator.UserCreatorService
	WebSocketTokenClient *vk.WebSocketTokenClient
	Lc                   fx.Lifecycle
}

type Transport struct {
	logger        *slog.Logger
	ownership     *Ownership
	tokens        webSocketTokenProvider
	userCreator   chatUserEnsurer
	chatMessages  messagePublisher
	commands      messagePublisher
	deduplicator  messageDeduplicator
	newConnection realtimeConnectionFactory

	mu       sync.Mutex
	bindings map[uuid.UUID]*activeBinding
}

type transportDependencies struct {
	logger        *slog.Logger
	ownership     *Ownership
	tokens        webSocketTokenProvider
	userCreator   chatUserEnsurer
	chatMessages  messagePublisher
	commands      messagePublisher
	deduplicator  messageDeduplicator
	newConnection realtimeConnectionFactory
}

func New(opts Opts) (*Transport, error) {
	ownership, err := NewOwnership(
		redsync.New(redsyncgoredis.NewPool(opts.Redis)),
		LeaseConfig{Expiry: leaseExpiry, RenewInterval: leaseRenewInterval},
	)
	if err != nil {
		return nil, fmt.Errorf("create VK Video ownership: %w", err)
	}

	oauthTokens := busTokenProvider{request: opts.Bus.Tokens.RequestUserToken}
	transport := newTransport(transportDependencies{
		logger:        opts.Logger,
		ownership:     ownership,
		tokens:        devAPIWebSocketTokenProvider{oauthTokens: oauthTokens, client: opts.WebSocketTokenClient},
		userCreator:   opts.UserCreator,
		chatMessages:  opts.Bus.ChatMessages,
		commands:      opts.Bus.Parser.ProcessMessageAsCommand,
		deduplicator:  redisDeduplicator{redis: opts.Redis},
		newConnection: newCentrifugoConnection,
	})
	opts.Lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), leaseExpiry)
			defer cancel()
			return transport.Shutdown(shutdownCtx)
		},
	})

	return transport, nil
}

func newTransport(deps transportDependencies) *Transport {
	return &Transport{
		logger:        deps.logger,
		ownership:     deps.ownership,
		tokens:        deps.tokens,
		userCreator:   deps.userCreator,
		chatMessages:  deps.chatMessages,
		commands:      deps.commands,
		deduplicator:  deps.deduplicator,
		newConnection: deps.newConnection,
		bindings:      make(map[uuid.UUID]*activeBinding),
	}
}

func (*Transport) Platform() platformentity.Platform {
	return platformentity.PlatformVKVideoLive
}

func (*Transport) Capabilities() platformentity.Capabilities {
	return platformentity.PlatformVKVideoLive.Capabilities()
}

func (*Transport) SetCallbackBaseURL(string) {}

func (t *Transport) Subscribe(ctx context.Context, binding channelplatformentity.ChannelPlatform) error {
	if !binding.Enabled {
		return nil
	}

	t.mu.Lock()
	active, exists := t.bindings[binding.ID]
	t.mu.Unlock()
	if exists {
		if vkVideoBindingSnapshotsEqual(active.binding, binding) {
			return nil
		}
		if err := t.Unsubscribe(ctx, binding); err != nil {
			return fmt.Errorf("restart VK Video chat binding: %w", err)
		}
	}

	owned := &ownedConnection{}
	lease, err := t.ownership.Acquire(ctx, binding.ID.String(), owned.Close)
	if err != nil {
		if errors.Is(err, redsync.ErrFailed) {
			return nil
		}
		var errTaken *redsync.ErrTaken
		if errors.As(err, &errTaken) {
			return nil
		}
		return fmt.Errorf("acquire VK Video chat lease: %w", err)
	}

	if err := t.startBinding(binding, lease, owned); err != nil {
		releaseCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), leaseExpiry)
		defer cancel()
		_ = lease.Release(releaseCtx)
		return err
	}

	return nil
}

func vkVideoBindingSnapshotsEqual(a, b channelplatformentity.ChannelPlatform) bool {
	if a.Platform != b.Platform || a.UserID != b.UserID || a.PlatformChannelID != b.PlatformChannelID {
		return false
	}
	if a.BotUserID == nil || b.BotUserID == nil {
		return a.BotUserID == b.BotUserID
	}
	return *a.BotUserID == *b.BotUserID
}

func (t *Transport) Unsubscribe(ctx context.Context, binding channelplatformentity.ChannelPlatform) error {
	t.mu.Lock()
	active, exists := t.bindings[binding.ID]
	if exists {
		delete(t.bindings, binding.ID)
	}
	t.mu.Unlock()
	if !exists {
		return nil
	}

	if err := active.lease.Release(ctx); err != nil && !errors.Is(err, ErrLeaseLost) {
		return fmt.Errorf("release VK Video chat lease: %w", err)
	}

	return nil
}

func (t *Transport) Shutdown(ctx context.Context) error {
	t.mu.Lock()
	bindings := t.bindings
	t.bindings = make(map[uuid.UUID]*activeBinding)
	t.mu.Unlock()

	var errs []error
	for bindingID, active := range bindings {
		if err := active.lease.Release(ctx); err != nil && !errors.Is(err, ErrLeaseLost) {
			errs = append(errs, fmt.Errorf("release VK Video chat lease for binding %s: %w", bindingID, err))
		}
	}

	return errors.Join(errs...)
}
