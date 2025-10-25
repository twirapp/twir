package handler

import (
	"context"
	"log/slog"

	"github.com/kvizyx/twitchy/eventsub"
	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/events"
)

func convertOutCome(outcomes []eventsub.ChannelPredictionEventOutcome) []events.PredictionOutcome {
	out := make([]events.PredictionOutcome, 0, len(outcomes))

	for _, outcome := range outcomes {
		topPredictors := make(
			[]events.PredictionTopPredictor,
			0,
			len(outcome.TopPredictors),
		)

		for _, predictor := range outcome.TopPredictors {
			won := uint64(predictor.ChannelPointsWon)
			topPredictors = append(
				topPredictors, events.PredictionTopPredictor{
					UserName:        predictor.UserLogin,
					UserDisplayName: predictor.UserName,
					UserID:          predictor.UserId,
					PointsUsed:      uint64(predictor.ChannelPointsUsed),
					PointsWin: lo.
						If(predictor.ChannelPointsWon > 0, &won).
						Else(nil),
				},
			)
		}

		out = append(
			out, events.PredictionOutcome{
				ID:            outcome.Id,
				Title:         outcome.Title,
				Color:         outcome.Color,
				Users:         uint64(outcome.Users),
				ChannelPoints: uint64(outcome.ChannelPoints),
				TopPredictors: topPredictors,
			},
		)
	}
	return out
}

func (c *Handler) HandleChannelPredictionBegin(
	ctx context.Context,
	event eventsub.ChannelPredictionBeginEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"Prediction begin",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("title", event.Title),
	)

	outComes := convertOutCome(event.Outcomes)

	err := c.twirBus.Events.PredictionBegin.Publish(
		ctx,
		events.PredictionBeginMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserName:        event.BroadcasterUserLogin,
			UserDisplayName: event.BroadcasterUserName,
			Info: events.PredictionInfo{
				Title:    event.Title,
				Outcomes: outComes,
			},
		},
	)

	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}

func (c *Handler) HandleChannelPredictionProgress(
	ctx context.Context,
	event eventsub.ChannelPredictionProgressEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"Prediction progress",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("title", event.Title),
	)

	outComes := convertOutCome(event.Outcomes)

	err := c.twirBus.Events.PredictionProgress.Publish(
		ctx,
		events.PredictionProgressMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserName:        event.BroadcasterUserLogin,
			UserDisplayName: event.BroadcasterUserName,
			Info: events.PredictionInfo{
				Title:    event.Title,
				Outcomes: outComes,
			},
		},
	)

	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}

func (c *Handler) HandleChannelPredictionLock(
	ctx context.Context,
	event eventsub.ChannelPredictionLockEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"Prediction lock",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("title", event.Title),
	)

	outComes := convertOutCome(event.Outcomes)

	err := c.twirBus.Events.PredictionLock.Publish(
		ctx,
		events.PredictionLockMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserName:        event.BroadcasterUserLogin,
			UserDisplayName: event.BroadcasterUserName,
			Info: events.PredictionInfo{
				Title:    event.Title,
				Outcomes: outComes,
			},
		},
	)

	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}

func (c *Handler) HandleChannelPredictionEnd(
	ctx context.Context,
	event eventsub.ChannelPredictionEndEvent,
	meta eventsub.WebsocketNotificationMetadata,
) {
	c.logger.Info(
		"Prediction end",
		slog.String("channelId", event.BroadcasterUserId),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("title", event.Title),
		slog.String("status", string(event.Status)),
	)

	if event.Status != "resolved" {
		return
	}

	outComes := convertOutCome(event.Outcomes)

	err := c.twirBus.Events.PredictionEnd.Publish(
		ctx,
		events.PredictionEndMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserId,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserName:        event.BroadcasterUserLogin,
			UserDisplayName: event.BroadcasterUserName,
			Info: events.PredictionInfo{
				Title:    event.Title,
				Outcomes: outComes,
			},
			WinningOutcomeID: event.WinningOutcomeId,
		},
	)

	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}
