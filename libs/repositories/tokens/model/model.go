package model

import (
	"time"

	"github.com/google/uuid"
)

type Token struct {
	ID                  uuid.UUID
	AccessToken         string
	RefreshToken        string
	ExpiresIn           int
	ObtainmentTimestamp time.Time
	Scopes              []string
}
