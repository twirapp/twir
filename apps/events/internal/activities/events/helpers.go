package events

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	channels "github.com/twirapp/twir/libs/repositories/channels"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
	"github.com/twirapp/twir/libs/twitch"
	"go.temporal.io/sdk/activity"
)

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
	return twitch.NewUserClientWithContext(ctx, twitchUserID, c.cfg, c.bus)
}

func (c *Activity) getHelixBotApiClient(ctx context.Context, botID string) (
	*helix.Client,
	error,
) {
	return twitch.NewBotClientWithContext(ctx, botID, c.cfg, c.bus)
}

// should be used with broadcaster channel client, otherwise it will return error
func (c *Activity) getChannelMods(client *helix.Client, twitchPlatformID string) (
	[]helix.Moderator,
	error,
) {
	var cursor string
	var moderators []helix.Moderator

	for {
		modsReq, err := client.GetModerators(
			&helix.GetModeratorsParams{
				BroadcasterID: twitchPlatformID,
				After:         cursor,
			},
		)
		if err != nil {
			return nil, err
		}
		if modsReq.ErrorMessage != "" {
			return nil, fmt.Errorf(modsReq.ErrorMessage)
		}

		moderators = append(moderators, modsReq.Data.Moderators...)

		if modsReq.Data.Pagination.Cursor == "" {
			break
		}

		cursor = modsReq.Data.Pagination.Cursor
	}

	return moderators, nil
}

func (c *Activity) getChannelVips(client *helix.Client, twitchPlatformID string) (
	[]helix.ChannelVips,
	error,
) {
	var cursor string
	var vips []helix.ChannelVips

	for {
		vipsReq, err := client.GetChannelVips(
			&helix.GetChannelVipsParams{
				BroadcasterID: twitchPlatformID,
				After:         cursor,
			},
		)
		if err != nil {
			return nil, err
		}
		if vipsReq.ErrorMessage != "" {
			return nil, fmt.Errorf(vipsReq.ErrorMessage)
		}

		vips = append(vips, vipsReq.Data.ChannelsVips...)
		cursor = vipsReq.Data.Pagination.Cursor

		if vipsReq.Data.Pagination.Cursor == "" {
			break
		}
	}

	return vips, nil
}

func (c *Activity) getChannelDbEntity(ctx context.Context, channelId string) (
	model.Channels,
	error,
) {
	channelInfo, err := c.getChannelRuntimeInfo(ctx, channelId)
	if err != nil {
		return model.Channels{}, err
	}

	return model.Channels{
		ID:    channelInfo.BroadcasterUserID,
		BotID: channelInfo.BotID,
	}, nil
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
	channel, err := c.channelsRepo.GetByID(ctx, channelUUID)
	if err != nil {
		if errors.Is(err, channels.ErrNotFound) {
			return channelRuntimeInfo{}, fmt.Errorf("channel not found")
		}

		return channelRuntimeInfo{}, err
	}

	var broadcasterUserID string
	if channel.TwitchUserID != nil {
		broadcasterUserID = *channel.TwitchPlatformID
	}

	var twitchPlatformID string
	if channel.TwitchPlatformID != nil {
		twitchPlatformID = *channel.TwitchPlatformID
	}

	return channelRuntimeInfo{
		ChannelID:         channel.ID.String(),
		BroadcasterUserID: broadcasterUserID,
		TwitchPlatformID:  twitchPlatformID,
		BotID:             channel.BotID,
	}, nil
}

func (c *Activity) getChannelRuntimeInfoByTwitchBroadcasterID(
	ctx context.Context,
	twitchBroadcasterID string,
) (channelRuntimeInfo, error) {
	user, err := c.usersRepo.GetByPlatformID(ctx, platform.PlatformTwitch, twitchBroadcasterID)
	if err != nil {
		if errors.Is(err, usersmodel.ErrNotFound) {
			return channelRuntimeInfo{}, fmt.Errorf("channel not found")
		}

		return channelRuntimeInfo{}, err
	}

	userUUID, err := uuid.Parse(user.ID)
	if err != nil {
		return channelRuntimeInfo{}, fmt.Errorf("parse user id: %w", err)
	}

	channel, err := c.channelsRepo.GetByTwitchUserID(ctx, userUUID)
	if err != nil {
		if errors.Is(err, channels.ErrNotFound) {
			return channelRuntimeInfo{}, fmt.Errorf("channel not found")
		}

		return channelRuntimeInfo{}, err
	}

	return channelRuntimeInfo{
		ChannelID:         channel.ID.String(),
		BroadcasterUserID: user.ID,
		TwitchPlatformID:  twitchBroadcasterID,
		BotID:             channel.BotID,
	}, nil
}

func (c *Activity) getHelixUserByLogin(client *helix.Client, userLogin string) (helix.User, error) {
	user, err := client.GetUsers(
		&helix.UsersParams{
			Logins: []string{userLogin},
		},
	)
	if err != nil {
		return helix.User{}, err
	}
	if user.ErrorMessage != "" {
		return helix.User{}, fmt.Errorf(user.ErrorMessage)
	}

	if len(user.Data.Users) == 0 {
		return helix.User{}, errors.New("user not found")
	}

	return user.Data.Users[0], nil
}

func (c *Activity) getHelixUserById(client *helix.Client, userId string) (helix.User, error) {
	user, err := client.GetUsers(
		&helix.UsersParams{
			IDs: []string{userId},
		},
	)
	if err != nil {
		return helix.User{}, err
	}
	if user.ErrorMessage != "" {
		return helix.User{}, fmt.Errorf(user.ErrorMessage)
	}

	if len(user.Data.Users) == 0 {
		return helix.User{}, errors.New("user not found")
	}

	return user.Data.Users[0], nil
}
