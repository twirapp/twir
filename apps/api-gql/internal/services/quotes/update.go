package quotes

import (
	"context"
	"log/slog"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/mappers"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/libs/audit"
	quotesrepository "github.com/twirapp/twir/libs/repositories/quotes"
)

type UpdateInput struct {
	ChannelID string
	ActorID   string
	ID        uuid.UUID
	Text      *string
}

func (c *Service) Update(ctx context.Context, input UpdateInput) (entity.Quote, error) {
	quote, err := c.quotesRepository.GetByID(ctx, input.ID)
	if err != nil {
		return entity.QuoteNil, err
	}

	if quote.ChannelID.String() != input.ChannelID {
		return entity.QuoteNil, ErrQuoteNotFound
	}

	newQuote, err := c.quotesRepository.Update(
		ctx,
		input.ID,
		quotesrepository.UpdateInput{
			Text: input.Text,
		},
	)
	if err != nil {
		return entity.QuoteNil, err
	}

	objectID := quote.ID.String()

	_ = c.auditRecorder.RecordUpdateOperation(
		ctx,
		audit.UpdateOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelQuote),
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  &objectID,
			},
			NewValue: newQuote,
			OldValue: quote,
		},
	)

	if err = c.quotesCacher.Invalidate(ctx, input.ChannelID); err != nil {
		c.logger.Error("failed to invalidate quotes cache", slog.Any("error", err))
	}

	return c.dbToModel(newQuote), nil
}
