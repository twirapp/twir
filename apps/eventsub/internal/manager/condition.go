package manager

import (
	model "github.com/twirapp/twir/libs/gomodels"
)

func GetTypeCondition(
	t model.EventsubConditionType,
	topic,
	channelID,
	botId string,
) map[string]string {
	switch t {
	case model.EventsubConditionTypeBroadcasterUserID:
		return map[string]string{
			"broadcaster_user_id": channelID,
		}
	case model.EventsubConditionTypeUserID:
		return map[string]string{
			"user_id": channelID,
		}
	case model.EventsubConditionTypeBroadcasterWithUserID:
		data := map[string]string{
			"broadcaster_user_id": channelID,
			"user_id":             botId,
		}
		if topic == "channel.follow" {
			data["user_id"] = channelID
		}
		return data
	case model.EventsubConditionTypeBroadcasterWithModeratorID:
		return map[string]string{
			"broadcaster_user_id": channelID,
			"moderator_user_id":   botId,
		}
	case model.EventsubConditionTypeToBroadcasterID:
		return map[string]string{
			"to_broadcaster_user_id": channelID,
		}
	default:
		return nil
	}
}
