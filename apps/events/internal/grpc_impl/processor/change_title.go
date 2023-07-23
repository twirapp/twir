package processor

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
)

func (c *Processor) ChangeTitle(newTitle string) error {
	req, err := c.streamerApiClient.EditChannelInformation(
		&helix.EditChannelInformationParams{
			BroadcasterID: c.channelId,
			Title:         newTitle,
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
