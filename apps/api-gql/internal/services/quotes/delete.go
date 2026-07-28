package quotes

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/libs/audit"
	"github.com/twirapp/twir/libs/errors"
)

func (c *Service) Delete(ctx context.Context, channelID, actorID string, id uuid.UUID) error {
	quote, err := c.quotesRepository.GetByID(ctx, id)
	if err != nil {
		return err
	}

	if quote.ChannelID.String() != channelID {
		return errors.NewNotFoundError("Quote with this ID was not found for your channel")
	}

	if err := c.quotesRepository.Delete(ctx, id); err != nil {
		return err
	}

	if err := c.quotesCacher.Invalidate(ctx, channelID); err != nil {
		c.logger.Error("failed to invalidate quotes cache", slog.Any("error", err))
	}

	objectID := quote.ID.String()

	_ = c.auditRecorder.RecordDeleteOperation(
		ctx,
		audit.DeleteOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelQuote),
				ActorID:   &actorID,
				ChannelID: &channelID,
				ObjectID:  &objectID,
			},
			OldValue: quote,
		},
	)

	return nil
}
