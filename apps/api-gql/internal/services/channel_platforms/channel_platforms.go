package channel_platforms

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
	authroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/auth"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/eventsub"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	channelplatformsrepo "github.com/twirapp/twir/libs/repositories/channel_platforms"
	channelplatformsmodel "github.com/twirapp/twir/libs/repositories/channel_platforms/model"
	channelsmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	usersrepo "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	channelservice "github.com/twirapp/twir/libs/services/channels"
	"go.uber.org/fx"
)

var (
	ErrPlatformUnavailable = errors.New("platform is not available")
	ErrLastBinding         = errors.New("cannot disconnect the last channel platform binding")
)

type Opts struct {
	fx.In

	ChannelService       *channelservice.ChannelService
	UsersRepository      usersrepo.Repository
	ChannelPlatformsRepo channelplatformsrepo.Repository
	Auth                 *authroutes.Auth
	PlatformRegistry     *appplatform.Registry
	TrmManager           trm.Manager
	TwirBus              *buscore.Bus
	Logger               *slog.Logger
}

func NewFx(opts Opts) *Service {
	service := New(
		opts.ChannelService,
		opts.UsersRepository,
		opts.ChannelPlatformsRepo,
		opts.Auth,
		opts.PlatformRegistry,
		opts.TrmManager,
	)
	service.eventSub = busEventSubPublisher{bus: opts.TwirBus}
	service.logger = opts.Logger
	return service
}

func New(
	channels ChannelReader,
	users UserLookup,
	bindings BindingRepository,
	oauth OAuthStarter,
	registry *appplatform.Registry,
	transactions transactionRunner,
) *Service {
	return &Service{
		channels:     channels,
		users:        users,
		bindings:     bindings,
		oauth:        oauth,
		registry:     registry,
		transactions: transactions,
	}
}

type Operations interface {
	List(context.Context, uuid.UUID) ([]Binding, error)
	Options() []Option
	Connect(context.Context, uuid.UUID, platformentity.Platform) (string, error)
	Disconnect(context.Context, uuid.UUID, platformentity.Platform) error
	SetEnabled(context.Context, uuid.UUID, platformentity.Platform, bool) (Binding, error)
}

type Service struct {
	channels     ChannelReader
	users        UserLookup
	bindings     BindingRepository
	oauth        OAuthStarter
	registry     *appplatform.Registry
	transactions transactionRunner
	eventSub     eventSubPublisher
	logger       *slog.Logger
}

var _ Operations = (*Service)(nil)

type Binding struct {
	Binding      channelplatformsmodel.ChannelPlatform
	Profile      usersmodel.User
	Capabilities platformentity.Capabilities
}

type Option struct {
	Platform     platformentity.Platform
	Capabilities platformentity.Capabilities
}

type ChannelReader interface {
	GetChannelByID(context.Context, uuid.UUID) (channelsmodel.Channel, error)
}

type UserLookup interface {
	GetByID(context.Context, uuid.UUID) (usersmodel.User, error)
}

type BindingRepository interface {
	GetByChannelAndPlatform(context.Context, uuid.UUID, platformentity.Platform) (channelplatformsmodel.ChannelPlatform, error)
	LockByChannelID(context.Context, uuid.UUID) error
	Patch(context.Context, uuid.UUID, channelplatformsrepo.PatchInput) (channelplatformsmodel.ChannelPlatform, error)
	Delete(context.Context, uuid.UUID) error
}

type transactionRunner interface {
	Do(context.Context, func(context.Context) error) error
}

type eventSubPublisher interface {
	Subscribe(context.Context, eventsub.EventsubSubscribeToAllEventsRequest) error
	Unsubscribe(context.Context, eventsub.EventsubUnsubscribeRequest) error
}

type busEventSubPublisher struct {
	bus *buscore.Bus
}

func (p busEventSubPublisher) Subscribe(ctx context.Context, request eventsub.EventsubSubscribeToAllEventsRequest) error {
	return p.bus.EventSub.SubscribeToAllEvents.Publish(ctx, request)
}

func (p busEventSubPublisher) Unsubscribe(ctx context.Context, request eventsub.EventsubUnsubscribeRequest) error {
	return p.bus.EventSub.Unsubscribe.Publish(ctx, request)
}

type OAuthStarter interface {
	StartPlatformAuthForChannel(context.Context, uuid.UUID, platformentity.Platform) (string, error)
}

func (s *Service) List(ctx context.Context, channelID uuid.UUID) ([]Binding, error) {
	if s.channels == nil {
		return nil, fmt.Errorf("channel service is not configured")
	}

	channel, err := s.channels.GetChannelByID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("get channel: %w", err)
	}

	result := make([]Binding, 0, len(channel.Bindings))
	for _, binding := range channel.Bindings {
		if !s.isAvailable(binding.Platform) {
			continue
		}

		mapped, err := s.withProfile(ctx, binding)
		if err != nil {
			return nil, err
		}
		result = append(result, mapped)
	}

	return result, nil
}

func (s *Service) Options() []Option {
	platforms := platformentity.All()
	options := make([]Option, 0, len(platforms))
	for _, platform := range platforms {
		if s.isAvailable(platform) {
			options = append(options, Option{
				Platform:     platform,
				Capabilities: platform.Capabilities(),
			})
		}
	}

	return options
}

func (s *Service) Connect(ctx context.Context, channelID uuid.UUID, platform platformentity.Platform) (string, error) {
	if err := s.requireAvailable(platform); err != nil {
		return "", err
	}
	if s.oauth == nil {
		return "", fmt.Errorf("platform OAuth service is not configured")
	}

	return s.oauth.StartPlatformAuthForChannel(ctx, channelID, platform)
}

