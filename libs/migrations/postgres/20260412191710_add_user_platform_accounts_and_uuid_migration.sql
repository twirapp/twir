-- +goose Up
-- +goose StatementBegin
ALTER TABLE users
    ADD COLUMN internal_id UUID DEFAULT uuidv7();

UPDATE users
SET internal_id = uuidv7()
WHERE internal_id IS NULL;

ALTER TABLE users
    ALTER COLUMN internal_id SET NOT NULL;

ALTER TABLE users
    ADD CONSTRAINT users_internal_id_unique UNIQUE (internal_id);

CREATE TEMP TABLE tmp_user_id_map ON COMMIT DROP AS
SELECT id AS old_user_id, internal_id AS new_user_id
FROM users;

CREATE FUNCTION pg_temp.map_user_id(old_id TEXT) RETURNS UUID
    LANGUAGE SQL
    STABLE
AS $$
    SELECT new_user_id
    FROM tmp_user_id_map
    WHERE old_user_id = old_id
$$;

CREATE TEMP TABLE tmp_user_fk_constraints ON COMMIT DROP AS
SELECT
    c.conrelid::regclass::text AS table_name,
    a.attname AS column_name,
    c.conname AS constraint_name,
    CASE c.confdeltype
        WHEN 'a' THEN 'NO ACTION'
        WHEN 'r' THEN 'RESTRICT'
        WHEN 'c' THEN 'CASCADE'
        WHEN 'n' THEN 'SET NULL'
        WHEN 'd' THEN 'SET DEFAULT'
    END AS on_delete,
    CASE c.confupdtype
        WHEN 'a' THEN 'NO ACTION'
        WHEN 'r' THEN 'RESTRICT'
        WHEN 'c' THEN 'CASCADE'
        WHEN 'n' THEN 'SET NULL'
        WHEN 'd' THEN 'SET DEFAULT'
    END AS on_update
FROM pg_constraint c
JOIN unnest(c.conkey) WITH ORDINALITY AS cols(attnum, ord) ON TRUE
JOIN pg_attribute a ON a.attrelid = c.conrelid AND a.attnum = cols.attnum
WHERE c.contype = 'f'
  AND c.confrelid = 'users'::regclass;

CREATE TEMP TABLE tmp_channel_fk_constraints ON COMMIT DROP AS
SELECT
    c.conrelid::regclass::text AS table_name,
    a.attname AS column_name,
    c.conname AS constraint_name,
    CASE c.confdeltype
        WHEN 'a' THEN 'NO ACTION'
        WHEN 'r' THEN 'RESTRICT'
        WHEN 'c' THEN 'CASCADE'
        WHEN 'n' THEN 'SET NULL'
        WHEN 'd' THEN 'SET DEFAULT'
    END AS on_delete,
    CASE c.confupdtype
        WHEN 'a' THEN 'NO ACTION'
        WHEN 'r' THEN 'RESTRICT'
        WHEN 'c' THEN 'CASCADE'
        WHEN 'n' THEN 'SET NULL'
        WHEN 'd' THEN 'SET DEFAULT'
    END AS on_update
FROM pg_constraint c
JOIN unnest(c.conkey) WITH ORDINALITY AS cols(attnum, ord) ON TRUE
JOIN pg_attribute a ON a.attrelid = c.conrelid AND a.attnum = cols.attnum
WHERE c.contype = 'f'
  AND c.confrelid = 'channels'::regclass;

DO $$
DECLARE
    row RECORD;
    users_pkey_name TEXT;
