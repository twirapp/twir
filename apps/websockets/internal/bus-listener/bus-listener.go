package bus_listener

import (
	"context"

	dotanamespace "github.com/twirapp/twir/apps/websockets/internal/namespaces/overlays/dota"
	"github.com/twirapp/twir/apps/websockets/internal/namespaces/overlays/dudes"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
	bus_core "github.com/twirapp/twir/libs/bus-core"
	busapi "github.com/twirapp/twir/libs/bus-core/api"
	"github.com/twirapp/twir/libs/bus-core/websockets"
	"gorm.io/gorm"
)

type BusListener struct {
	bus   *bus_core.Bus
	dudes *dudes.Dudes
	dota  *dotanamespace.Dota
	gorm  *gorm.DB
}

func New(
	lc *lifecycle.Lifecycle,
	bus *bus_core.Bus,
	dudesServer *dudes.Dudes,
	dotaServer *dotanamespace.Dota,
	gorm *gorm.DB,
) *BusListener {
	listener := &BusListener{
		bus:   bus,
		dudes: dudesServer,
		dota:  dotaServer,
		gorm:  gorm,
	}

	lc.Append(
		lifecycle.Hook{
			OnStart: func(_ context.Context) error {
				if err := listener.bus.Websocket.DudesUserSettings.SubscribeGroup(
					"websockets",
					func(ctx context.Context, data websockets.DudesChangeUserSettingsRequest) (
						struct{},
						error,
					) {
						return struct{}{}, listener.dudes.SendUserSettings(ctx, data.ChannelID, data.UserID)
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

				if err := listener.bus.Api.DotaStateUpdate.SubscribeGroup(
					"websockets",
					func(ctx context.Context, data busapi.DotaStateUpdateMessage) (struct{}, error) {
						return struct{}{}, listener.dota.SendEvent(data.ChannelID, "dotaStateUpdate", data)
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
				listener.bus.Api.DotaStateUpdate.Unsubscribe()
				return nil
			},
		},
	)

	return listener
}
