package grpc

import (
	"encoding/json"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/valyala/fasttemplate"
	"time"
)

type Data struct {
	UserName    string `json:"userName,omitempty"`
	RaidViewers int    `json:"raidViewers,omitempty"`
}

func hydrateStringWithData(str string, data Data) (string, error) {
	template := fasttemplate.New(str, "{", "}")

	bytes, err := json.Marshal(data)
	if err != nil {
		return "", err
	}

	m := make(map[string]any)

	if err = json.Unmarshal(bytes, &m); err != nil {
		return "", err
	}

	s := template.ExecuteString(m)

	return s, nil
}

func (c *EventsGrpcImplementation) processOperations(operations []model.EventOperation, data Data) {
	for _, operation := range operations {
		go func(operation model.EventOperation) {
			if operation.Delay.Valid {
				duration := time.Duration(operation.Delay.Int64) * time.Second
				time.Sleep(duration)
			}

			switch operation.Type {
			case "SEND_MESSAGE":
				if operation.Input.Valid {
					hydrateStringWithData(operation.Input.String, Data{
						UserName:    data.UserName,
						RaidViewers: 0,
					})
				}
			}
		}(operation)
	}
}
