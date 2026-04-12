-- +goose Up
-- +goose StatementBegin
ALTER TABLE users ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE user_platform_accounts ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE kick_bots ALTER COLUMN id SET DEFAULT uuidv7();
ALTER TABLE channels ALTER COLUMN id SET DEFAULT uuidv7();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE kick_bots ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE user_platform_accounts ALTER COLUMN id SET DEFAULT gen_random_uuid();
ALTER TABLE users ALTER COLUMN id SET DEFAULT gen_random_uuid();
-- +goose StatementEnd
