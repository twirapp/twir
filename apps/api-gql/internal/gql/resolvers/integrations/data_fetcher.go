package integrations

import (
	config "github.com/satont/twir/libs/config"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type DataFetcherOpts struct {
	fx.In

	Gorm   *gorm.DB
	Config config.Config
}

func NewIntegrationsDataFetcher(opts DataFetcherOpts) *DataFetcher {
	return &DataFetcher{
		gorm:   opts.Gorm,
		config: opts.Config,
	}
}

type DataFetcher struct {
	gorm   *gorm.DB
	config config.Config
}
