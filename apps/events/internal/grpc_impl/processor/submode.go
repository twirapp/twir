package processor

import (
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func (c *Processor) SwitchSubMode(operation model.EventOperationType) {
	c.streamerApiClient.UpdateChatSettings(&helix.UpdateChatSettingsParams{
		BroadcasterID: c.channelId,
		ModeratorID:   c.channelId,
		SubscriberMode: lo.ToPtr(lo.
			If(operation == model.OperationEnableSubMode, true).
			Else(false)),
	})
}
