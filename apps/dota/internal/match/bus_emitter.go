package match

import (
	"context"

	buscore "github.com/twirapp/twir/libs/bus-core"
	busapi "github.com/twirapp/twir/libs/bus-core/api"
	busdota "github.com/twirapp/twir/libs/bus-core/dota"
)

type BusEmitter struct {
	bus *buscore.Bus
}

var _ EventEmitter = (*BusEmitter)(nil)

func NewBusEmitter(bus *buscore.Bus) *BusEmitter {
	return &BusEmitter{bus: bus}
}

func (e *BusEmitter) MatchStarted(ctx context.Context, msg busdota.MatchStartedMessage) error {
	return e.bus.Dota.MatchStarted.Publish(ctx, msg)
}

func (e *BusEmitter) MatchEnded(ctx context.Context, msg busdota.MatchEndedMessage) error {
	return e.bus.Dota.MatchEnded.Publish(ctx, msg)
}

func (e *BusEmitter) RoshanKilled(ctx context.Context, msg busdota.RoshanKilledMessage) error {
	return e.bus.Dota.RoshanKilled.Publish(ctx, msg)
}

func (e *BusEmitter) AegisPickup(ctx context.Context, msg busdota.AegisPickupMessage) error {
	return e.bus.Dota.AegisPickup.Publish(ctx, msg)
}

func (e *BusEmitter) StateUpdate(ctx context.Context, msg busapi.DotaStateUpdateMessage) error {
	return e.bus.Api.DotaStateUpdate.Publish(ctx, msg)
}
