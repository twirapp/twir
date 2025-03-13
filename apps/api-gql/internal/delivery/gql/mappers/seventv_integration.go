package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func MapSevenTvIntegrationDataToGql(
	data entity.SevenTvIntegrationData,
) *gqlmodel.SevenTvIntegration {
	return &gqlmodel.SevenTvIntegration{
		IsEditor:                   data.IsEditor,
		BotSeventvProfile:          MapSevenTvProfileToGql(data.BotSeventvProfile),
		UserSeventvProfile:         MapSevenTvProfileToGql(data.UserSeventvProfile),
		RewardIDForAddEmote:        data.RewardIDForAddEmote,
		RewardIDForRemoveEmote:     data.RewardIDForRemoveEmote,
		EmoteSetID:                 data.EmoteSetID,
		DeleteEmotesOnlyAddedByApp: data.DeleteEmotesOnlyAddedByApp,
	}
}

func MapSevenTvProfileToGql(profile *entity.SevenTvProfile) *gqlmodel.SevenTvProfile {
	return &gqlmodel.SevenTvProfile{
		ID:          profile.ID,
		Username:    profile.Username,
		DisplayName: profile.DisplayName,
		AvatarURI:   profile.AvatarURI,
	}
}
