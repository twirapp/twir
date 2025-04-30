package messagehandler

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/imroc/req/v3"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/bots/internal/moderationhelpers"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	model "github.com/satont/twir/libs/gomodels"
	channelsmoderationsettingsmodel "github.com/twirapp/twir/libs/repositories/channels_moderation_settings/model"
)

type moderationHandleResult struct {
	Message string
	Time    int
	IsWarn  bool
}

var moderationFunctionsMapping = map[channelsmoderationsettingsmodel.ModerationSettingsType]func(
	c *MessageHandler,
	ctx context.Context,
	settings channelsmoderationsettingsmodel.ChannelModerationSettings,
	msg handleMessage,
) *moderationHandleResult{
	channelsmoderationsettingsmodel.ModerationSettingsTypeLinks:       (*MessageHandler).moderationLinksParser,
	channelsmoderationsettingsmodel.ModerationSettingsTypeDenylist:    (*MessageHandler).moderationDenyListParser,
	channelsmoderationsettingsmodel.ModerationSettingsTypeSymbols:     (*MessageHandler).moderationSymbolsParser,
	channelsmoderationsettingsmodel.ModerationSettingsTypeLongMessage: (*MessageHandler).moderationLongMessageParser,
	channelsmoderationsettingsmodel.ModerationSettingsTypeCaps:        (*MessageHandler).moderationCapsParser,
	channelsmoderationsettingsmodel.ModerationSettingsTypeEmotes:      (*MessageHandler).moderationEmotesParser,
	channelsmoderationsettingsmodel.ModerationSettingsTypeLanguage:    (*MessageHandler).moderationLanguageParser,
}

var excludedModerationBadges = []string{"BROADCASTER", "MODERATOR"}

func (c *MessageHandler) handleModeration(ctx context.Context, msg handleMessage) error {
	badges := createUserBadges(msg.Badges)

	for _, b := range badges {
		if slices.Contains(excludedModerationBadges, b) {
			return nil
		}
	}

	settings, err := c.getChannelModerationSettings(ctx, msg.BroadcasterUserId)
	if err != nil {
		return fmt.Errorf("cannot get moderation settings: %w", err)
	}

	for _, entity := range settings {
		function, ok := moderationFunctionsMapping[entity.Type]
		if !ok {
			continue
		}

		res := function(c, ctx, entity, msg)
		if res == nil {
			continue
		}

		if res.IsWarn {
			err := c.twitchActions.DeleteMessage(
				ctx,
				twitchactions.DeleteMessageOpts{
					BroadcasterID: msg.BroadcasterUserId,
					ModeratorID:   msg.EnrichedData.DbChannel.BotID,
					MessageID:     msg.MessageId,
				},
			)
			if err != nil {
				c.logger.Error(
					"cannot delete message",
					slog.String("userId", msg.ChatterUserId),
					slog.String("channelId", msg.BroadcasterUserId),
					slog.Any("err", err),
				)
			}

			// TODO: uncomment
			// err = c.twitchActions.WarnUser(
			// 	ctx, twitchactions.WarnUserOpts{
			// 		BroadcasterID: msg.BroadcasterUserId,
			// 		ModeratorID:   msg.EnrichedData.DbChannel.BotID,
			// 		UserID:        msg.ChatterUserId,
			// 		Reason:        entity.WarningMessage,
			// 	},
			// )
			// if err != nil {
			// 	c.logger.Error(
			// 		"cannot warn user",
			// 		slog.String("userId", msg.ChatterUserId),
			// 		slog.String("channelId", msg.BroadcasterUserId),
			// 		slog.Any("err", err),
			// 	)
			// }
		} else {
			err := c.twitchActions.Ban(
				ctx, twitchactions.BanOpts{
					Duration:      res.Time,
					Reason:        entity.BanMessage,
					BroadcasterID: msg.BroadcasterUserId,
					UserID:        msg.ChatterUserId,
					ModeratorID:   msg.EnrichedData.DbChannel.BotID,
				},
			)

			if err != nil {
				c.logger.Error(
					"cannot ban user",
					slog.String("userId", msg.ChatterUserId),
					slog.String("channelId", msg.BroadcasterUserId),
					slog.Any("err", err),
				)
			}
		}

		return nil
	}

	return nil
}

func (c *MessageHandler) getChannelModerationSettings(ctx context.Context, channelId string) (
	[]channelsmoderationsettingsmodel.ChannelModerationSettings,
	error,
) {
	settings, err := c.channelsModerationSettingsCacher.Get(ctx, channelId)
	if err != nil {
		return nil, err
	}

	return settings, nil
}

