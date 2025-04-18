package bus_listener

import (
	"context"
	"log/slog"

	"github.com/satont/twir/libs/logger"
	"github.com/twirapp/twir/apps/giveaways/internal/services"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/giveaways"
	"go.uber.org/fx"
)

type giveawaysListener struct {
	giveawaysService *services.Service
	bus              *buscore.Bus
	logger           logger.Logger
}

type Opts struct {
	fx.In
	Lc fx.Lifecycle

	Logger           logger.Logger
	GiveawaysService *services.Service
	Bus              *buscore.Bus
}

func New(opts Opts) error {
	impl := &giveawaysListener{
		giveawaysService: opts.GiveawaysService,
		bus:              opts.Bus,
		logger:           opts.Logger,
	}

	opts.Lc.Append(
		fx.Hook{
			OnStart: func(ctx context.Context) error {
				impl.bus.Giveaways.ChooseWinner.SubscribeGroup(
					"giveaways",
					impl.chooseWinner,
				)

				impl.bus.Giveaways.TryAddParticipant.SubscribeGroup(
					"giveaways",
					impl.tryAddParticipant,
				)

				return nil
			},
			OnStop: func(ctx context.Context) error {
				impl.bus.Giveaways.ChooseWinner.Unsubscribe()
				impl.bus.Giveaways.TryAddParticipant.Unsubscribe()

				return nil
			},
		},
	)

	return nil
}

func (c *giveawaysListener) tryAddParticipant(
	ctx context.Context,
	req giveaways.TryAddParticipantRequest,
) struct{} {
	if err := c.giveawaysService.TryAddParticipant(
		ctx,
		req.UserID,
		req.UserLogin,
		req.UserDisplayName,
		req.GiveawayID,
	); err != nil {
		c.logger.Error("failed to add participant to giveaways", slog.Any("err", err))
	}

	return struct{}{}
}

func (c *giveawaysListener) chooseWinner(
	ctx context.Context,
	req giveaways.ChooseWinnerRequest,
) giveaways.ChooseWinnerResponse {
	winners, err := c.giveawaysService.ChooseWinner(ctx, req.GiveawayID)
	if err != nil {
		return giveaways.ChooseWinnerResponse{}
	}

	mappedWinners := make([]giveaways.Winner, 0, len(winners))
	for _, winner := range winners {
		mappedWinners = append(
			mappedWinners,
			giveaways.Winner{
				UserID:          winner.UserID,
				UserLogin:       winner.UserLogin,
				UserDisplayName: winner.UserDisplayName,
			},
		)
	}

	return giveaways.ChooseWinnerResponse{
		Winners: mappedWinners,
	}
}
