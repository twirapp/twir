package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/entities/spotify_song_request"
	spotify_song_requests "github.com/twirapp/twir/libs/repositories/spotify_song_requests"
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

var _ spotify_song_requests.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

type scanModel struct {
	ID                     uuid.UUID                   `db:"id"`
	ChannelID              string                      `db:"channel_id"`
	TrackID                string                      `db:"track_id"`
	TrackURI               string                      `db:"track_uri"`
	Title                  string                      `db:"title"`
	Artist                 string                      `db:"artist"`
	Album                  string                      `db:"album"`
	DurationMs             int                         `db:"duration_ms"`
	RequesterUserID        *string                     `db:"requester_user_id"`
	RequesterName          string                      `db:"requester_name"`
	RequesterDisplayName   *string                     `db:"requester_display_name"`
	Source                 string                      `db:"source"`
	QueuePosition          int                         `db:"queue_position"`
	Status                 spotify_song_request.Status `db:"status"`
	QueuedAt               time.Time                   `db:"queued_at"`
	PlayingObservedAt      *time.Time                  `db:"playing_observed_at"`
	PlayedObservedAt       *time.Time                  `db:"played_observed_at"`
	CancelledPendingSkipAt *time.Time                  `db:"cancelled_pending_skip_at"`
	SkippedByTwirAt        *time.Time                  `db:"skipped_by_twir_at"`
	RemovedOrReconciledAt  *time.Time                  `db:"removed_or_reconciled_at"`
	UnknownAt              *time.Time                  `db:"unknown_at"`
	CreatedAt              time.Time                   `db:"created_at"`
	UpdatedAt              time.Time                   `db:"updated_at"`
}

func (s scanModel) toEntity() spotify_song_request.SpotifySongRequest {
	return spotify_song_request.SpotifySongRequest{
		ID:                     s.ID,
		ChannelID:              s.ChannelID,
		TrackID:                s.TrackID,
		TrackURI:               s.TrackURI,
		Title:                  s.Title,
		Artist:                 s.Artist,
		Album:                  s.Album,
		DurationMs:             s.DurationMs,
		RequesterUserID:        s.RequesterUserID,
		RequesterName:          s.RequesterName,
		RequesterDisplayName:   s.RequesterDisplayName,
		Source:                 s.Source,
		QueuePosition:          s.QueuePosition,
		Status:                 s.Status,
		QueuedAt:               s.QueuedAt,
		PlayingObservedAt:      s.PlayingObservedAt,
		PlayedObservedAt:       s.PlayedObservedAt,
		CancelledPendingSkipAt: s.CancelledPendingSkipAt,
		SkippedByTwirAt:        s.SkippedByTwirAt,
		RemovedOrReconciledAt:  s.RemovedOrReconciledAt,
		UnknownAt:              s.UnknownAt,
		CreatedAt:              s.CreatedAt,
		UpdatedAt:              s.UpdatedAt,
	}
}

const selectFields = `id,
	channel_id,
	track_id,
	track_uri,
	title,
	artist,
	album,
	duration_ms,
	requester_user_id,
	requester_name,
	requester_display_name,
	source,
	queue_position,
	status,
	queued_at,
	playing_observed_at,
	played_observed_at,
	cancelled_pending_skip_at,
	skipped_by_twir_at,
	removed_or_reconciled_at,
	unknown_at,
	created_at,
	updated_at`

