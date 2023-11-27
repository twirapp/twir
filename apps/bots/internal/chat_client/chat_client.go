package chat_client

import (
	"context"
	"fmt"
	"log/slog"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	ratelimiting "github.com/aidenwallis/go-ratelimiting/local"
	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/nicklaw5/helix/v2"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/bots/internal/moderation_helpers"
	"github.com/satont/twir/apps/bots/pkg/tlds"
	cfg "github.com/satont/twir/libs/config"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/gopool"
	"github.com/satont/twir/libs/grpc/generated/events"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"github.com/satont/twir/libs/grpc/generated/tokens"
	"github.com/satont/twir/libs/grpc/generated/websockets"
	"github.com/satont/twir/libs/logger"
	"github.com/satont/twir/libs/twitch"
	"github.com/satont/twir/libs/utils"
	"gorm.io/gorm"
)

type Channel struct {
	IsMod   bool
	Limiter ratelimiting.SlidingWindow
	ID      string
}

type ChannelsRateLimiter struct {
	sync.Mutex
	Items map[string]*Channel
}

type RateLimiters struct {
	Global   ratelimiting.SlidingWindow
	Channels ChannelsRateLimiter
}

type BotClientIrc struct {
	*irc.Client
	size            int8
	disconnectChann chan struct{}
	Connected       bool
}

type ChatClientPrometheus struct {
	messagesCounter prometheus.Counter
	readersCounter  prometheus.Gauge
	channelsCounter prometheus.Gauge
}

type ChatClient struct {
	services *services

	Readers []*BotClientIrc
	Writer  *BotClientIrc

	counters *ChatClientPrometheus

	// channelId:writer
	channelsToReader *utils.SyncMap[*BotClientIrc]
	joinMu           *sync.Mutex
	workersPool      *gopool.Pool
	joinRateLimiter  ratelimiting.SlidingWindow

	RateLimiters RateLimiters
	Model        *model.Bots
	TwitchUser   *helix.User

	moderationService *moderationService
}

type services struct {
	DB           *gorm.DB
	Cfg          cfg.Config
	Logger       logger.Logger
	Model        *model.Bots
	Redis        *redis.Client
	TwitchClient *helix.Client
	tlds         *tlds.TLDS

	ParserGrpc     parser.ParserClient
	TokensGrpc     tokens.TokensClient
	EventsGrpc     events.EventsClient
	WebsocketsGrpc websockets.WebsocketClient
}

type Opts struct {
	DB              *gorm.DB
	Cfg             cfg.Config
	Logger          logger.Logger
	Model           *model.Bots
	ParserGrpc      parser.ParserClient
	TokensGrpc      tokens.TokensClient
	EventsGrpc      events.EventsClient
	WebsocketsGrpc  websockets.WebsocketClient
	Redis           *redis.Client
	JoinRateLimiter ratelimiting.SlidingWindow
	Tlds            *tlds.TLDS
}

