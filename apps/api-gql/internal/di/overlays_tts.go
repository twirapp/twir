package di

import (
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/tts"
	"github.com/twirapp/twir/libs/repositories/overlays_tts"
	"github.com/twirapp/twir/libs/repositories/overlays_tts/pgx"
	"go.uber.org/fx"
)

var OverlaysTTSModule = fx.Options(
	// Register repository
	fx.Provide(
		fx.Annotate(
			pgx.NewFx,
			fx.As(new(overlays_tts.Repository)),
		),
	),

	// Register service
	fx.Provide(tts.New),
)

