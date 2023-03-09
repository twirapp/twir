package variables

import (
	"context"
	"net/http"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func (c *Variables) getService(channelId string) ([]model.ChannelsCustomvars, error) {
	variables := []model.ChannelsCustomvars{}
	err := c.services.Gorm.Where(`"channelId" = ?`, channelId).Find(&variables).Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot get variables")
	}

	return variables, nil
}

func (c *Variables) builtInService() ([]*parser.GetVariablesResponse_Variable, error) {
	req, err := c.services.Grpc.Parser.GetDefaultVariables(context.Background(), &emptypb.Empty{})
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot get builtin variables")
	}

	return req.List, nil
}

func (c *Variables) postService(
	channelId string,
	dto *variableDto,
) (*model.ChannelsCustomvars, error) {
	existedVariable := &model.ChannelsCustomvars{}
	err := c.services.Gorm.Where(`"channelId" = ? AND name = ?`, channelId, dto.Name).
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
		EvalValue:   dto.EvalValue,
		Response:    dto.Response,
		ChannelID:   channelId,
	}
	err = c.services.Gorm.Save(&newVariable).Error
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot create variable")
	}

	return &newVariable, nil
}

func (c *Variables) deleteService(channelId string, variableId string) error {
	variable := &model.ChannelsCustomvars{}
	err := c.services.Gorm.Where(`"channelId" = ? AND "id" = ?`, channelId, variableId).
		First(variable).
		Error
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(http.StatusNotFound, "variable not found")
	}

	err = c.services.Gorm.Delete(variable).Error
	if err != nil {
		c.services.Logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot delete variable")
	}

	return nil
}

func (c *Variables) putService(
	channelId string,
	variableId string,
	dto *variableDto,
) (*model.ChannelsCustomvars, error) {
	err := c.services.Gorm.Where("id = ?", variableId).First(&model.ChannelsCustomvars{}).Error
	if err == gorm.ErrRecordNotFound {
		return nil, fiber.NewError(http.StatusNotFound, "variable not found")
	}

	newData := model.ChannelsCustomvars{
		ID:          variableId,
		Name:        dto.Name,
		Description: null.StringFromPtr(dto.Description),
		Type:        dto.Type,
		EvalValue:   dto.EvalValue,
		Response:    dto.Response,
		ChannelID:   channelId,
	}

	err = c.services.Gorm.Select("*").Updates(&newData).Error
	if err != nil {
		c.services.Logger.Error(err)
		return nil, fiber.NewError(
			http.StatusInternalServerError,
			"something happend on our side, cannot update variable",
		)
	}

	return &newData, nil
}
