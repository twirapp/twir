package model

import (
	"github.com/guregu/null"
)

type UsersOnline struct {
	ID        string      `gorm:"primary_key;column:id;type:TEXT;"        json:"id"`
	ChannelId string      `gorm:"primary_key;column:channelId;type:TEXT;" json:"channelId"`
	UserId    null.String `gorm:"primary_key;column:userId;type:TEXT;"    json:"userId"`
	UserName  null.String `gorm:"primary_key;column:userName;type:TEXT;"  json:"userName"`
}

func (UsersOnline) TableName() string {
	return "users_online"
}
