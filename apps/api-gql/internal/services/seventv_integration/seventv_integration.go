package seventv_integration

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/integrations/seventv"
	"github.com/twirapp/twir/libs/repositories/bots"
	seventvintegrationrepository "github.com/twirapp/twir/libs/repositories/seventv_integration"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	SeventvRepository seventvintegrationrepository.Repository
	BotsRepository    bots.Repository
	Redis             *redis.Client
	Config            config.Config
}

func New(opts Opts) *Service {
	return &Service{
		seventvRepository: opts.SeventvRepository,
		botsRepository:    opts.BotsRepository,
		redis:             opts.Redis,
		config:            opts.Config,
	}
}

type Service struct {
	seventvRepository seventvintegrationrepository.Repository
	redis             *redis.Client
	botsRepository    bots.Repository
	config            config.Config
}

func (c *Service) getBotSevenTvProfile(ctx context.Context) (entity.SevenTvProfile, error) {
	defaultBot, err := c.botsRepository.GetDefault(ctx)
	if err != nil {
		return entity.SevenTvProfile{}, fmt.Errorf("failed to get default bot: %w", err)
	}

	client := seventv.NewClient(c.config.SevenTvToken)

	resp, err := client.GetProfileByTwitchId(ctx, defaultBot.ID)
	if err != nil {
		return entity.SevenTvProfile{}, fmt.Errorf("failed to get bot profile: %w", err)
	}

	if resp == nil || resp.Users.UserByConnection == nil {
		return entity.SevenTvProfile{}, fmt.Errorf("failed to get bot profile: %w", err)
	}

	editorFor := make([]entity.SevenTvProfileEditor, 0, len(resp.Users.UserByConnection.EditorFor))
	for _, editor := range resp.Users.UserByConnection.EditorFor {
		var hasEmotesPermissions bool
		if editor.Permissions.SuperAdmin || (editor.Permissions.Emote.Admin && editor.Permissions.EmoteSet.Admin) {
			hasEmotesPermissions = true
		}

		editorFor = append(
			editorFor,
			entity.SevenTvProfileEditor{
				ID:                   editor.EditorId,
				AddedAt:              editor.AddedAt.UnixMilli(),
				HasEmotesPermissions: hasEmotesPermissions,
			},
		)
	}

	var avatarUrl string
	if resp.Users.UserByConnection.MainConnection.PlatformAvatarUrl != nil {
		avatarUrl = *resp.Users.UserByConnection.MainConnection.PlatformAvatarUrl
	}

	profile := entity.SevenTvProfile{
		ID:          resp.Users.UserByConnection.Id,
		Username:    resp.Users.UserByConnection.MainConnection.PlatformUsername,
		DisplayName: resp.Users.UserByConnection.MainConnection.PlatformDisplayName,
		Editors:     nil,
		EditorFor:   editorFor,
		AvatarURI:   avatarUrl,
		EmoteSetID:  resp.Users.UserByConnection.Style.ActiveEmoteSetId,
	}

	return profile, nil
}

func (c *Service) getUserSevenTvResponse(ctx context.Context, userID string) (
	entity.SevenTvProfile,
	error,
) {
	client := seventv.NewClient(c.config.SevenTvToken)

	resp, err := client.GetProfileByTwitchId(ctx, userID)
	if err != nil {
		return entity.SevenTvProfile{}, fmt.Errorf("failed to get user profile: %w", err)
	}

	if resp == nil || resp.Users.UserByConnection == nil {
		return entity.SevenTvProfile{}, fmt.Errorf("failed to get user profile: %w", err)
	}

	editors := make([]entity.SevenTvProfileEditor, 0, len(resp.Users.UserByConnection.Editors))
	for _, editor := range resp.Users.UserByConnection.Editors {
		var hasEmotesPermissions bool
		if editor.Permissions.SuperAdmin || (editor.Permissions.Emote.Admin && editor.Permissions.EmoteSet.Admin) {
			hasEmotesPermissions = true
		}

		editors = append(
			editors,
			entity.SevenTvProfileEditor{
				ID:                   editor.EditorId,
				AddedAt:              editor.AddedAt.UnixMilli(),
				HasEmotesPermissions: hasEmotesPermissions,
			},
		)
	}

	var avatarUrl string
	if resp.Users.UserByConnection.MainConnection.PlatformAvatarUrl != nil {
		avatarUrl = *resp.Users.UserByConnection.MainConnection.PlatformAvatarUrl
	}

	profile := entity.SevenTvProfile{
		ID:          resp.Users.UserByConnection.Id,
		Username:    resp.Users.UserByConnection.MainConnection.PlatformUsername,
		DisplayName: resp.Users.UserByConnection.MainConnection.PlatformDisplayName,
		Editors:     editors,
		EditorFor:   nil,
		AvatarURI:   avatarUrl,
		EmoteSetID:  resp.Users.UserByConnection.Style.ActiveEmoteSetId,
	}

	return profile, nil
}

