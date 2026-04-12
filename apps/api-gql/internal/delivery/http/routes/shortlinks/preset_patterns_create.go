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

type createPresetPattern struct {
	service  *shortenedurls.Service
	sessions *auth.Auth
}

type CreatePresetPatternOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Sessions *auth.Auth
}

func newCreatePresetPattern(opts CreatePresetPatternOpts) *createPresetPattern {
	return &createPresetPattern{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type createPresetPatternInput struct {
	PresetID string `path:"presetId" minLength:"1" required:"true"`
	Body     struct {
		Pattern     string  `json:"pattern" required:"true" minLength:"1" maxLength:"512"`
		Description *string `json:"description,omitempty" maxLength:"256"`
	}
}

func (c *createPresetPattern) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-create-preset-pattern",
		Method:      http.MethodPost,
		Path:        "/v1/short-links/presets/{presetId}/patterns",
		Tags:        []string{"Short links"},
		Summary:     "Add pattern to preset",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *createPresetPattern) Handler(
	ctx context.Context,
	input *createPresetPatternInput,
) (*httpbase.BaseOutputJson[presetPatternDto], error) {
	_, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	item, err := c.service.CreatePresetPattern(
		ctx,
		shortlinksbanneduapresetpatternsrepository.CreateInput{
			PresetID:    input.PresetID,
			Pattern:     input.Body.Pattern,
			Description: input.Body.Description,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, shortlinksbanneduapresetpatternsrepository.ErrAlreadyExists):
			return nil, huma.NewError(http.StatusConflict, "Pattern already exists in this preset", err)
		default:
			return nil, huma.NewError(http.StatusBadRequest, "Cannot create preset pattern", err)
		}
	}

	return httpbase.CreateBaseOutputJson(presetPatternDto{
		ID:          item.ID,
		PresetID:    item.PresetID,
		Pattern:     item.Pattern,
		Description: item.Description,
		CreatedAt:   item.CreatedAt,
	}), nil
}

func (c *createPresetPattern) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
