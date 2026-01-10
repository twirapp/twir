package twir

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	twir_stats "github.com/twirapp/twir/apps/api-gql/internal/delivery/gql/twir-stats"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"go.uber.org/fx"
)

var FxModule = fx.Provide(
	httpbase.AsFxRoute(newStats),
)

type twirStatsRequestDto struct{}

type twirStatsResponseBody struct {
	Channels        int `json:"channels"`
	CreatedCommands int `json:"created_commands"`
	Viewers         int `json:"viewers"`
	Messages        int `json:"messages"`
	UsedEmotes      int `json:"used_emotes"`
	UsedCommands    int `json:"used_commands"`
	ShortUrls       int `json:"short_urls"`
	HasteBins       int `json:"haste_bins"`
	LiveChannels    int `json:"live_channels"`
}

type twirStatsResponseDto struct {
	Body twirStatsResponseBody
}

var _ httpbase.Route[*twirStatsRequestDto, *twirStatsResponseDto] = (*twirStats)(nil)

type StatsOpts struct {
	fx.In

	Service *twir_stats.TwirStats
}

func newStats(opts StatsOpts) *twirStats {
	return &twirStats{
		service: opts.Service,
	}
}

type twirStats struct {
	service *twir_stats.TwirStats
}

func (s *twirStats) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "twir-stats",
		Method:      http.MethodGet,
		Path:        "/v1/twir/stats",
		Tags:        []string{"Twir"},
		Summary:     "Twir Stats",
		Description: "Get Twir application statistics",
	}
}

func (s *twirStats) Handler(
	ctx context.Context,
	input *twirStatsRequestDto,
) (*twirStatsResponseDto, error) {
	cachedData := s.service.GetCachedData()

	return &twirStatsResponseDto{
		Body: twirStatsResponseBody{
			Channels:        cachedData.Channels,
			CreatedCommands: cachedData.CreatedCommands,
			Viewers:         cachedData.Viewers,
			Messages:        cachedData.Messages,
			UsedEmotes:      cachedData.UsedEmotes,
			UsedCommands:    cachedData.UsedCommands,
			ShortUrls:       cachedData.ShortUrls,
			HasteBins:       cachedData.HasteBins,
			LiveChannels:    cachedData.LiveChannels,
		},
	}, nil
}

func (s *twirStats) Register(api huma.API) {
	huma.Register(api, s.GetMeta(), s.Handler)
}
