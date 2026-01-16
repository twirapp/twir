package clickhouse

import (
	"context"
	"fmt"

	twirclickhouse "github.com/twirapp/twir/libs/baseapp/clickhouse"
	"github.com/twirapp/twir/libs/repositories/short_links_views"
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

var _ short_links_views.Repository = (*Clickhouse)(nil)

type Clickhouse struct {
	client *twirclickhouse.ClickhouseClient
}

func (c *Clickhouse) Create(ctx context.Context, input short_links_views.CreateInput) error {
	query := `
INSERT INTO short_links_views (short_link_id, user_id, ip, user_agent, country, city, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?);
`

	err := c.client.Exec(
		ctx,
		query,
		input.ShortLinkID,
		input.UserID,
		input.IP,
		input.UserAgent,
		input.Country,
		input.City,
		input.CreatedAt,
	)

	return err
}

func (c *Clickhouse) GetStatistics(
	ctx context.Context,
	input short_links_views.GetStatisticsInput,
) ([]short_links_views.StatisticsPoint, error) {
	var intervalFunc string
	switch input.Interval {
	case "hour":
		intervalFunc = "toStartOfHour(created_at)"
	case "day":
		intervalFunc = "toStartOfDay(created_at)"
	default:
		return nil, fmt.Errorf("invalid interval: %s", input.Interval)
	}

	query := fmt.Sprintf(`
SELECT
	toUnixTimestamp(%s) * 1000 as timestamp,
	count() as count
FROM short_links_views
WHERE short_link_id = ?
	AND created_at >= ?
	AND created_at <= ?
GROUP BY timestamp
ORDER BY timestamp ASC
`, intervalFunc)

	rows, err := c.client.Query(ctx, query, input.ShortLinkID, input.From, input.To)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []short_links_views.StatisticsPoint
	for rows.Next() {
		var point short_links_views.StatisticsPoint
		if err := rows.Scan(&point.Timestamp, &point.Count); err != nil {
			return nil, err
		}
		result = append(result, point)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}

func (c *Clickhouse) GetViews(
	ctx context.Context,
	input short_links_views.GetViewsInput,
) (short_links_views.GetViewsOutput, error) {
	offset := input.Page * input.PerPage

	// Get total count
	countQuery := `
SELECT count() as total
FROM short_links_views
WHERE short_link_id = ?
`
	var total uint64
	if err := c.client.QueryRow(ctx, countQuery, input.ShortLinkID).Scan(&total); err != nil {
		return short_links_views.GetViewsOutput{}, err
	}

	// Get views
	query := `
SELECT
	short_link_id,
	user_id,
	country,
	city,
	created_at
FROM short_links_views
WHERE short_link_id = ?
ORDER BY created_at DESC
LIMIT ? OFFSET ?
`

	rows, err := c.client.Query(ctx, query, input.ShortLinkID, input.PerPage, offset)
	if err != nil {
		return short_links_views.GetViewsOutput{}, err
	}
	defer rows.Close()

	var views []short_links_views.View
	for rows.Next() {
		var view short_links_views.View
		if err := rows.Scan(
			&view.ShortLinkID,
			&view.UserID,
			&view.Country,
			&view.City,
			&view.CreatedAt,
		); err != nil {
			return short_links_views.GetViewsOutput{}, err
		}
		views = append(views, view)
	}

	if rows.Err() != nil {
		return short_links_views.GetViewsOutput{}, rows.Err()
	}

	return short_links_views.GetViewsOutput{
		Views: views,
		Total: int(total),
	}, nil
}

func (c *Clickhouse) GetTopCountries(
	ctx context.Context,
	input short_links_views.GetTopCountriesInput,
) ([]short_links_views.CountryStats, error) {
	query := `
SELECT
	country,
	count() as count
FROM short_links_views
WHERE short_link_id = ?
	AND country IS NOT NULL
	AND country != ''
GROUP BY country
ORDER BY count DESC
LIMIT ?
`

	rows, err := c.client.Query(ctx, query, input.ShortLinkID, input.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var result []short_links_views.CountryStats
	for rows.Next() {
		var stats short_links_views.CountryStats
		if err := rows.Scan(&stats.Country, &stats.Count); err != nil {
			return nil, err
		}
		result = append(result, stats)
	}

	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return result, nil
}
