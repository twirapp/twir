package messagehandler

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/apps/bots/internal/moderationhelpers"
	"github.com/twirapp/twir/apps/bots/internal/twitchactions"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/logger"
	channelsmoderationsettingsmodel "github.com/twirapp/twir/libs/repositories/channels_moderation_settings/model"
	"github.com/twirapp/twir/libs/utils"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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
	msg twitch.TwitchChatMessage,
) *moderationHandleResult{
	channelsmoderationsettingsmodel.ModerationSettingsTypeLinks:       (*MessageHandler).moderationLinksParser,
	channelsmoderationsettingsmodel.ModerationSettingsTypeDenylist:    (*MessageHandler).moderationDenyListParser,
	channelsmoderationsettingsmodel.ModerationSettingsTypeSymbols:     (*MessageHandler).moderationSymbolsParser,
	channelsmoderationsettingsmodel.ModerationSettingsTypeLongMessage: (*MessageHandler).moderationLongMessageParser,
	channelsmoderationsettingsmodel.ModerationSettingsTypeCaps:        (*MessageHandler).moderationCapsParser,
	channelsmoderationsettingsmodel.ModerationSettingsTypeEmotes:      (*MessageHandler).moderationEmotesParser,
	channelsmoderationsettingsmodel.ModerationSettingsTypeLanguage:    (*MessageHandler).moderationLanguageParser,
	channelsmoderationsettingsmodel.ModerationSettingsTypeOneManSpam:  (*MessageHandler).moderationOneManSpam,
}

