package commands

import (
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
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
	isExists := isCommandWithThatNameExists(services.DB, channelId, dto.Name, dto.Aliases, nil)
	if isExists {
		return nil, fiber.NewError(400, "command with that name already exists")
	}

	newCommand := createCommandFromDto(dto, channelId, lo.ToPtr(uuid.NewV4().String()))

	err := services.DB.Save(newCommand).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "cannot create command")
	}

	responses := createResponsesFromDto(dto.Responses, newCommand.ID)
	err = services.DB.Save(&responses).Error
	if err != nil {
		services.DB.Where(`"id" = ?`, newCommand.ID).Delete(&model.ChannelsCommands{})

		return nil, fiber.NewError(500, "something went wrong on creating response")
	}

	newCommand.Responses = responses

	return newCommand, nil
}

func handleDelete(channelId string, commandId string, services types.Services) error {
	command, err := getChannelCommand(services.DB, channelId, commandId)
	if err != nil || command == nil {
		return fiber.NewError(404, "command not found")
	}

	err = services.DB.Delete(&command).Error
	if err != nil {
		return fiber.NewError(500, "cannot delete command")
	}

	return nil
}

func handleUpdate(
	channelId string,
	commandId string,
	dto *commandDto,
	services types.Services,
) (*model.ChannelsCommands, error) {
	command, err := getChannelCommand(services.DB, channelId, commandId)
	if err != nil || command == nil {
		return nil, fiber.NewError(404, "command not found")
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

	err = services.DB.Model(command).
		Select("*").
		Updates(createCommandFromDto(dto, channelId, lo.ToPtr(commandId))).
		Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, err
	}

	services.DB.Where(`"commandId" = ?`, command.ID).Delete(&model.ChannelsCommandsResponses{})
	responses := createResponsesFromDto(dto.Responses, commandId)
	err = services.DB.Save(&responses).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "something went wrong on creating response")
	}

	command.Responses = responses

	return command, nil
}