func (c *MessageHandler) moderationHandleResult(
	ctx context.Context,
	msg handleMessage,
	settings channelsmoderationsettingsmodel.ChannelModerationSettings,
) *moderationHandleResult {
	var channelRoles []model.ChannelRole
	if err := c.gorm.WithContext(ctx).Preload("Users", `"userId" = ?`, msg.ChatterUserId).Where(
		`"channelId" = ?`,
		settings.ChannelID,
	).
		Find(&channelRoles).
		Error; err != nil {
		c.logger.Error("cannot get channel roles", slog.Any("err", err))
		return nil
	}
	badges := createUserBadges(msg.Badges)

	for _, r := range channelRoles {
		if r.Type == model.ChannelRoleTypeCustom {
			continue
		}

		shouldExcludeRole := slices.Contains(settings.ExcludedRoles, r.ID)
		var userHasRole bool
		if len(r.Users) > 0 {
			userHasRole = true
		}

		if slices.Contains(badges, r.Type.String()) {
			userHasRole = true
		}

		if msg.DbUser.Stats != nil && !userHasRole {
			if r.RequiredWatchTime > 0 && msg.DbUser.Stats.Watched >= r.RequiredWatchTime {
				userHasRole = true
			}

			if r.RequiredMessages > 0 && msg.DbUser.Stats.Messages >= r.RequiredMessages {
				userHasRole = true
			}

			if r.RequiredUsedChannelPoints > 0 && msg.DbUser.Stats.UsedChannelPoints >= r.RequiredUsedChannelPoints {
				userHasRole = true
			}
		}

		if shouldExcludeRole && userHasRole {
			return nil
		}
	}

	warningRedisKey := fmt.Sprintf(
		"channels:%s:moderation_warns:%s:%s",
		settings.ChannelID,
		msg.ChatterUserId,
		settings.Type,
	)

	warningCount, err := c.redis.Get(ctx, warningRedisKey).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil
	}

	if settings.MaxWarnings > 0 && settings.MaxWarnings > warningCount {
		c.redis.Pipelined(
			ctx, func(pipe redis.Pipeliner) error {
				if pErr := pipe.Incr(ctx, warningRedisKey).Err(); pErr != nil {
					return pErr
				}

				if pErr := pipe.Persist(ctx, warningRedisKey).Err(); pErr != nil {
					return pErr
				}

				return nil
			},
		)

		return &moderationHandleResult{
			IsWarn:  true,
			Message: settings.WarningMessage,
		}
	} else {
		duration := time.Duration(settings.BanTime) * time.Second
		c.redis.Del(ctx, warningRedisKey)
		return &moderationHandleResult{
			IsWarn:  false,
			Time:    int(duration.Seconds()),
			Message: settings.BanMessage,
		}
	}
}

func (c *MessageHandler) moderationLinksParser(
	ctx context.Context,
	settings channelsmoderationsettingsmodel.ChannelModerationSettings,
	msg handleMessage,
) *moderationHandleResult {
	containLink := c.moderationHelpers.HasLink(msg.Message.Text, settings.CheckClips)

	if !containLink {
		return nil
	}

	permit := model.ChannelsPermits{}
	err := c.gorm.WithContext(ctx).Where(
		`"channelId" = ? AND "userId" = ?`,
		settings.ChannelID,
		msg.ChatterUserId,
	).
		Find(&permit).
		Error
	if err != nil {
		return nil
	}

	if permit.ID != "" {
		c.gorm.WithContext(ctx).Delete(&permit)
		return nil
	} else {
		return c.moderationHandleResult(ctx, msg, settings)
	}
}

func (c *MessageHandler) moderationDenyListParser(
	ctx context.Context,
	settings channelsmoderationsettingsmodel.ChannelModerationSettings,
	msg handleMessage,
) *moderationHandleResult {
	if len(settings.DenyList) == 0 {
		return nil
	}

	hasDeniedWord := c.moderationHelpers.HasDeniedWord(
		moderationhelpers.HasDeniedWordInput{
			Message:             msg.Message.Text,
			RulesList:           settings.DenyList,
			RegexpEnabled:       settings.DenyListRegexpEnabled,
			WordBoundaryEnabled: settings.DenyListWordBoundaryEnabled,
			SensitivityEnabled:  settings.DenyListSensitivityEnabled,
		},
	)
	if !hasDeniedWord {
		return nil
	}

	return c.moderationHandleResult(ctx, msg, settings)
}

