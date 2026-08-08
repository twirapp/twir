package pgx

import (
	"context"
	"errors"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/overlays_stream_stats"
	"github.com/twirapp/twir/libs/repositories/overlays_stream_stats/model"
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

var _ overlays_stream_stats.Repository = (*Pgx)(nil)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

const columns = `
	id,
	channel_id,
	design,
	variant,
	viewers_enabled,
	viewers_mode,
	platform_icons_enabled,
	messages_enabled,
	uptime_enabled,
	subscribers_enabled,
	followers_enabled,
	viewers_color,
	messages_color,
	uptime_color,
	subscribers_color,
	followers_color,
	counter_order,
	custom_html_enabled,
	custom_html,
	custom_css,
	created_at,
	updated_at
`

func (p *Pgx) GetByChannelID(
	ctx context.Context,
	channelID uuid.UUID,
) (model.StreamStatsOverlay, error) {
	query := `
SELECT ` + columns + `
FROM channels_overlays_stream_stats
WHERE channel_id = @channelID
LIMIT 1;
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, pgx.NamedArgs{"channelID": channelID})
	if err != nil {
		return model.Nil, fmt.Errorf("stream stats overlay get by channel ID: %w", err)
	}

	overlay, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.StreamStatsOverlay])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, overlays_stream_stats.ErrNotFound
		}

		return model.Nil, fmt.Errorf("stream stats overlay get by channel ID: %w", err)
	}

	return overlay, nil
}

func (p *Pgx) Create(
	ctx context.Context,
	input overlays_stream_stats.CreateInput,
) (model.StreamStatsOverlay, error) {
	query := `
INSERT INTO channels_overlays_stream_stats (
	channel_id,
	design,
	variant,
	viewers_enabled,
	viewers_mode,
	platform_icons_enabled,
	messages_enabled,
	uptime_enabled,
	subscribers_enabled,
	followers_enabled,
	viewers_color,
	messages_color,
	uptime_color,
	subscribers_color,
	followers_color,
	counter_order,
	custom_html_enabled,
	custom_html,
	custom_css
)
VALUES (
	@channelID,
	@design,
	@variant,
	@viewersEnabled,
	@viewersMode,
	@platformIconsEnabled,
	@messagesEnabled,
	@uptimeEnabled,
	@subscribersEnabled,
	@followersEnabled,
	@viewersColor,
	@messagesColor,
	@uptimeColor,
	@subscribersColor,
	@followersColor,
	@counterOrder,
	@customHTMLEnabled,
	@customHTML,
	@customCSS
)
RETURNING ` + columns + `;
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(
		ctx,
		query,
		pgx.NamedArgs{
			"channelID":            input.ChannelID,
			"design":               input.Design,
			"variant":              input.Variant,
			"viewersEnabled":       input.ViewersEnabled,
			"viewersMode":          input.ViewersMode,
			"platformIconsEnabled": input.PlatformIconsEnabled,
			"messagesEnabled":      input.MessagesEnabled,
			"uptimeEnabled":        input.UptimeEnabled,
			"subscribersEnabled":   input.SubscribersEnabled,
			"followersEnabled":     input.FollowersEnabled,
			"viewersColor":         input.ViewersColor,
			"messagesColor":        input.MessagesColor,
			"uptimeColor":          input.UptimeColor,
			"subscribersColor":     input.SubscribersColor,
			"followersColor":       input.FollowersColor,
			"counterOrder":         input.CounterOrder,
			"customHTMLEnabled":    input.CustomHTMLEnabled,
			"customHTML":           input.CustomHTML,
			"customCSS":            input.CustomCSS,
		},
	)
	if err != nil {
		return model.Nil, fmt.Errorf("stream stats overlay create: %w", err)
	}

	created, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.StreamStatsOverlay])
	if err != nil {
		return model.Nil, fmt.Errorf("stream stats overlay create: %w", err)
	}

	return created, nil
}

func (p *Pgx) Update(
	ctx context.Context,
	channelID uuid.UUID,
	input overlays_stream_stats.UpdateInput,
) (model.StreamStatsOverlay, error) {
	query := `
UPDATE channels_overlays_stream_stats
SET
	design = @design,
	variant = @variant,
	viewers_enabled = @viewersEnabled,
	viewers_mode = @viewersMode,
	platform_icons_enabled = @platformIconsEnabled,
	messages_enabled = @messagesEnabled,
	uptime_enabled = @uptimeEnabled,
	subscribers_enabled = @subscribersEnabled,
	followers_enabled = @followersEnabled,
	viewers_color = @viewersColor,
	messages_color = @messagesColor,
	uptime_color = @uptimeColor,
	subscribers_color = @subscribersColor,
	followers_color = @followersColor,
	counter_order = @counterOrder,
	custom_html_enabled = @customHTMLEnabled,
	custom_html = @customHTML,
	custom_css = @customCSS,
	updated_at = now()
WHERE channel_id = @channelID
RETURNING ` + columns + `;
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(
		ctx,
		query,
		pgx.NamedArgs{
			"channelID":            channelID,
			"design":               input.Design,
			"variant":              input.Variant,
			"viewersEnabled":       input.ViewersEnabled,
			"viewersMode":          input.ViewersMode,
			"platformIconsEnabled": input.PlatformIconsEnabled,
			"messagesEnabled":      input.MessagesEnabled,
			"uptimeEnabled":        input.UptimeEnabled,
			"subscribersEnabled":   input.SubscribersEnabled,
			"followersEnabled":     input.FollowersEnabled,
			"viewersColor":         input.ViewersColor,
			"messagesColor":        input.MessagesColor,
			"uptimeColor":          input.UptimeColor,
			"subscribersColor":     input.SubscribersColor,
			"followersColor":       input.FollowersColor,
			"counterOrder":         input.CounterOrder,
			"customHTMLEnabled":    input.CustomHTMLEnabled,
			"customHTML":           input.CustomHTML,
			"customCSS":            input.CustomCSS,
		},
	)
	if err != nil {
		return model.Nil, fmt.Errorf("stream stats overlay update: %w", err)
	}

	updated, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.StreamStatsOverlay])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, overlays_stream_stats.ErrNotFound
		}

		return model.Nil, fmt.Errorf("stream stats overlay update: %w", err)
	}

	return updated, nil
}
