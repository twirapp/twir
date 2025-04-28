package giveaways

const (
	TryAddParticipantSubject = "giveaways.try_add_participant"
	ChooseWinnerSubject      = "giveaways.choose_winner"
	NewParticipantsSubject   = "giveaways.new_participants"
)

type TryAddParticipantRequest struct {
	UserID          string
	UserLogin       string
	UserDisplayName string
	GiveawayID      string
}

type ChooseWinnerRequest struct {
	GiveawayID string
}

type ChooseWinnerResponse struct {
	Winners []Winner
}

type Winner struct {
	UserID          string
	UserLogin       string
	UserDisplayName string
}

type NewParticipant struct {
	GiveawayID      string
	UserID          string
	UserLogin       string
	UserDisplayName string
}
