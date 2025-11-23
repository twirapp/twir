package handle_message

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	googletranslate "cloud.google.com/go/translate"
	redislimiter "github.com/aidenwallis/go-ratelimiting/redis"
	redislimiteradapter "github.com/aidenwallis/go-ratelimiting/redis/adapters/go-redis"
	"github.com/lkretschmer/deepl-go"
	"github.com/redis/go-redis/v9"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	channelsrepository "github.com/twirapp/twir/libs/repositories/channels"
	channelmodel "github.com/twirapp/twir/libs/repositories/channels/model"
	channelschattrenslationsrepository "github.com/twirapp/twir/libs/repositories/chat_translation"
	"github.com/twirapp/twir/libs/repositories/chat_translation/model"
	"go.uber.org/fx"
	googleapioption "google.golang.org/api/option"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Config  config.Config
	Logger  logger.Logger
	TwirBus *buscore.Bus
	Redis   *redis.Client

	ChannelsRepository             channelsrepository.Repository
	ChannelsTranslationsRepository channelschattrenslationsrepository.Repository
	ChannelsTranslationsCache      *generic_cacher.GenericCacher[model.ChatTranslation]
	ChannelsCache                  *generic_cacher.GenericCacher[channelmodel.Channel]
}

type provider = func(c *Service, ctx context.Context, input translateRequest) (
	*translateResult,
	error,
)

var providers = []provider{
	(*Service).translateDeeplUnOfficial,
}

func New(opts Opts) *Service {
	s := &Service{
		config:                         opts.Config,
		logger:                         opts.Logger,
		twirBus:                        opts.TwirBus,
		redis:                          opts.Redis,
		channelsRepository:             opts.ChannelsRepository,
		channelsTranslationsRepository: opts.ChannelsTranslationsRepository,
		channelsTranslationsCache:      opts.ChannelsTranslationsCache,
		rateLimiter:                    redislimiter.NewSlidingWindow(redislimiteradapter.NewAdapter(opts.Redis)),
		channelsCache:                  opts.ChannelsCache,
	}

	if opts.Config.DeeplApiKey != "" {
		s.deeplClient = deepl.NewClient(opts.Config.DeeplApiKey)
		providers = append(
			[]provider{
				(*Service).translateDeeplOfficial,
			},
			providers...,
		)
	}

	if len(opts.Config.GoogleTranslateServiceAccountJson) > 0 {
		key, err := base64.StdEncoding.DecodeString(opts.Config.GoogleTranslateServiceAccountJson)
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

		time.Sleep(5 * time.Second)

		fmt.Println(
			s.translateGoogleOfficial(
				context.TODO(), translateRequest{
					Text:          "Hello world",
					SrcLang:       "en",
					DestLang:      "ru",
					ExcludedWords: nil,
				},
			),
		)
	}

	opts.LC.Append(
		fx.Hook{
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
	logger                         logger.Logger
	twirBus                        *buscore.Bus
	redis                          *redis.Client
	channelsRepository             channelsrepository.Repository
	channelsTranslationsRepository channelschattrenslationsrepository.Repository
	channelsTranslationsCache      *generic_cacher.GenericCacher[model.ChatTranslation]
	channelsCache                  *generic_cacher.GenericCacher[channelmodel.Channel]

	rateLimiter redislimiter.SlidingWindow

	deeplClient           *deepl.Client
	googleTranslateClient *googletranslate.Client
}

func (c *Service) Handle(ctx context.Context, msg twitch.TwitchChatMessage) (struct{}, error) {
	if msg.Message == nil || strings.HasPrefix(
		msg.Message.Text,
		msg.EnrichedData.ChannelCommandPrefix,
	) {
		return struct{}{}, nil
	}

	// if msg.ChatterUserId == msg.EnrichedData.DbChannel.BotID {
	// 	return struct{}{}
	// }

	resp, err := c.rateLimiter.Use(
		ctx, &redislimiter.SlidingWindowOptions{
			Key:             fmt.Sprintf("chat-translator:rate_limit:%s", msg.BroadcasterUserId),
			MaximumCapacity: 30,
			Window:          30 * time.Second,
		},
	)
	if err != nil {
		c.logger.Error("cannot use rate limiter", slog.Any("err", err))
		return struct{}{}, err
	}
	if !resp.Success {
		return struct{}{}, nil
	}

	channelTranslationSettings, err := c.channelsTranslationsCache.Get(
		ctx,
		msg.BroadcasterUserId,
	)
	if err != nil {
		if errors.Is(err, channelschattrenslationsrepository.ErrSettingsNotFound) {
			return struct{}{}, nil
		}
		c.logger.Error("cannot get channel translation settings", slog.Any("err", err))
		return struct{}{}, err
	}

	if channelTranslationSettings.ChannelID == "" ||
		!channelTranslationSettings.Enabled ||
		slices.Contains(channelTranslationSettings.ExcludedUsersIDs, msg.ChatterUserId) {
		return struct{}{}, nil
	}

	channel, err := c.channelsCache.Get(ctx, msg.BroadcasterUserId)
	if err != nil {
		c.logger.Error("cannot get channel", slog.Any("err", err))
		return struct{}{}, err
	}
	if c.config.IsProduction() && msg.ChatterUserId == channel.BotID {
		return struct{}{}, nil
	}

	textForDetect := msg.Message.Text
	for emoteName := range msg.EnrichedData.UsedEmotesWithThirdParty {
		textForDetect = strings.ReplaceAll(textForDetect, emoteName, "")
	}

	msgLang, err := c.detectLanguage(ctx, textForDetect)
	if err != nil {
		c.logger.Error("cannot detect language", slog.Any("err", err))
		return struct{}{}, err
	}

	if len(msgLang.DetectedLanguages) == 0 {
		return struct{}{}, nil
	}

	if msgLang.DetectedLanguages[0].Language == channelTranslationSettings.TargetLanguage {
		return struct{}{}, nil
	}

	bestDetected := msgLang.DetectedLanguages[0]

	if slices.Contains(channelTranslationSettings.ExcludedLanguages, bestDetected.Language) {
		return struct{}{}, nil
	}

	excludedWords := make([]string, 0, len(msg.EnrichedData.UsedEmotesWithThirdParty))
	for k := range msg.EnrichedData.UsedEmotesWithThirdParty {
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
		return struct{}{}, err
	}
	if res == nil || len(res.TranslatedText) == 0 {
		return struct{}{}, nil
	}

	if res.TranslatedText[0] == msg.Message.Text {
		return struct{}{}, nil
	}

	var resultText strings.Builder
	if channelTranslationSettings.UseItalic {
		resultText.WriteString("/me ")
	}
	resultText.WriteString(fmt.Sprintf("%s: %s", msg.ChatterUserName, res.TranslatedText[0]))

	if err := c.twirBus.Bots.SendMessage.Publish(
		ctx,
		bots.SendMessageRequest{
			ChannelName:       &msg.BroadcasterUserLogin,
			ChannelId:         msg.BroadcasterUserId,
			Message:           resultText.String(),
			SkipToxicityCheck: false,
		},
	); err != nil {
		c.logger.Error("cannot send message", slog.Any("err", err))
		return struct{}{}, err
	}

	return struct{}{}, nil
}
