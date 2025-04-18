package giveaways

const (
	newParticipantsSubscriptionKey = "api-gql.giveaways.newParticipants"
)

func CreateNewPariticipantSubscriptionKeyByGiveawayID(giveawayID string) string {
	return newParticipantsSubscriptionKey + "." + giveawayID
}
