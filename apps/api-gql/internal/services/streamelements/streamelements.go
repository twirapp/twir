package streamelements

import (
	"context"
	"fmt"
	"net/url"

	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/services/importer"
	buscore "github.com/twirapp/twir/libs/bus-core"
	busintegrations "github.com/twirapp/twir/libs/bus-core/integrations"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/integrations/oauthlock"
	streamelementsintegration "github.com/twirapp/twir/libs/integrations/streamelements"
	channelsintegrations "github.com/twirapp/twir/libs/repositories/channels_integrations"
	channelsintegrationsmodel "github.com/twirapp/twir/libs/repositories/channels_integrations/model"
	"github.com/twirapp/twir/libs/repositories/integrations"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Config                  config.Config
	TwirBus                 *buscore.Bus
	Importer                *importer.Service
	IntegrationsRepository  integrations.Repository
	ChannelIntegrationsRepo channelsintegrations.Repository
	Redis                   *redis.Client
}

func New(opts Opts) (*Service, error) {
	siteBaseURL, err := url.Parse(opts.Config.SiteBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("parse site base URL: %w", err)
	}

	return &Service{
		config:                  opts.Config,
		redirectURL:             siteBaseURL.JoinPath("/dashboard/integrations/callbacks/streamelements").String(),
		importer:                opts.Importer,
		integrationsRepo:        opts.IntegrationsRepository,
		channelIntegrationsRepo: opts.ChannelIntegrationsRepo,
		clientFactory:           streamElementsClientFactory{},
		locker:                  oauthlock.NewRedis(opts.Redis),
		events:                  busIntegrationEvents{bus: opts.TwirBus},
	}, nil
}

type providerImporter interface {
	ImportCommands(context.Context, string, string, []importer.Command) (importer.Report, error)
	ImportTimers(context.Context, string, string, []importer.Timer) (importer.Report, error)
}

type providerClient interface {
	GetAuthLinkWithState(redirectURL, state string) (string, error)
	ExchangeCode(context.Context, string, string) (*streamelementsintegration.TokenResponse, error)
	GetProfile(context.Context) (*streamelementsintegration.UserProfile, error)
	GetCommands(context.Context, string) ([]streamelementsintegration.Command, error)
	GetTimers(context.Context, string) ([]streamelementsintegration.Timer, error)
}

type providerClientFactory interface {
	NewStatic(clientID, clientSecret string) providerClient
	NewAuthorized(
		clientID, clientSecret, channelID, redirectURL string,
		tokens streamelementsintegration.Tokens,
		store streamelementsintegration.TokenStore,
		locker oauthlock.Locker,
	) providerClient
}

type streamElementsClientFactory struct{}

func (streamElementsClientFactory) NewStatic(clientID, clientSecret string) providerClient {
	return streamelementsintegration.NewStatic(clientID, clientSecret)
}

