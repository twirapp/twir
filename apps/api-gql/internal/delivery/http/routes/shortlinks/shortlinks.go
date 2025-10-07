package shortlinks

import (
	"context"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twirapp/twir/libs/repositories/shortened_urls/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Api      huma.API
	Config   config.Config
	Service  *shortenedurls.Service
	Sessions *auth.Auth
	Logger   logger.Logger
}

func New(opts Opts) {
	huma.Register(
		opts.Api,
		huma.Operation{
			OperationID: "short-url-create",
			Method:      http.MethodPost,
			Path:        "/v1/short-links",
			Tags:        []string{"Short links"},
			Summary:     "Create short url",
		}, func(
			ctx context.Context,
			input *createLinkInput,
		) (
			*createLinkOutput, error,
		) {
			if input.Body.Alias == "" {
				existedLink, err := opts.Service.GetByUrl(ctx, input.Body.Url)
				if err != nil {
					return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
				}

				if existedLink != model.Nil {
					baseUrl, _ := url.Parse(opts.Config.SiteBaseUrl)
					baseUrl.Path = "/s/" + existedLink.ShortID

					return &createLinkOutput{
						Body: linkOutputDto{
							Id:       existedLink.ShortID,
							Url:      existedLink.URL,
							ShortUrl: baseUrl.String(),
							Views:    existedLink.Views,
						},
					}, nil
				}
			}

			if input.Body.Alias != "" {
				existedLink, err := opts.Service.GetByShortID(ctx, input.Body.Alias)
				if err != nil {
					return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
				}

				if existedLink != model.Nil {
					return nil, huma.NewError(http.StatusConflict, "Alias already in use")
				}
			}

			var createdByUserID *string
			user, _ := opts.Sessions.GetAuthenticatedUserModel(ctx)
			if user != nil {
				createdByUserID = &user.ID
			}

			link, err := opts.Service.Create(
				ctx, shortenedurls.CreateInput{
					URL:             input.Body.Url,
					ShortID:         input.Body.Alias,
					CreatedByUserID: createdByUserID,
				},
			)
			if err != nil {
				return nil, huma.NewError(http.StatusNotFound, "Cannot generate short id", err)
			}

			baseUrl, _ := url.Parse(opts.Config.SiteBaseUrl)
			baseUrl.Path = "/s/" + link.ShortID

			if err := opts.Sessions.AddLatestShortenerUrlsId(ctx, link.ShortID); err != nil {
				opts.Logger.Warn("Cannot save latest short links ids to session: " + err.Error())
			}

			return &createLinkOutput{
				Body: linkOutputDto{
					Id:       link.ShortID,
					Url:      input.Body.Url,
					ShortUrl: baseUrl.String(),
					Views:    link.Views,
				},
			}, nil
		},
	)

	huma.Register(
		opts.Api,
		huma.Operation{
			OperationID: "short-url-get-info",
			Method:      http.MethodGet,
			Path:        "/v1/short-links/{shortId}/info",
			Tags:        []string{"Short links"},
			Summary:     "Get short url data",
		}, func(
			ctx context.Context,
			input *struct {
				ShortId string `query:"shortId" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
			},
		) (
			*createLinkOutput, error,
		) {
			link, err := opts.Service.GetByShortID(ctx, input.ShortId)
			if err != nil {
				return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
			}

			if link == model.Nil {
				return nil, huma.NewError(http.StatusNotFound, "Link not found")
			}

			baseUrl, _ := url.Parse(opts.Config.SiteBaseUrl)
			baseUrl.Path = "/s/" + input.ShortId

			return &createLinkOutput{
				Body: linkOutputDto{
					Id:       link.ShortID,
					Url:      link.URL,
					ShortUrl: baseUrl.String(),
					Views:    link.Views,
				},
			}, nil
		},
	)

	huma.Register(
		opts.Api,
		huma.Operation{
			OperationID:   "short-url-redirect",
			Method:        http.MethodGet,
			Path:          "/v1/short-links/{shortId}",
			Tags:          []string{"Short links"},
			Summary:       "Redirect to url",
			DefaultStatus: 301,
		},
		func(
			ctx context.Context, input *struct {
				ShortId string `path:"shortId" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
			},
		) (
			*linkRedirectOutput, error,
		) {
			link, err := opts.Service.GetByShortID(ctx, input.ShortId)
			if err != nil {
				return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
			}

			if link == model.Nil {
				return nil, huma.NewError(http.StatusNotFound, "Link not found")
			}

			newViews := link.Views + 1

			if err := opts.Service.Update(
				ctx,
				link.ShortID,
				shortenedurls.UpdateInput{
					Views: &newViews,
				},
			); err != nil {
				return nil, huma.NewError(http.StatusInternalServerError, "Cannot update link", err)
			}

			return &linkRedirectOutput{
				Status:   http.StatusPermanentRedirect,
				Location: link.URL,
			}, nil
		},
	)

	huma.Register(
		opts.Api,
		huma.Operation{
			OperationID: "short-url-profile",
			Method:      http.MethodGet,
			Path:        "/v1/short-links",
			Tags:        []string{"Short links"},
			Summary:     "Get user's short links from authenticated user and/or from browser session",
		}, func(
			ctx context.Context,
			input *struct {
				Page    int `query:"page" minimum:"0" default:"0"`
				PerPage int `query:"perPage" minimum:"1" maximum:"100" default:"20"`
			},
		) (
			*struct {
				Body linksProfileOutputDto
			}, error,
		) {
			var (
				links []linkOutputDto
				total int
			)

			user, err := opts.Sessions.GetAuthenticatedUserModel(ctx)
			if user != nil && err == nil {
				data, err := opts.Service.GetList(
					ctx, shortenedurls.GetListInput{
						Page:        input.Page,
						PerPage:     input.PerPage,
						OwnerUserID: &user.ID,
					},
				)
				if err != nil {
					return nil, huma.NewError(http.StatusNotFound, "Cannot get links", err)
				}
				total = data.Total

				baseUrl, _ := url.Parse(opts.Config.SiteBaseUrl)

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

			if linksIds, err := opts.Sessions.GetLatestShortenerUrlsIds(ctx); err == nil {
				data, err := opts.Service.GetManyByShortIDs(ctx, linksIds)
				if err != nil {
					return nil, huma.NewError(http.StatusNotFound, "Cannot get links", err)
				}

				for _, link := range data {
					baseUrl, _ := url.Parse(opts.Config.SiteBaseUrl)
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

			slices.SortFunc(
				links,
				func(a, b linkOutputDto) int {
					return b.CreatedAt.Compare(a.CreatedAt)
				},
			)

			if total == 0 {
				total = len(links)
			}

			return &struct {
				Body linksProfileOutputDto
			}{
				Body: linksProfileOutputDto{
					Total: total,
					Items: links,
				},
			}, nil
		},
	)
}

type createLinkInput struct {
	Body createLinkInputDto
}

type createLinkInputDto struct {
	Url   string `json:"url" required:"true" format:"uri" minLength:"1" maxLength:"2000" example:"https://example.com" pattern:"^https?://.*"`
	Alias string `json:"alias" required:"false" minLength:"3" maxLength:"30" example:"stream" pattern:"^[a-zA-Z0-9]+$"`
}

type createLinkOutput struct {
	Body linkOutputDto
}

type linkOutputDto struct {
	Id        string    `json:"id" example:"KKMEa"`
	Url       string    `json:"url" example:"https://example.com"`
	ShortUrl  string    `json:"short_url" example:"https://twir.app/s/KKMEa"`
	Views     int       `json:"views" example:"1"`
	CreatedAt time.Time `json:"created_at" format:"date-time" example:"2023-01-01T00:00:00Z"`
}

type linkRedirectOutput struct {
	Status   int
	Location string `header:"Location"`
}

type linksProfileOutputDto struct {
	Total int             `json:"total" example:"1"`
	Items []linkOutputDto `json:"items"`
}
