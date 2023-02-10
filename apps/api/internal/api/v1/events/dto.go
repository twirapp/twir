package events

import model "github.com/satont/tsuwari/libs/gomodels"

type operationDto struct {
	Type   model.EventOperationType `validate:"required" json:"type"`
	Input  *string                  `json:"input"`
	Delay  *int64                   `validate:"lt=60" json:"delay"`
	Repeat *int64                   `validate:"gte=1,lt=10" json:"repeat"`
}

type eventDto struct {
	Type        string         `validate:"required" json:"type"`
	RewardID    *string        `json:"rewardId"`
	CommandID   *string        `json:"commandId"`
	Description *string        `json:"description"`
	Operations  []operationDto `validate:"dive"`
}
