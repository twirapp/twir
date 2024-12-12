package variables

import (
	"context"

	dbmodels "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
)

func (c *Service) Delete(ctx context.Context, id, channelID, actorID string) error {
	entity := dbmodels.ChannelsCustomvars{}
	if err := c.gorm.
		WithContext(ctx).
		Where(`"channelId" = ? AND id = ?`, channelID, id).
		First(&entity).Error; err != nil {
		return err
	}

	if err := c.gorm.
		WithContext(ctx).
		Delete(&entity).Error; err != nil {
		return err
	}

	c.logger.Audit(
		"Variable delete",
		audit.Fields{
			OldValue:      entity,
			ActorID:       &actorID,
			ChannelID:     &channelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelVariable),
			OperationType: audit.OperationDelete,
			ObjectID:      &entity.ID,
		},
	)

	return nil
}
