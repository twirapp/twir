package repositories

import (
	"reflect"

	"github.com/Masterminds/squirrel"
)

func SquirrelApplyPatch(
	updateBuilder squirrel.UpdateBuilder,
	input map[string]interface{},
) squirrel.UpdateBuilder {
	for field, value := range input {
		// Check if the interface{} is nil or if the value inside it is nil
		if value != nil && !isNil(value) {
			updateBuilder = updateBuilder.Set(field, value)
		}
	}
	return updateBuilder
}

// isNil checks if an interface value is nil or contains a nil dynamic value
func isNil(value interface{}) bool {
	// Use reflection to detect nil dynamic value
	v := reflect.ValueOf(value)
	return !v.IsValid() || (v.Kind() == reflect.Ptr && v.IsNil())
}
