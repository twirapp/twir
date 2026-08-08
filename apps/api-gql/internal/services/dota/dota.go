package dota

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/url"
	"strings"

	"github.com/google/uuid"
	channelsservice "github.com/twirapp/twir/apps/api-gql/internal/services/channels"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	buscore "github.com/twirapp/twir/libs/bus-core"
	busapi "github.com/twirapp/twir/libs/bus-core/api"
	cfg "github.com/twirapp/twir/libs/config"
	dotaentity "github.com/twirapp/twir/libs/entities/dota"
	"github.com/twirapp/twir/libs/integrations/steam"
	"github.com/twirapp/twir/libs/logger"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	"github.com/twirapp/twir/libs/repositories/dota/model"
	"github.com/twirapp/twir/libs/wsrouter"
)

const steamCallbackPath = "/dashboard/modules?dotaSteamCallback=1"

type Service struct {
	repository dotarepository.Repository
	steam      *steam.Client
	config     cfg.Config
	logger     *slog.Logger
	wsRouter   wsrouter.WsRouter
	channels   *channelsservice.Service
	twirBus    *buscore.Bus
}

func New(
	repository dotarepository.Repository,
	steamClient *steam.Client,
	config cfg.Config,
	logger *slog.Logger,
	lc *lifecycle.Lifecycle,
	wsRouter wsrouter.WsRouter,
	channelsService *channelsservice.Service,
	twirBus *buscore.Bus,
) *Service {
	s := &Service{
		repository: repository,
		steam:      steamClient,
		config:     config,
		logger:     logger,
		wsRouter:   wsRouter,
		channels:   channelsService,
		twirBus:    twirBus,
	}

	lc.Append(
		lifecycle.Hook{
			OnStart: func(ctx context.Context) error {
				return twirBus.Api.DotaStateUpdate.SubscribeGroup("api", s.handleStateUpdate)
			},
			OnStop: func(ctx context.Context) error {
				twirBus.Api.DotaStateUpdate.Unsubscribe()
				return nil
			},
		},
	)

	return s
}

func (s *Service) handleStateUpdate(
	ctx context.Context,
	msg busapi.DotaStateUpdateMessage,
) (struct{}, error) {
	channelID, err := uuid.Parse(msg.ChannelID)
	if err != nil {
		return struct{}{}, fmt.Errorf("parse dota state update channel id: %w", err)
	}

	return struct{}{}, s.wsRouter.Publish(createStateSubscriptionKey(channelID), msg)
}

func (s *Service) GetOrCreate(ctx context.Context, channelID uuid.UUID) (dotaentity.ChannelDotaSettings, error) {
	settings, err := s.repository.GetByChannelID(ctx, channelID)
	if err == nil {
		return mapModelToEntity(settings), nil
	}
	if !errors.Is(err, dotarepository.ErrNotFound) {
		return dotaentity.Nil, fmt.Errorf("get dota settings: %w", err)
	}

	created, err := s.repository.Create(ctx, dotarepository.CreateInput{ChannelID: channelID})
	if err != nil {
		return dotaentity.Nil, fmt.Errorf("create dota settings: %w", err)
	}

	return mapModelToEntity(created), nil
}

type UpdateInput struct {
	Enabled            bool
	Mmr                int
	MmrDelta           int
	PredictionSettings dotaentity.PredictionSettings
	ChatEvents         dotaentity.ChatEvents
	CommandsSettings   dotaentity.CommandsSettings
}

func (s *Service) Update(
	ctx context.Context,
	channelID uuid.UUID,
	input UpdateInput,
) (dotaentity.ChannelDotaSettings, error) {
	current, err := s.GetOrCreate(ctx, channelID)
	if err != nil {
		return dotaentity.Nil, err
	}

	updated, err := s.repository.Update(
		ctx,
		channelID,
		dotarepository.UpdateInput{
			Enabled:            input.Enabled,
			SteamAccountID:     current.SteamAccountID,
			Mmr:                input.Mmr,
			MmrDelta:           input.MmrDelta,
			PredictionSettings: mapPredictionSettingsToModel(input.PredictionSettings),
			ChatEvents:         mapChatEventsToModel(input.ChatEvents),
			CommandsSettings:   mapCommandsSettingsToModel(input.CommandsSettings),
		},
	)
	if err != nil {
		return dotaentity.Nil, fmt.Errorf("update dota settings: %w", err)
	}

	return mapModelToEntity(updated), nil
}

