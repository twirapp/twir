package queue

import (
	"context"
	"errors"
	"github.com/satont/twir/apps/timers/internal/repositories/streams"
	"github.com/satont/twir/apps/timers/internal/repositories/timers"
	"github.com/satont/twir/libs/grpc/generated/bots"
	"github.com/satont/twir/libs/grpc/generated/parser"
	"log/slog"
)

func (c *Queue) handle(t *timer) {
	channel, err := c.channelsRepository.GetById(t.ChannelID)
	if err != nil {
		c.logger.Error("error on getting channel", slog.Any("error", err))
		return
	}

	if !channel.Enabled {
		c.logger.Info("channel not enabled", slog.String("channelId", t.ChannelID))
		return
	}

	if !channel.IsBotMod {
		c.logger.Info("bot is not moderator, response wont be sent", slog.String("channelId", t.ChannelID))
		return
	}

	stream, err := c.streamsRepository.GetByChannelId(t.ChannelID)
	if err != nil && errors.Is(err, streams.NotFound) && c.config.AppEnv != "development" {
		c.logger.Info("stream not found, probably channel is offline", slog.String("channelId", t.ChannelID))
		return
	} else if err != nil && c.config.AppEnv != "development" {
		c.logger.Info("error on getting stream", slog.String("channelId", t.ChannelID), slog.Any("error", err))
		return
	}

	var response timers.TimerResponse
	for index, r := range t.Responses {
		if index == t.LastResponse {
			response = r
			break
		}
	}

	// not found
	if response.ID == "" || response.Text == "" {
		c.logger.Info(
			"timer text or id is empty",
			slog.String("responseId", response.ID),
			slog.String("timerId", t.ID),
			slog.String("channelId", t.ChannelID),
			slog.Bool("isAnnounce", response.IsAnnounce),
			slog.String("text", response.Text),
			slog.String("timerName", t.Name),
		)
		return
	}

	// send
	err = c.sendMessage(stream, t.ChannelID, response.Text, response.IsAnnounce)
	if err != nil {
		c.logger.Error(
			"cannot send message", slog.String("responseId", response.ID),
			slog.String("timerId", t.ID),
			slog.String("channelId", t.ChannelID),
			slog.String("text", response.Text),
			slog.Any("error", err),
			slog.String("timerName", t.Name),
		)
	} else {
		c.logger.Info(
			"sent message", slog.String("responseId", response.ID),
			slog.String("timerId", t.ID),
			slog.String("channelId", t.ChannelID),
			slog.String("text", response.Text),
			slog.Any("error", err),
			slog.String("timerName", t.Name),
		)
	}

	// set next index
	nextIndex := t.LastResponse + 1

	if nextIndex < len(t.Responses) {
		t.LastResponse = nextIndex
	} else {
		t.LastResponse = 0
	}
}

func (c *Queue) sendMessage(stream streams.Stream, channelId, text string, isAnnounce bool) error {
	ctx := context.Background()

	parseReq, err := c.parserGrpc.ParseTextResponse(
		ctx,
		&parser.ParseTextRequestData{
			Sender: &parser.Sender{
				Id:          "",
				Name:        "bot",
				DisplayName: "Bot",
				Badges:      []string{"BROADCASTER"},
			},
			Channel: &parser.Channel{Id: stream.UserID, Name: stream.UserLogin},
			Message: &parser.Message{Text: text},
		},
	)
	if err != nil {
		return err
	}

	for _, message := range parseReq.Responses {
		_, err = c.botsGrpc.SendMessage(
			ctx,
			&bots.SendMessageRequest{
				ChannelId:      channelId,
				ChannelName:    &stream.UserLogin,
				Message:        message,
				IsAnnounce:     &isAnnounce,
				SkipRateLimits: true,
			},
		)
		if err != nil {
			return err
		}
	}

	return nil
}
