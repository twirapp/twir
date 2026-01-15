package clickhouse

import (
	"context"

	twirclickhouse "github.com/twirapp/twir/libs/baseapp/clickhouse"
	"github.com/twirapp/twir/libs/repositories/short_links_views"
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

var _ short_links_views.Repository = (*Clickhouse)(nil)

type Clickhouse struct {
	client *twirclickhouse.ClickhouseClient
}

func (c *Clickhouse) Create(ctx context.Context, input short_links_views.CreateInput) error {
	query := `
INSERT INTO short_links_views (short_link_id, user_id, ip, user_agent, created_at)
VALUES (?, ?, ?, ?, ?);
`

	err := c.client.Exec(
		ctx,
		query,
		input.ShortLinkID,
		input.UserID,
		input.IP,
		input.UserAgent,
		input.CreatedAt,
	)

	return err
}
