package shortlinks

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	"go.uber.org/fx"
)

type topCountriesRequestDto struct {
	ShortId string `path:"shortId" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
	Limit   int    `query:"limit" default:"10" minimum:"1" maximum:"50"`
}

type countryStatsDto struct {
	Country string `json:"country" example:"US"`
	Count   uint64 `json:"count" example:"42"`
}

var _ httpbase.Route[*topCountriesRequestDto, *httpbase.BaseOutputJson[[]countryStatsDto]] = (*topCountries)(nil)

type TopCountriesOpts struct {
	fx.In

	Service *shortenedurls.Service
}

type topCountries struct {
	service *shortenedurls.Service
}

func newTopCountries(opts TopCountriesOpts) *topCountries {
	return &topCountries{
		service: opts.Service,
	}
}

func (s *topCountries) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-url-get-top-countries",
		Method:      http.MethodGet,
		Path:        "/v1/short-links/{shortId}/top-countries",
		Tags:        []string{"Short links"},
		Summary:     "Get top countries by views for short url",
	}
}

func (s *topCountries) Register(api huma.API) {
	huma.Register(
		api,
		s.GetMeta(),
		s.Handler,
	)
}

func (s *topCountries) Handler(
	ctx context.Context,
	input *topCountriesRequestDto,
) (*httpbase.BaseOutputJson[[]countryStatsDto], error) {
	limit := input.Limit
	if limit == 0 {
		limit = 10
	}

	countries, err := s.service.GetTopCountries(
		ctx,
		shortenedurls.GetTopCountriesInput{
			ShortLinkID: input.ShortId,
			Limit:       limit,
		},
	)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get top countries", err)
	}

	result := make([]countryStatsDto, len(countries))
	for i, c := range countries {
		result[i] = countryStatsDto{
			Country: c.Country,
			Count:   c.Count,
		}
	}

	return httpbase.CreateBaseOutputJson(result), nil
}
