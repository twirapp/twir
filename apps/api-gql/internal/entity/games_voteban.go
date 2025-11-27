package entity

import (
	"github.com/google/uuid"
)

type VotingMode string

const (
	VotingModeChat        VotingMode = "chat"
	VotingModeTwitchPolls VotingMode = "twitch_polls"
)

type GamesVoteBan struct {
	ID                       uuid.UUID
	ChannelID                string
	Enabled                  bool
	TimeoutSeconds           int
	TimeoutModerators        bool
	InitMessage              string
	BanMessage               string
	BanMessageModerators     string
	SurviveMessage           string
	SurviveMessageModerators string
	NeededVotes              int
	VoteDuration             int
	VotingMode               VotingMode
	ChatVotesWordsPositive   []string
	ChatVotesWordsNegative   []string
}

var GamesVoteBanNil = GamesVoteBan{}
