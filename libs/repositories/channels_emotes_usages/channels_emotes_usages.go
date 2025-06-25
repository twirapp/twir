package channels_emotes_usages

import (
	"context"

	"github.com/twirapp/twir/libs/repositories/channels_emotes_usages/model"
)

type Repository interface {
	CreateMany(ctx context.Context, inputs []ChannelEmoteUsageInput) error
	Count(ctx context.Context, input CountInput) (uint64, error)
	GetEmotesStatistics(ctx context.Context, input GetEmotesStatisticsInput) (
		[]model.EmoteStatistic,
		error,
	)
	GetEmotesRanges(
		ctx context.Context,
		channelID string,
		emotesNames []string,
		rangeType EmoteStatisticRange,
	) (map[string][]model.EmoteRange, error)
	GetChannelEmoteUsageHistory(
		ctx context.Context,
		input EmotesUsersTopOrHistoryInput,
	) ([]model.EmoteUsage, uint64, error)
	GetChannelEmoteUsageTopUsers(
		ctx context.Context,
		input EmotesUsersTopOrHistoryInput,
	) ([]model.EmoteUsageTopUser, uint64, error)
}

type ChannelEmoteUsageInput struct {
	ChannelID string
	UserID    string
	Emote     string
}

type CountInput struct {
	ChannelID *string
	UserID    *string
}

type GetEmotesStatisticsInput struct {
	ChannelID string
	Search    *string
	Sort      Sort
	Page      int
	PerPage   int
}

type GetEmotesRangesInput struct {
	ChannelID string
	EmoteName string
}

type EmoteStatisticRange string

const (
	EmoteStatisticRangeLastDay        EmoteStatisticRange = "LAST_DAY"
	EmoteStatisticRangeLastWeek       EmoteStatisticRange = "LAST_WEEK"
	EmoteStatisticRangeLastMonth      EmoteStatisticRange = "LAST_MONTH"
	EmoteStatisticRangeLastThreeMonth EmoteStatisticRange = "LAST_THREE_MONTH"
	EmoteStatisticRangeLastYear       EmoteStatisticRange = "LAST_YEAR"
)

type EmotesUsersTopOrHistoryInput struct {
	ChannelID string
	EmoteName string
	Page      int
	PerPage   int
}

type Sort string

const (
	SortAsc  Sort = "ASC"
	SortDesc Sort = "DESC"
)
