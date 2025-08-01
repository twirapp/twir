package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upDonatepaySplitTable, downDonatepaySplitTable)
}

type donatePaySplitTableChannelsIntegrations20250801200912 struct {
	ID        string
	ApiKey    string
	ChannelID string
}

func upDonatepaySplitTable(ctx context.Context, tx *sql.Tx) error {
	createQuery := `
CREATE TABLE channels_integrations_donatepay (
	id ulid PRIMARY KEY DEFAULT gen_ulid(),
	channel_id text NOT NULL REFERENCES channels(id),
	api_key text,
	enabled boolean NOT NULL DEFAULT false
);

CREATE UNIQUE INDEX channels_integrations_donatepay_channel_id_key on channels_integrations_donatepay(channel_id);
`

	if _, err := tx.ExecContext(ctx, createQuery); err != nil {
		return err
	}

	var donatePayIntegrationId string
	if err := tx.QueryRowContext(
		ctx,
		"SELECT id from integrations WHERE service = 'DONATEPAY'",
	).Scan(&donatePayIntegrationId); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		return err
	}

	rows, err := tx.QueryContext(
		ctx,
		`SELECT id, "apiKey", "channelId" from channels_integrations WHERE "integrationId" = $1 AND enabled = $2 AND "apiKey" IS NOT NULL`,
		donatePayIntegrationId,
		true,
	)
	if err != nil {
		return err
	}

	defer rows.Close()

	var channelsIntegrations []donatePaySplitTableChannelsIntegrations20250801200912
	for rows.Next() {
		integration := donatePaySplitTableChannelsIntegrations20250801200912{}

		if err := rows.Scan(&integration.ID, &integration.ApiKey, &integration.ChannelID); err != nil {
			return err
		}

		channelsIntegrations = append(channelsIntegrations, integration)
	}

	if rows.Err() != nil {
		return err
	}

	insertQuery := `INSERT INTO channels_integrations_donatepay(channel_id, enabled, api_key) VALUES ($1, $2, $3)`
	for _, integration := range channelsIntegrations {
		if _, err := tx.ExecContext(
			ctx,
			insertQuery,
			integration.ChannelID,
			true,
			integration.ApiKey,
		); err != nil {
			return err
		}
	}

	// This code is executed when the migration is applied.
	return nil
}

func downDonatepaySplitTable(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}
