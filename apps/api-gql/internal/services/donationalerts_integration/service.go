package donationalerts_integration

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"github.com/goccy/go-json"
	"github.com/imroc/req/v3"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/integrations"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/repositories/donationalerts_integration"
	"github.com/twirapp/twir/libs/repositories/donationalerts_integration/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	DonationAlertsRepository donationalerts_integration.Repository
	TwirBus                  *buscore.Bus
	Config                   config.Config
}

func New(opts Opts) *Service {
	return &Service{
		donationAlertsRepository: opts.DonationAlertsRepository,
		twirBus:                  opts.TwirBus,
		config:                   opts.Config,
	}
}

type Service struct {
	donationAlertsRepository donationalerts_integration.Repository
	twirBus                  *buscore.Bus
	config                   config.Config
}

type AuthLinkResponse struct {
	Link string `json:"link"`
}

// mapModelToEntity converts repository model to service entity
func (s *Service) mapModelToEntity(m model.DonationAlertsIntegration) entity.DonationAlertsIntegration {
	var data *entity.ChannelsIntegrationsData

	// Parse JSONB data if it exists
	if m.Data != nil {
		var parsedData entity.ChannelsIntegrationsData
		if err := json.Unmarshal(m.Data, &parsedData); err == nil {
			data = &parsedData
		}
	}

	return entity.DonationAlertsIntegration{
		ID:            m.ID,
		Enabled:       m.Enabled,
		ChannelID:     m.ChannelID,
		IntegrationID: m.IntegrationID,
		AccessToken:   m.AccessToken,
		RefreshToken:  m.RefreshToken,
		ClientID:      m.ClientID,
		ClientSecret:  m.ClientSecret,
		APIKey:        m.APIKey,
		Data:          data,
	}
}

func (s *Service) GetIntegrationData(ctx context.Context, channelID string) (
	entity.DonationAlertsIntegration,
	error,
) {
	integration, err := s.donationAlertsRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		if errors.Is(err, donationalerts_integration.ErrNotFound) {
			// Return default data if integration doesn't exist
			return entity.DonationAlertsIntegration{
				ChannelID: channelID,
				Enabled:   false,
			}, nil
		}
		return entity.DonationAlertsIntegration{}, fmt.Errorf(
			"failed to get donationalerts integration: %w",
			err,
		)
	}

	return s.mapModelToEntity(integration), nil
}

func (s *Service) getCallbackUrl() (string, error) {
	u, err := url.Parse(s.config.SiteBaseUrl)
	if err != nil {
		return "", fmt.Errorf("invalid site base URL: %w", err)
	}

	return u.JoinPath("integrations", "donationalerts", "callback").String(), nil
}

func (s *Service) GetAuthLink(ctx context.Context) (*AuthLinkResponse, error) {
	if s.config.DonationAlertsClientId == "" || s.config.DonationAlertsSecret == "" {
		return nil, errors.New("donationalerts integration not properly configured")
	}

	redirectUrl, err := s.getCallbackUrl()
	if err != nil {
		return nil, fmt.Errorf("failed to get redirect URL: %w", err)
	}

	// Build OAuth authorization URL
	authURL := "https://www.donationalerts.com/oauth/authorize"
	params := url.Values{}
	params.Add("client_id", s.config.DonationAlertsClientId)
	params.Add(
		"redirect_uri",
		redirectUrl,
	)
	params.Add("response_type", "code")
	params.Add(
		"scope",
		"oauth-user-show oauth-donation-subscribe",
	)

	fullURL := fmt.Sprintf("%s?%s", authURL, params.Encode())

	return &AuthLinkResponse{
		Link: fullURL,
	}, nil
}

