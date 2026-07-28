package postgres

import (
	"context"
	"errors"

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
		pool: opts.PgxPool,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ channelpublicsettingsrepo.Repository = (*Pgx)(nil)

type Pgx struct {
	pool *pgxpool.Pool
}

func (c *Pgx) GetByChannelID(
	ctx context.Context,
	channelID uuid.UUID,
) (channelpublicsettings.ChannelPublicSettings, error) {
	settings := channelpublicsettings.ChannelPublicSettings{}

	err := c.pool.QueryRow(
		ctx,
		`
SELECT id, channel_id, description
FROM channels_public_settings
WHERE channel_id = $1
LIMIT 1;
`,
		channelID,
	).Scan(&settings.ID, &settings.ChannelID, &settings.Description)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return channelpublicsettings.Nil, channelpublicsettingsrepo.ErrNotFound
		}
		return channelpublicsettings.Nil, err
	}

	rows, err := c.pool.Query(
		ctx,
		`
SELECT id, title, href
FROM channels_public_settings_links
WHERE settings_id = $1
ORDER BY id;
`,
		settings.ID,
	)
	if err != nil {
		return channelpublicsettings.Nil, err
	}
	defer rows.Close()

	settings.SocialLinks = make([]channelpublicsettings.SocialLink, 0)
	for rows.Next() {
		var link channelpublicsettings.SocialLink
		if err := rows.Scan(&link.ID, &link.Title, &link.Href); err != nil {
			return channelpublicsettings.Nil, err
		}
		settings.SocialLinks = append(settings.SocialLinks, link)
	}
	if err := rows.Err(); err != nil {
		return channelpublicsettings.Nil, err
	}

	return settings, nil
}

func (c *Pgx) Upsert(ctx context.Context, input channelpublicsettingsrepo.UpsertInput) error {
	tx, err := c.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	var settingsID uuid.UUID
	err = tx.QueryRow(
		ctx,
		`
INSERT INTO channels_public_settings (id, channel_id, description)
VALUES (uuidv7(), $1, $2)
ON CONFLICT (channel_id) DO UPDATE
SET description = CASE
	WHEN $3::boolean THEN $2
	ELSE channels_public_settings.description
END
RETURNING id;
`,
		input.ChannelID,
		input.Description,
		input.DescriptionSet,
	).Scan(&settingsID)
	if err != nil {
		return err
	}

	if input.SocialLinksSet {
		if _, err := tx.Exec(
			ctx,
			`DELETE FROM channels_public_settings_links WHERE settings_id = $1;`,
			settingsID,
		); err != nil {
			return err
		}

		for _, link := range input.SocialLinks {
			if _, err := tx.Exec(
				ctx,
				`
INSERT INTO channels_public_settings_links (id, settings_id, title, href)
VALUES (uuidv7(), $1, $2, $3);
`,
				settingsID,
				link.Title,
				link.Href,
			); err != nil {
				return err
			}
		}
	}

	return tx.Commit(ctx)
}
