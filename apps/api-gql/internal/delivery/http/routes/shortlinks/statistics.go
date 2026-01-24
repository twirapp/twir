package shortlinks

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	humahelpers "github.com/twirapp/twir/apps/api-gql/internal/server/huma_helpers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	shortlinkscustomdomains "github.com/twirapp/twir/apps/api-gql/internal/services/shortlinkscustomdomains"
	"go.uber.org/fx"
)

type statisticsRequestDto struct {
	ShortId  string `path:"shortId" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
	From     int64  `query:"from" required:"true"`
	To       int64  `query:"to" required:"true"`
	Interval string `query:"interval" enum:"hour,day" default:"day"`
}

type statisticsPointDto struct {
	Timestamp int64 `json:"timestamp"`
	Count     int64 `json:"count"`
}

var _ httpbase.Route[*statisticsRequestDto, *httpbase.BaseOutputJson[[]statisticsPointDto]] = (*statistics)(nil)

type StatisticsOpts struct {
	fx.In

	Service              *shortenedurls.Service
	Sessions             *auth.Auth
	CustomDomainsService *shortlinkscustomdomains.Service
}

type statistics struct {
	service              *shortenedurls.Service
	sessions             *auth.Auth
	customDomainsService *shortlinkscustomdomains.Service
}

func newStatistics(opts StatisticsOpts) *statistics {
	return &statistics{
		service:              opts.Service,
		sessions:             opts.Sessions,
		customDomainsService: opts.CustomDomainsService,
	}
}

func (s *statistics) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-url-get-statistics",
		Method:      http.MethodGet,
		Path:        "/v1/short-links/{shortId}/statistics",
		Tags:        []string{"Short links"},
		Summary:     "Get short url statistics",
	}
}

func (s *statistics) Register(api huma.API) {
	huma.Register(
		api,
		s.GetMeta(),
		s.Handler,
	)
}

func (s *statistics) Handler(
	ctx context.Context,
	input *statisticsRequestDto,
) (*httpbase.BaseOutputJson[[]statisticsPointDto], error) {
	var domain *string
	if host, err := humahelpers.GetHostFromCtx(ctx); err == nil && !isDefaultDomain(host) {
		domain = &host
	} else if user, err := s.sessions.GetAuthenticatedUserModel(ctx); err == nil && user != nil {
		if userDomain, err := s.customDomainsService.GetByUserID(ctx, user.ID); err == nil && !userDomain.IsNil() && userDomain.Verified {
			domain = &userDomain.Domain
			link, err := s.service.GetByShortID(ctx, domain, input.ShortId)
			if err == nil && link.IsNil() {
				domain = nil
			}
		}
	}

	from := time.Unix(input.From/1000, (input.From%1000)*1000000)
	to := time.Unix(input.To/1000, (input.To%1000)*1000000)

	interval := input.Interval
	if interval == "" {
		interval = "day"
	}

	points, err := s.service.GetStatistics(
		ctx,
		shortenedurls.GetStatisticsInput{
			ShortLinkID: input.ShortId,
			Domain:      domain,
			From:        from,
			To:          to,
			Interval:    interval,
		},
	)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get statistics", err)
	}

	result := make([]statisticsPointDto, len(points))
	for i, p := range points {
		result[i] = statisticsPointDto{
			Timestamp: p.Timestamp,
			Count:     p.Count,
		}
	}

	return httpbase.CreateBaseOutputJson(result), nil
}
