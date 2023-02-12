package processor

import (
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func (c *Processor) SwitchEmoteOnly(operation model.EventOperationType) {
	resp, err := c.streamerApiClient.UpdateChatSettings(&helix.UpdateChatSettingsParams{
		BroadcasterID: c.channelId,
		ModeratorID:   c.channelId,
		EmoteMode: lo.ToPtr(lo.
			If(operation == model.OperationEnableEmoteOnly, true).
			Else(false)),
	})
	if err != nil {
		c.services.Logger.Sugar().Error(err)
	}

	if resp.ErrorMessage != "" {
		c.services.Logger.Sugar().Error("Cannot update channel chat settings", c.channelId, resp.ErrorMessage)
	}
}
