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

type createPreset struct {
	service  *shortenedurls.Service
	sessions *auth.Auth
}

type CreatePresetOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Sessions *auth.Auth
}

func newCreatePreset(opts CreatePresetOpts) *createPreset {
	return &createPreset{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type createPresetInput struct {
	Body struct {
		Name        string  `json:"name" required:"true" minLength:"1" maxLength:"100"`
		Description *string `json:"description,omitempty" maxLength:"256"`
	}
}

func (c *createPreset) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-create-preset",
		Method:      http.MethodPost,
		Path:        "/v1/short-links/presets",
		Tags:        []string{"Short links"},
		Summary:     "Create banned UA preset",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *createPreset) Handler(
	ctx context.Context,
	input *createPresetInput,
) (*httpbase.BaseOutputJson[presetDto], error) {
	user, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	item, err := c.service.CreatePreset(
		ctx,
		shortlinksbanneduapresetsrepository.CreateInput{
			UserID:      user.ID,
			Name:        input.Body.Name,
			Description: input.Body.Description,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, shortlinksbanneduapresetsrepository.ErrAlreadyExists):
			return nil, huma.NewError(http.StatusConflict, "Preset with this name already exists", err)
		default:
			return nil, huma.NewError(http.StatusBadRequest, "Cannot create preset", err)
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

func (c *createPreset) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
