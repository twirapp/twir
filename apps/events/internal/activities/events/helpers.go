package events

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/nicklaw5/helix/v2"
	"github.com/twirapp/twir/apps/events/internal/shared"
	model "github.com/twirapp/twir/libs/gomodels"
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

func (c *Activity) getHelixChannelApiClient(ctx context.Context, channelId string) (
	*helix.Client,
	error,
) {
	return twitch.NewUserClientWithContext(ctx, channelId, c.cfg, c.bus)
}

func (c *Activity) getHelixBotApiClient(ctx context.Context, botID string) (
	*helix.Client,
	error,
) {
	return twitch.NewBotClientWithContext(ctx, botID, c.cfg, c.bus)
}

// should be used with broadcaster channel client, otherwise it will return error
func (c *Activity) getChannelMods(client *helix.Client, channelId string) (
	[]helix.Moderator,
	error,
) {
	var cursor string
	var moderators []helix.Moderator

	for {
		modsReq, err := client.GetModerators(
			&helix.GetModeratorsParams{
				BroadcasterID: channelId,
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

func (c *Activity) getChannelVips(client *helix.Client, channelId string) (
	[]helix.ChannelVips,
	error,
) {
	var cursor string
	var vips []helix.ChannelVips

	for {
		vipsReq, err := client.GetChannelVips(
			&helix.GetChannelVipsParams{
				BroadcasterID: channelId,
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
	channel := model.Channels{}
	err := c.db.WithContext(ctx).Where(`"id" = ?`, channelId).First(&channel).Error
	if err != nil {
		return channel, err
	}

	return channel, nil
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
