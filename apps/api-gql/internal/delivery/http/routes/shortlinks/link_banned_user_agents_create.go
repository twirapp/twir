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

type createLinkBannedUserAgent struct {
	service  *shortenedurls.Service
	sessions *auth.Auth
}

type CreateLinkBannedUserAgentOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Sessions *auth.Auth
}

func newCreateLinkBannedUserAgent(opts CreateLinkBannedUserAgentOpts) *createLinkBannedUserAgent {
	return &createLinkBannedUserAgent{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type createLinkBannedUserAgentInput struct {
	LinkID string `path:"linkId" minLength:"1" required:"true"`
	Body   struct {
		Pattern     string  `json:"pattern" required:"true" minLength:"1" maxLength:"512" example:"Chatterino.*"`
		Description *string `json:"description,omitempty" maxLength:"256"`
	}
}

func (c *createLinkBannedUserAgent) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-create-link-banned-user-agent",
		Method:      http.MethodPost,
		Path:        "/v1/short-links/by-id/{linkId}/banned-user-agents",
		Tags:        []string{"Short links"},
		Summary:     "Create banned user agent pattern for a specific link",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *createLinkBannedUserAgent) Handler(
	ctx context.Context,
	input *createLinkBannedUserAgentInput,
) (*httpbase.BaseOutputJson[linkBannedUserAgentDto], error) {
	_, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	item, err := c.service.CreateLinkBannedUserAgent(
		ctx,
		shortlinkslinkbannedusaragentsrepository.CreateInput{
			LinkID:      input.LinkID,
			Pattern:     input.Body.Pattern,
			Description: input.Body.Description,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, shortlinkslinkbannedusaragentsrepository.ErrAlreadyExists):
			return nil, huma.NewError(http.StatusConflict, "Pattern already exists", err)
		default:
			return nil, huma.NewError(http.StatusBadRequest, "Cannot create banned user agent", err)
		}
	}

	return httpbase.CreateBaseOutputJson(linkBannedUserAgentDto{
		ID:          item.ID,
		LinkID:      item.LinkID,
		Pattern:     item.Pattern,
		Description: item.Description,
		CreatedAt:   item.CreatedAt,
	}), nil
}

func (c *createLinkBannedUserAgent) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
