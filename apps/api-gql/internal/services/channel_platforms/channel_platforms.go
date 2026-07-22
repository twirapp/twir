package channel_platforms

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	authroutes "github.com/twirapp/twir/apps/api-gql/internal/delivery/http/routes/auth"
	appplatform "github.com/twirapp/twir/apps/api-gql/internal/platform"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
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
}

func NewFx(opts Opts) *Service {
	return New(
		opts.ChannelService,
		opts.UsersRepository,
		opts.ChannelPlatformsRepo,
		opts.Auth,
		opts.PlatformRegistry,
	)
}

func New(
	channels ChannelReader,
	users UserLookup,
	bindings BindingRepository,
	oauth OAuthStarter,
	registry *appplatform.Registry,
) *Service {
	return &Service{
		channels: channels,
		users:    users,
		bindings: bindings,
		oauth:    oauth,
		registry: registry,
	}
}

type Operations interface {
	List(context.Context, uuid.UUID) ([]Binding, error)
	Connect(context.Context, platformentity.Platform, string) (string, error)
	Disconnect(context.Context, uuid.UUID, platformentity.Platform) error
	SetEnabled(context.Context, uuid.UUID, platformentity.Platform, bool) (Binding, error)
}

type Service struct {
	channels ChannelReader
	users    UserLookup
	bindings BindingRepository
	oauth    OAuthStarter
	registry *appplatform.Registry
}

var _ Operations = (*Service)(nil)

type Binding struct {
	Binding      channelplatformsmodel.ChannelPlatform
	Profile      usersmodel.User
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
	Patch(context.Context, uuid.UUID, channelplatformsrepo.PatchInput) (channelplatformsmodel.ChannelPlatform, error)
	Delete(context.Context, uuid.UUID) error
}

type OAuthStarter interface {
	StartPlatformAuth(context.Context, platformentity.Platform, string) (string, error)
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

func (s *Service) Connect(ctx context.Context, platform platformentity.Platform, redirectTo string) (string, error) {
	if err := s.requireAvailable(platform); err != nil {
		return "", err
	}
	if s.oauth == nil {
		return "", fmt.Errorf("platform OAuth service is not configured")
	}

	return s.oauth.StartPlatformAuth(ctx, platform, redirectTo)
}

func (s *Service) Disconnect(ctx context.Context, channelID uuid.UUID, platform platformentity.Platform) error {
	if err := s.requireAvailable(platform); err != nil {
		return err
	}
	if s.channels == nil || s.bindings == nil {
		return fmt.Errorf("channel platform binding service is not configured")
	}

	channel, err := s.channels.GetChannelByID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("get channel: %w", err)
	}
	if len(channel.Bindings) <= 1 {
		return ErrLastBinding
	}

	binding, err := s.bindings.GetByChannelAndPlatform(ctx, channelID, platform)
	if err != nil {
		return fmt.Errorf("get channel platform binding: %w", err)
	}
	if err := s.bindings.Delete(ctx, binding.ID); err != nil {
		return fmt.Errorf("delete channel platform binding: %w", err)
	}

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
	if s.bindings == nil {
		return Binding{}, fmt.Errorf("channel platform binding repository is not configured")
	}

	binding, err := s.bindings.GetByChannelAndPlatform(ctx, channelID, platform)
	if err != nil {
		return Binding{}, fmt.Errorf("get channel platform binding: %w", err)
	}
	updated, err := s.bindings.Patch(ctx, binding.ID, channelplatformsrepo.PatchInput{Enabled: &enabled})
	if err != nil {
		return Binding{}, fmt.Errorf("set channel platform binding enabled state: %w", err)
	}

	return s.withProfile(ctx, updated)
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