func (c *MessageHandler) handleModeration(ctx context.Context, msg twitch.TwitchChatMessage) error {
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(attribute.String("function.name", utils.GetFuncName()))

	if msg.IsChatterBroadcaster() || msg.IsChatterModerator() {
		return nil
	}

	settings, err := c.getChannelModerationSettings(ctx, msg.BroadcasterUserId)
	if err != nil {
		return fmt.Errorf("cannot get moderation settings: %w", err)
	}

	for _, entity := range settings {
		if !entity.Enabled {
			continue
		}

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
					logger.Error(err),
				)
			}

			var warningMessage = strings.TrimSpace(entity.WarningMessage)

			if warningMessage == "" && entity.Name != nil && *entity.Name != "" {
				warningMessage = fmt.Sprintf(
					"Reason: %s",
					*entity.Name,
				)
			} else if warningMessage == "" {
				warningMessage = fmt.Sprintf(
					"Reason: %s",
					strings.Join(strings.Split(entity.Type.String(), "_"), " "),
				)
			}

			err = c.twitchActions.WarnUser(
				ctx, twitchactions.WarnUserOpts{
					BroadcasterID: msg.BroadcasterUserId,
					ModeratorID:   msg.EnrichedData.DbChannel.BotID,
					UserID:        msg.ChatterUserId,
					Reason:        warningMessage,
				},
			)
			if err != nil {
				c.logger.Error(
					"cannot warn user",
					slog.String("userId", msg.ChatterUserId),
					slog.String("channelId", msg.BroadcasterUserId),
					logger.Error(err),
				)
			}
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
					logger.Error(err),
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
	msg twitch.TwitchChatMessage,
	settings channelsmoderationsettingsmodel.ChannelModerationSettings,
) *moderationHandleResult {
	var channelRoles []model.ChannelRole
	if err := c.gorm.WithContext(ctx).Preload("Users", `"userId" = ?`, msg.ChatterUserId).Where(
		`"channelId" = ?`,
		settings.ChannelID,
	).
		Find(&channelRoles).
		Error; err != nil {
		c.logger.Error("cannot get channel roles", logger.Error(err))
		return nil
	}

	for _, r := range channelRoles {
		if r.Type == model.ChannelRoleTypeCustom {
			continue
		}

		shouldExcludeRole := slices.Contains(settings.ExcludedRoles, r.ID)
		var userHasRole bool
		if len(r.Users) > 0 {
			userHasRole = true
		}

		if msg.HasRoleFromDbByType(r.Type.String()) {
			userHasRole = true
		}

		if msg.EnrichedData.DbUserChannelStat != nil && !userHasRole {
			if r.RequiredWatchTime > 0 && msg.EnrichedData.DbUserChannelStat.Watched >= r.RequiredWatchTime {
				userHasRole = true
			}

			if r.RequiredMessages > 0 && msg.EnrichedData.DbUserChannelStat.Messages >= r.RequiredMessages {
				userHasRole = true
			}

			if r.RequiredUsedChannelPoints > 0 && msg.EnrichedData.DbUserChannelStat.UsedChannelPoints >= r.RequiredUsedChannelPoints {
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
	msg twitch.TwitchChatMessage,
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
	msg twitch.TwitchChatMessage,
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
	msg twitch.TwitchChatMessage,
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
	msg twitch.TwitchChatMessage,
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
	msg twitch.TwitchChatMessage,
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
	msg twitch.TwitchChatMessage,
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
		reqUrl = "http://language-processor:3012/detect"
	} else {
		reqUrl = "http://localhost:3012/detect"
	}

	u, err := url.Parse(reqUrl)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("text", text)
	u.RawQuery = q.Encode()

	httpReq, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, err
	}

	res, err := http.DefaultClient.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	if res.StatusCode < 200 || res.StatusCode >= 300 {
		return nil, errors.New("cannot get response")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	resp := langDetectResult{}
	if err := json.Unmarshal(bodyBytes, &resp); err != nil {
		return nil, err
	}

	return &resp, nil
}

func (c *MessageHandler) moderationLanguageParser(
	ctx context.Context,
	settings channelsmoderationsettingsmodel.ChannelModerationSettings,
	msg twitch.TwitchChatMessage,
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
	text = strings.ToLower(text)

	if utf8.RuneCountInString(text) < settings.TriggerLength {
		return nil
	}

	for _, word := range settings.LanguageExcludedWords {
		text = strings.ReplaceAll(text, strings.ToLower(word), "")
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

func (c *MessageHandler) moderationOneManSpam(
	ctx context.Context,
	settings channelsmoderationsettingsmodel.ChannelModerationSettings,
	msg twitch.TwitchChatMessage,
) *moderationHandleResult {
	if len(msg.Message.Text) < settings.TriggerLength {
		return nil
	}

	if settings.OneManSpamMessageMemorySeconds == 0 || settings.OneManSpamMinimumStoredMessages == 0 {
		return nil
	}

	redisKey := fmt.Sprintf(
		"channels:%s:moderation:one_man_spam:%s",
		msg.BroadcasterUserId,
		msg.ChatterUserId,
	)

	defer func() {
		deferCtx, deferCtxCancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer deferCtxCancel()

		if err := c.redis.HSet(
			deferCtx,
			redisKey,
			msg.ID,
			msg.Message.Text,
		).Err(); err != nil {
			c.logger.Error("cannot set one man spam to redis", logger.Error(err))
			return
		}

		if err := c.redis.HExpire(
			deferCtx,
			redisKey,
			time.Duration(settings.OneManSpamMessageMemorySeconds)*time.Second,
			msg.ID,
		).Err(); err != nil {
			c.logger.Error("cannot expire one man spam redis key", logger.Error(err))
			return
		}
	}()

	storedData, err := c.redis.HGetAll(ctx, redisKey).Result()
	if err != nil {
		return nil
	}

	messages := make([]string, 0, len(storedData))
	for _, value := range storedData {
		messages = append(messages, value)
	}

	if len(messages) < settings.OneManSpamMinimumStoredMessages {
		return nil
	}

	if c.moderationHelpers.ContainsSimilar(
		msg.Message.Text,
		messages,
		float64(settings.MaxPercentage),
	) {
		return c.moderationHandleResult(ctx, msg, settings)
	}

	return nil
}
