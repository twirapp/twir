package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func QuotesFrom(quote entity.Quote) gqlmodel.Quote {
	return gqlmodel.Quote{
		ID:          quote.ID,
		Number:      quote.Number,
		Text:        quote.Text,
		CreatorID:   quote.CreatorID,
		CreatorName: quote.CreatorName,
		GameID:      quote.GameID,
		GameName:    quote.GameName,
		CreatedAt:   quote.CreatedAt,
		UpdatedAt:   quote.UpdatedAt,
	}
}
