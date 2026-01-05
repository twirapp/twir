package pgx

import (
	"context"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	channels_overlays "github.com/twirapp/twir/libs/repositories/channels_overlays"
	"github.com/twirapp/twir/libs/repositories/channels_overlays/model"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool: opts.PgxPool,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ channels_overlays.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

type overlayRow struct {
	ID        uuid.UUID `json:"id"`
	ChannelID string    `json:"channel_id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Width     int       `json:"width"`
	Height    int       `json:"height"`
	InstaSave bool      `json:"insta_save"`
}

type layerRow struct {
	ID                      uuid.UUID `json:"id"`
	Type                    string    `json:"type"`
	Settings                []byte    `json:"settings"`
	OverlayID               uuid.UUID `json:"overlay_id"`
	PosX                    int       `json:"pos_x"`
	PosY                    int       `json:"pos_y"`
	Width                   int       `json:"width"`
	Height                  int       `json:"height"`
	Rotation                int       `json:"rotation"`
	CreatedAt               time.Time `json:"created_at"`
	UpdatedAt               time.Time `json:"updated_at"`
	PeriodicallyRefetchData bool      `json:"periodically_refetch_data"`
}

func (c *Pgx) getLayers(ctx context.Context, overlayID uuid.UUID) ([]model.OverlayLayer, error) {
	query := `
SELECT id, type, settings, overlay_id, pos_x, pos_y, width, height, rotation, created_at, updated_at, periodically_refetch_data
FROM channels_overlays_layers
WHERE overlay_id = $1
ORDER BY created_at ASC
`
	rows, err := c.pool.Query(ctx, query, overlayID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var layers []model.OverlayLayer
	for rows.Next() {
		var row layerRow
		if err := rows.Scan(
			&row.ID,
			&row.Type,
			&row.Settings,
			&row.OverlayID,
			&row.PosX,
			&row.PosY,
			&row.Width,
			&row.Height,
			&row.Rotation,
			&row.CreatedAt,
			&row.UpdatedAt,
			&row.PeriodicallyRefetchData,
		); err != nil {
			return nil, err
		}

		var settings model.OverlayLayerSettings
		if err := json.Unmarshal(row.Settings, &settings); err != nil {
			return nil, err
		}

		layers = append(
			layers, model.OverlayLayer{
				ID:                      row.ID,
				Type:                    model.OverlayType(row.Type),
				Settings:                settings,
				OverlayID:               row.OverlayID,
				PosX:                    row.PosX,
				PosY:                    row.PosY,
				Width:                   row.Width,
				Height:                  row.Height,
				Rotation:                row.Rotation,
				CreatedAt:               row.CreatedAt,
				UpdatedAt:               row.UpdatedAt,
				PeriodicallyRefetchData: row.PeriodicallyRefetchData,
			},
		)
	}

	return layers, nil
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.Overlay, error) {
	query := `
SELECT id, channel_id, name, created_at, updated_at, width, height, insta_save
FROM channels_overlays
WHERE id = $1
`

	row := c.pool.QueryRow(ctx, query, id)

	var overlay overlayRow
	if err := row.Scan(
		&overlay.ID,
		&overlay.ChannelID,
		&overlay.Name,
		&overlay.CreatedAt,
		&overlay.UpdatedAt,
		&overlay.Width,
		&overlay.Height,
		&overlay.InstaSave,
	); err != nil {
		if err == pgx.ErrNoRows {
			return model.Nil, channels_overlays.ErrNotFound
		}
		return model.Nil, err
	}

	layers, err := c.getLayers(ctx, overlay.ID)
	if err != nil {
		return model.Nil, err
	}

	return model.Overlay{
		ID:        overlay.ID,
		ChannelID: overlay.ChannelID,
		Name:      overlay.Name,
		CreatedAt: overlay.CreatedAt,
		UpdatedAt: overlay.UpdatedAt,
		Width:     overlay.Width,
		Height:    overlay.Height,
		InstaSave: overlay.InstaSave,
		Layers:    layers,
	}, nil
}

func (c *Pgx) GetManyByChannelID(ctx context.Context, channelID string) ([]model.Overlay, error) {
	query := `
SELECT id, channel_id, name, created_at, updated_at, width, height, insta_save
FROM channels_overlays
WHERE channel_id = $1
ORDER BY created_at DESC
`

	rows, err := c.pool.Query(ctx, query, channelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var overlays []model.Overlay
	for rows.Next() {
		var overlay overlayRow
		if err := rows.Scan(
			&overlay.ID,
			&overlay.ChannelID,
			&overlay.Name,
			&overlay.CreatedAt,
			&overlay.UpdatedAt,
			&overlay.Width,
			&overlay.Height,
			&overlay.InstaSave,
		); err != nil {
			return nil, err
		}

		layers, err := c.getLayers(ctx, overlay.ID)
		if err != nil {
			return nil, err
		}

		overlays = append(
			overlays, model.Overlay{
				ID:        overlay.ID,
				ChannelID: overlay.ChannelID,
				Name:      overlay.Name,
				CreatedAt: overlay.CreatedAt,
				UpdatedAt: overlay.UpdatedAt,
				Width:     overlay.Width,
				Height:    overlay.Height,
				InstaSave: overlay.InstaSave,
				Layers:    layers,
			},
		)
	}

	return overlays, nil
}

func (c *Pgx) Create(ctx context.Context, input channels_overlays.CreateInput) (
	model.Overlay,
	error,
) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return model.Nil, err
	}
	defer tx.Rollback(ctx)

	overlayID := uuid.New()
	now := time.Now().UTC()

	overlayQuery := `
INSERT INTO channels_overlays (id, channel_id, name, created_at, updated_at, width, height, insta_save)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
`

	_, err = tx.Exec(
		ctx,
		overlayQuery,
		overlayID,
		input.ChannelID,
		input.Name,
		now,
		now,
		input.Width,
		input.Height,
		input.InstaSave,
	)
	if err != nil {
		return model.Nil, err
	}

	for _, layer := range input.Layers {
		if err := c.insertLayer(ctx, tx, overlayID, layer); err != nil {
			return model.Nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Nil, err
	}

	return c.GetByID(ctx, overlayID)
}

func (c *Pgx) insertLayer(
	ctx context.Context,
	tx pgx.Tx,
	overlayID uuid.UUID,
	layer channels_overlays.CreateLayerInput,
) error {
	layerID := uuid.New()
	now := time.Now().UTC()

	settingsJSON, err := json.Marshal(layer.Settings)
	if err != nil {
		return err
	}

	layerQuery := `
INSERT INTO channels_overlays_layers (id, type, settings, overlay_id, pos_x, pos_y, width, height, rotation, created_at, updated_at, periodically_refetch_data)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
`

	_, err = tx.Exec(
		ctx,
		layerQuery,
		layerID,
		string(layer.Type),
		settingsJSON,
		overlayID,
		layer.PosX,
		layer.PosY,
		layer.Width,
		layer.Height,
		layer.Rotation,
		now,
		now,
		layer.PeriodicallyRefetchData,
	)
	return err
}

func (c *Pgx) Update(ctx context.Context, id uuid.UUID, input channels_overlays.UpdateInput) (
	model.Overlay,
	error,
) {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return model.Nil, err
	}
	defer tx.Rollback(ctx)

	now := time.Now().UTC()

	overlayQuery := `
UPDATE channels_overlays
SET name = $1, updated_at = $2, width = $3, height = $4, insta_save = $5
WHERE id = $6
`

	result, err := tx.Exec(ctx, overlayQuery, input.Name, now, input.Width, input.Height, input.InstaSave, id)
	if err != nil {
		return model.Nil, err
	}

	if result.RowsAffected() == 0 {
		return model.Nil, channels_overlays.ErrNotFound
	}

	deleteLayersQuery := `DELETE FROM channels_overlays_layers WHERE overlay_id = $1`
	_, err = tx.Exec(ctx, deleteLayersQuery, id)
	if err != nil {
		return model.Nil, err
	}

	for _, layer := range input.Layers {
		if err := c.insertLayer(ctx, tx, id, layer); err != nil {
			return model.Nil, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Nil, err
	}

	return c.GetByID(ctx, id)
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM channels_overlays WHERE id = $1`
	result, err := c.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return channels_overlays.ErrNotFound
	}

	return nil
}

