package integrations

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/helpers"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/integrations_nightbot"
	"google.golang.org/protobuf/types/known/emptypb"
)

type nightbotTokensResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

type nightbotChannelResponse struct {
	User struct {
		DisplayName string `json:"displayName"`
		Avatar      string `json:"avatar"`
	} `json:"user"`
}

type nightbotCustomCommandsResponse struct {
	Commands []struct {
		Name      string `json:"name"`
		Message   string `json:"message"`
		CoolDown  int    `json:"coolDown"`
		Count     int    `json:"count"`
		UserLevel string `json:"userLevel"`
	} `json:"commands"`
	TotalCount int `json:"_total"`
}

type nightbotTimersResponse struct {
	TotalCount int `json:"_total"`
	Timers     []struct {
		ID       string `json:"_id"`
		Name     string `json:"name"`
		Message  string `json:"message"`
		Interval string `json:"interval"`
		Lines    int    `json:"lines"`
		Enabled  bool   `json:"enabled"`
	}
}

type nightbotRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
}

func (c *Integrations) refreshNightbotTokens(
	ctx context.Context,
	integration *model.ChannelsIntegrations,
) error {
	refreshData := nightbotRefreshResponse{}
	resp, err := req.R().
		SetContext(ctx).
		SetFormData(
			map[string]string{
				"grant_type":    "refresh_token",
				"client_id":     integration.ClientID.String,
				"client_secret": integration.ClientSecret.String,
				"refresh_token": integration.RefreshToken.String,
			},
		).
		SetSuccessResult(&refreshData).
		Post("https://api.nightbot.tv/oauth2/token")
	if err != nil {
		return err
	}

	if !resp.IsSuccessState() {
		return fmt.Errorf("nightbot integration error: %s", resp.String())
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

	commandsData := nightbotCustomCommandsResponse{}
	resp, err := req.R().
		SetContext(ctx).
		SetBearerAuthToken(integration.AccessToken.String).
		SetSuccessResult(&commandsData).
		Get("https://api.nightbot.tv/1/commands")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		if resp.StatusCode == http.StatusUnauthorized {
			err = c.refreshNightbotTokens(ctx, integration)
			if err != nil {
				return nil, err
			}
		}
		return nil, fmt.Errorf("nightbot integration error: %s", resp.String())
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
	broadcasterRole, ok := lo.Find(twirRoles, func(r model.ChannelRole) bool {
		return r.Name == "BROADCASTER"
	})
	if !ok {
		return nil, errors.New("twir internal error")
	}
	moderatorRole, ok := lo.Find(twirRoles, func(r model.ChannelRole) bool {
		return r.Name == "MODERATOR"
	})
	if !ok {
		return nil, errors.New("twir internal error")
	}
	subscriberRole, ok := lo.Find(twirRoles, func(r model.ChannelRole) bool {
		return r.Name == "SUBSCRIBER"
	})
	if !ok {
		return nil, errors.New("twir internal error")
	}
	vipRole, ok := lo.Find(twirRoles, func(r model.ChannelRole) bool {
		return r.Name == "VIP"
	})
	if !ok {
		return nil, errors.New("twir internal error")
	}

	importedCount := 0
	failedCount := 0
	failedCommandsNames := []string{}
	for _, command := range commandsData.Commands {
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
			Aliases:           pq.StringArray{},
			Visible:           true,
			ChannelID:         dashboardId,
			Description:       null.String{},
		}

		newCommand.Responses = append(newCommand.Responses, &model.ChannelsCommandsResponses{
			ID:    uuid.NewString(),
			Text:  null.StringFrom(commandResponse),
			Order: 0,
		})

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

	timersData := nightbotTimersResponse{}
	resp, err := req.R().
		SetContext(ctx).
		SetBearerAuthToken(integration.AccessToken.String).
		SetSuccessResult(&timersData).
		Get("https://api.nightbot.tv/1/timers")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		if resp.StatusCode == http.StatusUnauthorized {
			err = c.refreshNightbotTokens(ctx, integration)
			if err != nil {
				return nil, err
			}
		}
		return nil, fmt.Errorf("nightbot integration error: %s", resp.String())
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

	tokensData := nightbotTokensResponse{}
	resp, err := req.R().
		SetContext(ctx).
		SetFormData(
			map[string]string{
				"grant_type":    "authorization_code",
				"client_id":     channelIntegration.Integration.ClientID.String,
				"client_secret": channelIntegration.Integration.ClientSecret.String,
				"redirect_uri":  channelIntegration.Integration.RedirectURL.String,
				"code":          request.Code,
			},
		).
		SetSuccessResult(&tokensData).
		Post("https://api.nightbot.tv/oauth2/token")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("nightbot token request failed: %s", resp.String())
	}

	channelData := &nightbotChannelResponse{}
	resp, err = req.R().
		SetContext(ctx).
		SetSuccessResult(channelData).
		SetBearerAuthToken(tokensData.AccessToken).
		Get("https://api.nightbot.tv/1/me")
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("nightbot token request failed: %s", resp.String())
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
