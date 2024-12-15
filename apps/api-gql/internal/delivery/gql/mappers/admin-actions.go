package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

var eventSubConditionTypeGqlToEntity = map[gqlmodel.EventsubSubscribeConditionInput]entity.EventsubSubscribeCondition{
	gqlmodel.EventsubSubscribeConditionInputChannel:                entity.EventsubSubscribeConditionChannel,
	gqlmodel.EventsubSubscribeConditionInputUser:                   entity.EventsubSubscribeConditionUser,
	gqlmodel.EventsubSubscribeConditionInputChannelWithModeratorID: entity.EventsubSubscribeConditionChannelWithModeratorID,
	gqlmodel.EventsubSubscribeConditionInputChannelWithBotID:       entity.EventsubSubscribeConditionChannelWithBotID,
}

func ConditionTypeGqlToEntity(conditionType gqlmodel.EventsubSubscribeConditionInput) entity.EventsubSubscribeCondition {
	return eventSubConditionTypeGqlToEntity[conditionType]
}
