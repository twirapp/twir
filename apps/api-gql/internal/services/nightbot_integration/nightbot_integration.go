package nightbot_integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/avito-tech/go-transaction-manager/trm/v2"
	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/services/commands"
	"github.com/twirapp/twir/apps/api-gql/internal/services/timers"
	buscore "github.com/twirapp/twir/libs/bus-core"
	timersbusservice "github.com/twirapp/twir/libs/bus-core/timers"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	config "github.com/twirapp/twir/libs/config"
	channelsintegrations "github.com/twirapp/twir/libs/repositories/channels_integrations"
	channelsintegrationsmodel "github.com/twirapp/twir/libs/repositories/channels_integrations/model"
	commandswithgroupsandresponsesmodel "github.com/twirapp/twir/libs/repositories/commands_with_groups_and_responses/model"
	"github.com/twirapp/twir/libs/repositories/integrations"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
	"github.com/twirapp/twir/libs/repositories/plans"
	"github.com/twirapp/twir/libs/repositories/roles"
	rolesmodel "github.com/twirapp/twir/libs/repositories/roles/model"
	timersrepository "github.com/twirapp/twir/libs/repositories/timers"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Config                  config.Config
	Logger                  *slog.Logger
	TrManager               trm.Manager
	TwirBus                 *buscore.Bus
	IntegrationsRepository  integrations.Repository
	ChannelIntegrationsRepo channelsintegrations.Repository
	RolesRepository         roles.Repository
	CommandsService         *commands.Service
	TimersService           *timers.Service
	TimersRepository        timersrepository.Repository
	CachedCommandsClient    *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses]
	PlansRepository         plans.Repository
}

func New(opts Opts) (*Service, error) {
	s := &Service{
		config:                  opts.Config,
		logger:                  opts.Logger,
		trManager:               opts.TrManager,
		twirBus:                 opts.TwirBus,
		integrationsRepo:        opts.IntegrationsRepository,
		channelIntegrationsRepo: opts.ChannelIntegrationsRepo,
		rolesRepository:         opts.RolesRepository,
		commandsService:         opts.CommandsService,
		timersService:           opts.TimersService,
		timersRepository:        opts.TimersRepository,
		cachedCommandsClient:    opts.CachedCommandsClient,
		plansRepository:         opts.PlansRepository,
	}

	siteBaseUrl, err := url.Parse(opts.Config.SiteBaseUrl)
	if err != nil {
		return nil, fmt.Errorf("failed to parse site base URL: %w", err)
	}

	s.redirectURL = siteBaseUrl.JoinPath("/dashboard/integrations/nightbot").String()

	return s, nil
}

