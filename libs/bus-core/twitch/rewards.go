package twitch

import (
	"time"
)

const RedemptionAddSubject = "twitch.reward.activated"

type ActivatedRedemption struct {
	ID                   string
	BroadcasterUserID    string
	BroadcasterUserLogin string
	BroadcasterUserName  string
	UserID               string
	UserLogin            string
	UserName             string
	UserInput            string
	Status               string
	RedeemedAt           time.Time
	Reward               ActivatedRedemptionReward
}

type ActivatedRedemptionReward struct {
	ID     string
	Title  string
	Prompt string
	Cost   int
}
