package custom_var

import (
	"context"
	"errors"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/eval"
)

var CustomVar = &types.Variable{
	Name:        "customvar",
	Description: lo.ToPtr("Custom variable"),
	Visible:     lo.ToPtr(false),
	Handler: func(
		ctx context.Context, parseCtx *types.VariableParseContext, variableData *types.VariableData,
	) (*types.VariableHandlerResult, error) {
		result := &types.VariableHandlerResult{}

		if variableData.Params == nil {
			return result, nil
		}

		v := &model.ChannelsCustomvars{}
		err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"channelId" = ? AND "name" = ?`, parseCtx.Channel.ID, variableData.Params).
			Find(v).Error
		if err != nil {
			parseCtx.Services.Logger.Sugar().Error(err)
			return result, nil
		}

		if v.ID == "" || (v.Response == "" && v.EvalValue == "") {
			return result, nil
		}

		if v.Type == model.CustomVarScript {
			res, err := parseCtx.Services.Bus.Eval.Evaluate.Request(
				ctx,
				eval.EvalRequest{
					Expression: v.EvalValue,
				},
			)

			if err != nil {
				parseCtx.Services.Logger.Sugar().Error(err)

				return nil, errors.New(
					"cannot evaluate variable. This is internal error, please report this bug",
				)
			}

			result.Result = res.Data.Result
		}

		if v.Type == model.CustomVarText || v.Type == model.CustomVarNumber {
			result.Result = v.Response
		}

		result.Result = lo.Substring(result.Result, 0, 1000)

		return result, nil
	},
}