type Service struct {
	config                  config.Config
	logger                  *slog.Logger
	trManager               trm.Manager
	twirBus                 *buscore.Bus
	integrationsRepo        integrations.Repository
	channelIntegrationsRepo channelsintegrations.Repository
	rolesRepository         roles.Repository
	commandsService         *commands.Service
	timersService           *timers.Service
	timersRepository        timersrepository.Repository
	cachedCommandsClient    *generic_cacher.GenericCacher[[]commandswithgroupsandresponsesmodel.CommandWithGroupAndResponses]
	plansRepository         plans.Repository
	redirectURL             string
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

func (s *Service) GetAuthLink(ctx context.Context) (string, error) {
	integration, err := s.integrationsRepo.GetByService(ctx, integrationsmodel.ServiceNightbot)
	if err != nil {
		return "", fmt.Errorf("failed to get integration: %w", err)
	}

	if integration.ClientID == nil || integration.ClientSecret == nil || integration.RedirectURL == nil {
		return "", fmt.Errorf("nightbot not enabled on our side, please be patient")
	}

	link, _ := url.Parse("https://api.nightbot.tv/oauth2/authorize")

	q := link.Query()
	q.Add("response_type", "code")
	q.Add("client_id", *integration.ClientID)
	q.Add("scope", "commands commands_default timers regulars spam_protection")
	q.Add("redirect_uri", *integration.RedirectURL)
	link.RawQuery = q.Encode()

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

	tokenReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.nightbot.tv/oauth2/token",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to create token request: %w", err)
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenResp, err := http.DefaultClient.Do(tokenReq)
	if err != nil {
		return fmt.Errorf("failed to get tokens: %w", err)
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode < 200 || tokenResp.StatusCode >= 300 {
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

	if meResp.StatusCode < 200 || meResp.StatusCode >= 300 {
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

	existingIntegration, err := s.channelIntegrationsRepo.GetByChannelAndService(
		ctx,
		channelID,
		integrationsmodel.ServiceNightbot,
	)
	if err != nil {
		return fmt.Errorf("failed to get existing integration: %w", err)
	}

	if existingIntegration.ID != "" {
		err = s.channelIntegrationsRepo.Update(
			ctx, existingIntegration.ID, channelsintegrations.UpdateInput{
				Enabled:      lo.ToPtr(true),
				AccessToken:  &tokensData.AccessToken,
				RefreshToken: &tokensData.RefreshToken,
				Data: &channelsintegrationsmodel.Data{
					UserName: &channelData.User.DisplayName,
					Avatar:   &channelData.User.Avatar,
				},
			},
		)
		if err != nil {
			return fmt.Errorf("failed to update integration: %w", err)
		}
	} else {
		_, err = s.channelIntegrationsRepo.Create(
			ctx, channelsintegrations.CreateInput{
				ChannelID:     channelID,
				IntegrationID: integration.ID,
				Enabled:       true,
				AccessToken:   &tokensData.AccessToken,
				RefreshToken:  &tokensData.RefreshToken,
				Data: &channelsintegrationsmodel.Data{
					UserName: &channelData.User.DisplayName,
					Avatar:   &channelData.User.Avatar,
				},
			},
		)
		if err != nil {
			return fmt.Errorf("failed to create integration: %w", err)
		}
	}

	return nil
}

func (s *Service) Logout(ctx context.Context, channelID string) error {
	channelIntegration, err := s.channelIntegrationsRepo.GetByChannelAndService(
		ctx,
		channelID,
		integrationsmodel.ServiceNightbot,
	)
	if err != nil {
		return fmt.Errorf("failed to get channel integration: %w", err)
	}

	if channelIntegration.ID == "" {
		return nil
	}

	err = s.channelIntegrationsRepo.Update(
		ctx, channelIntegration.ID, channelsintegrations.UpdateInput{
			Enabled:      lo.ToPtr(false),
			AccessToken:  lo.ToPtr(""),
			RefreshToken: lo.ToPtr(""),
			Data: &channelsintegrationsmodel.Data{
				UserName: nil,
				Avatar:   nil,
			},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update integration: %w", err)
	}

	return nil
}

func (s *Service) refreshNightbotTokens(
	ctx context.Context,
	integration integrationsmodel.Integration,
	channelIntegration channelsintegrationsmodel.ChannelIntegration,
) error {
	if channelIntegration.RefreshToken == nil {
		return fmt.Errorf("no refresh token available")
	}

	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("client_id", *integration.ClientID)
	formData.Set("client_secret", *integration.ClientSecret)
	formData.Set("refresh_token", *channelIntegration.RefreshToken)

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.nightbot.tv/oauth2/token",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to create refresh request: %w", err)
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to refresh tokens: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("nightbot refresh error: %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read refresh response: %w", err)
	}

	refreshData := nightbotTokensResponse{}
	if err := json.Unmarshal(bodyBytes, &refreshData); err != nil {
		return fmt.Errorf("failed to unmarshal refresh response: %w", err)
	}

	err = s.channelIntegrationsRepo.Update(
		ctx, channelIntegration.ID, channelsintegrations.UpdateInput{
			Enabled:      lo.ToPtr(true),
			AccessToken:  &refreshData.AccessToken,
			RefreshToken: &refreshData.RefreshToken,
		},
	)
	if err != nil {
		return fmt.Errorf("failed to update tokens: %w", err)
	}

	return nil
}

func (s *Service) ImportCommands(
	ctx context.Context,
	channelID string,
	actorID string,
) (*ImportCommandsResult, error) {
	integration, err := s.integrationsRepo.GetByService(ctx, integrationsmodel.ServiceNightbot)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	channelIntegration, err := s.channelIntegrationsRepo.GetByChannelAndService(
		ctx,
		channelID,
		integrationsmodel.ServiceNightbot,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel integration: %w", err)
	}

	if channelIntegration.ID == "" || channelIntegration.AccessToken == nil {
		return nil, fmt.Errorf("enable nightbot integration first")
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.nightbot.tv/1/commands",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create commands request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+*channelIntegration.AccessToken)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get commands: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		if err := s.refreshNightbotTokens(ctx, integration, channelIntegration); err != nil {
			return nil, fmt.Errorf("failed to refresh tokens: %w", err)
		}
		return s.ImportCommands(ctx, channelID, actorID)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("nightbot commands error: %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read commands response: %w", err)
	}

	commandsData := nightbotCustomCommandsResponse{}
	if err := json.Unmarshal(bodyBytes, &commandsData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal commands: %w", err)
	}

	if len(commandsData.Commands) == 0 {
		return &ImportCommandsResult{
			ImportedCount:       0,
			FailedCount:         0,
			FailedCommandsNames: []string{},
		}, nil
	}

	channelRoles, err := s.rolesRepository.GetManyByChannelID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel roles: %w", err)
	}

	rolesByType := make(map[rolesmodel.ChannelRoleEnum]uuid.UUID)
	for _, role := range channelRoles {
		rolesByType[role.Type] = role.ID
	}

	broadcasterRole := rolesByType[rolesmodel.ChannelRoleTypeBroadcaster]
	moderatorRole := rolesByType[rolesmodel.ChannelRoleTypeModerator]
	subscriberRole := rolesByType[rolesmodel.ChannelRoleTypeSubscriber]
	vipRole := rolesByType[rolesmodel.ChannelRoleTypeVip]

	result := &ImportCommandsResult{
		ImportedCount:       0,
		FailedCount:         0,
		FailedCommandsNames: []string{},
	}

	var commandsToCreate []commands.CreateInput

	for _, command := range commandsData.Commands {
		if command.Alias != nil && *command.Alias != "" {
			continue
		}

		commandName := strings.ToLower(command.Name)
		if len(command.Name) > 0 && command.Name[0] == '!' {
			commandName = commandName[1:]
		}

		var commandRoles []string
		commandResponse := command.Message

		switch command.UserLevel {
		case "admin":
			result.FailedCount++
			result.FailedCommandsNames = append(
				result.FailedCommandsNames,
				command.Name+" (command userlevel is not supported)",
			)
			continue
		case "owner":
			commandRoles = append(commandRoles, broadcasterRole.String())
		case "moderator":
			commandRoles = append(commandRoles, broadcasterRole.String(), moderatorRole.String())
		case "twitch_vip":
			commandRoles = append(
				commandRoles,
				broadcasterRole.String(),
				moderatorRole.String(),
				vipRole.String(),
			)
		case "regular":
			result.FailedCount++
			result.FailedCommandsNames = append(
				result.FailedCommandsNames,
				command.Name+" (command userlevel is not supported)",
			)
			continue
		case "subscriber":
			commandRoles = append(
				commandRoles,
				broadcasterRole.String(),
				moderatorRole.String(),
				subscriberRole.String(),
			)
		case "everyone":
			commandRoles = []string{}
		default:
			result.FailedCount++
			result.FailedCommandsNames = append(
				result.FailedCommandsNames,
				command.Name+" (command userlevel is not supported)",
			)
			continue
		}

		var aliases []string
		for _, cmd := range commandsData.Commands {
			if cmd.Alias != nil && *cmd.Alias == command.Name && strings.HasPrefix(cmd.Name, "!") {
				aliases = append(aliases, cmd.Name[1:])
			}
		}

		commandsToCreate = append(
			commandsToCreate, commands.CreateInput{
				ChannelID:                 channelID,
				ActorID:                   actorID,
				Name:                      commandName,
				Cooldown:                  command.CoolDown,
				CooldownType:              "GLOBAL",
				Enabled:                   true,
				Aliases:                   aliases,
				Description:               "",
				Visible:                   true,
				IsReply:                   true,
				KeepResponsesOrder:        true,
				DeniedUsersIDS:            []string{},
				AllowedUsersIDS:           []string{},
				RolesIDS:                  commandRoles,
				OnlineOnly:                false,
				EnabledCategories:         []string{},
				RequiredWatchTime:         0,
				RequiredMessages:          0,
				RequiredUsedChannelPoints: 0,
				Responses: []commands.CreateInputResponse{
					{
						Text:              &commandResponse,
						Order:             0,
						TwitchCategoryIDs: []string{},
						OnlineOnly:        false,
						OfflineOnly:       false,
					},
				},
			},
		)
	}

	txErr := s.trManager.Do(
		ctx, func(txCtx context.Context) error {
			for _, cmd := range commandsToCreate {
				_, err := s.commandsService.Create(txCtx, cmd)
				if err != nil {
					if strings.Contains(err.Error(), "already exists") {
						result.FailedCount++
						result.FailedCommandsNames = append(
							result.FailedCommandsNames,
							cmd.Name+" (command with this name already exists)",
						)
						continue
					}
					if strings.Contains(err.Error(), "maximum commands limit") {
						result.FailedCount++
						result.FailedCommandsNames = append(
							result.FailedCommandsNames,
							cmd.Name+" (maximum commands limit reached)",
						)
						continue
					}
					result.FailedCount++
					result.FailedCommandsNames = append(
						result.FailedCommandsNames,
						cmd.Name+" (twir internal error)",
					)
					continue
				}
				result.ImportedCount++
			}
			return nil
		},
	)

	if txErr != nil {
		return nil, fmt.Errorf("failed to import commands: %w", txErr)
	}

	if err := s.cachedCommandsClient.Invalidate(ctx, channelID); err != nil {
		s.logger.Error("failed to invalidate commands cache after nightbot import", "error", err)
	}

	return result, nil
}

