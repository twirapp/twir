package giveaways

const (
	TryAddParticipantSubject = "giveaways.try_add_participant"
	ChooseWinnerSubject      = "giveaways.choose_winner"
)

type TryAddParticipantRequest struct {
	UserID      string
	DisplayName string
	GiveawayID  string
}

type ChooseWinnerRequest struct {
	GiveawayID string
}

type ChooseWinnerResponse struct {
	Winners []Winner
}

type Winner struct {
	UserID      string
	DisplayName string
}
