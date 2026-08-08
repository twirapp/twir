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
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	buscore "github.com/twirapp/twir/libs/bus-core"
	channelplatformentity "github.com/twirapp/twir/libs/entities/channel_platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/integrations/vk"
	"github.com/twirapp/twir/libs/repositories/channels"
)

const (
	leaseExpiry        = 30 * time.Second
	leaseRenewInterval = 10 * time.Second
)

type Opts struct {
	Logger               *slog.Logger
	Redis                *redis.Client
	Bus                  *buscore.Bus
	UserCreator          *user_creator.UserCreatorService
	WebSocketTokenClient *vk.WebSocketTokenClient
	ChannelsRepo         channels.Repository
	Lc                   *lifecycle.Lifecycle
	ProxyUrl             string
}

type Transport struct {
	logger           *slog.Logger
	ownership        *Ownership
	tokens           webSocketTokenProvider
	userCreator      chatUserEnsurer
	chatMessages     messagePublisher
	commands         messagePublisher
	deduplicator     messageDeduplicator
	newConnection    realtimeConnectionFactory
	databaseBindings bindingsProvider
	proxyUrl         string

	mu               sync.Mutex
	bindings         map[uuid.UUID]*activeBinding
	contentionLogged map[uuid.UUID]bool
}

type transportDependencies struct {
	logger           *slog.Logger
	ownership        *Ownership
	tokens           webSocketTokenProvider
	userCreator      chatUserEnsurer
	chatMessages     messagePublisher
	commands         messagePublisher
	deduplicator     messageDeduplicator
	newConnection    realtimeConnectionFactory
	databaseBindings bindingsProvider
	proxyUrl         string
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
		logger:           opts.Logger,
		ownership:        ownership,
		tokens:           devAPIWebSocketTokenProvider{oauthTokens: oauthTokens, client: opts.WebSocketTokenClient},
		userCreator:      opts.UserCreator,
		chatMessages:     opts.Bus.ChatMessages,
		commands:         opts.Bus.Parser.ProcessMessageAsCommand,
		deduplicator:     redisDeduplicator{redis: opts.Redis},
		newConnection:    newCentrifugoConnection,
		databaseBindings: newDatabaseBindingsProvider(opts.ChannelsRepo),
		proxyUrl:         opts.ProxyUrl,
	})
	reconcileCtx, stopReconcile := context.WithCancel(context.Background())
	go transport.reconcileLoop(reconcileCtx)
	opts.Lc.Append(lifecycle.Hook{
		OnStop: func(ctx context.Context) error {
			stopReconcile()
			shutdownCtx, cancel := context.WithTimeout(context.WithoutCancel(ctx), leaseExpiry)
			defer cancel()
			return transport.Shutdown(shutdownCtx)
		},
	})

	return transport, nil
}

func newTransport(deps transportDependencies) *Transport {
	return &Transport{
		logger:           deps.logger,
		ownership:        deps.ownership,
		tokens:           deps.tokens,
		userCreator:      deps.userCreator,
		chatMessages:     deps.chatMessages,
		commands:         deps.commands,
		deduplicator:     deps.deduplicator,
		newConnection:    deps.newConnection,
		databaseBindings: deps.databaseBindings,
		proxyUrl:         deps.proxyUrl,
		bindings:         make(map[uuid.UUID]*activeBinding),
		contentionLogged: make(map[uuid.UUID]bool),
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
		if isLeaseContended(err) {
			t.logLeaseContended(ctx, binding.ID)
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

	t.mu.Lock()
	delete(t.contentionLogged, binding.ID)
	t.mu.Unlock()

	return nil
}

func isLeaseContended(err error) bool {
	if errors.Is(err, redsync.ErrFailed) {
		return true
	}
	var errTaken *redsync.ErrTaken
	return errors.As(err, &errTaken)
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
	delete(t.contentionLogged, binding.ID)
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
	t.contentionLogged = make(map[uuid.UUID]bool)
	t.mu.Unlock()

	var errs []error
	for bindingID, active := range bindings {
		if err := active.lease.Release(ctx); err != nil && !errors.Is(err, ErrLeaseLost) {
			errs = append(errs, fmt.Errorf("release VK Video chat lease for binding %s: %w", bindingID, err))
		}
	}

	return errors.Join(errs...)
}

// reconcileLoop periodically syncs active bindings with the database state,
// e.g. when the chat lease was still held by a previous (already stopped)
// eventsub instance at startup or when an unsubscribe event was handled
// by another replica.
func (t *Transport) reconcileLoop(ctx context.Context) {
	ticker := time.NewTicker(leaseRenewInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			t.reconcileWithDatabase(ctx)
		}
	}
}

func (t *Transport) reconcileWithDatabase(ctx context.Context) {
	if t.databaseBindings == nil {
		return
	}

	databaseBindings, err := t.databaseBindings(ctx)
	if err != nil {
		if t.logger != nil {
			t.logger.WarnContext(ctx, "VK Video chat bindings reconcile query failed", slog.Any("error", err))
		}
		return
	}

	keep := make(map[uuid.UUID]struct{}, len(databaseBindings))
	for _, binding := range databaseBindings {
		keep[binding.ID] = struct{}{}
	}

	t.mu.Lock()
	stale := make([]channelplatformentity.ChannelPlatform, 0, len(t.bindings))
	for id, active := range t.bindings {
		if _, ok := keep[id]; !ok {
			stale = append(stale, active.binding)
		}
	}
	t.mu.Unlock()

	for _, binding := range stale {
		if err := t.Unsubscribe(ctx, binding); err != nil && t.logger != nil {
			t.logger.WarnContext(
				ctx,
				"VK Video chat forget removed binding failed",
				slog.String("binding_id", binding.ID.String()),
				slog.Any("error", err),
			)
		}
	}

	for _, binding := range databaseBindings {
		if err := t.Subscribe(ctx, binding); err != nil && t.logger != nil {
			t.logger.WarnContext(
				ctx,
				"VK Video chat binding reconcile failed",
				slog.String("binding_id", binding.ID.String()),
				slog.Any("error", err),
			)
		}
	}
}

func (t *Transport) logLeaseContended(ctx context.Context, bindingID uuid.UUID) {
	if t.logger == nil {
		return
	}

	t.mu.Lock()
	alreadyLogged := t.contentionLogged[bindingID]
	t.contentionLogged[bindingID] = true
	t.mu.Unlock()

	message := "VK Video chat binding is owned by another instance, will retry in background"
	if alreadyLogged {
		t.logger.DebugContext(ctx, message, slog.String("binding_id", bindingID.String()))
		return
	}
	t.logger.InfoContext(ctx, message, slog.String("binding_id", bindingID.String()))
}
