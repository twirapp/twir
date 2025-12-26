package brb

import (
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"go.uber.org/fx"
)

var FxModule = fx.Provide(
	httpbase.AsFxRoute(newStop),
	httpbase.AsFxRoute(newStart),
)