func (s *Service) ImportTimers(
	ctx context.Context,
	channelID string,
	actorID string,
) (*ImportTimersResult, error) {
	integration, err := s.integrationsRepo.GetByService(ctx, integrationsmodel.ServiceNightbot)
	if err != nil {
		return nil, fmt.Errorf("failed to get integration: %w", err)
	}

	channelIntegration, err := s.channelIntegrationsRepo.GetByChannelAndService(
		ctx,
		channelID,
		integrationsmodel.ServiceNightbot,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get channel integration: %w", err)
	}

	if channelIntegration.ID == "" || channelIntegration.AccessToken == nil {
		return nil, fmt.Errorf("enable nightbot integration first")
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodGet,
		"https://api.nightbot.tv/1/timers",
		nil,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create timers request: %w", err)
	}
	httpReq.Header.Set("Authorization", "Bearer "+*channelIntegration.AccessToken)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("failed to get timers: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		if err := s.refreshNightbotTokens(ctx, integration, channelIntegration); err != nil {
			return nil, fmt.Errorf("failed to refresh tokens: %w", err)
		}
		return s.ImportTimers(ctx, channelID, actorID)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("nightbot timers error: %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read timers response: %w", err)
	}

	timersData := nightbotTimersResponse{}
	if err := json.Unmarshal(bodyBytes, &timersData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal timers: %w", err)
	}

	if len(timersData.Timers) == 0 {
		return &ImportTimersResult{
			ImportedCount:     0,
			FailedCount:       0,
			FailedTimersNames: []string{},
		}, nil
	}

	result := &ImportTimersResult{
		ImportedCount:     0,
		FailedCount:       0,
		FailedTimersNames: []string{},
	}

	currentCount, err := s.timersRepository.CountByChannelID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("cannot get timers count: %w", err)
	}

	plan, err := s.plansRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get plan: %w", err)
	}
	if plan.IsNil() {
		return nil, fmt.Errorf("plan not found for channel")
	}

	spaceLeft := plan.MaxTimers - currentCount
	re := regexp.MustCompile(`\*/(\d+)|(\d+) \* \* \* \*`)

	var timersToCreate []timers.CreateInput
	var createdTimerIDs []uuid.UUID

	for _, timer := range timersData.Timers {
		if spaceLeft <= 0 {
			result.FailedCount++
			result.FailedTimersNames = append(result.FailedTimersNames, timer.Name+" (no space left)")
			continue
		}

		var interval string
		match := re.FindStringSubmatch(timer.Interval)
		for i := 1; i < len(match); i++ {
			if match[i] != "" {
				interval = match[i]
				break
			}
		}

		if interval == "" {
			result.FailedCount++
			result.FailedTimersNames = append(
				result.FailedTimersNames,
				timer.Name+" (invalid timer interval)",
			)
			continue
		}

		parsedInterval, err := strconv.Atoi(interval)
		if err != nil {
			result.FailedCount++
			result.FailedTimersNames = append(
				result.FailedTimersNames,
				timer.Name+" (invalid timer interval)",
			)
			continue
		}
		if parsedInterval == 0 {
			parsedInterval = 60
		}

		timersToCreate = append(
			timersToCreate, timers.CreateInput{
				ChannelID:       channelID,
				ActorID:         actorID,
				Name:            timer.Name,
				Enabled:         timer.Enabled,
				TimeInterval:    parsedInterval,
				MessageInterval: timer.Lines,
				Responses: []timers.CreateResponse{
					{
						Text:       timer.Message,
						IsAnnounce: false,
						Count:      1,
					},
				},
			},
		)
		spaceLeft--
	}

	txErr := s.trManager.Do(
		ctx, func(txCtx context.Context) error {
			for _, timerInput := range timersToCreate {
				createdTimer, err := s.timersService.Create(txCtx, timerInput)
				if err != nil {
					if strings.Contains(err.Error(), "already exists") {
						result.FailedCount++
						result.FailedTimersNames = append(
							result.FailedTimersNames,
							timerInput.Name+" (timer already exists)",
						)
						continue
					}
					result.FailedCount++
					result.FailedTimersNames = append(
						result.FailedTimersNames,
						timerInput.Name+" (twir internal error)",
					)
					continue
				}
				result.ImportedCount++
				createdTimerIDs = append(createdTimerIDs, createdTimer.ID)
			}
			return nil
		},
	)

	if txErr != nil {
		return nil, fmt.Errorf("failed to import timers: %w", txErr)
	}

	for _, timerID := range createdTimerIDs {
		if err := s.twirBus.Timers.AddTimer.Publish(
			ctx,
			timersbusservice.AddOrRemoveTimerRequest{TimerID: timerID.String()},
		); err != nil {
			s.logger.Error("failed to publish add timer bus event after nightbot import", "error", err)
		}
	}

	return result, nil
}
