package shortlinks

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	shortlinksbannedusaragentsrepository "github.com/twirapp/twir/libs/repositories/short_links_banned_user_agents"
	"go.uber.org/fx"
)

type deleteBannedUserAgent struct {
	service  *shortenedurls.Service
	sessions *auth.Auth
}

type DeleteBannedUserAgentOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Sessions *auth.Auth
}

func newDeleteBannedUserAgent(opts DeleteBannedUserAgentOpts) *deleteBannedUserAgent {
	return &deleteBannedUserAgent{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type deleteBannedUserAgentInput struct {
	ID string `path:"id" minLength:"1" required:"true"`
}

func (c *deleteBannedUserAgent) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-delete-banned-user-agent",
		Method:      http.MethodDelete,
		Path:        "/v1/short-links/banned-user-agents/{id}",
		Tags:        []string{"Short links"},
		Summary:     "Delete banned user agent pattern",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *deleteBannedUserAgent) Handler(
	ctx context.Context,
	input *deleteBannedUserAgentInput,
) (*httpbase.BaseOutputJson[any], error) {
	user, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	if err := c.service.DeleteBannedUserAgent(ctx, input.ID, user.ID); err != nil {
		switch {
		case errors.Is(err, shortlinksbannedusaragentsrepository.ErrNotFound):
			return nil, huma.NewError(http.StatusNotFound, "Banned user agent not found")
		default:
			return nil, huma.NewError(http.StatusInternalServerError, "Cannot delete banned user agent", err)
		}
	}

	return httpbase.CreateBaseOutputJson[any](nil), nil
}

func (c *deleteBannedUserAgent) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
