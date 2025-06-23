package clickhouse

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/twirapp/twir/libs/baseapp"
	channelsemotesusagesrepository "github.com/twirapp/twir/libs/repositories/channels_emotes_usages"
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

var _ channelsemotesusagesrepository.Repository = (*Clickhouse)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

type Clickhouse struct {
	client *baseapp.ClickhouseClient
}

func (c *Clickhouse) createBatch(
	ctx context.Context,
	input []channelsemotesusagesrepository.ChannelEmoteUsageInput,
) error {
	if len(input) == 0 {
		return nil
	}

	batch, err := c.client.PrepareBatch(
		ctx,
		"INSERT INTO channels_emotes_usages(channel_id, user_id, emote)",
	)
	if err != nil {
		return fmt.Errorf("prepare batch failed: %w", err)
	}

	for _, i := range input {
		err := batch.Append(
			i.ChannelID,
			i.UserID,
			i.Emote,
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

const batchSize = 1000

func (c *Clickhouse) CreateMany(
	ctx context.Context,
	inputs []channelsemotesusagesrepository.ChannelEmoteUsageInput,
) error {
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

func (c *Clickhouse) Count(ctx context.Context, input channelsemotesusagesrepository.CountInput) (
	uint64,
	error,
) {
	selectBuilder := sq.Select("COUNT(*)").From("channels_emotes_usages")

	if input.ChannelID != nil && *input.ChannelID != "" {
		selectBuilder = selectBuilder.Where(squirrel.Eq{"channel_id": *input.ChannelID})
	}

	if input.UserID != nil && *input.UserID != "" {
		selectBuilder = selectBuilder.Where(squirrel.Eq{"user_id": *input.UserID})
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return 0, fmt.Errorf("build query error: %w", err)
	}

	var count uint64
	err = c.client.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("query error: %w", err)
	}

	return count, nil
}
