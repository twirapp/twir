package natshandler

import (
	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/libs/nats/parser"
)

func (c *NatsServiceImpl) ParseResponse(data parser.ParseResponseRequest) string {
	isCommand := lo.IfF(data.ParseVariables != nil, func() bool {
		return *data.ParseVariables
	}).ElseF(func() bool { return false })

	cacheService := variables_cache.New(variables_cache.VariablesCacheOpts{
		Text:        &data.Message.Text,
		SenderId:    data.Sender.Id,
		ChannelName: data.Channel.Name,
		ChannelId:   data.Channel.Id,
		SenderName:  &data.Sender.DisplayName,
		Redis:       c.redis,
		Regexp:      nil,
		Twitch:      c.commands.Twitch,
		DB:          c.commands.Db,
		Nats:        c.commands.Nats,
		IsCommand:   isCommand,
	})

	return c.variables.ParseInput(cacheService, data.Message.Text)
}
