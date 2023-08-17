package handlers

import (
	"fmt"
	"log/slog"
	"strings"
	"time"

	"github.com/satont/twir/libs/twitch"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/bots/internal/bots/handlers/moderation"
	model "github.com/satont/twir/libs/gomodels"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type handleResult struct {
	IsDelete bool
	Time     *int
	Message  string
}

type parsers struct {
	db *gorm.DB
}

var functionsMapping = map[string]func(
	c *parsers, settings *model.ChannelsModerationSettings, ircMsg Message, badges []string,
) *handleResult{
	"links":       (*parsers).linksParser,
	"blacklists":  (*parsers).blacklistsParser,
	"symbols":     (*parsers).symbolsParser,
	"longMessage": (*parsers).longMessageParser,
	"caps":        (*parsers).capsParser,
	"emotes":      (*parsers).emotesParser,
}

func (c *Handlers) moderateMessage(msg Message, badges []string) bool {
	if lo.Some(badges, []string{"BROADCASTER", "MODERATOR"}) {
		return false
	}

	settings := []model.ChannelsModerationSettings{}
	if err := c.db.Where(`"channelId" = ? AND "enabled" = ?`, msg.Channel.ID, true).Find(&settings).Error; err != nil {
		c.logger.Error("cannot find moderation settings", slog.Any("err", err))
		return false
	}

	prsrs := parsers{db: c.db}

	for _, s := range settings {
		if !s.Enabled {
			continue
		}
		if !s.Vips && lo.Contains(badges, "VIP") {
			continue
		}
		if !s.Subscribers && lo.Contains(badges, "SUBSCRIBER") {
			continue
		}
		function, ok := functionsMapping[s.Type]
		if !ok {
			continue
		}

		res := function(&prsrs, &s, msg, badges)

		if res != nil {
			twitchClient, err := twitch.NewBotClient(c.BotClient.Model.ID, c.cfg, c.tokensGrpc)
			if err != nil {
				return false
			}

			if res.IsDelete {
				res, err := twitchClient.DeleteChatMessage(
					&helix.DeleteChatMessageParams{
						BroadcasterID: msg.Channel.ID,
						ModeratorID:   c.BotClient.Model.ID,
						MessageID:     msg.ID,
					},
				)

				if res.StatusCode != 200 {
					c.logger.Error("cannot delete message", slog.String("err", res.ErrorMessage))
				}
				if err != nil {
					c.logger.Error("cannot delete message", slog.Any("err", err))
				}
			} else {
				if res.Time != nil {
					res, err := twitchClient.BanUser(
						&helix.BanUserParams{
							BroadcasterID: msg.Channel.ID,
							ModeratorId:   c.BotClient.Model.ID,
							Body: helix.BanUserRequestBody{
								Duration: *res.Time,
								Reason:   res.Message,
								UserId:   msg.User.ID,
							},
						},
					)

					if res.StatusCode != 200 {
						c.logger.Error(
							"cannot ban user",
							slog.String("err", res.ErrorMessage),
							slog.String("userId", msg.User.ID),
							slog.String("channelId", msg.Channel.ID),
						)
					}
					if err != nil {
						c.logger.Error(
							"cannot ban user",
							slog.Any("err", err),
							slog.String("userId", msg.User.ID),
							slog.String("channelId", msg.Channel.ID),
						)
					}
				}
			}

			return true
		}
	}

	return false
}

func (c *parsers) returnByWarnedState(
	reason string,
	userID string,
	settings *model.ChannelsModerationSettings,
) *handleResult {
	warning := model.ChannelModerationWarn{}
	err := c.db.
		Where(
			`"channelId" = ? AND "userId" = ? AND "reason" = ?`,
			settings.ChannelID,
			userID,
			reason,
		).
		First(&warning).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		fmt.Println(err)
		return nil
	}

	if warning.ID == "" {
		c.db.Save(
			&model.ChannelModerationWarn{
				ID:        uuid.NewV4().String(),
				ChannelID: settings.ChannelID,
				UserID:    userID,
				Reason:    reason,
			},
		)
		return &handleResult{
			IsDelete: true,
			Message:  settings.WarningMessage,
		}
	} else {
		duration := time.Duration(settings.BanTime) * time.Second
		c.db.Delete(&warning)
		return &handleResult{
			IsDelete: false,
			Time:     lo.ToPtr(int(duration.Seconds())),
			Message:  settings.BanMessage,
		}
	}
}

func (c *parsers) linksParser(
	settings *model.ChannelsModerationSettings,
	ircMsg Message,
	badges []string,
) *handleResult {
	containLink := moderation.HasLink(ircMsg.Message, true)

	if !containLink {
		return nil
	}

	permit := model.ChannelsPermits{}
	err := c.db.Where(`"channelId" = ? AND "userId" = ?`, settings.ChannelID, ircMsg.User.ID).
		Find(&permit).
		Error
	if err != nil {
		return nil
	}

	if permit.ID != "" {
		c.db.Delete(&permit)
		return nil
	} else {
		return c.returnByWarnedState(settings.Type, ircMsg.User.ID, settings)
	}
}

func (c *parsers) blacklistsParser(
	settings *model.ChannelsModerationSettings,
	ircMsg Message,
	badges []string,
) *handleResult {
	hasBlackListedWord := moderation.HasBlackListedWord(ircMsg.Message, settings.BlackListSentences)

	if hasBlackListedWord {
		return c.returnByWarnedState(settings.Type, ircMsg.User.ID, settings)
	}

	return nil
}

func (c *parsers) symbolsParser(
	settings *model.ChannelsModerationSettings,
	ircMsg Message,
	badges []string,
) *handleResult {
	if !settings.MaxPercentage.Valid {
		return nil
	}

	isTooMuchSymbols := moderation.IsToMuchSymbols(
		ircMsg.Message,
		int(settings.MaxPercentage.Int64),
	)

	if isTooMuchSymbols {
		return c.returnByWarnedState(settings.Type, ircMsg.User.ID, settings)
	}

	return nil
}

func (c *parsers) longMessageParser(
	settings *model.ChannelsModerationSettings,
	ircMsg Message,
	badges []string,
) *handleResult {
	if !settings.TriggerLength.Valid {
		return nil
	}

	isTooLong := moderation.IsTooLong(ircMsg.Message, int(settings.TriggerLength.Int64))

	if isTooLong {
		return c.returnByWarnedState(settings.Type, ircMsg.User.ID, settings)
	}

	return nil
}

func (c *parsers) capsParser(
	settings *model.ChannelsModerationSettings,
	ircMsg Message,
	badges []string,
) *handleResult {
	if !settings.MaxPercentage.Valid {
		return nil
	}

	msg := ircMsg.Message
	for _, v := range ircMsg.Emotes {
		msg = strings.ReplaceAll(msg, v.Name, "")
	}

	isTooLong := moderation.IsTooMuchCaps(msg, int(settings.MaxPercentage.Int64))

	if isTooLong {
		return c.returnByWarnedState(settings.Type, ircMsg.User.ID, settings)
	}

	return nil
}

func (c *parsers) emotesParser(
	settings *model.ChannelsModerationSettings,
	ircMsg Message,
	badges []string,
) *handleResult {
	if !settings.TriggerLength.Valid {
		return nil
	}

	length := 0

	for _, e := range ircMsg.Emotes {
		length += len(e.Positions)
	}

	if length < int(settings.TriggerLength.Int64) {
		return nil
	}

	return c.returnByWarnedState(settings.Type, ircMsg.User.ID, settings)
}
