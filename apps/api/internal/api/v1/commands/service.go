package commands

import (
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"gorm.io/gorm"
	"net/http"
	"strings"

	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/samber/lo"
	uuid "github.com/satori/go.uuid"
)

func handleGet(channelId string) []model.ChannelsCommands {
	cmds := getChannelCommands(channelId)

	return cmds
}

func handlePost(
	channelId string,
	dto *commandDto,
) (*model.ChannelsCommands, error) {
	db := do.MustInvoke[*gorm.DB](di.Injector)
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	dto.Name = strings.ToLower(dto.Name)
	dto.Aliases = lo.Map(dto.Aliases, func(a string, _ int) string {
		return strings.ToLower(a)
	})

	isExists := isCommandWithThatNameExists(channelId, dto.Name, dto.Aliases, nil)
	if isExists {
		return nil, fiber.NewError(400, "command with that name already exists")
	}

	if len(dto.Responses) == 0 {
		return nil, fiber.NewError(400, "responses cannot be empty")
	}

	newCommand := createCommandFromDto(dto, channelId, lo.ToPtr(uuid.NewV4().String()))

	err := db.Save(newCommand).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot create command")
	}

	responses := createResponsesFromDto(dto.Responses, newCommand.ID)
	err = db.Save(&responses).Error
	if err != nil {
		db.Where(`"id" = ?`, newCommand.ID).Delete(&model.ChannelsCommands{})

		return nil, fiber.NewError(
			http.StatusInternalServerError,
			"something went wrong on creating response",
		)
	}

	newCommand.Responses = responses

	return newCommand, nil
}

func handleDelete(channelId string, commandId string) error {
	db := do.MustInvoke[*gorm.DB](di.Injector)
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	command, err := getChannelCommand(channelId, commandId)
	if err != nil || command == nil {
		return fiber.NewError(http.StatusNotFound, "command not found")
	}

	err = db.Delete(&command).Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot delete command")
	}

	return nil
}

func handleUpdate(
	channelId string,
	commandId string,
	dto *commandDto,
) (*model.ChannelsCommands, error) {
	db := do.MustInvoke[*gorm.DB](di.Injector)
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	dto.Name = strings.ToLower(dto.Name)
	dto.Aliases = lo.Map(dto.Aliases, func(a string, _ int) string {
		return strings.ToLower(a)
	})

	command, err := getChannelCommand(channelId, commandId)
	if err != nil || command == nil {
		return nil, fiber.NewError(http.StatusNotFound, "command not found")
	}

	if len(dto.Responses) == 0 && !command.Default {
		return nil, fiber.NewError(400, "responses cannot be empty")
	}

	isExists := isCommandWithThatNameExists(
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

	err = db.
		Select("*").
		Updates(command).
		Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if !command.Default {
		db.Where(`"commandId" = ?`, command.ID).Delete(&model.ChannelsCommandsResponses{})
		responses := createResponsesFromDto(dto.Responses, commandId)
		err = db.Save(&responses).Error
		if err != nil {
			logger.Error(err)
			return nil, fiber.NewError(
				http.StatusInternalServerError,
				"something went wrong on creating response",
			)
		}

		command.Responses = responses
	}

	return command, nil
}

func handlePatch(
	channelId, commandId string,
	dto *commandPatchDto,
) (*model.ChannelsCommands, error) {
	db := do.MustInvoke[*gorm.DB](di.Injector)
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	command, err := getChannelCommand(channelId, commandId)
	if err != nil || command == nil {
		return nil, fiber.NewError(http.StatusNotFound, "command not found")
	}

	command.Enabled = *dto.Enabled

	err = db.
		Select("*").
		Updates(command).
		Error
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	newCommand, _ := getChannelCommand(channelId, commandId)
	return newCommand, nil
}
