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

type createBannedUserAgent struct {
	service  *shortenedurls.Service
	sessions *auth.Auth
}

type CreateBannedUserAgentOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Sessions *auth.Auth
}

func newCreateBannedUserAgent(opts CreateBannedUserAgentOpts) *createBannedUserAgent {
	return &createBannedUserAgent{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type createBannedUserAgentInput struct {
	Body struct {
		Pattern     string  `json:"pattern" required:"true" minLength:"1" maxLength:"512" example:"Chatterino.*"`
		Description *string `json:"description,omitempty" maxLength:"256"`
	}
}

func (c *createBannedUserAgent) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-create-banned-user-agent",
		Method:      http.MethodPost,
		Path:        "/v1/short-links/banned-user-agents",
		Tags:        []string{"Short links"},
		Summary:     "Create banned user agent pattern",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *createBannedUserAgent) Handler(
	ctx context.Context,
	input *createBannedUserAgentInput,
) (*httpbase.BaseOutputJson[bannedUserAgentDto], error) {
	user, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	item, err := c.service.CreateBannedUserAgent(
		ctx,
		shortlinksbannedusaragentsrepository.CreateInput{
			UserID:      user.ID,
			Pattern:     input.Body.Pattern,
			Description: input.Body.Description,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, shortlinksbannedusaragentsrepository.ErrAlreadyExists):
			return nil, huma.NewError(http.StatusConflict, "Pattern already exists", err)
		default:
			return nil, huma.NewError(http.StatusBadRequest, "Cannot create banned user agent", err)
		}
	}

	return httpbase.CreateBaseOutputJson(bannedUserAgentDto{
		ID:          item.ID,
		Pattern:     item.Pattern,
		Description: item.Description,
		CreatedAt:   item.CreatedAt,
	}), nil
}

func (c *createBannedUserAgent) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
