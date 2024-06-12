package redis_keys

func StreamParsedMessages(streamID string) string {
	return "stream:parsedMessages:" + streamID
}
