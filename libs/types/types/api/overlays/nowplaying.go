package overlays

type ChannelOverlayNowPlayingPreset string

const (
	ChannelOverlayNowPlayingPresetTransparent   ChannelOverlayNowPlayingPreset = "TRANSPARENT"
	ChannelOverlayNowPlayingPresetAidenRedesign ChannelOverlayNowPlayingPreset = "AIDEN_REDESIGN"
	ChannelOverlayNowPLayingPresetSimpleLine    ChannelOverlayNowPlayingPreset = "SIMPLE_LINE"
	ChannelOverlayNowPlayingPresetPulseStrip    ChannelOverlayNowPlayingPreset = "PULSE_STRIP"
	ChannelOverlayNowPlayingPresetAuraStack     ChannelOverlayNowPlayingPreset = "AURA_STACK"
	ChannelOverlayNowPlayingPresetVinylHaze     ChannelOverlayNowPlayingPreset = "VINYL_HAZE"
	ChannelOverlayNowPlayingPresetSignalDeck    ChannelOverlayNowPlayingPreset = "SIGNAL_DECK"
)

var AllPresets = []ChannelOverlayNowPlayingPreset{
	ChannelOverlayNowPlayingPresetTransparent,
	ChannelOverlayNowPlayingPresetAidenRedesign,
	ChannelOverlayNowPLayingPresetSimpleLine,
	ChannelOverlayNowPlayingPresetPulseStrip,
	ChannelOverlayNowPlayingPresetAuraStack,
	ChannelOverlayNowPlayingPresetVinylHaze,
	ChannelOverlayNowPlayingPresetSignalDeck,
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
	case ChannelOverlayNowPLayingPresetSimpleLine:
		return ChannelOverlayNowPLayingPresetSimpleLine.String()
	case ChannelOverlayNowPlayingPresetPulseStrip:
		return ChannelOverlayNowPlayingPresetPulseStrip.String()
	case ChannelOverlayNowPlayingPresetAuraStack:
		return ChannelOverlayNowPlayingPresetAuraStack.String()
	case ChannelOverlayNowPlayingPresetVinylHaze:
		return ChannelOverlayNowPlayingPresetVinylHaze.String()
	case ChannelOverlayNowPlayingPresetSignalDeck:
		return ChannelOverlayNowPlayingPresetSignalDeck.String()
	default:
		return ""
	}
}
