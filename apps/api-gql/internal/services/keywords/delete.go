package keywords

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/libs/audit"
)

func (c *Service) Delete(ctx context.Context, channelID, actorID string, id uuid.UUID) error {
	keyword, err := c.keywordsRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if keyword.ChannelID != channelID {
		return fmt.Errorf("keyword not found")
	}

	if err := c.keywordsRepository.Delete(ctx, id); err != nil {
		return err
	}

	if err := c.keywordsCacher.Invalidate(ctx, channelID); err != nil {
		c.logger.Error("failed to invalidate keywords cache", err)
	}

	_ = c.auditRecorder.RecordDeleteOperation(
		ctx,
		audit.DeleteOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelKeyword),
				ActorID:   &actorID,
				ChannelID: &channelID,
				ObjectID:  lo.ToPtr(keyword.ID.String()),
			},
			OldValue: keyword,
		},
	)

	return nil
}
