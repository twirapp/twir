package pgx

import (
	"context"
	"errors"
	"fmt"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	entity "github.com/twirapp/twir/libs/entities/song_request_overlay_settings"
	"github.com/twirapp/twir/libs/repositories/song_request_overlay_settings"
)

type row struct {
	ID                    uuid.UUID
	ChannelID             string
	Style                 entity.Style
	AccentColor           string
	TickerBackgroundColor string
	TickerTextColor       string
	TickerSpeed           int
	HideOnPause           bool
	CreatedAt             time.Time
	UpdatedAt             time.Time
}

func (r row) toEntity() entity.SongRequestOverlaySettings {
	return entity.SongRequestOverlaySettings{
		ID:                    r.ID,
		ChannelID:             r.ChannelID,
		Style:                 r.Style,
		AccentColor:           r.AccentColor,
		TickerBackgroundColor: r.TickerBackgroundColor,
		TickerTextColor:       r.TickerTextColor,
		TickerSpeed:           r.TickerSpeed,
		HideOnPause:           r.HideOnPause,
		CreatedAt:             r.CreatedAt,
		UpdatedAt:             r.UpdatedAt,
	}
}

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return &Pgx{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

var _ song_request_overlay_settings.Repository = (*Pgx)(nil)

func scanRow(rows pgx.Rows) (entity.SongRequestOverlaySettings, error) {
	data, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[row])
	if err != nil {
		return entity.Nil, err
	}

	return data.toEntity(), nil
}

func (r *Pgx) GetByChannelID(
	ctx context.Context,
	channelID string,
) (entity.SongRequestOverlaySettings, error) {
	query := `SELECT
	id,
	channel_id,
	style,
	accent_color,
	ticker_background_color,
	ticker_text_color,
	ticker_speed,
	hide_on_pause,
	created_at,
	updated_at
FROM channels_song_requests_overlay_settings
WHERE channel_id = @channel_id`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(
		ctx,
		query,
		pgx.NamedArgs{"channel_id": channelID},
	)
	if err != nil {
		return entity.Nil, fmt.Errorf("query song request overlay settings: %w", err)
	}

	settings, err := scanRow(rows)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return entity.Nil, song_request_overlay_settings.ErrNotFound
		}
		return entity.Nil, fmt.Errorf("collect song request overlay settings: %w", err)
	}

	return settings, nil
}

func (r *Pgx) Upsert(
	ctx context.Context,
	input song_request_overlay_settings.UpsertInput,
) (entity.SongRequestOverlaySettings, error) {
	query := `INSERT INTO channels_song_requests_overlay_settings (
	channel_id,
	style,
	accent_color,
	ticker_background_color,
	ticker_text_color,
	ticker_speed,
	hide_on_pause
) VALUES (
	@channel_id,
	@style,
	@accent_color,
	@ticker_background_color,
	@ticker_text_color,
	@ticker_speed,
	@hide_on_pause
)
ON CONFLICT (channel_id) DO UPDATE SET
	style = EXCLUDED.style,
	accent_color = EXCLUDED.accent_color,
	ticker_background_color = EXCLUDED.ticker_background_color,
	ticker_text_color = EXCLUDED.ticker_text_color,
	ticker_speed = EXCLUDED.ticker_speed,
	hide_on_pause = EXCLUDED.hide_on_pause,
	updated_at = now()
RETURNING
	id,
	channel_id,
	style,
	accent_color,
	ticker_background_color,
	ticker_text_color,
	ticker_speed,
	hide_on_pause,
	created_at,
	updated_at`

	rows, err := r.getter.DefaultTrOrDB(ctx, r.pool).Query(
		ctx,
		query,
		pgx.NamedArgs{
			"channel_id":              input.ChannelID,
			"style":                   input.Style,
			"accent_color":            input.AccentColor,
			"ticker_background_color": input.TickerBackgroundColor,
			"ticker_text_color":       input.TickerTextColor,
			"ticker_speed":            input.TickerSpeed,
			"hide_on_pause":           input.HideOnPause,
		},
	)
	if err != nil {
		return entity.Nil, fmt.Errorf("upsert song request overlay settings: %w", err)
	}

	settings, err := scanRow(rows)
	if err != nil {
		return entity.Nil, fmt.Errorf("collect upserted song request overlay settings: %w", err)
	}

	return settings, nil
}
