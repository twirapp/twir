package entity

type EventsubSubscribeCondition string

const (
	EventsubSubscribeConditionChannel                EventsubSubscribeCondition = "CHANNEL"
	EventsubSubscribeConditionUser                   EventsubSubscribeCondition = "USER"
	EventsubSubscribeConditionChannelWithModeratorID EventsubSubscribeCondition = "CHANNEL_WITH_MODERATOR_ID"
	EventsubSubscribeConditionChannelWithBotID       EventsubSubscribeCondition = "CHANNEL_WITH_BOT_ID"
)
