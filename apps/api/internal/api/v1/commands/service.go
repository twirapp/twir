package commands

import (
	"net/http"
	"strings"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api/internal/types"
	uuid "github.com/satori/go.uuid"
)

func handleGet(channelId string, services types.Services) []model.ChannelsCommands {
	cmds := getChannelCommands(services.DB, channelId)

	return cmds
}

func handlePost(
	channelId string,
	services types.Services,
	dto *commandDto,
) (*model.ChannelsCommands, error) {
	dto.Name = strings.ToLower(dto.Name)
	dto.Aliases = lo.Map(dto.Aliases, func(a string, _ int) string {
		return strings.ToLower(a)
	})

	isExists := isCommandWithThatNameExists(services.DB, channelId, dto.Name, dto.Aliases, nil)
	if isExists {
		return nil, fiber.NewError(400, "command with that name already exists")
	}

	if len(dto.Responses) == 0 {
		return nil, fiber.NewError(400, "responses cannot be empty")
	}

	newCommand := createCommandFromDto(dto, channelId, lo.ToPtr(uuid.NewV4().String()))

	err := services.DB.Save(newCommand).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot create command")
	}

	for _, r := range dto.Restrictions {
		services.DB.Save(&model.CommandRestriction{
			ID:        uuid.NewV4().String(),
			CommandID: newCommand.ID,
			Type:      r.Type,
			Value:     r.Value,
		})
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
	dto.Name = strings.ToLower(dto.Name)
	dto.Aliases = lo.Map(dto.Aliases, func(a string, _ int) string {
		return strings.ToLower(a)
	})

	command, err := getChannelCommand(services.DB, channelId, commandId)
	if err != nil || command == nil {
		return nil, fiber.NewError(http.StatusNotFound, "command not found")
	}

	if len(dto.Responses) == 0 && !command.Default {
		return nil, fiber.NewError(400, "responses cannot be empty")
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

	command.Aliases = dto.Aliases
	command.Cooldown = null.IntFrom(int64(dto.Cooldown))
	command.CooldownType = dto.CooldownType
	command.Description = null.StringFromPtr(dto.Description)
	command.Enabled = *dto.Enabled
	command.IsReply = *dto.IsReply
	command.KeepResponsesOrder = *dto.KeepResponsesOrder
	command.Name = dto.Name
	command.Permission = dto.Permission
	command.Visible = *dto.Visible

	err = services.DB.
		Select("*").
		Updates(command).
		Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, err
	}

	if !command.Default {
		services.DB.Where(`"commandId" = ?`, command.ID).Delete(&model.ChannelsCommandsResponses{})
		responses := createResponsesFromDto(dto.Responses, commandId)
		err = services.DB.Save(&responses).Error
		if err != nil {
			services.Logger.Sugar().Error(err)
			return nil, fiber.NewError(
				http.StatusInternalServerError,
				"something went wrong on creating response",
			)
		}

		command.Responses = responses
	}

	return command, nil
}
