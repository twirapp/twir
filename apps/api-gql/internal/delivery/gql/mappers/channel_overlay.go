package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	customoverlayentity "github.com/twirapp/twir/libs/entities/custom_overlay"
)

func ChannelOverlayLayerTypeEntityToGql(t customoverlayentity.ChannelOverlayType) gqlmodel.ChannelOverlayLayerType {
	switch t {
	case customoverlayentity.ChannelOverlayTypeHTML:
		return gqlmodel.ChannelOverlayLayerTypeHTML
	case customoverlayentity.ChannelOverlayTypeIMAGE:
		return gqlmodel.ChannelOverlayLayerTypeImage
	default:
		return gqlmodel.ChannelOverlayLayerTypeHTML
	}
}

func ChannelOverlayLayerTypeGqlToEntity(t gqlmodel.ChannelOverlayLayerType) customoverlayentity.ChannelOverlayType {
	switch t {
	case gqlmodel.ChannelOverlayLayerTypeHTML:
		return customoverlayentity.ChannelOverlayTypeHTML
	case gqlmodel.ChannelOverlayLayerTypeImage:
		return customoverlayentity.ChannelOverlayTypeIMAGE
	default:
		return customoverlayentity.ChannelOverlayTypeHTML
	}
}

func ChannelOverlayLayerSettingsEntityToGql(s customoverlayentity.ChannelOverlayLayerSettings) *gqlmodel.ChannelOverlayLayerSettings {
	return &gqlmodel.ChannelOverlayLayerSettings{
		HTMLOverlayHTML:                    s.HtmlOverlayHTML,
		HTMLOverlayCSS:                     s.HtmlOverlayCSS,
		HTMLOverlayJs:                      s.HtmlOverlayJS,
		HTMLOverlayDataPollSecondsInterval: s.HtmlOverlayDataPollSecondsInterval,
		ImageURL:                           s.ImageUrl,
	}
}

func ChannelOverlayLayerEntityToGql(l customoverlayentity.ChannelOverlayLayer) gqlmodel.ChannelOverlayLayer {
	return gqlmodel.ChannelOverlayLayer{
		ID:                      l.ID,
		Type:                    ChannelOverlayLayerTypeEntityToGql(l.Type),
		Settings:                ChannelOverlayLayerSettingsEntityToGql(l.Settings),
		OverlayID:               l.OverlayID,
		PosX:                    l.PosX,
		PosY:                    l.PosY,
		Width:                   l.Width,
		Height:                  l.Height,
		CreatedAt:               l.CreatedAt,
		UpdatedAt:               l.UpdatedAt,
		PeriodicallyRefetchData: l.PeriodicallyRefetchData,
		Rotation:                l.Rotation,
		Locked:                  l.Locked,
		Visible:                 l.Visible,
		Opacity:                 l.Opacity,
	}
}

func ChannelOverlayEntityToGql(o customoverlayentity.ChannelOverlay) gqlmodel.ChannelOverlay {
	layers := make([]gqlmodel.ChannelOverlayLayer, len(o.Layers))
	for i, l := range o.Layers {
		layers[i] = ChannelOverlayLayerEntityToGql(l)
	}

	return gqlmodel.ChannelOverlay{
		ID:        o.ID,
		ChannelID: o.ChannelID,
		Name:      o.Name,
		CreatedAt: o.CreatedAt,
		UpdatedAt: o.UpdatedAt,
		Width:     o.Width,
		Height:    o.Height,
		InstaSave: o.InstaSave,
		Layers:    layers,
	}
}