func (s *Service) ResetSession(ctx context.Context, channelID uuid.UUID) (dotaentity.ChannelDotaSettings, error) {
	if _, err := s.GetOrCreate(ctx, channelID); err != nil {
		return dotaentity.Nil, err
	}

	settings, err := s.repository.ResetSession(ctx, channelID)
	if err != nil {
		return dotaentity.Nil, fmt.Errorf("reset dota session: %w", err)
	}

	return mapModelToEntity(settings), nil
}

func (s *Service) RegenerateGsiToken(
	ctx context.Context,
	channelID uuid.UUID,
) (dotaentity.ChannelDotaSettings, error) {
	if _, err := s.GetOrCreate(ctx, channelID); err != nil {
		return dotaentity.Nil, err
	}

	settings, err := s.repository.RegenerateGsiToken(ctx, channelID)
	if err != nil {
		return dotaentity.Nil, fmt.Errorf("regenerate dota gsi token: %w", err)
	}

	return mapModelToEntity(settings), nil
}

func (s *Service) SteamAuthLink() (string, error) {
	return steam.BuildAuthURL(strings.TrimRight(s.config.SiteBaseUrl, "/") + steamCallbackPath)
}

// SteamLink verifies a Steam OpenID assertion (the raw query string from the
// return_to callback) and binds the account to the channel.
func (s *Service) SteamLink(
	ctx context.Context,
	channelID uuid.UUID,
	queryString string,
) (dotaentity.ChannelDotaSettings, error) {
	current, err := s.GetOrCreate(ctx, channelID)
	if err != nil {
		return dotaentity.Nil, err
	}

	query, err := url.ParseQuery(strings.TrimPrefix(queryString, "?"))
	if err != nil {
		return dotaentity.Nil, fmt.Errorf("parse steam openid params: %w", err)
	}

	steamID64, err := s.steam.VerifyAssertion(ctx, query)
	if err != nil {
		return dotaentity.Nil, fmt.Errorf("verify steam assertion: %w", err)
	}

	updated, err := s.repository.Update(
		ctx,
		channelID,
		dotarepository.UpdateInput{
			Enabled:            current.Enabled,
			SteamAccountID:     &steamID64,
			Mmr:                current.Mmr,
			MmrDelta:           current.MmrDelta,
			PredictionSettings: mapPredictionSettingsToModel(current.PredictionSettings),
			ChatEvents:         mapChatEventsToModel(current.ChatEvents),
			CommandsSettings:   mapCommandsSettingsToModel(current.CommandsSettings),
		},
	)
	if err != nil {
		return dotaentity.Nil, fmt.Errorf("link steam account: %w", err)
	}

	return mapModelToEntity(updated), nil
}

func (s *Service) SteamUnlink(ctx context.Context, channelID uuid.UUID) (dotaentity.ChannelDotaSettings, error) {
	current, err := s.GetOrCreate(ctx, channelID)
	if err != nil {
		return dotaentity.Nil, err
	}

	updated, err := s.repository.Update(
		ctx,
		channelID,
		dotarepository.UpdateInput{
			Enabled:            current.Enabled,
			SteamAccountID:     nil,
			Mmr:                current.Mmr,
			MmrDelta:           current.MmrDelta,
			PredictionSettings: mapPredictionSettingsToModel(current.PredictionSettings),
			ChatEvents:         mapChatEventsToModel(current.ChatEvents),
			CommandsSettings:   mapCommandsSettingsToModel(current.CommandsSettings),
		},
	)
	if err != nil {
		return dotaentity.Nil, fmt.Errorf("unlink steam account: %w", err)
	}

	return mapModelToEntity(updated), nil
}

func (s *Service) SteamProfile(ctx context.Context, steamID64 string) (*steam.PlayerSummary, error) {
	summaries, err := s.steam.GetPlayerSummaries(ctx, []string{steamID64})
	if err != nil {
		s.logger.WarnContext(ctx, "failed to fetch steam player summary", logger.Error(err))
		return nil, nil
	}
	if len(summaries) == 0 {
		return nil, nil
	}

	return &summaries[0], nil
}

