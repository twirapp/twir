package commands

import (
	"net/http"
	"strings"

	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/libs/twitch"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/samber/lo"
	uuid "github.com/satori/go.uuid"
)

func (c *Commands) getService(channelId string) ([]model.ChannelsCommands, error) {
	twitchClient, err := twitch.NewAppClient(*c.services.Config, c.services.Grpc.Tokens)
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	cmds := c.getChannelCommands(channelId)

	usersForReq := []string{}

	for _, cmd := range cmds {
		usersForReq = append(usersForReq, cmd.DeniedUsersIDS...)
		usersForReq = append(usersForReq, cmd.AllowedUsersIDS...)
	}

	if len(usersForReq) == 0 {
		return cmds, nil
	}

	twitchUsersReq, err := twitchClient.GetUsers(&helix.UsersParams{
		IDs: usersForReq,
	})

	if err == nil && twitchUsersReq.ErrorMessage == "" {
		for i, cmd := range cmds {
			for userIdx, deniedUser := range cmd.DeniedUsersIDS {
				twitchUser, ok := lo.Find(twitchUsersReq.Data.Users, func(u helix.User) bool {
					return u.ID == deniedUser
				})
				if !ok {
					continue
				}

				cmds[i].DeniedUsersIDS[userIdx] = twitchUser.Login
			}

			for userIdx, allowedUser := range cmd.AllowedUsersIDS {
				twitchUser, ok := lo.Find(twitchUsersReq.Data.Users, func(u helix.User) bool {
					return u.ID == allowedUser
				})
				if !ok {
					continue
				}

				cmds[i].AllowedUsersIDS[userIdx] = twitchUser.Login
			}
		}
	} else {
		if err != nil {
			c.services.Logger.Error(err)
		}

		if twitchUsersReq.ErrorMessage != "" {
			c.services.Logger.Error(twitchUsersReq.ErrorMessage)
		}
	}

	return cmds, nil
}

func (c *Commands) postService(
	channelId string,
	dto *commandDto,
) (*model.ChannelsCommands, error) {
	twitchClient, err := twitch.NewAppClient(*c.services.Config, c.services.Grpc.Tokens)
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	dto.Name = strings.Replace(strings.ToLower(dto.Name), "!", "", 1)
	dto.Aliases = lo.Map(dto.Aliases, func(a string, _ int) string {
		return strings.Replace(strings.ToLower(a), "!", "", 1)
	})

	isExists := c.isCommandWithThatNameExists(channelId, dto.Name, dto.Aliases, nil)
	if isExists {
		return nil, fiber.NewError(400, "command with that name already exists")
	}

	if len(dto.Responses) == 0 {
		return nil, fiber.NewError(400, "responses cannot be empty")
	}

	newCommand := createCommandFromDto(dto, channelId, lo.ToPtr(uuid.NewV4().String()))

	newCommand.DeniedUsersIDS = []string{}
	newCommand.AllowedUsersIDS = []string{}

	if len(dto.DeniedUsersIds) > 0 || len(dto.AllowedUsersIds) > 0 {
		twitchUsersReq, err := twitchClient.GetUsers(&helix.UsersParams{
			Logins: append(dto.DeniedUsersIds, dto.AllowedUsersIds...),
		})
		if err != nil {
			c.services.Logger.Error(err)
			return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
		}
		if twitchUsersReq.ErrorMessage != "" {
			c.services.Logger.Error(twitchUsersReq.ErrorMessage)
			return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
		}

		for _, deniedUser := range dto.DeniedUsersIds {
			twitchUser, ok := lo.Find(twitchUsersReq.Data.Users, func(u helix.User) bool {
				return u.Login == strings.ToLower(deniedUser)
			})

			if !ok {
				continue
			}
			newCommand.DeniedUsersIDS = append(newCommand.DeniedUsersIDS, twitchUser.ID)
		}

		for _, allowedUser := range dto.AllowedUsersIds {
			twitchUser, ok := lo.Find(twitchUsersReq.Data.Users, func(u helix.User) bool {
				return u.Login == strings.ToLower(allowedUser)
			})

			if !ok {
				continue
			}
			newCommand.AllowedUsersIDS = append(newCommand.AllowedUsersIDS, twitchUser.ID)
		}
	}

	err = c.services.Gorm.Save(newCommand).Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot create command")
	}

	responses := createResponsesFromDto(dto.Responses, newCommand.ID)
	err = c.services.Gorm.Save(&responses).Error
	if err != nil {
		c.services.Gorm.Where(`"id" = ?`, newCommand.ID).Delete(&model.ChannelsCommands{})

		return nil, fiber.NewError(
			http.StatusInternalServerError,
			"something went wrong on creating response",
		)
	}

	newCommand.Responses = responses

	return newCommand, nil
}

