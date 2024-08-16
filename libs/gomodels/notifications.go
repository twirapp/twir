package model

import (
	"database/sql"
	"time"

	"github.com/guregu/null"
	uuid "github.com/satori/go.uuid"
)

var (
	_ = time.Second
	_ = sql.LevelDefault
	_ = null.Bool{}
	_ = uuid.UUID{}
)

/*
DB Table Details
-------------------------------------


Table: notifications
[ 0] id                                             TEXT                 null: false  primary: true   isArray: false  auto: false  col: TEXT            len: -1      default: [gen_random_uuid()]
[ 1] imageSrc                                       TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []
[ 2] createdAt                                      TIMESTAMP            null: false  primary: false  isArray: false  auto: false  col: TIMESTAMP       len: -1      default: [CURRENT_TIMESTAMP]
[ 3] userId                                         TEXT                 null: true   primary: false  isArray: false  auto: false  col: TEXT            len: -1      default: []


JSON Sample
-------------------------------------
{    "id": "DuIiohglsRVQIhgkKFFGbtNdw",    "image_src": "UAorPURptbNKQdQQUJEoGxmQt",    "created_at": "2215-01-06T23:22:41.53073624+03:00",    "user_id": "dTqGMFCvscREQKqVfOSdJlbtU"}



*/

type Notifications struct {
	ID           string      `gorm:"primaryKey;AUTO_INCREMENT;column:id;type:TEXT;"            json:"id"`
	CreatedAt    time.Time   `gorm:"column:createdAt;type:TIMESTAMP;default:CURRENT_TIMESTAMP;" json:"created_at"`
	UserID       null.String `gorm:"column:userId;type:TEXT;"                                   json:"user_id"`
	Message      null.String `gorm:"column:message;type:TEXT;"                                  json:"message"`
	EditorJsJson null.String `gorm:"column:editor_js_json;type:TEXT;"                                  json:"editorJsJson"`
}

func (n *Notifications) TableName() string {
	return "notifications"
}
