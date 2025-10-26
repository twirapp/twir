package valorant

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/goccy/go-json"
	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	valorantintegration "github.com/twirapp/twir/apps/api-gql/internal/services/valorant_integration"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/integrations/valorant"
	"go.uber.org/fx"
	"golang.org/x/sync/errgroup"
)

type Opts struct {
	fx.In

	Api      huma.API
	Config   config.Config
	Sessions *auth.Auth
	Service  *valorantintegration.Service
	KV       kv.KV
}

func New(opts Opts) {
	huma.Register(
		opts.Api,
		huma.Operation{
			OperationID: "integrations-valorant-stats",
			Summary:     "Get valorant stats data",
			Description: "Requires api-key header.",
			Method:      http.MethodGet,
			Tags:        []string{"Valorant"},
			Path:        "/v1/integrations/valorant/stats",
			Security: []map[string][]string{
				{"api-key": {}},
			},
		},
		func(
			ctx context.Context,
			input *struct{},
		) (*httpdelivery.BaseOutputJson[integrationsValorantStatsOutput], error) {
			user, err := opts.Sessions.GetAuthenticatedUserModel(ctx)
			if user == nil || err != nil {
				return nil, huma.NewError(http.StatusUnauthorized, "Not authenticated", err)
			}

			selectedDashboardId, err := opts.Sessions.GetSelectedDashboard(ctx)
			if err != nil {
				return nil, huma.NewError(http.StatusUnauthorized, "Not authenticated", err)
			}

			var output *integrationsValorantStatsOutput
			if cachedBytes, _ := opts.KV.Get(
				ctx,
				"valorant_stats_"+selectedDashboardId,
			).Bytes(); cachedBytes != nil {
				if err := json.Unmarshal(cachedBytes, &output); err != nil {
					return nil, huma.NewError(
						http.StatusInternalServerError,
						"Failed to get cached valorant stats",
						err,
					)
				}

				return httpdelivery.CreateBaseOutputJson(*output), nil
			}

			wg, wgCtx := errgroup.WithContext(ctx)
			var (
				matches []valorant.StoredMatchesResponseMatch
				mmr     *valorant.MmrResponseData
			)

			wg.Go(
				func() error {
					m, err := opts.Service.GetChannelStoredMatchesByChannelID(ctx, selectedDashboardId)
					if err != nil {
						return err
					}
					matches = m.Data
					return nil
				},
			)

			wg.Go(
				func() error {
					m, err := opts.Service.GetChannelMmr(wgCtx, selectedDashboardId)
					if err != nil {
						return err
					}
					mmr = m.Data
					return nil
				},
			)

			if err := wg.Wait(); err != nil {
				return nil, err
			}

			output = &integrationsValorantStatsOutput{
				Matches: matches,
				MMR:     mmr,
			}

			bytes, err := json.Marshal(output)
			if err == nil {
				if err := opts.KV.Set(
					ctx,
					"valorant_stats_"+selectedDashboardId,
					bytes,
					kvoptions.WithExpire(10*time.Second),
				); err != nil {
					return nil, huma.NewError(
						http.StatusInternalServerError,
						"Failed to cache valorant stats",
						err,
					)
				}
			}

			return httpdelivery.CreateBaseOutputJson(*output), nil
		},
	)
}

type integrationsValorantStatsOutput struct {
	Matches []valorant.StoredMatchesResponseMatch `json:"matches"`
	MMR     *valorant.MmrResponseData             `json:"mmr"`
}
