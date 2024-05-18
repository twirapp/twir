package model

import (
	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ChannelGamesVoteBanVotingMode string

const (
	ChannelGamesVoteBanVotingModeChat        ChannelGamesVoteBanVotingMode = "chat"
	ChannelGamesVoteBanVotingModeTwitchPolls ChannelGamesVoteBanVotingMode = "twitch_polls"
)

type ChannelGamesVoteBan struct {
	ID                       uuid.UUID                     `gorm:"column:id;type:uuid;default:gen_random_uuid()"`
	ChannelID                string                        `gorm:"column:channel_id;type:text"`
	Enabled                  bool                          `gorm:"column:enabled;type:bool"`
	TimeoutSeconds           int8                          `gorm:"column:timeout_seconds;type:int2"`
	TimeoutModerators        bool                          `gorm:"column:timeout_moderators;type:bool"`
	InitMessage              string                        `gorm:"column:init_message;type:text"`
	BanMessage               string                        `gorm:"column:ban_message;type:text"`
	BanMessageModerators     string                        `gorm:"column:ban_message_moderators;type:text"`
	SurviveMessage           string                        `gorm:"column:survive_message;type:text"`
	SurviveMessageModerators string                        `gorm:"column:survive_message_moderators;type:text"`
	NeededVotes              int                           `gorm:"column:needed_votes;type:int"`
	VoteDuration             int                           `gorm:"column:vote_duration;type:int"`
	VotingMode               ChannelGamesVoteBanVotingMode `gorm:"column:voting_mode;type:channel_games_voteban_voting_mode"`
	ChatVotesWordsPositive   pq.StringArray                `gorm:"column:chat_votes_words_positive;type:text[]"`
	ChatVotesWordsNegative   pq.StringArray                `gorm:"column:chat_votes_words_negative;type:text[]"`
}

func (ChannelGamesVoteBan) TableName() string {
	return "channels_games_voteban"
}

type ChannelGamesVoteBanRedisStruct struct {
	TargetUserId   string `redis:"target_user_id"`
	TargetUserName string `redis:"target_user_name"`
	TotalVotes     int    `redis:"total_votes"`
	PositiveVotes  int    `redis:"positive_votes"`
	NegativeVotes  int    `redis:"negative_votes"`
}
