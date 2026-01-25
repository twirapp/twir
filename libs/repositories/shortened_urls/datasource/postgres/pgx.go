package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgconn"
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

func (c *Pgx) GetManyByShortIDs(
	ctx context.Context,
	domain *string,
	ids []string,
) ([]model.ShortenedUrl, error) {
	query := `
SELECT short_id, created_at, updated_at, url, created_by_user_id, views, user_ip, user_agent, domain
FROM shortened_urls
WHERE short_id = ANY($1)
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	args := []any{ids}
	if domain == nil {
		query += "\nAND domain IS NULL"
	} else {
		query += "\nAND domain = $2"
		args = append(args, *domain)
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[model.ShortenedUrl])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, shortened_urls.ErrNotFound
		}

		return nil, err
	}

	return models, nil
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

func (c *Pgx) Delete(ctx context.Context, domain *string, id string) error {
	query := `
DELETE FROM shortened_urls
WHERE short_id = $1
`

	args := []any{id}
	if domain == nil {
		query += "\nAND domain IS NULL"
	} else {
		query += "\nAND domain = $2"
		args = append(args, *domain)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, args...)
	return err
}

func (c *Pgx) Update(
	ctx context.Context,
	domain *string,
	id string,
	input shortened_urls.UpdateInput,
) (model.ShortenedUrl, error) {
	updateBuilder := sq.Update("shortened_urls").
		Where(squirrel.Eq{"short_id": id}).
		Set("updated_at", squirrel.Expr("NOW()")).
		Suffix("RETURNING short_id, created_at, updated_at, url, created_by_user_id, views, user_ip, user_agent, domain")

	if domain == nil {
		updateBuilder = updateBuilder.Where("domain IS NULL")
	} else {
		updateBuilder = updateBuilder.Where(squirrel.Eq{"domain": *domain})
	}

	if input.Views != nil {
		updateBuilder = updateBuilder.Set("views", *input.Views)
	}

	if input.ShortID != nil {
		updateBuilder = updateBuilder.Set("short_id", *input.ShortID)
	}

	if input.URL != nil {
		updateBuilder = updateBuilder.Set("url", *input.URL)
	}

	if input.Domain != nil {
		updateBuilder = updateBuilder.Set("domain", *input.Domain)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, mapUniqueError(err)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ShortenedUrl])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) GetByUrl(ctx context.Context, domain *string, url string) (model.ShortenedUrl, error) {
	query := `
SELECT short_id, created_at, updated_at, url, created_by_user_id, views, user_ip, user_agent, domain
FROM shortened_urls
WHERE url = $1
`

	args := []any{url}
	if domain == nil {
		query += "\nAND domain IS NULL"
	} else {
		query += "\nAND domain = $2"
		args = append(args, *domain)
	}
	query += "\nLIMIT 1"

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
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

func (c *Pgx) GetByShortID(ctx context.Context, domain *string, id string) (model.ShortenedUrl, error) {
	query := `
SELECT short_id, created_at, updated_at, url, created_by_user_id, views, user_ip, user_agent, domain
FROM shortened_urls
WHERE short_id = $1
`

	args := []any{id}
	if domain == nil {
		query += "\nAND domain IS NULL"
	} else {
		query += "\nAND domain = $2"
		args = append(args, *domain)
	}
	query += "\nLIMIT 1"

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
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
INSERT INTO shortened_urls (short_id, url, created_by_user_id, user_ip, user_agent, domain)
VALUES (@short_id, @url, @created_by_user_id, @user_ip, @user_agent, @domain)
RETURNING short_id, created_at, updated_at, url, created_by_user_id, views, user_ip, user_agent, domain
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		pgx.NamedArgs{
			"short_id":           input.ShortID,
			"url":                input.URL,
			"created_by_user_id": input.CreatedByUserID,
			"user_ip":            input.UserIp,
			"user_agent":         input.UserAgent,
			"domain":             input.Domain,
		},
	)
	if err != nil {
		return model.Nil, mapUniqueError(err)
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
	queryBuilder := sq.Select(
		"short_id",
		"created_at",
		"updated_at",
		"url",
		"created_by_user_id",
		"views",
		"user_ip",
		"user_agent",
		"domain",
	).
		From("shortened_urls")

	// Set sorting order based on SortBy parameter
	sortBy := input.SortBy
	if sortBy == "" {
		sortBy = "created_at" // Default sorting
	}

	if sortBy == "views" {
		queryBuilder = queryBuilder.OrderBy("views DESC")
	} else {
		queryBuilder = queryBuilder.OrderBy("created_at DESC")
	}

	countQueryBuilder := sq.Select("COUNT(*)").From("shortened_urls")

	if input.UserID != nil {
		queryBuilder = queryBuilder.Where("created_by_user_id = ?", *input.UserID)
		countQueryBuilder = countQueryBuilder.Where("created_by_user_id = ?", *input.UserID)
	}

	countQuery, countArgs, err := countQueryBuilder.ToSql()
	if err != nil {
		return shortened_urls.GetListOutput{}, fmt.Errorf("count query builder error: %w", err)
	}

	var count int64
	err = c.pool.QueryRow(ctx, countQuery, countArgs...).Scan(&count)
	if err != nil {
		return shortened_urls.GetListOutput{}, fmt.Errorf("count queryrow error: %w", err)
	}

	perPage := input.PerPage
	if perPage == 0 {
		perPage = 20
	}

	if input.Page > 0 || perPage > 0 {
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

func mapUniqueError(err error) error {
	var pgErr *pgconn.PgError
	if !errors.As(err, &pgErr) {
		return err
	}

	if pgErr.Code != "23505" {
		return err
	}

	switch pgErr.ConstraintName {
	case "shortened_urls_domain_short_id_unique_idx", "shortened_urls_pkey", "shortened_urls_short_id_key":
		return shortened_urls.ErrShortIDAlreadyExists
	default:
		return fmt.Errorf("unique constraint error: %w", err)
	}
}
