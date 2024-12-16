package dashboard_widget_events

import (
	"context"
	"fmt"

	dbmodel "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type DashboardWidgetsEventsOpts struct {
	fx.In

	Db *gorm.DB
}

func New(opts DashboardWidgetsEventsOpts) *Service {
	return &Service{
		db: opts.Db,
	}
}

type Service struct {
	db *gorm.DB
}

func (d *Service) GetDashboardWidgetsEvents(
	ctx context.Context,
	channelID string,
	limit int,
) ([]entity.DashboardWidgetEvent, error) {

	var entities []dbmodel.ChannelsEventsListItem
	if err := d.db.
		WithContext(ctx).
		Where(
			"channel_id = ?",
			channelID,
		).Order("created_at desc").Limit(limit).Find(&entities).Error; err != nil {
		return nil, fmt.Errorf("cannot get events: %w", err)
	}

	events := make([]entity.DashboardWidgetEvent, 0, limit)

	for _, e := range entities {
		t := entity.DashboardWidgetEventType(e.Type)

		events = append(
			events, entity.DashboardWidgetEvent{
				ID:        e.ID,
				UserID:    e.UserID,
				ChannelID: channelID,
				Type:      t,
				Data: entity.DashboardWidgetEventData{
					DonationAmount:                  e.Data.DonationAmount,
					DonationCurrency:                e.Data.DonationCurrency,
					DonationMessage:                 e.Data.DonationMessage,
					DonationUsername:                e.Data.DonationUsername,
					RaidedViewersCount:              e.Data.RaidedViewersCount,
					RaidedFromUserName:              e.Data.RaidedFromUserName,
					RaidedFromDisplayName:           e.Data.RaidedFromDisplayName,
					FollowUserName:                  e.Data.FollowUserName,
					FollowUserDisplayName:           e.Data.FollowUserDisplayName,
					RedemptionTitle:                 e.Data.RedemptionTitle,
					RedemptionInput:                 e.Data.RedemptionInput,
					RedemptionUserName:              e.Data.RedemptionUserName,
					RedemptionUserDisplayName:       e.Data.RedemptionUserDisplayName,
					RedemptionCost:                  e.Data.RedemptionCost,
					SubLevel:                        e.Data.SubLevel,
					SubUserName:                     e.Data.SubUserName,
					SubUserDisplayName:              e.Data.SubUserDisplayName,
					ReSubLevel:                      e.Data.ReSubLevel,
					ReSubUserName:                   e.Data.ReSubUserName,
					ReSubUserDisplayName:            e.Data.ReSubUserDisplayName,
					ReSubMonths:                     e.Data.ReSubMonths,
					ReSubStreak:                     e.Data.ReSubStreak,
					SubGiftLevel:                    e.Data.SubGiftLevel,
					SubGiftUserName:                 e.Data.SubGiftUserName,
					SubGiftUserDisplayName:          e.Data.SubGiftUserDisplayName,
					SubGiftTargetUserName:           e.Data.SubGiftTargetUserName,
					SubGiftTargetUserDisplayName:    e.Data.SubGiftTargetUserDisplayName,
					FirstUserMessageUserName:        e.Data.FirstUserMessageUserName,
					FirstUserMessageUserDisplayName: e.Data.FirstUserMessageUserDisplayName,
					FirstUserMessageMessage:         e.Data.FirstUserMessageMessage,
					BanReason:                       e.Data.BanReason,
					BanEndsInMinutes:                e.Data.BanEndsInMinutes,
					BannedUserName:                  e.Data.BannedUserName,
					BannedUserLogin:                 e.Data.BannedUserLogin,
					ModeratorName:                   e.Data.ModeratorName,
					ModeratorDisplayName:            e.Data.ModeratorDisplayName,
					Message:                         e.Data.Message,
					UserLogin:                       e.Data.UserLogin,
					UserDisplayName:                 e.Data.UserDisplayName,
				},
				CreatedAt: e.CreatedAt,
			},
		)
	}

	return events, nil
}
