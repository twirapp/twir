package pgx

import (
	"context"
	"errors"
	"fmt"
	"net/netip"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	uploadedfile "github.com/twirapp/twir/libs/entities/uploaded_file"
	"github.com/twirapp/twir/libs/repositories/uploaded_files"
)

var (
	_             uploadedfiles.Repository = (*Pgx)(nil)
	sq                                     = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
	selectColumns                          = []string{
		"id",
		"public_id",
		"uploaded_by_user_id",
		"file_name",
		"mime_type",
		"extension",
		"size_bytes",
		"s3_key",
		"delete_key",
		"user_agent",
		"user_ip",
		"expires_at",
		"created_at",
	}
	selectColumnsString = strings.Join(selectColumns, ", ")
)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return &Pgx{
		pool:   pool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

type scanModel struct {
	ID               uuid.UUID
	PublicID         string
	UploadedByUserId *string
	FileName         *string
	MimeType         string
	Extension        string
	SizeBytes        int64
	S3Key            string
	DeleteKey        string
	UserAgent        *string
	UserIp           *netip.Addr
	ExpiresAt        time.Time
	CreatedAt        time.Time

	isNil bool
}

func (m scanModel) toEntity() uploadedfile.Entity {
	return uploadedfile.Entity{
		ID:               m.ID,
		PublicID:         m.PublicID,
		UploadedByUserID: m.UploadedByUserId,
		FileName:         m.FileName,
		MimeType:         m.MimeType,
		Extension:        m.Extension,
		SizeBytes:        m.SizeBytes,
		S3Key:            m.S3Key,
		DeleteKey:        m.DeleteKey,
		UserAgent:        m.UserAgent,
		UserIP:           m.UserIp,
		ExpiresAt:        m.ExpiresAt,
		CreatedAt:        m.CreatedAt,
	}
}

func (c *Pgx) Create(ctx context.Context, input uploadedfiles.CreateInput) (uploadedfile.Entity, error) {
	query := `
INSERT INTO uploaded_files (
		public_id,
		uploaded_by_user_id,
		file_name,
		mime_type,
		extension,
		size_bytes,
		s3_key,
		delete_key,
		user_agent,
		user_ip,
		expires_at
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
RETURNING ` + selectColumnsString

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.PublicID,
		input.UploadedByUserID,
		input.FileName,
		input.MimeType,
		input.Extension,
		input.SizeBytes,
		input.S3Key,
		input.DeleteKey,
		input.UserAgent,
		input.UserIP,
		input.ExpiresAt,
	)
	if err != nil {
		return uploadedfile.Nil, fmt.Errorf("create uploaded file query: %w", err)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return uploadedfile.Nil, fmt.Errorf("create uploaded file row: %w", err)
	}

	return result.toEntity(), nil
}

func (c *Pgx) GetByPublicID(ctx context.Context, publicID string) (uploadedfile.Entity, error) {
	query := `
SELECT ` + selectColumnsString + `
FROM uploaded_files
WHERE public_id = $1
LIMIT 1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, publicID)
	if err != nil {
		return uploadedfile.Nil, fmt.Errorf("get uploaded file query: %w", err)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return uploadedfile.Nil, uploadedfiles.ErrNotFound
		}
		return uploadedfile.Nil, fmt.Errorf("get uploaded file row: %w", err)
	}

	return result.toEntity(), nil
}

func (c *Pgx) GetManyByPublicIDs(ctx context.Context, publicIDs []string) ([]uploadedfile.Entity, error) {
	query := `
SELECT ` + selectColumnsString + `
FROM uploaded_files
WHERE public_id = ANY($1)
ORDER BY created_at DESC
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, publicIDs)
	if err != nil {
		return nil, fmt.Errorf("get uploaded files query: %w", err)
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return nil, fmt.Errorf("get uploaded files rows: %w", err)
	}

	items := make([]uploadedfile.Entity, 0, len(models))
	for _, model := range models {
		items = append(items, model.toEntity())
	}

	return items, nil
}

func (c *Pgx) GetList(ctx context.Context, input uploadedfiles.GetListInput) (uploadedfiles.GetListOutput, error) {
	queryBuilder := sq.Select(selectColumns...).
		From("uploaded_files").
		Where(squirrel.Eq{"uploaded_by_user_id": input.UserID}).
		OrderBy("created_at DESC")
	countQueryBuilder := sq.Select("COUNT(*)").
		From("uploaded_files").
		Where(squirrel.Eq{"uploaded_by_user_id": input.UserID})

	countQuery, countArgs, err := countQueryBuilder.ToSql()
	if err != nil {
		return uploadedfiles.GetListOutput{}, fmt.Errorf("build uploaded files count query: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	var count int64
	if err := conn.QueryRow(ctx, countQuery, countArgs...).Scan(&count); err != nil {
		return uploadedfiles.GetListOutput{}, fmt.Errorf("query uploaded files count: %w", err)
	}

	perPage := input.PerPage
	if perPage == 0 {
		perPage = 20
	}
	queryBuilder = queryBuilder.Limit(uint64(perPage)).Offset(uint64(input.Page * perPage))
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return uploadedfiles.GetListOutput{}, fmt.Errorf("build uploaded files query: %w", err)
	}

	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return uploadedfiles.GetListOutput{}, fmt.Errorf("query uploaded files: %w", err)
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return uploadedfiles.GetListOutput{}, fmt.Errorf("collect uploaded files: %w", err)
	}

	items := make([]uploadedfile.Entity, 0, len(models))
	for _, model := range models {
		items = append(items, model.toEntity())
	}

	return uploadedfiles.GetListOutput{Items: items, Total: int(count)}, nil
}

func (c *Pgx) DeleteByID(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM uploaded_files
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	if _, err := conn.Exec(ctx, query, id); err != nil {
		return fmt.Errorf("delete uploaded file: %w", err)
	}

	return nil
}

func (c *Pgx) GetExpired(ctx context.Context, limit int) ([]uploadedfile.Entity, error) {
	query := `
SELECT ` + selectColumnsString + `
FROM uploaded_files
WHERE expires_at < now()
ORDER BY expires_at ASC
LIMIT $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, limit)
	if err != nil {
		return nil, fmt.Errorf("get expired uploaded files query: %w", err)
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[scanModel])
	if err != nil {
		return nil, fmt.Errorf("get expired uploaded files rows: %w", err)
	}

	items := make([]uploadedfile.Entity, 0, len(models))
	for _, model := range models {
		items = append(items, model.toEntity())
	}

	return items, nil
}
