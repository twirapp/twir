package handler

import (
	"log/slog"

	"github.com/samber/lo"
	"github.com/twirapp/twir/libs/bus-core/events"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
)

func convertOutCome(outcomes []eventsub_bindings.PredictionOutcome) []events.PredictionOutcome {
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
					UserID:          predictor.UserID,
					PointsUsed:      uint64(predictor.ChannelPointsUsed),
					PointsWin: lo.
						If(predictor.ChannelPointsWon > 0, &won).
						Else(nil),
				},
			)
		}

		out = append(
			out, events.PredictionOutcome{
				ID:            outcome.ID,
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

func (c *Handler) handleChannelPredictionBegin(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPredictionBegin,
) {
	c.logger.Info(
		"Prediction begin",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("title", event.Title),
	)

	outComes := convertOutCome(event.Outcomes)

	err := c.twirBus.Events.PredictionBegin.Publish(
		events.PredictionBeginMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserID,
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

func (c *Handler) handleChannelPredictionProgress(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPredictionProgress,
) {
	c.logger.Info(
		"Prediction progress",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("title", event.Title),
	)

	outComes := convertOutCome(event.Outcomes)

	err := c.twirBus.Events.PredictionProgress.Publish(
		events.PredictionProgressMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserID,
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

func (c *Handler) handleChannelPredictionLock(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPredictionLock,
) {
	c.logger.Info(
		"Prediction lock",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("title", event.Title),
	)

	outComes := convertOutCome(event.Outcomes)

	err := c.twirBus.Events.PredictionLock.Publish(
		events.PredictionLockMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserID,
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

func (c *Handler) handleChannelPredictionEnd(
	_ *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPredictionEnd,
) {
	c.logger.Info(
		"Prediction end",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("title", event.Title),
		slog.String("status", event.Status),
	)

	if event.Status != "resolved" {
		return
	}

	outComes := convertOutCome(event.Outcomes)

	err := c.twirBus.Events.PredictionEnd.Publish(
		events.PredictionEndMessage{
			BaseInfo: events.BaseInfo{
				ChannelID:   event.BroadcasterUserID,
				ChannelName: event.BroadcasterUserLogin,
			},
			UserName:        event.BroadcasterUserLogin,
			UserDisplayName: event.BroadcasterUserName,
			Info: events.PredictionInfo{
				Title:    event.Title,
				Outcomes: outComes,
			},
			WinningOutcomeID: event.WinningOutcomeID,
		},
	)

	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}
