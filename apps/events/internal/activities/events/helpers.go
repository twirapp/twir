package events

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/kvizyx/twitchy/helix"
	"github.com/twirapp/twir/apps/events/internal/shared"
	channelentity "github.com/twirapp/twir/libs/entities/channel"
	"github.com/twirapp/twir/libs/entities/platform"
	channels "github.com/twirapp/twir/libs/repositories/channels"
	"github.com/twirapp/twir/libs/twitch"
	"go.temporal.io/sdk/activity"
)

type twitchUserClientFactory func(context.Context, uuid.UUID) (*helix.Client, error)
type twitchBotClientFactory func(context.Context, string) (*helix.Client, error)
type twitchChannelBotClientFactory func(context.Context, string, string) (*helix.Client, error)

func (c *Activity) getWorkflowExecutionState(ctx context.Context) (
	shared.EventsWorkflowExecutionState,
	error,
) {
	info := activity.GetInfo(ctx)
	workflowID := info.WorkflowExecution.ID
	if workflowID == "" {
		return shared.EventsWorkflowExecutionState{}, errors.New("workflow id is empty")
	}

	state := shared.EventsWorkflowExecutionState{}

	cachedState, err := c.redis.Get(ctx, "events:workflows:"+workflowID).Bytes()
	if err != nil {
		return state, err
	}

	if err := json.Unmarshal(cachedState, &state); err != nil {
		return state, err
	}

	return state, nil
}

func (c *Activity) setWorkflowExecutionState(
	ctx context.Context,
	state shared.EventsWorkflowExecutionState,
) error {
	info := activity.GetInfo(ctx)
	workflowID := info.WorkflowExecution.ID
	if workflowID == "" {
		return errors.New("workflow id is empty")
	}

	bytes, err := json.Marshal(state)
	if err != nil {
		return err
	}

	_, err = c.redis.Set(ctx, "events:workflows:"+workflowID, bytes, 7*24*time.Hour).Result()
	if err != nil {
		return err
	}

	return nil
}

func (c *Activity) getHelixChannelApiClient(ctx context.Context, twitchUserID string) (
	*helix.Client,
	error,
) {
	parsedUserID, err := uuid.Parse(twitchUserID)
	if err != nil {
		user, userErr := c.usersRepo.GetByPlatformID(ctx, platform.PlatformTwitch, twitchUserID)
		if userErr != nil {
			return nil, fmt.Errorf("resolve twitch user id: %w", userErr)
		}
		parsedUserID = user.ID
	}

	if c.newTwitchUserClient != nil {
		return c.newTwitchUserClient(ctx, parsedUserID)
	}

	return twitch.NewUserClientWithContext(ctx, parsedUserID, c.cfg, c.bus)
}

func (c *Activity) getHelixBotApiClient(ctx context.Context, botID string) (
	*helix.Client,
	error,
) {
	if c.newTwitchBotClient != nil {
		return c.newTwitchBotClient(ctx, botID)
	}

	return twitch.NewBotClientWithContext(ctx, botID, c.cfg, c.bus)
}

func (c *Activity) getHelixChannelBotApiClient(ctx context.Context, botID string, channelID string) (
	*helix.Client,
	error,
) {
	if c.newChannelBotClient != nil {
		return c.newChannelBotClient(ctx, botID, channelID)
	}

	return twitch.NewChannelBotClientWithContext(ctx, botID, channelID, c.cfg)
}

// should be used with broadcaster channel client, otherwise it will return error
func (c *Activity) getChannelMods(ctx context.Context, client *helix.Client, twitchPlatformID string) (
	[]helix.Moderator,
	error,
) {
	pager, err := client.Moderation.GetModeratorsPager(
		helix.GetModeratorsRequest{BroadcasterID: twitchPlatformID},
	)
	if err != nil {
		return nil, err
	}

	var moderators []helix.Moderator
	for pager.Next(ctx) {
		moderators = append(moderators, pager.Page().Data...)
	}

	return moderators, pager.Err()
}

func (c *Activity) getChannelVips(ctx context.Context, client *helix.Client, twitchPlatformID string) (
	[]helix.VIP,
	error,
) {
	pager, err := client.Moderation.GetVIPsPager(
		helix.GetVIPsRequest{BroadcasterID: twitchPlatformID},
	)
	if err != nil {
		return nil, err
	}

	var vips []helix.VIP
	for pager.Next(ctx) {
		vips = append(vips, pager.Page().Data...)
	}

	return vips, pager.Err()
}

