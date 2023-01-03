package bot

type BotInfo struct {
	IsMod   bool   `json:"isMod"`
	BotID   string `json:"botId"`
	BotName string `json:"botName"`
	Enabled bool   `json:"enabled"`
}

type Bot struct {
	GET BotInfo
}