func (c *Pgx) Create(
	ctx context.Context,
	req spotify_song_request.SpotifySongRequest,
) (spotify_song_request.SpotifySongRequest, error) {
	query := `INSERT INTO spotify_song_requests (
	id,
	channel_id,
	track_id,
	track_uri,
	title,
	artist,
	album,
	duration_ms,
	requester_user_id,
	requester_name,
	requester_display_name,
	source,
	queue_position,
	status,
	queued_at,
	playing_observed_at,
	played_observed_at,
	cancelled_pending_skip_at,
	skipped_by_twir_at,
	removed_or_reconciled_at,
	unknown_at,
	created_at,
	updated_at
) VALUES (
	@id,
	@channel_id,
	@track_id,
	@track_uri,
	@title,
	@artist,
	@album,
	@duration_ms,
	@requester_user_id,
	@requester_name,
	@requester_display_name,
	@source,
	@queue_position,
	@status,
	@queued_at,
	@playing_observed_at,
	@played_observed_at,
	@cancelled_pending_skip_at,
	@skipped_by_twir_at,
	@removed_or_reconciled_at,
	@unknown_at,
	@created_at,
	@updated_at
) RETURNING ` + selectFields

	rows, err := c.pool.Query(
		ctx,
		query,
		pgx.NamedArgs{
			"id":                        req.ID,
			"channel_id":                req.ChannelID,
			"track_id":                  req.TrackID,
			"track_uri":                 req.TrackURI,
			"title":                     req.Title,
			"artist":                    req.Artist,
			"album":                     req.Album,
			"duration_ms":               req.DurationMs,
			"requester_user_id":         req.RequesterUserID,
			"requester_name":            req.RequesterName,
			"requester_display_name":    req.RequesterDisplayName,
			"source":                    req.Source,
			"queue_position":            req.QueuePosition,
			"status":                    req.Status,
			"queued_at":                 req.QueuedAt,
			"playing_observed_at":       req.PlayingObservedAt,
			"played_observed_at":        req.PlayedObservedAt,
			"cancelled_pending_skip_at": req.CancelledPendingSkipAt,
			"skipped_by_twir_at":        req.SkippedByTwirAt,
			"removed_or_reconciled_at":  req.RemovedOrReconciledAt,
			"unknown_at":                req.UnknownAt,
			"created_at":                req.CreatedAt,
			"updated_at":                req.UpdatedAt,
		},
	)
	if err != nil {
		return spotify_song_request.Nil, fmt.Errorf("create spotify song request: %w", err)
	}

	created, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return spotify_song_request.Nil, fmt.Errorf("collect created spotify song request: %w", err)
	}

	return created.toEntity(), nil
}

func (c *Pgx) GetByID(
	ctx context.Context,
	id string,
) (spotify_song_request.SpotifySongRequest, error) {
	query := `SELECT ` + selectFields + `
FROM spotify_song_requests
WHERE id = CAST(@id AS uuid)
LIMIT 1`

	rows, err := c.pool.Query(ctx, query, pgx.NamedArgs{"id": id})
	if err != nil {
		return spotify_song_request.Nil, fmt.Errorf("query spotify song request by id: %w", err)
	}

	request, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return spotify_song_request.Nil, spotify_song_requests.ErrNotFound
		}

		return spotify_song_request.Nil, fmt.Errorf("collect spotify song request by id: %w", err)
	}

	return request.toEntity(), nil
}

func (c *Pgx) GetActiveByChannel(
	ctx context.Context,
	channelID string,
) ([]spotify_song_request.SpotifySongRequest, error) {
	query := `SELECT ` + selectFields + `
FROM spotify_song_requests
WHERE channel_id = @channel_id
	AND status IN ('queued', 'playing', 'cancelled_pending_skip')
ORDER BY queue_position ASC, created_at ASC`

	rows, err := c.pool.Query(ctx, query, pgx.NamedArgs{"channel_id": channelID})
	if err != nil {
		return nil, fmt.Errorf("query active spotify song requests by channel: %w", err)
	}

	requests, err := pgx.CollectRows(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return nil, fmt.Errorf("collect active spotify song requests by channel: %w", err)
	}

	result := make([]spotify_song_request.SpotifySongRequest, 0, len(requests))
	for _, request := range requests {
		result = append(result, request.toEntity())
	}

	return result, nil
}

func (c *Pgx) GetActiveByRequester(
	ctx context.Context,
	channelID, requesterName string,
) ([]spotify_song_request.SpotifySongRequest, error) {
	query := `SELECT ` + selectFields + `
FROM spotify_song_requests
WHERE channel_id = @channel_id
	AND requester_name = @requester_name
	AND status IN ('queued', 'playing', 'cancelled_pending_skip')
ORDER BY queue_position ASC, created_at ASC`

	rows, err := c.pool.Query(
		ctx,
		query,
		pgx.NamedArgs{"channel_id": channelID, "requester_name": requesterName},
	)
	if err != nil {
		return nil, fmt.Errorf("query active spotify song requests by requester: %w", err)
	}

	requests, err := pgx.CollectRows(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return nil, fmt.Errorf("collect active spotify song requests by requester: %w", err)
	}

	result := make([]spotify_song_request.SpotifySongRequest, 0, len(requests))
	for _, request := range requests {
		result = append(result, request.toEntity())
	}

	return result, nil
}