func (c *Activity) getChannelDbEntity(ctx context.Context, channelId string) (
	channelRuntimeInfo,
	error,
) {
	return c.getChannelRuntimeInfo(ctx, channelId)
}

func (c *Activity) getTwitchChannelDbEntity(ctx context.Context, data shared.EventData) (
	channelRuntimeInfo,
	error,
) {
	broadcasterID := twitchBroadcasterID(data)
	if broadcasterID == "" {
		return channelRuntimeInfo{}, errors.New("twitch broadcaster id is empty")
	}

	return c.getChannelDbEntity(ctx, broadcasterID)
}

func (c *Activity) getEventTwitchBotApiClient(ctx context.Context, data shared.EventData) (
	*helix.Client,
	error,
) {
	channel, err := c.getTwitchChannelDbEntity(ctx, data)
	if err != nil {
		return nil, err
	}

	return c.getHelixBotApiClient(ctx, channel.BotID)
}

func twitchBroadcasterID(data shared.EventData) string {
	if data.ChannelTwitchPlatformID != "" {
		return data.ChannelTwitchPlatformID
	}
	if data.Platform == platform.PlatformTwitch {
		return data.ChannelID
	}

	return ""
}

func (c *Activity) getChannelRuntimeInfo(ctx context.Context, channelId string) (channelRuntimeInfo, error) {
	channelUUID, err := uuid.Parse(channelId)
	if err == nil {
		return c.getChannelRuntimeInfoByChannelUUID(ctx, channelUUID)
	}

	return c.getChannelRuntimeInfoByTwitchBroadcasterID(ctx, channelId)
}

func (c *Activity) getChannelRuntimeInfoByChannelUUID(
	ctx context.Context,
	channelUUID uuid.UUID,
) (channelRuntimeInfo, error) {
	channel, err := c.channelService.GetChannelByID(ctx, channelUUID)
	if err != nil {
		if errors.Is(err, channels.ErrNotFound) {
			return channelRuntimeInfo{}, fmt.Errorf("channel not found")
		}

		return channelRuntimeInfo{}, err
	}

	return getTwitchChannelRuntimeInfo(channel)
}

func getTwitchChannelRuntimeInfo(channel channelentity.Channel) (channelRuntimeInfo, error) {
	twitchBinding, botConfig, ok, err := channel.TwitchBinding()
	if err != nil {
		return channelRuntimeInfo{}, fmt.Errorf("parse Twitch bot config: %w", err)
	}
	if !ok {
		return channelRuntimeInfo{}, errors.New("twitch channel binding not found")
	}

	return channelRuntimeInfo{
		ChannelID:         channel.ID.String(),
		BroadcasterUserID: twitchBinding.PlatformChannelID,
		TwitchPlatformID:  twitchBinding.PlatformChannelID,
		BotID:             botConfig.BotID,
		IsBotMod:          botConfig.IsBotMod,
		IsTwitchBanned:    botConfig.IsTwitchBanned,
	}, nil
}

func (c *Activity) getChannelRuntimeInfoByTwitchBroadcasterID(
	ctx context.Context,
	twitchBroadcasterID string,
) (channelRuntimeInfo, error) {
	channel, err := c.channelService.GetChannelByPlatformChannelID(
		ctx,
		platform.PlatformTwitch,
		twitchBroadcasterID,
	)
	if err != nil {
		if errors.Is(err, channels.ErrNotFound) {
			return channelRuntimeInfo{}, fmt.Errorf("channel not found")
		}

		return channelRuntimeInfo{}, err
	}

	return getTwitchChannelRuntimeInfo(channel)
}

func (c *Activity) getHelixUserByLogin(ctx context.Context, client *helix.Client, userLogin string) (helix.User, error) {
	user, err := client.Users.GetUsers(
		ctx, helix.GetUsersRequest{
			Logins: []string{userLogin},
		},
	)
	if err != nil {
		return helix.User{}, err
	}

	if len(user.Data) == 0 {
		return helix.User{}, errors.New("user not found")
	}

	return user.Data[0], nil
}

func (c *Activity) getHelixUserById(ctx context.Context, client *helix.Client, userId string) (helix.User, error) {
	user, err := client.Users.GetUsers(
		ctx, helix.GetUsersRequest{
			IDs: []string{userId},
		},
	)
	if err != nil {
		return helix.User{}, err
	}

	if len(user.Data) == 0 {
		return helix.User{}, errors.New("user not found")
	}

	return user.Data[0], nil
}
