package commands

import (
	"errors"
	"fmt"
	model "tsuwari/models"

	"github.com/guregu/null"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
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

	newCommand := &model.ChannelsCommands{
		Name:               dto.Name,
		Cooldown:           null.IntFrom(int64(dto.Cooldown)),
		CooldownType:       dto.CooldownType,
		Enabled:            *dto.Enabled,
		Aliases:            dto.Aliases,
		Description:        null.StringFromPtr(dto.Description),
		Visible:            *dto.Visible,
		ChannelID:          channelId,
		Permission:         dto.Permission,
		Default:            false,
		DefaultName:        null.String{},
		Module:             "CUSTOM",
		IsReply:            *dto.IsReply,
		KeepResponsesOrder: *dto.KeepOrder,
	}

	err := services.DB.Save(newCommand).Error
	if err != nil {
		fmt.Println(err)
		return nil, errors.New("cannot create command")
	}

	responses := []model.ChannelsCommandsResponses{}
	for _, r := range dto.Responses {
		response := &model.ChannelsCommandsResponses{
			ID:        uuid.NewV4().String(),
			Text:      null.NewString(r.Text, true),
			Order:     int(r.Order),
			CommandID: newCommand.ID,
		}
		err := services.DB.Save(response).Error
		if err != nil {
			services.DB.Where(`"id" = ?`, newCommand.ID).Delete(&model.ChannelsCommands{})

			return nil, errors.New("something went wrong on creating response")
		}

		responses = append(responses, *response)
	}

	newCommand.Responses = responses

	return newCommand, nil
}

func HandleDelete(channelId string, commandId string, services types.Services) error {
	command := &model.ChannelsCommands{}
	err := services.DB.Where(`"channelId" = ? AND "id" = ?`, channelId, commandId).
		First(&command).
		Error
	if err != nil || command == nil {
		return errors.New("command not found")
	}

	err = services.DB.Delete(&command).Error
	if err != nil {
		return errors.New("cannot delete command")
	}

	return nil
}

func getChannelCommands(db *gorm.DB, channelId string) []model.ChannelsCommands {
	cmds := []model.ChannelsCommands{}
	db.Preload("Responses").
		Where(`"channelId" = ?`, channelId).
		Find(&cmds)

	return cmds
}

func isCommandWithThatNameExists(
	db *gorm.DB,
	channelId string,
	name string,
	aliases []string,
	exceptCommandId *string,
) bool {
	cmds := getChannelCommands(db, channelId)

	if len(cmds) == 0 {
		return false
	}

	strings := []string{}
	for _, v := range cmds {
		if exceptCommandId != nil && v.ID == *exceptCommandId {
			continue
		}
		strings = append(strings, v.Name)
		for _, a := range v.Aliases {
			strings = append(strings, a)
		}
	}

	for _, str := range strings {
		if str == name {
			return true
		}

		for _, a := range aliases {
			if a == str {
				return true
			}
		}
	}

	return false
}
