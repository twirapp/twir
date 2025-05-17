package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func VariableModelToGql(variable entity.CustomVariable) gqlmodel.Variable {
	return gqlmodel.Variable{
		ID:             variable.ID,
		Name:           variable.Name,
		Description:    variable.Description,
		Type:           gqlmodel.VariableType(variable.Type),
		EvalValue:      variable.EvalValue,
		Response:       variable.Response,
		ScriptLanguage: VariableScriptLanguageToGql(variable.ScriptLanguage),
	}
}

var variableScriptLanguageToGql = map[entity.CustomVarScriptLanguage]gqlmodel.VariableScriptLanguage{
	entity.ScriptLanguageJavaScript: gqlmodel.VariableScriptLanguageJavascript,
	entity.ScriptLanguagePython:     gqlmodel.VariableScriptLanguagePython,
}

func VariableScriptLanguageToGql(language entity.CustomVarScriptLanguage) gqlmodel.VariableScriptLanguage {
	return variableScriptLanguageToGql[language]
}

func VariableScriptLanguageToEntity(language gqlmodel.VariableScriptLanguage) entity.CustomVarScriptLanguage {
	for k, v := range variableScriptLanguageToGql {
		if v == language {
			return k
		}
	}

	return entity.ScriptLanguageJavaScript
}
