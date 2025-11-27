package model

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type VotingMode string

const (
	VotingModeChat        VotingMode = "chat"
	VotingModeTwitchPolls VotingMode = "twitch_polls"
)

type VoteBan struct {
	ID                       uuid.UUID      `db:"id"`
	ChannelID                string         `db:"channel_id"`
	Enabled                  bool           `db:"enabled"`
	TimeoutSeconds           int            `db:"timeout_seconds"`
	TimeoutModerators        bool           `db:"timeout_moderators"`
	InitMessage              string         `db:"init_message"`
	BanMessage               string         `db:"ban_message"`
	BanMessageModerators     string         `db:"ban_message_moderators"`
	SurviveMessage           string         `db:"survive_message"`
	SurviveMessageModerators string         `db:"survive_message_moderators"`
	NeededVotes              int            `db:"needed_votes"`
	VoteDuration             int            `db:"vote_duration"`
	VotingMode               VotingMode     `db:"voting_mode"`
	ChatVotesWordsPositive   pq.StringArray `db:"chat_votes_words_positive"`
	ChatVotesWordsNegative   pq.StringArray `db:"chat_votes_words_negative"`
}

var Nil = VoteBan{}
