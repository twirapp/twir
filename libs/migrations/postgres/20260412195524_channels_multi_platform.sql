-- +goose Up
-- +goose StatementBegin
CREATE TEMP TABLE tmp_channel_fk_constraints ON COMMIT DROP AS
SELECT
    c.conrelid::regclass::text AS table_name,
    a.attname AS column_name,
    c.conname AS constraint_name,
    a.attnotnull AS is_not_null,
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
JOIN pg_attribute a
    ON a.attrelid = c.conrelid
   AND a.attnum = c.conkey[1]
WHERE c.contype = 'f'
  AND c.confrelid = 'channels'::regclass
  AND array_length(c.conkey, 1) = 1;

DO $$
DECLARE
    row RECORD;
    channels_pkey_name TEXT;
    channels_user_fk_name TEXT;
BEGIN
    FOR row IN SELECT * FROM tmp_channel_fk_constraints LOOP
        EXECUTE format(
            'ALTER TABLE %I DROP CONSTRAINT %I',
            row.table_name,
            row.constraint_name
        );
    END LOOP;

    SELECT conname
    INTO channels_user_fk_name
    FROM pg_constraint
    WHERE conrelid = 'channels'::regclass
      AND contype = 'f'
      AND confrelid = 'users'::regclass
    LIMIT 1;

    IF channels_user_fk_name IS NOT NULL THEN
        EXECUTE format('ALTER TABLE channels DROP CONSTRAINT %I', channels_user_fk_name);
    END IF;

    SELECT conname
    INTO channels_pkey_name
    FROM pg_constraint
    WHERE conrelid = 'channels'::regclass
      AND contype = 'p'
    LIMIT 1;

    EXECUTE format('ALTER TABLE channels DROP CONSTRAINT %I', channels_pkey_name);
    ALTER TABLE channels RENAME COLUMN id TO user_id;

    ALTER TABLE channels ADD COLUMN id UUID;

    UPDATE channels
    SET id = gen_random_uuid()
    WHERE id IS NULL;

    ALTER TABLE channels
        ALTER COLUMN id SET NOT NULL,
        ALTER COLUMN id SET DEFAULT gen_random_uuid();

    ALTER TABLE channels ADD PRIMARY KEY (id);

    ALTER TABLE channels ADD COLUMN platform platform;

    UPDATE channels
    SET platform = 'twitch'::platform
    WHERE platform IS NULL;

    ALTER TABLE channels
        ALTER COLUMN platform SET NOT NULL,
        ALTER COLUMN platform SET DEFAULT 'twitch'::platform,
        ALTER COLUMN user_id SET NOT NULL;

    ALTER TABLE channels
        ADD CONSTRAINT channels_user_id_fkey
            FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;

    ALTER TABLE channels
        ADD CONSTRAINT channels_user_platform_unique UNIQUE (user_id, platform);

    CREATE TEMP TABLE tmp_channel_id_map ON COMMIT DROP AS
    SELECT user_id AS old_channel_id, id AS new_channel_id
    FROM channels;

    FOR row IN SELECT * FROM tmp_channel_fk_constraints LOOP
        EXECUTE format(
            'UPDATE %I AS child '
            || 'SET %I = mapping.new_channel_id '
            || 'FROM tmp_channel_id_map AS mapping '
            || 'WHERE child.%I = mapping.old_channel_id',
            row.table_name,
            row.column_name,
            row.column_name
        );

        IF row.is_not_null THEN
            EXECUTE format(
                'ALTER TABLE %I ALTER COLUMN %I SET NOT NULL',
                row.table_name,
                row.column_name
            );
        END IF;

        EXECUTE format(
            'ALTER TABLE %I ADD CONSTRAINT %I FOREIGN KEY (%I) REFERENCES channels(id) ON UPDATE %s ON DELETE %s',
            row.table_name,
            row.constraint_name,
            row.column_name,
            row.on_update,
            row.on_delete
        );
    END LOOP;

    ALTER TABLE channels DROP COLUMN twitch_id;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'channels_multi_platform is not reversible';
-- +goose StatementEnd
