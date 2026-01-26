package pgx

import (
	"cmp"
	"context"
	"database/sql"
	"fmt"
	"slices"

	"github.com/Masterminds/squirrel"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	timersentity "github.com/twirapp/twir/libs/entities/timers"
	"github.com/twirapp/twir/libs/repositories/timers"
)

type Opts struct {
	Pgx *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool: opts.Pgx,
	}
}

func NewFx(pgxpool *pgxpool.Pool) *Pgx {
	return New(
		Opts{
			Pgx: pgxpool,
		},
	)
}

var (
	_  timers.Repository = (*Pgx)(nil)
	sq                   = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type Pgx struct {
	pool *pgxpool.Pool
}

type scanModel struct {
	ID                       uuid.UUID
	ChannelID                string
	Name                     string
	Enabled                  bool
	OfflineEnabled           bool
	OnlineEnabled            bool
	TimeInterval             int
	MessageInterval          int
	LastTriggerMessageNumber int
	Responses                []scanModelResponse `db:"responses"`
}

type scanModelResponse struct {
	ID            uuid.UUID
	Text          string
	IsAnnounce    bool
	TimerID       uuid.UUID
	Count         int
	AnnounceColor int
}

func (c scanModel) toEntity() timersentity.Timer {
	responses := make([]timersentity.Response, len(c.Responses))
	for i, r := range c.Responses {
		responses[i] = timersentity.Response{
			ID:            r.ID,
			Text:          r.Text,
			IsAnnounce:    r.IsAnnounce,
			TimerID:       r.TimerID,
			Count:         r.Count,
			AnnounceColor: timersentity.AnnounceColor(r.AnnounceColor),
		}
	}

	return timersentity.Timer{
		ID:                       c.ID,
		ChannelID:                c.ChannelID,
		Name:                     c.Name,
		Enabled:                  c.Enabled,
		OfflineEnabled:           c.OfflineEnabled,
		OnlineEnabled:            c.OnlineEnabled,
		TimeInterval:             c.TimeInterval,
		MessageInterval:          c.MessageInterval,
		LastTriggerMessageNumber: c.LastTriggerMessageNumber,
		Responses:                responses,
	}
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (timersentity.Timer, error) {
	query := `
SELECT t."id", t."channelId", t."name", t."enabled", t."offline_enabled", t."online_enabled", t."timeInterval", t."messageInterval", t."lastTriggerMessageNumber",
			 r."id" response_id, r."text" response_text, r."isAnnounce" response_is_announce, r."timerId" response_timer_id, r.count response_count, r."announce_color" response_announce_color
FROM "channels_timers" t
LEFT JOIN "channels_timers_responses" r ON t."id" = r."timerId"
WHERE
   t."id" = $1
ORDER BY t.id;
`
	rows, err := c.pool.Query(ctx, query, id)
	if err != nil {
		return timersentity.Nil, err
	}
	defer rows.Close()

	var timer scanModel
	for rows.Next() {
		var (
			responseID, responseTimerID sql.Null[uuid.UUID]
			responseText                sql.Null[string]
			responseIsAnnounce          sql.Null[bool]
			responseCount               sql.Null[int]
			responseAnnounceColor       sql.Null[int]
		)

		if err := rows.Scan(
			&timer.ID,
			&timer.ChannelID,
			&timer.Name,
			&timer.Enabled,
			&timer.OfflineEnabled,
			&timer.OnlineEnabled,
			&timer.TimeInterval,
			&timer.MessageInterval,
			&timer.LastTriggerMessageNumber,
			&responseID,
			&responseText,
			&responseIsAnnounce,
			&responseTimerID,
			&responseCount,
			&responseAnnounceColor,
		); err != nil {
			return timersentity.Nil, err
		}

		if responseID.Valid {
			timer.Responses = append(
				timer.Responses, scanModelResponse{
					ID:            responseID.V,
					Text:          responseText.V,
					IsAnnounce:    responseIsAnnounce.V,
					TimerID:       responseTimerID.V,
					Count:         responseCount.V,
					AnnounceColor: responseAnnounceColor.V,
				},
			)
		}
	}

	if rows.Err() != nil {
		return timersentity.Nil, rows.Err()
	}

	return timer.toEntity(), nil
}

func (c *Pgx) CountByChannelID(ctx context.Context, channelID string) (int, error) {
	query := `SELECT count(*) from "channels_timers" where "channelId" = $1`

	var count int
	err := c.pool.QueryRow(ctx, query, channelID).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Pgx) Create(ctx context.Context, data timers.CreateInput) (timersentity.Timer, error) {
	createQuery := `
INSERT INTO "channels_timers" ("channelId", "name", "enabled", "offline_enabled", "online_enabled", "timeInterval", "messageInterval")
VALUES ($1, $2, $3, $4, $5, $6, $7)
RETURNING "id", "channelId", "name", "enabled", "offline_enabled", "online_enabled", "timeInterval", "messageInterval"
`
	createResponseQuery := `
INSERT INTO "channels_timers_responses" ("id", "text", "isAnnounce", "timerId", count, "announce_color")
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING "id", "text", "isAnnounce", "timerId", count, "announce_color"
`
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return timersentity.Nil, err
	}

	defer tx.Rollback(ctx)

	var newTimer scanModel
	if err := tx.QueryRow(
		ctx,
		createQuery,
		data.ChannelID,
		data.Name,
		data.Enabled,
		data.OfflineEnabled,
		data.OnlineEnabled,
		data.TimeInterval,
		data.MessageInterval,
	).Scan(
		&newTimer.ID,
		&newTimer.ChannelID,
		&newTimer.Name,
		&newTimer.Enabled,
		&newTimer.OfflineEnabled,
		&newTimer.OnlineEnabled,
		&newTimer.TimeInterval,
		&newTimer.MessageInterval,
	); err != nil {
		return timersentity.Nil, fmt.Errorf("cannot create timer: %w", err)
	}

	for _, r := range data.Responses {
		var newResponse scanModelResponse

		if err := tx.QueryRow(
			ctx,
			createResponseQuery,
			uuid.New(),
			r.Text,
			r.IsAnnounce,
			newTimer.ID,
			r.Count,
			int(r.AnnounceColor),
		).Scan(
			&newResponse.ID,
			&newResponse.Text,
			&newResponse.IsAnnounce,
			&newResponse.TimerID,
			&newResponse.Count,
			&newResponse.AnnounceColor,
		); err != nil {
			return timersentity.Nil, fmt.Errorf("cannot create response for timer: %w", err)
		}
		newTimer.Responses = append(newTimer.Responses, newResponse)
	}

	if err := tx.Commit(ctx); err != nil {
		return timersentity.Nil, err
	}

	return newTimer.toEntity(), nil
}

func (c *Pgx) GetAllByChannelID(ctx context.Context, channelID string) ([]timersentity.Timer, error) {
	query := `
SELECT t."id", t."channelId", t."name", t."enabled", t."offline_enabled", t."online_enabled", t."timeInterval", t."messageInterval", t."lastTriggerMessageNumber",
			 r."id" response_id, r."text" response_text, r."isAnnounce" response_is_announce, r."timerId" response_timer_id, r.count response_count, r."announce_color" response_announce_color
FROM "channels_timers" t
LEFT JOIN "channels_timers_responses" r ON t."id" = r."timerId"
WHERE t."channelId" = $1
ORDER BY t."id";
`

	rows, err := c.pool.Query(ctx, query, channelID)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	timersMap := make(map[uuid.UUID]*timersentity.Timer)
	for rows.Next() {
		var timer scanModel
		var (
			responseID, responseTimerID sql.Null[uuid.UUID]
			responseText                sql.Null[string]
			responseIsAnnounce          sql.Null[bool]
			responseCount               sql.Null[int]
			responseAnnounceColor       sql.Null[int]
		)

		if err := rows.Scan(
			&timer.ID,
			&timer.ChannelID,
			&timer.Name,
			&timer.Enabled,
			&timer.OfflineEnabled,
			&timer.OnlineEnabled,
			&timer.TimeInterval,
			&timer.MessageInterval,
			&timer.LastTriggerMessageNumber,
			&responseID,
			&responseText,
			&responseIsAnnounce,
			&responseTimerID,
			&responseCount,
			&responseAnnounceColor,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %v, %w", timer.ID, err)
		}
		if _, ok := timersMap[timer.ID]; !ok {
			e := timer.toEntity()
			timersMap[timer.ID] = &e
		}

		if responseID.Valid {
			timersMap[timer.ID].Responses = append(
				timersMap[timer.ID].Responses, timersentity.Response{
					ID:            responseID.V,
					Text:          responseText.V,
					IsAnnounce:    responseIsAnnounce.V,
					TimerID:       responseTimerID.V,
					Count:         responseCount.V,
					AnnounceColor: timersentity.AnnounceColor(responseAnnounceColor.V),
				},
			)
		}
	}

	result := make([]timersentity.Timer, 0, len(timersMap))
	for _, timer := range timersMap {
		result = append(result, *timer)
	}

	slices.SortFunc(
		result, func(i, j timersentity.Timer) int {
			return cmp.Compare(i.ID.String(), j.ID.String())
		},
	)

	return result, nil
}

func (c *Pgx) UpdateByID(ctx context.Context, id uuid.UUID, data timers.UpdateInput) (
	timersentity.Timer,
	error,
) {
	updateBuilder := sq.Update("channels_timers")

	if data.Name != nil {
		updateBuilder = updateBuilder.Set("name", *data.Name)
	}

	if data.Enabled != nil {
		updateBuilder = updateBuilder.Set("enabled", *data.Enabled)
	}

	if data.OfflineEnabled != nil {
		updateBuilder = updateBuilder.Set(`"offline_enabled"`, *data.OfflineEnabled)
	}

	if data.OnlineEnabled != nil {
		updateBuilder = updateBuilder.Set(`"online_enabled"`, *data.OnlineEnabled)
	}

	if data.TimeInterval != nil {
		updateBuilder = updateBuilder.Set(`"timeInterval"`, *data.TimeInterval)
	}

	if data.MessageInterval != nil {
		updateBuilder = updateBuilder.Set(`"messageInterval"`, *data.MessageInterval)
	}

	updateBuilder = updateBuilder.Where(squirrel.Eq{"id": id})

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return timersentity.Nil, err
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return timersentity.Nil, err
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return timersentity.Nil, err
	}
	if result.RowsAffected() == 0 {
		return timersentity.Nil, timers.ErrTimerNotFound
	}

	if len(data.Responses) > 0 {
		_, err = tx.Exec(ctx, `DELETE FROM channels_timers_responses WHERE "timerId" = $1`, id)
		if err != nil {
			return timersentity.Nil, err
		}
		for _, r := range data.Responses {
			_, err := tx.Exec(
				ctx,
				`INSERT INTO "channels_timers_responses" ("id", "text", "isAnnounce", "timerId", count, "announce_color") VALUES ($1, $2, $3, $4, $5, $6)`,
				uuid.New(),
				r.Text,
				r.IsAnnounce,
				id,
				r.Count,
				int(r.AnnounceColor),
			)
			if err != nil {
				return timersentity.Nil, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return timersentity.Nil, err
	}

	return c.GetByID(ctx, id)
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM "channels_timers" WHERE "id" = $1`

	rows, err := c.pool.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if rows.RowsAffected() != 1 {
		return timers.ErrTimerNotFound
	}

	return nil
}

func (c *Pgx) Count(ctx context.Context, input timers.CountInput) (int64, error) {
	qb := sq.Select("COUNT(*)").From("channels_timers")

	if input.ChannelID != nil {
		qb = qb.Where(squirrel.Eq{`"channelId"`: *input.ChannelID})
	}

	if input.Enabled != nil {
		qb = qb.Where(squirrel.Eq{"enabled": *input.Enabled})
	}

	query, args, err := qb.ToSql()
	if err != nil {
		return 0, fmt.Errorf("failed to build count query: %w", err)
	}

	var count int64
	err = c.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count timers: %w", err)
	}

	return count, nil
}

func (c *Pgx) GetMany(ctx context.Context, input timers.GetManyInput) ([]timersentity.Timer, error) {
	qb := sq.Select(
		"t.id",
		`t."channelId"`,
		"t.name",
		"t.enabled",
		`t."offline_enabled"`,
		`t."online_enabled"`,
		`t."timeInterval"`,
		`t."messageInterval"`,
		`t."lastTriggerMessageNumber"`,
		`COALESCE(
		json_agg(
			json_build_object(
				'id',            tr.id,
				'text',          tr.text,
				'isAnnounce',    tr."isAnnounce",
				'timerId',       tr."timerId",
				'count',         tr.count,
				'announce_color',tr.announce_color
			)
		),
		'[]'::json
	) AS responses`,
	).
		From("channels_timers t").
		LeftJoin(`channels_timers_responses tr ON t.id = tr."timerId"`).
		GroupBy("t.id")

	if input.ChannelID != nil {
		qb = qb.Where(squirrel.Eq{"t.channelId": *input.ChannelID})
	}
	if input.Enabled != nil {
		qb = qb.Where(squirrel.Eq{"t.enabled": *input.Enabled})
	}

	if input.Limit > 0 {
		qb = qb.Limit(uint64(input.Limit))
	}
	if input.Offset > 0 {
		qb = qb.Offset(uint64(input.Offset))
	}

	qb = qb.OrderBy("t.id")

	query, args, err := qb.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build get many query: %w", err)
	}

	rows, err := c.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query timers: %w", err)
	}
	defer rows.Close()

	var result []timersentity.Timer

	for rows.Next() {
		var (
			timer        scanModel
			rawResponses []byte
		)

		if err := rows.Scan(
			&timer.ID,
			&timer.ChannelID,
			&timer.Name,
			&timer.Enabled,
			&timer.OfflineEnabled,
			&timer.OnlineEnabled,
			&timer.TimeInterval,
			&timer.MessageInterval,
			&timer.LastTriggerMessageNumber,
			&rawResponses, // responses JSON
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var responses []scanModelResponse
		if err := json.Unmarshal(rawResponses, &responses); err != nil {
			return nil, fmt.Errorf("failed to unmarshal responses JSON: %w", err)
		}
		timer.Responses = responses

		result = append(result, timer.toEntity())
	}

	return result, nil
}
