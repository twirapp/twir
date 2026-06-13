-- +goose Up
-- +goose StatementBegin
-- userId is already UUID, just rename to snake_case and add FK
ALTER TABLE channels_roles_users RENAME COLUMN "userId" TO user_id;

ALTER TABLE channels_roles_users
	ADD CONSTRAINT channels_roles_users_user_id_fkey
		FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'roles_users_userid_uuid is not reversible';
-- +goose StatementEnd