func (c *Commands) deleteService(channelId string, commandId string) error {
	command, err := c.getChannelCommand(channelId, commandId)
	if err != nil || command == nil {
		return fiber.NewError(http.StatusNotFound, "command not found")
	}

	err = c.services.Gorm.Delete(&command).Error
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "cannot delete command")
	}

	return nil
}

func (c *Commands) putService(
	channelId string,
	commandId string,
	dto *commandDto,
) (*model.ChannelsCommands, error) {
	twitchClient, err := twitch.NewAppClient(*c.services.Config, c.services.Grpc.Tokens)
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	command, err := c.getChannelCommand(channelId, commandId)
	if err != nil || command == nil {
		return nil, fiber.NewError(http.StatusNotFound, "command not found")
	}

	if len(dto.Responses) == 0 && !command.Default {
		return nil, fiber.NewError(400, "responses cannot be empty")
	}

	dto.Name = strings.Replace(strings.ToLower(dto.Name), "!", "", 1)
	dto.Aliases = lo.Map(dto.Aliases, func(a string, _ int) string {
		return strings.Replace(strings.ToLower(a), "!", "", 1)
	})

	isExists := c.isCommandWithThatNameExists(
		channelId,
		dto.Name,
		dto.Aliases,
		&command.ID,
	)
	if isExists {
		return nil, fiber.NewError(400, "command with that name already exists")
	}

	command.Name = dto.Name
	command.Aliases = dto.Aliases
	command.Cooldown = null.IntFrom(int64(dto.Cooldown))
	command.CooldownType = dto.CooldownType
	command.Description = null.StringFromPtr(dto.Description)
	command.Enabled = *dto.Enabled
	command.IsReply = *dto.IsReply
	command.KeepResponsesOrder = *dto.KeepResponsesOrder
	command.Visible = *dto.Visible
	command.GroupID = null.StringFromPtr(dto.GroupID)
	command.RolesIDS = dto.RolesIDS

	command.DeniedUsersIDS = []string{}
	command.AllowedUsersIDS = []string{}
	if len(dto.DeniedUsersIds) > 0 || len(dto.AllowedUsersIds) > 0 {
		twitchUsersReq, err := twitchClient.GetUsers(&helix.UsersParams{
			Logins: append(dto.DeniedUsersIds, dto.AllowedUsersIds...),
		})
		if err != nil {
			c.services.Logger.Error(err)
			return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
		}
		if twitchUsersReq.ErrorMessage != "" {
			c.services.Logger.Error(twitchUsersReq.ErrorMessage)
			return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
		}
		for _, deniedUser := range dto.DeniedUsersIds {
			twitchUser, ok := lo.Find(twitchUsersReq.Data.Users, func(u helix.User) bool {
				return u.Login == strings.ToLower(deniedUser)
			})

			if !ok {
				continue
			}
			command.DeniedUsersIDS = append(command.DeniedUsersIDS, twitchUser.ID)
		}

		for _, allowedUser := range dto.AllowedUsersIds {
			twitchUser, ok := lo.Find(twitchUsersReq.Data.Users, func(u helix.User) bool {
				return u.Login == strings.ToLower(allowedUser)
			})

			if !ok {
				continue
			}
			command.AllowedUsersIDS = append(command.AllowedUsersIDS, twitchUser.ID)
		}
	}

	if dto.GroupID == nil {
		command.Group = nil
	}

	err = c.services.Gorm.
		Select("*").
		Updates(command).
		Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, err
	}

	if !command.Default {
		c.services.Gorm.Where(`"commandId" = ?`, command.ID).Delete(&model.ChannelsCommandsResponses{})
		responses := createResponsesFromDto(dto.Responses, commandId)
		err = c.services.Gorm.Save(&responses).Error
		if err != nil {
			c.services.Logger.Error(err)
			return nil, fiber.NewError(
				http.StatusInternalServerError,
				"something went wrong on creating response",
			)
		}

		command.Responses = responses
	}

	newCmd, err := c.getChannelCommand(channelId, commandId)
	if err != nil || command == nil {
		return nil, fiber.NewError(http.StatusNotFound, "command not found")
	}

	return newCmd, nil
}

func (c *Commands) patchService(
	channelId, commandId string,
	dto *commandPatchDto,
) (*model.ChannelsCommands, error) {
	command, err := c.getChannelCommand(channelId, commandId)
	if err != nil || command == nil {
		return nil, fiber.NewError(http.StatusNotFound, "command not found")
	}

	command.Enabled = *dto.Enabled

	err = c.services.Gorm.
		Select("*").
		Updates(command).
		Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, err
	}

	newCommand, _ := c.getChannelCommand(channelId, commandId)
	return newCommand, nil
}
