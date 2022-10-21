package commands

import (
	"errors"
	"fmt"
	model "tsuwari/models"

	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func HandleGet(channelId string, services types.Services) []model.ChannelsCommands {
	cmds := getChannelCommands(services.DB, channelId)

	return cmds
}

func HandlePost(
	channelId string,
	services types.Services,
	dto *commandDto,
) (*model.ChannelsCommands, error) {
	isExists := isCommandWithThatNameExists(services.DB, channelId, dto.Name, dto.Aliases, nil)
	if isExists {
		return nil, errors.New("command with that name already exists")
	}

	newCommand := createCommandFromDto(dto, channelId)

	err := services.DB.Save(newCommand).Error
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("cannot create command")
	}

	responses := createResponsesFromDto(dto.Responses, newCommand.ID)
	err = services.DB.Save(&responses).Error
	if err != nil {
		services.DB.Where(`"id" = ?`, newCommand.ID).Delete(&model.ChannelsCommands{})

		return nil, errors.New("something went wrong on creating response")
	}

	newCommand.Responses = responses

	return newCommand, nil
}

func HandleDelete(channelId string, commandId string, services types.Services) error {
	command, err := getChannelCommand(services.DB, channelId, commandId)
	if err != nil || command == nil {
		return errors.New("command not found")
	}

	err = services.DB.Delete(&command).Error
	if err != nil {
		return errors.New("cannot delete command")
	}

	return nil
}

func HandleUpdate(
	channelId string,
	commandId string,
	dto *commandDto,
	services types.Services,
) (*model.ChannelsCommands, error) {
	command, err := getChannelCommand(services.DB, channelId, commandId)
	if err != nil || command == nil {
		return nil, errors.New("command not found")
	}

	isExists := isCommandWithThatNameExists(
		services.DB,
		channelId,
		dto.Name,
		dto.Aliases,
		&command.ID,
	)
	if isExists {
		return nil, errors.New("command with that name already exists")
	}

	err = services.DB.Model(command).Updates(createCommandFromDto(dto, channelId)).Error
	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	services.DB.Where(`"commandId" = ?`, command.ID).Delete(&model.ChannelsCommandsResponses{})
	responses := createResponsesFromDto(dto.Responses, commandId)
	err = services.DB.Save(&responses).Error
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("something went wrong on creating response")
	}

	command.Responses = responses

	return command, nil
}
