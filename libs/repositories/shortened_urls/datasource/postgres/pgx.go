package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/shortened_urls"
	"github.com/twirapp/twir/libs/repositories/shortened_urls/model"
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

var _ shortened_urls.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) Count(ctx context.Context, input shortened_urls.CountInput) (int64, error) {
	selectBuilder := sq.Select("COUNT(*)").From("shortened_urls")
	if input.UserID != "" {
		selectBuilder = selectBuilder.Where("created_by_user_id", input.UserID)
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var count int64
	err = c.pool.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Pgx) Delete(ctx context.Context, id string) error {
	query := `
DELETE FROM shortened_urls
WHERE short_id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id)
	return err
}

func (c *Pgx) Update(
	ctx context.Context,
	id string,
	input shortened_urls.UpdateInput,
) (model.ShortenedUrl, error) {
	updateBuilder := sq.Update("shortened_urls").
		Where(squirrel.Eq{"short_id": id}).
		Suffix("RETURNING short_id, created_at, updated_at, url, created_by_user_id, views")

	if input.Views != nil {
		updateBuilder = updateBuilder.Set("views", *input.Views)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ShortenedUrl])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) GetByUrl(ctx context.Context, url string) (model.ShortenedUrl, error) {
	query := `
SELECT short_id, created_at, updated_at, url, created_by_user_id, views
FROM shortened_urls
WHERE url = $1
LIMIT 1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, url)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ShortenedUrl])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, shortened_urls.ErrNotFound
		}
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) GetByShortID(ctx context.Context, id string) (model.ShortenedUrl, error) {
	query := `
SELECT short_id, created_at, updated_at, url, created_by_user_id, views
FROM shortened_urls
WHERE short_id = $1
LIMIT 1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, id)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ShortenedUrl])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, shortened_urls.ErrNotFound
		}
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) Create(ctx context.Context, input shortened_urls.CreateInput) (
	model.ShortenedUrl,
	error,
) {
	query := `
INSERT INTO shortened_urls (short_id, url, created_by_user_id)
VALUES ($1, $2, $3)
RETURNING short_id, created_at, updated_at, url, created_by_user_id, views
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, input.ShortID, input.URL, input.CreatedByUserID)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ShortenedUrl])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) GetList(ctx context.Context, input shortened_urls.GetListInput) (
	shortened_urls.GetListOutput,
	error,
) {
	queryBuilder := sq.Select("short_id, created_at, updated_at, url, created_by_user_id, views").
		From("shortened_urls").
		OrderBy("created_at DESC")

	countQueryBUilder := sq.Select("COUNT(*)").From("shortened_urls")

	if input.UserID != nil {
		queryBuilder = queryBuilder.Where("created_by_user_id", *input.UserID)
		countQueryBUilder = countQueryBUilder.Where("created_by_user_id", *input.UserID)
	}

	countQuery, countArgs, err := countQueryBUilder.ToSql()
	if err != nil {
		return shortened_urls.GetListOutput{}, fmt.Errorf("count query error: %w", err)
	}

	var count int64
	err = c.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&count)
	if err != nil {
		return shortened_urls.GetListOutput{}, fmt.Errorf("count query error: %w", err)
	}

	perPage := input.PerPage
	if perPage == 0 {
		perPage = 20
	}

	if input.Page > 0 && perPage > 0 {
		offset := input.Page * perPage
		queryBuilder = queryBuilder.Limit(uint64(perPage)).Offset(uint64(offset))
	}

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return shortened_urls.GetListOutput{}, fmt.Errorf("query error: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return shortened_urls.GetListOutput{}, fmt.Errorf("query error: %w", err)
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ShortenedUrl])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return shortened_urls.GetListOutput{}, shortened_urls.ErrNotFound
		}

		return shortened_urls.GetListOutput{}, fmt.Errorf("query error: %w", err)
	}

	return shortened_urls.GetListOutput{
		Items: models,
		Total: int(count),
	}, nil
}
