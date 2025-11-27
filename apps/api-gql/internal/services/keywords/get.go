package keywords

import (
	"context"

	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func (c *Service) GetAllByChannelID(ctx context.Context, channelID string) (
	[]entity.Keyword,
	error,
) {
	keywords, err := c.keywordsRepository.GetAllByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	converted := make([]entity.Keyword, len(keywords))
	for i, keyword := range keywords {
		converted[i] = c.dbToModel(keyword)
	}

	return converted, nil
}
