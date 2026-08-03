package postgres

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/entities/notification"
	"github.com/twirapp/twir/libs/repositories/notifications"
)

type Pgx struct {
	pool *pgxpool.Pool
}

func New(pool *pgxpool.Pool) *Pgx {
	return &Pgx{pool: pool}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(pool)
}

var _ notifications.Repository = (*Pgx)(nil)

func (p *Pgx) UpsertDiscord(
	ctx context.Context,
	input notifications.UpsertDiscordInput,
) (notifications.UpsertDiscordResult, error) {
	const query = `
		WITH previous AS (
			SELECT discord_attachment_keys
			FROM notifications
			WHERE discord_message_id = $1
		), upserted AS (
			INSERT INTO notifications (
				id,
				"createdAt",
				message,
				editor_js_json,
				discord_message_id,
				discord_channel_id,
				discord_attachment_keys,
				updated_at
			) VALUES (
				uuidv7()::text, $2, $3, $4, $1, $5, $6, $7
			)
			ON CONFLICT (discord_message_id) WHERE discord_message_id IS NOT NULL
			DO UPDATE SET
				message = EXCLUDED.message,
				editor_js_json = EXCLUDED.editor_js_json,
				discord_channel_id = EXCLUDED.discord_channel_id,
				discord_attachment_keys = EXCLUDED.discord_attachment_keys,
				updated_at = EXCLUDED.updated_at
			RETURNING
				id,
				"userId",
				message,
				editor_js_json,
				"createdAt",
				updated_at,
				discord_message_id,
				discord_channel_id,
				discord_attachment_keys,
				(xmax = 0) AS created
		)
		SELECT
			u.id,
			u."userId",
			u.message,
			u.editor_js_json,
			u."createdAt",
			u.updated_at,
			u.discord_message_id,
			u.discord_channel_id,
			u.discord_attachment_keys,
			COALESCE(p.discord_attachment_keys, '{}'),
			u.created
		FROM upserted u
		LEFT JOIN previous p ON TRUE
	`

	var entity notification.Notification
	var previousAttachmentKeys []string
	var created bool
	err := p.pool.QueryRow(
		ctx,
		query,
		input.DiscordMessageID,
		input.CreatedAt,
		input.Text,
		input.EditorJSJSON,
		input.DiscordChannelID,
		input.DiscordAttachmentKeys,
		input.UpdatedAt,
	).Scan(
		&entity.ID,
		&entity.UserID,
		&entity.Text,
		&entity.EditorJSJSON,
		&entity.CreatedAt,
		&entity.UpdatedAt,
		&entity.DiscordMessageID,
		&entity.DiscordChannelID,
		&entity.DiscordAttachmentKeys,
		&previousAttachmentKeys,
		&created,
	)
	if err != nil {
		return notifications.UpsertDiscordResult{}, err
	}

	return notifications.UpsertDiscordResult{
		Notification:           entity,
		PreviousAttachmentKeys: previousAttachmentKeys,
		Created:                created,
	}, nil
}

func (p *Pgx) DeleteDiscord(
	ctx context.Context,
	channelID string,
	messageIDs []string,
) ([]string, error) {
	if len(messageIDs) == 0 {
		return nil, nil
	}

	const query = `
		DELETE FROM notifications
		WHERE discord_channel_id = $1
			AND discord_message_id = ANY($2)
		RETURNING discord_attachment_keys
	`

	rows, err := p.pool.Query(ctx, query, channelID, messageIDs)
	if err != nil {
		return nil, err
	}

	attachmentKeyGroups, err := pgx.CollectRows(rows, pgx.RowTo[[]string])
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, err
	}

	var attachmentKeys []string
	for _, keys := range attachmentKeyGroups {
		attachmentKeys = append(attachmentKeys, keys...)
	}

	return attachmentKeys, nil
}
