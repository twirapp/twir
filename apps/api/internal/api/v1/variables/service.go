package variables

import (
	"context"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/api/internal/di"
	"github.com/satont/tsuwari/apps/api/internal/interfaces"
	"net/http"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/gofiber/fiber/v2"
	"github.com/guregu/null"
	"github.com/satont/tsuwari/apps/api/internal/types"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func handleGet(channelId string, services types.Services) ([]model.ChannelsCustomvars, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	variables := []model.ChannelsCustomvars{}
	err := services.DB.Where(`"channelId" = ?`, channelId).Find(&variables).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot get variables")
	}

	return variables, nil
}

func handleGetBuiltIn(services types.Services) ([]*parser.GetVariablesResponse_Variable, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)
	grpc := do.MustInvoke[parser.ParserClient](di.Injector)

	req, err := grpc.GetDefaultVariables(context.Background(), &emptypb.Empty{})
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot get builtin variables")
	}

	return req.List, nil
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
		EvalValue:   dto.EvalValue,
		Response:    dto.Response,
		ChannelID:   channelId,
	}
	err = services.DB.Save(&newVariable).Error
	if err != nil {
		return nil, fiber.NewError(http.StatusInternalServerError, "cannot create variable")
	}

	return &newVariable, nil
}

func handleDelete(channelId string, variableId string, services types.Services) error {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	variable := &model.ChannelsCustomvars{}
	err := services.DB.Where(`"channelId" = ? AND "id" = ?`, channelId, variableId).
		First(variable).
		Error
	if err == gorm.ErrRecordNotFound {
		return fiber.NewError(http.StatusNotFound, "variable not found")
	}

	err = services.DB.Delete(variable).Error
	if err != nil {
		logger.Error(err)
		return fiber.NewError(http.StatusInternalServerError, "cannot delete variable")
	}

	return nil
}

func handleUpdate(
	channelId string,
	variableId string,
	dto *variableDto,
	services types.Services,
) (*model.ChannelsCustomvars, error) {
	logger := do.MustInvoke[interfaces.Logger](di.Injector)

	err := services.DB.Where("id = ?", variableId).First(&model.ChannelsCustomvars{}).Error
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

	err = services.DB.Select("*").Updates(&newData).Error
	if err != nil {
		logger.Error(err)
		return nil, fiber.NewError(
			http.StatusInternalServerError,
			"something happend on our side, cannot update variable",
		)
	}

	return &newData, nil
}
