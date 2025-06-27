package clickhouse

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/twirapp/twir/libs/baseapp"
	"github.com/twirapp/twir/libs/repositories/chat_messages"
	"github.com/twirapp/twir/libs/repositories/chat_messages/model"
)

type Opts struct {
	Client *baseapp.ClickhouseClient
}

func New(opts Opts) *Clickhouse {
	return &Clickhouse{
		client: opts.Client,
	}
}

func NewFx(client *baseapp.ClickhouseClient) *Clickhouse {
	return New(Opts{Client: client})
}

var _ chat_messages.Repository = (*Clickhouse)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

type Clickhouse struct {
	client *baseapp.ClickhouseClient
}

func (c *Clickhouse) Create(ctx context.Context, input chat_messages.CreateInput) error {
	query := `
INSERT INTO chat_messages (id, channel_id, user_id, text, user_name, user_display_name, user_color)
VALUES (?, ?, ?, ?, ?, ?, ?);
`

	err := c.client.Exec(
		ctx,
		query,
		input.ID,
		input.ChannelID,
		input.UserID,
		input.Text,
		input.UserName,
		input.UserDisplayName,
		input.UserColor,
	)
	return err
}

const batchSize = 1000

func (c *Clickhouse) CreateMany(ctx context.Context, inputs []chat_messages.CreateInput) error {
	if len(inputs) == 0 {
		return nil
	}

	for i := 0; i < len(inputs); i += batchSize {
		end := i + batchSize
		if end > len(inputs) {
			end = len(inputs)
		}
		if err := c.createBatch(ctx, inputs[i:end]); err != nil {
			return err
		}
	}

	return nil
}

func (c *Clickhouse) createBatch(ctx context.Context, input []chat_messages.CreateInput) error {
	if len(input) == 0 {
		return nil
	}

	batch, err := c.client.PrepareBatch(
		ctx,
		"INSERT INTO chat_messages (id, channel_id, user_id, text, user_name, user_display_name, user_color)",
	)
	if err != nil {
		return fmt.Errorf("prepare batch failed: %w", err)
	}

	for _, i := range input {
		err := batch.Append(
			i.ID,
			i.ChannelID,
			i.UserID,
			i.Text,
			i.UserName,
			i.UserDisplayName,
			i.UserColor,
		)
		if err != nil {
			return fmt.Errorf("append to batch failed: %w", err)
		}
	}

	err = batch.Send()
	if err != nil {
		return fmt.Errorf("send batch failed: %w", err)
	}

	return nil
}

func (c *Clickhouse) GetMany(
	ctx context.Context,
	input chat_messages.GetManyInput,
) ([]model.ChatMessage, error) {
	perPage := input.PerPage
	if perPage == 0 {
		perPage = 20
	}

	if perPage > 1000 {
		perPage = 20
	}

	builder := sq.Select(
		"id",
		"channel_id",
		"user_id",
		"user_name",
		"user_display_name",
		"user_color",
		"text",
		"created_at",
		"updated_at",
	).From("chat_messages")

	if input.ChannelID != nil {
		builder = builder.Where(squirrel.Eq{"channel_id": *input.ChannelID})
	}

	if input.UserNameLike != nil && *input.UserNameLike != "" {
		builder = builder.Where(squirrel.ILike{"user_name": fmt.Sprintf("%%%s%", *input.UserNameLike)})
	}

	if input.TextLike != nil && *input.TextLike != "" {
		builder = builder.Where(squirrel.ILike{"text": fmt.Sprintf("%%%s%%", *input.TextLike)})
	}

	if len(input.UserIDs) > 0 {
		builder = builder.Where(squirrel.Eq{"user_id": input.UserIDs})
	}

	builder = builder.OrderBy("created_at DESC").
		Offset(uint64(input.Page * perPage)).
		Limit(uint64(perPage))

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("failed to build query: %w", err)
	}

	rows, err := c.client.Query(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query: %w", err)
	}
	defer rows.Close()

	var result []model.ChatMessage
	for rows.Next() {
		var m model.ChatMessage
		err := rows.Scan(
			&m.ID,
			&m.ChannelID,
			&m.UserID,
			&m.UserName,
			&m.UserDisplayName,
			&m.UserColor,
			&m.Text,
			&m.CreatedAt,
			&m.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan: %w", err)
		}

		result = append(result, m)
	}

	if rows.Err() != nil {
		return nil, fmt.Errorf("failed to iterate rows: %w", rows.Err())
	}

	return result, nil
}
