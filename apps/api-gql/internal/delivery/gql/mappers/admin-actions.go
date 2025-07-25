package mappers

import (
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
)

var eventSubConditionTypeGqlToEntity = map[gqlmodel.EventsubSubscribeConditionInput]model.EventsubConditionType{
	gqlmodel.EventsubSubscribeConditionInputChannel:                model.EventsubConditionTypeBroadcasterUserID,
	gqlmodel.EventsubSubscribeConditionInputUser:                   model.EventsubConditionTypeUserID,
	gqlmodel.EventsubSubscribeConditionInputChannelWithModeratorID: model.EventsubConditionTypeBroadcasterWithModeratorID,
	gqlmodel.EventsubSubscribeConditionInputChannelWithBotID:       model.EventsubConditionTypeBroadcasterWithUserID,
}

func ConditionTypeGqlToEntity(conditionType gqlmodel.EventsubSubscribeConditionInput) model.EventsubConditionType {
	return eventSubConditionTypeGqlToEntity[conditionType]
}
