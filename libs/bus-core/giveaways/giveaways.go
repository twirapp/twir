package giveaways

import "github.com/oklog/ulid/v2"

const (
	TryAddParticipantSubject = "giveaways.try_add_participant"
	ChooseWinnerSubject      = "giveaways.choose_winner"
)

type TryAddParticipantRequest struct {
	UserID      string
	DisplayName string
	GiveawayID  ulid.ULID
}

type ChooseWinnerRequest struct {
	GiveawayID ulid.ULID
}

type ChooseWinnerResponse struct {
	Winners []Winner
}

type Winner struct {
	UserID      string
	DisplayName string
}
