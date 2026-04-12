package shortlinks

import (
	"context"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	"go.uber.org/fx"
)

type bannedUserAgentDto struct {
	ID          string    `json:"id"`
	Pattern     string    `json:"pattern"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at" format:"date-time"`
}

type listBannedUserAgents struct {
	service  *shortenedurls.Service
	sessions *auth.Auth
}

type ListBannedUserAgentsOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Sessions *auth.Auth
}

func newListBannedUserAgents(opts ListBannedUserAgentsOpts) *listBannedUserAgents {
	return &listBannedUserAgents{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

func (c *listBannedUserAgents) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-list-banned-user-agents",
		Method:      http.MethodGet,
		Path:        "/v1/short-links/banned-user-agents",
		Tags:        []string{"Short links"},
		Summary:     "List banned user agent patterns",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *listBannedUserAgents) Handler(
	ctx context.Context,
	input *struct{},
) (*httpbase.BaseOutputJson[[]bannedUserAgentDto], error) {
	user, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	items, err := c.service.GetBannedUserAgents(ctx, user.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get banned user agents", err)
	}

	result := make([]bannedUserAgentDto, 0, len(items))
	for _, item := range items {
		result = append(result, bannedUserAgentDto{
			ID:          item.ID,
			Pattern:     item.Pattern,
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
		})
	}

	return httpbase.CreateBaseOutputJson(result), nil
}

func (c *listBannedUserAgents) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
