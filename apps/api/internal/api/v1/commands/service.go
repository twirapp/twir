package commands

import (
	"net/http"
	"strings"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"github.com/satont/tsuwari/libs/twitch"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api/internal/types"
	uuid "github.com/satori/go.uuid"
)

func handleGet(channelId string, services types.Services) ([]model.ChannelsCommands, error) {
	config := do.MustInvoke[cfg.Config](di.Provider)
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	twitchClient, err := twitch.NewAppClient(config, tokensGrpc)
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	cmds := getChannelCommands(services.DB, channelId)

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
			logger.Error(err)
		}

		if twitchUsersReq.ErrorMessage != "" {
			logger.Error(twitchUsersReq.ErrorMessage)
		}
	}

	return cmds, nil
}

func handlePost(
	channelId string,
	services types.Services,
	dto *commandDto,
) (*model.ChannelsCommands, error) {
	config := do.MustInvoke[cfg.Config](di.Provider)
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	twitchClient, err := twitch.NewAppClient(config, tokensGrpc)
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	dto.Name = strings.TrimSpace(dto.Name)
	dto.Name = strings.ToLower(dto.Name)
	dto.Name = strings.Replace(dto.Name, "!", "", 1)
	if len(dto.Name) == 0 {
		return nil, fiber.NewError(400, "name cannot be empty")
	}

	dto.Aliases = lo.Map(dto.Aliases, func(a string, _ int) string {
		a = strings.TrimSpace(a)
		a = strings.ToLower(a)
		a = strings.Replace(a, "!", "", 1)
		return a
	})
	if lo.SomeBy(dto.Aliases, func(a string) bool {
		return len(a) == 0
	}) {
		return nil, fiber.NewError(400, "aliase cannot be empty")
	}

	isExists := isCommandWithThatNameExists(services.DB, channelId, dto.Name, dto.Aliases, nil)
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
			logger.Error(err)
			return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
		}
		if twitchUsersReq.ErrorMessage != "" {
			logger.Error(twitchUsersReq.ErrorMessage)
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

	err = services.DB.Save(newCommand).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot create command")
	}

	responses := createResponsesFromDto(dto.Responses, newCommand.ID)
	err = services.DB.Save(&responses).Error
	if err != nil {
		services.DB.Where(`"id" = ?`, newCommand.ID).Delete(&model.ChannelsCommands{})

		return nil, fiber.NewError(
			http.StatusInternalServerError,
			"something went wrong on creating response",
		)
	}

	newCommand.Responses = responses

	return newCommand, nil
}

func handleDelete(channelId string, commandId string, services types.Services) error {
	command, err := getChannelCommand(services.DB, channelId, commandId)
	if err != nil || command == nil {
		return fiber.NewError(http.StatusNotFound, "command not found")
	}

	err = services.DB.Delete(&command).Error
	if err != nil {
		return fiber.NewError(http.StatusInternalServerError, "cannot delete command")
	}

	return nil
}

func handleUpdate(
	channelId string,
	commandId string,
	dto *commandDto,
	services types.Services,
) (*model.ChannelsCommands, error) {
	config := do.MustInvoke[cfg.Config](di.Provider)
	logger := do.MustInvoke[interfaces.Logger](di.Provider)
	tokensGrpc := do.MustInvoke[tokens.TokensClient](di.Provider)
	twitchClient, err := twitch.NewAppClient(config, tokensGrpc)
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
	}

	command, err := getChannelCommand(services.DB, channelId, commandId)
	if err != nil || command == nil {
		return nil, fiber.NewError(http.StatusNotFound, "command not found")
	}

	if len(dto.Responses) == 0 && !command.Default {
		return nil, fiber.NewError(400, "responses cannot be empty")
	}

	dto.Name = strings.TrimSpace(dto.Name)
	dto.Name = strings.ToLower(dto.Name)
	dto.Name = strings.Replace(dto.Name, "!", "", 1)
	if len(dto.Name) == 0 {
		return nil, fiber.NewError(400, "name cannot be empty")
	}

	dto.Aliases = lo.Map(dto.Aliases, func(a string, _ int) string {
		a = strings.TrimSpace(a)
		a = strings.ToLower(a)
		a = strings.Replace(a, "!", "", 1)
		return a
	})
	if lo.SomeBy(dto.Aliases, func(a string) bool {
		return len(a) == 0
	}) {
		return nil, fiber.NewError(400, "aliase cannot be empty")
	}

	isExists := isCommandWithThatNameExists(
		services.DB,
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
	command.OnlineOnly = *dto.OnlineOnly
	command.RequiredUsedChannelPoints = dto.RequiredUsedChannelPoints
	command.RequiredWatchTime = dto.RequiredWatchTime
	command.RequiredMessages = dto.RequiredMessages

	command.DeniedUsersIDS = []string{}
	command.AllowedUsersIDS = []string{}
	if len(dto.DeniedUsersIds) > 0 || len(dto.AllowedUsersIds) > 0 {
		twitchUsersReq, err := twitchClient.GetUsers(&helix.UsersParams{
			Logins: append(dto.DeniedUsersIds, dto.AllowedUsersIds...),
		})
		if err != nil {
			logger.Error(err)
			return nil, fiber.NewError(http.StatusInternalServerError, "internal error")
		}
		if twitchUsersReq.ErrorMessage != "" {
			logger.Error(twitchUsersReq.ErrorMessage)
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

	err = services.DB.
		Select("*").
		Updates(command).
		Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if !command.Default {
		services.DB.Where(`"commandId" = ?`, command.ID).Delete(&model.ChannelsCommandsResponses{})
		responses := createResponsesFromDto(dto.Responses, commandId)
		err = services.DB.Save(&responses).Error
		if err != nil {
			logger.Error(err)
			return nil, fiber.NewError(
				http.StatusInternalServerError,
				"something went wrong on creating response",
			)
		}

		command.Responses = responses
	}

	newCmd, err := getChannelCommand(services.DB, channelId, commandId)
	if err != nil || command == nil {
		return nil, fiber.NewError(http.StatusNotFound, "command not found")
	}

	return newCmd, nil
}

func handlePatch(
	channelId, commandId string,
	dto *commandPatchDto,
	services types.Services,
) (*model.ChannelsCommands, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Provider)

	command, err := getChannelCommand(services.DB, channelId, commandId)
	if err != nil || command == nil {
		return nil, fiber.NewError(http.StatusNotFound, "command not found")
	}

	command.Enabled = *dto.Enabled

	err = services.DB.
		Select("*").
		Updates(command).
		Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	newCommand, _ := getChannelCommand(services.DB, channelId, commandId)
	return newCommand, nil
}
