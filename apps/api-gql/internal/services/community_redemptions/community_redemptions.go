package community_redemptions

import (
	"context"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"

	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/twitch"
)

func New(lc *lifecycle.Lifecycle, twirBus *buscore.Bus) *Service {
	s := &Service{
		twirBus: twirBus,
		subs:    make(map[string]chan twitch.ActivatedRedemption),
	}

	lc.Append(
		lifecycle.Hook{
			OnStart: func(ctx context.Context) error {
				return twirBus.RedemptionAdd.Subscribe(s.handleBusEvent)
			},
			OnStop: func(ctx context.Context) error {
				twirBus.RedemptionAdd.Unsubscribe()
				return nil
			},
		},
	)

	return s
}

type Service struct {
	twirBus *buscore.Bus

	subs map[string]chan twitch.ActivatedRedemption
}

func (s *Service) handleBusEvent(_ context.Context, data twitch.ActivatedRedemption) (
	struct{},
	error,
) {
	if ch, ok := s.subs[data.BroadcasterUserID]; ok {
		ch <- data
	}

	return struct{}{}, nil
}

func (s *Service) Subscribe(channelID string) <-chan twitch.ActivatedRedemption {
	if _, ok := s.subs[channelID]; !ok {
		s.subs[channelID] = make(chan twitch.ActivatedRedemption)
	}

	return s.subs[channelID]
}

func (s *Service) Unsubscribe(channelID string) {
	if ch, ok := s.subs[channelID]; ok {
		close(ch)
		delete(s.subs, channelID)
	}
}
