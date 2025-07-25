package variables

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/logger/audit"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
)

func (c *Service) Delete(ctx context.Context, id uuid.UUID, channelID, actorID string) error {
	variable, err := c.variablesRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if variable.ChannelID != channelID {
		return ErrNotFound
	}

	err = c.variablesRepository.Delete(ctx, id)
	if err != nil {
		return err
	}

	c.logger.Audit(
		"Variable delete",
		audit.Fields{
			OldValue:      variable,
			ActorID:       &actorID,
			ChannelID:     &channelID,
			System:        mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelVariable),
			OperationType: audit.OperationDelete,
			ObjectID:      lo.ToPtr(variable.ID.String()),
		},
	)

	return nil
}
