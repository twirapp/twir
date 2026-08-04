package di

import (
	"github.com/goforj/wire"
	"github.com/twirapp/twir/apps/api-gql/internal/services/overlays/be_right_back"
	"github.com/twirapp/twir/libs/repositories/overlays_be_right_back"
	"github.com/twirapp/twir/libs/repositories/overlays_be_right_back/pgx"
)

var OverlaysBeRightBackProviderSet = wire.NewSet(
	pgx.NewFx,
	wire.Bind(new(overlays_be_right_back.Repository), new(*pgx.Pgx)),
	wire.Struct(new(be_right_back.Opts), "*"),
	be_right_back.New,
)
