package keywords

import (
	"context"
	"github.com/google/uuid"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func (c *Service) GetAllByChannelID(ctx context.Context, channelID string) (
	[]entity.Keyword,
	error,
) {
	parsedChannelID, err := uuid.Parse(channelID)
	if err != nil {
		return nil, err
	}

	keywords, err := c.keywordsRepository.GetAllByChannelID(ctx, parsedChannelID)
	if err != nil {
		return nil, err
	}

	converted := make([]entity.Keyword, len(keywords))
	for i, keyword := range keywords {
		converted[i] = c.dbToModel(keyword)
	}

	return converted, nil
}
