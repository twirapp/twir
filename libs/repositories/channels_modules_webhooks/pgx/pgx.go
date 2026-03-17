package pgx

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/entities/webhook_notifications"
	channelsmoduleswebhooks "github.com/twirapp/twir/libs/repositories/channels_modules_webhooks"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool:   opts.PgxPool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var (
	_  channelsmoduleswebhooks.Repository = (*Pgx)(nil)
	sq                                    = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

var selectColumns = []string{
	"id",
	"channel_id",
	"enabled",
	"github_issues_enabled",
	"github_pull_requests_enabled",
	"github_commits_enabled",
	"created_at",
	"updated_at",
	"github_issues_online_enabled",
	"github_issues_offline_enabled",
	"github_pull_requests_online_enabled",
	"github_pull_requests_offline_enabled",
	"github_commits_online_enabled",
	"github_commits_offline_enabled",
}

var selectColumnsStr = strings.Join(selectColumns, ", ")

func mapModelToEntity(m channelsmoduleswebhooks.Settings) webhook_notifications.Settings {
	if m.IsNil() {
		return webhook_notifications.Nil
	}

	return webhook_notifications.Settings{
		ID:                               m.ID,
		ChannelID:                        m.ChannelID,
		Enabled:                          m.Enabled,
		GithubIssuesEnabled:              m.GithubIssuesEnabled,
		GithubPullRequestsEnabled:        m.GithubPullRequestsEnabled,
		GithubCommitsEnabled:             m.GithubCommitsEnabled,
		GithubIssuesOnlineEnabled:        m.GithubIssuesOnlineEnabled,
		GithubIssuesOfflineEnabled:       m.GithubIssuesOfflineEnabled,
		GithubPullRequestsOnlineEnabled:  m.GithubPullRequestsOnlineEnabled,
		GithubPullRequestsOfflineEnabled: m.GithubPullRequestsOfflineEnabled,
		GithubCommitsOnlineEnabled:       m.GithubCommitsOnlineEnabled,
		GithubCommitsOfflineEnabled:      m.GithubCommitsOfflineEnabled,
		CreatedAt:                        m.CreatedAt,
		UpdatedAt:                        m.UpdatedAt,
	}
}

func (c *Pgx) GetByID(ctx context.Context, id string) (webhook_notifications.Settings, error) {
	query := `
SELECT ` + selectColumnsStr + `
FROM channels_modules_webhooks
WHERE id = $1
LIMIT 1;
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return webhook_notifications.Nil, fmt.Errorf("query err: %w", err)
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[channelsmoduleswebhooks.Settings],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return webhook_notifications.Nil, channelsmoduleswebhooks.ErrSettingsNotFound
		}
		return webhook_notifications.Nil, fmt.Errorf("collect err: %w", err)
	}

	return mapModelToEntity(result), nil
}

func (c *Pgx) GetByChannelID(
	ctx context.Context,
	channelID string,
) (webhook_notifications.Settings, error) {
	query := `
SELECT ` + selectColumnsStr + `
FROM channels_modules_webhooks
WHERE channel_id = $1
LIMIT 1;
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return webhook_notifications.Nil, fmt.Errorf("query err: %w", err)
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[channelsmoduleswebhooks.Settings],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return webhook_notifications.Nil, channelsmoduleswebhooks.ErrSettingsNotFound
		}
		return webhook_notifications.Nil, fmt.Errorf("collect err: %w", err)
	}

	return mapModelToEntity(result), nil
}

func (c *Pgx) GetEnabledChannels(
	ctx context.Context,
	input channelsmoduleswebhooks.GetEnabledChannelsInput,
) ([]string, error) {
	builder := sq.Select("channel_id").
		From("channels_modules_webhooks").
		Where(squirrel.Eq{"enabled": true})

	if input.GithubIssuesEnabled != nil {
		builder = builder.Where(squirrel.Eq{"github_issues_enabled": *input.GithubIssuesEnabled})
	}

	if input.GithubPullRequestsEnabled != nil {
		builder = builder.Where(
			squirrel.Eq{"github_pull_requests_enabled": *input.GithubPullRequestsEnabled},
		)
	}

	if input.GithubCommitsEnabled != nil {
		builder = builder.Where(squirrel.Eq{"github_commits_enabled": *input.GithubCommitsEnabled})
	}

	if input.GithubIssuesOnlineEnabled != nil {
		builder = builder.Where(
			squirrel.Eq{"github_issues_online_enabled": *input.GithubIssuesOnlineEnabled},
		)
	}

	if input.GithubIssuesOfflineEnabled != nil {
		builder = builder.Where(
			squirrel.Eq{"github_issues_offline_enabled": *input.GithubIssuesOfflineEnabled},
		)
	}

	if input.GithubPullRequestsOnlineEnabled != nil {
		builder = builder.Where(
			squirrel.Eq{"github_pull_requests_online_enabled": *input.GithubPullRequestsOnlineEnabled},
		)
	}

	if input.GithubPullRequestsOfflineEnabled != nil {
		builder = builder.Where(
			squirrel.Eq{
				"github_pull_requests_offline_enabled": *input.GithubPullRequestsOfflineEnabled,
			},
		)
	}

	if input.GithubCommitsOnlineEnabled != nil {
		builder = builder.Where(
			squirrel.Eq{"github_commits_online_enabled": *input.GithubCommitsOnlineEnabled},
		)
	}

	if input.GithubCommitsOfflineEnabled != nil {
		builder = builder.Where(
			squirrel.Eq{"github_commits_offline_enabled": *input.GithubCommitsOfflineEnabled},
		)
	}

	builder = builder.OrderBy("channel_id")

	if input.PerPage > 0 {
		page := input.Page
		if page < 0 {
			page = 0
		}
		builder = builder.Limit(uint64(input.PerPage)).Offset(uint64(page * input.PerPage))
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build err: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query err: %w", err)
	}
	defer rows.Close()

	var channels []string
	for rows.Next() {
		var channelID string
		if err := rows.Scan(&channelID); err != nil {
			return nil, fmt.Errorf("scan err: %w", err)
		}
		channels = append(channels, channelID)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("rows err: %w", rows.Err())
	}

	return channels, nil
}

