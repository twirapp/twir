package model

import (
	"github.com/guregu/null"
)

type IgnoredUser struct {
	ID          string      `gorm:"primary_key;AUTO_INCREMENT;column:id;type:TEXT;" json:"id"`
	Login       null.String `gorm:"column:login;type:TEXT;"                         json:"login"`
	DisplayName null.String `gorm:"column:displayName;type:TEXT;"                   json:"displayName"`
	Force       bool        `gorm:"column:force;type:BOOLEAN;default:false;"        json:"force"`
}

func (i *IgnoredUser) TableName() string {
	return "users_ignored"
}
