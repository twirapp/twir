package dataloader

import (
	"context"
	"fmt"

	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
)

func (c *dataLoader) getEmoteStatistic(
	ctx context.Context,
	emotes []EmoteRangeKey,
) ([][]gqlmodel.EmoteStatisticUsage, []error) {
	dashboardID, err := c.deps.AuthService.GetSelectedDashboard(ctx)
	if err != nil {
		return nil, []error{err}
	}

	emotesNames := make([]string, 0, len(emotes))
	for _, e := range emotes {
		emotesNames = append(emotesNames, e.EmoteName)
	}

	var emoteRange channels_emotes_usages.EmoteStatisticRange
	switch emotes[0].Range {
	case gqlmodel.EmoteStatisticRangeLastDay:
		emoteRange = channels_emotes_usages.EmoteStatisticRangeLastDay
	case gqlmodel.EmoteStatisticRangeLastWeek:
		emoteRange = channels_emotes_usages.EmoteStatisticRangeLastWeek
	case gqlmodel.EmoteStatisticRangeLastMonth:
		emoteRange = channels_emotes_usages.EmoteStatisticRangeLastMonth
	case gqlmodel.EmoteStatisticRangeLastThreeMonth:
		emoteRange = channels_emotes_usages.EmoteStatisticRangeLastThreeMonth
	case gqlmodel.EmoteStatisticRangeLastYear:
		emoteRange = channels_emotes_usages.EmoteStatisticRangeLastYear
	default:
		return nil, []error{fmt.Errorf("unknown range: %s", emotes[0].Range)}
	}

	ranges, err := c.deps.EmoteStatisticService.GetEmotesRanges(
		ctx,
		dashboardID,
		emotesNames,
		emoteRange,
	)
	if err != nil {
		return nil, []error{err}
	}

	result := make([][]gqlmodel.EmoteStatisticUsage, 0, len(emotes))
	for _, e := range emotes {
		foundRange, ok := ranges[e.EmoteName]
		if !ok {
			result = append(result, nil)
			continue
		}

		mappedRange := make([]gqlmodel.EmoteStatisticUsage, 0, len(foundRange))
		for _, r := range foundRange {
			mappedRange = append(
				mappedRange,
				gqlmodel.EmoteStatisticUsage{
					Count:     int(r.Count),
					Timestamp: int(r.TimeStamp),
				},
			)
		}

		result = append(result, mappedRange)
	}

	return result, nil
}

type EmoteRangeKey struct {
	ChannelID string
	EmoteName string
	Range     gqlmodel.EmoteStatisticRange
}

func GetEmoteStatisticByName(ctx context.Context, key EmoteRangeKey) (
	[]gqlmodel.EmoteStatisticUsage,
	error,
) {
	loaders := GetLoaderForRequest(ctx)
	return loaders.emoteStatistic.Load(ctx, key)
}
