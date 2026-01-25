package shortlinks

import (
	"context"
	"errors"
	"net/http"
	"net/url"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/server/gincontext"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	shortlinkscustomdomains "github.com/twirapp/twir/apps/api-gql/internal/services/shortlinkscustomdomains"
	config "github.com/twirapp/twir/libs/config"
	shortenedurlsrepository "github.com/twirapp/twir/libs/repositories/shortened_urls"
	"go.uber.org/fx"
)

type updateRequestDto struct {
	ShortId string `path:"shortId" minLength:"1" pattern:"^[a-zA-Z0-9]+$" required:"true"`
	Body    struct {
		NewShortId      *string `json:"new_short_id,omitempty" minLength:"3" maxLength:"50" pattern:"^[a-zA-Z0-9]+$"`
		Url             *string `json:"url,omitempty" minLength:"1" maxLength:"2048" format:"uri"`
		UseCustomDomain *bool   `json:"use_custom_domain,omitempty"`
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

	var (
		domain            *string
		hasCustomDomain   bool
		customDomainValue string
	)
	if userDomain, err := u.customDomainsService.GetByUserID(ctx, user.ID); err == nil && !userDomain.IsNil() && userDomain.Verified {
		hasCustomDomain = true
		customDomainValue = userDomain.Domain
		domain = &customDomainValue
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
	currentDomain := link.Domain

	// Check if user owns this link
	if link.CreatedByUserId == nil || *link.CreatedByUserId != user.ID {
		return nil, huma.NewError(http.StatusForbidden, "You don't have permission to update this link")
	}

	// Validate at least one field is provided
	if input.Body.NewShortId == nil && input.Body.Url == nil && input.Body.UseCustomDomain == nil {
		return nil, huma.NewError(http.StatusBadRequest, "At least one field must be provided")
	}

	targetDomain := currentDomain
	if input.Body.UseCustomDomain != nil {
		if *input.Body.UseCustomDomain {
			if !hasCustomDomain {
				return nil, huma.NewError(http.StatusBadRequest, "Custom domain is not available")
			}
			targetDomain = &customDomainValue
		} else {
			targetDomain = nil
		}
	}

	targetShortID := link.ShortID
	shortIDChanged := input.Body.NewShortId != nil && *input.Body.NewShortId != link.ShortID
	if shortIDChanged {
		targetShortID = *input.Body.NewShortId
	}

	domainChanged := (currentDomain == nil) != (targetDomain == nil)
	if !domainChanged && currentDomain != nil && targetDomain != nil {
		domainChanged = *currentDomain != *targetDomain
	}

	if shortIDChanged || domainChanged {
		existingLink, err := u.service.GetByShortID(ctx, targetDomain, targetShortID)
		if err != nil {
			return nil, huma.NewError(http.StatusInternalServerError, "Cannot get link", err)
		}
		if !existingLink.IsNil() {
			return nil, huma.NewError(http.StatusConflict, "Short ID already exists")
		}
	}

	// Update the link
	updateInput := shortenedurls.UpdateInput{
		URL: input.Body.Url,
	}
	if shortIDChanged {
		updateInput.ShortID = input.Body.NewShortId
	}
	if domainChanged {
		if targetDomain == nil {
			updateInput.ClearDomain = true
		} else {
			domainValue := *targetDomain
			updateInput.Domain = &domainValue
		}
	}

	updatedLink, err := u.service.Update(ctx, currentDomain, input.ShortId, updateInput)
	if err != nil {
		if errors.Is(err, shortenedurlsrepository.ErrShortIDAlreadyExists) {
			return nil, huma.NewError(http.StatusConflict, "Short ID already exists", err)
		}
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
