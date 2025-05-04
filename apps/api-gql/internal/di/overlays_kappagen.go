package di

import (
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/kappagen"
	"github.com/twirapp/twir/libs/repositories/overlays_kappagen"
	"github.com/twirapp/twir/libs/repositories/overlays_kappagen/pgx"
	"go.uber.org/fx"
)

var OverlaysKappagenModule = fx.Options(
	// Register repository
	fx.Provide(
		fx.Annotate(
			pgx.NewFx,
			fx.As(new(overlays_kappagen.Repository)),
		),
	),

	// Register service
	fx.Provide(kappagen.New),
)
