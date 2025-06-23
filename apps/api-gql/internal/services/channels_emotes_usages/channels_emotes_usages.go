package channels_emotes_usages

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	channelsemotesusagesrepository "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	ChannelsEmotesUsagesRepository channelsemotesusagesrepository.Repository
}

func New(opts Opts) *Service {
	return &Service{
		channelsEmotesUsagesRepository: opts.ChannelsEmotesUsagesRepository,
	}
}

type Service struct {
	channelsEmotesUsagesRepository channelsemotesusagesrepository.Repository
}

type CountInput struct {
	ChannelID *string
	UserID    *string
}

func (c *Service) Count(ctx context.Context, input CountInput) (uint64, error) {
	return c.channelsEmotesUsagesRepository.Count(
		ctx,
		channelsemotesusagesrepository.CountInput{
			ChannelID: input.ChannelID,
			UserID:    input.UserID,
		},
	)
}

type GetEmotesStatisticsInput struct {
	ChannelID   string
	EmoteSearch *string
	Page        int
	PerPage     int
}

func (c *Service) GetEmotesStatistics(ctx context.Context, input GetEmotesStatisticsInput) (
	[]entity.EmoteStatistic,
	error,
) {
	entities, err := c.channelsEmotesUsagesRepository.GetEmotesStatistics(
		ctx, channelsemotesusagesrepository.GetEmotesStatisticsInput{
			ChannelID: input.ChannelID,
			Search:    input.EmoteSearch,
			Sort:      channelsemotesusagesrepository.SortDesc,
			Page:      input.Page,
			PerPage:   input.PerPage,
		},
	)
	if err != nil {
		return nil, err
	}

	convertedEntities := make([]entity.EmoteStatistic, 0, len(entities))
	for _, e := range entities {
		convertedEntities = append(
			convertedEntities,
			entity.EmoteStatistic{
				EmoteName:         e.EmoteName,
				TotalUsages:       e.TotalUsages,
				LastUsedTimestamp: e.LastUsedTimestamp.UnixMilli(),
			},
		)
	}

	return convertedEntities, nil
}

func (c *Service) GetEmotesRanges(
	ctx context.Context,
	channelID string,
	emotesNames []string,
	rangeType channelsemotesusagesrepository.EmoteStatisticRange,
) (map[string][]entity.EmoteRange, error) {
	ranges, err := c.channelsEmotesUsagesRepository.GetEmotesRanges(
		ctx,
		channelID,
		emotesNames,
		rangeType,
	)
	if err != nil {
		return nil, err
	}

	convertedRanges := make(map[string][]entity.EmoteRange)
	for emoteName, emoteRanges := range ranges {
		convertedRanges[emoteName] = make([]entity.EmoteRange, len(emoteRanges))
		for i, emoteRange := range emoteRanges {
			convertedRanges[emoteName][i] = entity.EmoteRange{
				Count:     emoteRange.Count,
				TimeStamp: emoteRange.TimeStamp.UnixMilli(),
			}
		}
	}

	return convertedRanges, nil
}
