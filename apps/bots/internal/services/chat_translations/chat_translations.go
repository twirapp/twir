package chat_translations

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"
	"unicode/utf8"

	googletranslate "cloud.google.com/go/translate"
	redislimiter "github.com/aidenwallis/go-ratelimiting/redis"
	redislimiteradapter "github.com/aidenwallis/go-ratelimiting/redis/adapters/go-redis"
	"github.com/google/uuid"
	"github.com/lkretschmer/deepl-go"
	"github.com/redis/go-redis/v9"
	"github.com/twirapp/kv"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/generic"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	channelschattrenslationsrepository "github.com/twirapp/twir/libs/repositories/chat_translation"
	"github.com/twirapp/twir/libs/repositories/chat_translation/model"
	googleapioption "google.golang.org/api/option"
)

type provider = func(c *Service, ctx context.Context, input translateRequest) (
	*translateResult,
	error,
)

var providers = []provider{
	(*Service).translateDeeplUnOfficial,
}

func New(
	lc *lifecycle.Lifecycle,
	cfg config.Config,
	logger *slog.Logger,
	twirBus *buscore.Bus,
	redisClient *redis.Client,
	kv kv.KV,
	channelsTranslationsRepository channelschattrenslationsrepository.Repository,
	channelsTranslationsCache *generic_cacher.GenericCacher[model.ChatTranslation],
) *Service {
	s := &Service{
		config:                         cfg,
		logger:                         logger,
		twirBus:                        twirBus,
		redis:                          redisClient,
		channelsTranslationsRepository: channelsTranslationsRepository,
		channelsTranslationsCache:      channelsTranslationsCache,
		rateLimiter:                    redislimiter.NewSlidingWindow(redislimiteradapter.NewAdapter(redisClient)),
		kv:                             kv,
	}

	if cfg.DeeplApiKey != "" {
		s.deeplClient = deepl.NewClient(cfg.DeeplApiKey)
		providers = append(
			[]provider{
				(*Service).translateDeeplOfficial,
			},
			providers...,
		)
	}

	if len(cfg.GoogleTranslateServiceAccountJson) > 0 {
		key, err := base64.StdEncoding.DecodeString(cfg.GoogleTranslateServiceAccountJson)
		if err != nil {
			panic(err)
		}

		cl, err := googletranslate.NewClient(
			context.TODO(),
			googleapioption.WithCredentialsJSON(key),
		)
		if err != nil {
			panic(err)
		}

		s.googleTranslateClient = cl

		providers = append(
			[]provider{
				(*Service).translateGoogleOfficial,
			},
			providers...,
		)
	}

	lc.Append(
		lifecycle.Hook{
			OnStop: func(ctx context.Context) error {
				if s.googleTranslateClient != nil {
					return s.googleTranslateClient.Close()
				}

				return nil
			},
		},
	)

	return s
}

type Service struct {
	config                         config.Config
	logger                         *slog.Logger
	twirBus                        *buscore.Bus
	redis                          *redis.Client
	kv                             kv.KV
	channelsTranslationsRepository channelschattrenslationsrepository.Repository
	channelsTranslationsCache      *generic_cacher.GenericCacher[model.ChatTranslation]

	rateLimiter redislimiter.SlidingWindow

	deeplClient           *deepl.Client
	googleTranslateClient *googletranslate.Client
}

type ChatMessageInput struct {
	Message                  generic.ChatMessage
	ChannelID                uuid.UUID
	BotUserID                *uuid.UUID
	ChannelCommandPrefix     string
	UsedEmotesWithThirdParty map[string]int
}

