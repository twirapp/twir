package pgx

import (
	"context"
	"errors"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
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

func (p *Pgx) GetByChannelID(ctx context.Context, channelID string) (model.StreamStatsOverlay, error) {
	const query = `
SELECT
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
FROM channels_overlays_stream_stats
WHERE channel_id = $1
LIMIT 1;
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	overlay := model.StreamStatsOverlay{}
	err := conn.QueryRow(ctx, query, channelID).Scan(
		&overlay.ID,
		&overlay.ChannelID,
		&overlay.Design,
		&overlay.Variant,
		&overlay.ViewersEnabled,
		&overlay.ViewersMode,
		&overlay.PlatformIconsEnabled,
		&overlay.MessagesEnabled,
		&overlay.UptimeEnabled,
		&overlay.SubscribersEnabled,
		&overlay.FollowersEnabled,
		&overlay.ViewersColor,
		&overlay.MessagesColor,
		&overlay.UptimeColor,
		&overlay.SubscribersColor,
		&overlay.FollowersColor,
		&overlay.CounterOrder,
		&overlay.CustomHTMLEnabled,
		&overlay.CustomHTML,
		&overlay.CustomCSS,
		&overlay.CreatedAt,
		&overlay.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, overlays_stream_stats.ErrNotFound
		}

		return model.Nil, fmt.Errorf("stream stats overlay get by channel ID: %w", err)
	}

	return overlay, nil
}

func (p *Pgx) Create(ctx context.Context, input overlays_stream_stats.CreateInput) (model.StreamStatsOverlay, error) {
	const query = `
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
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, $16, $17, $18, $19);
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(
		ctx,
		query,
		input.ChannelID,
		input.Design,
		input.Variant,
		input.ViewersEnabled,
		input.ViewersMode,
		input.PlatformIconsEnabled,
		input.MessagesEnabled,
		input.UptimeEnabled,
		input.SubscribersEnabled,
		input.FollowersEnabled,
		input.ViewersColor,
		input.MessagesColor,
		input.UptimeColor,
		input.SubscribersColor,
		input.FollowersColor,
		input.CounterOrder,
		input.CustomHTMLEnabled,
		input.CustomHTML,
		input.CustomCSS,
	)
	if err != nil {
		return model.Nil, fmt.Errorf("stream stats overlay create: %w", err)
	}

	return p.GetByChannelID(ctx, input.ChannelID)
}

func (p *Pgx) Update(ctx context.Context, channelID string, input overlays_stream_stats.UpdateInput) (model.StreamStatsOverlay, error) {
	const query = `
UPDATE channels_overlays_stream_stats
SET
	design = $1,
	variant = $2,
	viewers_enabled = $3,
	viewers_mode = $4,
	platform_icons_enabled = $5,
	messages_enabled = $6,
	uptime_enabled = $7,
	subscribers_enabled = $8,
	followers_enabled = $9,
	viewers_color = $10,
	messages_color = $11,
	uptime_color = $12,
	subscribers_color = $13,
	followers_color = $14,
	counter_order = $15,
	custom_html_enabled = $16,
	custom_html = $17,
	custom_css = $18,
	updated_at = now()
WHERE channel_id = $19
RETURNING channel_id;
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	var updatedChannelID string
	err := conn.QueryRow(
		ctx,
		query,
		input.Design,
		input.Variant,
		input.ViewersEnabled,
		input.ViewersMode,
		input.PlatformIconsEnabled,
		input.MessagesEnabled,
		input.UptimeEnabled,
		input.SubscribersEnabled,
		input.FollowersEnabled,
		input.ViewersColor,
		input.MessagesColor,
		input.UptimeColor,
		input.SubscribersColor,
		input.FollowersColor,
		input.CounterOrder,
		input.CustomHTMLEnabled,
		input.CustomHTML,
		input.CustomCSS,
		channelID,
	).Scan(&updatedChannelID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, overlays_stream_stats.ErrNotFound
		}

		return model.Nil, fmt.Errorf("stream stats overlay update: %w", err)
	}

	return p.GetByChannelID(ctx, updatedChannelID)
}
