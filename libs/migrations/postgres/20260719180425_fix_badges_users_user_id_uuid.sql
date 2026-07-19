-- +goose Up
-- +goose StatementBegin
ALTER TABLE badges_users
	DROP CONSTRAINT IF EXISTS badges_users_user_id_fkey;

DO $$
BEGIN
	IF EXISTS (
		SELECT 1
		FROM information_schema.columns
		WHERE table_schema = 'public'
			AND table_name = 'badges_users'
			AND column_name = 'user_id'
			AND data_type = 'text'
	) THEN
		EXECUTE $update$
			UPDATE badges_users bu
			SET user_id = u.id::text
			FROM users u
			WHERE u.platform = 'twitch'
				AND u.platform_id = bu.user_id
		$update$;

		EXECUTE 'ALTER TABLE badges_users ALTER COLUMN user_id TYPE UUID USING user_id::uuid';
	END IF;
END $$;

ALTER TABLE badges_users
	ADD CONSTRAINT badges_users_user_id_fkey
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'fix_badges_users_user_id_uuid is not reversible';
-- +goose StatementEnd
