package streamlabs_integration

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	buscore "github.com/twirapp/twir/libs/bus-core"
	busintegrations "github.com/twirapp/twir/libs/bus-core/integrations"
	config "github.com/twirapp/twir/libs/config"
	streamlabsentity "github.com/twirapp/twir/libs/entities/streamlabs_integration"
	"github.com/twirapp/twir/libs/integrations/oauthlock"
	provider "github.com/twirapp/twir/libs/integrations/streamlabs"
	repository "github.com/twirapp/twir/libs/repositories/streamlabs_integration"
	"github.com/twirapp/twir/libs/repositories/streamlabs_integration/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	StreamlabsRepository repository.Repository
	TwirBus              *buscore.Bus
	Config               config.Config
	Redis                *redis.Client
}

func New(opts Opts) *Service {
	return &Service{
		streamlabsRepository: opts.StreamlabsRepository,
		config:               opts.Config,
		clientFactory:        streamlabsClientFactory{},
		locker:               oauthlock.NewRedis(opts.Redis),
		events:               busIntegrationEvents{bus: opts.TwirBus},
	}
}

type providerClient interface {
	GetAuthLink(state string) string
	ExchangeCode(context.Context, string) (*provider.TokenResponse, error)
	GetProfile(context.Context) (*provider.UserProfile, error)
}

type providerClientFactory interface {
	New(clientID, clientSecret, redirectURL string) providerClient
	NewAuthorized(
		clientID, clientSecret, channelID, redirectURL string,
		tokens provider.Tokens,
		store provider.TokenStore,
		locker oauthlock.Locker,
	) providerClient
}

type streamlabsClientFactory struct{}

func (streamlabsClientFactory) New(clientID, clientSecret, redirectURL string) providerClient {
	return provider.New(clientID, clientSecret, redirectURL)
}

func (streamlabsClientFactory) NewAuthorized(
	clientID, clientSecret, channelID, redirectURL string,
	tokens provider.Tokens,
	store provider.TokenStore,
	locker oauthlock.Locker,
) providerClient {
	return provider.NewAuthorized(
		clientID,
		clientSecret,
		channelID,
		redirectURL,
		tokens,
		store,
		locker,
	)
}

type integrationEvents interface {
	PublishAdd(context.Context, busintegrations.Request) error
	PublishRemove(context.Context, busintegrations.Request) error
}

type busIntegrationEvents struct {
	bus *buscore.Bus
}

func (b busIntegrationEvents) PublishAdd(
	ctx context.Context,
	request busintegrations.Request,
) error {
	return b.bus.Integrations.Add.Publish(ctx, request)
}

func (b busIntegrationEvents) PublishRemove(
	ctx context.Context,
	request busintegrations.Request,
) error {
	return b.bus.Integrations.Remove.Publish(ctx, request)
}

type Service struct {
	streamlabsRepository repository.Repository
	config               config.Config
	clientFactory        providerClientFactory
	locker               oauthlock.Locker
	events               integrationEvents
}

type AuthLinkResponse struct {
	Link string `json:"link"`
}

func (s *Service) mapModelToEntity(m model.StreamlabsIntegration) streamlabsentity.Entity {
	return streamlabsentity.Entity{
		ID:           m.ID,
		Enabled:      m.Enabled,
		ChannelID:    m.ChannelID,
		AccessToken:  m.AccessToken,
		RefreshToken: m.RefreshToken,
		UserName:     m.UserName,
		Avatar:       m.Avatar,
		CreatedAt:    m.CreatedAt,
		UpdatedAt:    m.UpdatedAt,
	}
}

func (s *Service) GetIntegrationData(ctx context.Context, channelID string) (
	streamlabsentity.Entity,
	error,
) {
	integration, err := s.streamlabsRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return streamlabsentity.Entity{ChannelID: channelID, Enabled: false}, nil
		}
		return streamlabsentity.Entity{}, fmt.Errorf("failed to get streamlabs integration: %w", err)
	}

	return s.mapModelToEntity(integration), nil
}

func (s *Service) getCallbackURL(ctx context.Context) (string, error) {
	baseURL, _ := gincontext.GetBaseUrlFromContext(ctx, s.config.SiteBaseUrl)
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", fmt.Errorf("invalid site base URL: %w", err)
	}

	return u.JoinPath("dashboard", "integrations", "streamlabs").String(), nil
}

func (s *Service) GetAuthLink(
	ctx context.Context,
	state string,
) (*AuthLinkResponse, error) {
	if err := s.ensureConfigured(); err != nil {
		return nil, err
	}
	if strings.TrimSpace(state) == "" {
		return nil, errors.New("streamlabs OAuth state is required")
	}

	redirectURL, err := s.getCallbackURL(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get redirect URL: %w", err)
	}
	client := s.clientFactory.New(
		s.config.StreamlabsClientId,
		s.config.StreamlabsClientSecret,
		redirectURL,
	)

	return &AuthLinkResponse{Link: client.GetAuthLink(state)}, nil
}