func (c *MessageHandler) moderationSymbolsParser(
	ctx context.Context,
	settings channelsmoderationsettingsmodel.ChannelModerationSettings,
	msg handleMessage,
) *moderationHandleResult {
	if utf8.RuneCountInString(msg.Message.Text) < settings.TriggerLength {
		return nil
	}

	isToMuchSymbols, _ := c.moderationHelpers.IsToMuchSymbols(
		msg.Message.Text,
		settings.MaxPercentage+1,
	)
	if !isToMuchSymbols {
		return nil
	}

	return c.moderationHandleResult(ctx, msg, settings)
}

func (c *MessageHandler) moderationLongMessageParser(
	ctx context.Context,
	settings channelsmoderationsettingsmodel.ChannelModerationSettings,
	msg handleMessage,
) *moderationHandleResult {
	isToLong := c.moderationHelpers.IsTooLong(msg.Message.Text, settings.TriggerLength)

	if !isToLong {
		return nil
	}

	return c.moderationHandleResult(ctx, msg, settings)
}

func (c *MessageHandler) moderationCapsParser(
	ctx context.Context,
	settings channelsmoderationsettingsmodel.ChannelModerationSettings,
	msg handleMessage,
) *moderationHandleResult {
	text := msg.Message.Text

	for emote, _ := range msg.EnrichedData.UsedEmotesWithThirdParty {
		text = strings.ReplaceAll(text, emote, "")
	}

	if utf8.RuneCountInString(text) < settings.TriggerLength {
		return nil
	}

	isTooLong, _ := c.moderationHelpers.IsTooMuchCaps(text, settings.MaxPercentage)

	if !isTooLong {
		return nil
	}

	return c.moderationHandleResult(ctx, msg, settings)
}

func (c *MessageHandler) moderationEmotesParser(
	ctx context.Context,
	settings channelsmoderationsettingsmodel.ChannelModerationSettings,
	msg handleMessage,
) *moderationHandleResult {
	if settings.TriggerLength == 0 {
		return nil
	}

	var totalEmotesInMessage int
	for _, count := range msg.EnrichedData.UsedEmotesWithThirdParty {
		totalEmotesInMessage += count
	}

	if totalEmotesInMessage < settings.TriggerLength+1 {
		return nil
	}

	return c.moderationHandleResult(ctx, msg, settings)
}

type detectedLang struct {
	Language    string  `json:"language"`
	Probability float64 `json:"probability"`
}

type langDetectResult struct {
	Text              string         `json:"text"`
	CleanedText       string         `json:"cleaned_text"`
	DetectedLanguages []detectedLang `json:"detected_languages"`
}

func (c *MessageHandler) moderationDetectLanguage(text string) (*langDetectResult, error) {
	var reqUrl string
	if c.config.AppEnv == "production" {
		reqUrl = fmt.Sprint("http://language-processor:3012/detect")
	} else {
		reqUrl = "http://localhost:3012/detect"
	}

	resp := langDetectResult{}
	res, err := req.R().SetQueryParam("text", text).SetSuccessResult(&resp).Get(reqUrl)
	if err != nil {
		return nil, err
	}
	if !res.IsSuccessState() {
		return nil, errors.New("cannot get response")
	}

	return &resp, nil
}

func (c *MessageHandler) moderationLanguageParser(
	ctx context.Context,
	settings channelsmoderationsettingsmodel.ChannelModerationSettings,
	msg handleMessage,
) *moderationHandleResult {
	text := msg.Message.Text
	for emote, _ := range msg.EnrichedData.UsedEmotesWithThirdParty {
		text = strings.ReplaceAll(text, emote, "")
	}

	for _, fragment := range msg.Message.Fragments {
		if fragment.Mention != nil {
			text = strings.ReplaceAll(text, fragment.Text, "")
		}
	}

	text = strings.TrimSpace(text)

	if utf8.RuneCountInString(text) < 10 {
		return nil
	}

	detected, err := c.moderationDetectLanguage(text)
	if err != nil || detected == nil || len(detected.DetectedLanguages) == 0 {
		return nil
	}

	bestDetected := detected.DetectedLanguages[0]
	if !slices.Contains(settings.DeniedChatLanguages, bestDetected.Language) {
		return nil
	}

	return c.moderationHandleResult(ctx, msg, settings)
}
