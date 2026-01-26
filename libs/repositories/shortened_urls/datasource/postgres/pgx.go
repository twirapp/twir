package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

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

var (
	_             shortened_urls.Repository = (*Pgx)(nil)
	sq                                      = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	selectColumns                           = []string{
		"shortened_urls.short_id",
		"shortened_urls.created_at",
		"shortened_urls.updated_at",
		"shortened_urls.url",
		"shortened_urls.created_by_user_id",
		"shortened_urls.views",
		"shortened_urls.user_ip",
		"shortened_urls.user_agent",
		"shortened_urls.domain_id",
		"domains.domain",
	}
	selectColumnsStr = strings.Join(selectColumns, ", ")
	selectColumnsCte = []string{
		"updated.short_id",
		"updated.created_at",
		"updated.updated_at",
		"updated.url",
		"updated.created_by_user_id",
		"updated.views",
		"updated.user_ip",
		"updated.user_agent",
		"updated.domain_id",
		"domains.domain",
	}
	selectColumnsCteStr  = strings.Join(selectColumnsCte, ", ")
	selectColumnsCreated = []string{
		"created.short_id",
		"created.created_at",
		"created.updated_at",
		"created.url",
		"created.created_by_user_id",
		"created.views",
		"created.user_ip",
		"created.user_agent",
		"created.domain_id",
		"domains.domain",
	}
	selectColumnsCreatedStr = strings.Join(selectColumnsCreated, ", ")
	baseSelectQuery         = `
SELECT ` + selectColumnsStr + `
FROM shortened_urls
LEFT JOIN short_links_custom_domains AS domains
	ON domains.id = shortened_urls.domain_id
`
)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetManyByShortIDs(
	ctx context.Context,
	domainID *string,
	ids []string,
) ([]model.ShortenedUrl, error) {
	query := baseSelectQuery + `
WHERE shortened_urls.short_id = ANY($1)
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	args := []any{ids}
	if domainID == nil {
		query += "\nAND shortened_urls.domain_id IS NULL"
	} else {
		query += "\nAND shortened_urls.domain_id = $2"
		args = append(args, *domainID)
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

func (c *Pgx) Delete(ctx context.Context, domainID *string, id string) error {
	query := `
DELETE FROM shortened_urls
WHERE short_id = $1
`

	args := []any{id}
	if domainID == nil {
		query += "\nAND domain_id IS NULL"
	} else {
		query += "\nAND domain_id = $2"
		args = append(args, *domainID)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, args...)
	return err
}

func (c *Pgx) ClearDomainForUser(ctx context.Context, domainID string, userID string) error {
	query := `
UPDATE shortened_urls
SET domain_id = NULL, updated_at = NOW()
WHERE domain_id = $1 AND created_by_user_id = $2
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, domainID, userID)
	return err
}

func (c *Pgx) CountDomainShortIDConflicts(ctx context.Context, domainID string, userID string) (int64, error) {
	query := `
SELECT COUNT(*)
FROM shortened_urls AS custom
JOIN shortened_urls AS defaults
	ON defaults.domain_id IS NULL
	AND defaults.short_id = custom.short_id
WHERE custom.domain_id = $1 AND custom.created_by_user_id = $2
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	var count int64
	if err := conn.QueryRow(ctx, query, domainID, userID).Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Pgx) Update(
	ctx context.Context,
	domainID *string,
	id string,
	input shortened_urls.UpdateInput,
) (model.ShortenedUrl, error) {
	updateBuilder := sq.Update("shortened_urls").
		Where(squirrel.Eq{"short_id": id}).
		Set("updated_at", squirrel.Expr("NOW()")).
		Suffix("RETURNING short_id, created_at, updated_at, url, created_by_user_id, views, user_ip, user_agent, domain_id")

	if domainID == nil {
		updateBuilder = updateBuilder.Where("domain_id IS NULL")
	} else {
		updateBuilder = updateBuilder.Where(squirrel.Eq{"domain_id": *domainID})
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

	if input.DomainID != nil {
		updateBuilder = updateBuilder.Set("domain_id", *input.DomainID)
	} else if input.ClearDomain {
		updateBuilder = updateBuilder.Set("domain_id", nil)
	}

	query, args, err := updateBuilder.ToSql()
	if err != nil {
		return model.Nil, err
	}

	query = `
WITH updated AS (
` + query + `
)
SELECT ` + selectColumnsCteStr + `
FROM updated
LEFT JOIN short_links_custom_domains AS domains
	ON domains.id = updated.domain_id
`

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

func (c *Pgx) GetByUrl(ctx context.Context, domainID *string, url string) (model.ShortenedUrl, error) {
	query := baseSelectQuery + `
WHERE shortened_urls.url = $1
`

	args := []any{url}
	if domainID == nil {
		query += "\nAND shortened_urls.domain_id IS NULL"
	} else {
		query += "\nAND shortened_urls.domain_id = $2"
		args = append(args, *domainID)
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

func (c *Pgx) GetByShortID(ctx context.Context, domainID *string, id string) (model.ShortenedUrl, error) {
	query := baseSelectQuery + `
WHERE shortened_urls.short_id = $1
`

	args := []any{id}
	if domainID == nil {
		query += "\nAND shortened_urls.domain_id IS NULL"
	} else {
		query += "\nAND shortened_urls.domain_id = $2"
		args = append(args, *domainID)
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
WITH created AS (
	INSERT INTO shortened_urls (short_id, url, created_by_user_id, user_ip, user_agent, domain_id)
	VALUES (@short_id, @url, @created_by_user_id, @user_ip, @user_agent, @domain_id)
	RETURNING short_id, created_at, updated_at, url, created_by_user_id, views, user_ip, user_agent, domain_id
)
SELECT ` + selectColumnsCreatedStr + `
FROM created
LEFT JOIN short_links_custom_domains AS domains
	ON domains.id = created.domain_id
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
			"domain_id":          input.DomainID,
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
		selectColumns...,
	).
		From("shortened_urls").
		LeftJoin("short_links_custom_domains AS domains ON domains.id = shortened_urls.domain_id")

	// Set sorting order based on SortBy parameter
	sortBy := input.SortBy
	if sortBy == "" {
		sortBy = "created_at" // Default sorting
	}

	if sortBy == "views" {
		queryBuilder = queryBuilder.OrderBy("shortened_urls.views DESC")
	} else {
		queryBuilder = queryBuilder.OrderBy("shortened_urls.created_at DESC")
	}

	countQueryBuilder := sq.Select("COUNT(*)").From("shortened_urls")

	if input.UserID != nil {
		queryBuilder = queryBuilder.Where("shortened_urls.created_by_user_id = ?", *input.UserID)
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
	case "shortened_urls_domain_id_short_id_unique_idx", "shortened_urls_domain_short_id_unique_idx", "shortened_urls_pkey", "shortened_urls_short_id_key":
		return shortened_urls.ErrShortIDAlreadyExists
	default:
		return fmt.Errorf("unique constraint error: %w", err)
	}
}
