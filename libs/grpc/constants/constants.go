package constants

const (
	PARSER_SERVER_PORT = 46061 + iota
	BOTS_SERVER_PORT
	SCHEDULER_SERVER_PORT
	TIMERS_SERVER_PORT
	EVENTSUB_SERVER_PORT
	INTEGRATIONS_SERVER_PORT
	DOTA_SERVER_PORT
	EVAL_SERVER_PORT
	WATCHED_SERVER_PORT
	WEBSOCKET_SERVER_PORT
	TOKENS_SERVER_PORT
	EMOTES_CACHER_SERVER_PORT
	EVENTS_SERVER_PORT
	YTSR_SERVER_PORT
	DISCORD_SERVER_PORT
	LANGUAGE_DETECTOR_SERVER_PORT
)

var (
	ServerPorts = []int{
		PARSER_SERVER_PORT,
		BOTS_SERVER_PORT,
		SCHEDULER_SERVER_PORT,
		TIMERS_SERVER_PORT,
		EVENTSUB_SERVER_PORT,
		INTEGRATIONS_SERVER_PORT,
		DOTA_SERVER_PORT,
		EVAL_SERVER_PORT,
		WATCHED_SERVER_PORT,
		WEBSOCKET_SERVER_PORT,
		TOKENS_SERVER_PORT,
		EMOTES_CACHER_SERVER_PORT,
		EVENTS_SERVER_PORT,
		YTSR_SERVER_PORT,
		DISCORD_SERVER_PORT,
		LANGUAGE_DETECTOR_SERVER_PORT,
	}
)
