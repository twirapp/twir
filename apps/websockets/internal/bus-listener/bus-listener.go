package bus_listener

import (
	"context"

	"github.com/satont/twir/apps/websockets/internal/namespaces/overlays/dudes"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/websockets"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Bus   *bus_core.Bus
	Dudes *dudes.Dudes
}

type BusListener struct {
	bus   *bus_core.Bus
	dudes *dudes.Dudes
}

func New(opts Opts) *BusListener {
	listener := &BusListener{
		bus:   opts.Bus,
		dudes: opts.Dudes,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				if err := listener.bus.WebsocketsDudesChangeColor.SubscribeGroup(
					func(ctx context.Context, data websockets.DudesChangeColorRequest) struct{} {
						listener.dudes.SendEvent(data.ChannelID, "dudes:changeColor", data)

						return struct{}{}
					},
				); err != nil {
					return err
				}
				if err := listener.bus.WebsocketsDudesGrow.SubscribeGroup(
					func(ctx context.Context, data websockets.DudesGrowRequest) struct{} {
						listener.dudes.SendEvent(data.ChannelID, "dudes:grow", data)

						return struct{}{}
					},
				); err != nil {
					return err
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				listener.bus.WebsocketsDudesChangeColor.Unsubscribe()
				listener.bus.WebsocketsDudesGrow.Unsubscribe()
				return nil
			},
		},
	)

	return listener
}
