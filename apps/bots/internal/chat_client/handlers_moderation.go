package chat_client

import (
	"log/slog"
	"regexp"

	"github.com/samber/lo"
	model "github.com/satont/twir/libs/gomodels"
)

type moderationService struct {
	*services

	linksRegexp           *regexp.Regexp
	linksWithSpacesRegexp *regexp.Regexp
}

type moderationHandleResult struct {
	IsDelete bool
	Time     *int
	Message  string
}

var moderationFunctionsMapping = map[string]func(
	c *moderationService,
	settings model.ChannelsModerationSettings,
	ircMsg Message,
	badges []string,
) *moderationHandleResult{
	"links":       (*moderationService).linksParser,
	"denylist":    (*moderationService).denyListParser,
	"symbols":     (*moderationService).symbolsParser,
	"longMessage": (*moderationService).longMessageParser,
	"caps":        (*moderationService).capsParser,
	"emotes":      (*moderationService).emotesParser,
}

func (c *ChatClient) handleModeration(msg Message) bool {
	userBadges := lo.Keys(msg.User.Badges)

	if lo.Some(userBadges, []string{"BROADCASTER", "MODERATOR"}) {
		return false
	}

	var settings []model.ChannelsModerationSettings
	if err := c.services.DB.Where(
		`"channelId" = ? AND "enabled" = ?`,
		msg.Channel.ID,
		true,
	).Find(&settings).Error; err != nil {
		c.services.Logger.Error("cannot find moderation settings", slog.Any("err", err))
		return false
	}

	for _, setting := range settings {
		if !setting.Enabled {
			continue
		}

		function, ok := moderationFunctionsMapping[setting.Type]
		if !ok {
			continue
		}

		res := function(c.moderationService, setting, msg, userBadges)

		if res == nil {
			continue
		}
	}

	return false
}

func (c *moderationService) linksParser(
	settings model.ChannelsModerationSettings,
	ircMsg Message,
	badges []string,
) *moderationHandleResult {
	return nil
}

func (c *moderationService) denyListParser(
	settings model.ChannelsModerationSettings,
	ircMsg Message,
	badges []string,
) *moderationHandleResult {
	return nil
}

func (c *moderationService) symbolsParser(
	settings model.ChannelsModerationSettings,
	ircMsg Message,
	badges []string,
) *moderationHandleResult {
	return nil
}

func (c *moderationService) longMessageParser(
	settings model.ChannelsModerationSettings,
	ircMsg Message,
	badges []string,
) *moderationHandleResult {
	return nil
}

func (c *moderationService) capsParser(
	settings model.ChannelsModerationSettings,
	ircMsg Message,
	badges []string,
) *moderationHandleResult {
	return nil
}

func (c *moderationService) emotesParser(
	settings model.ChannelsModerationSettings,
	ircMsg Message,
	badges []string,
) *moderationHandleResult {
	return nil
}
