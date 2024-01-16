package handler

import (
	"context"

	"github.com/dnsge/twitch-eventsub-bindings"
	"github.com/satont/twir/libs/grpc/events"
	"go.uber.org/zap"
)

func convertChoices(choices []eventsub_bindings.PollChoice) []*events.PollInfo_Choice {
	converted := make([]*events.PollInfo_Choice, 0, len(choices))

	for _, choice := range choices {
		converted = append(
			converted, &events.PollInfo_Choice{
				Id:                  choice.ID,
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
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPollBegin,
) {
	zap.S().Infow(
		"Poll begin",
		"channelId", event.BroadcasterUserID,
		"channelName", event.BroadcasterUserLogin,
		"pollTitle", event.Title,
	)

	choices := convertChoices(event.Choices)

	msg := &events.PollBeginMessage{
		BaseInfo:        &events.BaseInfo{ChannelId: event.BroadcasterUserID},
		UserName:        event.BroadcasterUserLogin,
		UserDisplayName: event.BroadcasterUserName,
		Info: &events.PollInfo{
			Title:   event.Title,
			Choices: choices,
			ChannelsPointsVoting: &events.PollInfo_ChannelPointsVotes{
				Enabled:       event.ChannelPointsVoting.IsEnabled,
				AmountPerVote: uint64(event.ChannelPointsVoting.AmountPerVote),
			},
			BitsVoting: &events.PollInfo_BitsVotes{
				Enabled:       event.BitsVoting.IsEnabled,
				AmountPerVote: uint64(event.BitsVoting.AmountPerVote),
			},
		},
	}

	_, err := c.services.Grpc.Events.PollBegin(context.Background(), msg)
	if err != nil {
		zap.S().Error(err)
	}
}

func (c *Handler) handleChannelPollProgress(
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPollProgress,
) {
	zap.S().Infow(
		"Poll Progress",
		"channelId", event.BroadcasterUserID,
		"channelName", event.BroadcasterUserLogin,
		"pollTitle", event.Title,
	)

	choices := convertChoices(event.Choices)

	msg := &events.PollProgressMessage{
		BaseInfo:        &events.BaseInfo{ChannelId: event.BroadcasterUserID},
		UserName:        event.BroadcasterUserLogin,
		UserDisplayName: event.BroadcasterUserName,
		Info: &events.PollInfo{
			Title:   event.Title,
			Choices: choices,
			ChannelsPointsVoting: &events.PollInfo_ChannelPointsVotes{
				Enabled:       event.ChannelPointsVoting.IsEnabled,
				AmountPerVote: uint64(event.ChannelPointsVoting.AmountPerVote),
			},
			BitsVoting: &events.PollInfo_BitsVotes{
				Enabled:       event.BitsVoting.IsEnabled,
				AmountPerVote: uint64(event.BitsVoting.AmountPerVote),
			},
		},
	}

	_, err := c.services.Grpc.Events.PollProgress(context.Background(), msg)
	if err != nil {
		zap.S().Error(err)
	}
}

func (c *Handler) handleChannelPollEnd(
	h *eventsub_bindings.ResponseHeaders,
	event *eventsub_bindings.EventChannelPollEnd,
) {
	zap.S().Infow(
		"Poll end",
		"channelId", event.BroadcasterUserID,
		"channelName", event.BroadcasterUserLogin,
		"pollTitle", event.Title,
	)

	choices := convertChoices(event.Choices)

	msg := &events.PollEndMessage{
		BaseInfo:        &events.BaseInfo{ChannelId: event.BroadcasterUserID},
		UserName:        event.BroadcasterUserLogin,
		UserDisplayName: event.BroadcasterUserName,
		Info: &events.PollInfo{
			Title:   event.Title,
			Choices: choices,
			ChannelsPointsVoting: &events.PollInfo_ChannelPointsVotes{
				Enabled:       event.ChannelPointsVoting.IsEnabled,
				AmountPerVote: uint64(event.ChannelPointsVoting.AmountPerVote),
			},
			BitsVoting: &events.PollInfo_BitsVotes{
				Enabled:       event.BitsVoting.IsEnabled,
				AmountPerVote: uint64(event.BitsVoting.AmountPerVote),
			},
		},
	}

	_, err := c.services.Grpc.Events.PollEnd(context.Background(), msg)
	if err != nil {
		zap.S().Error(err)
	}
}
