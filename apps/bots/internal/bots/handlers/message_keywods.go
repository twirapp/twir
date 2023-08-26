package handlers

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/satont/twir/libs/grpc/generated/events"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/parser"
)

func (c *Handlers) handleKeywords(
	msg *Message,
	userBadges []string,
) {
	var keywords []model.ChannelsKeywords
	err := c.db.Where(`"channelId" = ? AND "enabled" = ?`, msg.Channel.ID, true).Find(&keywords).Error
	if err != nil {
		c.logger.Error(
			"cannot get keywords",
			slog.Any("err", err),
			slog.String("channelId", msg.Channel.ID),
		)
		return
	}

	if len(keywords) == 0 {
		return
	}

	wg := sync.WaitGroup{}
	wg.Add(len(keywords))

	message := msg.Message
	var timesInMessage int

	for _, k := range keywords {
		go func(k model.ChannelsKeywords) {
			defer wg.Done()

			if k.IsRegular {
				regx, err := regexp.Compile(k.Text)
				if err != nil {
					c.BotClient.SayWithRateLimiting(
						msg.Channel.Name,
						fmt.Sprintf("regular expression is wrong for keyword %s", k.Text),
						nil,
					)
					return
				}

				if !regx.MatchString(message) {
					return
				} else {
					timesInMessage = len(regx.FindAllStringSubmatch(message, -1))
				}
			} else {
				if !strings.Contains(strings.ToLower(message), strings.ToLower(k.Text)) {
					return
				} else {
					timesInMessage = strings.Count(strings.ToLower(message), strings.ToLower(k.Text))
				}
			}

			defer c.keywordsCounter.Inc()

			isOnCooldown := false
			if k.Cooldown != 0 && k.CooldownExpireAt.Valid {
				isOnCooldown = k.CooldownExpireAt.Time.After(time.Now().UTC())
			}

			if isOnCooldown {
				return
			}

			query := make(map[string]any)

			var responses []string
			if k.Response != "" {
				requestStruct := &parser.ParseTextRequestData{
					Channel: &parser.Channel{
						Id:   msg.Channel.ID,
						Name: msg.Channel.Name,
					},
					Sender: &parser.Sender{
						Id:          msg.User.ID,
						Name:        msg.User.Name,
						DisplayName: msg.User.DisplayName,
						Badges:      userBadges,
					},
					Message: &parser.Message{
						Text: k.Response,
						Id:   msg.ID,
					},
					ParseVariables: lo.ToPtr(true),
				}

				res, err := c.parserGrpc.ParseTextResponse(context.Background(), requestStruct)
				if err != nil {
					c.logger.Error(
						"cannot parse keyword response",
						slog.Any("err", err),
						slog.String("channelId", msg.Channel.ID),
					)
				}

				responses = res.Responses
			}

			_, err = c.eventsGrpc.KeywordMatched(
				context.Background(),
				&events.KeywordMatchedMessage{
					BaseInfo:        &events.BaseInfo{ChannelId: msg.Channel.ID},
					KeywordId:       k.ID,
					KeywordName:     k.Text,
					KeywordResponse: strings.Join(responses, " "),
					UserId:          msg.User.ID,
					UserName:        msg.User.Name,
					UserDisplayName: msg.User.DisplayName,
				},
			)
			if err != nil {
				c.logger.Error(
					"cannot send keywords matched event",
					slog.Any("err", err),
					slog.String("channelId", msg.Channel.ID),
					slog.String("userId", msg.User.ID),
				)
			}

			for _, r := range responses {
				validateResponseErr := ValidateResponseSlashes(r)
				if validateResponseErr != nil {
					c.BotClient.SayWithRateLimiting(
						msg.Channel.Name,
						validateResponseErr.Error(),
						nil,
					)
				} else {
					c.BotClient.SayWithRateLimiting(
						msg.Channel.Name,
						r,
						lo.If(k.IsReply, &msg.ID).Else(nil),
					)
				}
			}

			query["cooldownExpireAt"] = time.Now().
				Add(time.Duration(k.Cooldown) * time.Second).
				UTC()

			query["usages"] = k.Usages + timesInMessage
			err = c.db.Model(&k).Where("id = ?", k.ID).Select("*").Updates(query).Error
			if err != nil {
				c.logger.Error(
					"cannot update keyword usages",
					slog.Any("err", err),
					slog.String("channelId", msg.Channel.ID),
				)
			}
		}(k)
	}

	wg.Wait()
}
