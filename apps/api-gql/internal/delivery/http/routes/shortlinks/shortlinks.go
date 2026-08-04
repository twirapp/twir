package shortlinks

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortlinkscustomdomains"
	config "github.com/twirapp/twir/libs/config"
	shortlinksbanneduapresetsrepository "github.com/twirapp/twir/libs/repositories/short_links_banned_ua_presets"
	"github.com/twirapp/twir/libs/repositories/shortened_urls/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Api      huma.API
	Config   config.Config
	Service  *shortenedurls.Service
	Sessions *auth.Auth
	Logger   *slog.Logger
}

var FxModule = fx.Provide(
	httpbase.AsFxRoute(newCreate),
	httpbase.AsFxRoute(newInfo),
	httpbase.AsFxRoute(newRedirect),
	httpbase.AsFxRoute(newProfile),
	httpbase.AsFxRoute(newStatistics),
	httpbase.AsFxRoute(newTopCountries),
	httpbase.AsFxRoute(newUpdate),
	httpbase.AsFxRoute(newDelete),
	httpbase.AsFxRoute(newGetCustomDomain),
	httpbase.AsFxRoute(newCreateCustomDomain),
	httpbase.AsFxRoute(newVerifyCustomDomain),
	httpbase.AsFxRoute(newDeleteCustomDomain),
	httpbase.AsFxRoute(newAllowCustomDomain),
	httpbase.AsFxRoute(newListPresets),
	httpbase.AsFxRoute(newCreatePreset),
	httpbase.AsFxRoute(newUpdatePreset),
	httpbase.AsFxRoute(newDeletePreset),
	httpbase.AsFxRoute(newListPresetPatterns),
	httpbase.AsFxRoute(newCreatePresetPattern),
	httpbase.AsFxRoute(newDeletePresetPattern),
	httpbase.AsFxRoute(newListLinkPresets),
	httpbase.AsFxRoute(newApplyPresetToLink),
	httpbase.AsFxRoute(newRemovePresetFromLink),
	httpbase.AsFxRoute(newListLinkBannedUserAgents),
	httpbase.AsFxRoute(newCreateLinkBannedUserAgent),
	httpbase.AsFxRoute(newDeleteLinkBannedUserAgent),
)

type linkOutputDto struct {
	Id        string    `json:"id" example:"KKMEa"`
	Url       string    `json:"url" example:"https://example.com"`
	ShortUrl  string    `json:"short_url" example:"https://twir.app/s/KKMEa"`
	Views     int       `json:"views" example:"1"`
	CreatedAt time.Time `json:"created_at" format:"date-time" example:"2023-01-01T00:00:00Z"`
}

var (
	errShortLinkNotFound  = errors.New("short link not found")
	errShortLinkForbidden = errors.New("you don't have permission to manage this link")
	errPresetNotFound     = errors.New("preset not found")
)

func resolveOwnedShortLink(
	ctx context.Context,
	service *shortenedurls.Service,
	customDomainsService *shortlinkscustomdomains.Service,
	userID string,
	linkID string,
) (model.ShortenedUrl, error) {
	var domain *string
	if userDomain, err := customDomainsService.GetByUserID(ctx, userID); err == nil &&
		!userDomain.IsNil() && userDomain.Verified {
		domain = &userDomain.Domain
	}

	link, err := service.GetByShortID(ctx, domain, linkID)
	if err != nil {
		return model.Nil, err
	}
	if link.IsNil() && domain != nil {
		link, err = service.GetByShortID(ctx, nil, linkID)
		if err != nil {
			return model.Nil, err
		}
	}
	if link.IsNil() {
		return model.Nil, errShortLinkNotFound
	}
	if link.CreatedByUserId == nil || *link.CreatedByUserId != userID {
		return model.Nil, errShortLinkForbidden
	}

	return link, nil
}

func resolveOwnedPreset(
	ctx context.Context,
	service *shortenedurls.Service,
	userID string,
	presetID string,
) error {
	preset, err := service.GetPresetByID(ctx, presetID)
	if err != nil {
		if errors.Is(err, shortlinksbanneduapresetsrepository.ErrNotFound) {
			return errPresetNotFound
		}
		return err
	}
	if preset.IsNil() || preset.UserID != userID {
		return errPresetNotFound
	}

	return nil
}
