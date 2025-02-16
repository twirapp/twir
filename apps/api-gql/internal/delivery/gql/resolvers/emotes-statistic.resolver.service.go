package resolvers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
)

func (r *queryResolver) getEmoteStatisticUsagesForRange(
	ctx context.Context,
	emoteName string,
	timeRange gqlmodel.EmoteStatisticRange,
) ([]gqlmodel.EmoteStatisticUsage, error) {
	dashboardId, err := r.deps.Sessions.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, err
	}

	var usages []emoteStatisticUsageModel

	var interval string
	var truncateBy string
	switch timeRange {
	case gqlmodel.EmoteStatisticRangeLastDay:
		interval = "24 hours"
		truncateBy = "hour"
	case gqlmodel.EmoteStatisticRangeLastWeek:
		interval = "7 days"
		truncateBy = "day"
	case gqlmodel.EmoteStatisticRangeLastMonth:
		interval = "30 days"
		truncateBy = "day"
	case gqlmodel.EmoteStatisticRangeLastThreeMonth:
		interval = "90 days"
		truncateBy = "day"
	case gqlmodel.EmoteStatisticRangeLastYear:
		interval = "365 days"
		truncateBy = "day"
	default:
	}

	query := fmt.Sprintf(
		`
SELECT
hh AS time, COUNT(emote) AS count
FROM (select
				generate_series(
								DATE_TRUNC(@truncate_by, NOW() - INTERVAL '%s'),
								DATE_TRUNC(@truncate_by, NOW()),
								INTERVAL '1 %s'
				) as hh
) s
left join channels_emotes_usages on DATE_TRUNC(@truncate_by, "createdAt") = hh AND "channelId" = @dashboard_id AND emote = @emote_name
GROUP BY
	time
ORDER BY
	time asc;
	`, interval, truncateBy,
	)

	if err := r.deps.Gorm.
		WithContext(ctx).
		Raw(
			query,
			sql.Named("truncate_by", truncateBy),
			sql.Named("my_interval", interval),
			sql.Named("dashboard_id", dashboardId),
			sql.Named("emote_name", emoteName),
		).
		Find(&usages).Error; err != nil {
		return nil, err
	}

	result := make([]gqlmodel.EmoteStatisticUsage, 0, len(usages))
	for _, usage := range usages {
		result = append(
			result,
			gqlmodel.EmoteStatisticUsage{
				Count:     usage.Count,
				Timestamp: int(usage.Time.UTC().UnixMilli()),
			},
		)
	}

	return result, nil
}

type emoteEntityModelWithCount struct {
	model.ChannelEmoteUsage
	Count int `gorm:"column:count"`
}

type emoteStatisticUsageModel struct {
	Time  time.Time `gorm:"column:time"`
	Emote string    `gorm:"column:emote"`
	Count int       `gorm:"column:count"`
}