func (streamElementsClientFactory) NewAuthorized(
	clientID, clientSecret, channelID, redirectURL string,
	tokens streamelementsintegration.Tokens,
	store streamelementsintegration.TokenStore,
	locker oauthlock.Locker,
) providerClient {
	return streamelementsintegration.NewAuthorized(
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

func (b busIntegrationEvents) PublishAdd(ctx context.Context, request busintegrations.Request) error {
	return b.bus.Integrations.Add.Publish(ctx, request)
}

func (b busIntegrationEvents) PublishRemove(ctx context.Context, request busintegrations.Request) error {
	return b.bus.Integrations.Remove.Publish(ctx, request)
}

type Service struct {
	config                  config.Config
	redirectURL             string
	importer                providerImporter
	integrationsRepo        integrations.Repository
	channelIntegrationsRepo channelsintegrations.Repository
	clientFactory           providerClientFactory
	locker                  oauthlock.Locker
	events                  integrationEvents
}

var _ streamelementsintegration.TokenStore = (*Service)(nil)

func (s *Service) GetAuthLink(ctx context.Context, state string) (string, error) {
	if err := s.ensureConfigured(); err != nil {
		return "", err
	}

	client := s.clientFactory.NewStatic(
		s.config.StreamElementsClientId,
		s.config.StreamElementsClientSecret,
	)
	link, err := client.GetAuthLinkWithState(s.redirectURL, state)
	if err != nil {
		return "", fmt.Errorf("create StreamElements authorization URL: %w", err)
	}
	return link, nil
}

func (s *Service) GetData(ctx context.Context, channelID string) (*IntegrationData, error) {
	channelIntegration, err := s.channelIntegrationsRepo.GetByChannelAndService(
		ctx,
		channelID,
		integrationsmodel.ServiceStreamElements,
	)
	if err != nil {
		return nil, fmt.Errorf("get StreamElements integration: %w", err)
	}
	if channelIntegration.IsNil() || channelIntegration.ID == "" || channelIntegration.Data == nil {
		return nil, nil
	}

	data := &IntegrationData{}
	if channelIntegration.Data.UserName != nil {
		data.UserName = *channelIntegration.Data.UserName
	}
	if channelIntegration.Data.Avatar != nil {
		data.Avatar = *channelIntegration.Data.Avatar
	}
	return data, nil
}

func (s *Service) PostCode(ctx context.Context, channelID, code string) error {
	if err := s.ensureConfigured(); err != nil {
		return err
	}

	integration, err := s.integrationsRepo.GetByService(ctx, integrationsmodel.ServiceStreamElements)
	if err != nil {
		return fmt.Errorf("get StreamElements configuration: %w", err)
	}
	if integration.ID == "" {
		return fmt.Errorf("StreamElements integration is not configured")
	}

	client := s.clientFactory.NewStatic(
		s.config.StreamElementsClientId,
		s.config.StreamElementsClientSecret,
	)
	tokens, err := client.ExchangeCode(ctx, code, s.redirectURL)
	if err != nil {
		return fmt.Errorf("exchange StreamElements authorization code: %w", err)
	}
	if tokens == nil || tokens.AccessToken == "" || tokens.RefreshToken == "" {
		return fmt.Errorf("exchange StreamElements authorization code: provider returned incomplete credentials")
	}

	profile, err := client.GetProfile(ctx)
	if err != nil {
		return fmt.Errorf("get StreamElements profile: %w", err)
	}
	if profile == nil {
		return fmt.Errorf("get StreamElements profile: provider returned no profile")
	}

	channelIntegration, err := s.channelIntegrationsRepo.GetByChannelAndService(
		ctx,
		channelID,
		integrationsmodel.ServiceStreamElements,
	)
	if err != nil {
		return fmt.Errorf("get existing StreamElements integration: %w", err)
	}

	userName := profile.DisplayName
	if userName == "" {
		userName = profile.Username
	}
	avatar := profile.Avatar
	if avatar == "" {
		avatar = profile.ProfileImage
	}
	data := &channelsintegrationsmodel.Data{UserName: &userName, Avatar: &avatar}

	integrationID := channelIntegration.ID
	if channelIntegration.IsNil() || channelIntegration.ID == "" {
		created, err := s.channelIntegrationsRepo.Create(ctx, channelsintegrations.CreateInput{
			ChannelID:     channelID,
			IntegrationID: integration.ID,
			Enabled:       true,
			AccessToken:   &tokens.AccessToken,
			RefreshToken:  &tokens.RefreshToken,
			Data:          data,
		})
		if err != nil {
			return fmt.Errorf("create StreamElements integration: %w", err)
		}
		integrationID = created.ID
	} else {
		if err := s.channelIntegrationsRepo.Update(ctx, channelIntegration.ID, channelsintegrations.UpdateInput{
			Enabled:      lo.ToPtr(true),
			AccessToken:  &tokens.AccessToken,
			RefreshToken: &tokens.RefreshToken,
			Data:         data,
		}); err != nil {
			return fmt.Errorf("update StreamElements integration: %w", err)
		}
	}

	if integrationID == "" {
		return fmt.Errorf("persist StreamElements integration: repository returned an empty ID")
	}
	if err := s.events.PublishAdd(ctx, busintegrations.Request{
		ID:      integrationID,
		Service: busintegrations.StreamElements,
	}); err != nil {
		return fmt.Errorf("publish StreamElements add event: %w", err)
	}

	return nil
}

func (s *Service) Logout(ctx context.Context, channelID string) error {
	channelIntegration, err := s.channelIntegrationsRepo.GetByChannelAndService(
		ctx,
		channelID,
		integrationsmodel.ServiceStreamElements,
	)
	if err != nil {
		return fmt.Errorf("get StreamElements integration: %w", err)
	}
	if channelIntegration.IsNil() || channelIntegration.ID == "" {
		return nil
	}

	if err := s.channelIntegrationsRepo.Update(ctx, channelIntegration.ID, channelsintegrations.UpdateInput{
		Enabled:           lo.ToPtr(false),
		ClearAccessToken:  true,
		ClearRefreshToken: true,
		Data:              &channelsintegrationsmodel.Data{},
	}); err != nil {
		return fmt.Errorf("disable StreamElements integration: %w", err)
	}

	if err := s.events.PublishRemove(ctx, busintegrations.Request{
		ID:      channelID,
		Service: busintegrations.StreamElements,
	}); err != nil {
		return fmt.Errorf("publish StreamElements remove event: %w", err)
	}
	return nil
}

func (s *Service) ImportCommands(
	ctx context.Context,
	channelID, actorID string,
) (importer.Report, error) {
	client, err := s.authorizedClient(ctx, channelID)
	if err != nil {
		return importer.Report{}, err
	}
	profile, err := client.GetProfile(ctx)
	if err != nil {
		return importer.Report{}, fmt.Errorf("get StreamElements profile for command import: %w", err)
	}
	if profile == nil || profile.ID == "" {
		return importer.Report{}, fmt.Errorf("get StreamElements profile for command import: provider returned no channel")
	}
	providerCommands, err := client.GetCommands(ctx, profile.ID)
	if err != nil {
		return importer.Report{}, fmt.Errorf("get StreamElements commands: %w", err)
	}

	commands, normalizationFailures := NormalizeCommands(providerCommands)
	report, err := s.importer.ImportCommands(ctx, channelID, actorID, commands)
	if err != nil {
		return importer.Report{}, fmt.Errorf("import StreamElements commands: %w", err)
	}
	return mergeImportReport(normalizationFailures, report), nil
}

func (s *Service) ImportTimers(
	ctx context.Context,
	channelID, actorID string,
) (importer.Report, error) {
	client, err := s.authorizedClient(ctx, channelID)
	if err != nil {
		return importer.Report{}, err
	}
	profile, err := client.GetProfile(ctx)
	if err != nil {
		return importer.Report{}, fmt.Errorf("get StreamElements profile for timer import: %w", err)
	}
	if profile == nil || profile.ID == "" {
		return importer.Report{}, fmt.Errorf("get StreamElements profile for timer import: provider returned no channel")
	}
	providerTimers, err := client.GetTimers(ctx, profile.ID)
	if err != nil {
		return importer.Report{}, fmt.Errorf("get StreamElements timers: %w", err)
	}

	timers, normalizationFailures := NormalizeTimers(providerTimers)
	report, err := s.importer.ImportTimers(ctx, channelID, actorID, timers)
	if err != nil {
		return importer.Report{}, fmt.Errorf("import StreamElements timers: %w", err)
	}
	return mergeImportReport(normalizationFailures, report), nil
}

func (s *Service) GetTokens(
	ctx context.Context,
	channelID string,
) (streamelementsintegration.Tokens, error) {
	channelIntegration, err := s.channelIntegrationsRepo.GetByChannelAndService(
		ctx,
		channelID,
		integrationsmodel.ServiceStreamElements,
	)
	if err != nil {
		return streamelementsintegration.Tokens{}, fmt.Errorf("get StreamElements tokens: %w", err)
	}
	if channelIntegration.IsNil() || channelIntegration.ID == "" || !channelIntegration.Enabled ||
		channelIntegration.AccessToken == nil || channelIntegration.RefreshToken == nil {
		return streamelementsintegration.Tokens{}, fmt.Errorf("StreamElements integration is not connected")
	}
	return streamelementsintegration.Tokens{
		AccessToken:  *channelIntegration.AccessToken,
		RefreshToken: *channelIntegration.RefreshToken,
	}, nil
}

func (s *Service) UpdateTokens(
	ctx context.Context,
	channelID string,
	tokens streamelementsintegration.Tokens,
) error {
	channelIntegration, err := s.channelIntegrationsRepo.GetByChannelAndService(
		ctx,
		channelID,
		integrationsmodel.ServiceStreamElements,
	)
	if err != nil {
		return fmt.Errorf("get StreamElements integration for token update: %w", err)
	}
	if channelIntegration.IsNil() || channelIntegration.ID == "" {
		return fmt.Errorf("StreamElements integration is not connected")
	}
	if tokens.AccessToken == "" || tokens.RefreshToken == "" {
		return fmt.Errorf("refuse to persist incomplete StreamElements credentials")
	}
	if err := s.channelIntegrationsRepo.Update(ctx, channelIntegration.ID, channelsintegrations.UpdateInput{
		AccessToken:  &tokens.AccessToken,
		RefreshToken: &tokens.RefreshToken,
	}); err != nil {
		return fmt.Errorf("persist StreamElements tokens: %w", err)
	}
	return nil
}

func (s *Service) authorizedClient(ctx context.Context, channelID string) (providerClient, error) {
	if err := s.ensureConfigured(); err != nil {
		return nil, err
	}
	tokens, err := s.GetTokens(ctx, channelID)
	if err != nil {
		return nil, err
	}
	return s.clientFactory.NewAuthorized(
		s.config.StreamElementsClientId,
		s.config.StreamElementsClientSecret,
		channelID,
		s.redirectURL,
		tokens,
		s,
		s.locker,
	), nil
}

func (s *Service) ensureConfigured() error {
	if s.config.StreamElementsClientId == "" || s.config.StreamElementsClientSecret == "" || s.redirectURL == "" {
		return fmt.Errorf("StreamElements integration is not configured")
	}
	return nil
}

func mergeImportReport(normalizationFailures []importer.Failure, report importer.Report) importer.Report {
	failures := make([]importer.Failure, 0, len(normalizationFailures)+len(report.Failures))
	failures = append(failures, normalizationFailures...)
	failures = append(failures, report.Failures...)
	report.Failures = failures
	report.FailedCount = len(failures)
	return report
}
