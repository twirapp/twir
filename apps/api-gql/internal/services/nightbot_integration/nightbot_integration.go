package nightbot_integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/services/importer"
	buscore "github.com/twirapp/twir/libs/bus-core"
	buscoretokens "github.com/twirapp/twir/libs/bus-core/tokens"
	channelsintegrations "github.com/twirapp/twir/libs/repositories/channels_integrations"
	channelsintegrationsmodel "github.com/twirapp/twir/libs/repositories/channels_integrations/model"
	"github.com/twirapp/twir/libs/repositories/integrations"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TwirBus                 *buscore.Bus
	Importer                *importer.Service
	IntegrationsRepository  integrations.Repository
	ChannelIntegrationsRepo channelsintegrations.Repository
}

func New(opts Opts) *Service {
	return &Service{
		twirBus:                 opts.TwirBus,
		importer:                opts.Importer,
		httpClient:              http.DefaultClient,
		integrationsRepo:        opts.IntegrationsRepository,
		channelIntegrationsRepo: opts.ChannelIntegrationsRepo,
	}
}

type providerImporter interface {
	ImportCommands(context.Context, string, string, []importer.Command) (importer.Report, error)
	ImportTimers(context.Context, string, string, []importer.Timer) (importer.Report, error)
}

type Service struct {
	twirBus                 *buscore.Bus
	importer                providerImporter
	httpClient              *http.Client
	integrationsRepo        integrations.Repository
	channelIntegrationsRepo channelsintegrations.Repository
}

type nightbotTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
}

type nightbotChannelResponse struct {
	User struct {
		DisplayName string `json:"displayName"`
		Avatar      string `json:"avatar"`
	} `json:"user"`
}

type nightbotCustomCommandsResponse struct {
	Commands []struct {
		Alias     *string `json:"alias,omitempty"`
		Name      string  `json:"name"`
		Message   string  `json:"message"`
		UserLevel string  `json:"userLevel"`
		CoolDown  int     `json:"coolDown"`
		Count     int     `json:"count"`
	} `json:"commands"`
	TotalCount int `json:"_total"`
}

type nightbotTimersResponse struct {
	Timers []struct {
		ID       string `json:"_id"`
		Name     string `json:"name"`
		Message  string `json:"message"`
		Interval string `json:"interval"`
		Lines    int    `json:"lines"`
		Enabled  bool   `json:"enabled"`
	}
	TotalCount int `json:"_total"`
}

type IntegrationData struct {
	UserName string
	Avatar   string
}

type ImportCommandsResult struct {
	ImportedCount       int
	FailedCount         int
	FailedCommandsNames []string
}

type ImportTimersResult struct {
	ImportedCount     int
	FailedCount       int
	FailedTimersNames []string
}

func (s *Service) GetAuthLink(ctx context.Context, states ...string) (string, error) {
	integration, err := s.integrationsRepo.GetByService(ctx, integrationsmodel.ServiceNightbot)
	if err != nil {
		return "", fmt.Errorf("failed to get integration: %w", err)
	}

	if integration.ClientID == nil || integration.ClientSecret == nil || integration.RedirectURL == nil {
		return "", fmt.Errorf("nightbot not enabled on our side, please be patient")
	}

	link, _ := url.Parse("https://api.nightbot.tv/oauth2/authorize")
	query := link.Query()
	query.Add("response_type", "code")
	query.Add("client_id", *integration.ClientID)
	query.Add("scope", "commands commands_default timers regulars spam_protection")
	query.Add("redirect_uri", *integration.RedirectURL)
	if len(states) > 0 {
		query.Add("state", states[0])
	}
	link.RawQuery = query.Encode()

	return link.String(), nil
}

func (s *Service) GetData(ctx context.Context, channelID string) (*IntegrationData, error) {
	channelIntegration, err := s.channelIntegrationsRepo.GetByChannelAndService(
		ctx,
		channelID,
		integrationsmodel.ServiceNightbot,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel integration: %w", err)
	}

	if channelIntegration.ID == "" || channelIntegration.Data == nil {
		return nil, nil
	}

	result := &IntegrationData{}
	if channelIntegration.Data.UserName != nil {
		result.UserName = *channelIntegration.Data.UserName
	}
	if channelIntegration.Data.Avatar != nil {
		result.Avatar = *channelIntegration.Data.Avatar
	}

	return result, nil
}

