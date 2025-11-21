package pgx

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/goccy/go-json"
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
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (p *Pgx) GetByChannelID(ctx context.Context, channelID string) (
	brbmodel.BeRightBackOverlay,
	error,
) {
	query := `
SELECT id, channel_id, created_at, updated_at, data
FROM channels_overlays_be_right_back
WHERE channel_id = $1
LIMIT 1;
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	row := conn.QueryRow(ctx, query, channelID)
	overlay := brbmodel.BeRightBackOverlay{}

	var brbOverlaySettings []byte

	err := row.Scan(
		&overlay.ID,
		&overlay.ChannelID,
		&overlay.CreatedAt,
		&overlay.UpdatedAt,
		&brbOverlaySettings,
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

	if len(brbOverlaySettings) > 0 {
		if err := json.Unmarshal(brbOverlaySettings, &overlay.Settings); err != nil {
			return brbmodel.BeRightBackOverlay{}, err
		}
	}

	return overlay, nil
}

func (p *Pgx) Create(
	ctx context.Context,
	input overlays_be_right_back.CreateInput,
) (brbmodel.BeRightBackOverlay, error) {
	query := `
INSERT INTO channels_overlays_be_right_back (channel_id, data)
VALUES ($1, $2);
`

	settingsBytes, err := json.Marshal(input.Settings)
	if err != nil {
		return brbmodel.BeRightBackOverlay{}, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err = conn.Exec(ctx, query, input.ChannelID, string(settingsBytes))
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
SET data = $1, updated_at = now()
WHERE channel_id = $2
RETURNING channel_id
	`

	settingsBytes, err := json.Marshal(input.Settings)
	if err != nil {
		return brbmodel.BeRightBackOverlay{}, err
	}

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	row := conn.QueryRow(ctx, query, string(settingsBytes), channelID)
	var channelId string

	err = row.Scan(&channelId)
	if err != nil {
		return brbmodel.BeRightBackOverlay{}, fmt.Errorf("be right back overlay update: %w", err)
	}

	return p.GetByChannelID(ctx, channelId)
}
