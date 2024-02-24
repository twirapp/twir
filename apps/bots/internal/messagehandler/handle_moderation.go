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

	"github.com/goccy/go-json"
	"github.com/imroc/req/v3"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/bots/internal/twitchactions"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/utils"
	"github.com/twirapp/twir/libs/bus-core/twitch"
)

type moderationHandleResult struct {
	IsDelete bool
	Time     int
	Message  string
}

var messagesTimeouterStore = utils.NewTtlSyncMap[struct{}](10 * time.Second)
var moderationFunctionsMapping = map[model.ModerationSettingsType]func(
	c *MessageHandler,
	ctx context.Context,
	settings model.ChannelModerationSettings,
	msg handleMessage,
) *moderationHandleResult{
	model.ModerationSettingsTypeLinks:       (*MessageHandler).moderationLinksParser,
	model.ModerationSettingsTypeDenylist:    (*MessageHandler).moderationDenyListParser,
	model.ModerationSettingsTypeSymbols:     (*MessageHandler).moderationSymbolsParser,
	model.ModerationSettingsTypeLongMessage: (*MessageHandler).moderationLongMessageParser,
	model.ModerationSettingsTypeCaps:        (*MessageHandler).moderationCapsParser,
	model.ModerationSettingsTypeEmotes:      (*MessageHandler).moderationEmotesParser,
	model.ModerationSettingsTypeLanguage:    (*MessageHandler).moderationLanguageParser,
}

