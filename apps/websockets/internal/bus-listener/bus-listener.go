package bus_listener

import (
	"context"

	"github.com/satont/twir/apps/websockets/internal/namespaces/overlays/dudes"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/websockets"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In
	LC fx.Lifecycle

	Bus   *bus_core.Bus
	Dudes *dudes.Dudes
	Gorm  *gorm.DB
}

type BusListener struct {
	bus   *bus_core.Bus
	dudes *dudes.Dudes
	gorm  *gorm.DB
}

func New(opts Opts) *BusListener {
	listener := &BusListener{
		bus:   opts.Bus,
		dudes: opts.Dudes,
		gorm:  opts.Gorm,
	}

	opts.LC.Append(
		fx.Hook{
			OnStart: func(_ context.Context) error {
				if err := listener.bus.Websocket.DudesUserSettings.SubscribeGroup(
					"websockets",
					func(ctx context.Context, data websockets.DudesChangeUserSettingsRequest) (
						struct{},
						error,
					) {
						return struct{}{}, listener.dudes.SendUserSettings(data.ChannelID, data.UserID)
					},
				); err != nil {
					return err
				}
				if err := listener.bus.Websocket.DudesGrow.SubscribeGroup(
					"websockets",
					func(ctx context.Context, data websockets.DudesGrowRequest) (struct{}, error) {
						return struct{}{}, listener.dudes.SendEvent(data.ChannelID, "grow", data)
					},
				); err != nil {
					return err
				}

				if err := listener.bus.Websocket.DudesLeave.SubscribeGroup(
					"websockets",
					func(ctx context.Context, data websockets.DudesLeaveRequest) (struct{}, error) {
						return struct{}{}, listener.dudes.SendEvent(data.ChannelID, "leave", data)
					},
				); err != nil {
					return err
				}

				return nil
			},
			OnStop: func(ctx context.Context) error {
				listener.bus.Websocket.DudesUserSettings.Unsubscribe()
				listener.bus.Websocket.DudesGrow.Unsubscribe()
				listener.bus.Websocket.DudesLeave.Unsubscribe()
				return nil
			},
		},
	)

	return listener
}