func (s *Service) Disconnect(ctx context.Context, channelID uuid.UUID, platform platformentity.Platform) error {
	if err := s.requireAvailable(platform); err != nil {
		return err
	}
	if s.channels == nil || s.bindings == nil || s.transactions == nil {
		return fmt.Errorf("channel platform binding service is not configured")
	}

	var binding channelplatformsmodel.ChannelPlatform
	err := s.transactions.Do(ctx, func(txCtx context.Context) error {
		if err := s.bindings.LockByChannelID(txCtx, channelID); err != nil {
			return fmt.Errorf("lock channel platform bindings: %w", err)
		}

		channel, err := s.channels.GetChannelByID(txCtx, channelID)
		if err != nil {
			return fmt.Errorf("get channel: %w", err)
		}
		availableBindings := 0
		for _, binding := range channel.Bindings {
			if s.isAvailable(binding.Platform) {
				availableBindings++
			}
		}
		if availableBindings <= 1 {
			return ErrLastBinding
		}

		binding, err = s.bindings.GetByChannelAndPlatform(txCtx, channelID, platform)
		if err != nil {
			return fmt.Errorf("get channel platform binding: %w", err)
		}
		if err := s.bindings.Delete(txCtx, binding.ID); err != nil {
			return fmt.Errorf("delete channel platform binding: %w", err)
		}

		return nil
	})
	if err != nil {
		return err
	}

	s.publishBindingUnsubscribe(ctx, channelID, binding)
	return nil
}

func (s *Service) SetEnabled(
	ctx context.Context,
	channelID uuid.UUID,
	platform platformentity.Platform,
	enabled bool,
) (Binding, error) {
	if err := s.requireAvailable(platform); err != nil {
		return Binding{}, err
	}
	if s.bindings == nil || s.transactions == nil {
		return Binding{}, fmt.Errorf("channel platform binding service is not configured")
	}

	var updated channelplatformsmodel.ChannelPlatform
	err := s.transactions.Do(ctx, func(txCtx context.Context) error {
		binding, err := s.bindings.GetByChannelAndPlatform(txCtx, channelID, platform)
		if err != nil {
			return fmt.Errorf("get channel platform binding: %w", err)
		}
		updated, err = s.bindings.Patch(txCtx, binding.ID, channelplatformsrepo.PatchInput{Enabled: &enabled})
		if err != nil {
			return fmt.Errorf("set channel platform binding enabled state: %w", err)
		}
		return nil
	})
	if err != nil {
		return Binding{}, err
	}

	s.publishBindingLifecycle(ctx, channelID, updated)
	return s.withProfile(ctx, updated)
}

func (s *Service) publishBindingLifecycle(
	ctx context.Context,
	channelID uuid.UUID,
	binding channelplatformsmodel.ChannelPlatform,
) {
	if binding.Enabled {
		s.publishBindingSubscribe(ctx, channelID, binding.Platform)
		return
	}

	s.publishBindingUnsubscribe(ctx, channelID, binding)
}

func (s *Service) publishBindingSubscribe(ctx context.Context, channelID uuid.UUID, platform platformentity.Platform) {
	if s.eventSub == nil {
		return
	}
	if err := s.eventSub.Subscribe(ctx, eventsub.EventsubSubscribeToAllEventsRequest{
		ChannelID: channelID.String(),
		Platform:  platform,
	}); err != nil && s.logger != nil {
		s.logger.ErrorContext(ctx, "cannot publish eventsub subscribe", logger.Error(err), slog.String("channel_id", channelID.String()), slog.String("platform", platform.String()))
	}
}

func (s *Service) publishBindingUnsubscribe(ctx context.Context, channelID uuid.UUID, binding channelplatformsmodel.ChannelPlatform) {
	if s.eventSub == nil {
		return
	}
	if err := s.eventSub.Unsubscribe(ctx, eventsub.EventsubUnsubscribeRequest{
		ChannelID: channelID.String(),
		Platform:  binding.Platform,
		Binding: &eventsub.EventsubBindingSnapshot{
			ID:                binding.ID.String(),
			UserID:            binding.UserID.String(),
			PlatformChannelID: binding.PlatformChannelID,
		},
	}); err != nil && s.logger != nil {
		s.logger.ErrorContext(ctx, "cannot publish eventsub unsubscribe", logger.Error(err), slog.String("channel_id", channelID.String()), slog.String("platform", binding.Platform.String()))
	}
}

func (s *Service) withProfile(ctx context.Context, binding channelplatformsmodel.ChannelPlatform) (Binding, error) {
	if s.users == nil {
		return Binding{}, fmt.Errorf("users repository is not configured")
	}

	profile, err := s.users.GetByID(ctx, binding.UserID)
	if err != nil {
		return Binding{}, fmt.Errorf("get channel platform profile: %w", err)
	}

	return Binding{
		Binding:      binding,
		Profile:      profile,
		Capabilities: binding.Platform.Capabilities(),
	}, nil
}

func (s *Service) requireAvailable(platform platformentity.Platform) error {
	if !s.isAvailable(platform) {
		return fmt.Errorf("%w: %s", ErrPlatformUnavailable, platform)
	}

	return nil
}

func (s *Service) isAvailable(platform platformentity.Platform) bool {
	if s.registry == nil {
		return false
	}

	provider, ok := s.registry.Get(platform)
	return ok && provider != nil
}
