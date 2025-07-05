package clickhouse

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	twirclickhouse "github.com/twirapp/twir/libs/baseapp/clickhouse"
	channelsemotesusagesrepository "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
	"github.com/twirapp/twir/libs/repositories/channels_emotes_usages/model"
)

type Opts struct {
	Client *twirclickhouse.ClickhouseClient
}

func New(opts Opts) *Clickhouse {
	return &Clickhouse{
		client: opts.Client,
	}
}

func NewFx(client *twirclickhouse.ClickhouseClient) *Clickhouse {
	return New(Opts{Client: client})
}

var _ channelsemotesusagesrepository.Repository = (*Clickhouse)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

type Clickhouse struct {
	client *twirclickhouse.ClickhouseClient
}

func (c *Clickhouse) GetUserMostUsedEmotes(
	ctx context.Context,
	input channelsemotesusagesrepository.UserMostUsedEmotesInput,
) ([]model.UserMostUsedEmote, error) {
	limit := input.Limit
	if limit == 0 || limit > 50 {
		limit = 10
	}

	query := `
SELECT emote, COUNT(*)
FROM channels_emotes_usages
WHERE "channel_id" = ? AND "user_id" = ?
GROUP BY emote
ORDER BY COUNT(*)
DESC LIMIT ?
`

	rows, err := c.client.Query(ctx, query, input.ChannelID, input.UserID, limit)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var result []model.UserMostUsedEmote
	for rows.Next() {
		var emote string
		var count uint64
		err := rows.Scan(&emote, &count)
		if err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		result = append(
			result,
			model.UserMostUsedEmote{
				Emote: emote,
				Count: count,
			},
		)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows error: %w", rows.Err())
	}

	return result, nil
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

	if input.TimeAfter != nil {
		selectBuilder = selectBuilder.Where(squirrel.Gt{"created_at": *input.TimeAfter})
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
		return nil, nil // Return empty map or nil, nil as per original logic
	}

	// Define time range and bucketing parameters
	var startTime time.Time       // The start of the historical period (e.g., 1 day ago)
	var timeBucketFunc string     // ClickHouse function like 'toStartOfHour' or 'toStartOfDay'
	var intervalUnit string       // 'HOUR' or 'DAY' for ClickHouse INTERVAL clause
	var startBucketTime time.Time // The start of the very first time bucket to generate
	var endBucketTime time.Time   // The start of the very last time bucket to generate (usually current hour/day)

	now := time.Now()

	switch rangeType {
	case channelsemotesusagesrepository.EmoteStatisticRangeLastDay:
		startTime = now.AddDate(0, 0, -1)
		timeBucketFunc = "toStartOfHour"
		intervalUnit = "HOUR"
		// Ensure buckets align with hour boundaries
		startBucketTime = startTime.Truncate(time.Hour)
		endBucketTime = now.Truncate(time.Hour)
	case channelsemotesusagesrepository.EmoteStatisticRangeLastWeek:
		startTime = now.AddDate(0, 0, -7)
		timeBucketFunc = "toStartOfDay"
		intervalUnit = "DAY"
		// Ensure buckets align with day boundaries
		startBucketTime = startTime.Truncate(24 * time.Hour)
		endBucketTime = now.Truncate(24 * time.Hour)
	case channelsemotesusagesrepository.EmoteStatisticRangeLastMonth:
		startTime = now.AddDate(0, -1, 0)
		timeBucketFunc = "toStartOfDay"
		intervalUnit = "DAY"
		startBucketTime = startTime.Truncate(24 * time.Hour)
		endBucketTime = now.Truncate(24 * time.Hour)
	case channelsemotesusagesrepository.EmoteStatisticRangeLastThreeMonth:
		startTime = now.AddDate(0, -3, 0)
		timeBucketFunc = "toStartOfDay"
		intervalUnit = "DAY"
		startBucketTime = startTime.Truncate(24 * time.Hour)
		endBucketTime = now.Truncate(24 * time.Hour)
	case channelsemotesusagesrepository.EmoteStatisticRangeLastYear:
		startTime = now.AddDate(-1, 0, 0)
		timeBucketFunc = "toStartOfDay"
		intervalUnit = "DAY"
		startBucketTime = startTime.Truncate(24 * time.Hour)
		endBucketTime = now.Truncate(24 * time.Hour)
	default:
		return nil, fmt.Errorf("invalid range type: %s", rangeType)
	}

	// Calculate the number of intervals needed for numbers() function
	var durationInUnits int
	if intervalUnit == "HOUR" {
		durationInUnits = int(endBucketTime.Sub(startBucketTime).Hours())
	} else { // DAY
		durationInUnits = int(endBucketTime.Sub(startBucketTime).Hours() / 24)
	}
	numbersCount := durationInUnits + 1 // +1 to include both start and end buckets

	if numbersCount <= 0 {
		// No intervals to generate (e.g., startBucketTime is after endBucketTime)
		return make(map[string][]model.EmoteRange), nil
	}

	// Dynamically build the dimensions CTE (channel_id, emote combinations)
	var emoteDimensionsBuilder strings.Builder
	for i, emote := range emotesNames {
		if i > 0 {
			emoteDimensionsBuilder.WriteString(" UNION ALL ")
		}
		// Using single quotes for string literals in SQL.
		// ClickHouse driver should handle escaping channelID and emote if they contain single quotes.
		// If not, consider `strings.ReplaceAll(value, "'", "''")` for robustness.
		emoteDimensionsBuilder.WriteString(
			fmt.Sprintf(
				"SELECT '%s' AS channel_id, '%s' AS emote",
				channelID,
				emote,
			),
		)
	}
	emoteDimensions := emoteDimensionsBuilder.String()

	// Construct the full ClickHouse query
	// This query generates a full time series, cross joins with specified dimensions,
	// and left joins the actual aggregated data to fill missing points with 0.
	query := fmt.Sprintf(
		`
    WITH time_series AS (
        SELECT
            %s(toDateTime(?)) + INTERVAL number %s AS time_bucket
        FROM numbers(%d)
    ),
    dimensions AS (
        %s
    )
    SELECT
        ts.time_bucket,
        d.channel_id,
        d.emote,
        COALESCE(tce.cnt, 0) AS cnt
    FROM time_series AS ts
    CROSS JOIN dimensions AS d
    LEFT JOIN (
        SELECT
            channel_id,
            emote,
            %s(created_at) AS time_bucket,
            COUNT(*) AS cnt
        FROM twir.channels_emotes_usages
        WHERE created_at >= toDateTime(?) -- Use original startTime for filtering raw data
          AND channel_id = ?
          AND emote IN (?)
        GROUP BY channel_id, emote, time_bucket
    ) AS tce ON ts.time_bucket = tce.time_bucket
            AND d.channel_id = tce.channel_id
            AND d.emote = tce.emote
    ORDER BY ts.time_bucket ASC, d.channel_id ASC, d.emote ASC;
    `,
		timeBucketFunc,  // 1. For time_series CTE (e.g., toStartOfHour)
		intervalUnit,    // 2. For time_series CTE (e.g., HOUR)
		numbersCount,    // 3. For numbers() table function (e.g., 25)
		emoteDimensions, // 4. For dimensions CTE (dynamically built UNION ALL)
		timeBucketFunc,  // 5. For the inner aggregation query (e.g., toStartOfHour)
	)

	// Execute query
	// The parameters map to the '?' placeholders in order:
	// 1. `startBucketTime` for the `toDateTime(?)` in the `time_series` CTE.
	// 2. `startTime` (original) for the `created_at >= toDateTime(?)` in the inner subquery.
	// 3. `channelID` for the `channel_id = ?` in the inner subquery.
	// 4. `emotesNames` for the `emote IN (?)` in the inner subquery.
	rows, err := c.client.Query(ctx, query, startBucketTime, startTime, channelID, emotesNames)
	if err != nil {
		return nil, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	// Process results
	result := make(map[string][]model.EmoteRange)

	for rows.Next() {
		var currentChannelID, emote string // Scan into separate variables to avoid modifying the outer channelID
		var timestamp time.Time
		var count uint64

		if err := rows.Scan(&timestamp, &currentChannelID, &emote, &count); err != nil {
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
) ([]model.EmoteUsage, uint64, error) {
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
			channel_id,
			emote,
			user_id,
			created_at
		FROM channels_emotes_usages
		WHERE channel_id = ? AND emote = ?
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := c.client.Query(ctx, query, input.ChannelID, input.EmoteName, perPage, offset)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var result []model.EmoteUsage
	for rows.Next() {
		var usage model.EmoteUsage
		err := rows.Scan(&usage.ChannelID, &usage.Emote, &usage.UserID, &usage.CreatedAt)
		if err != nil {
			return nil, 0, fmt.Errorf("scan failed: %w", err)
		}

		result = append(result, usage)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	totalQuery := `
		SELECT COUNT(*) FROM channels_emotes_usages
		WHERE channel_id = ? AND emote = ?
`

	var total uint64
	err = c.client.QueryRow(ctx, totalQuery, input.ChannelID, input.EmoteName).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("total query error: %w", err)
	}

	return result, total, nil
}

func (c *Clickhouse) GetChannelUsageTopUsers(
	ctx context.Context,
	input channelsemotesusagesrepository.EmotesUsersTopOrHistoryInput,
) ([]model.EmoteUsageTopUser, uint64, error) {
	var (
		page    = input.Page
		perPage = input.PerPage
	)

	if perPage == 0 || perPage > 1000 {
		perPage = 20
	}

	offset := page * perPage

	queryBuilder := sq.
		Select(
			"channel_id",
			"user_id",
			"COUNT(*) AS count",
		).
		From("channels_emotes_usages").
		Where(squirrel.Eq{"channel_id": input.ChannelID}).
		Where(squirrel.Eq{"emote": input.EmoteName}).
		GroupBy("channel_id", "user_id").
		OrderBy("count DESC").
		Limit(uint64(perPage)).
		Offset(uint64(offset))

	if input.EmoteName != "" {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"emote": input.EmoteName})
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return nil, 0, fmt.Errorf("build query error: %w", err)
	}

	rows, err := c.client.Query(ctx, query, args...)
	if err != nil {
		return nil, 0, fmt.Errorf("query execution failed: %w", err)
	}
	defer rows.Close()

	var result []model.EmoteUsageTopUser
	for rows.Next() {
		var usage model.EmoteUsageTopUser
		err := rows.Scan(&usage.ChannelID, &usage.UserID, &usage.Count)
		if err != nil {
			return nil, 0, fmt.Errorf("scan failed: %w", err)
		}

		result = append(result, usage)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, fmt.Errorf("rows error: %w", err)
	}

	totalQuery := `
		SELECT COUNT(DISTINCT user_id) FROM channels_emotes_usages
		WHERE channel_id = ? AND emote = ?
`

	var total uint64
	err = c.client.QueryRow(ctx, totalQuery, input.ChannelID, input.EmoteName).Scan(&total)
	if err != nil {
		return nil, 0, fmt.Errorf("total query error: %w", err)
	}

	return result, total, nil
}

func (c *Clickhouse) DeleteRowsByChannelID(ctx context.Context, channelID string) error {
	query := `
		DELETE FROM channels_emotes_usages
		WHERE channel_id = ?
	`

	err := c.client.Exec(ctx, query, channelID)
	return err
}
