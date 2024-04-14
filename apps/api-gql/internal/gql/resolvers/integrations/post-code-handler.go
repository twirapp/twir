package integrations

import (
	config "github.com/satont/twir/libs/config"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type PostCodeHandlerOpts struct {
	fx.In

	Gorm   *gorm.DB
	Config config.Config
}

func NewPostCodeHandler(opts PostCodeHandlerOpts) *PostCodeHandler {
	return &PostCodeHandler{
		gorm:   opts.Gorm,
		config: opts.Config,
	}
}

type PostCodeHandler struct {
	gorm   *gorm.DB
	config config.Config
}
