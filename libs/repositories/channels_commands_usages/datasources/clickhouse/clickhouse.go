package clickhouse

import (
	"context"

	"github.com/Masterminds/squirrel"
	twirclickhouse "github.com/twirapp/twir/libs/baseapp/clickhouse"
	channelscommandsusages "github.com/twirapp/twir/libs/repositories/channels_commands_usages"
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

var _ channelscommandsusages.Repository = (*Clickhouse)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Question)

type Clickhouse struct {
	client *twirclickhouse.ClickhouseClient
}

func (c *Clickhouse) Count(ctx context.Context, input channelscommandsusages.CountInput) (
	uint64,
	error,
) {
	selectBuilder := sq.Select("COUNT(*)").From("channels_commands_usages")

	if input.ChannelID != nil && *input.ChannelID != "" {
		selectBuilder = selectBuilder.Where(squirrel.Eq{"channel_id": *input.ChannelID})
	}

	if input.UserID != nil && *input.UserID != "" {
		selectBuilder = selectBuilder.Where(squirrel.Eq{"user_id": *input.UserID})
	}

	if input.CommandID != nil {
		selectBuilder = selectBuilder.Where(squirrel.Eq{"command_id": *input.CommandID})
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

func (c *Clickhouse) Create(ctx context.Context, input channelscommandsusages.CreateInput) error {
	query := `
INSERT INTO channels_commands_usages (channel_id, user_id, command_id)
VALUES (?, ?, ?);
`

	err := c.client.Exec(ctx, query, input.ChannelID, input.UserID, input.CommandID)
	return err
}
