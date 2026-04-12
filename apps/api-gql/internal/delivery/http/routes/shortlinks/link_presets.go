package shortlinks

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	shortlinkslinkpresetsrepository "github.com/twirapp/twir/libs/repositories/short_links_link_presets"
	"go.uber.org/fx"
)

type linkPresetDto struct {
	ID        string    `json:"id"`
	LinkID    string    `json:"link_id"`
	PresetID  string    `json:"preset_id"`
	CreatedAt time.Time `json:"created_at" format:"date-time"`
}

type listLinkPresets struct {
	service  *shortenedurls.Service
	sessions *auth.Auth
}

type ListLinkPresetsOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Sessions *auth.Auth
}

func newListLinkPresets(opts ListLinkPresetsOpts) *listLinkPresets {
	return &listLinkPresets{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type listLinkPresetsInput struct {
	LinkID string `path:"linkId" minLength:"1" required:"true"`
}

func (c *listLinkPresets) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-list-link-presets",
		Method:      http.MethodGet,
		Path:        "/v1/short-links/by-id/{linkId}/presets",
		Tags:        []string{"Short links"},
		Summary:     "List presets applied to a link",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *listLinkPresets) Handler(
	ctx context.Context,
	input *listLinkPresetsInput,
) (*httpbase.BaseOutputJson[[]linkPresetDto], error) {
	_, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	items, err := c.service.GetLinkPresets(ctx, input.LinkID)
	if err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot get link presets", err)
	}

	result := make([]linkPresetDto, 0, len(items))
	for _, item := range items {
		result = append(result, linkPresetDto{
			ID:        item.ID,
			LinkID:    item.LinkID,
			PresetID:  item.PresetID,
			CreatedAt: item.CreatedAt,
		})
	}

	return httpbase.CreateBaseOutputJson(result), nil
}

func (c *listLinkPresets) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}

type applyPresetToLink struct {
	service  *shortenedurls.Service
	sessions *auth.Auth
}

type ApplyPresetToLinkOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Sessions *auth.Auth
}

func newApplyPresetToLink(opts ApplyPresetToLinkOpts) *applyPresetToLink {
	return &applyPresetToLink{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type applyPresetToLinkInput struct {
	LinkID string `path:"linkId" minLength:"1" required:"true"`
	Body   struct {
		PresetID string `json:"preset_id" required:"true"`
	}
}

func (c *applyPresetToLink) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-apply-preset-to-link",
		Method:      http.MethodPost,
		Path:        "/v1/short-links/by-id/{linkId}/presets",
		Tags:        []string{"Short links"},
		Summary:     "Apply preset to link",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *applyPresetToLink) Handler(
	ctx context.Context,
	input *applyPresetToLinkInput,
) (*httpbase.BaseOutputJson[linkPresetDto], error) {
	_, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	item, err := c.service.ApplyPresetToLink(
		ctx,
		shortlinkslinkpresetsrepository.CreateInput{
			LinkID:   input.LinkID,
			PresetID: input.Body.PresetID,
		},
	)
	if err != nil {
		switch {
		case errors.Is(err, shortlinkslinkpresetsrepository.ErrAlreadyExists):
			return nil, huma.NewError(http.StatusConflict, "Preset already applied to this link", err)
		default:
			return nil, huma.NewError(http.StatusBadRequest, "Cannot apply preset to link", err)
		}
	}

	return httpbase.CreateBaseOutputJson(linkPresetDto{
		ID:        item.ID,
		LinkID:    item.LinkID,
		PresetID:  item.PresetID,
		CreatedAt: item.CreatedAt,
	}), nil
}

func (c *applyPresetToLink) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}

type removePresetFromLink struct {
	service  *shortenedurls.Service
	sessions *auth.Auth
}

type RemovePresetFromLinkOpts struct {
	fx.In

	Service  *shortenedurls.Service
	Sessions *auth.Auth
}

func newRemovePresetFromLink(opts RemovePresetFromLinkOpts) *removePresetFromLink {
	return &removePresetFromLink{
		service:  opts.Service,
		sessions: opts.Sessions,
	}
}

type removePresetFromLinkInput struct {
	LinkID   string `path:"linkId" minLength:"1" required:"true"`
	PresetID string `path:"presetId" minLength:"1" required:"true"`
}

func (c *removePresetFromLink) GetMeta() huma.Operation {
	return huma.Operation{
		OperationID: "short-links-remove-preset-from-link",
		Method:      http.MethodDelete,
		Path:        "/v1/short-links/by-id/{linkId}/presets/{presetId}",
		Tags:        []string{"Short links"},
		Summary:     "Remove preset from link",
		Security: []map[string][]string{
			{"bearer": {}},
		},
	}
}

func (c *removePresetFromLink) Handler(
	ctx context.Context,
	input *removePresetFromLinkInput,
) (*httpbase.BaseOutputJson[any], error) {
	_, err := c.sessions.GetAuthenticatedUserModel(ctx)
	if err != nil {
		return nil, huma.NewError(http.StatusUnauthorized, "Unauthorized")
	}

	if err := c.service.RemovePresetFromLink(ctx, input.LinkID, input.PresetID); err != nil {
		return nil, huma.NewError(http.StatusInternalServerError, "Cannot remove preset from link", err)
	}

	return httpbase.CreateBaseOutputJson[any](nil), nil
}

func (c *removePresetFromLink) Register(api huma.API) {
	huma.Register(api, c.GetMeta(), c.Handler)
}
