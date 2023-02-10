package processor

import "github.com/satont/go-helix/v2"

func (c *Processor) ChangeTitle(newTitle string) {
	c.streamerApiClient.EditChannelInformation(&helix.EditChannelInformationParams{
		BroadcasterID: c.channelId,
		Title:         newTitle,
	})
}
