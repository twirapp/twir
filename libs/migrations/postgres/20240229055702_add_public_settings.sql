-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

CREATE TABLE "channels_public_settings" (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	channel_id text NOT NULL,
	description varchar(1000),

	FOREIGN KEY (channel_id) REFERENCES channels (id)
);

CREATE TABLE "channels_public_settings_links"
(
	id          TEXT PRIMARY KEY DEFAULT uuid_generate_v4(),
	settings_id uuid         NOT NULL,
	href        varchar(500) NOT NULL,
	title       varchar(30)  NOT NULL,

	FOREIGN KEY (settings_id) REFERENCES channels_public_settings (id) ON DELETE CASCADE
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
