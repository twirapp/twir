-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
ALTER TYPE "channels_modules_settings_type_enum" ADD VALUE 'be_right_back_overlay';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
