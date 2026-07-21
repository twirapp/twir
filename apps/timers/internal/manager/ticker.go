package manager

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/redis/go-redis/v9"
	timersentity "github.com/twirapp/twir/libs/entities/timers"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/redis_keys"
)

func (c *Manager) tryTick(id TimerID) {
	t, ok := c.timers[id]
	if !ok {
		return
	}

	ctx := context.Background()

	channel, err := c.channelCachedRepo.Get(ctx, t.dbRow.ChannelID.String())
	if err != nil {
		c.logger.Error(
			"[tick] cannot get channel",
			logger.Error(err),
			slog.String("channelId", t.dbRow.ChannelID.String()),
			slog.String("timerId", id.String()),
		)
		return
	}

	if !channel.IsBotMod || !channel.IsEnabled {
		return
	}

	streams, err := c.channelservice.GetChannelStreams(ctx, channel.ID)
	if err != nil {
		c.logger.Error(
			"[tick] cannot get channel stream",
			logger.Error(err),
			slog.String("channelId", t.dbRow.ChannelID.String()),
			slog.String("timerId", id.String()),
		)
		return
	}

	isOffline := len(streams) == 0

	if isOffline {
		if !t.dbRow.OfflineEnabled {
			return
		}
		t.offlineMessageNumber++
	} else {
		if !t.dbRow.OnlineEnabled {
			return
		}
		t.offlineMessageNumber = 0
		t.lastTriggerOfflineNumber = 0
	}

	var (
		currentMessageNumber    int
		lastTriggerMessageCount int
	)

	if isOffline {
		currentMessageNumber = t.offlineMessageNumber
		lastTriggerMessageCount = t.lastTriggerOfflineNumber
	} else {
		var streamParsedMessages uint64
		for _, stream := range streams {
			streamParsedMessages += stream.ParsedChatLines
		}

		currentMessageNumber = int(streamParsedMessages)
		lastTriggerMessageCount = t.lastTriggerMessageNumber
	}

	var (
		now              = time.Now()
		shouldSend       bool
		timeInterval     = time.Duration(t.dbRow.TimeInterval) * time.Minute
		messageInterval  = t.dbRow.MessageInterval
		secondsSinceLast = now.Sub(t.lastTriggerTimestamp).
					Seconds() -
			5 // https://go.dev/pkg/time/?m=old#hdr-Timer_Resolution
	)

	if currentMessageNumber < lastTriggerMessageCount {
		if isOffline {
			t.lastTriggerOfflineNumber = currentMessageNumber
		} else {
			t.lastTriggerMessageNumber = currentMessageNumber
		}
		lastTriggerMessageCount = currentMessageNumber
	}

	messagesSinceLast := currentMessageNumber - lastTriggerMessageCount

	switch {
	case timeInterval == 0 && messageInterval > 0:
		if messagesSinceLast >= messageInterval {
			shouldSend = true
		}
		break
	case timeInterval > 0 && messageInterval == 0:
		if secondsSinceLast >= timeInterval.Seconds() {
			shouldSend = true
		}
		break
	case timeInterval > 0 && messageInterval > 0:
		var (
			timeTriggered bool
			msgTriggered  bool
		)
		if secondsSinceLast >= timeInterval.Seconds() {
			timeTriggered = true
		}
		if messagesSinceLast >= messageInterval {
			msgTriggered = true
		}
		shouldSend = timeTriggered && msgTriggered
	}

	if !shouldSend {
		return
	}

	targets := getTimerSendTargets(channel, t.dbRow.Platforms)
	if len(targets) == 0 {
		return
	}

	var response timersentity.Response
	for index, r := range t.dbRow.Responses {
		if index == t.currentResponseIndex {
			response = r
			break
		}
	}

	var twitchUserID string
	if channel.TwitchUserID != nil {
		twitchUserID = channel.TwitchUserID.String()
	}

	wasSent := false
	for _, target := range targets {
		err = c.sendMessage(
			ctx,
			target.channelID,
			twitchUserID,
			channel.ID.String(),
			response.Text,
			response.IsAnnounce,
			response.AnnounceColor,
			response.Count,
			string(target.platform),
		)
		if err != nil {
			c.logger.Error(
				"[tick] cannot send timer message",
				logger.Error(err),
				slog.String("channelId", t.dbRow.ChannelID.String()),
				slog.String("timerId", id.String()),
				slog.String("platform", target.platform.String()),
			)
			continue
		}

		wasSent = true
	}

	if !wasSent {
		return
	}

	if isOffline {
		t.lastTriggerOfflineNumber = currentMessageNumber
	} else {
		t.lastTriggerMessageNumber = currentMessageNumber
	}
	t.lastTriggerTimestamp = now

	nextIndex := t.currentResponseIndex + 1

	if nextIndex >= len(t.dbRow.Responses) {
		nextIndex = 0
	}

	t.currentResponseIndex = nextIndex
}

func (c *Manager) getStreamChatLines(ctx context.Context, streamID string) (int, error) {
	streamParsedMessages, err := c.redis.Get(
		ctx,
		redis_keys.StreamParsedMessages(streamID),
	).Int()
	if err != nil && !errors.Is(err, redis.Nil) {
		return 0, fmt.Errorf("failed to get stream parsed messages: %w", err)
	}

	return streamParsedMessages, nil
}
