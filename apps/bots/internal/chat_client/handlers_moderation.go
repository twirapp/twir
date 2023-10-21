package chat_client

import (
	"context"
	"fmt"
	"log/slog"
	"regexp"
	"slices"
	"strings"
	"time"

	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/bots/internal/moderation_helpers"
	model "github.com/satont/twir/libs/gomodels"
	language_detector "github.com/satont/twir/libs/grpc/generated/language-detector"
)

type moderationService struct {
	*services

	linksRegexp           *regexp.Regexp
	linksWithSpacesRegexp *regexp.Regexp
}

type moderationHandleResult struct {
	IsDelete bool
	Time     int
	Message  string
}

var moderationFunctionsMapping = map[model.ModerationSettingsType]func(
	c *moderationService,
	ctx context.Context,
	settings model.ChannelModerationSettings,
	ircMsg Message,
) *moderationHandleResult{
	model.ModerationSettingsTypeLinks:       (*moderationService).linksParser,
	model.ModerationSettingsTypeDenylist:    (*moderationService).denyListParser,
	model.ModerationSettingsTypeSymbols:     (*moderationService).symbolsParser,
	model.ModerationSettingsTypeLongMessage: (*moderationService).longMessageParser,
	model.ModerationSettingsTypeCaps:        (*moderationService).capsParser,
	model.ModerationSettingsTypeEmotes:      (*moderationService).emotesParser,
	model.ModerationSettingsTypeLanguage:    (*moderationService).languageParser,
}

func (c *moderationService) getChannelSettings(ctx context.Context, channelId string) (
	[]model.ChannelModerationSettings,
	error,
) {
	cacheKey := fmt.Sprintf("channels:%s:moderation_settings:*", channelId)

	cachedKeys, err := c.Redis.Keys(ctx, cacheKey).Result()
	if err != nil {
		return nil, err
	}

	var settings []model.ChannelModerationSettings

	if len(cachedKeys) > 0 {
		for _, key := range cachedKeys {
			var item model.ChannelModerationSettings
			if err := c.Redis.Get(ctx, key).Scan(&item); err != nil {
				return nil, err
			}

			settings = append(settings, item)
		}

		return settings, nil
	}

	if err := c.services.DB.Where(
		`"channelId" = ? AND "enabled" = ?`,
		channelId,
		true,
	).Find(&settings).Error; err != nil {
		c.services.Logger.Error("cannot find moderation settings", slog.Any("err", err))
		return nil, err
	}

	return settings, nil
}

