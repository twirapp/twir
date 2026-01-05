package postgres

import (
	"context"
	"errors"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/chat_translation"
	"github.com/twirapp/twir/libs/repositories/chat_translation/model"
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
	_  chat_translation.Repository = (*Pgx)(nil)
	sq                             = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) GetByChannelID(ctx context.Context, channelID string) (model.ChatTranslation, error) {
	query := `
SELECT id, channel_id, created_at, updated_at, enabled, target_language, excluded_languages, use_italic, excluded_users_ids
FROM channels_chat_translation_settings
WHERE channel_id = $1
LIMIT 1;
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return model.ChatTranslationNil, fmt.Errorf("query err: %w", err)
	}

	translation, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ChatTranslation])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ChatTranslationNil, chat_translation.ErrSettingsNotFound
		}
		return model.ChatTranslationNil, fmt.Errorf("collect err: %w", err)
	}

	return translation, nil
}

func (c *Pgx) Create(
	ctx context.Context,
	input chat_translation.CreateInput,
) (model.ChatTranslation, error) {
	query := `
INSERT INTO channels_chat_translation_settings (channel_id, enabled, target_language, excluded_languages, use_italic, excluded_users_ids)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING id, channel_id, created_at, updated_at, enabled, target_language, excluded_languages, use_italic, excluded_users_ids
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(
		ctx,
		query,
		input.ChannelID,
		input.Enabled,
		input.TargetLanguage,
		input.ExcludedLanguages,
		input.UseItalic,
		input.ExcludedUsersIDs,
	)
	if err != nil {
		return model.ChatTranslationNil, fmt.Errorf("query err: %w", err)
	}

	translation, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ChatTranslation])
	if err != nil {
		return model.ChatTranslationNil, fmt.Errorf("collect err: %w", err)
	}

	return translation, nil
}

func (c *Pgx) Update(
	ctx context.Context,
	id uuid.UUID,
	input chat_translation.UpdateInput,
) (model.ChatTranslation, error) {
	builder := sq.Update("channels_chat_translation_settings").
		Where(squirrel.Eq{"id": id.String()}).
		Suffix("RETURNING id, channel_id, created_at, updated_at, enabled, target_language, excluded_languages, use_italic, excluded_users_ids")

	if input.Enabled != nil {
		builder = builder.Set("enabled", *input.Enabled)
	}

	if input.TargetLanguage != nil {
		builder = builder.Set("target_language", *input.TargetLanguage)
	}

	if input.ExcludedLanguages != nil {
		builder = builder.Set("excluded_languages", *input.ExcludedLanguages)
	}

	if input.UseItalic != nil {
		builder = builder.Set("use_italic", *input.UseItalic)
	}

	if input.ExcludedUsersIDs != nil {
		builder = builder.Set("excluded_users_ids", *input.ExcludedUsersIDs)
	}

	query, args, err := builder.ToSql()
	if err != nil {
		return model.ChatTranslationNil, fmt.Errorf("build err: %w", err)
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.ChatTranslationNil, chat_translation.ErrSettingsNotFound
		}
		return model.ChatTranslationNil, fmt.Errorf("query err: %w", err)
	}

	translation, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[model.ChatTranslation])
	if err != nil {
		return model.ChatTranslationNil, fmt.Errorf("collect err: %w", err)
	}

	return translation, nil
}

func (c *Pgx) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
DELETE FROM channels_chat_translation_settings
WHERE id = $1
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err := conn.Exec(ctx, query, id.String())
	if errors.Is(err, pgx.ErrNoRows) {
		return chat_translation.ErrSettingsNotFound
	}
	return err
}
