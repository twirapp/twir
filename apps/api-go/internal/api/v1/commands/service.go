package commands

import (
	"errors"
	model "tsuwari/models"

	"github.com/guregu/null"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
)

func HandleGet(channelId string, services types.Services) []model.ChannelsCommands {
	cmds := []model.ChannelsCommands{}
	services.DB.Preload("Responses").Where(`"channelId" = ?`, channelId).Find(&cmds)

	return cmds
}

func HandlePost(
	channelId string,
	services types.Services,
	dto *commandDto,
) (*model.ChannelsCommands, error) {
	newCommand := &model.ChannelsCommands{
		Name:         dto.Name,
		Cooldown:     null.IntFromPtr(dto.Cooldown),
		CooldownType: dto.CooldownType,
		Enabled:      *dto.Enabled,
		Aliases:      dto.Aliases,
		Description:  null.StringFromPtr(dto.Description),
		Visible:      *dto.Visible,
		ChannelID:    channelId,
		Permission:   dto.Permission,
		Default:      false,
		DefaultName:  null.String{},
		Module:       "CUSTOM",
		IsReply:      *dto.IsReply,
		KeepOrder:    *dto.KeepOrder,
	}

	err := services.DB.Save(newCommand).Error
	if err != nil {
		return nil, errors.New("cannot create command")
	}

	return newCommand, nil
}
