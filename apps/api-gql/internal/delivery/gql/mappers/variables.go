package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func VariableModelToGql(variable entity.CustomVariable) gqlmodel.Variable {
	return gqlmodel.Variable{
		ID:          variable.ID,
		Name:        variable.Name,
		Description: variable.Description,
		Type:        gqlmodel.VariableType(variable.Type),
		EvalValue:   variable.EvalValue,
		Response:    variable.Response,
	}
}