func (c *Pgx) UpdateLayer(ctx context.Context, layerId uuid.UUID, input channels_overlays.LayerUpdateInput) (model.OverlayLayer, error) {
	query := `
UPDATE channels_overlays_layers
SET pos_x = COALESCE($1, pos_x),
		pos_y = COALESCE($2, pos_y),
		updated_at = $3,
		width = COALESCE($4, width),
		height = COALESCE($5, height),
		rotation = COALESCE($6, rotation)
WHERE id = $7
RETURNING id, type, settings, overlay_id, pos_x, pos_y, width, height, rotation, created_at, updated_at, periodically_refetch_data
`

	now := time.Now().UTC()
	result, err := c.pool.Exec(
		ctx,
		query,
		input.PosX,
		input.PosY,
		now,
		input.Width,
		input.Height,
		input.Rotation,
		layerId,
	)
	if err != nil {
		return model.OverlayLayer{}, err
	}

	if result.RowsAffected() == 0 {
		return model.OverlayLayer{}, channels_overlays.ErrNotFound
	}

	getQuery := `
SELECT id, type, settings, overlay_id, pos_x, pos_y, width, height, rotation, created_at, updated_at, periodically_refetch_data
FROM channels_overlays_layers
WHERE id = $1
`

	row := c.pool.QueryRow(ctx, getQuery, layerId)

	var layer layerRow
	if err := row.Scan(
		&layer.ID,
		&layer.Type,
		&layer.Settings,
		&layer.OverlayID,
		&layer.PosX,
		&layer.PosY,
		&layer.Width,
		&layer.Height,
		&layer.Rotation,
		&layer.CreatedAt,
		&layer.UpdatedAt,
		&layer.PeriodicallyRefetchData,
	); err != nil {
		return model.OverlayLayer{}, err
	}

	var settings model.OverlayLayerSettings
	if err := json.Unmarshal(layer.Settings, &settings); err != nil {
		return model.OverlayLayer{}, err
	}

	return model.OverlayLayer{
		ID:                      layer.ID,
		Type:                    model.OverlayType(layer.Type),
		Settings:                settings,
		OverlayID:               layer.OverlayID,
		PosX:                    layer.PosX,
		PosY:                    layer.PosY,
		Width:                   layer.Width,
		Height:                  layer.Height,
		Rotation:                layer.Rotation,
		CreatedAt:               layer.CreatedAt,
		UpdatedAt:               layer.UpdatedAt,
		PeriodicallyRefetchData: layer.PeriodicallyRefetchData,
	}, nil
}
