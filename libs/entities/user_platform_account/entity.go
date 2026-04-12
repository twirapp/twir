package user_platform_account

import (
	"time"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/entities/platform"
)

type UserPlatformAccount struct {
	ID                  uuid.UUID
	UserID              uuid.UUID
	Platform            platform.Platform
	PlatformUserID      string
	PlatformLogin       string
	PlatformDisplayName string
	PlatformAvatar      string
	AccessToken         string
	RefreshToken        string
	Scopes              []string
	ExpiresIn           int
	ObtainmentTimestamp time.Time

	isNil bool
}

func (u UserPlatformAccount) IsNil() bool { return u.isNil }

var Nil = UserPlatformAccount{isNil: true}