func (s *Service) PostCode(ctx context.Context, channelID string, code string) error {
	integration, err := s.integrationsRepo.GetByService(ctx, integrationsmodel.ServiceNightbot)
	if err != nil {
		return fmt.Errorf("failed to get integration: %w", err)
	}

	if integration.ClientID == nil || integration.ClientSecret == nil || integration.RedirectURL == nil {
		return fmt.Errorf("nightbot not enabled on our side, please be patient")
	}

	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("client_id", *integration.ClientID)
	formData.Set("client_secret", *integration.ClientSecret)
	formData.Set("redirect_uri", *integration.RedirectURL)
	formData.Set("code", code)

	tokenReq, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.nightbot.tv/oauth2/token", strings.NewReader(formData.Encode()))
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenResp, err := http.DefaultClient.Do(tokenReq)
	if err != nil {
		return fmt.Errorf("failed to get tokens: %w", err)
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode < http.StatusOK || tokenResp.StatusCode >= http.StatusMultipleChoices {
		bodyBytes, _ := io.ReadAll(tokenResp.Body)
		return fmt.Errorf("nightbot token request failed: %s", string(bodyBytes))
	}

	tokenBodyBytes, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read token response: %w", err)
	}

	tokensData := nightbotTokensResponse{}
	if err := json.Unmarshal(tokenBodyBytes, &tokensData); err != nil {
		return fmt.Errorf("failed to unmarshal token response: %w", err)
	}

	meReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.nightbot.tv/1/me", nil)
	if err != nil {
		return fmt.Errorf("failed to create me request: %w", err)
	}
	meReq.Header.Set("Authorization", "Bearer "+tokensData.AccessToken)

	meResp, err := http.DefaultClient.Do(meReq)
	if err != nil {
		return fmt.Errorf("failed to get user info: %w", err)
	}
	defer meResp.Body.Close()

	if meResp.StatusCode < http.StatusOK || meResp.StatusCode >= http.StatusMultipleChoices {
		bodyBytes, _ := io.ReadAll(meResp.Body)
		return fmt.Errorf("nightbot me request failed: %s", string(bodyBytes))
	}

	meBodyBytes, err := io.ReadAll(meResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read me response: %w", err)
	}

	channelData := &nightbotChannelResponse{}
	if err := json.Unmarshal(meBodyBytes, channelData); err != nil {
		return fmt.Errorf("failed to unmarshal me response: %w", err)
	}

	existingIntegration, err := s.channelIntegrationsRepo.GetByChannelAndService(ctx, channelID, integrationsmodel.ServiceNightbot)
	if err != nil {
		return fmt.Errorf("failed to get existing integration: %w", err)
	}

	data := &channelsintegrationsmodel.Data{UserName: &channelData.User.DisplayName, Avatar: &channelData.User.Avatar}
	if existingIntegration.ID != "" {
		if err := s.channelIntegrationsRepo.Update(ctx, existingIntegration.ID, channelsintegrations.UpdateInput{
			Enabled: lo.ToPtr(true), AccessToken: &tokensData.AccessToken, RefreshToken: &tokensData.RefreshToken, Data: data,
		}); err != nil {
			return fmt.Errorf("failed to update integration: %w", err)
		}
	} else {
		if _, err := s.channelIntegrationsRepo.Create(ctx, channelsintegrations.CreateInput{
			ChannelID: channelID, IntegrationID: integration.ID, Enabled: true, AccessToken: &tokensData.AccessToken, RefreshToken: &tokensData.RefreshToken, Data: data,
		}); err != nil {
			return fmt.Errorf("failed to create integration: %w", err)
		}
	}

	return nil
}

func (s *Service) Logout(ctx context.Context, channelID string) error {
	channelIntegration, err := s.channelIntegrationsRepo.GetByChannelAndService(ctx, channelID, integrationsmodel.ServiceNightbot)
	if err != nil {
		return fmt.Errorf("failed to get channel integration: %w", err)
	}
	if channelIntegration.ID == "" {
		return nil
	}

	if err := s.channelIntegrationsRepo.Update(ctx, channelIntegration.ID, channelsintegrations.UpdateInput{
		Enabled: lo.ToPtr(false), AccessToken: lo.ToPtr(""), RefreshToken: lo.ToPtr(""), Data: &channelsintegrationsmodel.Data{},
	}); err != nil {
		return fmt.Errorf("failed to update integration: %w", err)
	}

	return nil
}

