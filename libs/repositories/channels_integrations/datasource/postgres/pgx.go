package postgres

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	channelsintegrations "github.com/twirapp/twir/libs/repositories/channels_integrations"
	"github.com/twirapp/twir/libs/repositories/channels_integrations/model"
	integrationsmodel "github.com/twirapp/twir/libs/repositories/integrations/model"
)

var _ channelsintegrations.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

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

func (p *Pgx) conn(ctx context.Context) trmpgx.Tr {
	return p.getter.DefaultTrOrDB(ctx, p.pool)
}

type channelIntegrationRow struct {
	ID            string         `db:"id"`
	ChannelID     string         `db:"channelId"`
	IntegrationID string         `db:"integrationId"`
	Enabled       bool           `db:"enabled"`
	AccessToken   sql.NullString `db:"accessToken"`
	RefreshToken  sql.NullString `db:"refreshToken"`
	UserName      sql.NullString `db:"userName"`
	Avatar        sql.NullString `db:"avatar"`
}

func (p *Pgx) GetByChannelAndService(
	ctx context.Context,
	channelID string,
	service integrationsmodel.Service,
) (model.ChannelIntegration, error) {
	query := `
SELECT
    ci.id,
    ci."channelId",
    ci."integrationId",
    ci.enabled,
    ci."accessToken",
    ci."refreshToken",
    ci.data->>'username' AS "userName",
    ci.data->>'avatar' AS "avatar"
FROM channels_integrations ci
JOIN integrations i ON ci."integrationId" = i.id
WHERE ci."channelId" = $1 AND i.service = $2
LIMIT 1
`

	rows, err := p.conn(ctx).Query(ctx, query, channelID, service)
	if err != nil {
		return model.Nil, fmt.Errorf("GetByChannelAndService: failed to execute query: %w", err)
	}
	defer rows.Close()

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[channelIntegrationRow])
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Nil, nil
		}
		return model.Nil, fmt.Errorf("GetByChannelAndService: failed to collect row: %w", err)
	}

	return p.rowToModel(result), nil
}

func (p *Pgx) Create(
	ctx context.Context,
	input channelsintegrations.CreateInput,
) (model.ChannelIntegration, error) {
	query := sq.Insert("channels_integrations").
		Columns(
			`"channelId"`,
			`"integrationId"`,
			"enabled",
			`"accessToken"`,
			`"refreshToken"`,
			"data",
		).
		Values(
			input.ChannelID,
			input.IntegrationID,
			input.Enabled,
			input.AccessToken,
			input.RefreshToken,
			input.Data,
		).
		Suffix(`RETURNING id, "channelId", "integrationId", enabled, "accessToken", "refreshToken", data->>'username' AS "userName", data->>'avatar' AS "avatar"`)

	sqlStr, args, err := query.ToSql()
	if err != nil {
		return model.Nil, fmt.Errorf("Create: failed to build query: %w", err)
	}

	rows, err := p.conn(ctx).Query(ctx, sqlStr, args...)
	if err != nil {
		return model.Nil, fmt.Errorf("Create: failed to execute query: %w", err)
	}

	result, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[channelIntegrationRow])
	if err != nil {
		return model.Nil, fmt.Errorf("Create: failed to collect row: %w", err)
	}

	return p.rowToModel(result), nil
}

func (p *Pgx) Update(
	ctx context.Context,
	id string,
	input channelsintegrations.UpdateInput,
) error {
	builder := sq.Update("channels_integrations").Where(squirrel.Eq{"id": id})

	if input.Enabled != nil {
		builder = builder.Set("enabled", *input.Enabled)
	}
	if input.AccessToken != nil {
		builder = builder.Set(`"accessToken"`, *input.AccessToken)
	}
	if input.RefreshToken != nil {
		builder = builder.Set(`"refreshToken"`, *input.RefreshToken)
	}
	if input.Data != nil {
		data := map[string]any{}
		if input.Data.UserName != nil {
			data["username"] = *input.Data.UserName
		}
		if input.Data.Avatar != nil {
			data["avatar"] = *input.Data.Avatar
		}
		dataBytes, err := json.Marshal(data)
		if err != nil {
			return fmt.Errorf("Update: failed to marshal data: %w", err)
		}
		builder = builder.Set("data", dataBytes)
	}

	sqlStr, args, err := builder.ToSql()
	if err != nil {
		return fmt.Errorf("Update: failed to build query: %w", err)
	}

	_, err = p.conn(ctx).Exec(ctx, sqlStr, args...)
	if err != nil {
		return fmt.Errorf("Update: failed to execute query: %w", err)
	}

	return nil
}

func (p *Pgx) rowToModel(row channelIntegrationRow) model.ChannelIntegration {
	result := model.ChannelIntegration{
		ID:            row.ID,
		ChannelID:     row.ChannelID,
		IntegrationID: row.IntegrationID,
		Enabled:       row.Enabled,
	}

	if row.AccessToken.Valid {
		result.AccessToken = &row.AccessToken.String
	}
	if row.RefreshToken.Valid {
		result.RefreshToken = &row.RefreshToken.String
	}

	if row.UserName.Valid || row.Avatar.Valid {
		result.Data = &model.Data{}
		if row.UserName.Valid {
			result.Data.UserName = &row.UserName.String
		}
		if row.Avatar.Valid {
			result.Data.Avatar = &row.Avatar.String
		}
	}

	return result
}
