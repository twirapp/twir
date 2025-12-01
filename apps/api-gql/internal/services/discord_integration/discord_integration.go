package discord_integration

import (
	"cmp"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/discord"
	cfg "github.com/twirapp/twir/libs/config"
	channelsintegrationsdiscord "github.com/twirapp/twir/libs/repositories/channels_integrations_discord"
	"github.com/twirapp/twir/libs/repositories/channels_integrations_discord/model"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

type Opts struct {
	fx.In

	DiscordRepository channelsintegrationsdiscord.Repository
	Config            cfg.Config
	Bus               *buscore.Bus
}

func New(opts Opts) *Service {
	return &Service{
		repo:   opts.DiscordRepository,
		config: opts.Config,
		bus:    opts.Bus,
	}
}

type Service struct {
	repo   channelsintegrationsdiscord.Repository
	config cfg.Config
	bus    *buscore.Bus
}

func (s *Service) GetAuthLink(_ context.Context) (string, error) {
	if s.config.DiscordClientID == "" || s.config.DiscordClientSecret == "" {
		return "", fmt.Errorf("discord not enabled on our side, please be patient")
	}

	u, _ := url.Parse("https://discord.com/oauth2/authorize")

	redirectUrl := fmt.Sprintf("%s/dashboard/integrations/discord", s.config.SiteBaseUrl)

	q := u.Query()
	q.Add("client_id", s.config.DiscordClientID)
	q.Add("response_type", "code")
	q.Add("permissions", "1497333180438")
	q.Add("scope", "bot applications.commands")
	q.Add("redirect_uri", redirectUrl)
	u.RawQuery = q.Encode()

	return u.String(), nil
}

func (s *Service) GetData(ctx context.Context, channelID string) (
	*entity.DiscordIntegrationData,
	error,
) {
	guilds, err := s.repo.GetByChannelID(ctx, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to get discord guilds: %w", err)
	}

	if len(guilds) == 0 {
		return &entity.DiscordIntegrationData{
			Guilds: []entity.DiscordGuild{},
		}, nil
	}

	result := make([]entity.DiscordGuild, 0, len(guilds))
	resultMu := sync.Mutex{}

	guildsGroup, gCtx := errgroup.WithContext(ctx)
	for _, guild := range guilds {
		guild := guild

		guildsGroup.Go(
			func() error {
				guildInfoResp, err := s.bus.Discord.GetGuildInfo.Request(
					gCtx, discord.GetGuildInfoRequest{
						GuildID: guild.GuildID,
					},
				)
				if err != nil {
					// Skip guild if we can't get info (bot might have been kicked)
					return nil
				}

				guildInfo := guildInfoResp.Data

				var icon *string
				if guildInfo.Icon != "" {
					icon = &guildInfo.Icon
				}

				resultMu.Lock()
				result = append(
					result, entity.DiscordGuild{
						ID:                               guild.GuildID,
						Name:                             guildInfo.Name,
						Icon:                             icon,
						LiveNotificationEnabled:          guild.LiveNotificationEnabled,
						LiveNotificationChannelsIds:      guild.LiveNotificationChannelsIds,
						LiveNotificationShowTitle:        guild.LiveNotificationShowTitle,
						LiveNotificationShowCategory:     guild.LiveNotificationShowCategory,
						LiveNotificationShowViewers:      guild.LiveNotificationShowViewers,
						LiveNotificationMessage:          guild.LiveNotificationMessage,
						LiveNotificationShowPreview:      guild.LiveNotificationShowPreview,
						LiveNotificationShowProfileImage: guild.LiveNotificationShowProfileImage,
						OfflineNotificationMessage:       guild.OfflineNotificationMessage,
						ShouldDeleteMessageOnOffline:     guild.ShouldDeleteMessageOnOffline,
						AdditionalUsersIdsForLiveCheck:   guild.AdditionalUsersIdsForLiveCheck,
					},
				)
				resultMu.Unlock()

				return nil
			},
		)
	}

	if err := guildsGroup.Wait(); err != nil {
		return nil, fmt.Errorf("failed to get guilds info: %w", err)
	}

	slices.SortFunc(
		result, func(a, b entity.DiscordGuild) int {
			return cmp.Compare(strings.ToLower(a.Name), strings.ToLower(b.Name))
		},
	)

	return &entity.DiscordIntegrationData{
		Guilds: result,
	}, nil
}

type discordPostCodeResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refresh_token"`
	Guild        struct {
		Id string `json:"id"`
	} `json:"guild"`
}

