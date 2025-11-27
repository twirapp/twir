package integrations

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api/internal/helpers"
	"github.com/twirapp/twir/libs/api/messages/integrations_nightbot"
	model "github.com/twirapp/twir/libs/gomodels"
	"google.golang.org/protobuf/types/known/emptypb"
)

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

type nightbotRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	Scope        string `json:"scope"`
	ExpiresIn    int    `json:"expires_in"`
}

func (c *Integrations) refreshNightbotTokens(
	ctx context.Context,
	integration *model.ChannelsIntegrations,
) error {
	formData := url.Values{}
	formData.Set("grant_type", "refresh_token")
	formData.Set("client_id", integration.ClientID.String)
	formData.Set("client_secret", integration.ClientSecret.String)
	formData.Set("refresh_token", integration.RefreshToken.String)

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.nightbot.tv/oauth2/token",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return err
	}
	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("nightbot integration error: %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	refreshData := nightbotRefreshResponse{}
	if err := json.Unmarshal(bodyBytes, &refreshData); err != nil {
		return err
	}

	integration.AccessToken = null.StringFrom(refreshData.AccessToken)
	integration.RefreshToken = null.StringFrom(refreshData.RefreshToken)
	integration.Enabled = true

	err = c.Db.WithContext(ctx).Save(integration).Error
	if err != nil {
		return err
	}

	return nil
}

