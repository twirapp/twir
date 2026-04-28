-- +goose Up
-- +goose StatementBegin
-- Add platform and platform_id columns to users
ALTER TABLE users
    ADD COLUMN platform  platform NOT NULL DEFAULT 'twitch',
    ADD COLUMN platform_id TEXT    NOT NULL DEFAULT '';

-- Backfill platform_id from old TEXT id (Twitch user ID)
UPDATE users SET platform_id = id;

-- Remove defaults now that data is populated
ALTER TABLE users
    ALTER COLUMN platform_id DROP DEFAULT;

-- Add new UUID column for internal id
ALTER TABLE users
    ADD COLUMN new_id UUID DEFAULT uuidv7();

-- Ensure every row has a new UUID
UPDATE users SET new_id = uuidv7() WHERE new_id IS NULL;

ALTER TABLE users ALTER COLUMN new_id SET NOT NULL;

-- Need a unique constraint on new_id before we can make it PK
-- and before FK tables can reference it
ALTER TABLE users ADD CONSTRAINT users_new_id_unique UNIQUE (new_id);

-- Capture FK constraints referencing users(id) for re-attachment later
CREATE TEMP TABLE tmp_users_fk ON COMMIT DROP AS
SELECT
    c.conrelid::regclass::text AS tbl,
    a.attname                  AS col,
    c.conname                  AS cname,
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

-- Build old_id -> new_id mapping
CREATE TEMP TABLE tmp_user_id_map ON COMMIT DROP AS
SELECT id AS old_id, new_id FROM users;

CREATE TEMP TABLE tmp_channels_downstream_fk ON COMMIT DROP AS
SELECT
    c.conrelid::regclass::text AS tbl,
    a.attname                  AS col,
    c.conname                  AS cname,
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
    r RECORD;
    users_pkey TEXT;
BEGIN
    FOR r IN SELECT * FROM tmp_users_fk LOOP
        EXECUTE format('ALTER TABLE %I DROP CONSTRAINT %I', r.tbl, r.cname);
    END LOOP;

    FOR r IN SELECT * FROM tmp_channels_downstream_fk LOOP
        EXECUTE format('ALTER TABLE %I DROP CONSTRAINT %I', r.tbl, r.cname);
    END LOOP;

    FOR r IN SELECT * FROM tmp_channels_downstream_fk LOOP
        EXECUTE format(
            'UPDATE %I SET %I = m.new_id::text FROM tmp_user_id_map m WHERE m.old_id = %I.%I::text',
            r.tbl, r.col, r.tbl, r.col
        );
        EXECUTE format(
            'ALTER TABLE %I ALTER COLUMN %I TYPE UUID USING %I::uuid',
            r.tbl, r.col, r.col
        );
    END LOOP;

    FOR r IN SELECT * FROM tmp_users_fk LOOP
        EXECUTE format(
            'UPDATE %I SET %I = m.new_id::text FROM tmp_user_id_map m WHERE m.old_id = %I.%I::text',
            r.tbl, r.col, r.tbl, r.col
        );
        EXECUTE format(
            'ALTER TABLE %I ALTER COLUMN %I TYPE UUID USING %I::uuid',
            r.tbl, r.col, r.col
        );
    END LOOP;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'tokens' AND column_name = 'user_id'
        AND udt_name <> 'uuid'
    ) THEN
        EXECUTE 'UPDATE tokens SET user_id = m.new_id::text FROM tmp_user_id_map m WHERE m.old_id = tokens.user_id::text';
        EXECUTE 'ALTER TABLE tokens ALTER COLUMN user_id TYPE UUID USING user_id::uuid';
    END IF;

    SELECT conname INTO users_pkey FROM pg_constraint
    WHERE conrelid = 'users'::regclass AND contype = 'p';
    EXECUTE format('ALTER TABLE users DROP CONSTRAINT %I', users_pkey);

    ALTER TABLE users DROP CONSTRAINT users_new_id_unique;

    ALTER TABLE users DROP COLUMN id;
    ALTER TABLE users RENAME COLUMN new_id TO id;

    ALTER TABLE users ALTER COLUMN id SET DEFAULT uuidv7();
    ALTER TABLE users ADD PRIMARY KEY (id);

    FOR r IN SELECT * FROM tmp_users_fk LOOP
        EXECUTE format(
            'ALTER TABLE %I ADD CONSTRAINT %I FOREIGN KEY (%I) REFERENCES users(id) ON UPDATE %s ON DELETE %s',
            r.tbl, r.cname, r.col, r.on_update, r.on_delete
        );
    END LOOP;

    FOR r IN SELECT * FROM tmp_channels_downstream_fk LOOP
        EXECUTE format(
            'ALTER TABLE %I ADD CONSTRAINT %I FOREIGN KEY (%I) REFERENCES channels(id) ON UPDATE %s ON DELETE %s',
            r.tbl, r.cname, r.col, r.on_update, r.on_delete
        );
    END LOOP;

    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'tokens' AND column_name = 'user_id'
    ) THEN
        IF NOT EXISTS (
            SELECT 1 FROM pg_constraint
            WHERE conrelid = 'tokens'::regclass AND conname = 'tokens_user_id_fkey'
        ) THEN
            ALTER TABLE tokens
                ADD CONSTRAINT tokens_user_id_fkey
                FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
        END IF;
    END IF;
END $$;

-- Add unique constraint on (platform, platform_id) — each platform user maps to exactly one internal user
ALTER TABLE users ADD CONSTRAINT users_platform_platform_id_unique UNIQUE (platform, platform_id);

CREATE INDEX idx_users_platform_id ON users(platform, platform_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'users_multi_platform is not reversible';
-- +goose StatementEnd
