package dashboard_widget_events

import (
	"context"
	"fmt"

	dbmodel "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/services/dashboard-widget-events/model"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type DashboardWidgetsEventsOpts struct {
	fx.In

	Db *gorm.DB
}

func New(opts DashboardWidgetsEventsOpts) *DashboardWidgetsEvents {
	return &DashboardWidgetsEvents{
		db: opts.Db,
	}
}

type DashboardWidgetsEvents struct {
	db *gorm.DB
}

func (d *DashboardWidgetsEvents) GetDashboardWidgetsEvents(
	ctx context.Context,
	channelID string,
	limit int,
) ([]model.Event, error) {

	var entities []dbmodel.ChannelsEventsListItem
	if err := d.db.
		WithContext(ctx).
		Where(
			"channel_id = ?",
			channelID,
		).Order("created_at desc").Limit(limit).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("cannot get events: %w", err)
	}

	events := make([]model.Event, 0, limit)

	for _, entity := range entities {
		t := model.EventType(entity.Type)

		events = append(
			events, model.Event{
				ID:        entity.ID,
				UserID:    entity.UserID,
				ChannelID: channelID,
				Type:      t,
				Data: model.Data{
					DonationAmount:                  entity.Data.DonationAmount,
					DonationCurrency:                entity.Data.DonationCurrency,
					DonationMessage:                 entity.Data.DonationMessage,
					DonationUsername:                entity.Data.DonationUsername,
					RaidedViewersCount:              entity.Data.RaidedViewersCount,
					RaidedFromUserName:              entity.Data.RaidedFromUserName,
					RaidedFromDisplayName:           entity.Data.RaidedFromDisplayName,
					FollowUserName:                  entity.Data.FollowUserName,
					FollowUserDisplayName:           entity.Data.FollowUserDisplayName,
					RedemptionTitle:                 entity.Data.RedemptionTitle,
					RedemptionInput:                 entity.Data.RedemptionInput,
					RedemptionUserName:              entity.Data.RedemptionUserName,
					RedemptionUserDisplayName:       entity.Data.RedemptionUserDisplayName,
					RedemptionCost:                  entity.Data.RedemptionCost,
					SubLevel:                        entity.Data.SubLevel,
					SubUserName:                     entity.Data.SubUserName,
					SubUserDisplayName:              entity.Data.SubUserDisplayName,
					ReSubLevel:                      entity.Data.ReSubLevel,
					ReSubUserName:                   entity.Data.ReSubUserName,
					ReSubUserDisplayName:            entity.Data.ReSubUserDisplayName,
					ReSubMonths:                     entity.Data.ReSubMonths,
					ReSubStreak:                     entity.Data.ReSubStreak,
					SubGiftLevel:                    entity.Data.SubGiftLevel,
					SubGiftUserName:                 entity.Data.SubGiftUserName,
					SubGiftUserDisplayName:          entity.Data.SubGiftUserDisplayName,
					SubGiftTargetUserName:           entity.Data.SubGiftTargetUserName,
					SubGiftTargetUserDisplayName:    entity.Data.SubGiftTargetUserDisplayName,
					FirstUserMessageUserName:        entity.Data.FirstUserMessageUserName,
					FirstUserMessageUserDisplayName: entity.Data.FirstUserMessageUserDisplayName,
					FirstUserMessageMessage:         entity.Data.FirstUserMessageMessage,
					BanReason:                       entity.Data.BanReason,
					BanEndsInMinutes:                entity.Data.BanEndsInMinutes,
					BannedUserName:                  entity.Data.BannedUserName,
					BannedUserLogin:                 entity.Data.BannedUserLogin,
					ModeratorName:                   entity.Data.ModeratorName,
					ModeratorDisplayName:            entity.Data.ModeratorDisplayName,
					Message:                         entity.Data.Message,
					UserLogin:                       entity.Data.UserLogin,
					UserDisplayName:                 entity.Data.UserDisplayName,
				},
				CreatedAt: entity.CreatedAt,
			},
		)
	}

	return events, nil
}
