package baseapp

import (
	"context"
	"fmt"
	"reflect"

	"github.com/goccy/go-json"
	"gorm.io/gorm"
)

func getKeyFromData(key string, data map[string]interface{}) string {
	objId, ok := data[key]
	if !ok {
		return ""
	}
	return objId.(string)
}

func prepareData(data map[string]interface{}) string {
	dataByte, _ := json.Marshal(&data)
	return string(dataByte)
}

func getDataBeforeOperation(db *gorm.DB) (map[string]interface{}, error) {
	objMap := map[string]interface{}{}
	if db.Error == nil && !db.DryRun {
		primaryKeyValue := ""
		value := db.Statement.ReflectValue

		// Check if the value is a struct
		if value.Kind() == reflect.Struct {
			Id, ID := value.FieldByName("Id"), value.FieldByName("ID")
			if Id.IsValid() {
				primaryKeyValue = Id.String()
			} else if ID.IsValid() {
				primaryKeyValue = ID.String()
			}
		}

		if primaryKeyValue == "" {
			return nil, nil
		}

		// objectType := reflect.TypeOf(db.Statement.ReflectValue.Interface())
		// // Create a new instance of the object type
		// targetObj := reflect.New(objectType).Interface()
		// // Fetch the target object separately
		// if err := db.Session(&gorm.Session{SkipHooks: true, NewDB: true}).Where(
		// 	"id = ?",
		// 	primaryKeyValue,
		// ).First(&targetObj).Error; err != nil {
		// 	return nil, fmt.Errorf("gorm callback: error while finding target object: %s", err.Error())
		// }

		jsonBytes, err := json.Marshal(db.Statement.Model)
		if err != nil {
			return nil, fmt.Errorf("gorm callback: error while marshalling json data: %s", err.Error())
		}
		json.Unmarshal(jsonBytes, &objMap)
	}
	return objMap, nil
}

func getUserIDFromContext(ctx context.Context) *string {
	val, ok := ctx.Value(RequesterUserIdContextKey).(string)
	if ok {
		return &val
	}

	return nil
}

func getDashboardIDFromContext(ctx context.Context) *string {
	val, ok := ctx.Value(SelectedDashboardContextKey).(string)
	if ok {
		return &val
	}

	return nil
}