// GsiConfig renders the gamestate_integration cfg file for the channel, or an
// empty string when the public GSI endpoint is not configured.
func (s *Service) GsiConfig(settings dotaentity.ChannelDotaSettings) string {
	base := strings.TrimRight(s.config.DotaGsiPublicUrl, "/")
	if base == "" || settings.GsiToken == "" {
		return ""
	}

	return fmt.Sprintf(`"twir"
{
	"uri"		"%s/gsi/%s"
	"timeout"	"5.0"
	"buffer"	"0.5"
	"throttle"	"0.5"
	"heartbeat"	"30.0"
	"auth"
	{
		"token"		"%s"
	}
	"data"
	{
		"provider"		"1"
		"map"			"1"
		"player_id"		"1"
		"player_state"		"1"
		"player_match_stats"	"1"
		"hero"			"1"
		"events"		"1"
	}
}
`, base, settings.GsiToken, settings.GsiToken)
}

func mapModelToEntity(m model.ChannelDotaSettings) dotaentity.ChannelDotaSettings {
	if m.IsNil() {
		return dotaentity.Nil
	}

	return dotaentity.ChannelDotaSettings{
		ID:             m.ID,
		ChannelID:      m.ChannelID,
		Enabled:        m.Enabled,
		SteamAccountID: m.SteamAccountID,
		GsiToken:       m.GsiToken,
		Mmr:            m.Mmr,
		MmrDelta:       m.MmrDelta,
		SessionWins:    m.SessionWins,
		SessionLosses:  m.SessionLosses,
		PredictionSettings: dotaentity.PredictionSettings{
			Enabled:       m.PredictionSettings.Enabled,
			TitleTemplate: m.PredictionSettings.TitleTemplate,
			WindowSeconds: m.PredictionSettings.WindowSeconds,
		},
		ChatEvents: dotaentity.ChatEvents{
			MatchStarted: mapChatEventToEntity(m.ChatEvents.MatchStarted),
			MatchEnded:   mapChatEventToEntity(m.ChatEvents.MatchEnded),
			RoshanKilled: mapChatEventToEntity(m.ChatEvents.RoshanKilled),
			AegisPickup:  mapChatEventToEntity(m.ChatEvents.AegisPickup),
		},
		CommandsSettings: dotaentity.CommandsSettings{
			Mmr: m.CommandsSettings.Mmr,
			Wl:  m.CommandsSettings.Wl,
			Lg:  m.CommandsSettings.Lg,
			Gm:  m.CommandsSettings.Gm,
			Np:  m.CommandsSettings.Np,
			Wp:  m.CommandsSettings.Wp,
		},
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func mapChatEventToEntity(m model.ChatEventSettings) dotaentity.ChatEventSettings {
	return dotaentity.ChatEventSettings{
		Enabled:  m.Enabled,
		Template: m.Template,
		Cooldown: m.Cooldown,
	}
}

func mapPredictionSettingsToModel(e dotaentity.PredictionSettings) model.PredictionSettings {
	return model.PredictionSettings{
		Enabled:       e.Enabled,
		TitleTemplate: e.TitleTemplate,
		WindowSeconds: e.WindowSeconds,
	}
}

func mapChatEventsToModel(e dotaentity.ChatEvents) model.ChatEvents {
	return model.ChatEvents{
		MatchStarted: mapChatEventToModel(e.MatchStarted),
		MatchEnded:   mapChatEventToModel(e.MatchEnded),
		RoshanKilled: mapChatEventToModel(e.RoshanKilled),
		AegisPickup:  mapChatEventToModel(e.AegisPickup),
	}
}

func mapChatEventToModel(e dotaentity.ChatEventSettings) model.ChatEventSettings {
	return model.ChatEventSettings{
		Enabled:  e.Enabled,
		Template: e.Template,
		Cooldown: e.Cooldown,
	}
}

func mapCommandsSettingsToModel(e dotaentity.CommandsSettings) model.CommandsSettings {
	return model.CommandsSettings{
		Mmr: e.Mmr,
		Wl:  e.Wl,
		Lg:  e.Lg,
		Gm:  e.Gm,
		Np:  e.Np,
		Wp:  e.Wp,
	}
}