func (s *Service) ConnectGuild(ctx context.Context, channelID, code string) error {
	if s.config.DiscordClientID == "" || s.config.DiscordClientSecret == "" {
		return fmt.Errorf("discord not enabled on our side, please be patient")
	}

	res := discordPostCodeResponse{}

	formData := url.Values{}
	formData.Set("grant_type", "authorization_code")
	formData.Set("code", code)
	formData.Set(
		"redirect_uri", fmt.Sprintf(
			"%s/dashboard/integrations/discord",
			s.config.SiteBaseUrl,
		),
	)

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"https://discord.com/api/oauth2/token",
		strings.NewReader(formData.Encode()),
	)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	authStr := base64.StdEncoding.EncodeToString(
		[]byte(s.config.DiscordClientID + ":" + s.config.DiscordClientSecret),
	)
	httpReq.Header.Set("Authorization", "Basic "+authStr)

	r, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to get token: %w", err)
	}
	defer r.Body.Close()

	if r.StatusCode < 200 || r.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(r.Body)
		return fmt.Errorf("failed to get token: %s", string(bodyBytes))
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if err := json.Unmarshal(bodyBytes, &res); err != nil {
		return fmt.Errorf("failed to unmarshal response: %w", err)
	}

	if res.Guild.Id == "" {
		return fmt.Errorf("failed to get guild id")
	}

	// Check if the guild is already connected
	existing, err := s.repo.GetByChannelIDAndGuildID(ctx, channelID, res.Guild.Id)
	if err != nil {
		return fmt.Errorf("failed to check existing guild: %w", err)
	}

	if !existing.IsNil() {
		// Guild already connected
		return nil
	}

	// Create new guild connection
	_, err = s.repo.Create(
		ctx, channelsintegrationsdiscord.CreateInput{
			ChannelID:                        channelID,
			GuildID:                          res.Guild.Id,
			LiveNotificationEnabled:          false,
			LiveNotificationChannelsIds:      []string{},
			LiveNotificationShowTitle:        true,
			LiveNotificationShowCategory:     true,
			LiveNotificationShowViewers:      true,
			LiveNotificationMessage:          "",
			LiveNotificationShowPreview:      true,
			LiveNotificationShowProfileImage: true,
			OfflineNotificationMessage:       "",
			ShouldDeleteMessageOnOffline:     false,
			AdditionalUsersIdsForLiveCheck:   []string{},
		},
	)
	if err != nil {
		return fmt.Errorf("failed to create guild connection: %w", err)
	}

	return nil
}

func (s *Service) DisconnectGuild(ctx context.Context, channelID, guildID string) error {
	return s.repo.DeleteByChannelIDAndGuildID(ctx, channelID, guildID)
}

func (s *Service) UpdateGuild(
	ctx context.Context,
	channelID,
	guildID string,
	input UpdateGuildInput,
) error {
	existing, err := s.repo.GetByChannelIDAndGuildID(ctx, channelID, guildID)
	if err != nil {
		return fmt.Errorf("failed to get guild: %w", err)
	}

	if existing.IsNil() {
		return fmt.Errorf("guild not found")
	}

	return s.repo.Update(
		ctx, existing.ID, channelsintegrationsdiscord.UpdateInput{
			LiveNotificationEnabled:          input.LiveNotificationEnabled,
			LiveNotificationChannelsIds:      input.LiveNotificationChannelsIds,
			LiveNotificationShowTitle:        input.LiveNotificationShowTitle,
			LiveNotificationShowCategory:     input.LiveNotificationShowCategory,
			LiveNotificationShowViewers:      input.LiveNotificationShowViewers,
			LiveNotificationMessage:          input.LiveNotificationMessage,
			LiveNotificationShowPreview:      input.LiveNotificationShowPreview,
			LiveNotificationShowProfileImage: input.LiveNotificationShowProfileImage,
			OfflineNotificationMessage:       input.OfflineNotificationMessage,
			ShouldDeleteMessageOnOffline:     input.ShouldDeleteMessageOnOffline,
			AdditionalUsersIdsForLiveCheck:   input.AdditionalUsersIdsForLiveCheck,
		},
	)
}