func (c *Service) Handle(ctx context.Context, input ChatMessageInput) error {
	msg := input.Message

	if msg.Message == nil || strings.HasPrefix(
		msg.Message.Text,
		input.ChannelCommandPrefix,
	) {
		return nil
	}

	if utf8.RuneCountInString(msg.Message.Text) < 5 {
		return nil
	}

	if input.BotUserID != nil && msg.UserID == input.BotUserID.String() {
		return nil
	}

	resp, err := c.rateLimiter.Use(
		ctx, &redislimiter.SlidingWindowOptions{
			Key:             fmt.Sprintf("chat-translator:rate_limit:%s", msg.BroadcasterUserId),
			MaximumCapacity: 30,
			Window:          30 * time.Second,
		},
	)
	if err != nil {
		c.logger.Error("cannot use rate limiter", logger.Error(err))
		return err
	}
	if !resp.Success {
		return nil
	}

	channelTranslationSettings, err := c.channelsTranslationsCache.Get(
		ctx,
		input.ChannelID.String(),
	)
	if err != nil {
		if errors.Is(err, channelschattrenslationsrepository.ErrSettingsNotFound) {
			return nil
		}
		c.logger.Error("cannot get channel translation settings", logger.Error(err))
		return err
	}

	if channelTranslationSettings.IsNil() {
		return nil
	}

	if channelTranslationSettings.ChannelID == "" ||
		!channelTranslationSettings.Enabled ||
		slices.Contains(channelTranslationSettings.ExcludedUsersIDs, msg.ChatterUserId) {
		return nil
	}

	if c.config.IsProduction() && input.BotUserID != nil && msg.UserID == input.BotUserID.String() {
		return nil
	}

	textForTranslate := msg.Message.Text
	for emoteName := range input.UsedEmotesWithThirdParty {
		textForTranslate = strings.ReplaceAll(textForTranslate, emoteName, "")
	}
	textForTranslate = prepareTextForTranslation(textForTranslate)

	if utf8.RuneCountInString(textForTranslate) < 5 || !hasTranslatableText(textForTranslate) {
		return nil
	}

	msgLang, err := c.detectLanguage(ctx, textForTranslate)
	if err != nil {
		c.logger.Error("cannot detect language", logger.Error(err))
		return err
	}

	if msgLang.Language == channelTranslationSettings.TargetLanguage {
		return nil
	}

	if isAmbiguousCyrillicDetection(
		textForTranslate,
		msgLang.Language,
		channelTranslationSettings.TargetLanguage,
	) {
		c.logger.Info(
			"skipping translation: ambiguous cyrillic language detection",
			slog.String("detected_lang", msgLang.Language),
			slog.Float64("confidence", msgLang.Confidence),
			slog.String("target_lang", channelTranslationSettings.TargetLanguage),
			slog.String("text", textForTranslate),
		)
		return nil
	}

	if slices.Contains(channelTranslationSettings.ExcludedLanguages, msgLang.Language) {
		return nil
	}

	excludedWords := make([]string, 0, len(input.UsedEmotesWithThirdParty))
	for k := range input.UsedEmotesWithThirdParty {
		excludedWords = append(excludedWords, k)
	}

	res, err := c.translate(
		ctx,
		translateRequest{
			Text:          textForTranslate,
			SrcLang:       msgLang.Language,
			DestLang:      channelTranslationSettings.TargetLanguage,
			ExcludedWords: excludedWords,
		},
	)
	if err != nil {
		c.logger.Error("cannot translate message", logger.Error(err))
		return err
	}
	if res == nil || len(res.TranslatedText) == 0 {
		return nil
	}

	if translationsEquivalent(textForTranslate, res.TranslatedText[0]) {
		return nil
	}

	var resultText strings.Builder
	if channelTranslationSettings.UseItalic {
		resultText.WriteString("/me ")
	}
	resultText.WriteString(fmt.Sprintf("%s: %s", msg.ChatterUserName, res.TranslatedText[0]))

	platformSource := platform.Platform(msg.Platform)
	if platformSource == "" {
		platformSource = platform.PlatformTwitch
	}

	if err := c.twirBus.Bots.SendMessage.Publish(
		ctx,
		bots.SendMessageRequest{
			ChannelID:         input.ChannelID,
			Platforms:         []platform.Platform{platformSource},
			Message:           resultText.String(),
			SkipToxicityCheck: false,
		},
	); err != nil {
		c.logger.Error("cannot send message", logger.Error(err))
		return err
	}

	return nil
}
