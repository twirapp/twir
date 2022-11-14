package handlers

import (
	"fmt"
	"strings"
	"time"

	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/samber/lo"
	"github.com/satont/go-helix/v2"
	"github.com/satont/tsuwari/apps/bots/internal/bots/handlers/moderation"
	model "github.com/satont/tsuwari/libs/gomodels"
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

var functionsMapping = map[string]func(c *parsers, settings *model.ChannelsModerationSettings, ircMsg irc.PrivateMessage, badges []string) *handleResult{
	"links":       (*parsers).linksParser,
	"blacklists":  (*parsers).blacklistsParser,
	"symbols":     (*parsers).symbolsParser,
	"longMessage": (*parsers).longMessageParser,
	"caps":        (*parsers).capsParser,
	"emotes":      (*parsers).emotesParser,
}

func (c *Handlers) moderateMessage(msg irc.PrivateMessage, badges []string) bool {
	if lo.Some(badges, []string{"BROADCASTER", "MODERATOR"}) {
		return false
	}

	settings := []model.ChannelsModerationSettings{}
	if err := c.db.Where(`"channelId" = ? AND "enabled" = ?`, msg.RoomID, true).Find(&settings).Error; err != nil {
		c.logger.Sugar().Error(err)
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
			if res.IsDelete {
				res, err := c.BotClient.Api.Client.DeleteMessage(&helix.DeleteMessageParams{
					BroadcasterID: msg.RoomID,
					ModeratorID:   c.BotClient.Model.ID,
					MessageID:     msg.ID,
				})

				if res.StatusCode != 200 {
					c.logger.Sugar().Info(res.ErrorMessage)
				}
				if err != nil {
					c.logger.Sugar().Error(err)
				}
			} else {
				if res.Time != nil {
					fmt.Println(*res.Time)
					res, err := c.BotClient.Api.Client.BanUser(&helix.BanUserParams{
						BroadcasterID: msg.RoomID,
						ModeratorId:   c.BotClient.Model.ID,
						Body: helix.BanUserRequestBody{
							Duration: int(*res.Time),
							Reason:   res.Message,
							UserId:   msg.User.ID,
						},
					})

					if res.StatusCode != 200 {
						c.logger.Sugar().Info(res.ErrorMessage)
					}
					if err != nil {
						c.logger.Sugar().Error(err)
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
		c.db.Save(&model.ChannelModerationWarn{
			ID:        uuid.NewV4().String(),
			ChannelID: settings.ChannelID,
			UserID:    userID,
			Reason:    reason,
		})
		return &handleResult{
			IsDelete: true,
			Message:  lo.If(settings.WarningMessage.Valid, settings.WarningMessage.String).Else(""),
		}
	} else {
		duration := time.Duration(settings.BanTime) * time.Second
		c.db.Delete(&warning)
		return &handleResult{
			IsDelete: false,
			Time:     lo.ToPtr(int(duration.Seconds())),
			Message:  lo.If(settings.BanMessage.Valid, settings.BanMessage.String).Else(""),
		}
	}
}

func (c *parsers) linksParser(
	settings *model.ChannelsModerationSettings,
	ircMsg irc.PrivateMessage,
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
	ircMsg irc.PrivateMessage,
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
	ircMsg irc.PrivateMessage,
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
	ircMsg irc.PrivateMessage,
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
	ircMsg irc.PrivateMessage,
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
	ircMsg irc.PrivateMessage,
	badges []string,
) *handleResult {
	if !settings.TriggerLength.Valid {
		return nil
	}

	if len(ircMsg.Emotes) < int(settings.TriggerLength.Int64) {
		return nil
	}

	return c.returnByWarnedState(settings.Type, ircMsg.User.ID, settings)
}
