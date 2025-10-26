package redis_keys

func CreateDistributedModTaskKey(channelId, userId string) string {
	return "mod_task:distributed:" + channelId + ":" + userId
}
