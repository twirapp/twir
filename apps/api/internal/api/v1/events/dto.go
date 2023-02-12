package events

import model "github.com/satont/tsuwari/libs/gomodels"

type operationDto struct {
	Type        model.EventOperationType `validate:"required" json:"type"`
	Input       *string                  `json:"input"`
	Delay       int                      `validate:"lt=600" json:"delay"`
	Repeat      int                      `validate:"gte=1,lt=10" json:"repeat"`
	UseAnnounce *bool                    `json:"useAnnounce"`
}

type eventDto struct {
	Type        string         `validate:"required" json:"type"`
	RewardID    *string        `json:"rewardId"`
	CommandID   *string        `json:"commandId"`
	Description string         `validate:"required" json:"description"`
	Operations  []operationDto `validate:"dive"`
}

type eventPatchDto struct {
	Enabled *bool `validate:"omitempty" json:"enabled,omitempty"`
}
