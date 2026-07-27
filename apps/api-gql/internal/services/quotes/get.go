package quotes

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func (c *Service) GetAllByChannelID(ctx context.Context, channelID string) ([]entity.Quote, error) {
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, err
	}

	quotes, err := c.quotesRepository.GetAllByChannelID(ctx, parsedChannelID)
	if err != nil {
		return nil, err
	}

	converted := make([]entity.Quote, len(quotes))
	for index, quote := range quotes {
		converted[index] = c.dbToModel(quote)
	}

	return converted, nil
}
