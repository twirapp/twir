package valorantintegration

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/integrations/valorant"
	channelsintegrationsvalorant "github.com/twirapp/twir/libs/repositories/channels_integrations_valorant"
	channelsintegrationsvalorantmodel "github.com/twirapp/twir/libs/repositories/channels_integrations_valorant/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	ChannelsIntegrationsValorantRepository channelsintegrationsvalorant.Repository
	Config                                 cfg.Config
	HenrikApi                              *valorant.HenrikValorantApiClient
}

func New(opts Opts) *Service {
	return &Service{
		repo:      opts.ChannelsIntegrationsValorantRepository,
		henrikApi: opts.HenrikApi,
		config:    opts.Config,
	}
}

type Service struct {
	repo      channelsintegrationsvalorant.Repository
	henrikApi *valorant.HenrikValorantApiClient
	config    cfg.Config
}

type valorantTokenResponse struct {
	Scope        string `json:"scope"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	IdToken      string `json:"id_token"`
	SubSid       string `json:"sub_sid"`
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
}

type valorantAccountResp struct {
	Puuid    string `json:"puuid"`
	UserName string `json:"gameName"`
	TagLine  string `json:"tagLine"`
}

type valorantShardResponse struct {
	Puuid       string `json:"puuid"`
	Game        string `json:"game"`
	ActiveShard string `json:"activeShard"`
}

type valorantHenrikResponse struct {
	Data struct {
		Puuid        string `json:"puuid"`
		Region       string `json:"region"`
		Name         string `json:"name"`
		Tag          string `json:"tag"`
		AccountLevel int    `json:"account_level"`
	} `json:"data"`
	Status int `json:"status"`
}

const apiBase = "https://europe.api.riotgames.com"

func (c *Service) GetChannelStoredMatchesByChannelID(
	ctx context.Context,
	channelID string,
) (*valorant.StoredMatchesResponse, error) {
	integration, err := c.repo.GetByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	if integration == channelsintegrationsvalorantmodel.Nil {
		return nil, fmt.Errorf("no valorant integration found for channel id %s", channelID)
	}

	if integration.ValorantPuuid == nil || integration.ValorantActiveRegion == nil {
		return nil, fmt.Errorf("valorant integration data is incomplete for channel id %s", channelID)
	}

	response, err := c.henrikApi.GetProfileStoredMatches(
		ctx,
		*integration.ValorantActiveRegion,
		*integration.ValorantPuuid,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Service) GetChannelMmr(ctx context.Context, channelID string) (
	*valorant.MmrResponse,
	error,
) {
	integration, err := c.repo.GetByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}
	if integration == channelsintegrationsvalorantmodel.Nil {
		return nil, fmt.Errorf("no valorant integration found for channel id %s", channelID)
	}

	if integration.ValorantPuuid == nil || integration.ValorantActiveRegion == nil {
		return nil, fmt.Errorf("valorant integration data is incomplete for channel id %s", channelID)
	}

	response, err := c.henrikApi.GetValorantProfileMmr(
		ctx,
		"pc",
		*integration.ValorantActiveRegion,
		*integration.ValorantPuuid,
	)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (c *Service) getCallbackUrl(ctx context.Context) (string, error) {
	baseUrl, _ := gincontext.GetBaseUrlFromContext(ctx, c.config.SiteBaseUrl)
	u, err := url.Parse(baseUrl)
	if err != nil {
		return "", fmt.Errorf("invalid site base URL: %w", err)
	}

	return u.JoinPath("dashboard", "integrations", "valorant").String(), nil
}

func (c *Service) GetAuthLink(ctx context.Context) (string, error) {
	if c.config.Valorant.ClientID == "" || c.config.Valorant.ClientSecret == "" {
		return "", fmt.Errorf("valorant not enabled on our side, please be patient")
	}

	redirectUrl, err := c.getCallbackUrl(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get redirect URL: %w", err)
	}

	link, _ := url.Parse("https://auth.riotgames.com/authorize")

	q := link.Query()
	q.Add("response_type", "code")
	q.Add("client_id", c.config.Valorant.ClientID)
	q.Add("scope", strings.Join([]string{"openid", "offline_access", "cpid"}, "+"))
	q.Add("redirect_uri", redirectUrl)
	link.RawQuery = q.Encode()

	return link.String(), nil
}

func (c *Service) GetData(
	ctx context.Context,
	channelID string,
) (*entity.ValorantIntegrationData, error) {
	integration, err := c.repo.GetByChannelID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get valorant data: %w", err)
	}

	if integration.IsNil() {
		return &entity.ValorantIntegrationData{
			Enabled: false,
		}, nil
	}

	result := &entity.ValorantIntegrationData{
		Enabled:  integration.Enabled,
		UserName: integration.UserName,
		// Note: Valorant doesn't provide avatar in the same way, so we'll leave it nil
	}

	return result, nil
}

func (c *Service) PostCode(
	ctx context.Context,
	channelID string,
	code string,
) error {
	if c.config.Valorant.ClientID == "" || c.config.Valorant.ClientSecret == "" {
		return fmt.Errorf("valorant not enabled on our side, please be patient")
	}

	redirectUrl, err := c.getCallbackUrl(ctx)
	if err != nil {
		return fmt.Errorf("failed to get redirect URL: %w", err)
	}

	// Exchange code for tokens
	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("redirect_uri", redirectUrl)
	formData.Set("code", code)

	tokenReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://auth.riotgames.com/token",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	authStr := base64.StdEncoding.EncodeToString(
		[]byte(c.config.Valorant.ClientID + ":" + c.config.Valorant.ClientSecret),
	)
	tokenReq.Header.Set("Authorization", "Basic "+authStr)

	tokenResp, err := http.DefaultClient.Do(tokenReq)
	if err != nil {
		return fmt.Errorf("failed to get valorant tokens: %w", err)
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode < 200 || tokenResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(tokenResp.Body)
		return fmt.Errorf("failed to get valorant tokens: %s", string(bodyBytes))
	}

	tokenBodyBytes, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read token response: %w", err)
	}

	var tokenResponse valorantTokenResponse
	if err := json.Unmarshal(tokenBodyBytes, &tokenResponse); err != nil {
		return fmt.Errorf("failed to parse token response: %w", err)
	}

	// Get account info
	accountReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		apiBase+"/riot/account/v1/accounts/me",
		nil,
	)
	if err != nil {
		return fmt.Errorf("failed to create account request: %w", err)
	}
	accountReq.Header.Set("Authorization", "Bearer "+tokenResponse.AccessToken)

	accountResp, err := http.DefaultClient.Do(accountReq)
	if err != nil {
		return fmt.Errorf("failed to get account info: %w", err)
	}
	defer accountResp.Body.Close()

	if accountResp.StatusCode < 200 || accountResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(accountResp.Body)
		return fmt.Errorf("cannot get valorant account info: %s", string(bodyBytes))
	}

	accountBodyBytes, err := io.ReadAll(accountResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read account response: %w", err)
	}

	var accountResponse valorantAccountResp
	if err := json.Unmarshal(accountBodyBytes, &accountResponse); err != nil {
		return fmt.Errorf("failed to parse account response: %w", err)
	}

	// Get shard info
	shardUrl := fmt.Sprintf(
		"%s/riot/account/v1/active-shards/by-game/%s/by-puuid/%s",
		apiBase,
		"val",
		accountResponse.Puuid,
	)
	shardReq, err := http.NewRequestWithContext(ctx, http.MethodGet, shardUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create shard request: %w", err)
	}
	shardReq.Header.Set("X-Riot-Token", c.config.Valorant.RiotApiKey)

	shardResp, err := http.DefaultClient.Do(shardReq)
	if err != nil {
		return fmt.Errorf("failed to get shard info: %w", err)
	}
	defer shardResp.Body.Close()

	if shardResp.StatusCode < 200 || shardResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(shardResp.Body)
		return fmt.Errorf(
			"cannot get valorant shard info: %v, %s",
			shardResp.StatusCode,
			string(bodyBytes),
		)
	}

	shardBodyBytes, err := io.ReadAll(shardResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read shard response: %w", err)
	}

	var shardResponse valorantShardResponse
	if err := json.Unmarshal(shardBodyBytes, &shardResponse); err != nil {
		return fmt.Errorf("failed to parse shard response: %w", err)
	}

	// Get Henrik profile data
	henrikUrl := fmt.Sprintf(
		"https://api.henrikdev.xyz/valorant/v1/account/%s/%s",
		accountResponse.UserName,
		accountResponse.TagLine,
	)
	henrikReq, err := http.NewRequestWithContext(ctx, http.MethodGet, henrikUrl, nil)
	if err != nil {
		return fmt.Errorf("failed to create henrik request: %w", err)
	}
	henrikReq.Header.Set("Authorization", c.config.Valorant.HenrikApiKey)

	henrikResp, err := http.DefaultClient.Do(henrikReq)
	if err != nil {
		return fmt.Errorf("failed to get henrik profile: %w", err)
	}
	defer henrikResp.Body.Close()

	if henrikResp.StatusCode < 200 || henrikResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(henrikResp.Body)
		return fmt.Errorf(
			"cannot get valorant henrik profile: %d %s",
			henrikResp.StatusCode,
			string(bodyBytes),
		)
	}

	henrikBodyBytes, err := io.ReadAll(henrikResp.Body)
	if err != nil {
		return fmt.Errorf("failed to read henrik response: %w", err)
	}

	var henrikResponse valorantHenrikResponse
	if err := json.Unmarshal(henrikBodyBytes, &henrikResponse); err != nil {
		return fmt.Errorf("failed to parse henrik response: %w", err)
	}

	userName := fmt.Sprintf(
		"%s#%s",
		accountResponse.UserName,
		accountResponse.TagLine,
	)

	// Check if integration already exists
	existing, err := c.repo.GetByChannelID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to check existing integration: %w", err)
	}

	if !existing.IsNil() {
		// Update existing integration
		err = c.repo.Update(
			ctx, existing.ID, channelsintegrationsvalorant.UpdateInput{
				Enabled:              boolPtr(true),
				AccessToken:          &tokenResponse.AccessToken,
				RefreshToken:         &tokenResponse.RefreshToken,
				UserName:             &userName,
				ValorantActiveRegion: &shardResponse.ActiveShard,
				ValorantPuuid:        &henrikResponse.Data.Puuid,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to update valorant integration: %w", err)
		}
	} else {
		// Create new integration
		_, err = c.repo.Create(
			ctx, channelsintegrationsvalorant.CreateInput{
				ChannelID:            channelID,
				Enabled:              true,
				AccessToken:          &tokenResponse.AccessToken,
				RefreshToken:         &tokenResponse.RefreshToken,
				UserName:             &userName,
				ValorantActiveRegion: &shardResponse.ActiveShard,
				ValorantPuuid:        &henrikResponse.Data.Puuid,
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create valorant integration: %w", err)
		}
	}

	return nil
}

func (c *Service) Logout(
	ctx context.Context,
	channelID string,
) error {
	integration, err := c.repo.GetByChannelID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to get integration: %w", err)
	}

	if integration.IsNil() {
		return nil
	}

	return c.repo.Delete(ctx, integration.ID)
}

func boolPtr(b bool) *bool {
	return &b
}
