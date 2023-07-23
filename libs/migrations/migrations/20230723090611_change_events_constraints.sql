-- +goose Up
-- +goose StatementBegin
SELECT 'up SQL query';

ALTER TABLE channels_events_operations
	DROP CONSTRAINT "FK_b2e27e84fa5bfbf8fd27ac9e948";
ALTER TABLE channels_events_operations_filters
	DROP CONSTRAINT "FK_d56a4ca65d44fc4adbaf66a2e80";

ALTER TABLE channels_events_operations
	ADD CONSTRAINT fk_channels_events_operations_event_id
		FOREIGN KEY ("eventId")
			REFERENCES channels_events ("id")
			ON DELETE CASCADE;

ALTER TABLE channels_events_operations_filters
	ADD CONSTRAINT fk_channels_events_operations_filters_operation_id
		FOREIGN KEY ("operationId")
			REFERENCES channels_events_operations ("id")
			ON DELETE CASCADE;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
