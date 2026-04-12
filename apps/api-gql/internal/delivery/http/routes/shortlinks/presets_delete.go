package shortlinks

import (
	"context"
	"errors"
	"net/http"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	shortlinksbanneduapresetsrepository "github.com/twirapp/twir/libs/repositories/short_links_banned_ua_presets"
	"go.uber.org/fx"
)

type deletePreset struct {
	service  *shortenedurls.Service
	sessions *auth.Auth
}

type DeletePresetOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Sessions *auth.Auth
}

func newDeletePreset(opts DeletePresetOpts) *deletePreset {
	return &deletePreset{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type deletePresetInput struct {
	PresetID string `path:"presetId" minLength:"1" required:"true"`
}

func (c *deletePreset) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-delete-preset",
		Method:      http.MethodDelete,
		Path:        "/v1/short-links/presets/{presetId}",
		Tags:        []string{"Short links"},
		Summary:     "Delete banned UA preset",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *deletePreset) Handler(
	ctx context.Context,
	input *deletePresetInput,
) (*httpbase.BaseOutputJson[any], error) {
	user, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	if err := c.service.DeletePreset(ctx, input.PresetID, user.ID); err != nil {
		switch {
		case errors.Is(err, shortlinksbanneduapresetsrepository.ErrNotFound):
			return nil, huma.NewError(http.StatusNotFound, "Preset not found")
		default:
			return nil, huma.NewError(http.StatusInternalServerError, "Cannot delete preset", err)
		}
	}

	return httpbase.CreateBaseOutputJson[any](nil), nil
}

func (c *deletePreset) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
