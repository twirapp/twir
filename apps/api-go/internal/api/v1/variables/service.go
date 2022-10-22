package variables

import (
	model "tsuwari/models"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/satont/tsuwari/apps/api-go/internal/types"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func handleGet(channelId string, services types.Services) ([]model.ChannelsCustomvars, error) {
	variables := []model.ChannelsCustomvars{}
	err := services.DB.Where(`"channelId" = ?`, channelId).Find(&variables).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "cannot get variables")
	}

	return variables, nil
}

func handlePost(
	channelId string,
	dto *variableDto,
	services types.Services,
) (*model.ChannelsCustomvars, error) {
	existedVariable := &model.ChannelsCustomvars{}
	err := services.DB.Where(`"channelId" = ? AND name = ?`, channelId, dto.Name).
		First(existedVariable).
		Error
	if err == nil && existedVariable != nil {
		return nil, fiber.NewError(400, "variable with name name already exists")
	}

	newVariable := model.ChannelsCustomvars{
		ID:          uuid.NewV4().String(),
		Name:        dto.Name,
		Description: null.StringFromPtr(dto.Description),
		Type:        dto.Type,
		EvalValue:   null.StringFromPtr(dto.EvalValue),
		Response:    null.StringFromPtr(dto.Response),
		ChannelID:   channelId,
	}
	err = services.DB.Save(&newVariable).Error
	if err != nil {
		return nil, fiber.NewError(500, "cannot create variable")
	}

	return &newVariable, nil
}

func handleDelete(channelId string, variableId string, services types.Services) error {
	variable := &model.ChannelsCustomvars{}
	err := services.DB.Where(`"channelId" = ? AND "id" = ?`, channelId, variableId).
		First(variable).
		Error
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(404, "variable not found")
	}

	err = services.DB.Delete(variable).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return fiber.NewError(500, "cannot delete variable")
	}

	return nil
}

func handleUpdate(
	channelId string,
	variableId string,
	dto *variableDto,
	services types.Services,
) (*model.ChannelsCustomvars, error) {
	err := services.DB.Where("id = ?", variableId).First(&model.ChannelsCustomvars{}).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(404, "variable not found")
	}

	newData := model.ChannelsCustomvars{
		ID:          variableId,
		Name:        dto.Name,
		Description: null.StringFromPtr(dto.Description),
		Type:        dto.Type,
		EvalValue:   null.StringFromPtr(dto.EvalValue),
		Response:    null.StringFromPtr(dto.Response),
		ChannelID:   channelId,
	}

	err = services.DB.Select("*").Updates(&newData).Error
	if err != nil {
		services.Logger.Sugar().Error(err)
		return nil, fiber.NewError(500, "something happend on our side, cannot update variable")
	}

	return &newData, nil
}
