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
	case customoverlayentity.ChannelOverlayTypeTEXT:
		return gqlmodel.ChannelOverlayLayerTypeText
	case customoverlayentity.ChannelOverlayTypeVIDEO:
		return gqlmodel.ChannelOverlayLayerTypeVideo
	case customoverlayentity.ChannelOverlayTypeIFRAME:
		return gqlmodel.ChannelOverlayLayerTypeIframe
	case customoverlayentity.ChannelOverlayTypeYOUTUBE:
		return gqlmodel.ChannelOverlayLayerTypeYoutube
	case customoverlayentity.ChannelOverlayTypeEMOTE:
		return gqlmodel.ChannelOverlayLayerTypeEmote
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
	case gqlmodel.ChannelOverlayLayerTypeText:
		return customoverlayentity.ChannelOverlayTypeTEXT
	case gqlmodel.ChannelOverlayLayerTypeVideo:
		return customoverlayentity.ChannelOverlayTypeVIDEO
	case gqlmodel.ChannelOverlayLayerTypeIframe:
		return customoverlayentity.ChannelOverlayTypeIFRAME
	case gqlmodel.ChannelOverlayLayerTypeYoutube:
		return customoverlayentity.ChannelOverlayTypeYOUTUBE
	case gqlmodel.ChannelOverlayLayerTypeEmote:
		return customoverlayentity.ChannelOverlayTypeEMOTE
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
		TextContent:                        s.TextContent,
		TextFontFamily:                     s.TextFontFamily,
		TextFontSize:                       s.TextFontSize,
		TextFontWeight:                     s.TextFontWeight,
		TextColor:                          s.TextColor,
		TextAlign:                          s.TextAlign,
		VideoURL:                           s.VideoUrl,
		VideoLoop:                          s.VideoLoop,
		VideoMuted:                         s.VideoMuted,
		IframeURL:                          s.IframeUrl,
		IframeScale:                        s.IframeScale,
		WidgetKey:                          s.WidgetKey,
		YoutubeVideoID:                     s.YoutubeVideoID,
		YoutubeAutoplay:                    s.YoutubeAutoplay,
		YoutubeLoop:                        s.YoutubeLoop,
		YoutubeMuted:                       s.YoutubeMuted,
		EmoteURL:                           s.EmoteUrl,
		EmoteName:                          s.EmoteName,
		EmoteProvider:                      s.EmoteProvider,
	}
}

func ChannelOverlayLayerEntityToGql(l customoverlayentity.ChannelOverlayLayer) gqlmodel.ChannelOverlayLayer {
	return gqlmodel.ChannelOverlayLayer{
		ID:                      l.ID,
		Type:                    ChannelOverlayLayerTypeEntityToGql(l.Type),
		Name:                    l.Name,
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
