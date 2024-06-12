package model

import (
	"time"

	"github.com/lib/pq"
)

type ChannelsStreams struct {
	ID           string          `gorm:"primary_key;column:id;type:TEXT;"            json:"id"`
	UserId       string          `gorm:"column:userId;type:TEXT;"                    json:"userId"`
	UserLogin    string          `gorm:"column:userLogin;type:TEXT;"                 json:"userLogin"`
	UserName     string          `gorm:"column:userName;type:TEXT;"                  json:"userName"`
	GameId       string          `gorm:"column:gameId;type:TEXT;default:TEXT;"       json:"gameId"`
	GameName     string          `gorm:"column:gameName;type:TEXT;"                  json:"gameName"`
	CommunityIds pq.StringArray  `gorm:"column:communityIds;type:text[];default:[];" json:"communityIds"`
	Type         string          `gorm:"column:type;type:TEXT;"                      json:"type"`
	Title        string          `gorm:"column:title;type:TEXT;"                     json:"title"`
	ViewerCount  int             `gorm:"column:viewerCount;type:INT;"                json:"viewerCount"`
	StartedAt    time.Time       `gorm:"column:startedAt;type:TIMESTAMP;"            json:"startedAt"`
	Language     string          `gorm:"column:language;type:TEXT;"                  json:"language"`
	ThumbnailUrl string          `gorm:"column:thumbnailUrl;type:TEXT;"              json:"thumbnailUrl"`
	TagIds       *pq.StringArray `gorm:"column:tagIds;type:text[];default:[];"       json:"tagIds"`
	Tags         *pq.StringArray `gorm:"column:tags;type:text[];default:[];"         json:"tags"`
	IsMature     bool            `gorm:"column:isMature;type:BOOL;"                  json:"isMature"`
	Channel      *Channels       `gorm:"foreignKey:UserId"                           json:"channel"`
}

func (ChannelsStreams) TableName() string {
	return "channels_streams"
}
