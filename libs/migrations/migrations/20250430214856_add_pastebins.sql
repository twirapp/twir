-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE pastebins (
	id varchar(5) PRIMARY KEY NOT NULL, --nanoid
	created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
	content TEXT NOT NULL,
	expire_at TIMESTAMPTZ,

	owner_user_id TEXT REFERENCES users(id) ON DELETE SET NULL
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
