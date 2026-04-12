package shortlinks

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	shortlinkslinkbannedusaragentsrepository "github.com/twirapp/twir/libs/repositories/short_links_link_banned_user_agents"
	"go.uber.org/fx"
)

type deleteLinkBannedUserAgent struct {
	service  *shortenedurls.Service
	sessions *auth.Auth
}

type DeleteLinkBannedUserAgentOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Sessions *auth.Auth
}

func newDeleteLinkBannedUserAgent(opts DeleteLinkBannedUserAgentOpts) *deleteLinkBannedUserAgent {
	return &deleteLinkBannedUserAgent{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type deleteLinkBannedUserAgentInput struct {
	LinkID string `path:"linkId" minLength:"1" required:"true"`
	ID     string `path:"id" minLength:"1" required:"true"`
}

func (c *deleteLinkBannedUserAgent) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-delete-link-banned-user-agent",
		Method:      http.MethodDelete,
		Path:        "/v1/short-links/{linkId}/banned-user-agents/{id}",
		Tags:        []string{"Short links"},
		Summary:     "Delete banned user agent pattern for a specific link",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *deleteLinkBannedUserAgent) Handler(
	ctx context.Context,
	input *deleteLinkBannedUserAgentInput,
) (*httpbase.BaseOutputJson[any], error) {
	_, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	if err := c.service.DeleteLinkBannedUserAgent(ctx, input.ID, input.LinkID); err != nil {
		switch {
		case errors.Is(err, shortlinkslinkbannedusaragentsrepository.ErrNotFound):
			return nil, huma.NewError(http.StatusNotFound, "Banned user agent not found")
		default:
			return nil, huma.NewError(http.StatusInternalServerError, "Cannot delete banned user agent", err)
		}
	}

	return httpbase.CreateBaseOutputJson[any](nil), nil
}

func (c *deleteLinkBannedUserAgent) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
