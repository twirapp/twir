package pgx

import (
	"context"
	"time"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/entities/dashboard_widget"
	"github.com/twirapp/twir/libs/repositories/dashboard_widgets"
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

var _ dashboard_widgets.Repository = (*Pgx)(nil)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

type model struct {
	ID         string    `db:"id"`
	ChannelID  string    `db:"channel_id"`
	WidgetID   string    `db:"widget_id"`
	X          int       `db:"x"`
	Y          int       `db:"y"`
	W          int       `db:"w"`
	H          int       `db:"h"`
	MinW       int       `db:"min_w"`
	MinH       int       `db:"min_h"`
	Visible    bool      `db:"visible"`
	StackId    *string   `db:"stack_id"`
	StackOrder int       `db:"stack_order"`
	CreatedAt  time.Time `db:"created_at"`
	UpdatedAt  time.Time `db:"updated_at"`

	isNil bool
}

func (m model) IsNil() bool {
	return m.isNil
}

func (c *Pgx) GetByChannelID(ctx context.Context, channelID string) ([]dashboard_widget.DashboardWidget, error) {
	query := `
SELECT id, channel_id, widget_id, x, y, w, h, min_w, min_h, visible, stack_id, stack_order, created_at, updated_at
FROM channels_dashboard_widgets
WHERE channel_id = $1
ORDER BY widget_id
`

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return nil, err
	}

	models, err := pgx.CollectRows(rows, pgx.RowToStructByName[model])
	if err != nil {
		return nil, err
	}

	result := make([]dashboard_widget.DashboardWidget, 0, len(models))
	for _, m := range models {
		result = append(result, m.toEntity())
	}

	return result, nil
}

func (c *Pgx) UpsertMany(ctx context.Context, channelID string, widgets []dashboard_widget.DashboardWidget) error {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)

	// Use ON CONFLICT to handle upserts properly
	if len(widgets) > 0 {
		query := `
INSERT INTO channels_dashboard_widgets (channel_id, widget_id, x, y, w, h, min_w, min_h, visible, stack_id, stack_order)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
ON CONFLICT (channel_id, widget_id) DO UPDATE SET
	x = EXCLUDED.x,
	y = EXCLUDED.y,
	w = EXCLUDED.w,
	h = EXCLUDED.h,
	min_w = EXCLUDED.min_w,
	min_h = EXCLUDED.min_h,
	visible = EXCLUDED.visible,
	stack_id = EXCLUDED.stack_id,
	stack_order = EXCLUDED.stack_order,
	updated_at = NOW()
`
		for _, widget := range widgets {
			_, err := conn.Exec(
				ctx,
				query,
				channelID,
				widget.WidgetID,
				widget.X,
				widget.Y,
				widget.W,
				widget.H,
				widget.MinW,
				widget.MinH,
				widget.Visible,
				widget.StackId,
				widget.StackOrder,
			)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func (m model) toEntity() dashboard_widget.DashboardWidget {
	return dashboard_widget.DashboardWidget{
		ID:         m.ID,
		ChannelID:  m.ChannelID,
		WidgetID:   m.WidgetID,
		X:          m.X,
		Y:          m.Y,
		W:          m.W,
		H:          m.H,
		MinW:       m.MinW,
		MinH:       m.MinH,
		Visible:    m.Visible,
		StackId:    m.StackId,
		StackOrder: m.StackOrder,
	}
}
