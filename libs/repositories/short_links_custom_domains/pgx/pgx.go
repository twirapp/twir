package pgx

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
	shortlinkscustomdomains "github.com/twirapp/twir/libs/repositories/short_links_custom_domains"
	"github.com/twirapp/twir/libs/repositories/short_links_custom_domains/model"
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
	_  shortlinkscustomdomains.Repository = (*Pgx)(nil)
	sq                                    = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

var selectColumns = []string{
	"id",
	"user_id",
	"domain",
	"verified",
	"verification_token",
	"created_at",
	"updated_at",
}

var selectColumnsStr = strings.Join(selectColumns, ", ")

func (c *Pgx) GetByUserID(ctx context.Context, userID string) (model.CustomDomain, error) {
	query := `
SELECT ` + selectColumnsStr + `
FROM short_links_custom_domains
WHERE user_id = $1
LIMIT 1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, userID)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.CustomDomain])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, shortlinkscustomdomains.ErrNotFound
		}

		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) GetByDomain(ctx context.Context, domain string) (model.CustomDomain, error) {
	query := `
SELECT ` + selectColumnsStr + `
FROM short_links_custom_domains
WHERE domain = $1
LIMIT 1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, domain)
	if err != nil {
		return model.Nil, err
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.CustomDomain])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, shortlinkscustomdomains.ErrNotFound
		}

		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) CountByUserID(ctx context.Context, userID string) (int, error) {
	query := `
SELECT COUNT(*)
FROM short_links_custom_domains
WHERE user_id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	row := conn.QueryRow(ctx, query, userID)

	var count int
	if err := row.Scan(&count); err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Pgx) Create(ctx context.Context, input shortlinkscustomdomains.CreateInput) (model.CustomDomain, error) {
	query := `
INSERT INTO short_links_custom_domains (user_id, domain, verification_token)
VALUES (@user_id, @domain, @verification_token)
RETURNING ` + selectColumnsStr + `
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		pgx.NamedArgs{
			"user_id":            input.UserID,
			"domain":             input.Domain,
			"verification_token": input.VerificationToken,
		},
	)
	if err != nil {
		return model.Nil, mapUniqueError(err)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.CustomDomain])
	if err != nil {
		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) Update(ctx context.Context, id string, input shortlinkscustomdomains.UpdateInput) (model.CustomDomain, error) {
	updateBuilder := sq.Update("short_links_custom_domains").
		Set("updated_at", squirrel.Expr("NOW()")).
		Where(squirrel.Eq{"id": id}).
		Suffix("RETURNING " + selectColumnsStr)

	if input.Domain != nil {
		updateBuilder = updateBuilder.Set("domain", *input.Domain)
	}

	if input.Verified != nil {
		updateBuilder = updateBuilder.Set("verified", *input.Verified)
	}

	if input.VerificationToken != nil {
		updateBuilder = updateBuilder.Set("verification_token", *input.VerificationToken)
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

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.CustomDomain])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, shortlinkscustomdomains.ErrNotFound
		}

		return model.Nil, err
	}

	return result, nil
}

func (c *Pgx) Delete(ctx context.Context, id string) error {
	query := `
DELETE FROM short_links_custom_domains
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	result, err := conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return shortlinkscustomdomains.ErrNotFound
	}

	return nil
}

func (c *Pgx) VerifyDomain(ctx context.Context, id string) error {
	query := `
UPDATE short_links_custom_domains
SET verified = true,
	updated_at = NOW()
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	result, err := conn.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return shortlinkscustomdomains.ErrNotFound
	}

	return nil
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
	case "short_links_custom_domains_domain_key":
		return shortlinkscustomdomains.ErrDomainAlreadyExists
	case "short_links_custom_domains_user_id_key":
		return shortlinkscustomdomains.ErrUserAlreadyHasDomain
	default:
		return fmt.Errorf("unique constraint error: %w", err)
	}
}
