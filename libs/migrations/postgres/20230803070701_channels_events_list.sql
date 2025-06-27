-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';
DROP TABLE "channel_events_donations";
DROP TABLE "channel_events_follows";
DROP TABLE "channel_events_list";
DROP type "channel_events_list_type_enum";

CREATE TYPE "channel_events_list_type_enum" AS ENUM (
    'DONATION',
    'FOLLOW',
    'RAIDED',
    'SUBSCRIBE',
    'RESUBSCRIBE',
    'SUBGIFT',
    'FIRST_USER_MESSAGE',
    'CHAT_CLEAR',
    'REDEMPTION_CREATED'
    );

CREATE TABLE "channels_events_list"
(
    "id"         uuid      default uuid_generate_v4() not null,
    "channel_id" text                                 not null
        constraint "channels_events_list_channel_id_fk" references channels ("id"),
    "user_id"    text,
    "type"       channel_events_list_type_enum        not null,
    "data"       jsonb     default '{}'               not null,
    "created_at" timestamp default now()              not null
);

CREATE INDEX channels_events_list_user_id_index ON channels_events_list ("user_id");

delete
from "channels_integrations"
where "integrationId" = (SELECT id from integrations where service = 'DONATIONALERTS');
delete
from "channels_integrations"
where "integrationId" = (SELECT id from integrations where service = 'STREAMLABS');
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
