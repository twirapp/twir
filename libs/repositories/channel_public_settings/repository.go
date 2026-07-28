package channel_public_settings

import (
	"context"

	"github.com/google/uuid"
	channelpublicsettings "github.com/twirapp/twir/libs/entities/channel_public_settings"
)

type Repository interface {
	GetByChannelID(
		ctx context.Context,
		channelID uuid.UUID,
	) (channelpublicsettings.ChannelPublicSettings, error)
	Upsert(ctx context.Context, input UpsertInput) error
}

type SocialLinkInput struct {
	Title string
	Href  string
}

type UpsertInput struct {
	ChannelID uuid.UUID

	// Description is applied only when DescriptionSet is true.
	// A nil Description with DescriptionSet=true clears the description (sets NULL).
	Description    *string
	DescriptionSet bool

	// SocialLinks fully replaces existing links only when SocialLinksSet is true.
	SocialLinks    []SocialLinkInput
	SocialLinksSet bool
}
