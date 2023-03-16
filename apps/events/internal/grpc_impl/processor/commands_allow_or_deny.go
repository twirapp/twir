package processor

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	model "github.com/satont/tsuwari/libs/gomodels"
)

func (c *Processor) AllowOrRemoveAllowCommandToUser(operation model.EventOperationType, commandId, input string) error {
	hydratedName, err := c.hydrateStringWithData(input, c.data)

	if err != nil || len(hydratedName) == 0 {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	user, err := c.streamerApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{hydratedName},
	})

	if err != nil || user.ErrorMessage != "" || len(user.Data.Users) == 0 {
		if err != nil {
			return err
		}

		return fmt.Errorf("user not found")
	}

	command := &model.ChannelsCommands{}
	err = c.services.DB.Where("id = ?", commandId).First(command).Error
	if err != nil {
		return fmt.Errorf("command not found")
	}

	if operation == model.OperationAllowCommandToUser {
		if lo.SomeBy(command.AllowedUsersIDS, func(item string) bool {
			return item == user.Data.Users[0].ID
		}) {
			return InternalError
		}

		command.AllowedUsersIDS = append(command.AllowedUsersIDS, user.Data.Users[0].ID)
	} else {
		command.AllowedUsersIDS = lo.Filter(command.AllowedUsersIDS, func(item string, _ int) bool {
			return item != user.Data.Users[0].ID
		})
	}

	err = c.services.DB.Save(command).Error

	return err
}

func (c *Processor) DenyOrRemoveDenyCommandToUser(operation model.EventOperationType, commandId, input string) error {
	hydratedName, err := c.hydrateStringWithData(input, c.data)

	if err != nil || len(hydratedName) == 0 {
		return fmt.Errorf("cannot hydrate string %w", err)
	}

	user, err := c.streamerApiClient.GetUsers(&helix.UsersParams{
		Logins: []string{hydratedName},
	})

	if err != nil || user.ErrorMessage != "" || len(user.Data.Users) == 0 {
		if err != nil {
			return err
		}

		return fmt.Errorf("user not found")
	}

	command := &model.ChannelsCommands{}
	err = c.services.DB.Where("id = ?", commandId).First(command).Error
	if err != nil {
		return fmt.Errorf("command not found")
	}

	if operation == model.OperationDenyCommandToUser {
		if lo.SomeBy(command.DeniedUsersIDS, func(item string) bool {
			return item == user.Data.Users[0].ID
		}) {
			return InternalError
		}

		command.DeniedUsersIDS = append(command.DeniedUsersIDS, user.Data.Users[0].ID)
	} else {
		command.DeniedUsersIDS = lo.Filter(command.DeniedUsersIDS, func(item string, _ int) bool {
			return item != user.Data.Users[0].ID
		})
	}

	err = c.services.DB.Save(command).Error

	return err
}
