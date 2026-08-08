package di

import (
	"github.com/goforj/wire"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/kappagen"
	"github.com/twirapp/twir/libs/repositories/overlays_kappagen"
	"github.com/twirapp/twir/libs/repositories/overlays_kappagen/pgx"
)

var OverlaysKappagenProviderSet = wire.NewSet(
	pgx.NewFx,
	wire.Bind(new(overlays_kappagen.Repository), new(*pgx.Pgx)),
	kappagen.New,
)
