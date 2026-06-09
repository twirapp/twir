package chat_alerts

import (
	"context"
	"fmt"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/entities/platform"
	deprecatedgormmodel "github.com/twirapp/twir/libs/gomodels"
	channelseventslist "github.com/twirapp/twir/libs/repositories/channels_events_list"
	"github.com/twirapp/twir/libs/repositories/channels_events_list/model"
	"github.com/twirapp/twir/libs/twitch"
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
	internalChannelID := req.BaseInfo.ChannelDBID
	if internalChannelID == "" {
		internalChannelID = req.BaseInfo.ChannelID
	}
	platformChannelID := req.BaseInfo.ChannelID
	eventPlatform := req.BaseInfo.Platform
	if eventPlatform == "" {
		eventPlatform = platform.PlatformTwitch
	}
	streamUserID := platformChannelID
	if eventPlatform == platform.PlatformKick {
		streamUserID = internalChannelID
	}

	var stream *deprecatedgormmodel.ChannelsStreams
	if err := c.db.Where(`"userId" = ?`, streamUserID).
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
				ChannelID:    &internalChannelID,
				Platform:     &eventPlatform,
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

	if eventPlatform == platform.PlatformTwitch {
		user, err := c.usersRepo.GetByPlatformID(ctx, platform.PlatformTwitch, platformChannelID)
		if err != nil {
			return fmt.Errorf("cannot get user by platform id: %w", err)
		}

		twitchClient, err := twitch.NewUserClientWithContext(ctx, user.ID, c.cfg, c.bus)
		if err != nil {
			return err
		}

		followersReq, err := twitchClient.GetChannelFollows(
			&helix.GetChannelFollowsParams{
				BroadcasterID: platformChannelID,
			},
		)
		if err != nil {
			return err
		}

		text = strings.ReplaceAll(text, "{followers}", fmt.Sprint(followersReq.Data.Total))
	} else {
		text = strings.ReplaceAll(text, "{followers}", "0")
	}

	if text == "" {
		return nil
	}

	return c.bus.Bots.SendMessage.Publish(
		ctx,
		bots.SendMessageRequest{
			ChannelName:       lo.If(req.BaseInfo.ChannelName != "", &req.BaseInfo.ChannelName).Else(nil),
			ChannelId:         internalChannelID,
			PlatformChannelID: platformChannelID,
			Platform:          eventPlatform.String(),
			Message:           text,
			SkipRateLimits:    true,
		},
	)
}
