package shortlinks

import (
	"context"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	humahelpers "github.com/twirapp/twir/apps/api-gql/internal/server/huma_helpers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	shortlinkscustomdomains "github.com/twirapp/twir/apps/api-gql/internal/services/shortlinkscustomdomains"
	config "github.com/twirapp/twir/libs/config"
	"go.uber.org/fx"
)

type topCountriesRequestDto struct {
	ShortId string `path:"shortId" minLength:"1" pattern:"^[a-zA-Z0-9-_]+$" required:"true"`
	Limit   int    `                                                                      query:"limit" default:"10" minimum:"1" maximum:"50"`
}

type countryStatsDto struct {
	Country string `json:"country" example:"US"`
	Count   uint64 `json:"count"   example:"42"`
}

var _ httpbase.Route[*topCountriesRequestDto, *httpbase.BaseOutputJson[[]countryStatsDto]] = (*topCountries)(nil)

type TopCountriesOpts struct {
	fx.In

	Service              *shortenedurls.Service
	Sessions             *auth.Auth
	CustomDomainsService *shortlinkscustomdomains.Service
	Config               config.Config
}

type topCountries struct {
	service              *shortenedurls.Service
	sessions             *auth.Auth
	customDomainsService *shortlinkscustomdomains.Service
	config               config.Config
}

func newTopCountries(opts TopCountriesOpts) *topCountries {
	return &topCountries{
		service:              opts.Service,
		sessions:             opts.Sessions,
		customDomainsService: opts.CustomDomainsService,
		config:               opts.Config,
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

	var domain *string
	if host, err := humahelpers.GetHostFromCtx(ctx); err == nil && !isDefaultDomain(
		s.config.SiteBaseUrl,
		host,
	) {
		domain = &host
	} else if user, err := s.sessions.GetAuthenticatedUserModel(ctx); err == nil && user != nil {
		if userDomain, err := s.customDomainsService.GetByUserID(
			ctx,
			user.ID,
		); err == nil && !userDomain.IsNil() && userDomain.Verified {
			domain = &userDomain.Domain
			link, err := s.service.GetByShortID(ctx, domain, input.ShortId)
			if err == nil && link.IsNil() {
				domain = nil
			}
		}
	}

	countries, err := s.service.GetTopCountries(
		ctx,
		shortenedurls.GetTopCountriesInput{
			ShortLinkID: input.ShortId,
			Domain:      domain,
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
