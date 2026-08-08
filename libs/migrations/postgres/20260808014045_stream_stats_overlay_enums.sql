-- +goose Up
-- +goose StatementBegin
CREATE TYPE channels_overlays_stream_stats_design AS ENUM (
    'GLASS',
    'CARDS',
    'NEON',
    'SOLID',
    'MINIMAL',
    'TERMINAL',
    'OUTLINE'
);

CREATE TYPE channels_overlays_stream_stats_variant AS ENUM (
    'HORIZONTAL',
    'HORIZONTAL_COMPACT',
    'VERTICAL',
    'VERTICAL_COMPACT',
    'LARGE'
);

CREATE TYPE channels_overlays_stream_stats_viewers_mode AS ENUM (
    'CUMULATIVE',
    'SEPARATE'
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TYPE IF EXISTS channels_overlays_stream_stats_viewers_mode;
DROP TYPE IF EXISTS channels_overlays_stream_stats_variant;
DROP TYPE IF EXISTS channels_overlays_stream_stats_design;
-- +goose StatementEnd
