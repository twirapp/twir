package mappers

import (
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/gqlmodel"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
)

func ChannelOverlayLayerTypeEntityToGql(t entity.ChannelOverlayType) gqlmodel.ChannelOverlayLayerType {
	switch t {
	case entity.ChannelOverlayTypeHTML:
		return gqlmodel.ChannelOverlayLayerTypeHTML
	default:
		return gqlmodel.ChannelOverlayLayerTypeHTML
	}
}

func ChannelOverlayLayerTypeGqlToEntity(t gqlmodel.ChannelOverlayLayerType) entity.ChannelOverlayType {
	switch t {
	case gqlmodel.ChannelOverlayLayerTypeHTML:
		return entity.ChannelOverlayTypeHTML
	default:
		return entity.ChannelOverlayTypeHTML
	}
}

func ChannelOverlayLayerSettingsEntityToGql(s entity.ChannelOverlayLayerSettings) *gqlmodel.ChannelOverlayLayerSettings {
	return &gqlmodel.ChannelOverlayLayerSettings{
		HTMLOverlayHTML:                    s.HtmlOverlayHTML,
		HTMLOverlayCSS:                     s.HtmlOverlayCSS,
		HTMLOverlayJs:                      s.HtmlOverlayJS,
		HTMLOverlayDataPollSecondsInterval: s.HtmlOverlayDataPollSecondsInterval,
	}
}

func ChannelOverlayLayerEntityToGql(l entity.ChannelOverlayLayer) gqlmodel.ChannelOverlayLayer {
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
	}
}

func ChannelOverlayEntityToGql(o entity.ChannelOverlay) gqlmodel.ChannelOverlay {
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
		Layers:    layers,
	}
}
