package handlers

import (
	"context"
	"github.com/gempir/go-twitch-irc/v3"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
)

func (c *Handlers) handleCommand(msg Message, userBadges []string) {
	requestStruct := &parser.ProcessCommandRequest{
		Sender: &parser.Sender{
			Id:          msg.User.ID,
			Name:        msg.User.Name,
			DisplayName: msg.User.DisplayName,
			Badges:      userBadges,
		},
		Channel: &parser.Channel{
			Id:   msg.Channel.ID,
			Name: msg.Channel.Name,
		},
		Message: &parser.Message{
			Id:   msg.ID,
			Text: msg.Message,
			Emotes: lo.Map(msg.Emotes, func(item *twitch.Emote, _ int) *parser.Message_Emote {
				return &parser.Message_Emote{
					Name:  item.Name,
					Id:    item.ID,
					Count: int64(item.Count),
					Positions: lo.Map(item.Positions, func(item twitch.EmotePosition, _ int) *parser.Message_EmotePosition {
						return &parser.Message_EmotePosition{
							Start: int64(item.Start),
							End:   int64(item.End),
						}
					}),
				}
			}),
		},
	}

	res, err := c.parserGrpc.ProcessCommand(context.Background(), requestStruct)
	if err != nil {
		return
	}

	if res.KeepOrder != nil && !*res.KeepOrder {
		for _, response := range res.Responses {
			r := response
			if r == "" || r == " " {
				continue
			}
			go c.BotClient.SayWithRateLimiting(
				msg.Channel.Name,
				r,
				lo.If(res.IsReply, lo.ToPtr(msg.ID)).Else(nil),
			)
		}
	} else {
		for _, r := range res.Responses {
			validateResposeErr := ValidateResponseSlashes(r)

			if r == "" || r == " " {
				continue
			}

			if validateResposeErr != nil {
				c.BotClient.SayWithRateLimiting(
					msg.Channel.Name,
					validateResposeErr.Error(),
					nil,
				)
			} else {
				c.BotClient.SayWithRateLimiting(
					msg.Channel.Name,
					r,
					lo.If(res.IsReply, &msg.ID).Else(nil),
				)
			}
		}
	}
}
