package model

import (
	"time"

	"github.com/lib/pq"
)

type ChannelsStreams struct {
	ID             string          `gorm:"primary_key;column:id;type:TEXT;"            json:"id"`
	UserId         string          `gorm:"column:userId;type:TEXT;"                    json:"user_id"`
	UserLogin      string          `gorm:"column:userLogin;type:TEXT;"                 json:"user_login"`
	UserName       string          `gorm:"column:userName;type:TEXT;"                  json:"user_name"`
	GameId         string          `gorm:"column:gameId;type:TEXT;default:TEXT;"       json:"game_id"`
	GameName       string          `gorm:"column:gameName;type:TEXT;"                  json:"game_name"`
	CommunityIds   pq.StringArray  `gorm:"column:communityIds;type:text[];default:[];" json:"community_ids"`
	Type           string          `gorm:"column:type;type:TEXT;"                      json:"type"`
	Title          string          `gorm:"column:title;type:TEXT;"                     json:"title"`
	ViewerCount    int             `gorm:"column:viewerCount;type:INT;"                json:"viewer_count"`
	StartedAt      time.Time       `gorm:"column:startedAt;type:TIMESTAMP;"            json:"started_at"`
	Language       string          `gorm:"column:language;type:TEXT;"                  json:"language"`
	ThumbnailUrl   string          `gorm:"column:thumbnailUrl;type:TEXT;"              json:"thumbnail_url"`
	TagIds         *pq.StringArray `gorm:"column:tagIds;type:text[];default:[];"       json:"tag_ids"`
	IsMature       bool            `gorm:"column:isMature;type:BOOL;"                  json:"is_mature"`
	ParsedMessages int             `gorm:"column:parsedMessages;type:INT;"             json:"parsedMessages"`
	Channel        *Channels       `gorm:"foreignKey:UserId"                           json:"channel"`
}

func (ChannelsStreams) TableName() string {
	return "channels_streams"
}
