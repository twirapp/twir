package botssettings

import (
	"time"

	"github.com/google/uuid"
)

const UpdatePrefixSubject = "bots_settings.update_prefix"

type UpdatePrefixRequest struct {
	ID        uuid.UUID
	ChannelID string
	Prefix    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
