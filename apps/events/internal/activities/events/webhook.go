package events

import (
	"bytes"
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/goccy/go-json"
	"github.com/twirapp/twir/apps/events/internal/shared"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/events/model"
	"go.temporal.io/sdk/activity"
)

type sendHttpRequestData struct {
	EventType string           `json:"event_type"`
	ChannelID string           `json:"channel_id"`
	EventData shared.EventData `json:"event_data"`
}

func (c *Activity) SendHttpRequest(
	ctx context.Context,
	operation model.EventOperation,
	data shared.EventData,
) error {
	activity.RecordHeartbeat(ctx, nil)

	if operation.Input == nil || *operation.Input == "" {
		return fmt.Errorf("input (webhook URL) is required for send http request operation")
	}

	body, err := json.Marshal(
		sendHttpRequestData{
			EventType: data.EventType,
			ChannelID: data.ChannelID,
			EventData: data,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot marshal event data: %w", err)
	}

	for i := range 3 {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			*operation.Input,
			bytes.NewReader(body),
		)
		if err != nil {
			return fmt.Errorf("cannot create http request: %w", err)
		}

		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("X-Twir-Event-Type", data.EventType)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			c.logger.Error(
				"failed to send http request, will retry",
				logger.Error(err),
				slog.String("url", *operation.Input),
				slog.String("eventType", data.EventType),
				slog.String("channelId", data.ChannelID),
				slog.Int("bodySize", len(body)),
				slog.Int("statusCode", resp.StatusCode),
			)
			activity.RecordHeartbeat(ctx, nil)
			time.Sleep(time.Second * 2 * time.Duration(i+1))
			continue
		}
		resp.Body.Close()

		if resp.StatusCode < 200 || resp.StatusCode >= 300 {
			c.logger.Error(
				"received non-success status code, will retry",
				slog.Int("statusCode", resp.StatusCode),
				slog.String("url", *operation.Input),
				slog.String("eventType", data.EventType),
				slog.String("channelId", data.ChannelID),
				slog.Int("bodySize", len(body)),
				slog.Int("statusCode", resp.StatusCode),
			)
			activity.RecordHeartbeat(ctx, nil)
			time.Sleep(time.Second * 2 * time.Duration(i+1))
			continue
		}

		c.logger.Info(
			"successfully sent http request",
			slog.String("url", *operation.Input),
			slog.String("eventType", data.EventType),
			slog.String("channelId", data.ChannelID),
			slog.Int("bodySize", len(body)),
			slog.Int("statusCode", resp.StatusCode),
		)
		return nil
	}

	c.logger.Warn(
		"failed to send http request after three attempts",
		slog.String("url", *operation.Input),
		slog.String("url", *operation.Input),
		slog.String("eventType", data.EventType),
		slog.String("channelId", data.ChannelID),
		slog.Int("bodySize", len(body)),
	)

	return fmt.Errorf("failed to send http request after three attempts")
}
