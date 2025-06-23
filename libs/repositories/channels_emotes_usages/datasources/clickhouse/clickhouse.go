package clickhouse

import (
	"context"
	"fmt"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/twirapp/twir/libs/baseapp"
	channelsemotesusagesrepository "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	"github.com/twirapp/twir/libs/repositories/channels_emotes_usages/model"
)

type Opts struct {
	Client *baseapp.ClickhouseClient
}

func New(opts Opts) *Clickhouse {
	return &Clickhouse{
		client: opts.Client,
	}
}

func NewFx(client *baseapp.ClickhouseClient) *Clickhouse {
	return New(Opts{Client: client})
}

var _ channelsemotesusagesrepository.Repository = (*Clickhouse)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

type Clickhouse struct {
	client *baseapp.ClickhouseClient
}

func (c *Clickhouse) createBatch(
	ctx context.Context,
	input []channelsemotesusagesrepository.ChannelEmoteUsageInput,
) error {
	if len(input) == 0 {
		return nil
	}

	batch, err := c.client.PrepareBatch(
		ctx,
		"INSERT INTO channels_emotes_usages(channel_id, user_id, emote)",
	)
	if err != nil {
		return fmt.Errorf("prepare batch failed: %w", err)
	}

	for _, i := range input {
		err := batch.Append(
			i.ChannelID,
			i.UserID,
			i.Emote,
		)
		if err != nil {
			return fmt.Errorf("append to batch failed: %w", err)
		}
	}

	err = batch.Send()
	if err != nil {
		return fmt.Errorf("send batch failed: %w", err)
	}

	return nil
}

const batchSize = 1000

func (c *Clickhouse) CreateMany(
	ctx context.Context,
	inputs []channelsemotesusagesrepository.ChannelEmoteUsageInput,
) error {
	if len(inputs) == 0 {
		return nil
	}

	for i := 0; i < len(inputs); i += batchSize {
		end := i + batchSize
		if end > len(inputs) {
			end = len(inputs)
		}
		if err := c.createBatch(ctx, inputs[i:end]); err != nil {
			return err
		}
	}

	return nil
}

func (c *Clickhouse) Count(ctx context.Context, input channelsemotesusagesrepository.CountInput) (
	uint64,
	error,
) {
	selectBuilder := sq.Select("COUNT(*)").From("channels_emotes_usages")

	if input.ChannelID != nil && *input.ChannelID != "" {
		selectBuilder = selectBuilder.Where(squirrel.Eq{"channel_id": *input.ChannelID})
	}

	if input.UserID != nil && *input.UserID != "" {
		selectBuilder = selectBuilder.Where(squirrel.Eq{"user_id": *input.UserID})
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query error: %w", err)
	}

	var count uint64
	err = c.client.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("query error: %w", err)
	}

	return count, nil
}

func (c *Clickhouse) GetEmotesStatistics(
	ctx context.Context,
	input channelsemotesusagesrepository.GetEmotesStatisticsInput,
) ([]model.EmoteStatistic, error) {
	selectBuilder := sq.
		Select("emote", "count(*) as count", "max(created_at) as last_used").
		From("channels_emotes_usages").
		Where(squirrel.Eq{"channel_id": input.ChannelID}).
		GroupBy("emote")

	if input.Sort == channelsemotesusagesrepository.SortAsc {
		selectBuilder = selectBuilder.OrderBy("count ASC")
	} else {
		selectBuilder = selectBuilder.OrderBy("count DESC")
	}

	if input.Search != nil && *input.Search != "" {
		selectBuilder = selectBuilder.Where(
			squirrel.ILike{
				"emote": fmt.Sprintf(
					"%%%s%%",
					*input.Search,
				),
			},
		)
	}

	var (
		page    = input.Page
		perPage = input.PerPage
	)

	if perPage == 0 || perPage > 1000 {
		perPage = 20
	}

	offset := page * perPage

	selectBuilder = selectBuilder.Limit(uint64(perPage)).Offset(uint64(offset))

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query error: %w", err)
	}

	rows, err := c.client.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var result []model.EmoteStatistic
	for rows.Next() {
		var stat model.EmoteStatistic
		err := rows.Scan(&stat.EmoteName, &stat.TotalUsages, &stat.LastUsedTimestamp)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		result = append(result, stat)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", rows.Err())
	}

	return result, nil
}

