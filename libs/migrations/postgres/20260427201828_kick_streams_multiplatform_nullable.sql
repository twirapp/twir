-- +goose Up
-- +goose StatementBegin
ALTER TABLE channels_streams
    ALTER COLUMN "gameId" DROP NOT NULL,
    ALTER COLUMN "gameName" DROP NOT NULL,
    ALTER COLUMN "communityIds" DROP NOT NULL,
    ALTER COLUMN "type" DROP NOT NULL,
    ALTER COLUMN "title" DROP NOT NULL,
    ALTER COLUMN "viewerCount" DROP NOT NULL,
    ALTER COLUMN "startedAt" DROP NOT NULL,
    ALTER COLUMN "language" DROP NOT NULL,
    ALTER COLUMN "thumbnailUrl" DROP NOT NULL,
    ALTER COLUMN "tagIds" DROP NOT NULL,
    ALTER COLUMN "tags" DROP NOT NULL,
    ALTER COLUMN "isMature" DROP NOT NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE channels_streams
    ALTER COLUMN "gameId" SET NOT NULL,
    ALTER COLUMN "gameName" SET NOT NULL,
    ALTER COLUMN "communityIds" SET NOT NULL,
    ALTER COLUMN "type" SET NOT NULL,
    ALTER COLUMN "title" SET NOT NULL,
    ALTER COLUMN "viewerCount" SET NOT NULL,
    ALTER COLUMN "startedAt" SET NOT NULL,
    ALTER COLUMN "language" SET NOT NULL,
    ALTER COLUMN "thumbnailUrl" SET NOT NULL,
    ALTER COLUMN "tagIds" SET NOT NULL,
    ALTER COLUMN "tags" SET NOT NULL,
    ALTER COLUMN "isMature" SET NOT NULL;
-- +goose StatementEnd
