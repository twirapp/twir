package voteban

import "github.com/google/uuid"

type Voteban struct {
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

	isNil bool
}

func (c Voteban) IsNil() bool {
	return c.isNil
}

type VotingMode string

const (
	VotingModeChat        VotingMode = "chat"
	VotingModeTwitchPolls VotingMode = "twitch_polls"
)

var Nil = Voteban{
	isNil: true,
}
