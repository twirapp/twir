package processor

import "github.com/satont/go-helix/v2"

func (c *Processor) ChangeCategory(newCategory string) {
	searchCategory, err := c.streamerApiClient.SearchCategories(&helix.SearchCategoriesParams{
		Query: newCategory,
	})

	if len(searchCategory.Data.Categories) == 0 || err != nil {
		return
	}

	c.streamerApiClient.EditChannelInformation(&helix.EditChannelInformationParams{
		BroadcasterID: c.channelId,
		GameID:        searchCategory.Data.Categories[0].ID,
	})
}