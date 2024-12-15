package keywords

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
)

func (c *Service) Delete(ctx context.Context, channelID, actorID, id string) error {
	parsedUuid, err := uuid.Parse(id)
	if err != nil {
		return err
	}

	keyword, err := c.keywordsRepository.GetByID(ctx, parsedUuid)
	if err != nil {
		return err
	}

	if keyword.ChannelID != channelID {
		return fmt.Errorf("keyword not found")
	}

	if err := c.keywordsRepository.Delete(ctx, parsedUuid); err != nil {
		return err
	}

	if err := c.keywordsCacher.Invalidate(ctx, channelID); err != nil {
		c.logger.Error("failed to invalidate keywords cache", err)
	}

	c.logger.Audit(
		"Keywords remove",
		audit.Fields{
			OldValue:      keyword,
			ActorID:       &actorID,
			ChannelID:     &channelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelKeyword),
			OperationType: audit.OperationDelete,
			ObjectID:      lo.ToPtr(keyword.ID.String()),
		},
	)

	return nil
}
