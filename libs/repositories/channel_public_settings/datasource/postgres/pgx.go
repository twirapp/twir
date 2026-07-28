package postgres

import (
	"context"
	"errors"
	"fmt"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	channelpublicsettings "github.com/twirapp/twir/libs/entities/channel_public_settings"
	channelpublicsettingsrepo "github.com/twirapp/twir/libs/repositories/channel_public_settings"
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

var _ channelpublicsettingsrepo.Repository = (*Pgx)(nil)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

type settingsModel struct {
	ID          uuid.UUID
	ChannelID   uuid.UUID
	Description *string
}

type socialLinkModel struct {
	ID    uuid.UUID
	Title string
	Href  string
}

func (c *Pgx) GetByChannelID(
	ctx context.Context,
	channelID uuid.UUID,
) (channelpublicsettings.ChannelPublicSettings, error) {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)

	settingsRows, err := conn.Query(
		ctx,
		`
SELECT id, channel_id, description
FROM channels_public_settings
WHERE channel_id = $1
LIMIT 1
`,
		channelID,
	)
	if err != nil {
		return channelpublicsettings.Nil, fmt.Errorf(
			"failed to query channel_public_settings: %w",
			err,
		)
	}

	settings, err := pgx.CollectExactlyOneRow(
		settingsRows,
		pgx.RowToStructByName[settingsModel],
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return channelpublicsettings.Nil, channelpublicsettingsrepo.ErrNotFound
		}
		return channelpublicsettings.Nil, fmt.Errorf(
			"failed to collect channel_public_settings: %w",
			err,
		)
	}

	linksRows, err := conn.Query(
		ctx,
		`
SELECT id, title, href
FROM channels_public_settings_links
WHERE settings_id = $1
ORDER BY id
`,
		settings.ID,
	)
	if err != nil {
		return channelpublicsettings.Nil, fmt.Errorf(
			"failed to query channel_public_settings_links: %w",
			err,
		)
	}

	linksModels, err := pgx.CollectRows(
		linksRows,
		pgx.RowToStructByName[socialLinkModel],
	)
	if err != nil {
		return channelpublicsettings.Nil, fmt.Errorf(
			"failed to collect channel_public_settings_links: %w",
			err,
		)
	}

	socialLinks := make([]channelpublicsettings.SocialLink, 0, len(linksModels))
	for _, link := range linksModels {
		socialLinks = append(socialLinks, channelpublicsettings.SocialLink{
			ID:    link.ID,
			Title: link.Title,
			Href:  link.Href,
		})
	}

	return channelpublicsettings.ChannelPublicSettings{
		ID:          settings.ID,
		ChannelID:   settings.ChannelID,
		Description: settings.Description,
		SocialLinks: socialLinks,
	}, nil
}

func (c *Pgx) Upsert(ctx context.Context, input channelpublicsettingsrepo.UpsertInput) error {
	conn := c.getter.DefaultTrOrDB(ctx, c.pool)

	query := `
INSERT INTO channels_public_settings (id, channel_id, description)
VALUES (uuidv7(), $1, $2)
ON CONFLICT (channel_id) DO UPDATE
SET description = CASE
	WHEN $3::boolean THEN $2
	ELSE channels_public_settings.description
END
RETURNING id
`

	var settingsID uuid.UUID
	err := conn.QueryRow(
		ctx,
		query,
		input.ChannelID,
		input.Description,
		input.DescriptionSet,
	).Scan(&settingsID)
	if err != nil {
		return fmt.Errorf("failed to upsert channel_public_settings: %w", err)
	}

	if input.SocialLinksSet {
		_, err = conn.Exec(
			ctx,
			`DELETE FROM channels_public_settings_links WHERE settings_id = $1`,
			settingsID,
		)
		if err != nil {
			return fmt.Errorf("failed to delete old social links: %w", err)
		}

		for _, link := range input.SocialLinks {
			_, err = conn.Exec(
				ctx,
				`INSERT INTO channels_public_settings_links (id, settings_id, title, href) VALUES (uuidv7(), $1, $2, $3)`,
				settingsID,
				link.Title,
				link.Href,
			)
			if err != nil {
				return fmt.Errorf("failed to insert social link: %w", err)
			}
		}
	}

	return nil
}
