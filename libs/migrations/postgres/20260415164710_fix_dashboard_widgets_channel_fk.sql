-- +goose Up
-- +goose StatementBegin

UPDATE channels_dashboard_widgets cdw
SET channel_id = c.id
FROM channels c
WHERE c.twitch_user_id = cdw.channel_id;

DELETE FROM channels_dashboard_widgets
WHERE channel_id NOT IN (SELECT id FROM channels);

DROP INDEX IF EXISTS channels_dashboard_widgets_channel_id_idx;
DROP INDEX IF EXISTS channels_dashboard_widgets_stack_id_idx;

ALTER TABLE channels_dashboard_widgets
    DROP CONSTRAINT channels_dashboard_widgets_channel_id_fkey;

ALTER TABLE channels_dashboard_widgets
    ADD CONSTRAINT channels_dashboard_widgets_channel_id_fkey
        FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE CASCADE;

CREATE INDEX channels_dashboard_widgets_channel_id_idx
    ON channels_dashboard_widgets (channel_id);

CREATE INDEX channels_dashboard_widgets_stack_id_idx
    ON channels_dashboard_widgets (channel_id, stack_id)
    WHERE stack_id IS NOT NULL;

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'fix_dashboard_widgets_channel_fk is not reversible';
-- +goose StatementEnd
