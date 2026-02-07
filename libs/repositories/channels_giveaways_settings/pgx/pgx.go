package pgx

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/entities/channels_giveaways_settings"
)

type model struct {
	ID            uuid.UUID `db:"id"`
	ChannelID     string    `db:"channel_id"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
	WinnerMessage string    `db:"winner_message"`

	isNil bool
}

func (m model) IsNil() bool {
	return m.isNil
}

var Nil = model{
	isNil: true,
}

func modelToEntity(model model) channels_giveaways_settings.Settings {
	if model.IsNil() {
		return channels_giveaways_settings.Nil
	}

	return channels_giveaways_settings.Settings{
		ID:            model.ID,
		ChannelID:     model.ChannelID,
		CreatedAt:     model.CreatedAt,
		UpdatedAt:     model.UpdatedAt,
		WinnerMessage: model.WinnerMessage,
	}
}

type Repository struct {
	db *pgxpool.Pool
}

func New(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetByChannelID(
	ctx context.Context,
	channelID string,
) (channels_giveaways_settings.Settings, error) {
	const query = `
		INSERT INTO channels_giveaways_settings (channel_id)
		VALUES ($1)
		ON CONFLICT (channel_id) DO UPDATE SET channel_id = EXCLUDED.channel_id
		RETURNING id, channel_id, created_at, updated_at, winner_message
	`

	var model model
	err := r.db.QueryRow(ctx, query, channelID).Scan(
		&model.ID,
		&model.ChannelID,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.WinnerMessage,
	)
	if err != nil {
		if err == pgx.ErrNoRows {
			return channels_giveaways_settings.Nil, nil
		}
		return channels_giveaways_settings.Nil, err
	}

	return modelToEntity(model), nil
}

func (r *Repository) Update(
	ctx context.Context,
	channelID string,
	settings channels_giveaways_settings.Settings,
) (channels_giveaways_settings.Settings, error) {
	const query = `
		UPDATE channels_giveaways_settings
		SET winner_message = $2, updated_at = NOW()
		WHERE channel_id = $1
		RETURNING id, channel_id, created_at, updated_at, winner_message
	`

	var model model
	err := r.db.QueryRow(ctx, query, channelID, settings.WinnerMessage).Scan(
		&model.ID,
		&model.ChannelID,
		&model.CreatedAt,
		&model.UpdatedAt,
		&model.WinnerMessage,
	)
	if err != nil {
		return channels_giveaways_settings.Nil, err
	}

	return modelToEntity(model), nil
}
