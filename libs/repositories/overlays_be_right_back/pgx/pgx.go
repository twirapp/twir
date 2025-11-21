package pgx

import (
	"context"
	"errors"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/overlays_be_right_back"
	brbmodel "github.com/twirapp/twir/libs/repositories/overlays_be_right_back/model"
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

var _ overlays_be_right_back.Repository = (*Pgx)(nil)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (p *Pgx) GetByChannelID(ctx context.Context, channelID string) (
	brbmodel.BeRightBackOverlay,
	error,
) {
	query := `
SELECT
	id,
	channel_id,
	created_at,
	updated_at,
	text,
	late_enabled,
	late_text,
	late_display_brb_time,
	background_color,
	font_size,
	font_color,
	font_family
FROM channels_overlays_be_right_back
WHERE channel_id = $1
LIMIT 1;
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	row := conn.QueryRow(ctx, query, channelID)
	overlay := brbmodel.BeRightBackOverlay{
		Settings: &brbmodel.BeRightBackOverlaySettings{},
	}

	err := row.Scan(
		&overlay.ID,
		&overlay.ChannelID,
		&overlay.CreatedAt,
		&overlay.UpdatedAt,
		&overlay.Settings.Text,
		&overlay.Settings.Late.Enabled,
		&overlay.Settings.Late.Text,
		&overlay.Settings.Late.DisplayBrbTime,
		&overlay.Settings.BackgroundColor,
		&overlay.Settings.FontSize,
		&overlay.Settings.FontColor,
		&overlay.Settings.FontFamily,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return brbmodel.BeRightBackOverlay{}, overlays_be_right_back.ErrNotFound
		}
		return brbmodel.BeRightBackOverlay{}, fmt.Errorf(
			"be right back overlay get by channel ID: %w",
			err,
		)
	}

	return overlay, nil
}

func (p *Pgx) Create(
	ctx context.Context,
	input overlays_be_right_back.CreateInput,
) (brbmodel.BeRightBackOverlay, error) {
	query := `
INSERT INTO channels_overlays_be_right_back (
	channel_id,
	text,
	late_enabled,
	late_text,
	late_display_brb_time,
	background_color,
	font_size,
	font_color,
	font_family
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9);
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(
		ctx,
		query,
		input.ChannelID,
		input.Settings.Text,
		input.Settings.Late.Enabled,
		input.Settings.Late.Text,
		input.Settings.Late.DisplayBrbTime,
		input.Settings.BackgroundColor,
		input.Settings.FontSize,
		input.Settings.FontColor,
		input.Settings.FontFamily,
	)
	if err != nil {
		return brbmodel.BeRightBackOverlay{}, fmt.Errorf("be right back overlay create: %w", err)
	}

	return p.GetByChannelID(ctx, input.ChannelID)
}

func (p *Pgx) Update(
	ctx context.Context,
	channelID string,
	input overlays_be_right_back.UpdateInput,
) (brbmodel.BeRightBackOverlay, error) {
	query := `
UPDATE channels_overlays_be_right_back
SET
	text = $1,
	late_enabled = $2,
	late_text = $3,
	late_display_brb_time = $4,
	background_color = $5,
	font_size = $6,
	font_color = $7,
	font_family = $8,
	updated_at = now()
WHERE channel_id = $9
RETURNING channel_id
	`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	row := conn.QueryRow(
		ctx,
		query,
		input.Settings.Text,
		input.Settings.Late.Enabled,
		input.Settings.Late.Text,
		input.Settings.Late.DisplayBrbTime,
		input.Settings.BackgroundColor,
		input.Settings.FontSize,
		input.Settings.FontColor,
		input.Settings.FontFamily,
		channelID,
	)
	var channelId string

	err := row.Scan(&channelId)
	if err != nil {
		return brbmodel.BeRightBackOverlay{}, fmt.Errorf("be right back overlay update: %w", err)
	}

	return p.GetByChannelID(ctx, channelId)
}
