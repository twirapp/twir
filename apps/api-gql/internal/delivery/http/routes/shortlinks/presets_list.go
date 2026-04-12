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

type presetDto struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at" format:"date-time"`
	UpdatedAt   time.Time `json:"updated_at" format:"date-time"`
}

type listPresets struct {
	service  *shortenedurls.Service
	sessions *auth.Auth
}

type ListPresetsOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Sessions *auth.Auth
}

func newListPresets(opts ListPresetsOpts) *listPresets {
	return &listPresets{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

func (c *listPresets) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-list-presets",
		Method:      http.MethodGet,
		Path:        "/v1/short-links/presets",
		Tags:        []string{"Short links"},
		Summary:     "List banned UA presets",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *listPresets) Handler(
	ctx context.Context,
	input *struct{},
) (*httpbase.BaseOutputJson[[]presetDto], error) {
	user, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	items, err := c.service.GetPresets(ctx, user.ID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get presets", err)
	}

	result := make([]presetDto, 0, len(items))
	for _, item := range items {
		result = append(result, presetDto{
			ID:          item.ID,
			Name:        item.Name,
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
			UpdatedAt:   item.UpdatedAt,
		})
	}

	return httpbase.CreateBaseOutputJson(result), nil
}

func (c *listPresets) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
