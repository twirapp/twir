package variables

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/libs/audit"
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

	_ = c.auditRecorder.RecordDeleteOperation(
		ctx,
		audit.DeleteOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelVariable),
				ActorID:   &actorID,
				ChannelID: &channelID,
				ObjectID:  lo.ToPtr(variable.ID.String()),
			},
			OldValue: variable,
		},
	)

	return nil
}
