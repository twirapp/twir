-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels
    ADD COLUMN new_id UUID DEFAULT uuidv7();

UPDATE channels SET new_id = uuidv7() WHERE new_id IS NULL;

ALTER TABLE channels ALTER COLUMN new_id SET NOT NULL;
ALTER TABLE channels ADD CONSTRAINT channels_new_id_unique UNIQUE (new_id);

ALTER TABLE channels
    ADD COLUMN twitch_user_id UUID NULL REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    ADD COLUMN kick_user_id   UUID NULL REFERENCES users(id) ON DELETE RESTRICT ON UPDATE CASCADE;

UPDATE channels SET twitch_user_id = id;

CREATE TEMP TABLE tmp_channel_fk ON COMMIT DROP AS
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

CREATE TEMP TABLE tmp_channel_id_map ON COMMIT DROP AS
SELECT id AS old_id, new_id FROM channels;

DO $$
DECLARE
    r         RECORD;
    chan_pkey TEXT;
    chan_to_users_fk TEXT;
BEGIN
    FOR r IN SELECT * FROM tmp_channel_fk LOOP
        EXECUTE format('ALTER TABLE %I DROP CONSTRAINT %I', r.tbl, r.cname);
    END LOOP;

    FOR r IN SELECT * FROM tmp_channel_fk LOOP
        EXECUTE format(
            'UPDATE %I SET %I = m.new_id FROM tmp_channel_id_map m WHERE m.old_id = %I.%I',
            r.tbl, r.col, r.tbl, r.col
        );
    END LOOP;

    SELECT conname INTO chan_to_users_fk FROM pg_constraint
    WHERE conrelid = 'channels'::regclass AND contype = 'f'
      AND confrelid = 'users'::regclass
      AND conkey = ARRAY[(SELECT attnum FROM pg_attribute WHERE attrelid = 'channels'::regclass AND attname = 'id')];
    IF chan_to_users_fk IS NOT NULL THEN
        EXECUTE format('ALTER TABLE channels DROP CONSTRAINT %I', chan_to_users_fk);
    END IF;

    SELECT conname INTO chan_pkey FROM pg_constraint
    WHERE conrelid = 'channels'::regclass AND contype = 'p';
    IF chan_pkey IS NOT NULL THEN
        EXECUTE format('ALTER TABLE channels DROP CONSTRAINT %I', chan_pkey);
    END IF;

    ALTER TABLE channels DROP CONSTRAINT channels_new_id_unique;
    ALTER TABLE channels DROP COLUMN id;
    ALTER TABLE channels RENAME COLUMN new_id TO id;
    ALTER TABLE channels ALTER COLUMN id SET DEFAULT uuidv7();
    ALTER TABLE channels ADD PRIMARY KEY (id);

    FOR r IN SELECT * FROM tmp_channel_fk LOOP
        EXECUTE format(
            'ALTER TABLE %I ADD CONSTRAINT %I FOREIGN KEY (%I) REFERENCES channels(id) ON UPDATE %s ON DELETE %s',
            r.tbl, r.cname, r.col, r.on_update, r.on_delete
        );
    END LOOP;
END $$;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'channels_multi_platform is not reversible';
-- +goose StatementEnd
