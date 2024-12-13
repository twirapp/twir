package mappers

import (
	"fmt"

	"github.com/guregu/null"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

var dashboardEventsTypeMap = map[entity.DashboardWidgetEventType]gqlmodel.DashboardEventType{
	entity.TypeDonation:                   gqlmodel.DashboardEventTypeDonation,
	entity.TypeFollow:                     gqlmodel.DashboardEventTypeFollow,
	entity.TypeRaided:                     gqlmodel.DashboardEventTypeRaided,
	entity.TypeSubscribe:                  gqlmodel.DashboardEventTypeSubscribe,
	entity.TypeReSubscribe:                gqlmodel.DashboardEventTypeResubscribe,
	entity.TypeSubGift:                    gqlmodel.DashboardEventTypeSubgift,
	entity.TypeFirstUserMessage:           gqlmodel.DashboardEventTypeFirstUserMessage,
	entity.TypeChatClear:                  gqlmodel.DashboardEventTypeChatClear,
	entity.TypeRedemptionCreated:          gqlmodel.DashboardEventTypeRedemptionCreated,
	entity.TypeChannelBan:                 gqlmodel.DashboardEventTypeChannelBan,
	entity.TypeChannelUnbanRequestCreate:  gqlmodel.DashboardEventTypeChannelUnbanRequestCreate,
	entity.TypeChannelUnbanRequestResolve: gqlmodel.DashboardEventTypeChannelUnbanRequestResolve,
}

func DashboardEventsTypeToGql(t entity.DashboardWidgetEventType) (
	gqlmodel.DashboardEventType,
	error,
) {
	if v, ok := dashboardEventsTypeMap[t]; ok {
		return v, nil
	}

	return "", fmt.Errorf("unknown dashboard event type: %v", t)
}

func DashboardEventsTypeToDb(t gqlmodel.DashboardEventType) (
	entity.DashboardWidgetEventType,
	error,
) {
	for k, v := range dashboardEventsTypeMap {
		if v == t {
			return k, nil
		}
	}

	return "", fmt.Errorf("unknown dashboard event type: %v", t)
}

func DashboardEventsDbToGql(e entity.DashboardWidgetEvent) (gqlmodel.DashboardEventPayload, error) {
	t, err := DashboardEventsTypeToGql(e.Type)
	if err != nil {
		return gqlmodel.DashboardEventPayload{}, err
	}

	return gqlmodel.DashboardEventPayload{
		UserID:    e.UserID,
		Type:      t,
		CreatedAt: e.CreatedAt,
		Data: &gqlmodel.DashboardEventData{
			DonationAmount:                  null.StringFrom(e.Data.DonationAmount).Ptr(),
			DonationCurrency:                null.StringFrom(e.Data.DonationCurrency).Ptr(),
			DonationMessage:                 null.StringFrom(e.Data.DonationMessage).Ptr(),
			DonationUserName:                null.StringFrom(e.Data.DonationUsername).Ptr(),
			RaidedViewersCount:              null.StringFrom(e.Data.RaidedViewersCount).Ptr(),
			RaidedFromUserName:              null.StringFrom(e.Data.RaidedFromUserName).Ptr(),
			RaidedFromDisplayName:           null.StringFrom(e.Data.RaidedFromDisplayName).Ptr(),
			FollowUserName:                  null.StringFrom(e.Data.FollowUserName).Ptr(),
			FollowUserDisplayName:           null.StringFrom(e.Data.FollowUserDisplayName).Ptr(),
			RedemptionTitle:                 null.StringFrom(e.Data.RedemptionTitle).Ptr(),
			RedemptionInput:                 null.StringFrom(e.Data.RedemptionInput).Ptr(),
			RedemptionUserName:              null.StringFrom(e.Data.RedemptionUserName).Ptr(),
			RedemptionUserDisplayName:       null.StringFrom(e.Data.RedemptionUserDisplayName).Ptr(),
			RedemptionCost:                  null.StringFrom(e.Data.RedemptionCost).Ptr(),
			SubLevel:                        null.StringFrom(e.Data.SubLevel).Ptr(),
			SubUserName:                     null.StringFrom(e.Data.SubUserName).Ptr(),
			SubUserDisplayName:              null.StringFrom(e.Data.SubUserDisplayName).Ptr(),
			ReSubLevel:                      null.StringFrom(e.Data.ReSubLevel).Ptr(),
			ReSubUserName:                   null.StringFrom(e.Data.ReSubUserName).Ptr(),
			ReSubUserDisplayName:            null.StringFrom(e.Data.ReSubUserDisplayName).Ptr(),
			ReSubMonths:                     null.StringFrom(e.Data.ReSubMonths).Ptr(),
			ReSubStreak:                     null.StringFrom(e.Data.ReSubStreak).Ptr(),
			SubGiftLevel:                    null.StringFrom(e.Data.SubGiftLevel).Ptr(),
			SubGiftUserName:                 null.StringFrom(e.Data.SubGiftUserName).Ptr(),
			SubGiftUserDisplayName:          null.StringFrom(e.Data.SubGiftUserDisplayName).Ptr(),
			SubGiftTargetUserName:           null.StringFrom(e.Data.SubGiftTargetUserName).Ptr(),
			SubGiftTargetUserDisplayName:    null.StringFrom(e.Data.SubGiftTargetUserDisplayName).Ptr(),
			FirstUserMessageUserName:        null.StringFrom(e.Data.FirstUserMessageUserName).Ptr(),
			FirstUserMessageUserDisplayName: null.StringFrom(e.Data.FirstUserMessageUserDisplayName).Ptr(),
			FirstUserMessageMessage:         null.StringFrom(e.Data.FirstUserMessageMessage).Ptr(),
			BanReason:                       null.StringFrom(e.Data.BanReason).Ptr(),
			BanEndsInMinutes:                null.StringFrom(e.Data.BanEndsInMinutes).Ptr(),
			BannedUserName:                  null.StringFrom(e.Data.BannedUserName).Ptr(),
			BannedUserLogin:                 null.StringFrom(e.Data.BannedUserLogin).Ptr(),
			ModeratorName:                   null.StringFrom(e.Data.ModeratorName).Ptr(),
			ModeratorDisplayName:            null.StringFrom(e.Data.ModeratorDisplayName).Ptr(),
			Message:                         null.StringFrom(e.Data.Message).Ptr(),
			UserLogin:                       null.StringFrom(e.Data.UserLogin).Ptr(),
			UserName:                        null.StringFrom(e.Data.UserDisplayName).Ptr(),
		},
	}, nil
}
