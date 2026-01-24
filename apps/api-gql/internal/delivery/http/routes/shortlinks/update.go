package shortlinks

import (
	"context"
	"net/http"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	shortlinkscustomdomains "github.com/twirapp/twir/apps/api-gql/internal/services/shortlinkscustomdomains"
	config "github.com/twirapp/twir/libs/config"
	"go.uber.org/fx"
)

type updateRequestDto struct {
	ShortId string `path:"shortId" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
	Body    struct {
		NewShortId *string `json:"new_short_id,omitempty" minLength:"3" maxLength:"50" pattern:"^[a-zA-Z0-9]+$"`
		Url        *string `json:"url,omitempty" minLength:"1" maxLength:"2048" format:"uri"`
	}
}

var _ httpbase.Route[*updateRequestDto, *httpbase.BaseOutputJson[linkOutputDto]] = (*updateRoute)(nil)

type UpdateOpts struct {
	fx.In

	Service              *shortenedurls.Service
	CustomDomainsService *shortlinkscustomdomains.Service
	Sessions             *auth.Auth
	Config               config.Config
}

type updateRoute struct {
	service              *shortenedurls.Service
	customDomainsService *shortlinkscustomdomains.Service
	sessions             *auth.Auth
	config               config.Config
}

func newUpdate(opts UpdateOpts) *updateRoute {
	return &updateRoute{
		service:              opts.Service,
		customDomainsService: opts.CustomDomainsService,
		sessions:             opts.Sessions,
		config:               opts.Config,
	}
}

func (u *updateRoute) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-url-update",
		Method:      http.MethodPatch,
		Path:        "/v1/short-links/{shortId}",
		Tags:        []string{"Short links"},
		Summary:     "Update short url",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (u *updateRoute) Register(api huma.API) {
	huma.Register(
		api,
		u.GetMeta(),
		u.Handler,
	)
}

func (u *updateRoute) Handler(
	ctx context.Context,
	input *updateRequestDto,
) (*httpbase.BaseOutputJson[linkOutputDto], error) {
	user, err := u.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized", err)
	}

	var domain *string
	if userDomain, err := u.customDomainsService.GetByUserID(ctx, user.ID); err == nil && !userDomain.IsNil() && userDomain.Verified {
		domain = &userDomain.Domain
	}

	// Get the link to verify ownership
	link, err := u.service.GetByShortID(ctx, domain, input.ShortId)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get link", err)
	}
	if link.IsNil() && domain != nil {
		domain = nil
		link, err = u.service.GetByShortID(ctx, domain, input.ShortId)
		if err != nil {
			return nil, huma.NewError(http.StatusInternalServerError, "Cannot get link", err)
		}
	}
	if link.IsNil() {
		return nil, huma.NewError(http.StatusNotFound, "Link not found")
	}

	// Check if user owns this link
	if link.CreatedByUserId == nil || *link.CreatedByUserId != user.ID {
		return nil, huma.NewError(http.StatusForbidden, "You don't have permission to update this link")
	}

	// Validate at least one field is provided
	if input.Body.NewShortId == nil && input.Body.Url == nil {
		return nil, huma.NewError(http.StatusBadRequest, "At least one field must be provided")
	}

	// Check if new short ID already exists
	if input.Body.NewShortId != nil && *input.Body.NewShortId != input.ShortId {
		existingLink, err := u.service.GetByShortID(ctx, domain, *input.Body.NewShortId)
		if err == nil && !existingLink.IsNil() {
			return nil, huma.NewError(http.StatusConflict, "Short ID already exists")
		}
	}

	// Update the link
	updateInput := shortenedurls.UpdateInput{
		ShortID: input.Body.NewShortId,
		URL:     input.Body.Url,
	}

	updatedLink, err := u.service.Update(ctx, domain, input.ShortId, updateInput)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot update link", err)
	}

	// Build short URL
	baseUrl, err := gincontext.GetBaseUrlFromContext(ctx, u.config.SiteBaseUrl)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get base URL", err)
	}

	var shortURL string
	if updatedLink.Domain != nil {
		shortURL = "https://" + *updatedLink.Domain + "/" + updatedLink.ID
	} else {
		parsedBaseUrl, _ := url.Parse(baseUrl)
		parsedBaseUrl.Path = "/s/" + updatedLink.ID
		shortURL = parsedBaseUrl.String()
	}

	return httpbase.CreateBaseOutputJson(
		linkOutputDto{
			Id:        updatedLink.ID,
			Url:       updatedLink.Link,
			ShortUrl:  shortURL,
			Views:     updatedLink.Views,
			CreatedAt: updatedLink.CreatedAt,
		},
	), nil
}
