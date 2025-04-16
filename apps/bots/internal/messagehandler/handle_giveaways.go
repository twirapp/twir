package messagehandler

import (
	"context"
	"strings"

	giveawaysbus "github.com/twirapp/twir/libs/bus-core/giveaways"
	"github.com/twirapp/twir/libs/repositories/giveaways/model"
)

func (c *MessageHandler) handleGiveaways(ctx context.Context, msg handleMessage) error {
	// TODO: uncomment in prod
	// if msg.DbStream == nil {
	// return nil
	// }

	giveaways, err := c.giveawaysCacher.Get(ctx, msg.BroadcasterUserId)
	if err != nil {
		return err
	}

	if len(giveaways) == 0 {
		return nil
	}

	for _, giveaway := range giveaways {
		if giveaway == model.ChannelGiveawayNil {
			continue
		}

		if !giveaway.IsRunning {
			continue
		}

		if giveaway.IsStopped {
			continue
		}

		if giveaway.IsFinished {
			_ = c.giveawaysCacher.Invalidate(ctx, msg.BroadcasterUserId)
			continue
		}

		if giveaway.IsArchived {
			_ = c.giveawaysCacher.Invalidate(ctx, msg.BroadcasterUserId)
			continue
		}

		if !strings.Contains(strings.ToLower(msg.Message.Text), strings.ToLower(giveaway.Keyword)) {
			return nil
		}

		err = c.bus.Giveaways.TryAddParticipant.Publish(giveawaysbus.TryAddParticipantRequest{
			UserID:     msg.ChatterUserId,
			GiveawayID: giveaway.ID,
		})
		if err != nil {
			return err
		}
	}

	return nil
}
