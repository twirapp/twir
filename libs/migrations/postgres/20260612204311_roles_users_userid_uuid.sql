-- +goose Up
-- +goose StatementBegin
-- channels_roles_users.userId still stores Twitch platform IDs (TEXT).
-- Convert it to UUID referencing users.id, matching the users_multi_platform migration.
-- Also rename to snake_case for consistency.

-- 1. Rename old column
ALTER TABLE channels_roles_users RENAME COLUMN "userId" TO old_user_id;

-- 2. Add new UUID column (snake_case)
ALTER TABLE channels_roles_users ADD COLUMN user_id UUID;

-- 3. Populate from users table: old_user_id (platform_id) → users.id (UUID)
UPDATE channels_roles_users cru
SET user_id = u.id
FROM users u
WHERE u.platform_id = cru.old_user_id
  AND u.platform = 'twitch';

-- 4. Drop rows that couldn't be resolved (shouldn't happen, but safety)
DELETE FROM channels_roles_users WHERE user_id IS NULL;

-- 5. Set NOT NULL
ALTER TABLE channels_roles_users ALTER COLUMN user_id SET NOT NULL;

-- 6. Drop old column
ALTER TABLE channels_roles_users DROP COLUMN old_user_id;

-- 7. Add FK constraint
ALTER TABLE channels_roles_users
    ADD CONSTRAINT channels_roles_users_user_id_fkey
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'roles_users_userid_uuid is not reversible';
-- +goose StatementEnd
