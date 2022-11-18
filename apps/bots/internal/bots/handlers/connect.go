package handlers

func (c *Handlers) OnConnect() {
	c.logger.Sugar().
		Infow("Connected to twitch servers", "botName", c.BotClient.TwitchUser.Login, "botId", c.BotClient.Model.ID)
}
