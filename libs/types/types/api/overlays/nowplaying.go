package overlays

type ChannelOverlayNowPlayingPreset string

const (
	ChannelOverlayNowPlayingPresetTransparent   ChannelOverlayNowPlayingPreset = "TRANSPARENT"
	ChannelOverlayNowPlayingPresetAidenRedesign ChannelOverlayNowPlayingPreset = "AIDEN_REDESIGN"
)

var AllPresets = []ChannelOverlayNowPlayingPreset{
	ChannelOverlayNowPlayingPresetTransparent,
	ChannelOverlayNowPlayingPresetAidenRedesign,
}

func (c ChannelOverlayNowPlayingPreset) String() string {
	return string(c)
}

func (c ChannelOverlayNowPlayingPreset) TSName() string {
	switch c {
	case ChannelOverlayNowPlayingPresetTransparent:
		return ChannelOverlayNowPlayingPresetTransparent.String()
	case ChannelOverlayNowPlayingPresetAidenRedesign:
		return ChannelOverlayNowPlayingPresetAidenRedesign.String()
	default:
		return ""
	}
}