func (c *ChatClient) handleModeration(ctx context.Context, msg Message) bool {
	userBadges := lo.Keys(msg.User.Badges)

	if lo.Some(userBadges, []string{"BROADCASTER", "MODERATOR"}) {
		return false
	}

	settings, err := c.moderationService.getChannelSettings(ctx, msg.Channel.ID)
	if err != nil {
		c.services.Logger.Error("cannot get channel settings", slog.Any("err", err))
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

		// TODO: check if user has role which must be excluded from moderation

		res := function(c.moderationService, ctx, setting, msg)

		if res == nil {
			continue
		}

		if res.IsDelete {
			res, err := c.services.TwitchClient.DeleteChatMessage(
				&helix.DeleteChatMessageParams{
					BroadcasterID: msg.Channel.ID,
					ModeratorID:   c.Model.ID,
					MessageID:     msg.ID,
				},
			)

			if res.StatusCode != 200 {
				c.services.Logger.Error("cannot delete message", slog.String("err", res.ErrorMessage))
			}
			if err != nil {
				c.services.Logger.Error("cannot delete message", slog.Any("err", err))
			}
		} else {
			if res.Time != 0 {
				res, err := c.services.TwitchClient.BanUser(
					&helix.BanUserParams{
						BroadcasterID: msg.Channel.ID,
						ModeratorId:   c.Model.ID,
						Body: helix.BanUserRequestBody{
							Duration: res.Time,
							Reason:   res.Message,
							UserId:   msg.User.ID,
						},
					},
				)

				if res.StatusCode != 200 {
					c.services.Logger.Error(
						"cannot ban user",
						slog.String("err", res.ErrorMessage),
						slog.String("userId", msg.User.ID),
						slog.String("channelId", msg.Channel.ID),
					)
				}
				if err != nil {
					c.services.Logger.Error(
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

	return false
}

func (c *moderationService) returnByWarnedState(
	ctx context.Context,
	userID string,
	settings model.ChannelModerationSettings,
) *moderationHandleResult {
	warningRedisKey := fmt.Sprintf(
		"channels:%s:moderation_warns:%s:%s:*", settings.ChannelID,
		userID, settings.Type,
	)
	warningsKeys, err := c.services.Redis.Keys(ctx, warningRedisKey).Result()
	if err != nil {
		c.services.Logger.Error("cannot get warnings", slog.Any("err", err))
		return nil
	}

	if settings.MaxWarnings > 0 && settings.MaxWarnings > len(warningsKeys) {
		c.Redis.Set(
			ctx,
			fmt.Sprintf(
				"channels:%s:moderation_warns:%s:%s:%v",
				settings.ChannelID,
				userID,
				settings.Type,
				time.Now().Unix(),
			),
			"",
			24*time.Hour,
		)

		return &moderationHandleResult{
			IsDelete: true,
			Message:  settings.WarningMessage,
		}
	} else {
		duration := time.Duration(settings.BanTime) * time.Second

		for _, key := range warningsKeys {
			c.Redis.Del(ctx, key)
		}

		return &moderationHandleResult{
			IsDelete: false,
			Time:     int(duration.Seconds()),
			Message:  settings.BanMessage,
		}
	}
}

func (c *moderationService) linksParser(
	ctx context.Context,
	settings model.ChannelModerationSettings,
	ircMsg Message,
) *moderationHandleResult {
	containLink := moderation_helpers.HasLink(c.linksWithSpacesRegexp, ircMsg.Message)

	if !containLink {
		return nil
	}

	permit := model.ChannelsPermits{}
	err := c.services.DB.Where(
		`"channelId" = ? AND "userId" = ?`,
		settings.ChannelID,
		ircMsg.User.ID,
	).
		Find(&permit).
		Error
	if err != nil {
		return nil
	}

	if permit.ID != "" {
		c.services.DB.Delete(&permit)
		return nil
	} else {
		return c.returnByWarnedState(ctx, ircMsg.User.ID, settings)
	}
}

func (c *moderationService) denyListParser(
	ctx context.Context,
	settings model.ChannelModerationSettings,
	ircMsg Message,
) *moderationHandleResult {
	hasDeniedWord := moderation_helpers.HasDeniedWord(ircMsg.Message, settings.DenyList)
	if !hasDeniedWord {
		return nil
	}

	return c.returnByWarnedState(ctx, ircMsg.User.ID, settings)
}

func (c *moderationService) symbolsParser(
	ctx context.Context,
	settings model.ChannelModerationSettings,
	ircMsg Message,
) *moderationHandleResult {
	if len(ircMsg.Message) < settings.TriggerLength {
		return nil
	}

	isToMuchSymbols, _ := moderation_helpers.IsToMuchSymbols(ircMsg.Message, settings.MaxPercentage)
	if !isToMuchSymbols {
		return nil
	}

	return c.returnByWarnedState(ctx, ircMsg.User.ID, settings)
}

func (c *moderationService) longMessageParser(
	ctx context.Context,
	settings model.ChannelModerationSettings,
	ircMsg Message,
) *moderationHandleResult {
	isToLong := moderation_helpers.IsTooLong(ircMsg.Message, settings.TriggerLength)

	if !isToLong {
		return nil
	}

	return c.returnByWarnedState(ctx, ircMsg.User.ID, settings)
}

func (c *moderationService) capsParser(
	ctx context.Context,
	settings model.ChannelModerationSettings,
	ircMsg Message,
) *moderationHandleResult {
	msg := ircMsg.Message
	for _, v := range ircMsg.Emotes {
		msg = strings.ReplaceAll(msg, v.Name, "")
	}

	isTooLong, _ := moderation_helpers.IsTooMuchCaps(msg, settings.MaxPercentage)

	if !isTooLong {
		return nil
	}

	return c.returnByWarnedState(ctx, ircMsg.User.ID, settings)
}

func (c *moderationService) emotesParser(
	ctx context.Context,
	settings model.ChannelModerationSettings,
	ircMsg Message,
) *moderationHandleResult {
	if settings.TriggerLength == 0 {
		return nil
	}

	length := 0

	for _, e := range ircMsg.Emotes {
		length += len(e.Positions)
	}

	if length < settings.TriggerLength {
		return nil
	}

	channelEmotesKeys, err := c.Redis.Keys(
		ctx,
		fmt.Sprintf("emotes:channel:%s:*", settings.ChannelID),
	).Result()
	if err != nil {
		c.services.Logger.Error("cannot get channel emotes", slog.Any("err", err))
		return nil
	}
	for _, key := range channelEmotesKeys {
		key = strings.Replace(key, fmt.Sprintf("emotes:channel:%s:", settings.ChannelID), "", 1)
	}

	splittedMsg := strings.Split(ircMsg.Message, " ")

	for _, word := range splittedMsg {
		if slices.Contains(channelEmotesKeys, word) {
			length++
		}
	}

	globalEmotesKeys, err := c.Redis.Keys(
		ctx,
		fmt.Sprintf("emotes:global:*"),
	).Result()
	if err != nil {
		c.services.Logger.Error("cannot get global emotes", slog.Any("err", err))
		return nil
	}
	for _, key := range globalEmotesKeys {
		key = strings.Replace(key, fmt.Sprintf("emotes:global:"), "", 1)
	}

	for _, word := range splittedMsg {
		if slices.Contains(globalEmotesKeys, word) {
			length++
		}
	}

	return c.returnByWarnedState(ctx, ircMsg.User.ID, settings)
}

func (c *moderationService) languageParser(
	ctx context.Context,
	settings model.ChannelModerationSettings,
	ircMsg Message,
) *moderationHandleResult {
	detected, err := c.LanguageDetector.Detect(
		ctx, &language_detector.Request{
			Text: ircMsg.Message,
		},
	)

	if err != nil {
		c.services.Logger.Error("cannot detect language", slog.Any("err", err))
		return nil
	}

	if len(detected.Languages) == 0 {
		return nil
	}

	hasDeniedLanguage := lo.SomeBy(
		detected.Languages, func(item *language_detector.Response_Language) bool {
			return !slices.Contains(settings.AcceptedChatLanguages, item.Name)
		},
	)

	if !hasDeniedLanguage {
		return nil
	}

	return c.returnByWarnedState(ctx, ircMsg.User.ID, settings)
}