func (s *Service) ImportCommands(ctx context.Context, channelID string, actorID string) (*ImportCommandsResult, error) {
	commandsData, err := s.getCommands(ctx, channelID)
	if err != nil {
		return nil, err
	}

	commands, normalizationFailures := NormalizeCommands(commandsData)
	report, err := s.importer.ImportCommands(ctx, channelID, actorID, commands)
	if err != nil {
		return nil, fmt.Errorf("import Nightbot commands: %w", err)
	}

	failures := append(normalizationFailures, report.Failures...)
	return &ImportCommandsResult{
		ImportedCount:       report.ImportedCount,
		FailedCount:         len(failures),
		FailedCommandsNames: failureNames(failures),
	}, nil
}

func (s *Service) ImportTimers(ctx context.Context, channelID string, actorID string) (*ImportTimersResult, error) {
	timersData, err := s.getTimers(ctx, channelID)
	if err != nil {
		return nil, err
	}

	timers, normalizationFailures := NormalizeTimers(timersData)
	report, err := s.importer.ImportTimers(ctx, channelID, actorID, timers)
	if err != nil {
		return nil, fmt.Errorf("import Nightbot timers: %w", err)
	}

	failures := append(normalizationFailures, report.Failures...)
	return &ImportTimersResult{
		ImportedCount:     report.ImportedCount,
		FailedCount:       len(failures),
		FailedTimersNames: failureNames(failures),
	}, nil
}

func (s *Service) getCommands(ctx context.Context, channelID string) (nightbotCustomCommandsResponse, error) {
	accessToken, err := s.accessToken(ctx, channelID)
	if err != nil {
		return nightbotCustomCommandsResponse{}, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.nightbot.tv/1/commands", nil)
	if err != nil {
		return nightbotCustomCommandsResponse{}, fmt.Errorf("failed to create commands request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)

	response, err := s.httpClient.Do(request)
	if err != nil {
		return nightbotCustomCommandsResponse{}, fmt.Errorf("failed to get commands: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(response.Body)
		return nightbotCustomCommandsResponse{}, fmt.Errorf("nightbot commands error: %s", string(body))
	}

	var result nightbotCustomCommandsResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nightbotCustomCommandsResponse{}, fmt.Errorf("failed to unmarshal commands: %w", err)
	}
	return result, nil
}

func (s *Service) getTimers(ctx context.Context, channelID string) (nightbotTimersResponse, error) {
	accessToken, err := s.accessToken(ctx, channelID)
	if err != nil {
		return nightbotTimersResponse{}, err
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.nightbot.tv/1/timers", nil)
	if err != nil {
		return nightbotTimersResponse{}, fmt.Errorf("failed to create timers request: %w", err)
	}
	request.Header.Set("Authorization", "Bearer "+accessToken)

	response, err := s.httpClient.Do(request)
	if err != nil {
		return nightbotTimersResponse{}, fmt.Errorf("failed to get timers: %w", err)
	}
	defer response.Body.Close()
	if response.StatusCode < http.StatusOK || response.StatusCode >= http.StatusMultipleChoices {
		body, _ := io.ReadAll(response.Body)
		return nightbotTimersResponse{}, fmt.Errorf("nightbot timers error: %s", string(body))
	}

	var result nightbotTimersResponse
	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nightbotTimersResponse{}, fmt.Errorf("failed to unmarshal timers: %w", err)
	}
	return result, nil
}

func (s *Service) accessToken(ctx context.Context, channelID string) (string, error) {
	channelIntegration, err := s.channelIntegrationsRepo.GetByChannelAndService(ctx, channelID, integrationsmodel.ServiceNightbot)
	if err != nil {
		return "", fmt.Errorf("failed to get channel integration: %w", err)
	}
	if channelIntegration.ID == "" || channelIntegration.AccessToken == nil {
		return "", fmt.Errorf("enable nightbot integration first")
	}

	accessToken := *channelIntegration.AccessToken
	if channelIntegration.RefreshToken == nil {
		return accessToken, nil
	}

	tokenResponse, err := s.twirBus.Tokens.RequestChannelIntegrationToken.Request(ctx, buscoretokens.GetChannelIntegrationTokenRequest{
		ChannelID: channelID, Service: integrationsmodel.ServiceNightbot,
	})
	if err != nil {
		return "", fmt.Errorf("failed to get nightbot token: %w", err)
	}

	return tokenResponse.Data.AccessToken, nil
}

func failureNames(failures []importer.Failure) []string {
	names := make([]string, 0, len(failures))
	for _, failure := range failures {
		names = append(names, failure.Name)
	}
	return names
}