func (c *Pgx) GetActiveChannels(ctx context.Context) ([]string, error) {
	query := `SELECT DISTINCT channel_id
FROM spotify_song_requests
WHERE status IN ('queued', 'playing', 'cancelled_pending_skip')
ORDER BY channel_id ASC`

	rows, err := c.pool.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query active spotify song request channels: %w", err)
	}

	channelIDs, err := pgx.CollectRows(rows, pgx.RowTo[string])
	if err != nil {
		return nil, fmt.Errorf("collect active spotify song request channels: %w", err)
	}

	result := make([]string, 0, len(channelIDs))
	for _, channelID := range channelIDs {
		result = append(result, channelID)
	}

	return result, nil
}

func (c *Pgx) CountActiveByChannel(
	ctx context.Context,
	channelID string,
) (int64, error) {
	query := `SELECT count(*)
FROM spotify_song_requests
WHERE channel_id = @channel_id
	AND status IN ('queued', 'playing', 'cancelled_pending_skip')`

	var count int64
	if err := c.pool.QueryRow(ctx, query, pgx.NamedArgs{"channel_id": channelID}).Scan(&count); err != nil {
		return 0, fmt.Errorf("count active spotify song requests by channel: %w", err)
	}

	return count, nil
}

func (c *Pgx) CountActiveByRequester(
	ctx context.Context,
	channelID, requesterName string,
) (int64, error) {
	query := `SELECT count(*)
FROM spotify_song_requests
WHERE channel_id = @channel_id
	AND requester_name = @requester_name
	AND status IN ('queued', 'playing', 'cancelled_pending_skip')`

	var count int64
	if err := c.pool.QueryRow(
		ctx,
		query,
		pgx.NamedArgs{"channel_id": channelID, "requester_name": requesterName},
	).Scan(&count); err != nil {
		return 0, fmt.Errorf("count active spotify song requests by requester: %w", err)
	}

	return count, nil
}

func (c *Pgx) ListByChannel(
	ctx context.Context,
	channelID string,
	limit int,
) ([]spotify_song_request.SpotifySongRequest, error) {
	query := `SELECT ` + selectFields + `
FROM spotify_song_requests
WHERE channel_id = @channel_id
ORDER BY queue_position ASC, created_at ASC
LIMIT @limit`

	rows, err := c.pool.Query(ctx, query, pgx.NamedArgs{"channel_id": channelID, "limit": limit})
	if err != nil {
		return nil, fmt.Errorf("query spotify song requests by channel: %w", err)
	}

	requests, err := pgx.CollectRows(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return nil, fmt.Errorf("collect spotify song requests by channel: %w", err)
	}

	result := make([]spotify_song_request.SpotifySongRequest, 0, len(requests))
	for _, request := range requests {
		result = append(result, request.toEntity())
	}

	return result, nil
}

func (c *Pgx) UpdateStatus(
	ctx context.Context,
	id string,
	status spotify_song_request.Status,
) error {
	clause, err := statusUpdateClause(status)
	if err != nil {
		return err
	}

	query := `UPDATE spotify_song_requests
SET ` + clause + `
WHERE id = CAST(@id AS uuid)`

	cmd, err := c.pool.Exec(ctx, query, pgx.NamedArgs{"id": id, "status": status})
	if err != nil {
		return fmt.Errorf("update spotify song request status: %w", err)
	}

	if cmd.RowsAffected() == 0 {
		return spotify_song_requests.ErrNotFound
	}

	return nil
}

func (c *Pgx) CancelPendingSkip(ctx context.Context, id string) error {
	return c.UpdateStatus(ctx, id, spotify_song_request.StatusCancelledPendingSkip)
}

func statusUpdateClause(status spotify_song_request.Status) (string, error) {
	switch status {
	case spotify_song_request.StatusQueued:
		return "status = @status, updated_at = now()", nil
	case spotify_song_request.StatusPlaying:
		return "status = @status, playing_observed_at = now(), updated_at = now()", nil
	case spotify_song_request.StatusPlayed:
		return "status = @status, played_observed_at = now(), updated_at = now()", nil
	case spotify_song_request.StatusCancelledPendingSkip:
		return "status = @status, cancelled_pending_skip_at = now(), updated_at = now()", nil
	case spotify_song_request.StatusSkippedByTwir:
		return "status = @status, skipped_by_twir_at = now(), updated_at = now()", nil
	case spotify_song_request.StatusRemovedOrReconciled:
		return "status = @status, removed_or_reconciled_at = now(), updated_at = now()", nil
	case spotify_song_request.StatusUnknown:
		return "status = @status, unknown_at = now(), updated_at = now()", nil
	default:
		return "", fmt.Errorf("%w: %q", spotify_song_request.ErrInvalidStatus, status)
	}
}
