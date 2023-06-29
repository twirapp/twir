package processor

import (
	"errors"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
)

func (c *Processor) SwitchSubMode(operation model.EventOperationType) error {
	resp, err := c.streamerApiClient.UpdateChatSettings(
		&helix.UpdateChatSettingsParams{
			BroadcasterID: c.channelId,
			ModeratorID:   c.channelId,
			SubscriberMode: lo.ToPtr(
				lo.
					If(operation == model.OperationEnableSubMode, true).
					Else(false),
			),
		},
	)
	if err != nil {
		return err
	}

	if resp.ErrorMessage != "" {
		return errors.New(resp.ErrorMessage)
	}

	return nil
}
