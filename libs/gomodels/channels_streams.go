package model

import "github.com/lib/pq"

type ChannelsStreams struct {
	ID             string          `gorm:"primary_key;column:id;type:TEXT;"            json:"id"`
	UserId         string          `gorm:"column:userId;type:TEXT;"                    json:"userId"`
	UserLogin      string          `gorm:"column:userLogin;type:TEXT;"                 json:"userLogin"`
	UserName       bool            `gorm:"column:userName;type:TEXT;"                  json:"userName"`
	GameId         int             `gorm:"column:gameId;type:INT;default:[];"          json:"gameId"`
	GameName       bool            `gorm:"column:gameName;type:TEXT;"                  json:"gameName"`
	CommunityIds   pq.StringArray  `gorm:"column:communityIds;type:text[];default:[];" json:"communityIds"`
	Type           string          `gorm:"column:type;type:TEXT;"                      json:"type"`
	Title          string          `gorm:"column:title;type:TEXT;"                     json:"title"`
	ViewerCount    int             `gorm:"column:viewerCount;type:INT;default:[];"     json:"viewerCount"`
	StartedAt      string          `gorm:"column:startedAt;type:TEXT;"                 json:"startedAt"`
	Language       string          `gorm:"column:language;type:TEXT;"                  json:"language"`
	ThumbnailUrl   string          `gorm:"column:thumbnailUrl;type:TEXT;"              json:"thumbnailUrl"`
	TagIds         *pq.StringArray `gorm:"column:tagIds;type:text[];default:[];"       json:"tagIds"`
	IsMature       bool            `gorm:"column:isMature;type:TEXT;"                  json:"isMature"`
	ParsedMessages int             `gorm:"column:parsedMessages;type:INT;default:[];"  json:"parsedMessages"`
	Channel        *Channels       `gorm:"foreignKey:userId"                           json:"channel"`
}
