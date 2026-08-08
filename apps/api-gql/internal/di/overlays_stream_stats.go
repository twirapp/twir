package di

import (
	"github.com/goforj/wire"
	streamstatsservice "github.com/twirapp/twir/apps/api-gql/internal/services/overlays/stream_stats"
	"github.com/twirapp/twir/libs/repositories/overlays_stream_stats"
	"github.com/twirapp/twir/libs/repositories/overlays_stream_stats/pgx"
)

var OverlaysStreamStatsProviderSet = wire.NewSet(
	pgx.NewFx,
	wire.Bind(new(overlays_stream_stats.Repository), new(*pgx.Pgx)),
	streamstatsservice.New,
)
