package redis_keys

func TimersCurrentResponse(timerId string) string {
	return "timers:current_response:" + timerId
}
