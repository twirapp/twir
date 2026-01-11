package streamlabsintegration

import (
	"time"

	"github.com/google/uuid"
)

type Entity struct {
	ID           uuid.UUID
	Enabled      bool
	ChannelID    string
	AccessToken  string
	RefreshToken string
	UserName     string
	Avatar       string
	CreatedAt    time.Time
	UpdatedAt    time.Time

	isNil bool
}

func (c Entity) IsNil() bool {
	return c.isNil
}

var Nil = Entity{
	isNil: true,
}
