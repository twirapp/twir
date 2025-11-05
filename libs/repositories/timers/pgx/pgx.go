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
	"github.com/twirapp/twir/libs/repositories/timers"
	"github.com/twirapp/twir/libs/repositories/timers/model"
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

var _ timers.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetByID(ctx context.Context, id uuid.UUID) (model.Timer, error) {
	query := `
SELECT t."id", t."channelId", t."name", t."enabled", t."timeInterval", t."messageInterval", t."lastTriggerMessageNumber",
			 r."id" response_id, r."text" response_text, r."isAnnounce" response_is_announce, r."timerId" response_timer_id, r.count response_count, r."announce_color" response_announce_color
FROM "channels_timers" t
LEFT JOIN "channels_timers_responses" r ON t."id" = r."timerId"
WHERE
   t."id" = $1
ORDER BY t.id;
`
	rows, err := c.pool.Query(ctx, query, id)
	if err != nil {
		return model.Nil, err
	}

	defer rows.Close()

	var timer model.Timer
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
			return model.Nil, err
		}

		if responseID.Valid {
			timer.Responses = append(
				timer.Responses, model.Response{
					ID:            responseID.V,
					Text:          responseText.V,
					IsAnnounce:    responseIsAnnounce.V,
					TimerID:       responseTimerID.V,
					Count:         responseCount.V,
					AnnounceColor: model.AnnounceColor(responseAnnounceColor.V),
				},
			)
		}
	}

	return timer, nil
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

func (c *Pgx) Create(ctx context.Context, data timers.CreateInput) (model.Timer, error) {
	createQuery := `
INSERT INTO "channels_timers" ("channelId", "name", "enabled", "timeInterval", "messageInterval")
VALUES ($1, $2, $3, $4, $5)
RETURNING "id", "channelId", "name", "enabled", "timeInterval", "messageInterval"
`
	createResponseQuery := `
INSERT INTO "channels_timers_responses" ("id", "text", "isAnnounce", "timerId", count, "announce_color")
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING "id", "text", "isAnnounce", "timerId", count, "announce_color"
`
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return model.Nil, err
	}

	defer tx.Rollback(ctx)

	var newTimer model.Timer

	if err := tx.QueryRow(
		ctx,
		createQuery,
		data.ChannelID,
		data.Name,
		data.Enabled,
		data.TimeInterval,
		data.MessageInterval,
	).Scan(
		&newTimer.ID,
		&newTimer.ChannelID,
		&newTimer.Name,
		&newTimer.Enabled,
		&newTimer.TimeInterval,
		&newTimer.MessageInterval,
	); err != nil {
		return model.Nil, fmt.Errorf("cannot create timer: %w", err)
	}

	for _, r := range data.Responses {
		var newResponse model.Response

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
			return model.Nil, fmt.Errorf("cannot create response for timer: %w", err)
		}
		newTimer.Responses = append(newTimer.Responses, newResponse)
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Nil, err
	}

	return newTimer, nil
}

func (c *Pgx) GetAllByChannelID(ctx context.Context, channelID string) ([]model.Timer, error) {
	query := `
SELECT t."id", t."channelId", t."name", t."enabled", t."timeInterval", t."messageInterval", t."lastTriggerMessageNumber",
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

	var timersMap = make(map[uuid.UUID]*model.Timer)
	for rows.Next() {
		var timer model.Timer
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
			timersMap[timer.ID] = &timer
		}

		if responseID.Valid {
			timersMap[timer.ID].Responses = append(
				timersMap[timer.ID].Responses, model.Response{
					ID:            responseID.V,
					Text:          responseText.V,
					IsAnnounce:    responseIsAnnounce.V,
					TimerID:       responseTimerID.V,
					Count:         responseCount.V,
					AnnounceColor: model.AnnounceColor(responseAnnounceColor.V),
				},
			)
		}
	}

	result := make([]model.Timer, 0, len(timersMap))
	for _, timer := range timersMap {
		result = append(result, *timer)
	}

	slices.SortFunc(
		result, func(i, j model.Timer) int {
			return cmp.Compare(i.ID.String(), j.ID.String())
		},
	)

	return result, nil
}

func (c *Pgx) UpdateByID(ctx context.Context, id uuid.UUID, data timers.UpdateInput) (
	model.Timer,
	error,
) {
	updateBuilder := sq.Update("channels_timers")

	if data.Name != nil {
		updateBuilder = updateBuilder.Set("name", *data.Name)
	}

	if data.Enabled != nil {
		updateBuilder = updateBuilder.Set("enabled", *data.Enabled)
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
		return model.Nil, err
	}

	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return model.Nil, err
	}
	defer tx.Rollback(ctx)

	result, err := tx.Exec(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}
	if result.RowsAffected() == 0 {
		return model.Nil, timers.ErrTimerNotFound
	}

	if len(data.Responses) > 0 {
		_, err = tx.Exec(ctx, `DELETE FROM channels_timers_responses WHERE "timerId" = $1`, id)
		if err != nil {
			return model.Nil, err
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
				return model.Nil, err
			}
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return model.Nil, err
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

func (c *Pgx) GetMany(ctx context.Context, input timers.GetManyInput) ([]model.Timer, error) {
	qb := sq.Select(
		"t.id",
		`t."channelId"`,
		"t.name",
		"t.enabled",
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

	var result []model.Timer

	for rows.Next() {
		var (
			timer        model.Timer
			rawResponses []byte
		)

		if err := rows.Scan(
			&timer.ID,
			&timer.ChannelID,
			&timer.Name,
			&timer.Enabled,
			&timer.TimeInterval,
			&timer.MessageInterval,
			&timer.LastTriggerMessageNumber,
			&rawResponses, // responses JSON
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		var responses []model.Response
		if err := json.Unmarshal(rawResponses, &responses); err != nil {
			return nil, fmt.Errorf("failed to unmarshal responses JSON: %w", err)
		}
		timer.Responses = responses

		result = append(result, timer)
	}

	return result, nil
}
