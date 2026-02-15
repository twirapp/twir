package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	entity "github.com/twirapp/twir/libs/entities/channels_giveaways"
	"github.com/twirapp/twir/libs/entities/channels_giveaways_settings"
)

func GiveawayEntityTo(e entity.Giveaway) gqlmodel.ChannelGiveaway {
	var minWatchedTime, minUsedChannelPoints, minFollowDuration *int
	var minMessages *int

	if e.MinWatchedTime != nil {
		val := int(*e.MinWatchedTime)
		minWatchedTime = &val
	}
	if e.MinMessages != nil {
		val := int(*e.MinMessages)
		minMessages = &val
	}
	if e.MinUsedChannelPoints != nil {
		val := int(*e.MinUsedChannelPoints)
		minUsedChannelPoints = &val
	}
	if e.MinFollowDuration != nil {
		val := int(*e.MinFollowDuration)
		minFollowDuration = &val
	}

	return gqlmodel.ChannelGiveaway{
		ID:                   e.ID.String(),
		ChannelID:            e.ChannelID,
		Type:                 gqlmodel.GiveawayType(e.Type),
		CreatedAt:            e.CreatedAt,
		UpdatedAt:            e.UpdatedAt,
		StartedAt:            e.StartedAt,
		StoppedAt:            e.StoppedAt,
		Keyword:              e.Keyword,
		MinWatchedTime:       minWatchedTime,
		MinMessages:          minMessages,
		MinUsedChannelPoints: minUsedChannelPoints,
		MinFollowDuration:    minFollowDuration,
		RequireSubscription:  e.RequireSubscription,
		CreatedByUserID:      e.CreatedByUserID,
	}
}

func GiveawayParticipantEntityTo(
	e entity.GiveawayParticipant,
) gqlmodel.ChannelGiveawayParticipants {
	return gqlmodel.ChannelGiveawayParticipants{
		DisplayName: e.DisplayName,
		UserID:      e.UserID,
		IsWinner:    e.IsWinner,
		ID:          e.ID.String(),
		GiveawayID:  e.GiveawayID.String(),
	}
}

func GiveawayWinnerEntityTo(
	e entity.GiveawayWinner,
) gqlmodel.ChannelGiveawayWinner {
	return gqlmodel.ChannelGiveawayWinner{
		DisplayName: e.DisplayName,
		UserID:      e.UserID,
		UserLogin:   e.UserLogin,
	}
}

func GiveawaySettingsEntityTo(
	e channels_giveaways_settings.Settings,
) gqlmodel.GiveawaysSettings {
	return gqlmodel.GiveawaysSettings{
		ID:            e.ID,
		ChannelID:     e.ChannelID,
		WinnerMessage: e.WinnerMessage,
	}
}
