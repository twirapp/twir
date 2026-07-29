package pgx

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/twirapp/twir/libs/entities/platform"
)

const assignYouTubeBotQuery = `
	WITH updated AS (
		UPDATE channel_platforms
		SET
			bot_user_id = $1,
			updated_at = NOW()
		WHERE platform = $2
			AND bot_user_id IS DISTINCT FROM $1
		RETURNING channel_id
	)
	SELECT channel_id
	FROM updated
	ORDER BY channel_id`

func (r *Pgx) AssignYouTubeBot(ctx context.Context, botUserID uuid.UUID) ([]uuid.UUID, error) {
	rows, err := r.connection(ctx).Query(ctx, assignYouTubeBotQuery, botUserID, platform.PlatformYouTube)
	if err != nil {
		return nil, fmt.Errorf("assign YouTube bot to channel bindings: %w", err)
	}

	channelIDs, err := pgx.CollectRows(rows, pgx.RowTo[uuid.UUID])
	if err != nil {
		return nil, fmt.Errorf("collect channel IDs assigned to YouTube bot: %w", err)
	}

	return channelIDs, nil
}
