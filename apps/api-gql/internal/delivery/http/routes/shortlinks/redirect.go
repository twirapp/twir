package shortlinks

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	humahelpers "github.com/twirapp/twir/apps/api-gql/internal/server/huma_helpers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	config "github.com/twirapp/twir/libs/config"
	shortlinksviewsrepository "github.com/twirapp/twir/libs/repositories/short_links_views"
	"github.com/twirapp/twir/libs/repositories/shortened_urls/model"
	"go.uber.org/fx"
)

type redirectRequestDto struct {
	ShortId string `path:"shortId" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
}

type redirectResponseDto struct {
	Status   int
	Location string `header:"Location"`
}

var _ httpbase.Route[*redirectRequestDto, *redirectResponseDto] = (*redirect)(nil)

type RedirectOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Config   config.Config
	Sessions *auth.Auth
	Logger   *slog.Logger
}

func newRedirect(opts RedirectOpts) *redirect {
	return &redirect{
		service:  opts.Service,
		config:   opts.Config,
		sessions: opts.Sessions,
		logger:   opts.Logger,
	}
}

type redirect struct {
	service  *shortenedurls.Service
	config   config.Config
	sessions *auth.Auth
	logger   *slog.Logger
}

func (r *redirect) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID:   "short-url-redirect",
		Method:        http.MethodGet,
		Path:          "/v1/short-links/{shortId}",
		Tags:          []string{"Short links"},
		Summary:       "Redirect to url",
		DefaultStatus: 301,
	}
}

func (r *redirect) Handler(ctx context.Context, input *redirectRequestDto) (
	*redirectResponseDto,
	error,
) {
	link, err := r.service.GetByShortID(ctx, input.ShortId)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
	}

	if link == model.Nil {
		return nil, huma.NewError(http.StatusNotFound, "Link not found")
	}

	var userID *string
	user, _ := r.sessions.GetAuthenticatedUserModel(ctx)
	if user != nil {
		userID = &user.ID
	}

	var clientIp *string
	if ip, err := humahelpers.GetClientIpFromCtx(ctx); err == nil {
		clientIp = &ip
	} else {
		r.logger.Warn("Cannot get client IP", "error", err)
	}

	var clientAgent *string
	if agent, err := humahelpers.GetClientUserAgentFromCtx(ctx); err == nil {
		clientAgent = &agent
	} else {
		r.logger.Warn("Cannot get client user agent", "error", err)
	}

	var country *string
	if cfCountry, err := humahelpers.GetCloudflareCountryFromCtx(ctx); err == nil && cfCountry != "" {
		country = &cfCountry
	}

	var city *string
	if cfCity, err := humahelpers.GetCloudflareCityFromCtx(ctx); err == nil && cfCity != "" {
		city = &cfCity
	}

	if err := r.service.RecordView(
		ctx,
		shortenedurls.RecordViewInput{
			ShortLinkID: link.ShortID,
			UserID:      userID,
			IP:          clientIp,
			UserAgent:   clientAgent,
			Country:     country,
			City:        city,
		},
	); err != nil {
		r.logger.Warn("Cannot record view", "error", err)
	}

	newViews := link.Views + 1

	if err := r.service.Update(
		ctx,
		link.ShortID,
		shortenedurls.UpdateInput{
			Views: &newViews,
		},
	); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot update link", err)
	}

	// Publish update to subscribers
	go func() {
		lastView := &shortlinksviewsrepository.View{
			ShortLinkID: link.ShortID,
			UserID:      userID,
			Country:     country,
			City:        city,
			CreatedAt:   time.Now(),
		}
		_ = r.service.PublishViewUpdate(link.ShortID, newViews, lastView)
	}()

	return &redirectResponseDto{
		Status:   http.StatusPermanentRedirect,
		Location: link.URL,
	}, nil
}

func (r *redirect) Register(api huma.API) {
	huma.Register(
		api,
		r.GetMeta(),
		r.Handler,
	)
}
