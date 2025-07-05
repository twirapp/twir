package clickhouse

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	twirclickhouse "github.com/twirapp/twir/libs/baseapp/clickhouse"
	"github.com/twirapp/twir/libs/repositories/channels_redemptions_history"
	"github.com/twirapp/twir/libs/repositories/channels_redemptions_history/model"
)

type Opts struct {
	Client *twirclickhouse.ClickhouseClient
}

func New(opts Opts) *Clickhouse {
	return &Clickhouse{
		client: opts.Client,
	}
}

func NewFx(client *twirclickhouse.ClickhouseClient) *Clickhouse {
	return New(Opts{Client: client})
}

var _ channelsredemptionshistory.Repository = (*Clickhouse)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

type Clickhouse struct {
	client *twirclickhouse.ClickhouseClient
}

func (c *Clickhouse) Count(
	ctx context.Context,
	input channelsredemptionshistory.CountInput,
) (uint64, error) {
	selectBuilder := sq.Select("COUNT(*)").From("channels_redemptions_history")

	if input.ChannelID != nil && *input.ChannelID != "" {
		selectBuilder = selectBuilder.Where(squirrel.Eq{"channel_id": *input.ChannelID})
	}

	if len(input.RewardsIDs) > 0 {
		selectBuilder = selectBuilder.Where(squirrel.Eq{"reward_id": input.RewardsIDs})
	}

	query, args, err := selectBuilder.ToSql()
	if err != nil {
		return 0, err
	}

	var count uint64
	err = c.client.QueryRow(ctx, query, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (c *Clickhouse) Create(
	ctx context.Context,
	input channelsredemptionshistory.CreateInput,
) error {
	query := `
INSERT INTO channels_redemptions_history(channel_id, user_id, reward_id, reward_title, reward_prompt, reward_cost)
VALUES (?, ?, ?, ?, ?, ?);
`

	err := c.client.Exec(
		ctx,
		query,
		input.ChannelID,
		input.UserID,
		input.RewardID,
		input.RewardTitle,
		input.RewardPrompt,
		input.RewardCost,
	)
	return err
}

func (c *Clickhouse) createBatch(
	ctx context.Context,
	input []channelsredemptionshistory.CreateInput,
) error {
	if len(input) == 0 {
		return nil
	}

	batch, err := c.client.PrepareBatch(
		ctx,
		"INSERT INTO channels_redemptions_history(channel_id, user_id, reward_id, reward_title, reward_prompt, reward_cost)",
	)
	if err != nil {
		return fmt.Errorf("prepare batch failed: %w", err)
	}

	for _, i := range input {
		err := batch.Append(
			i.ChannelID,
			i.UserID,
			i.RewardID,
			i.RewardTitle,
			i.RewardPrompt,
			i.RewardCost,
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
	inputs []channelsredemptionshistory.CreateInput,
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

func (c *Clickhouse) GetMany(
	ctx context.Context,
	input channelsredemptionshistory.GetManyInput,
) (channelsredemptionshistory.GetManyPayload, error) {
	queryBuilder := sq.Select(
		"channel_id",
		"user_id",
		"reward_id",
		"reward_prompt",
		"reward_title",
		"reward_cost",
		"created_at",
	).From("channels_redemptions_history").
		Where(squirrel.Eq{"channel_id": input.ChannelID}).
		OrderBy("created_at DESC")

	if len(input.UserIDs) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"user_id": input.UserIDs})
	}

	if len(input.RewardsIDs) > 0 {
		queryBuilder = queryBuilder.Where(squirrel.Eq{"reward_id": input.RewardsIDs})
	}

	perPage := input.PerPage
	page := input.Page

	if perPage == 0 || perPage > 1000 {
		perPage = 20
	}

	offset := page * perPage

	queryBuilder = queryBuilder.Limit(uint64(perPage)).Offset(uint64(offset))

	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return channelsredemptionshistory.GetManyPayload{}, err
	}

	rows, err := c.client.Query(ctx, query, args...)
	if err != nil {
		return channelsredemptionshistory.GetManyPayload{}, err
	}
	defer rows.Close()

	var result []model.ChannelsRedemptionHistoryItem
	for rows.Next() {
		var item model.ChannelsRedemptionHistoryItem
		err := rows.Scan(
			&item.ChannelID,
			&item.UserID,
			&item.RewardID,
			&item.RewardPrompt,
			&item.RewardTitle,
			&item.RewardCost,
			&item.CreatedAt,
		)
		if err != nil {
			return channelsredemptionshistory.GetManyPayload{}, err
		}

		result = append(result, item)
	}

	if rows.Err() != nil {
		return channelsredemptionshistory.GetManyPayload{}, err
	}

	totalQuery := `
SELECT COUNT(*) FROM channels_redemptions_history
WHERE channel_id = ?
`

	var total uint64
	err = c.client.QueryRow(ctx, totalQuery, input.ChannelID).Scan(&total)
	if err != nil {
		return channelsredemptionshistory.GetManyPayload{}, err
	}

	return channelsredemptionshistory.GetManyPayload{
		Items: result,
		Total: total,
	}, nil
}