func (c *Integrations) IntegrationsNightbotImportCommands(
	ctx context.Context,
	_ *emptypb.Empty,
) (*integrations_nightbot.ImportCommandsResponse, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	integration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceNightbot,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	if !integration.AccessToken.Valid {
		return nil, errors.New("enable nightbot integration first")
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.nightbot.tv/1/commands", nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+integration.AccessToken.String)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode == http.StatusUnauthorized {
			err = c.refreshNightbotTokens(ctx, integration)
			if err != nil {
				return nil, err
			}
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("nightbot integration error: %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	commandsData := nightbotCustomCommandsResponse{}
	if err := json.Unmarshal(bodyBytes, &commandsData); err != nil {
		return nil, err
	}

	if len(commandsData.Commands) == 0 {
		return &integrations_nightbot.ImportCommandsResponse{
			ImportedCount:       0,
			FailedCount:         0,
			FailedCommandsNames: []string{},
		}, nil
	}

	twirRoles := []model.ChannelRole{}
	err = c.Db.Where(`"channelId" = ?`, dashboardId).Find(&twirRoles).Error
	if err != nil {
		return nil, errors.New("twir internal error")
	}
	broadcasterRole, ok := lo.Find(
		twirRoles, func(r model.ChannelRole) bool {
			return r.Type == model.ChannelRoleTypeBroadcaster
		},
	)
	if !ok {
		return nil, errors.New("twir internal error")
	}
	moderatorRole, ok := lo.Find(
		twirRoles, func(r model.ChannelRole) bool {
			return r.Type == model.ChannelRoleTypeModerator
		},
	)
	if !ok {
		return nil, errors.New("twir internal error")
	}
	subscriberRole, ok := lo.Find(
		twirRoles, func(r model.ChannelRole) bool {
			return r.Type == model.ChannelRoleTypeSubscriber
		},
	)
	if !ok {
		return nil, errors.New("twir internal error")
	}
	vipRole, ok := lo.Find(
		twirRoles, func(r model.ChannelRole) bool {
			return r.Type == model.ChannelRoleTypeVip
		},
	)
	if !ok {
		return nil, errors.New("twir internal error")
	}

	importedCount := 0
	failedCount := 0
	failedCommandsNames := []string{}
	for _, command := range commandsData.Commands {
		if command.Alias != nil && *command.Alias != "" {
			continue
		}

		commandName := strings.ToLower(command.Name)
		if command.Name[0] == '!' {
			commandName = commandName[1:]
		}
		commandRoles := []string{}
		commandResponse := command.Message

		twirCommand := model.ChannelsCommands{}
		err = c.Db.Where(`"channelId" = ? AND "name" = ?`, dashboardId, commandName).
			Find(&twirCommand).
			Error
		if err != nil {
			failedCount++
			failedCommandsNames = append(
				failedCommandsNames,
				command.Name+" (twir internal error)",
			)

			continue
		}

		if twirCommand.ID != "" {
			failedCount++
			failedCommandsNames = append(
				failedCommandsNames,
				command.Name+" (command with this name already exists)",
			)
			continue
		}

		switch command.UserLevel {
		case "admin":
			failedCount++
			failedCommandsNames = append(
				failedCommandsNames,
				command.Name+" (command userlevel is not supported)",
			)
			continue
		case "owner":
			commandRoles = append(commandRoles, broadcasterRole.ID)
		case "moderator":
			commandRoles = append(commandRoles, broadcasterRole.ID, moderatorRole.ID)
		case "twitch_vip":
			commandRoles = append(commandRoles, broadcasterRole.ID, moderatorRole.ID, vipRole.ID)
		case "regular":
			failedCount++
			failedCommandsNames = append(
				failedCommandsNames,
				command.Name+" (command userlevel is not supported)",
			)
			continue
		case "subscriber":
			commandRoles = append(
				commandRoles,
				broadcasterRole.ID,
				moderatorRole.ID,
				subscriberRole.ID,
			)
		case "everyone":
			commandRoles = []string{}
		case "default":
			failedCount++
			failedCommandsNames = append(
				failedCommandsNames,
				command.Name+" (command userlevel is not supported)",
			)
		}

		aliases := pq.StringArray{}
		for _, cmd := range commandsData.Commands {
			if cmd.Alias != nil && *cmd.Alias == command.Name && strings.HasPrefix(cmd.Name, "!") {
				aliases = append(aliases, cmd.Name[1:])
			}
		}

		newCommand := model.ChannelsCommands{
			ID:                        uuid.NewString(),
			Name:                      commandName,
			Cooldown:                  null.IntFrom(int64(command.CoolDown)),
			CooldownType:              "GLOBAL",
			Default:                   false,
			DefaultName:               null.String{},
			Module:                    "CUSTOM",
			IsReply:                   true,
			KeepResponsesOrder:        true,
			DeniedUsersIDS:            []string{},
			AllowedUsersIDS:           []string{},
			RolesIDS:                  commandRoles,
			OnlineOnly:                false,
			RequiredWatchTime:         0,
			RequiredMessages:          0,
			RequiredUsedChannelPoints: 0,
			Responses: make(
				[]*model.ChannelsCommandsResponses,
				0,
				1,
			),
			GroupID:           null.String{},
			EnabledCategories: pq.StringArray{},
			CooldownRolesIDs:  pq.StringArray{},
			Enabled:           true,
			Aliases:           aliases,
			Visible:           true,
			ChannelID:         dashboardId,
			Description:       null.String{},
		}

		newCommand.Responses = append(
			newCommand.Responses, &model.ChannelsCommandsResponses{
				ID:                uuid.NewString(),
				Text:              null.StringFrom(commandResponse),
				Order:             0,
				TwitchCategoryIDs: pq.StringArray{},
			},
		)

		err = c.Db.WithContext(ctx).Create(&newCommand).Error
		if err != nil {
			if pgerr, ok := err.(*pgconn.PgError); ok {
				if pgerr.Code == "23505" {
					failedCount++
					failedCommandsNames = append(
						failedCommandsNames,
						command.Name+" (command with this name already exists)",
					)

					continue
				}
			}

			failedCount++
			failedCommandsNames = append(
				failedCommandsNames,
				command.Name+" (twir internal error)",
			)

			continue
		}
		importedCount++
	}

	return &integrations_nightbot.ImportCommandsResponse{
		ImportedCount:       int32(importedCount),
		FailedCount:         int32(failedCount),
		FailedCommandsNames: failedCommandsNames,
	}, nil
}

func (c *Integrations) IntegrationsNightbotImportTimers(
	ctx context.Context,
	_ *emptypb.Empty,
) (*integrations_nightbot.ImportTimersResponse, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	integration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceNightbot,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	if !integration.AccessToken.Valid {
		return nil, errors.New("enable nightbot integration first")
	}

	httpReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.nightbot.tv/1/timers", nil)
	if err != nil {
		return nil, err
	}
	httpReq.Header.Set("Authorization", "Bearer "+integration.AccessToken.String)

	resp, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		if resp.StatusCode == http.StatusUnauthorized {
			err = c.refreshNightbotTokens(ctx, integration)
			if err != nil {
				return nil, err
			}
		}
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("nightbot integration error: %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	timersData := nightbotTimersResponse{}
	if err := json.Unmarshal(bodyBytes, &timersData); err != nil {
		return nil, err
	}

	if len(timersData.Timers) == 0 {
		return &integrations_nightbot.ImportTimersResponse{
			ImportedCount:     0,
			FailedCount:       0,
			FailedTimersNames: []string{},
		}, nil
	}

	importedCount := 0
	failedCount := 0
	failedTimersNames := []string{}

	var currentCount int64
	if err := c.Db.Model(&model.ChannelsTimers{}).Where(
		`"channelId" = ?`,
		dashboardId,
	).Count(&currentCount).Error; err != nil {
		return nil, fmt.Errorf("cannot get timers count time: %w", err)
	}

	spaceLeft := 10 - currentCount
	re := regexp.MustCompile(`\*/(\d+)|(\d+) \* \* \* \*`)
	for _, timer := range timersData.Timers {
		if spaceLeft == 0 {
			failedCount++
			failedTimersNames = append(failedTimersNames, timer.Name+" (no space left)")
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
			failedCount++
			failedTimersNames = append(
				failedTimersNames,
				timer.Name+" (invalid timer interval)",
			)
			continue
		}
		parsedInterval, err := strconv.Atoi(interval)
		if parsedInterval == 0 {
			parsedInterval = 60
		}

		if err != nil {
			failedCount++
			failedTimersNames = append(
				failedTimersNames,
				timer.Name+" (invalid timer interval)",
			)
			continue
		}

		entity := &model.ChannelsTimers{
			ID:              uuid.NewString(),
			ChannelID:       dashboardId,
			Name:            timer.Name,
			Enabled:         timer.Enabled,
			TimeInterval:    int32(parsedInterval),
			MessageInterval: int32(timer.Lines),
			Responses: []*model.ChannelsTimersResponses{
				{
					ID:         uuid.NewString(),
					Text:       timer.Message,
					IsAnnounce: false,
				},
			},
		}

		if err := c.Db.WithContext(ctx).Create(&entity).Error; err != nil {
			if pgerr, ok := err.(*pgconn.PgError); ok {
				if pgerr.Code == "23505" {
					failedCount++
					failedTimersNames = append(
						failedTimersNames,
						timer.Name+" (timer already exists)",
					)
					continue
				}
			}

			failedCount++
			failedTimersNames = append(
				failedTimersNames,
				timer.Name+" (twir internal error)",
			)
			continue
		}

		importedCount++
		spaceLeft--
	}

	return &integrations_nightbot.ImportTimersResponse{
		ImportedCount:     int32(importedCount),
		FailedCount:       int32(failedCount),
		FailedTimersNames: failedTimersNames,
	}, nil
}

func (c *Integrations) IntegrationsNightbotGetAuthLink(
	ctx context.Context,
	_ *emptypb.Empty,
) (*integrations_nightbot.GetAuthLink, error) {
	integration, err := c.getIntegrationByService(ctx, model.IntegrationServiceNightbot)
	if err != nil {
		return nil, err
	}

	if !integration.ClientID.Valid || !integration.ClientSecret.Valid ||
		!integration.RedirectURL.Valid {
		return nil, errors.New("nightbot not enabled on our side, please be patient")
	}

	link, _ := url.Parse("https://api.nightbot.tv/oauth2/authorize")

	q := link.Query()
	q.Add("response_type", "code")
	q.Add("client_id", integration.ClientID.String)
	q.Add("scope", "commands commands_default timers regulars spam_protection")
	q.Add("redirect_uri", integration.RedirectURL.String)
	link.RawQuery = q.Encode()

	return &integrations_nightbot.GetAuthLink{
		Link: link.String(),
	}, nil
}

func (c *Integrations) IntegrationsNightbotGetData(
	ctx context.Context,
	_ *emptypb.Empty,
) (*integrations_nightbot.GetDataResponse, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	integration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceNightbot,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	return &integrations_nightbot.GetDataResponse{
		UserName: integration.Data.UserName,
		Avatar:   integration.Data.Avatar,
	}, nil
}

func (c *Integrations) IntegrationsNightbotPostCode(
	ctx context.Context,
	request *integrations_nightbot.PostCodeRequest,
) (*emptypb.Empty, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	channelIntegration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceNightbot,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("client_id", channelIntegration.Integration.ClientID.String)
	formData.Set("client_secret", channelIntegration.Integration.ClientSecret.String)
	formData.Set("redirect_uri", channelIntegration.Integration.RedirectURL.String)
	formData.Set("code", request.Code)

	tokenReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://api.nightbot.tv/oauth2/token",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return nil, err
	}
	tokenReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	tokenResp, err := http.DefaultClient.Do(tokenReq)
	if err != nil {
		return nil, err
	}
	defer tokenResp.Body.Close()

	if tokenResp.StatusCode < 200 || tokenResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(tokenResp.Body)
		return nil, fmt.Errorf("nightbot token request failed: %s", string(bodyBytes))
	}

	tokenBodyBytes, err := io.ReadAll(tokenResp.Body)
	if err != nil {
		return nil, err
	}

	tokensData := nightbotTokensResponse{}
	if err := json.Unmarshal(tokenBodyBytes, &tokensData); err != nil {
		return nil, err
	}

	meReq, err := http.NewRequestWithContext(ctx, http.MethodGet, "https://api.nightbot.tv/1/me", nil)
	if err != nil {
		return nil, err
	}
	meReq.Header.Set("Authorization", "Bearer "+tokensData.AccessToken)

	meResp, err := http.DefaultClient.Do(meReq)
	if err != nil {
		return nil, err
	}
	defer meResp.Body.Close()

	if meResp.StatusCode < 200 || meResp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(meResp.Body)
		return nil, fmt.Errorf("nightbot token request failed: %s", string(bodyBytes))
	}

	meBodyBytes, err := io.ReadAll(meResp.Body)
	if err != nil {
		return nil, err
	}

	channelData := &nightbotChannelResponse{}
	if err := json.Unmarshal(meBodyBytes, channelData); err != nil {
		return nil, err
	}

	channelIntegration.Data = &model.ChannelsIntegrationsData{
		UserName: &channelData.User.DisplayName,
		Avatar:   &channelData.User.Avatar,
	}
	channelIntegration.AccessToken = null.StringFrom(tokensData.AccessToken)
	channelIntegration.RefreshToken = null.StringFrom(tokensData.RefreshToken)
	channelIntegration.Enabled = true

	if err = c.Db.WithContext(ctx).Save(channelIntegration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Integrations) IntegrationsNightbotLogout(
	ctx context.Context,
	empty *emptypb.Empty,
) (*emptypb.Empty, error) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	integration, err := c.getChannelIntegrationByService(
		ctx,
		model.IntegrationServiceNightbot,
		dashboardId,
	)
	if err != nil {
		return nil, err
	}

	integration.Data = &model.ChannelsIntegrationsData{}
	integration.AccessToken = null.String{}
	integration.RefreshToken = null.String{}
	integration.Enabled = false

	if err = c.Db.WithContext(ctx).Save(&integration).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}
