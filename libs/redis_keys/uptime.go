package redis_keys

const UptimeScanPattern = "uptime:*"

func UptimeInstanceKey(service string, instance string) string {
	return "uptime:" + service + ":" + instance
}