func (c *Clickhouse) GetEmotesRanges(
	ctx context.Context,
	channelID string,
	emotesNames []string,
	rangeType channelsemotesusagesrepository.EmoteStatisticRange,
) (
	map[string][]model.EmoteRange,
	error,
) {
	if len(emotesNames) == 0 {
		return nil, nil
	}

	// Define time range based on input
	var startTime time.Time
	now := time.Now()
	switch rangeType {
	case channelsemotesusagesrepository.EmoteStatisticRangeLastDay:
		startTime = now.AddDate(0, 0, -1)
	case channelsemotesusagesrepository.EmoteStatisticRangeLastWeek:
		startTime = now.AddDate(0, 0, -7)
	case channelsemotesusagesrepository.EmoteStatisticRangeLastMonth:
		startTime = now.AddDate(0, -1, 0)
	case channelsemotesusagesrepository.EmoteStatisticRangeLastThreeMonth:
		startTime = now.AddDate(0, -3, 0)
	case channelsemotesusagesrepository.EmoteStatisticRangeLastYear:
		startTime = now.AddDate(-1, 0, 0)
	default:
		return nil, fmt.Errorf("invalid range type: %s", rangeType)
	}

	// Determine time bucket based on range
	var timeBucket string
	switch rangeType {
	case channelsemotesusagesrepository.EmoteStatisticRangeLastDay:
		timeBucket = "toStartOfHour(created_at)"
	case channelsemotesusagesrepository.EmoteStatisticRangeLastWeek, channelsemotesusagesrepository.EmoteStatisticRangeLastMonth:
		timeBucket = "toStartOfDay(created_at)"
	case channelsemotesusagesrepository.EmoteStatisticRangeLastThreeMonth, channelsemotesusagesrepository.EmoteStatisticRangeLastYear:
		timeBucket = "toStartOfMonth(created_at)"
	}

	// Construct query
	query := fmt.Sprintf(
		`
		SELECT
			channel_id,
			emote,
			%s AS time_bucket,
			COUNT(*) AS cnt
		FROM twir.channels_emotes_usages
		WHERE created_at >= ?
			AND channel_id = ?
			AND emote IN (?)
		GROUP BY channel_id, emote, time_bucket
		ORDER BY channel_id, emote, time_bucket
	`, timeBucket,
	)

	// Execute query
	rows, err := c.client.Query(ctx, query, startTime, channelID, emotesNames)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	// Process results
	result := make(map[string][]model.EmoteRange)

	for rows.Next() {
		var channelID, emote string
		var timestamp time.Time
		var count uint64

		if err := rows.Scan(&channelID, &emote, &timestamp, &count); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		result[emote] = append(
			result[emote],
			model.EmoteRange{
				Count:     count,
				TimeStamp: timestamp,
			},
		)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return result, nil
}

func (c *Clickhouse) GetChannelEmoteUsageHistory(
	ctx context.Context,
	input channelsemotesusagesrepository.EmotesUsersTopOrHistoryInput,
) ([]model.EmoteUsage, error) {
	var (
		page    = input.Page
		perPage = input.PerPage
	)

	if perPage == 0 || perPage > 1000 {
		perPage = 20
	}

	offset := page * perPage

	query := `
		SELECT
			id,
			channel_id,
			emote,
			user_id,
			created_at
		FROM channels_emotes_usages
		WHERE channel_id = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := c.client.Query(ctx, query, input.ChannelID, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var result []model.EmoteUsage
	for rows.Next() {
		var usage model.EmoteUsage
		err := rows.Scan(&usage.ID, &usage.ChannelID, &usage.Emote, &usage.UserID, &usage.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		result = append(result, usage)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return result, nil
}

func (c *Clickhouse) GetChannelEmoteUsageTopUsers(
	ctx context.Context,
	input channelsemotesusagesrepository.EmotesUsersTopOrHistoryInput,
) ([]model.EmoteUsageTopUser, error) {
	var (
		page    = input.Page
		perPage = input.PerPage
	)

	if perPage == 0 {
		perPage = 20
	}

	if perPage > 1000 {
		perPage = 20
	}

	offset := page * perPage

	query := `
		SELECT
			channel_id,
			user_id,
			COUNT(*) AS count
		FROM channels_emotes_usages
		WHERE channel_id = ?
		GROUP BY channel_id, user_id
		ORDER BY count DESC
		LIMIT ? OFFSET ?
	`

	rows, err := c.client.Query(ctx, query, input.ChannelID, perPage, offset)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var result []model.EmoteUsageTopUser
	for rows.Next() {
		var usage model.EmoteUsageTopUser
		err := rows.Scan(&usage.ChannelID, &usage.UserID, &usage.Count)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		result = append(result, usage)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return result, nil
}
