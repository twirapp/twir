package shortlinks

import (
	"context"
	"log/slog"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	humahelpers "github.com/twirapp/twir/apps/api-gql/internal/server/huma_helpers"
	"github.com/twirapp/twir/apps/api-gql/internal/services/clientinfo"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	config "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/logger"
	shortlinksviewsrepository "github.com/twirapp/twir/libs/repositories/short_links_views"
	"go.uber.org/fx"
)

type redirectRequestDto struct {
	ShortID string `path:"shortId" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
}

type redirectResponseDto struct {
	Status   int
	Location string `header:"Location"`
}

var _ httpbase.Route[*redirectRequestDto, *redirectResponseDto] = (*redirect)(nil)

type RedirectOpts struct {
	fx.In

	Service           *shortenedurls.Service
	Config            config.Config
	Sessions          *auth.Auth
	Logger            *slog.Logger
	ClientInfoService *clientinfo.Service
}

func newRedirect(opts RedirectOpts) *redirect {
	return &redirect{
		service:           opts.Service,
		config:            opts.Config,
		sessions:          opts.Sessions,
		logger:            opts.Logger,
		clientInfoService: opts.ClientInfoService,
	}
}

type redirect struct {
	service           *shortenedurls.Service
	config            config.Config
	sessions          *auth.Auth
	logger            *slog.Logger
	clientInfoService *clientinfo.Service
}

func (r *redirect) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID:   "short-url-redirect",
		Method:        http.MethodGet,
		Path:          "/v1/short-links/{shortId}",
		Tags:          []string{"Short links"},
		Summary:       "Redirect to url",
		DefaultStatus: http.StatusMovedPermanently,
	}
}

func (r *redirect) Handler(ctx context.Context, input *redirectRequestDto) (
	*redirectResponseDto,
	error,
) {
	host, err := humahelpers.GetHostFromCtx(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get host", err)
	}

	var domain *string
	if !isDefaultDomain(r.config.SiteBaseUrl, host) {
		domain = &host
	}

	link, err := r.service.GetByShortID(ctx, domain, input.ShortID)
	if err != nil {
		return nil, huma.NewError(http.StatusNotFound, "Cannot get link", err)
	}

	if link.IsNil() {
		return nil, huma.NewError(http.StatusNotFound, "Link not found")
	}

	var userID *string
	user, _ := r.sessions.GetAuthenticatedUserModel(ctx)
	if user != nil {
		userID = &user.ID
	}

	var clientIP, clientUserAgent *string
	if info, err := r.clientInfoService.GetClientInfo(ctx); err == nil {
		clientIP = &info.IP
		clientUserAgent = &info.UserAgent
	} else {
		r.logger.Warn("Cannot get client info", "error", err)
	}

	var country *string
	var city *string
	if clientIP != nil {
		location, err := r.clientInfoService.LookupIP(ctx, *clientIP)
		if err != nil {
			r.logger.WarnContext(ctx, "Cannot resolve GeoIP location", "error", err)
		} else {
			country = location.Country
			city = location.City
		}
	}

	if err := r.service.RecordView(
		ctx,
		shortenedurls.RecordViewInput{
			ShortLinkID: link.ShortID,
			Domain:      domain,
			UserID:      userID,
			IP:          clientIP,
			UserAgent:   clientUserAgent,
			Country:     country,
			City:        city,
		},
	); err != nil {
		r.logger.WarnContext(ctx, "Cannot record view", logger.Error(err))
	}

	newViews := link.Views + 1

	_, err = r.service.Update(
		ctx,
		domain,
		link.ShortID,
		shortenedurls.UpdateInput{
			Views: &newViews,
		},
	)
	if err != nil {
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
		_ = r.service.PublishViewUpdate(domain, link.ShortID, newViews, lastView)
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

func isDefaultDomain(defaultHost, host string) bool {
	baseHost := resolveBaseHost(defaultHost)
	if baseHost == "" || host == "" {
		return false
	}

	host = strings.ToLower(host)

	if host == baseHost {
		return true
	}

	if strings.HasPrefix(baseHost, "cf.") {
		return host == strings.TrimPrefix(baseHost, "cf.")
	}

	return host == "cf."+baseHost
}

func resolveBaseHost(siteBaseURL string) string {
	parsed, err := url.Parse(siteBaseURL)
	if err == nil && parsed.Hostname() != "" {
		return strings.ToLower(parsed.Hostname())
	}

	trimmed := strings.TrimSpace(siteBaseURL)
	trimmed = strings.TrimPrefix(trimmed, "http://")
	trimmed = strings.TrimPrefix(trimmed, "https://")
	if trimmed == "" {
		return ""
	}

	trimmed = strings.Split(trimmed, "/")[0]
	trimmed = strings.Split(trimmed, ":")[0]

	return strings.ToLower(trimmed)
}