BEGIN
    FOR row IN SELECT * FROM tmp_user_fk_constraints LOOP
        EXECUTE format(
            'ALTER TABLE %I DROP CONSTRAINT %I',
            row.table_name,
            row.constraint_name
        );
    END LOOP;

    FOR row IN SELECT * FROM tmp_channel_fk_constraints LOOP
        EXECUTE format(
            'ALTER TABLE %I DROP CONSTRAINT %I',
            row.table_name,
            row.constraint_name
        );
    END LOOP;

    ALTER TABLE channels
        ADD COLUMN twitch_id TEXT;

    UPDATE channels
    SET twitch_id = id
    WHERE twitch_id IS NULL;

    ALTER TABLE channels
        ALTER COLUMN twitch_id SET NOT NULL;

    ALTER TABLE channels
        ALTER COLUMN id TYPE UUID USING pg_temp.map_user_id(id);

    FOR row IN
        SELECT *
        FROM tmp_user_fk_constraints
        WHERE NOT (table_name = 'channels' AND column_name = 'id')
    LOOP
        EXECUTE format(
            'ALTER TABLE %I ALTER COLUMN %I TYPE UUID USING pg_temp.map_user_id(%I)',
            row.table_name,
            row.column_name,
            row.column_name
        );
    END LOOP;

    FOR row IN SELECT * FROM tmp_channel_fk_constraints LOOP
        EXECUTE format(
            'ALTER TABLE %I ALTER COLUMN %I TYPE UUID USING pg_temp.map_user_id(%I)',
            row.table_name,
            row.column_name,
            row.column_name
        );
    END LOOP;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'tokens'
          AND column_name = 'user_id'
    ) THEN
        UPDATE tokens t
        SET user_id = u.internal_id
        FROM users u
        WHERE u."tokenId" = t.id
          AND t.user_id IS NULL;
    ELSE
        ALTER TABLE tokens ADD COLUMN user_id UUID;

        UPDATE tokens t
        SET user_id = u.internal_id
        FROM users u
        WHERE u."tokenId" = t.id;
    END IF;

    IF EXISTS (
        SELECT 1
        FROM information_schema.columns
        WHERE table_schema = 'public'
          AND table_name = 'tokens'
          AND column_name = 'user_id'
          AND udt_name <> 'uuid'
    ) THEN
        IF EXISTS (
            SELECT 1
            FROM pg_constraint
            WHERE conrelid = 'tokens'::regclass
              AND conname = 'tokens_user_id_fkey'
        ) THEN
            ALTER TABLE tokens DROP CONSTRAINT tokens_user_id_fkey;
        END IF;

        ALTER TABLE tokens
            ALTER COLUMN user_id TYPE UUID USING pg_temp.map_user_id(user_id);
    END IF;

    SELECT conname
    INTO users_pkey_name
    FROM pg_constraint
    WHERE conrelid = 'users'::regclass
      AND contype = 'p';

    EXECUTE format('ALTER TABLE users DROP CONSTRAINT %I', users_pkey_name);

    ALTER TABLE users DROP CONSTRAINT users_internal_id_unique;

    ALTER TABLE users RENAME COLUMN id TO twitch_id;
    ALTER TABLE users RENAME COLUMN internal_id TO id;
    ALTER TABLE users ADD PRIMARY KEY (id);

    CREATE INDEX IF NOT EXISTS tokens_user_id_idx ON tokens(user_id);

    ALTER TABLE tokens
        ADD CONSTRAINT tokens_user_id_fkey
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

    CREATE TABLE user_platform_accounts (
        id UUID PRIMARY KEY DEFAULT uuidv7(),
        user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
        platform platform NOT NULL,
        platform_user_id TEXT NOT NULL,
        platform_login TEXT NOT NULL,
        platform_display_name TEXT NOT NULL DEFAULT '',
        platform_avatar TEXT NOT NULL DEFAULT '',
        access_token TEXT NOT NULL,
        refresh_token TEXT NOT NULL,
        scopes TEXT[] NOT NULL DEFAULT '{}',
        expires_in INT NOT NULL DEFAULT 0,
        obtainment_timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
        UNIQUE(platform, platform_user_id)
    );

    CREATE INDEX idx_upa_user_id ON user_platform_accounts(user_id);
    CREATE INDEX idx_upa_platform ON user_platform_accounts(platform, platform_user_id);

    INSERT INTO user_platform_accounts (
        user_id,
        platform,
        platform_user_id,
        platform_login,
        platform_display_name,
        platform_avatar,
        access_token,
        refresh_token,
        scopes,
        expires_in,
        obtainment_timestamp
    )
    SELECT
        t.user_id,
        'twitch'::platform,
        u.twitch_id,
        COALESCE(NULLIF(u.twitch_id, ''), u.twitch_id),
        '',
        '',
        t."accessToken",
        t."refreshToken",
        COALESCE(t.scopes, '{}'),
        t."expiresIn",
        t."obtainmentTimestamp"
    FROM tokens t
    JOIN users u ON u.id = t.user_id
    ON CONFLICT (platform, platform_user_id) DO NOTHING;

    FOR row IN SELECT * FROM tmp_user_fk_constraints LOOP
        EXECUTE format(
            'ALTER TABLE %I ADD CONSTRAINT %I FOREIGN KEY (%I) REFERENCES users(id) ON UPDATE %s ON DELETE %s',
            row.table_name,
            row.constraint_name,
            row.column_name,
            row.on_update,
            row.on_delete
        );
    END LOOP;

    FOR row IN SELECT * FROM tmp_channel_fk_constraints LOOP
        EXECUTE format(
            'ALTER TABLE %I ADD CONSTRAINT %I FOREIGN KEY (%I) REFERENCES channels(id) ON UPDATE %s ON DELETE %s',
            row.table_name,
            row.constraint_name,
            row.column_name,
            row.on_update,
            row.on_delete
        );
    END LOOP;

END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'add_user_platform_accounts_and_uuid_migration is not reversible';
-- +goose StatementEnd
