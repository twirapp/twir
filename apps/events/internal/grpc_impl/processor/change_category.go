package processor

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
)

func (c *Processor) ChangeCategory(newCategory string) error {
	searchCategory, err := c.streamerApiClient.SearchCategories(
		&helix.SearchCategoriesParams{
			Query: newCategory,
		},
	)
	if err != nil {
		return err
	}

	if len(searchCategory.Data.Categories) == 0 {
		return errors.New("cannot get category with that name")
	}

	editReq, err := c.streamerApiClient.EditChannelInformation(
		&helix.EditChannelInformationParams{
			BroadcasterID: c.channelId,
			GameID:        searchCategory.Data.Categories[0].ID,
		},
	)
	if err != nil {
		return err
	}

	if editReq.ErrorMessage != "" {
		return errors.New(editReq.ErrorMessage)
	}

	return nil
}
