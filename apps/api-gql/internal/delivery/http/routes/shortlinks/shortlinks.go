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

type Dependencies struct {
	API                  huma.API
	Config               config.Config
	Service              *shortenedurls.Service
	CustomDomainsService *shortlinkscustomdomains.Service
	Sessions             *auth.Auth
	Logger               *slog.Logger
	Middlewares          *middlewares.Middlewares
	ClientInfoService    *clientinfo.Service
}

type Registration struct{}

type registerRoute interface {
	Register(huma.API)
}

func RegisterRoutes(deps Dependencies) Registration {
	routes := []registerRoute{
		newCreate(CreateOpts{
			Config: deps.Config, Service: deps.Service, CustomDomainsService: deps.CustomDomainsService,
			Sessions: deps.Sessions, Logger: deps.Logger, Middlewares: deps.Middlewares,
			ClientInfoService: deps.ClientInfoService,
		}),
		newInfo(InfoOpts{Service: deps.Service, Config: deps.Config}),
		newRedirect(RedirectOpts{
			Service: deps.Service, Config: deps.Config, Sessions: deps.Sessions,
			Logger: deps.Logger, ClientInfoService: deps.ClientInfoService,
		}),
		newProfile(ProfileOpts{Service: deps.Service, Config: deps.Config, Sessions: deps.Sessions}),
		newStatistics(StatisticsOpts{
			Service: deps.Service, Sessions: deps.Sessions,
			CustomDomainsService: deps.CustomDomainsService, Config: deps.Config,
		}),
		newTopCountries(TopCountriesOpts{
			Service: deps.Service, Sessions: deps.Sessions,
			CustomDomainsService: deps.CustomDomainsService, Config: deps.Config,
		}),
		newUpdate(UpdateOpts{
			Service: deps.Service, CustomDomainsService: deps.CustomDomainsService,
			Sessions: deps.Sessions, Config: deps.Config,
		}),
		newDelete(DeleteOpts{
			Service: deps.Service, CustomDomainsService: deps.CustomDomainsService, Sessions: deps.Sessions,
		}),
		newGetCustomDomain(GetCustomDomainOpts{
			CustomDomainsService: deps.CustomDomainsService, Sessions: deps.Sessions, Config: deps.Config,
		}),
		newCreateCustomDomain(CreateCustomDomainOpts{
			CustomDomainsService: deps.CustomDomainsService, Sessions: deps.Sessions, Config: deps.Config,
		}),
		newVerifyCustomDomain(VerifyCustomDomainOpts{
			CustomDomainsService: deps.CustomDomainsService, Sessions: deps.Sessions, Config: deps.Config,
		}),
		newDeleteCustomDomain(DeleteCustomDomainOpts{
			CustomDomainsService: deps.CustomDomainsService,
			ShortenedUrlsService: deps.Service,
			Sessions:             deps.Sessions,
		}),
		newAllowCustomDomain(AllowCustomDomainOpts{CustomDomainsService: deps.CustomDomainsService}),
		newListPresets(ListPresetsOpts{Service: deps.Service, Sessions: deps.Sessions}),
		newCreatePreset(CreatePresetOpts{Service: deps.Service, Sessions: deps.Sessions}),
		newUpdatePreset(UpdatePresetOpts{Service: deps.Service, Sessions: deps.Sessions}),
		newDeletePreset(DeletePresetOpts{Service: deps.Service, Sessions: deps.Sessions}),
		newListPresetPatterns(ListPresetPatternsOpts{Service: deps.Service, Sessions: deps.Sessions}),
		newCreatePresetPattern(CreatePresetPatternOpts{Service: deps.Service, Sessions: deps.Sessions}),
		newDeletePresetPattern(DeletePresetPatternOpts{Service: deps.Service, Sessions: deps.Sessions}),
		newListLinkPresets(ListLinkPresetsOpts{Service: deps.Service, Sessions: deps.Sessions}),
		newApplyPresetToLink(ApplyPresetToLinkOpts{Service: deps.Service, Sessions: deps.Sessions}),
		newRemovePresetFromLink(RemovePresetFromLinkOpts{Service: deps.Service, Sessions: deps.Sessions}),
		newListLinkBannedUserAgents(ListLinkBannedUserAgentsOpts{Service: deps.Service, Sessions: deps.Sessions}),
		newCreateLinkBannedUserAgent(CreateLinkBannedUserAgentOpts{Service: deps.Service, Sessions: deps.Sessions}),
		newDeleteLinkBannedUserAgent(DeleteLinkBannedUserAgentOpts{Service: deps.Service, Sessions: deps.Sessions}),
	}
	for _, route := range routes {
		route.Register(deps.API)
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
