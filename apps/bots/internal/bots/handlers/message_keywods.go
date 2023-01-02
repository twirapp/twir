package handlers

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/samber/lo"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/parser"
)

func (c *Handlers) handleKeywords(
	msg Message,
	userBadges []string,
) {
	keywords := []model.ChannelsKeywords{}
	err := c.db.Where(`"channelId" = ? AND "enabled" = ?`, msg.Channel.ID, true).Find(&keywords).Error
	if err != nil {
		fmt.Println(err)
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

			query := make(map[string]any)

			if !isOnCooldown && k.Response != "" {
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
					fmt.Println(err)
					return
				}

				for _, r := range res.Responses {
					validateResposeErr := ValidateResponseSlashes(r)
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
							lo.If(k.IsReply, &msg.ID).Else(nil),
						)
					}
				}

				query["cooldownExpireAt"] = time.Now().
					Add(time.Duration(k.Cooldown) * time.Second).
					UTC()
			}

			query["usages"] = k.Usages + timesInMessage
			c.db.Model(&k).Where("id = ?", k.ID).Select("*").Updates(query)
		}(k)
	}

	wg.Wait()
}
