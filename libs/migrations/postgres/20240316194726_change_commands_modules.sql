-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

UPDATE channels_commands SET module = 'DUDES' WHERE "defaultName" in ('dudes color', 'dudes grow', 'dudes leave', 'dudes sprite', 'jump');
UPDATE channels_commands SET module = 'OVERLAYS' WHERE "defaultName" = 'kappagen';
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
