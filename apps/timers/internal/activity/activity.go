package activity

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/satont/twir/apps/timers/internal/repositories/channels"
	"github.com/satont/twir/apps/timers/internal/repositories/streams"
	"github.com/satont/twir/apps/timers/internal/repositories/timers"
	config "github.com/satont/twir/libs/config"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/bots"
	busparser "github.com/twirapp/twir/libs/bus-core/parser"
	"github.com/twirapp/twir/libs/grpc/parser"
	"github.com/twirapp/twir/libs/redis_keys"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	TimersRepository   timers.Repository
	ChannelsRepository channels.Repository
	StreamsRepository  streams.Repository
	Cfg                config.Config
	ParserGrpc         parser.ParserClient
	Redis              *redis.Client
	Bus                *buscore.Bus
}

func New(opts Opts) *Activity {
	return &Activity{
		timersRepository:   opts.TimersRepository,
		channelsRepository: opts.ChannelsRepository,
		streamsRepository:  opts.StreamsRepository,
		cfg:                opts.Cfg,
		parserGrpc:         opts.ParserGrpc,
		redis:              opts.Redis,
		bus:                opts.Bus,
	}
}

type Activity struct {
	timersRepository   timers.Repository
	channelsRepository channels.Repository
	streamsRepository  streams.Repository
	cfg                config.Config
	parserGrpc         parser.ParserClient
	redis              *redis.Client
	bus                *buscore.Bus
}

func getLastTriggerMessageNumberKey(timerId, streamId string) string {
	return fmt.Sprintf("timers:last_trigger_message_number:%s:%s", timerId, streamId)
}

func (c *Activity) SendMessage(ctx context.Context, timerId string) (
	int,
	error,
) {
	timer, err := c.timersRepository.GetById(timerId)
	if err != nil {
		return 0, err
	}

	channel, err := c.channelsRepository.GetById(timer.ChannelID)
	if err != nil {
		return 0, err
	}

	if !channel.Enabled {
		return 0, nil
	}

	if !channel.IsBotMod {
		return 0, nil
	}

	currentResponse, err := c.redis.Get(ctx, redis_keys.TimersCurrentResponse(timerId)).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return currentResponse, err
	}

	stream, err := c.streamsRepository.GetByChannelId(channel.ID)
	if err != nil && errors.Is(err, streams.NotFound) && c.cfg.AppEnv != "development" {
		return currentResponse, nil
	} else if err != nil && c.cfg.AppEnv != "development" {
		return currentResponse, err
	}

	lastTriggerMessageNumber, err := c.redis.Get(
		ctx,
		getLastTriggerMessageNumberKey(timerId, stream.ID),
	).
		Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, err
	}

	parsedMessages, err := c.redis.Get(
		ctx,
		redis_keys.StreamParsedMessages(stream.ID),
	).Int()

	if timer.MessageInterval != 0 &&
		lastTriggerMessageNumber-parsedMessages+timer.MessageInterval > 0 {
		return currentResponse, nil
	}

	var response timers.TimerResponse
	for index, r := range timer.Responses {
		if index == currentResponse {
			response = r
			break
		}
	}

	if response.ID == "" || response.Text == "" {
		return currentResponse, nil
	}

	err = c.sendMessage(ctx, stream, channel.ID, response.Text, response.IsAnnounce)
	if err != nil {
		return currentResponse, err
	}

	nextIndex := currentResponse + 1

	if nextIndex >= len(timer.Responses) {
		nextIndex = 0
	}

	err = c.redis.Set(ctx, redis_keys.TimersCurrentResponse(timerId), nextIndex, 24*time.Hour).Err()
	if err != nil {
		return nextIndex, err
	}

	err = c.redis.Set(
		ctx, getLastTriggerMessageNumberKey(timerId, stream.ID), parsedMessages,
		24*time.Hour,
	).Err()
	if err != nil {
		return nextIndex, err
	}

	return nextIndex, nil
}

func (c *Activity) sendMessage(
	ctx context.Context,
	stream streams.Stream,
	channelId, text string,
	isAnnounce bool,
) error {
	parseReq, err := c.bus.Parser.ParseVariablesInText.Request(
		ctx,
		busparser.ParseVariablesInTextRequest{
			ChannelID:   stream.UserID,
			ChannelName: stream.UserLogin,
			Text:        text,
		},
	)
	if err != nil {
		return err
	}

	err = c.bus.Bots.SendMessage.Publish(
		bots.SendMessageRequest{
			ChannelId:      channelId,
			ChannelName:    &stream.UserLogin,
			Message:        parseReq.Data.Text,
			IsAnnounce:     isAnnounce,
			SkipRateLimits: true,
		},
	)
	if err != nil {
		return err
	}

	return nil
}