func (s *Service) PostCode(ctx context.Context, channelID, code string) error {
	if err := s.ensureConfigured(); err != nil {
		return err
	}

	foundIntegration, err := s.streamlabsRepository.GetByChannelID(ctx, channelID)
	if err != nil && !errors.Is(err, repository.ErrNotFound) {
		return fmt.Errorf("failed to get streamlabs integration: %w", err)
	}

	redirectURL, err := s.getCallbackURL(ctx)
	if err != nil {
		return fmt.Errorf("failed to get redirect URL: %w", err)
	}
	client := s.clientFactory.New(
		s.config.StreamlabsClientId,
		s.config.StreamlabsClientSecret,
		redirectURL,
	)
	tokens, err := client.ExchangeCode(ctx, code)
	if err != nil {
		return fmt.Errorf("exchange streamlabs authorization code: %w", err)
	}
	if tokens == nil || tokens.AccessToken == "" || tokens.RefreshToken == "" {
		return errors.New("exchange streamlabs authorization code: provider returned incomplete credentials")
	}

	profile, err := client.GetProfile(ctx)
	if err != nil {
		return fmt.Errorf("get streamlabs profile: %w", err)
	}
	if profile == nil {
		return errors.New("get streamlabs profile: provider returned no profile")
	}

	if foundIntegration == model.Nil {
		if err := s.streamlabsRepository.Create(ctx, repository.CreateOpts{
			ChannelID:    channelID,
			AccessToken:  tokens.AccessToken,
			RefreshToken: tokens.RefreshToken,
			Enabled:      true,
			UserName:     profile.StreamLabs.DisplayName,
			Avatar:       profile.StreamLabs.ThumbNail,
		}); err != nil {
			return fmt.Errorf("failed to create streamlabs integration: %w", err)
		}
	} else {
		if err := s.streamlabsRepository.Update(ctx, repository.UpdateOpts{
			ChannelID:    channelID,
			AccessToken:  &tokens.AccessToken,
			RefreshToken: &tokens.RefreshToken,
			Enabled:      lo.ToPtr(true),
			UserName:     &profile.StreamLabs.DisplayName,
			Avatar:       &profile.StreamLabs.ThumbNail,
		}); err != nil {
			return fmt.Errorf("failed to update streamlabs integration: %w", err)
		}
	}

	newIntegration, err := s.streamlabsRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to get streamlabs integration after update: %w", err)
	}
	if err := s.events.PublishAdd(ctx, busintegrations.Request{
		ID:      newIntegration.ID.String(),
		Service: busintegrations.Streamlabs,
	}); err != nil {
		return fmt.Errorf("failed to publish add integration event: %w", err)
	}

	return nil
}

func (s *Service) authorizedClient(ctx context.Context, channelID string) (providerClient, error) {
	if err := s.ensureConfigured(); err != nil {
		return nil, err
	}
	tokens, err := s.streamlabsRepository.GetTokens(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("get streamlabs tokens: %w", err)
	}
	redirectURL, err := s.getCallbackURL(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get redirect URL: %w", err)
	}
	return s.clientFactory.NewAuthorized(
		s.config.StreamlabsClientId,
		s.config.StreamlabsClientSecret,
		channelID,
		redirectURL,
		tokens,
		s.streamlabsRepository,
		s.locker,
	), nil
}

func (s *Service) Logout(ctx context.Context, channelID string) error {
	if s.locker == nil {
		return errors.New("streamlabs logout refresh lock is not configured")
	}
	err := s.locker.WithLock(
		ctx,
		provider.RefreshLockKey(channelID),
		func(lockCtx context.Context) error {
			if err := s.streamlabsRepository.Delete(lockCtx, channelID); err != nil {
				return fmt.Errorf("failed to disable streamlabs integration: %w", err)
			}
			return nil
		},
	)
	if err != nil {
		return fmt.Errorf("lock streamlabs logout: %w", err)
	}

	if err := s.events.PublishRemove(ctx, busintegrations.Request{
		ID:      channelID,
		Service: busintegrations.Streamlabs,
	}); err != nil {
		return fmt.Errorf("failed to publish remove integration event: %w", err)
	}

	return nil
}

func (s *Service) ensureConfigured() error {
	if s.config.StreamlabsClientId == "" || s.config.StreamlabsClientSecret == "" {
		return errors.New("streamlabs integration not properly configured")
	}
	return nil
}
