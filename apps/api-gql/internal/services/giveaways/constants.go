package giveaways

const (
	newParticipantsSubscriptionKey = "api-gql.giveaways.newParticipants"
)

func CreateNewParticipantSubscriptionKeyByGiveawayID(giveawayID string) string {
	return newParticipantsSubscriptionKey + "." + giveawayID
}
