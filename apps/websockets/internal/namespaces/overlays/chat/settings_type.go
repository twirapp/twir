package chat

import (
	"encoding/json"
	"reflect"
	"strings"

	"github.com/nicklaw5/helix/v2"
	model "github.com/satont/twir/libs/gomodels"
)

type settings struct {
	model.ChatOverlaySettings
	ChannelID          string            `json:"channelId"`
	ChannelName        string            `json:"channelName"`
	ChannelDisplayName string            `json:"channelDisplayName"`
	GlobalBadges       []helix.ChatBadge `json:"globalBadges"`
	ChannelBadges      []helix.ChatBadge `json:"channelBadges"`
}

type field struct {
	Name  string
	Value interface{}
}

func (u settings) MarshalJSON() ([]byte, error) {
	var fields []field

	v := reflect.ValueOf(u)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)

		if fieldType.Anonymous {
			// If the field is anonymous (embedded), process its fields
			embeddedFields, err := processEmbeddedField(fieldValue.Interface())
			if err != nil {
				return nil, err
			}
			fields = append(fields, embeddedFields...)
		} else {
			fields = append(
				fields, field{
					Name:  fieldType.Tag.Get("json"),
					Value: fieldValue.Interface(),
				},
			)
		}
	}

	data := make(map[string]interface{})
	for _, f := range fields {
		data[f.Name] = f.Value
	}

	return json.Marshal(data)
}

func processEmbeddedField(embedded interface{}) ([]field, error) {
	var fields []field

	v := reflect.ValueOf(embedded)
	t := v.Type()

	for i := 0; i < v.NumField(); i++ {
		fieldValue := v.Field(i)
		fieldType := t.Field(i)
		tag := strings.Split(fieldType.Tag.Get("json"), ",")

		fields = append(
			fields, field{
				Name:  tag[0],
				Value: fieldValue.Interface(),
			},
		)
	}

	return fields, nil
}
