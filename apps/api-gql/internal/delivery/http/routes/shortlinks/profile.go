package shortlinks

import (
	"context"
	"net/http"
	"net/url"
	"slices"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	config "github.com/twirapp/twir/libs/config"
	"go.uber.org/fx"
)

type profileRequestDto struct {
	Page    int    `query:"page" minimum:"0" default:"0"`
	PerPage int    `query:"perPage" minimum:"1" maximum:"100" default:"20"`
	SortBy  string `query:"sortBy" enum:"views,created_at" default:"views"`
}

type linksProfileOutputDto struct {
	Total int             `json:"total" example:"1"`
	Items []linkOutputDto `json:"items"`
}

var _ httpbase.Route[*profileRequestDto, *httpbase.BaseOutputJson[linksProfileOutputDto]] = (*profile)(nil)

type ProfileOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Config   config.Config
	Sessions *auth.Auth
}

func newProfile(opts ProfileOpts) *profile {
	return &profile{
		service:  opts.Service,
		config:   opts.Config,
		sessions: opts.Sessions,
	}
}

type profile struct {
	service  *shortenedurls.Service
	config   config.Config
	sessions *auth.Auth
}

func (p *profile) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-url-profile",
		Method:      http.MethodGet,
		Path:        "/v1/short-links",
		Tags:        []string{"Short links"},
		Summary:     "Get user's short links from authenticated user and/or from browser session",
	}
}

func (p *profile) Handler(ctx context.Context, input *profileRequestDto) (
	*httpbase.BaseOutputJson[linksProfileOutputDto],
	error,
) {
	var (
		links []linkOutputDto
		total int
	)

	user, err := p.sessions.GetAuthenticatedUserModel(ctx)
	if user != nil && err == nil {
		data, err := p.service.GetList(
			ctx, shortenedurls.GetListInput{
				Page:        input.Page,
				PerPage:     input.PerPage,
				OwnerUserID: &user.ID,
				SortBy:      input.SortBy,
			},
		)
		if err != nil {
			return nil, huma.NewError(http.StatusNotFound, "Cannot get links", err)
		}
		total = data.Total

		baseUrl, _ := url.Parse(p.config.SiteBaseUrl)

		for _, link := range data.List {
			baseUrl.Path = "/s/" + link.ID

			links = append(
				links,
				linkOutputDto{
					Id:        link.ID,
					Url:       link.Link,
					ShortUrl:  baseUrl.String(),
					Views:     link.Views,
					CreatedAt: link.CreatedAt,
				},
			)
		}
	}

	if linksIds, err := p.sessions.GetLatestShortenerUrlsIds(ctx); err == nil {
		data, err := p.service.GetManyByShortIDs(ctx, linksIds)
		if err != nil {
			return nil, huma.NewError(http.StatusNotFound, "Cannot get links", err)
		}

		for _, link := range data {
			baseUrl, _ := url.Parse(p.config.SiteBaseUrl)
			baseUrl.Path = "/s/" + link.ShortID

			links = append(
				links,
				linkOutputDto{
					Id:        link.ShortID,
					Url:       link.URL,
					ShortUrl:  baseUrl.String(),
					Views:     link.Views,
					CreatedAt: link.CreatedAt,
				},
			)
		}
	}

	// Remove duplicates
	seen := make(map[string]bool)
	uniqueLinks := []linkOutputDto{}
	for _, link := range links {
		if !seen[link.Id] {
			seen[link.Id] = true
			uniqueLinks = append(uniqueLinks, link)
		}
	}

	// perPage limit
	if len(uniqueLinks) > input.PerPage {
		uniqueLinks = uniqueLinks[:input.PerPage]
	}

	links = uniqueLinks

	// Sort based on sortBy parameter
	if input.SortBy == "views" {
		slices.SortFunc(
			links,
			func(a, b linkOutputDto) int {
				if a.Views != b.Views {
					return b.Views - a.Views // descending
				}
				return b.CreatedAt.Compare(a.CreatedAt) // tie-breaker by date
			},
		)
	} else {
		slices.SortFunc(
			links,
			func(a, b linkOutputDto) int {
				return b.CreatedAt.Compare(a.CreatedAt)
			},
		)
	}

	if total == 0 {
		total = len(links)
	}

	return httpbase.CreateBaseOutputJson(
		linksProfileOutputDto{
			Total: total,
			Items: links,
		},
	), nil
}

func (p *profile) Register(api huma.API) {
	huma.Register(api, p.GetMeta(), p.Handler)
}