func (c *Pgx) Create(
	ctx context.Context,
	input channelsmoduleswebhooks.CreateInput,
) (webhook_notifications.Settings, error) {
	query := `
INSERT INTO channels_modules_webhooks (
	channel_id,
	enabled,
	github_issues_enabled,
	github_pull_requests_enabled,
	github_commits_enabled,
	github_issues_online_enabled,
	github_issues_offline_enabled,
	github_pull_requests_online_enabled,
	github_pull_requests_offline_enabled,
	github_commits_online_enabled,
	github_commits_offline_enabled
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING ` + selectColumnsStr + `;
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.ChannelID,
		input.Enabled,
		input.GithubIssuesEnabled,
		input.GithubPullRequestsEnabled,
		input.GithubCommitsEnabled,
		input.GithubIssuesOnlineEnabled,
		input.GithubIssuesOfflineEnabled,
		input.GithubPullRequestsOnlineEnabled,
		input.GithubPullRequestsOfflineEnabled,
		input.GithubCommitsOnlineEnabled,
		input.GithubCommitsOfflineEnabled,
	)
	if err != nil {
		return webhook_notifications.Nil, fmt.Errorf("query err: %w", err)
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[channelsmoduleswebhooks.Settings],
	)
	if err != nil {
		return webhook_notifications.Nil, fmt.Errorf("collect err: %w", err)
	}

	return mapModelToEntity(result), nil
}

func (c *Pgx) Update(
	ctx context.Context,
	id uuid.UUID,
	input channelsmoduleswebhooks.UpdateInput,
) (webhook_notifications.Settings, error) {
	builder := sq.Update("channels_modules_webhooks").
		Where(squirrel.Eq{"id": id.String()}).
		Suffix("RETURNING " + selectColumnsStr)

	if input.Enabled != nil {
		builder = builder.Set("enabled", *input.Enabled)
	}

	if input.GithubIssuesEnabled != nil {
		builder = builder.Set("github_issues_enabled", *input.GithubIssuesEnabled)
	}

	if input.GithubPullRequestsEnabled != nil {
		builder = builder.Set("github_pull_requests_enabled", *input.GithubPullRequestsEnabled)
	}

	if input.GithubCommitsEnabled != nil {
		builder = builder.Set("github_commits_enabled", *input.GithubCommitsEnabled)
	}

	if input.GithubIssuesOnlineEnabled != nil {
		builder = builder.Set("github_issues_online_enabled", *input.GithubIssuesOnlineEnabled)
	}

	if input.GithubIssuesOfflineEnabled != nil {
		builder = builder.Set("github_issues_offline_enabled", *input.GithubIssuesOfflineEnabled)
	}

	if input.GithubPullRequestsOnlineEnabled != nil {
		builder = builder.Set(
			"github_pull_requests_online_enabled",
			*input.GithubPullRequestsOnlineEnabled,
		)
	}

	if input.GithubPullRequestsOfflineEnabled != nil {
		builder = builder.Set(
			"github_pull_requests_offline_enabled",
			*input.GithubPullRequestsOfflineEnabled,
		)
	}

	if input.GithubCommitsOnlineEnabled != nil {
		builder = builder.Set("github_commits_online_enabled", *input.GithubCommitsOnlineEnabled)
	}

	if input.GithubCommitsOfflineEnabled != nil {
		builder = builder.Set("github_commits_offline_enabled", *input.GithubCommitsOfflineEnabled)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return webhook_notifications.Nil, fmt.Errorf("build err: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return webhook_notifications.Nil, channelsmoduleswebhooks.ErrSettingsNotFound
		}
		return webhook_notifications.Nil, fmt.Errorf("query err: %w", err)
	}

	result, err := pgx.CollectExactlyOneRow(
		rows,
		pgx.RowToStructByName[channelsmoduleswebhooks.Settings],
	)
	if err != nil {
		return webhook_notifications.Nil, fmt.Errorf("collect err: %w", err)
	}

	return mapModelToEntity(result), nil
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_modules_webhooks
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id.String())
	if errors.Is(err, pgx.ErrNoRows) {
		return channelsmoduleswebhooks.ErrSettingsNotFound
	}
	return err
}
