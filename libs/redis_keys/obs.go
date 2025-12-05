package redis_keys

func ObsOverlayConnection(channelID string) string {
	return "overlays:obs:connected:" + channelID
}
