package variables

import (
	"context"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/variables"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Variables struct {
	*impl_deps.Deps
}

func (c *Variables) rpcTypeToDb(t variables.VariableType) model.CustomVarType {
	switch t {
	case variables.VariableType_NUMBER:
		return model.CustomVarNumber
	case variables.VariableType_TEXT:
		return model.CustomVarText
	case variables.VariableType_SCRIPT:
		return model.CustomVarScript
	default:
		return model.CustomVarText
	}
}

func (c *Variables) convertEntity(entity *model.ChannelsCustomvars) *variables.Variable {
	var t variables.VariableType
	switch entity.Type {
	case model.CustomVarNumber:
		t = variables.VariableType_NUMBER
	case model.CustomVarText:
		t = variables.VariableType_TEXT
	case model.CustomVarScript:
		t = variables.VariableType_SCRIPT
	}

	return &variables.Variable{
		Id:          entity.ID,
		Name:        entity.Name,
		Description: entity.Description.Ptr(),
		Type:        t,
		Response:    entity.Response,
		EvalValue:   entity.EvalValue,
		ChannelId:   entity.ChannelID,
	}
}

func (c *Variables) VariablesGetAll(ctx context.Context, req *emptypb.Empty) (*variables.GetAllResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	var entities []*model.ChannelsCustomvars
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ?`, dashboardId).
		Find(&entities).Error; err != nil {
		return nil, err
	}

	return &variables.GetAllResponse{
		Variables: lo.Map(
			entities, func(v *model.ChannelsCustomvars, _ int) *variables.Variable {
				return c.convertEntity(v)
			},
		),
	}, nil
}

func (c *Variables) VariablesGetById(ctx context.Context, req *variables.GetByIdRequest) (*variables.Variable, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.ChannelsCustomvars{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? AND "id" = ?`, dashboardId, req.Id).
		First(entity).Error; err != nil {
		return nil, err
	}

	return c.convertEntity(entity), nil
}

func (c *Variables) VariablesCreate(ctx context.Context, req *variables.CreateRequest) (*variables.Variable, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := &model.ChannelsCustomvars{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: null.StringFromPtr(req.Description),
		Type:        c.rpcTypeToDb(req.Type),
		EvalValue:   req.EvalValue,
		Response:    req.Response,
		ChannelID:   dashboardId,
	}

	if err := c.Db.
		WithContext(ctx).
		Create(entity).Error; err != nil {
		return nil, err
	}

	return c.convertEntity(entity), nil
}

func (c *Variables) VariablesDelete(ctx context.Context, req *variables.DeleteRequest) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	if err := c.Db.Model(&model.ChannelsCustomvars{}).
		WithContext(ctx).
		Where(`"channelId" = ? AND "id" = ?`, dashboardId, req.Id).
		Delete(nil).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Variables) VariablesUpdate(ctx context.Context, req *variables.PutRequest) (*variables.Variable, error) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := &model.ChannelsCustomvars{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? AND "id" = ?`, dashboardId, req.Id).
		First(entity).Error; err != nil {
		return nil, err
	}

	entity.Name = req.Variable.Name
	entity.Description = null.StringFromPtr(req.Variable.Description)
	entity.Type = c.rpcTypeToDb(req.Variable.Type)
	entity.EvalValue = req.Variable.EvalValue
	entity.Response = req.Variable.Response

	if err := c.Db.
		WithContext(ctx).
		Save(entity).Error; err != nil {
		return nil, err
	}

	return c.convertEntity(entity), nil
}
