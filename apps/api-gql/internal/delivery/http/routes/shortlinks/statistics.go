package shortlinks

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
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

	Service *shortenedurls.Service
}

type statistics struct {
	service *shortenedurls.Service
}

func newStatistics(opts StatisticsOpts) *statistics {
	return &statistics{
		service: opts.Service,
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
