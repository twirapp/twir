package di

import (
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/be_right_back"
	"github.com/twirapp/twir/libs/repositories/overlays_be_right_back"
	"github.com/twirapp/twir/libs/repositories/overlays_be_right_back/pgx"
	"go.uber.org/fx"
)

var OverlaysBeRightBackModule = fx.Options(
	// Register repository
	fx.Provide(
		fx.Annotate(
			pgx.NewFx,
			fx.As(new(overlays_be_right_back.Repository)),
		),
	),

	// Register service
	fx.Provide(be_right_back.New),
)