func (c *Service) GetSevenTvData(
	ctx context.Context,
	channelID string,
) (entity.SevenTvIntegrationData, error) {
	botProfile, err := c.getBotSevenTvProfile(ctx)
	if err != nil {
		return entity.SevenTvIntegrationData{}, fmt.Errorf("failed to get bot profile: %w", err)
	}

	userProfile, err := c.getUserSevenTvResponse(ctx, channelID)
	if err != nil {
		return entity.SevenTvIntegrationData{}, fmt.Errorf("failed to get user profile: %w", err)
	}

	var isBotEditor bool
	if botProfile.ID == userProfile.ID {
		isBotEditor = true
	} else {
		for _, editor := range userProfile.Editors {
			if editor.ID == botProfile.ID {
				isBotEditor = true
				break
			}
		}
	}

	resp := entity.SevenTvIntegrationData{
		IsEditor:                   isBotEditor,
		BotSeventvProfile:          &botProfile,
		UserSeventvProfile:         &userProfile,
		EmoteSetID:                 userProfile.EmoteSetID,
		RewardIDForAddEmote:        nil,
		RewardIDForRemoveEmote:     nil,
		DeleteEmotesOnlyAddedByApp: false,
	}

	integrationData, err := c.seventvRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return entity.SevenTvIntegrationData{}, fmt.Errorf("failed to get integration data: %w", err)
	}

	if integrationData.ChannelID != "" {
		resp.RewardIDForAddEmote = integrationData.RewardIdForAddEmote
		resp.RewardIDForRemoveEmote = integrationData.RewardIdForRemoveEmote
		resp.DeleteEmotesOnlyAddedByApp = integrationData.DeleteEmotesOnlyAddedByApp
	}

	return resp, nil
}

type UpdateInput struct {
	RewardIDForAddEmote        *string
	RewardIDForRemoveEmote     *string
	DeleteEmotesOnlyAddedByApp *bool
}

func (c *Service) UpdateSevenTvData(
	ctx context.Context,
	channelID string,
	input UpdateInput,
) error {
	integrationData, err := c.seventvRepository.GetByChannelID(ctx, channelID)
	if err != nil {
		return fmt.Errorf("failed to get integration data: %w", err)
	}

	return c.seventvRepository.Update(
		ctx,
		integrationData.ID,
		seventvintegrationrepository.UpdateInput{
			RewardIdForAddEmote:        input.RewardIDForAddEmote,
			RewardIdForRemoveEmote:     input.RewardIDForRemoveEmote,
			DeleteEmotesOnlyAddedByApp: input.DeleteEmotesOnlyAddedByApp,
		},
	)
}

type CreateInput struct {
	ChannelID                  string
	RewardIDForAddEmote        *string
	RewardIDForRemoveEmote     *string
	DeleteEmotesOnlyAddedByApp bool
}

func (c *Service) CreateSevenTvData(ctx context.Context, input CreateInput) error {
	return c.seventvRepository.Create(
		ctx,
		seventvintegrationrepository.CreateInput{
			ChannelID:                  input.ChannelID,
			RewardIdForAddEmote:        input.RewardIDForAddEmote,
			RewardIdForRemoveEmote:     input.RewardIDForRemoveEmote,
			DeleteEmotesOnlyAddedByApp: &input.DeleteEmotesOnlyAddedByApp,
		},
	)
}

func (c *Service) CreateOrUpdateSevenTvData(ctx context.Context, input CreateInput) error {
	integrationData, err := c.seventvRepository.GetByChannelID(ctx, input.ChannelID)
	if err != nil {
		return fmt.Errorf("failed to get integration data: %w", err)
	}

	if integrationData.ChannelID == "" {
		return c.CreateSevenTvData(ctx, input)
	}

	return c.UpdateSevenTvData(
		ctx,
		input.ChannelID,
		UpdateInput{
			RewardIDForAddEmote:        input.RewardIDForAddEmote,
			RewardIDForRemoveEmote:     input.RewardIDForRemoveEmote,
			DeleteEmotesOnlyAddedByApp: &input.DeleteEmotesOnlyAddedByApp,
		},
	)
}
