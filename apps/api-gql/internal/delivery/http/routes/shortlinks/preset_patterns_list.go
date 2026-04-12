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

type presetPatternDto struct {
	ID          string    `json:"id"`
	PresetID    string    `json:"preset_id"`
	Pattern     string    `json:"pattern"`
	Description *string   `json:"description"`
	CreatedAt   time.Time `json:"created_at" format:"date-time"`
}

type listPresetPatterns struct {
	service  *shortenedurls.Service
	sessions *auth.Auth
}

type ListPresetPatternsOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Sessions *auth.Auth
}

func newListPresetPatterns(opts ListPresetPatternsOpts) *listPresetPatterns {
	return &listPresetPatterns{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type listPresetPatternsInput struct {
	PresetID string `path:"presetId" minLength:"1" required:"true"`
}

func (c *listPresetPatterns) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-list-preset-patterns",
		Method:      http.MethodGet,
		Path:        "/v1/short-links/presets/{presetId}/patterns",
		Tags:        []string{"Short links"},
		Summary:     "List patterns in a preset",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *listPresetPatterns) Handler(
	ctx context.Context,
	input *listPresetPatternsInput,
) (*httpbase.BaseOutputJson[[]presetPatternDto], error) {
	_, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	items, err := c.service.GetPresetPatterns(ctx, input.PresetID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get preset patterns", err)
	}

	result := make([]presetPatternDto, 0, len(items))
	for _, item := range items {
		result = append(result, presetPatternDto{
			ID:          item.ID,
			PresetID:    item.PresetID,
			Pattern:     item.Pattern,
			Description: item.Description,
			CreatedAt:   item.CreatedAt,
		})
	}

	return httpbase.CreateBaseOutputJson(result), nil
}

func (c *listPresetPatterns) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
