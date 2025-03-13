package seventv_integration

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
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
}

func New(opts Opts) *Service {
	return &Service{
		seventvRepository: opts.SeventvRepository,
		botsRepository:    opts.BotsRepository,
		redis:             opts.Redis,
	}
}

type Service struct {
	seventvRepository seventvintegrationrepository.Repository
	redis             *redis.Client
	botsRepository    bots.Repository
}

const botSevenTvProfileKey = "cache:api:seventv:bot_profile"

func (c *Service) getBotSevenTvProfile(ctx context.Context) (entity.SevenTvProfile, error) {
	cachedBytes, err := c.redis.Get(ctx, botSevenTvProfileKey).Bytes()
	if len(cachedBytes) > 0 {
		parsedCached := entity.SevenTvProfile{}
		if err := json.Unmarshal(cachedBytes, &parsedCached); err != nil {
			return entity.SevenTvProfile{}, err
		}

		return parsedCached, nil
	}

	defaultBot, err := c.botsRepository.GetDefault(ctx)
	if err != nil {
		return entity.SevenTvProfile{}, fmt.Errorf("failed to get default bot: %w", err)
	}

	resp, err := seventv.GetProfile(ctx, defaultBot.ID)
	if err != nil {
		return entity.SevenTvProfile{}, fmt.Errorf("failed to get bot profile: %w", err)
	}

	editors := make([]entity.SevenTvProfileEditor, 0, len(resp.User.Editors))
	for _, editor := range resp.User.Editors {
		editors = append(
			editors,
			entity.SevenTvProfileEditor{
				ID:          editor.Id,
				Permissions: editor.Permissions,
				Visible:     editor.Visible,
				AddedAt:     editor.AddedAt,
			},
		)
	}

	profile := entity.SevenTvProfile{
		ID:          resp.User.Id,
		Username:    resp.Username,
		DisplayName: resp.DisplayName,
		Editors:     editors,
		AvatarURI:   resp.User.AvatarUrl,
	}

	profileBytes, err := json.Marshal(profile)
	if err != nil {
		return entity.SevenTvProfile{}, fmt.Errorf("failed to marshal profile: %w", err)
	}

	if err := c.redis.Set(ctx, botSevenTvProfileKey, profileBytes, 1*time.Hour).Err(); err != nil {
		return entity.SevenTvProfile{}, fmt.Errorf("failed to cache profile: %w", err)
	}

	return profile, nil
}

func (c *Service) getUserSevenTvResponse(ctx context.Context, userID string) (
	entity.SevenTvProfile,
	error,
) {
	resp, err := seventv.GetProfile(ctx, userID)
	if err != nil {
		return entity.SevenTvProfile{}, fmt.Errorf("failed to get bot profile: %w", err)
	}

	editors := make([]entity.SevenTvProfileEditor, 0, len(resp.User.Editors))
	for _, editor := range resp.User.Editors {
		editors = append(
			editors,
			entity.SevenTvProfileEditor{
				ID:          editor.Id,
				Permissions: editor.Permissions,
				Visible:     editor.Visible,
				AddedAt:     editor.AddedAt,
			},
		)
	}

	profile := entity.SevenTvProfile{
		ID:          resp.User.Id,
		Username:    resp.Username,
		DisplayName: resp.DisplayName,
		Editors:     editors,
		AvatarURI:   resp.User.AvatarUrl,
	}

	if resp.EmoteSet != nil {
		profile.EmoteSetID = &resp.EmoteSet.Id
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
