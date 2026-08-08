package dota

import (
	"context"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	busapi "github.com/twirapp/twir/libs/bus-core/api"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
	"github.com/twirapp/twir/libs/logger"
)

func createStateSubscriptionKey(channelID uuid.UUID) string {
	return fmt.Sprintf("dota:state_update:%s", channelID)
}

// StateSubscriptionByApiKey streams dota state updates for the overlay pages.
// The current state is emitted immediately so overlays render without waiting
// for the next match event.
func (s *Service) StateSubscriptionByApiKey(
	ctx context.Context,
	apiKey string,
) (<-chan busapi.DotaStateUpdateMessage, error) {
	identity, err := s.channels.ResolveApiKeyChannelIdentityByUserOrChannelApiKey(ctx, apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve channel by api key: %w", err)
	}

	channelID, err := uuid.Parse(identity.InternalChannelID)
	if err != nil {
		return nil, fmt.Errorf("invalid channel id: %w", err)
	}

	wsRouterSub, err := s.wsRouter.Subscribe([]string{createStateSubscriptionKey(channelID)})
	if err != nil {
		return nil, err
	}

	outputChan := make(chan busapi.DotaStateUpdateMessage, 1)

	if response, err := s.twirBus.Dota.GetData.Request(
		ctx,
		busdota.GetDataRequest{ChannelID: channelID.String()},
	); err == nil && response != nil {
		data := response.Data
		outputChan <- busapi.DotaStateUpdateMessage{
			ChannelID:      channelID.String(),
			InGame:         data.InGame,
			Mmr:            data.Mmr,
			SessionWins:    data.SessionWins,
			SessionLosses:  data.SessionLosses,
			WinProbability: data.WinProbability,
			HeroName:       data.HeroName,
			MatchID:        data.MatchID,
			TeamIsRadiant:  data.TeamIsRadiant,
			TeamKnown:      data.TeamKnown,
		}
	}

	go func() {
		defer func() {
			wsRouterSub.Unsubscribe()
			close(outputChan)
		}()

		for {
			select {
			case <-ctx.Done():
				return
			case data, ok := <-wsRouterSub.GetChannel():
				if !ok {
					return
				}

				var msg busapi.DotaStateUpdateMessage
				if err := json.Unmarshal(data, &msg); err != nil {
					s.logger.ErrorContext(ctx, "failed to unmarshal dota state update", logger.Error(err))
					continue
				}

				select {
				case outputChan <- msg:
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	return outputChan, nil
}
