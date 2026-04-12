package kick_bot

import "time"

type KickBot struct {
	ID                  string
	Type                string
	AccessToken         string
	RefreshToken        string
	Scopes              []string
	ExpiresIn           int
	ObtainmentTimestamp time.Time
	KickUserID          string
	KickUserLogin       string
	CreatedAt           time.Time

	isNil bool
}

func (k KickBot) IsNil() bool { return k.isNil }

var Nil = KickBot{isNil: true}
