package notification

import "time"

type Notification struct {
	ID                    string
	UserID                *string
	Text                  *string
	EditorJSJSON          *string
	CreatedAt             time.Time
	UpdatedAt             *time.Time
	DiscordMessageID      *string
	DiscordChannelID      *string
	DiscordAttachmentKeys []string

	isNil bool
}

func (n Notification) IsNil() bool {
	return n.isNil
}

var Nil = Notification{isNil: true}
