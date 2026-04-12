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

type updatePreset struct {
	service  *shortenedurls.Service
	sessions *auth.Auth
}

type UpdatePresetOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Sessions *auth.Auth
}

func newUpdatePreset(opts UpdatePresetOpts) *updatePreset {
	return &updatePreset{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type updatePresetInput struct {
	PresetID string `path:"presetId" minLength:"1" required:"true"`
	Body     struct {
		Name        *string `json:"name,omitempty" minLength:"1" maxLength:"100"`
		Description *string `json:"description,omitempty" maxLength:"256"`
	}
}

func (c *updatePreset) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-update-preset",
		Method:      http.MethodPatch,
		Path:        "/v1/short-links/presets/{presetId}",
		Tags:        []string{"Short links"},
		Summary:     "Update banned UA preset",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *updatePreset) Handler(
	ctx context.Context,
	input *updatePresetInput,
) (*httpbase.BaseOutputJson[presetDto], error) {
	_, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	item, err := c.service.UpdatePreset(
		ctx,
		input.PresetID,
		shortlinksbanneduapresetsrepository.UpdateInput{
			Name:        input.Body.Name,
			Description: input.Body.Description,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, shortlinksbanneduapresetsrepository.ErrNotFound):
			return nil, huma.NewError(http.StatusNotFound, "Preset not found")
		case errors.Is(err, shortlinksbanneduapresetsrepository.ErrAlreadyExists):
			return nil, huma.NewError(http.StatusConflict, "Preset with this name already exists", err)
		default:
			return nil, huma.NewError(http.StatusBadRequest, "Cannot update preset", err)
		}
	}

	return httpbase.CreateBaseOutputJson(presetDto{
		ID:          item.ID,
		Name:        item.Name,
		Description: item.Description,
		CreatedAt:   item.CreatedAt,
		UpdatedAt:   item.UpdatedAt,
	}), nil
}

func (c *updatePreset) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
