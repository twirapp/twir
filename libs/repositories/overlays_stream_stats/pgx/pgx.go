package pgx

import (
	"context"
	"errors"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	entity "github.com/twirapp/twir/libs/entities/overlays_stream_stats"
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
) (entity.StreamStatsOverlay, error) {
	query := `
SELECT ` + columns + `
FROM channels_overlays_stream_stats
WHERE channel_id = @channelID
LIMIT 1;
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, pgx.NamedArgs{"channelID": channelID})
	if err != nil {
		return entity.Nil, fmt.Errorf("stream stats overlay get by channel ID: %w", err)
	}

	overlay, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.StreamStatsOverlay])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Nil, overlays_stream_stats.ErrNotFound
		}

		return entity.Nil, fmt.Errorf("stream stats overlay get by channel ID: %w", err)
	}

	return mapModelToEntity(overlay), nil
}

func (p *Pgx) Create(
	ctx context.Context,
	input overlays_stream_stats.CreateInput,
) (entity.StreamStatsOverlay, error) {
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
		return entity.Nil, fmt.Errorf("stream stats overlay create: %w", err)
	}

	created, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.StreamStatsOverlay])
	if err != nil {
		return entity.Nil, fmt.Errorf("stream stats overlay create: %w", err)
	}

	return mapModelToEntity(created), nil
}

func (p *Pgx) Update(
	ctx context.Context,
	channelID uuid.UUID,
	input overlays_stream_stats.UpdateInput,
) (entity.StreamStatsOverlay, error) {
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
		return entity.Nil, fmt.Errorf("stream stats overlay update: %w", err)
	}

	updated, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.StreamStatsOverlay])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Nil, overlays_stream_stats.ErrNotFound
		}

		return entity.Nil, fmt.Errorf("stream stats overlay update: %w", err)
	}

	return mapModelToEntity(updated), nil
}

func mapModelToEntity(m model.StreamStatsOverlay) entity.StreamStatsOverlay {
	counterOrder := make([]entity.StreamStatsOverlayCounter, 0, len(m.CounterOrder))
	for _, counter := range m.CounterOrder {
		counterOrder = append(counterOrder, entity.StreamStatsOverlayCounter(counter))
	}

	return entity.StreamStatsOverlay{
		ID:                   m.ID,
		ChannelID:            m.ChannelID.String(),
		Design:               entity.StreamStatsOverlayDesign(m.Design),
		Variant:              entity.StreamStatsOverlayVariant(m.Variant),
		ViewersEnabled:       m.ViewersEnabled,
		ViewersMode:          entity.StreamStatsOverlayViewersMode(m.ViewersMode),
		PlatformIconsEnabled: m.PlatformIconsEnabled,
		MessagesEnabled:      m.MessagesEnabled,
		UptimeEnabled:        m.UptimeEnabled,
		SubscribersEnabled:   m.SubscribersEnabled,
		FollowersEnabled:     m.FollowersEnabled,
		ViewersColor:         m.ViewersColor,
		MessagesColor:        m.MessagesColor,
		UptimeColor:          m.UptimeColor,
		SubscribersColor:     m.SubscribersColor,
		FollowersColor:       m.FollowersColor,
		CounterOrder:         counterOrder,
		CustomHTMLEnabled:    m.CustomHTMLEnabled,
		CustomHTML:           m.CustomHTML,
		CustomCSS:            m.CustomCSS,
		CreatedAt:            m.CreatedAt,
		UpdatedAt:            m.UpdatedAt,
	}
}
