package redis_keys

func StreamParsedMessages(streamID string) string {
	return "stream:parsedMessages:" + streamID
}

func StreamByChannelID(channelID string) string {
	return "cache:twir:streams:channel:" + channelID
}
