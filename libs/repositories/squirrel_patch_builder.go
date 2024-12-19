package repositories

import (
	"github.com/Masterminds/squirrel"
)

func SquirrelApplyPatch(
	updateBuilder squirrel.UpdateBuilder,
	input map[string]interface{},
) squirrel.UpdateBuilder {
	for field, value := range input {
		if value != nil {
			updateBuilder = updateBuilder.Set(field, value)
		}
	}
	return updateBuilder
}
