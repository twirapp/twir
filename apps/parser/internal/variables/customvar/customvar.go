package customvar

import (
	"context"
	"errors"
	"fmt"
	"github.com/samber/do"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	"gorm.io/gorm"

	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/eval"

	"github.com/samber/lo"
)

var Variable = types.Variable{
	Name:        "customvar",
	Description: lo.ToPtr("Custom variable"),
	Visible:     lo.ToPtr(false),
	Handler: func(ctx *variables_cache.VariablesCacheService, data types.VariableHandlerParams) (*types.VariableHandlerResult, error) {
		evalGrpc := do.MustInvoke[eval.EvalClient](di.Provider)
		result := &types.VariableHandlerResult{}

		if data.Params == nil {
			return result, nil
		}

		v := getVarByName(*data.Params)

		if v == nil || v.Response == "" || v.EvalValue == "" {
			return result, nil
		}

		if v.Type == "SCRIPT" {
			req, err := evalGrpc.Process(context.Background(), &eval.Evaluate{
				Script: v.EvalValue,
			})
			if err != nil {
				return nil, errors.New(
					"cannot evaluate variable. This is internal error, please report this bug",
				)
			}

			result.Result = req.Result
		} else {
			result.Result = v.Response
		}

		return result, nil
	},
}

type CustomVar struct {
	Type      *string `json:"type"`
	EvalValue *string `json:"evalValue"`
	Response  *string `json:"response"`
}

func getVarByName(name string) *model.ChannelsCustomvars {
	db := do.MustInvoke[gorm.DB](di.Provider)

	variable := &model.ChannelsCustomvars{}
	err := db.Where(`"name" = ?`, name).First(variable).Error
	if err != nil {
		fmt.Println(err)
		return nil
	}

	return variable
}
