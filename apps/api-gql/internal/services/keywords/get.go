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

	converted := make([]entity.Keyword, 0, len(keywords))
	for _, keyword := range keywords {
		converted = append(converted, c.dbToModel(keyword))
	}

	return converted, nil
}