func New(opts Opts) *ChatClient {
	globalChatRateLimiter, _ := ratelimiting.NewSlidingWindow(100, 30*time.Second)

	twitchClient, err := twitch.NewBotClient(opts.Model.ID, opts.Cfg, opts.TokensGrpc)
	if err != nil {
		panic(err)
	}

	go func() {
		ticker := time.NewTicker(30 * time.Second)
		for {
			select {
			case <-ticker.C:
				res, err := opts.TokensGrpc.RequestBotToken(
					context.Background(),
					&tokens.GetBotTokenRequest{BotId: opts.Model.ID},
				)
				if err != nil {
					opts.Logger.Error("cannot refresh bot token", slog.Any("err", err))
					continue
				}

				twitchClient.SetUserAccessToken(res.AccessToken)
			}
		}
	}()

	meReq, err := twitchClient.GetUsers(
		&helix.UsersParams{
			IDs: []string{opts.Model.ID},
		},
	)
	if err != nil {
		panic(err)
	}
	if meReq.ErrorMessage != "" {
		panic(meReq.ErrorMessage)
	}
	if len(meReq.Data.Users) == 0 {
		panic("No user found for bot " + opts.Model.ID)
	}

	me := meReq.Data.Users[0]

	linksR, linksWithSpacesR := moderation_helpers.BuildLinksModerationRegexps(opts.Tlds.List)

	serv := &services{
		DB:             opts.DB,
		Cfg:            opts.Cfg,
		Logger:         opts.Logger,
		Model:          opts.Model,
		ParserGrpc:     opts.ParserGrpc,
		TokensGrpc:     opts.TokensGrpc,
		EventsGrpc:     opts.EventsGrpc,
		WebsocketsGrpc: opts.WebsocketsGrpc,
		Redis:          opts.Redis,
		TwitchClient:   twitchClient,
		tlds:           opts.Tlds,
	}

	s := &ChatClient{
		joinMu:           &sync.Mutex{},
		channelsToReader: utils.NewSyncMap[*BotClientIrc](),
		services:         serv,
		TwitchUser:       &me,
		Model:            opts.Model,
		RateLimiters: RateLimiters{
			Global: globalChatRateLimiter,
			Channels: ChannelsRateLimiter{
				Items: make(map[string]*Channel),
			},
		},
		workersPool:     gopool.NewPool(1000),
		joinRateLimiter: opts.JoinRateLimiter,
		moderationService: &moderationService{
			services:               serv,
			linksRegexp:            linksR,
			linksWithSpacesRegexp:  linksWithSpacesR,
			messagesTimeouterStore: utils.NewTtlSyncMap[struct{}](10 * time.Second),
		},
		counters: &ChatClientPrometheus{
			messagesCounter: promauto.NewCounter(
				prometheus.CounterOpts{
					Name: "bots_messages_counter",
					Help: "The total number of processed messages",
					ConstLabels: prometheus.Labels{
						"botName": me.Login,
						"botId":   me.ID,
					},
				},
			),
			readersCounter: promauto.NewGauge(
				prometheus.GaugeOpts{
					Name:        "bots_readers_counter",
					ConstLabels: prometheus.Labels{"botName": me.Login, "botId": me.ID},
				},
			),
			channelsCounter: promauto.NewGauge(
				prometheus.GaugeOpts{
					Name:        "bots_channels_counter",
					ConstLabels: prometheus.Labels{"botName": me.Login, "botId": me.ID},
				},
			),
		},
	}
	s.CreateWriter()

	go func() {
		for {
			s.counters.readersCounter.Set(float64(len(s.Readers)))
			s.counters.channelsCounter.Set(float64(s.channelsToReader.Len()))

			time.Sleep(5 * time.Second)
		}
	}()

	channels, err := s.getChannels()
	if err != nil {
		panic(err)
	}

	opts.Logger.Info(fmt.Sprintf("Joining %v channels", len(channels)))
	for _, ch := range channels {
		s.Join(ch)
	}

	return s
}

type SayOpts struct {
	Channel   string
	Text      string
	ReplyTo   *string
	WithLimit bool
}

func (c *ChatClient) Say(opts SayOpts) {
	text := strings.ReplaceAll(opts.Text, "\n", " ")
	if len(text) == 0 {
		return
	}
	textParts := splitTextByLength(text)

	c.RateLimiters.Global.Wait(context.Background())

	if opts.WithLimit {
		channelLimiter, ok := c.RateLimiters.Channels.Items[strings.ToLower(opts.Channel)]
		if ok {

			// it should be separately
			if !channelLimiter.Limiter.TryTake() {
				return
			}
		}
	}

	if opts.ReplyTo != nil {
		for _, part := range textParts {
			text = validateResponseSlashes(text)
			c.Writer.Reply(opts.Channel, *opts.ReplyTo, part)
		}
	} else {
		for _, part := range textParts {
			text = validateResponseSlashes(text)
			c.Writer.Say(opts.Channel, part)
		}
	}
}

func validateResponseSlashes(response string) string {
	if strings.HasPrefix(response, "/me") || strings.HasPrefix(response, "/announce") {
		return response
	} else if strings.HasPrefix(response, "/") {
		return "Slash commands except /me and /announce is disallowed. This response wont be ever sended."
	} else if strings.HasPrefix(response, ".") {
		return `Message cannot start with "." symbol.`
	} else {
		return response
	}
}

func splitTextByLength(text string) []string {
	var parts []string

	i := 500
	for utf8.RuneCountInString(text) > 0 {
		if utf8.RuneCountInString(text) < 500 {
			parts = append(parts, text)
			break
		}
		runned := []rune(text)
		parts = append(parts, string(runned[:i]))
		text = string(runned[i:])
	}

	return parts
}
