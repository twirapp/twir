package resolvers

import (
	"context"
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
	if err := r.gorm.
		WithContext(ctx).
		Raw(
			`
	SELECT
    emote,
    DATE_TRUNC('hour', "createdAt") AS time,
    COUNT(*) AS count
FROM
    channels_emotes_usages
WHERE
    emote = ? AND "createdAt" >= NOW() - INTERVAL '24 hours'
GROUP BY
    emote, time
ORDER BY
    time DESC;
`, emoteName,
		).Scan(&usages).Error; err != nil {
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