func (s *Service) PostCode(ctx context.Context, channelID, code string) error {
	_, err := s.donationAlertsRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		if errors.Is(err, donationalerts_integration.ErrNotFound) {
			data := entity.ChannelsIntegrationsData{
				// Username and Avatar would be populated after successful OAuth
			}
			dataBytes, _ := json.Marshal(data)

			err = s.donationAlertsRepository.Create(
				ctx, donationalerts_integration.CreateOpts{
					ChannelID: channelID,
					Enabled:   true,
					Data:      dataBytes,
					// AccessToken and RefreshToken would be set after OAuth exchange
				},
			)
			if err != nil {
				return fmt.Errorf("failed to create donationalerts integration: %w", err)
			}
		}

		return fmt.Errorf("failed to get donationalerts integration: %w", err)
	}

	if s.config.DonationAlertsClientId == "" || s.config.DonationAlertsSecret == "" {
		return errors.New("donationalerts integration not properly configured")
	}

	redirectUrl, err := s.getCallbackUrl()
	if err != nil {
		return fmt.Errorf("failed to get redirect URL: %w", err)
	}

	profile, err := s.getProfileData(
		ctx,
		s.config.DonationAlertsClientId,
		s.config.DonationAlertsSecret,
		redirectUrl,
		code,
	)
	if err != nil {
		return fmt.Errorf("failed to get donationalerts profile data: %w", err)
	}

	// Update existing integration
	enabled := true
	data := entity.ChannelsIntegrationsData{
		// Username and Avatar would be populated after successful OAuth
		UserName: &profile.Data.Name,
		Avatar:   &profile.Data.Avatar,
	}
	dataBytes, _ := json.Marshal(data)

	err = s.donationAlertsRepository.Update(
		ctx, donationalerts_integration.UpdateOpts{
			ChannelID: channelID,
			Enabled:   &enabled,
			Data:      dataBytes,
			// AccessToken and RefreshToken would be updated after OAuth exchange
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update donationalerts integration: %w", err)
	}

	newIntegration, err := s.donationAlertsRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to get donationalerts integration after update: %w", err)
	}

	if err = s.twirBus.Integrations.Add.Publish(
		ctx, integrations.Request{
			ID:      newIntegration.ID.String(),
			Service: integrations.DonationAlerts,
		},
	); err != nil {
		return fmt.Errorf("failed to publish add integration event: %w", err)
	}

	return err
}

func (s *Service) Logout(ctx context.Context, channelID string) error {
	// Disable the integration and clear tokens
	enabled := false
	emptyString := ""
	data := entity.ChannelsIntegrationsData{
		UserName: nil,
		Avatar:   nil,
	}
	dataBytes, _ := json.Marshal(data)

	err := s.donationAlertsRepository.Update(
		ctx, donationalerts_integration.UpdateOpts{
			ChannelID:    channelID,
			Enabled:      &enabled,
			AccessToken:  &emptyString,
			RefreshToken: &emptyString,
			Data:         dataBytes,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to disable donationalerts integration: %w", err)
	}

	return nil
}

type donationAlertsTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
}

type donationAlertsProfileResponse struct {
	Data struct {
		Name   string `json:"name"`
		Code   string `json:"code"`
		Avatar string `json:"avatar"`
	} `json:"data"`
}

func (s *Service) getProfileData(
	ctx context.Context,
	clientId, clientSecret, redirectURL, code string,
) (
	*donationAlertsProfileResponse,
	error,
) {
	data := donationAlertsTokensResponse{}
	resp, err := req.R().
		SetContext(ctx).
		SetFormData(
			map[string]string{
				"grant_type":    "authorization_code",
				"client_id":     clientId,
				"client_secret": clientSecret,
				"redirect_uri":  redirectURL,
				"code":          code,
			},
		).
		SetSuccessResult(&data).
		SetContentType("application/x-www-form-urlencoded").
		Post("https://www.donationalerts.com/oauth/token")
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for tokens: %w", err)
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("failed to exchange code for tokens: %s", resp.String())
	}

	profile := donationAlertsProfileResponse{}
	profileResp, err := req.R().
		SetContext(ctx).
		SetSuccessResult(&profile).
		SetBearerAuthToken(data.AccessToken).
		Get("https://www.donationalerts.com/api/v1/user/oauth")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch donationalerts profile: %w", err)
	}
	if !profileResp.IsSuccessState() {
		return nil, fmt.Errorf("failed to fetch donationalerts profile: %s", profileResp.String())
	}

	return &profile, nil
}
