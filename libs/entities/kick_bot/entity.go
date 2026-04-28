package kick_bot

import (
	"time"

	"github.com/google/uuid"
)

type KickBot struct {
	ID                  string
	Type                string
	AccessToken         string
	RefreshToken        string
	Scopes              []string
	ExpiresIn           int
	ObtainmentTimestamp time.Time
	KickUserID          uuid.UUID
	KickUserLogin       string
	CreatedAt           time.Time

	isNil bool
}

func (k KickBot) IsNil() bool { return k.isNil }

var Nil = KickBot{isNil: true}
