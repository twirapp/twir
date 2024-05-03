package resolvers

import (
	"context"
	"fmt"
	"time"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/gql/gqlmodel"
)

func (r *queryResolver) getEmoteStatisticUsagesForRange(
	ctx context.Context,
	emoteName string,
	timeRange gqlmodel.EmoteStatisticRange,
) ([]gqlmodel.EmoteStatisticUsage, error) {
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
    emote,
    DATE_TRUNC('%s', "createdAt") AS time,
    COUNT(*) AS count
FROM
    channels_emotes_usages
WHERE
    emote = '%s' AND "createdAt" >= NOW() - INTERVAL '%s'
GROUP BY
    emote, time
ORDER BY
    time DESC;
`, truncateBy, emoteName, interval,
	)

	if err := r.gorm.
		Debug().
		WithContext(ctx).
		Raw(query).
		Scan(&usages).Error; err != nil {
		return nil, err
	}

	result := make([]gqlmodel.EmoteStatisticUsage, 0, len(usages))
	for _, usage := range usages {
		result = append(
			result,
			gqlmodel.EmoteStatisticUsage{
				Count:  usage.Count,
				UsedAt: usage.Time,
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
	Emote string    `gorm:"column:emote"`
	Count int       `gorm:"column:count"`
	Time  time.Time `gorm:"column:time"`
}
