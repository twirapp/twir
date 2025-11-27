package messagehandler

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/lib/pq"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/bots/internal/entity"
	"github.com/twirapp/twir/apps/bots/internal/services/keywords"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/libs/bus-core/events"
	"github.com/twirapp/twir/libs/bus-core/parser"
	deprecatedgormmodel "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/logger"
	rolesmodel "github.com/twirapp/twir/libs/repositories/roles/model"
	"github.com/twirapp/twir/libs/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func (c *MessageHandler) handleKeywords(ctx context.Context, msg handleMessage) error {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.String("function.name", utils.GetFuncName()))

	entities, err := c.keywordsService.GetManyByChannelID(ctx, msg.BroadcasterUserId)
	if err != nil {
		return err
	}

	if len(entities) == 0 {
		return nil
	}

	message := msg.Message.Text
	var messagesForSend []string

	matchedKeywords := make([]entity.Keyword, 0, len(entities))

	timesInMessage := map[string]int{}

	for _, k := range entities {
		if !k.Enabled {
			continue
		}

		if k.IsRegular {
			regx, err := regexp.Compile(k.Text)
			if err != nil {
				messagesForSend = append(
					messagesForSend,
					fmt.Sprintf("regular expression is wrong for keyword %s", k.Text),
				)
				continue
			}

			if !regx.MatchString(message) {
				continue
			} else {
				timesInMessage[k.ID.String()] = len(regx.FindAllStringSubmatch(message, -1))
			}
		} else {
			if !strings.Contains(strings.ToLower(message), strings.ToLower(k.Text)) {
				continue
			} else {
				timesInMessage[k.ID.String()] = strings.Count(
					strings.ToLower(message),
					strings.ToLower(k.Text),
				)
			}
		}

		isOnCooldown := false
		if k.Cooldown != 0 && k.CooldownExpireAt != nil {
			isOnCooldown = k.CooldownExpireAt.After(time.Now().UTC())
		}

		if isOnCooldown {
			continue
		}

		matchedKeywords = append(matchedKeywords, k)
	}

	var wg sync.WaitGroup

	for _, k := range matchedKeywords {
		wg.Add(1)

		go func() {
			defer wg.Done()

			if len(k.RolesIDs) > 0 {
				channelRoles, err := c.keywordsService.GetChannelRoles(ctx, msg.BroadcasterUserId)
				if err != nil {
					c.logger.Error(
						"cannot get channel roles",
						logger.Error(err),
						slog.String("channelId", msg.BroadcasterUserId),
					)
					return
				}

				userBadges := createUserBadges(msg.Badges)
				hasRole := slices.Contains(userBadges, "BROADCASTER")

				if msg.EnrichedData.DbUser.IsBotAdmin {
					hasRole = true
				}

				for _, badge := range userBadges {
					for _, r := range channelRoles {
						if r.Type != rolesmodel.ChannelRoleTypeCustom && badge == r.Type.String() {
							hasRole = true
							break
						}
					}
				}

				if !hasRole {
					userRoles, err := c.keywordsService.GetUserAccessibleRoles(
						ctx,
						msg.BroadcasterUserId,
						msg.ChatterUserId,
					)
					if err != nil {
						c.logger.Error(
							"cannot get user roles",
							logger.Error(err),
							slog.String("channelId", msg.BroadcasterUserId),
							slog.String("userId", msg.ChatterUserId),
						)
						return
					}

					for _, r := range userRoles {
						for _, id := range k.RolesIDs {
							fmt.Println(r.ID.String(), id.String())
							if id.String() == r.ID.String() {
								hasRole = true
								break
							}
						}
					}
				}

				if !hasRole {
					return
				}
			}

			response := c.keywordsParseResponse(ctx, msg, k)

			c.keywordsTriggerEvent(ctx, msg, k, response)
			c.twitchActions.SendMessage(
				ctx, twitchactions.SendMessageOpts{
					BroadcasterID:        msg.BroadcasterUserId,
					SenderID:             msg.EnrichedData.DbChannel.BotID,
					Message:              response,
					ReplyParentMessageID: lo.If(k.IsReply, msg.MessageId).Else(""),
				},
			)
			c.keywordsIncrementStats(ctx, k, timesInMessage[k.ID.String()])
			c.keywordsTriggerAlert(ctx, k)
		}()
	}

	wg.Wait()

	return nil
}

