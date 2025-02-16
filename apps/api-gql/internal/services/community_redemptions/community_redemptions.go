package community_redemptions

import (
	"context"

	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/twitch"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	TwirBus *buscore.Bus
}

func New(opts Opts) *Service {
	s := &Service{
		twirBus: opts.TwirBus,
		subs:    make(map[string]chan twitch.ActivatedRedemption),
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				return opts.TwirBus.RedemptionAdd.Subscribe(s.handleBusEvent)
			},
			OnStop: func(ctx context.Context) error {
				opts.TwirBus.RedemptionAdd.Unsubscribe()
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

func (s *Service) handleBusEvent(_ context.Context, data twitch.ActivatedRedemption) struct{} {
	if ch, ok := s.subs[data.BroadcasterUserID]; ok {
		ch <- data
	}

	return struct{}{}
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
