package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/entities/requested_song"
	"github.com/twirapp/twir/libs/repositories/requested_songs"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{pool: opts.PgxPool}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ requested_songs.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

type scanModel struct {
	ID                   uuid.UUID  `db:"id"`
	ChannelID            uuid.UUID  `db:"channelId"`
	OrderedByID          uuid.UUID  `db:"orderedById"`
	OrderedByName        string     `db:"orderedByName"`
	OrderedByDisplayName *string    `db:"orderedByDisplayName"`
	VideoID              string     `db:"videoId"`
	Title                string     `db:"title"`
	Duration             int32      `db:"duration"`
	QueuePosition        int        `db:"queuePosition"`
	SongLink             *string    `db:"songLink"`
	CreatedAt            time.Time  `db:"createdAt"`
	DeletedAt            *time.Time `db:"deletedAt"`
}

func (s scanModel) toEntity() requested_song.RequestedSong {
	return requested_song.RequestedSong{
		ID:                   s.ID,
		ChannelID:            s.ChannelID,
		OrderedByID:          s.OrderedByID,
		OrderedByName:        s.OrderedByName,
		OrderedByDisplayName: s.OrderedByDisplayName,
		VideoID:              s.VideoID,
		Title:                s.Title,
		Duration:             s.Duration,
		QueuePosition:        s.QueuePosition,
		SongLink:             s.SongLink,
		CreatedAt:            s.CreatedAt,
		DeletedAt:            s.DeletedAt,
	}
}

const selectFields = `id,
	"channelId",
	"orderedById",
	"orderedByName",
	"orderedByDisplayName",
	"videoId",
	title,
	duration,
	"queuePosition",
	"songLink",
	"createdAt",
	"deletedAt"`

func (c *Pgx) GetByVideoID(
	ctx context.Context,
	channelID string,
	videoID string,
) (requested_song.RequestedSong, error) {
	query := `SELECT ` + selectFields + `
FROM channels_requested_songs
WHERE "videoId" = @video_id
	AND "channelId" = CAST(@channel_id AS uuid)
	AND "deletedAt" IS NULL
LIMIT 1`

	rows, err := c.pool.Query(
		ctx,
		query,
		pgx.NamedArgs{"channel_id": channelID, "video_id": videoID},
	)
	if err != nil {
		return requested_song.Nil, fmt.Errorf("query requested song: %w", err)
	}

	song, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return requested_song.Nil, requested_songs.ErrNotFound
		}

		return requested_song.Nil, fmt.Errorf("collect requested song: %w", err)
	}

	return song.toEntity(), nil
}

func (c *Pgx) GetQueue(
	ctx context.Context,
	channelID string,
) ([]requested_song.RequestedSong, error) {
	query := `SELECT ` + selectFields + `
FROM channels_requested_songs
WHERE "channelId" = CAST(@channel_id AS uuid)
	AND "deletedAt" IS NULL
ORDER BY "queuePosition" ASC`

	rows, err := c.pool.Query(
		ctx,
		query,
		pgx.NamedArgs{"channel_id": channelID},
	)
	if err != nil {
		return nil, fmt.Errorf("query requested songs queue: %w", err)
	}

	songs, err := pgx.CollectRows(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return nil, fmt.Errorf("collect requested songs queue: %w", err)
	}

	result := make([]requested_song.RequestedSong, 0, len(songs))
	for _, song := range songs {
		result = append(result, song.toEntity())
	}

	return result, nil
}

func (c *Pgx) CountByChannelID(
	ctx context.Context,
	channelID string,
	createdAfter time.Time,
) (int64, error) {
	query := `SELECT count(*)
FROM channels_requested_songs
WHERE "channelId" = CAST(@channel_id AS uuid)
	AND "createdAt" >= @created_after`

	var count int64
	if err := c.pool.QueryRow(
		ctx,
		query,
		pgx.NamedArgs{"channel_id": channelID, "created_after": createdAfter},
	).Scan(&count); err != nil {
		return 0, fmt.Errorf("count requested songs: %w", err)
	}

	return count, nil
}

func (c *Pgx) SoftDeleteByVideoID(ctx context.Context, channelID string, videoID string) error {
	query := `UPDATE channels_requested_songs
SET "deletedAt" = now()
WHERE "videoId" = @video_id
	AND "channelId" = CAST(@channel_id AS uuid)
	AND "deletedAt" IS NULL`

	cmd, err := c.pool.Exec(
		ctx,
		query,
		pgx.NamedArgs{"channel_id": channelID, "video_id": videoID},
	)
	if err != nil {
		return fmt.Errorf("soft delete requested song: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return requested_songs.ErrNotFound
	}

	return nil
}

func (c *Pgx) SoftDeleteAll(ctx context.Context, channelID string) error {
	query := `UPDATE channels_requested_songs
SET "deletedAt" = now()
WHERE "channelId" = CAST(@channel_id AS uuid)
	AND "deletedAt" IS NULL`

	if _, err := c.pool.Exec(
		ctx,
		query,
		pgx.NamedArgs{"channel_id": channelID},
	); err != nil {
		return fmt.Errorf("soft delete all requested songs: %w", err)
	}

	return nil
}

func (c *Pgx) UpdateQueuePositions(
	ctx context.Context,
	channelID string,
	videoIDs []string,
) error {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	query := `UPDATE channels_requested_songs
SET "queuePosition" = @position
WHERE "videoId" = @video_id
	AND "channelId" = CAST(@channel_id AS uuid)
	AND "deletedAt" IS NULL`

	for position, videoID := range videoIDs {
		cmd, err := tx.Exec(
			ctx,
			query,
			pgx.NamedArgs{
				"channel_id": channelID,
				"video_id":   videoID,
				"position":   position,
			},
		)
		if err != nil {
			return fmt.Errorf("update queue position for song %s: %w", videoID, err)
		}

		if cmd.RowsAffected() == 0 {
			return fmt.Errorf("song %s: %w", videoID, requested_songs.ErrNotFound)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}
