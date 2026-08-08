package di

import (
	"github.com/goforj/wire"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/tts"
	"github.com/twirapp/twir/libs/repositories/overlays_tts"
	"github.com/twirapp/twir/libs/repositories/overlays_tts/pgx"
)

var OverlaysTTSProviderSet = wire.NewSet(
	pgx.NewFx,
	wire.Bind(new(overlays_tts.Repository), new(*pgx.Pgx)),
	tts.New,
)
