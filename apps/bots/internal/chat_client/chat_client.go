package chat_client

import (
	"context"
	"fmt"
	"strings"
	"sync"
	"time"
	"unicode/utf8"

	ratelimiting "github.com/aidenwallis/go-ratelimiting/local"
	irc "github.com/gempir/go-twitch-irc/v3"
	"github.com/nicklaw5/helix/v2"
	"github.com/redis/go-redis/v9"
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
}

type ChatClient struct {
	services *services

	Readers []*BotClientIrc
	Writer  *BotClientIrc

	// channelId:writer
	channelsToReader *utils.SyncMap[*BotClientIrc]
	joinMu           *sync.Mutex
	workersPool      *gopool.Pool

	RateLimiters RateLimiters
	Model        *model.Bots
	TwitchUser   *helix.User
}

type services struct {
	DB             *gorm.DB
	Cfg            cfg.Config
	Logger         logger.Logger
	Model          *model.Bots
	ParserGrpc     parser.ParserClient
	TokensGrpc     tokens.TokensClient
	EventsGrpc     events.EventsClient
	WebsocketsGrpc websockets.WebsocketClient
	Redis          *redis.Client
	TwitchClient   *helix.Client
}

type Opts struct {
	DB             *gorm.DB
	Cfg            cfg.Config
	Logger         logger.Logger
	Model          *model.Bots
	ParserGrpc     parser.ParserClient
	TokensGrpc     tokens.TokensClient
	EventsGrpc     events.EventsClient
	WebsocketsGrpc websockets.WebsocketClient
	Redis          *redis.Client
}

func New(opts Opts) *ChatClient {
	globalRateLimiter, _ := ratelimiting.NewSlidingWindow(100, 30*time.Second)

	twitchClient, err := twitch.NewBotClient(opts.Model.ID, opts.Cfg, opts.TokensGrpc)
	if err != nil {
		panic(err)
	}

	meReq, _ := twitchClient.GetUsers(
		&helix.UsersParams{
			IDs: []string{opts.Model.ID},
		},
	)
	if len(meReq.Data.Users) == 0 {
		panic("No user found for bot " + opts.Model.ID)
	}

	me := meReq.Data.Users[0]

	s := &ChatClient{
		joinMu:           &sync.Mutex{},
		channelsToReader: utils.NewSyncMap[*BotClientIrc](),
		services: &services{
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
		},
		TwitchUser: &me,
		Model:      opts.Model,
		RateLimiters: RateLimiters{
			Global: globalRateLimiter,
			Channels: ChannelsRateLimiter{
				Items: make(map[string]*Channel),
			},
		},
		workersPool: gopool.NewPool(1000),
	}
	s.CreateWriter()

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
