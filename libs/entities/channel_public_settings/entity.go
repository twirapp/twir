package channelpublicsettings

import "github.com/google/uuid"

type SocialLink struct {
	ID    uuid.UUID
	Title string
	Href  string
}

type ChannelPublicSettings struct {
	ID          uuid.UUID
	ChannelID   uuid.UUID
	Description *string
	SocialLinks []SocialLink

	isNil bool
}

func (c ChannelPublicSettings) IsNil() bool {
	return c.isNil
}

var Nil = ChannelPublicSettings{isNil: true}
