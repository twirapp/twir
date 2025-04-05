package handle_message

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	redislimiter "github.com/aidenwallis/go-ratelimiting/redis"
	redislimiteradapter "github.com/aidenwallis/go-ratelimiting/redis/adapters/go-redis"
	"github.com/redis/go-redis/v9"
	config "github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/logger"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelschattrenslationsrepository "github.com/twirapp/twir/libs/repositories/chat_translation"
	"github.com/twirapp/twir/libs/repositories/chat_translation/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Config  config.Config
	Logger  logger.Logger
	TwirBus *buscore.Bus
	Redis   *redis.Client

	ChannelsRepository             channelsrepository.Repository
	ChannelsTranslationsRepository channelschattrenslationsrepository.Repository
	ChannelsTranslationsCache      *generic_cacher.GenericCacher[model.ChatTranslation]
}

func New(opts Opts) *Service {
	return &Service{
		config:                         opts.Config,
		logger:                         opts.Logger,
		twirBus:                        opts.TwirBus,
		redis:                          opts.Redis,
		channelsRepository:             opts.ChannelsRepository,
		channelsTranslationsRepository: opts.ChannelsTranslationsRepository,
		channelsTranslationsCache:      opts.ChannelsTranslationsCache,
		rateLimiter:                    redislimiter.NewSlidingWindow(redislimiteradapter.NewAdapter(opts.Redis)),
	}
}

type Service struct {
	config                         config.Config
	logger                         logger.Logger
	twirBus                        *buscore.Bus
	redis                          *redis.Client
	channelsRepository             channelsrepository.Repository
	channelsTranslationsRepository channelschattrenslationsrepository.Repository
	channelsTranslationsCache      *generic_cacher.GenericCacher[model.ChatTranslation]

	rateLimiter redislimiter.SlidingWindow
}

func (c *Service) Handle(ctx context.Context, msg twitch.TwitchChatMessage) struct{} {
	channel, err := c.getCachedChannel(ctx, msg.BroadcasterUserId)
	if err != nil {
		if errors.Is(err, channelsrepository.ErrNotFound) {
			return struct{}{}
		}
		c.logger.Error("cannot get channel", slog.Any("err", err))
		return struct{}{}
	}

	if msg.Message == nil || (msg.ChatterUserId == channel.BotID && c.config.AppEnv == "production") {
		return struct{}{}
	}

	resp, err := c.rateLimiter.Use(
		ctx, &redislimiter.SlidingWindowOptions{
			Key:             fmt.Sprintf("chat-translator:rate_limit:%s", msg.BroadcasterUserId),
			MaximumCapacity: 30,
			Window:          30 * time.Second,
		},
	)
	if err != nil {
		c.logger.Error("cannot use rate limiter", slog.Any("err", err))
		return struct{}{}
	}
	if !resp.Success {
		return struct{}{}
	}

	channelTranslationSettings, err := c.channelsTranslationsCache.Get(
		ctx,
		msg.BroadcasterUserId,
	)
	if err != nil {
		if errors.Is(err, channelschattrenslationsrepository.ErrSettingsNotFound) {
			return struct{}{}
		}
		c.logger.Error("cannot get channel translation settings", slog.Any("err", err))
		return struct{}{}
	}

	if channelTranslationSettings.ChannelID == "" ||
		!channelTranslationSettings.Enabled ||
		slices.Contains(channelTranslationSettings.ExcludedUsersIDs, msg.ChatterUserId) {
		return struct{}{}
	}

	textForDetect := msg.Message.Text
	for emoteName := range msg.UsedEmotesWithThirdParty {
		textForDetect = strings.ReplaceAll(textForDetect, emoteName, "")
	}

	msgLang, err := c.detectLanguage(ctx, textForDetect)
	if err != nil {
		c.logger.Error("cannot detect language", slog.Any("err", err))
		return struct{}{}
	}

	if len(msgLang.DetectedLanguages) == 0 {
		return struct{}{}
	}

	bestDetected := msgLang.DetectedLanguages[0]

	if slices.Contains(channelTranslationSettings.ExcludedLanguages, bestDetected.Language) {
		return struct{}{}
	}

	excludedWords := make([]string, 0, len(msg.UsedEmotesWithThirdParty))
	for k := range msg.UsedEmotesWithThirdParty {
		excludedWords = append(excludedWords, k)
	}

	res, err := c.translate(
		ctx,
		translateRequest{
			Text:          msg.Message.Text,
			SrcLang:       bestDetected.Language,
			DestLang:      channelTranslationSettings.TargetLanguage,
			ExcludedWords: excludedWords,
		},
	)
	if err != nil {
		c.logger.Error("cannot translate message", slog.Any("err", err))
		return struct{}{}
	}
	if res == nil || len(res.TranslatedText) == 0 {
		return struct{}{}
	}

	var resultText strings.Builder
	if channelTranslationSettings.UseItalic {
		resultText.WriteString("/me ")
	}
	resultText.WriteString(fmt.Sprintf("%s: %s", msg.ChatterUserName, res.TranslatedText[0]))

	if err := c.twirBus.Bots.SendMessage.Publish(
		bots.SendMessageRequest{
			ChannelName:       &msg.BroadcasterUserLogin,
			ChannelId:         msg.BroadcasterUserId,
			Message:           resultText.String(),
			SkipToxicityCheck: false,
		},
	); err != nil {
		c.logger.Error("cannot send message", slog.Any("err", err))
	}

	return struct{}{}
}
