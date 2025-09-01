package chat_alerts

import (
	"context"
	"fmt"
	"strings"

	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/events"
	deprecatedgormmodel "github.com/twirapp/twir/libs/gomodels"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
)

func (c *ChatAlerts) follow(
	ctx context.Context,
	settings deprecatedgormmodel.ChatAlertsSettings,
	req events.FollowMessage,
) error {
	if !settings.Followers.Enabled {
		return nil
	}

	if len(settings.Followers.Messages) == 0 {
		return nil
	}

	sample := lo.Sample(settings.Followers.Messages)

	var stream *deprecatedgormmodel.ChannelsStreams
	if err := c.db.Where(`"userId" = ?`, req.BaseInfo.ChannelID).
		Find(&stream).Error; err != nil {
		return err
	}

	text := strings.ReplaceAll(sample.Text, "{user}", req.UserName)

	var followersCount int64
	if stream != nil && stream.ID != "" {
		t := model.ChannelEventListItemTypeFollow
		count, err := c.channelEventListsRepo.CountBy(
			ctx,
			channelseventslist.CountByInput{
				ChannelID:    &req.BaseInfo.ChannelID,
				CreatedAtGTE: &stream.StartedAt,
				Type:         &t,
			},
		)
		if err != nil {
			return err
		}

		followersCount = count
	}

	text = strings.ReplaceAll(text, "{streamFollowers}", fmt.Sprint(followersCount))

	if text == "" {
		return nil
	}

	return c.bus.Bots.SendMessage.Publish(
		ctx,
		bots.SendMessageRequest{
			ChannelId:      req.BaseInfo.ChannelID,
			Message:        text,
			SkipRateLimits: true,
		},
	)
}
