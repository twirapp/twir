package commands

import (
	model "tsuwari/models"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func getChannelCommands(db *gorm.DB, channelId string) []model.ChannelsCommands {
	cmds := []model.ChannelsCommands{}
	db.Preload("Responses").
		Where(`"channelId" = ?`, channelId).
		Find(&cmds)

	return cmds
}

func getChannelCommand(
	db *gorm.DB,
	channelId string,
	commandId string,
) (*model.ChannelsCommands, error) {
	command := &model.ChannelsCommands{}
	err := db.Where(`"channelId" = ? AND "id" = ?`, channelId, commandId).
		First(&command).
		Error
	if err != nil {
		return nil, err
	}
	return command, nil
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

func createCommandFromDto(dto *commandDto, channelId string) *model.ChannelsCommands {
	return &model.ChannelsCommands{
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
}

func createResponsesFromDto(
	responsesDto []responsesDto,
	commandId string,
) []model.ChannelsCommandsResponses {
	responses := []model.ChannelsCommandsResponses{}
	for _, r := range responsesDto {
		response := model.ChannelsCommandsResponses{
			ID:        uuid.NewV4().String(),
			Text:      null.NewString(r.Text, true),
			Order:     int(r.Order),
			CommandID: commandId,
		}
		responses = append(responses, response)
	}

	return responses
}
