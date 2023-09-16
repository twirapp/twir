package processor

import (
	"errors"
	"fmt"

	"github.com/nicklaw5/helix/v2"
)

func (c *Processor) ChangeTitle(newTitle string) error {
	hydratedTitle, err := c.HydrateStringWithData(newTitle, c.data)

	if err != nil || len(hydratedTitle) == 0 {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	req, err := c.streamerApiClient.EditChannelInformation(
		&helix.EditChannelInformationParams{
			BroadcasterID: c.channelId,
			Title:         hydratedTitle,
		},
	)
	if err != nil {
		return err
	}

	if req.ErrorMessage != "" {
		return errors.New(req.ErrorMessage)
	}

	return nil
}
