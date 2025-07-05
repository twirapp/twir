package handler

import (
	"context"
	"log/slog"

	"github.com/twirapp/twir/libs/bus-core/events"
	eventsub_bindings "github.com/twirapp/twitch-eventsub-framework/esb"
	"go.uber.org/zap"
)

func convertChoices(choices []eventsub_bindings.PollChoice) []events.PollChoice {
	converted := make([]events.PollChoice, 0, len(choices))

	for _, choice := range choices {
		converted = append(
			converted,
			events.PollChoice{
				ID:                  choice.ID,
				Title:               choice.Title,
				BitsVotes:           uint64(choice.BitsVotes),
				ChannelsPointsVotes: uint64(choice.ChannelPointsVotes),
				Votes:               uint64(choice.Votes),
			},
		)
	}

	return converted
}

func (c *Handler) handleChannelPollBegin(
	ctx context.Context,
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPollBegin,
) {
	c.logger.Info(
		"Poll begin",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("pollTitle", event.Title),
	)

	choices := convertChoices(event.Choices)

	msg := events.PollBeginMessage{
		BaseInfo: events.BaseInfo{
			ChannelID:   event.BroadcasterUserID,
			ChannelName: event.BroadcasterUserLogin,
		},
		UserName:        event.BroadcasterUserLogin,
		UserDisplayName: event.BroadcasterUserName,
		Info: events.PollInfo{
			Title:   event.Title,
			Choices: choices,
			ChannelsPointsVoting: events.PollChannelPointsVotes{
				Enabled:       event.ChannelPointsVoting.IsEnabled,
				AmountPerVote: uint64(event.ChannelPointsVoting.AmountPerVote),
			},
			BitsVoting: events.PollBitsVotes{
				Enabled:       event.BitsVoting.IsEnabled,
				AmountPerVote: uint64(event.BitsVoting.AmountPerVote),
			},
		},
	}

	err := c.twirBus.Events.PollBegin.Publish(ctx, msg)
	if err != nil {
		zap.S().Error(err)
	}
}

func (c *Handler) handleChannelPollProgress(
	ctx context.Context,
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPollProgress,
) {
	c.logger.Info(
		"Poll Progress",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("pollTitle", event.Title),
	)

	choices := convertChoices(event.Choices)

	msg := events.PollProgressMessage{
		BaseInfo: events.BaseInfo{
			ChannelID:   event.BroadcasterUserID,
			ChannelName: event.BroadcasterUserLogin,
		},
		UserName:        event.BroadcasterUserLogin,
		UserDisplayName: event.BroadcasterUserName,
		Info: events.PollInfo{
			Title:   event.Title,
			Choices: choices,
			ChannelsPointsVoting: events.PollChannelPointsVotes{
				Enabled:       event.ChannelPointsVoting.IsEnabled,
				AmountPerVote: uint64(event.ChannelPointsVoting.AmountPerVote),
			},
			BitsVoting: events.PollBitsVotes{
				Enabled:       event.BitsVoting.IsEnabled,
				AmountPerVote: uint64(event.BitsVoting.AmountPerVote),
			},
		},
	}

	err := c.twirBus.Events.PollProgress.Publish(ctx, msg)
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}

func (c *Handler) handleChannelPollEnd(
	ctx context.Context,
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPollEnd,
) {
	c.logger.Info(
		"Poll end",
		slog.String("channelId", event.BroadcasterUserID),
		slog.String("channelName", event.BroadcasterUserLogin),
		slog.String("pollTitle", event.Title),
	)

	choices := convertChoices(event.Choices)

	msg := events.PollEndMessage{
		BaseInfo: events.BaseInfo{
			ChannelID:   event.BroadcasterUserID,
			ChannelName: event.BroadcasterUserLogin,
		},
		UserName:        event.BroadcasterUserLogin,
		UserDisplayName: event.BroadcasterUserName,
		Info: events.PollInfo{
			Title:   event.Title,
			Choices: choices,
			ChannelsPointsVoting: events.PollChannelPointsVotes{
				Enabled:       event.ChannelPointsVoting.IsEnabled,
				AmountPerVote: uint64(event.ChannelPointsVoting.AmountPerVote),
			},
			BitsVoting: events.PollBitsVotes{
				Enabled:       event.BitsVoting.IsEnabled,
				AmountPerVote: uint64(event.BitsVoting.AmountPerVote),
			},
		},
	}

	err := c.twirBus.Events.PollEnd.Publish(ctx, msg)
	if err != nil {
		c.logger.Error(err.Error(), slog.Any("err", err))
	}
}
