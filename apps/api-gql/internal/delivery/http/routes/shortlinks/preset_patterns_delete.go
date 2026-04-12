package shortlinks

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	shortlinksbanneduapresetpatternsrepository "github.com/twirapp/twir/libs/repositories/short_links_banned_ua_preset_patterns"
	"go.uber.org/fx"
)

type deletePresetPattern struct {
	service  *shortenedurls.Service
	sessions *auth.Auth
}

type DeletePresetPatternOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Sessions *auth.Auth
}

func newDeletePresetPattern(opts DeletePresetPatternOpts) *deletePresetPattern {
	return &deletePresetPattern{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type deletePresetPatternInput struct {
	PresetID string `path:"presetId" minLength:"1" required:"true"`
	ID       string `path:"id" minLength:"1" required:"true"`
}

func (c *deletePresetPattern) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-delete-preset-pattern",
		Method:      http.MethodDelete,
		Path:        "/v1/short-links/presets/{presetId}/patterns/{id}",
		Tags:        []string{"Short links"},
		Summary:     "Remove pattern from preset",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *deletePresetPattern) Handler(
	ctx context.Context,
	input *deletePresetPatternInput,
) (*httpbase.BaseOutputJson[any], error) {
	_, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	if err := c.service.DeletePresetPattern(ctx, input.ID, input.PresetID); err != nil {
		switch {
		case errors.Is(err, shortlinksbanneduapresetpatternsrepository.ErrNotFound):
			return nil, huma.NewError(http.StatusNotFound, "Pattern not found")
		default:
			return nil, huma.NewError(http.StatusInternalServerError, "Cannot delete pattern", err)
		}
	}

	return httpbase.CreateBaseOutputJson[any](nil), nil
}

func (c *deletePresetPattern) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
