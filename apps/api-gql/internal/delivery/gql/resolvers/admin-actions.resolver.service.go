package resolvers

import (
	"fmt"

	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
)

func (r *mutationResolver) eventSubGqlToCondition(
	input gqlmodel.EventsubSubscribeConditionInput,
) (model.EventsubConditionType, error) {
	switch input {
	case gqlmodel.EventsubSubscribeConditionInputChannel:
		return model.EventsubConditionTypeBroadcasterUserID, nil
	case gqlmodel.EventsubSubscribeConditionInputUser:
		return model.EventsubConditionTypeUserID, nil
	case gqlmodel.EventsubSubscribeConditionInputChannelWithModeratorID:
		return model.EventsubConditionTypeBroadcasterWithModeratorID, nil
	case gqlmodel.EventsubSubscribeConditionInputChannelWithBotID:
		return model.EventsubConditionTypeBroadcasterWithUserID, nil
	default:
		return "", fmt.Errorf("unknown input type")
	}
}