func (c *MessageHandler) keywordsIncrementStats(
	ctx context.Context,
	keyword entity.Keyword,
	count int,
) {
	input := keywords.UpdateInput{}

	usages := keyword.Usages + count
	input.Usages = &usages

	if keyword.Cooldown != 0 {
		expires := time.Now().
			Add(time.Duration(keyword.Cooldown) * time.Second).
			UTC()
		input.CooldownExpireAt = &expires
	}

	_, err := c.keywordsService.Update(ctx, keyword.ID, keyword.ChannelID, input)
	if err != nil {
		c.logger.Error(
			"cannot update keyword usages",
			logger.Error(err),
			slog.String("channelId", keyword.ChannelID),
		)
	}
}

func (c *MessageHandler) keywordsTriggerEvent(
	ctx context.Context,
	msg handleMessage,
	keyword entity.Keyword,
	response string,
) {
	err := c.twirBus.Events.KeywordMatched.Publish(
		ctx,
		events.KeywordMatchedMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   msg.BroadcasterUserId,
				ChannelName: msg.BroadcasterUserLogin,
			},
			KeywordID:       keyword.ID.String(),
			KeywordName:     keyword.Text,
			KeywordResponse: response,
			UserID:          msg.ChatterUserId,
			UserName:        msg.ChatterUserLogin,
			UserDisplayName: msg.ChatterUserName,
		},
	)
	if err != nil {
		c.logger.Error(
			"cannot send keywords matched event",
			logger.Error(err),
			slog.String("channelId", msg.BroadcasterUserId),
			slog.String("userId", msg.ChatterUserId),
		)
	}
}

func (c *MessageHandler) keywordsParseResponse(
	ctx context.Context,
	msg handleMessage,
	keyword entity.Keyword,
) string {
	if keyword.Response == "" {
		return ""
	}

	res, err := c.twirBus.Parser.ParseVariablesInText.Request(
		ctx, parser.ParseVariablesInTextRequest{
			ChannelID:   msg.BroadcasterUserId,
			ChannelName: msg.BroadcasterUserLogin,
			Text:        keyword.Response,
			UserID:      msg.ChatterUserId,
			UserLogin:   msg.ChatterUserLogin,
			UserName:    msg.ChatterUserName,
		},
	)
	if err != nil {
		c.logger.Error(
			"cannot parse keyword response",
			logger.Error(err),
			slog.String("channelId", msg.BroadcasterUserId),
		)
	}

	return res.Data.Text
}

func (c *MessageHandler) keywordsTriggerAlert(
	ctx context.Context,
	keyword entity.Keyword,
) {
	alert := deprecatedgormmodel.ChannelAlert{}
	if err := c.gorm.WithContext(ctx).Where(
		"channel_id = ? AND keywords_ids && ?",
		keyword.ChannelID,
		pq.StringArray{keyword.ID.String()},
	).Find(&alert).Error; err != nil {
		c.logger.Error(
			"cannot get alert",
			logger.Error(err),
			slog.String("channelId", keyword.ChannelID),
		)
		return
	}

	if alert.ID == "" {
		return
	}

	if _, err := c.websocketsGrpc.TriggerAlert(
		context.Background(),
		&websockets.TriggerAlertRequest{
			ChannelId: keyword.ChannelID,
			AlertId:   alert.ID,
		},
	); err != nil {
		c.logger.Error(
			"cannot trigger alert",
			logger.Error(err),
			slog.String("channelId", keyword.ChannelID),
		)
	}
}
