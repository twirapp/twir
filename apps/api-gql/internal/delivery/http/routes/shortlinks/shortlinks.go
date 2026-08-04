package shortlinks

import (
	"log/slog"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/delivery/http/middlewares"
	"github.com/twirapp/twir/apps/api-gql/internal/services/clientinfo"
	"github.com/twirapp/twir/apps/api-gql/internal/services/shortenedurls"
	shortlinkscustomdomains "github.com/twirapp/twir/apps/api-gql/internal/services/shortlinkscustomdomains"
	config "github.com/twirapp/twir/libs/config"
)

type Registration struct{}

type registerRoute interface {
	Register(huma.API)
}

func RegisterRoutes(api huma.API, config config.Config, service *shortenedurls.Service, customDomainsService *shortlinkscustomdomains.Service, sessions *auth.Auth, logger *slog.Logger, middlewaresService *middlewares.Middlewares, clientInfoService *clientinfo.Service) Registration {
	routes := []registerRoute{
		newCreate(CreateOpts{
			Config: config, Service: service, CustomDomainsService: customDomainsService,
			Sessions: sessions, Logger: logger, Middlewares: middlewaresService,
			ClientInfoService: clientInfoService,
		}),
		newInfo(InfoOpts{Service: service, Config: config}),
		newRedirect(RedirectOpts{
			Service: service, Config: config, Sessions: sessions,
			Logger: logger, ClientInfoService: clientInfoService,
		}),
		newProfile(ProfileOpts{Service: service, Config: config, Sessions: sessions}),
		newStatistics(StatisticsOpts{
			Service: service, Sessions: sessions,
			CustomDomainsService: customDomainsService, Config: config,
		}),
		newTopCountries(TopCountriesOpts{
			Service: service, Sessions: sessions,
			CustomDomainsService: customDomainsService, Config: config,
		}),
		newUpdate(UpdateOpts{
			Service: service, CustomDomainsService: customDomainsService,
			Sessions: sessions, Config: config,
		}),
		newDelete(DeleteOpts{
			Service: service, CustomDomainsService: customDomainsService, Sessions: sessions,
		}),
		newGetCustomDomain(GetCustomDomainOpts{
			CustomDomainsService: customDomainsService, Sessions: sessions, Config: config,
		}),
		newCreateCustomDomain(CreateCustomDomainOpts{
			CustomDomainsService: customDomainsService, Sessions: sessions, Config: config,
		}),
		newVerifyCustomDomain(VerifyCustomDomainOpts{
			CustomDomainsService: customDomainsService, Sessions: sessions, Config: config,
		}),
		newDeleteCustomDomain(DeleteCustomDomainOpts{
			CustomDomainsService: customDomainsService,
			ShortenedUrlsService: service,
			Sessions:             sessions,
		}),
		newAllowCustomDomain(AllowCustomDomainOpts{CustomDomainsService: customDomainsService}),
		newListPresets(ListPresetsOpts{Service: service, Sessions: sessions}),
		newCreatePreset(CreatePresetOpts{Service: service, Sessions: sessions}),
		newUpdatePreset(UpdatePresetOpts{Service: service, Sessions: sessions}),
		newDeletePreset(DeletePresetOpts{Service: service, Sessions: sessions}),
		newListPresetPatterns(ListPresetPatternsOpts{Service: service, Sessions: sessions}),
		newCreatePresetPattern(CreatePresetPatternOpts{Service: service, Sessions: sessions}),
		newDeletePresetPattern(DeletePresetPatternOpts{Service: service, Sessions: sessions}),
		newListLinkPresets(ListLinkPresetsOpts{Service: service, Sessions: sessions}),
		newApplyPresetToLink(ApplyPresetToLinkOpts{Service: service, Sessions: sessions}),
		newRemovePresetFromLink(RemovePresetFromLinkOpts{Service: service, Sessions: sessions}),
		newListLinkBannedUserAgents(ListLinkBannedUserAgentsOpts{Service: service, Sessions: sessions}),
		newCreateLinkBannedUserAgent(CreateLinkBannedUserAgentOpts{Service: service, Sessions: sessions}),
		newDeleteLinkBannedUserAgent(DeleteLinkBannedUserAgentOpts{Service: service, Sessions: sessions}),
	}
	for _, route := range routes {
		route.Register(api)
	}

	return Registration{}
}

type linkOutputDto struct {
	Id        string    `json:"id" example:"KKMEa"`
	Url       string    `json:"url" example:"https://example.com"`
	ShortUrl  string    `json:"short_url" example:"https://twir.app/s/KKMEa"`
	Views     int       `json:"views" example:"1"`
	CreatedAt time.Time `json:"created_at" format:"date-time" example:"2023-01-01T00:00:00Z"`
}
