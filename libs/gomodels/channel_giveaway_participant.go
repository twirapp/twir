package model

import (
	"github.com/guregu/null"
)

type ChannelGiveawayParticipant struct {
	ID string `gorm:"primary_key;column:id;type:TEXT;" json:"id"`

	GiveawayID string           `gorm:"column:giveaway_id;type:TEXT;" json:"giveaway_id"`
	Giveaway   *ChannelGiveaway `gorm:"foreignKey:GiveawayID" json:"giveaway"`

	IsWinner bool `gorm:"column:is_winner;type:BOOLEAN;" json:"is_winner"`

	UserID string `gorm:"column:user_id;type:TEXT;" json:"user_id"`
	User   *Users `gorm:"foreignKey:UserID" json:"user"`

	IsSubscriber   bool `gorm:"column:is_subscriber;type:BOOLEAN;" json:"is_subscriber"`
	SubscriberTier int  `gorm:"column:subscriber_tier;type:INTEGER;" json:"subscriber_tier"`

	UserFollowSince      null.Time `gorm:"column:user_follow_since;type:TIMESTAMP;" json:"user_follow_since"`
	UserStatsWatchedTime int64     `gorm:"column:user_stats_watched_time;type:BIGINT;" json:"user_stats_watched_time"`
	UserStatsMessages    int       `gorm:"column:user_stats_messages;type:INTEGER;" json:"user_stats_messages"`
}
