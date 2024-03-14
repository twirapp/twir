package events

import (
	"context"
	"fmt"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/events/internal/shared"
	model "github.com/satont/twir/libs/gomodels"
	"go.temporal.io/sdk/activity"
)

func (c *Activity) CommandAllowOrRemoveUserPermission(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	hydratedName, hydrateErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		operation.Input.String,
		data,
	)
	if hydrateErr != nil || len(hydratedName) == 0 {
		return fmt.Errorf("cannot hydrate string %w", hydrateErr)
	}

	twitchClient, twitchClientErr := c.getHelixBotApiClient(ctx, data.ChannelID)
	if twitchClientErr != nil {
		return fmt.Errorf("cannot get twitch client %w", twitchClientErr)
	}

	user, userErr := c.getHelixUserByLogin(twitchClient, hydratedName)
	if userErr != nil {
		return fmt.Errorf("cannot get user %w", userErr)
	}

	command := &model.ChannelsCommands{}
	commandErr := c.db.Where("id = ?", operation.Target.String).First(command).Error
	if commandErr != nil {
		return fmt.Errorf("command not found")
	}

	if operation.Type == model.OperationAllowCommandToUser {
		for _, allowedUserId := range command.AllowedUsersIDS {
			if allowedUserId == user.ID {
				return nil
			}
		}

		command.AllowedUsersIDS = append(command.AllowedUsersIDS, user.ID)
	} else {
		command.AllowedUsersIDS = lo.Filter(
			command.AllowedUsersIDS,
			func(item string, _ int) bool {
				return item != user.ID
			},
		)
	}

	return c.db.Save(command).Error
}

func (c *Activity) CommandDenyOrRemoveUserPermission(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	hydratedName, hydrateErr := c.hydrator.HydrateStringWithData(
		data.ChannelID,
		operation.Input.String,
		data,
	)
	if hydrateErr != nil || len(hydratedName) == 0 {
		return fmt.Errorf("cannot hydrate string %w", hydrateErr)
	}

	twitchClient, twitchClientErr := c.getHelixBotApiClient(ctx, data.ChannelID)
	if twitchClientErr != nil {
		return fmt.Errorf("cannot get twitch client %w", twitchClientErr)
	}

	user, userErr := c.getHelixUserByLogin(twitchClient, hydratedName)
	if userErr != nil {
		return userErr
	}

	command := &model.ChannelsCommands{}
	commandErr := c.db.Where("id = ?", operation.Target.String).First(command).Error
	if commandErr != nil {
		return fmt.Errorf("command not found: %w", commandErr)
	}

	if operation.Type == model.OperationDenyCommandToUser {
		for _, deniedUserId := range command.DeniedUsersIDS {
			if deniedUserId == user.ID {
				return nil
			}
		}

		command.DeniedUsersIDS = append(command.DeniedUsersIDS, user.ID)
	} else {
		command.DeniedUsersIDS = lo.Filter(
			command.DeniedUsersIDS,
			func(item string, _ int) bool {
				return item != user.ID
			},
		)
	}

	return c.db.Save(command).Error
}
