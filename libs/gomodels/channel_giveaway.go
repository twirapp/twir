package model

import (
	"github.com/guregu/null"
	"time"
)

type ChannelGiveAwayType string

func (c ChannelGiveAwayType) String() string {
	return string(c)
}

const (
	ChannelGiveAwayTypeByKeyword      ChannelGiveAwayType = "BY_KEYWORD"
	ChannelGiveAwayTypeByRandomNumber ChannelGiveAwayType = "BY_RANDOM_NUMBER"
)

type ChannelGiveaway struct {
	ID          string              `gorm:"primary_key;column:id;type:TEXT;" json:"id"`
	Description string              `gorm:"column:description;type:TEXT;" json:"description"`
	Type        ChannelGiveAwayType `gorm:"column:type;type:VARCHAR;" json:"type"`

	ChannelID string    `gorm:"column:channel_id;type:TEXT;" json:"channel_id"`
	Channel   *Channels `gorm:"foreignKey:ChannelID" json:"channel"`

	CreatedAt time.Time `gorm:"column:created_at;type:TIMESTAMP;" json:"created_at"`
	StartAt   time.Time `gorm:"column:start_at;type:TIMESTAMP;" json:"start_at"`
	EndAt     null.Time `gorm:"column:end_at;type:TIMESTAMP;" json:"end_at"`
	ClosedAt  null.Time `gorm:"column:closed_at;type:TIMESTAMP;" json:"closed_at"`

	RequiredMinWatchTime      null.Int `gorm:"column:required_min_watch_time;type:INTEGER;" json:"required_min_watch_time"`
	RequiredMinFollowTime     null.Int `gorm:"column:required_min_follow_time;type:INTEGER;" json:"required_min_follow_time"`
	RequiredMinMessages       null.Int `gorm:"column:required_min_messages;type:INTEGER;" json:"required_min_messages"`
	RequiredMinSubscriberTier null.Int `gorm:"column:required_min_subscriber_tier;type:INTEGER;" json:"required_min_subscriber_tier"`
	RequiredMinSubscribeTime  null.Int `gorm:"column:required_min_subscribe_time;type:INTEGER;" json:"required_min_subscribe_time"`

	EligibleUserGroups []string `gorm:"column:eligible_user_groups;type:TEXT;" json:"eligible_user_groups"`

	Keyword             null.String `gorm:"column:keyword;type:TEXT;" json:"keyword"`
	RandomNumberFrom    null.Int    `gorm:"column:random_number_from;type:INTEGER;" json:"random_number_from"`
	RandomNumberTo      null.Int    `gorm:"column:random_number_to;type:INTEGER;" json:"random_number_to"`
	WinningRandomNumber null.Int    `gorm:"column:winning_random_number;type:INTEGER;" json:"winning_random_number"`

	WinnersCount int `gorm:"column:winners_count;type:INTEGER;" json:"winners_count"`

	SubscribersLuck      null.Int `gorm:"column:subscribers_luck;type:INTEGER;" json:"subscribers_luck"`
	SubscribersTier1Luck null.Int `gorm:"column:subscribers_tier1_luck;type:INTEGER;" json:"subscribers_tier1_luck"`
	SubscribersTier2Luck null.Int `gorm:"column:subscribers_tier2_luck;type:INTEGER;" json:"subscribers_tier2_luck"`
	SubscribersTier3Luck null.Int `gorm:"column:subscribers_tier3_luck;type:INTEGER;" json:"subscribers_tier3_luck"`

	WatchedTimeLucks       []ChannelGiveawayConfigurableLuck `gorm:"column:watched_time_lucks;type:TEXT;" json:"watched_time_lucks"`
	MessagesLucks          []ChannelGiveawayConfigurableLuck `gorm:"column:messages_lucks;type:TEXT;" json:"messages_lucks"`
	UsedChannelPointsLucks []ChannelGiveawayConfigurableLuck `gorm:"column:used_channel_points_lucks;type:TEXT;" json:"used_channel_points_lucks"`
}

type ChannelGiveawayConfigurableLuck struct {
	Value int `json:"value"`
	Luck  int `json:"luck"`
}

func (g *ChannelGiveaway) TableName() string {
	return "channels_giveaways"
}
