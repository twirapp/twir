-- +goose Up
-- +goose StatementBegin
ALTER TABLE tokens ADD COLUMN "deviceID" TEXT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE tokens DROP COLUMN "deviceID";
-- +goose StatementEnd
