package shortlinks

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortlinkscustomdomains"
)

type linkBannedUserAgentDto struct {
	ID          string    `json:"id"`
	LinkID      string    `json:"link_id"`
	Pattern     string    `json:"pattern"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at" format:"date-time"`
}

type listLinkBannedUserAgents struct {
	service              *shortenedurls.Service
	customDomainsService *shortlinkscustomdomains.Service
	sessions             *auth.Auth
}

type ListLinkBannedUserAgentsOpts struct {
	Service              *shortenedurls.Service
	CustomDomainsService *shortlinkscustomdomains.Service
	Sessions             *auth.Auth
}

func newListLinkBannedUserAgents(opts ListLinkBannedUserAgentsOpts) *listLinkBannedUserAgents {
	return &listLinkBannedUserAgents{
		service:              opts.Service,
		customDomainsService: opts.CustomDomainsService,
		sessions:             opts.Sessions,
	}
}

type listLinkBannedUserAgentsInput struct {
	LinkID string `path:"linkId" minLength:"1" required:"true"`
}

func (c *listLinkBannedUserAgents) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-list-link-banned-user-agents",
		Method:      http.MethodGet,
		Path:        "/v1/short-links/by-id/{linkId}/banned-user-agents",
		Tags:        []string{"Short links"},
		Summary:     "List banned user agent patterns for a specific link",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *listLinkBannedUserAgents) Handler(
	ctx context.Context,
	input *listLinkBannedUserAgentsInput,
) (*httpbase.BaseOutputJson[[]linkBannedUserAgentDto], error) {
	user, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	link, err := resolveOwnedShortLink(ctx, c.service, c.customDomainsService, user.ID, input.LinkID)
	if err != nil {
		switch {
		case errors.Is(err, errShortLinkNotFound):
			return nil, huma.NewError(http.StatusNotFound, "Link not found")
		case errors.Is(err, errShortLinkForbidden):
			return nil, huma.NewError(http.StatusForbidden, "You don't have permission to manage this link")
		default:
			return nil, huma.NewError(http.StatusInternalServerError, "Cannot get link", err)
		}
	}

	items, err := c.service.GetLinkBannedUserAgents(ctx, link.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get banned user agents", err)
	}

	result := make([]linkBannedUserAgentDto, 0, len(items))
	for _, item := range items {
		result = append(result, linkBannedUserAgentDto{
			ID:          item.ID,
			LinkID:      item.LinkID,
			Pattern:     item.Pattern,
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
		})
	}

	return httpbase.CreateBaseOutputJson(result), nil
}

func (c *listLinkBannedUserAgents) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
