package di

import (
	"github.com/goforj/wire"

	dotaservice "github.com/twirapp/twir/apps/api-gql/internal/services/dota"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/integrations/steam"
	dotarepository "github.com/twirapp/twir/libs/repositories/dota"
	dotarepositorypgx "github.com/twirapp/twir/libs/repositories/dota/pgx"
)

var DotaProviderSet = wire.NewSet(
	dotarepositorypgx.NewFx,
	wire.Bind(new(dotarepository.Repository), new(*dotarepositorypgx.Pgx)),
	NewSteamClient,
	dotaservice.New,
)

func NewSteamClient(config cfg.Config) *steam.Client {
	return steam.New(config.DotaSteamAPIKey)
}