func (c *MessageHandler) handleModeration(ctx context.Context, msg handleMessage) error {
	badges := createUserBadges(msg.Badges)

	if lo.Some(badges, []string{"broadcaster", "moderator"}) {
		return nil
	}

	settings, err := c.getChannelModerationSettings(ctx, msg.BroadcasterUserId)
	if err != nil {
		return err
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

		if _, exists := messagesTimeouterStore.Get(msg.BroadcasterUserId); !exists {
			opts := twitchactions.SendMessageOpts{
				Message:       entity.BanMessage,
				SenderID:      msg.DbChannel.BotID,
				BroadcasterID: msg.BroadcasterUserId,
			}
			if res.IsDelete {
				opts.Message = entity.WarningMessage
			}
			if opts.Message != "" {
				c.twitchActions.SendMessage(ctx, opts)
			}
			messagesTimeouterStore.Add(msg.BroadcasterUserId, struct{}{})
		}

		if res.IsDelete {
			err := c.twitchActions.DeleteMessage(
				ctx,
				twitchactions.DeleteMessageOpts{
					BroadcasterID: msg.BroadcasterUserId,
					ModeratorID:   msg.DbChannel.BotID,
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
		} else {
			err := c.twitchActions.Ban(
				ctx, twitchactions.BanOpts{
					Duration:      res.Time,
					Reason:        entity.BanMessage,
					BroadcasterID: msg.BroadcasterUserId,
					UserID:        msg.ChatterUserId,
					ModeratorID:   msg.DbChannel.BotID,
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
	[]model.ChannelModerationSettings,
	error,
) {
	cacheKey := fmt.Sprintf("channels:%s:moderation_settings", channelId)

	cached, err := c.redis.Get(ctx, cacheKey).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}

	var settings []model.ChannelModerationSettings

	if len(cached) > 0 {
		if err := json.Unmarshal(cached, &settings); err != nil {
			return nil, err
		}

		return settings, nil
	}

	if err := c.gorm.
		WithContext(ctx).
		Where(
			`"channel_id" = ? AND "enabled" = ?`,
			channelId,
			true,
		).Find(&settings).Error; err != nil {
		return nil, err
	}

	c.redis.Set(ctx, cacheKey, settings, 24*time.Hour)

	return settings, nil
}

func (c *MessageHandler) moderationHandleResult(
	ctx context.Context,
	msg handleMessage,
	settings model.ChannelModerationSettings,
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
			roleSettings := model.ChannelRoleSettings{}
			if err := json.Unmarshal(r.Settings, &roleSettings); err == nil {
				if msg.DbUser.Stats.Watched >= roleSettings.RequiredWatchTime {
					userHasRole = true
				}

				if msg.DbUser.Stats.Messages >= roleSettings.RequiredMessages {
					userHasRole = true
				}

				if msg.DbUser.Stats.UsedChannelPoints >= roleSettings.RequiredUsedChannelPoints {
					userHasRole = true
				}
			}
		}

		if shouldExcludeRole && userHasRole {
			return nil
		}
	}

	warningRedisKey := fmt.Sprintf(
		"channels:%s:moderation_warns:%s:%s:*", settings.ChannelID,
		msg.ChatterUserId, settings.Type,
	)
	warningsKeys, err := c.redis.Keys(ctx, warningRedisKey).Result()
	if err != nil {
		c.logger.Error("cannot get warnings", slog.Any("err", err))
		return nil
	}

	if settings.MaxWarnings > 0 && settings.MaxWarnings > len(warningsKeys) {
		c.redis.Set(
			ctx,
			fmt.Sprintf(
				"channels:%s:moderation_warns:%s:%s:%v",
				settings.ChannelID,
				msg.ChatterUserId,
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
			c.redis.Del(ctx, key)
		}

		return &moderationHandleResult{
			IsDelete: false,
			Time:     int(duration.Seconds()),
			Message:  settings.BanMessage,
		}
	}
}

func (c *MessageHandler) moderationLinksParser(
	ctx context.Context,
	settings model.ChannelModerationSettings,
	msg handleMessage,
) *moderationHandleResult {
	containLink := c.moderationHelpers.HasLink(msg.Message.Text)

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
	settings model.ChannelModerationSettings,
	msg handleMessage,
) *moderationHandleResult {
	if len(settings.DenyList) == 0 {
		return nil
	}

	hasDeniedWord := c.moderationHelpers.HasDeniedWord(msg.Message.Text, settings.DenyList)
	if !hasDeniedWord {
		return nil
	}

	return c.moderationHandleResult(ctx, msg, settings)
}

func (c *MessageHandler) moderationSymbolsParser(
	ctx context.Context,
	settings model.ChannelModerationSettings,
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
	settings model.ChannelModerationSettings,
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
	settings model.ChannelModerationSettings,
	msg handleMessage,
) *moderationHandleResult {
	text := msg.Message.Text

	for _, f := range msg.Message.Fragments {
		if f.Type != twitch.FragmentType_EMOTE && f.Type != twitch.FragmentType_CHEERMOTE {
			continue
		}

		text = strings.ReplaceAll(text, f.Text, "")
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
	settings model.ChannelModerationSettings,
	msg handleMessage,
) *moderationHandleResult {
	if settings.TriggerLength == 0 {
		return nil
	}

	length := 0

	for _, f := range msg.Message.Fragments {
		if f.Type != twitch.FragmentType_EMOTE {
			continue
		}

		length += 1
	}

	if length < settings.TriggerLength+1 {
		return nil
	}

	channelEmotesKeys, err := c.redis.Keys(
		ctx,
		fmt.Sprintf("emotes:channel:%s:*", settings.ChannelID),
	).Result()
	if err != nil {
		c.logger.Error("cannot get channel emotes", slog.Any("err", err))
		return nil
	}
	for _, key := range channelEmotesKeys {
		key = strings.Replace(key, fmt.Sprintf("emotes:channel:%s:", settings.ChannelID), "", 1)
	}

	splittedMsg := strings.Split(msg.Message.Text, " ")

	for _, word := range splittedMsg {
		if slices.Contains(channelEmotesKeys, word) {
			length++
		}
	}

	globalEmotesKeys, err := c.redis.Keys(
		ctx,
		fmt.Sprintf("emotes:global:*"),
	).Result()
	if err != nil {
		c.logger.Error("cannot get global emotes", slog.Any("err", err))
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

	return c.moderationHandleResult(ctx, msg, settings)
}

type langDetectLang struct {
	Code    int    `json:"code"`
	Iso6933 int    `json:"iso_693_3"`
	Name    string `json:"name"`
}

func (c *MessageHandler) moderationDetectLanguage(text string) ([]langDetectLang, error) {
	var reqUrl string
	if c.config.AppEnv == "production" {
		reqUrl = fmt.Sprint("http://language-detector:3012")
	} else {
		reqUrl = "http://localhost:3012"
	}

	var resp []langDetectLang
	res, err := req.R().SetQueryParam("text", text).SetSuccessResult(&resp).Get(reqUrl)
	if err != nil {
		return nil, err
	}
	if !res.IsSuccessState() {
		return nil, errors.New("cannot get response")
	}

	return resp, nil
}

func (c *MessageHandler) moderationLanguageParser(
	ctx context.Context,
	settings model.ChannelModerationSettings,
	msg handleMessage,
) *moderationHandleResult {

	detected, err := c.moderationDetectLanguage(msg.Message.Text)
	if err != nil || len(detected) == 0 {
		return nil
	}

	hasDeniedLanguage := lo.SomeBy(
		detected,
		func(item langDetectLang) bool {
			return slices.Contains(settings.DeniedChatLanguages, string(item.Code))
		},
	)

	if !hasDeniedLanguage {
		return nil
	}

	return c.moderationHandleResult(ctx, msg, settings)
}
