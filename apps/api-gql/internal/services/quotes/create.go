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

type CreateInput struct {
	ChannelID   string
	ActorID     string
	CreatorName *string
	Text        string
}

func (c *Service) Create(ctx context.Context, input CreateInput) (entity.Quote, error) {
	parsedChannelID, err := uuid.Parse(input.ChannelID)
	if err != nil {
		return entity.QuoteNil, err
	}

	quote, err := c.quotesRepository.Create(
		ctx,
		quotesrepository.CreateInput{
			ChannelID:   parsedChannelID,
			Text:        input.Text,
			CreatorID:   &input.ActorID,
			CreatorName: input.CreatorName,
		},
	)
	if err != nil {
		return entity.QuoteNil, err
	}

	objectID := quote.ID.String()

	_ = c.auditRecorder.RecordCreateOperation(
		ctx,
		audit.CreateOperation{
			Metadata: audit.OperationMetadata{
				System:    mappers.AuditSystemToTableName(gqlmodel.AuditLogSystemChannelQuote),
				ActorID:   &input.ActorID,
				ChannelID: &input.ChannelID,
				ObjectID:  &objectID,
			},
			NewValue: quote,
		},
	)

	if err := c.quotesCacher.Invalidate(ctx, input.ChannelID); err != nil {
		c.logger.Error("failed to invalidate quotes cache", slog.Any("error", err))
	}

	return c.dbToModel(quote), nil
}