type UpdateGuildInput struct {
	LiveNotificationEnabled          *bool
	LiveNotificationChannelsIds      *[]string
	LiveNotificationShowTitle        *bool
	LiveNotificationShowCategory     *bool
	LiveNotificationShowViewers      *bool
	LiveNotificationMessage          *string
	LiveNotificationShowPreview      *bool
	LiveNotificationShowProfileImage *bool
	OfflineNotificationMessage       *string
	ShouldDeleteMessageOnOffline     *bool
	AdditionalUsersIdsForLiveCheck   *[]string
}

func discordChannelTypeToEntity(t discord.ChannelType) entity.DiscordChannelType {
	switch t {
	case discord.ChannelTypeVoice:
		return entity.DiscordChannelTypeVoice
	case discord.ChannelTypeText:
		return entity.DiscordChannelTypeText
	default:
		return entity.DiscordChannelTypeText
	}
}

func (s *Service) GetGuildChannels(
	ctx context.Context,
	guildID string,
) ([]entity.DiscordGuildChannel, error) {
	channelsResp, err := s.bus.Discord.GetGuildChannels.Request(
		ctx, discord.GetGuildChannelsRequest{
			GuildID: guildID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get guild channels: %w", err)
	}

	channels := make([]entity.DiscordGuildChannel, 0, len(channelsResp.Data.Channels))
	for _, channel := range channelsResp.Data.Channels {
		channels = append(
			channels, entity.DiscordGuildChannel{
				ID:              channel.ID,
				Name:            channel.Name,
				Type:            discordChannelTypeToEntity(channel.Type),
				CanSendMessages: channel.CanSendMessages,
			},
		)
	}

	return channels, nil
}

func (s *Service) GetGuildInfo(ctx context.Context, guildID string) (
	*entity.DiscordGuildInfo,
	error,
) {
	guildInfoResp, err := s.bus.Discord.GetGuildInfo.Request(
		ctx, discord.GetGuildInfoRequest{
			GuildID: guildID,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to get guild info: %w", err)
	}

	guildInfo := guildInfoResp.Data

	channels := make([]entity.DiscordGuildChannel, 0, len(guildInfo.Channels))
	for _, channel := range guildInfo.Channels {
		channels = append(
			channels, entity.DiscordGuildChannel{
				ID:              channel.ID,
				Name:            channel.Name,
				Type:            discordChannelTypeToEntity(channel.Type),
				CanSendMessages: channel.CanSendMessages,
			},
		)
	}

	roles := make([]entity.DiscordGuildRole, 0, len(guildInfo.Roles))
	for _, role := range guildInfo.Roles {
		roles = append(
			roles, entity.DiscordGuildRole{
				ID:    role.ID,
				Name:  role.Name,
				Color: role.Color,
			},
		)
	}

	var icon *string
	if guildInfo.Icon != "" {
		icon = &guildInfo.Icon
	}

	return &entity.DiscordGuildInfo{
		ID:       guildInfo.ID,
		Name:     guildInfo.Name,
		Icon:     icon,
		Channels: channels,
		Roles:    roles,
	}, nil
}

// GetByChannelID returns all guilds connected to a channel (for use by other services like discord app)
func (s *Service) GetByChannelID(
	ctx context.Context,
	channelID string,
) ([]model.ChannelIntegrationDiscord, error) {
	return s.repo.GetByChannelID(ctx, channelID)
}

// GetByGuildID returns all channel integrations for a guild (for use by other services like discord app)
func (s *Service) GetByGuildID(
	ctx context.Context,
	guildID string,
) ([]model.ChannelIntegrationDiscord, error) {
	return s.repo.GetByGuildID(ctx, guildID)
}

// GetByChannelIDAndGuildID returns a specific guild connection (for use by other services like discord app)
func (s *Service) GetByChannelIDAndGuildID(
	ctx context.Context,
	channelID, guildID string,
) (model.ChannelIntegrationDiscord, error) {
	return s.repo.GetByChannelIDAndGuildID(ctx, channelID, guildID)
}

// DeleteByGuildID deletes all channel integrations for a guild (for use when bot is kicked from guild)
func (s *Service) DeleteByGuildID(ctx context.Context, guildID string) error {
	guilds, err := s.repo.GetByGuildID(ctx, guildID)
	if err != nil {
		return err
	}

	for _, guild := range guilds {
		if err := s.repo.Delete(ctx, guild.ID); err != nil {
			return err
		}
	}

	return nil
}
