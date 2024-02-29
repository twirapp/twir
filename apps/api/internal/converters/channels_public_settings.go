package converters

import (
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/api/messages/channels_public_settings"
)

func ChannelsPublicSettingsModelToRpc(
	settings *model.ChannelPublicSettings,
) *channels_public_settings.Settings {
	links := make([]*channels_public_settings.SocialLink, 0, len(settings.SocialLinks))
	for _, link := range settings.SocialLinks {
		links = append(
			links,
			&channels_public_settings.SocialLink{
				Title: link.Title,
				Href:  link.Href,
			},
		)
	}

	return &channels_public_settings.Settings{
		Description: settings.Description.Ptr(),
		SocialLinks: links,
	}
}
