package overlays

import (
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/websockets"
)

type DudesChangeColorRequest struct {
	websockets.DudesChangeColorRequest
}

type DudesGrowRequest struct {
	websockets.DudesGrowRequest
}

type DudesUserSettings struct {
	model.ChannelsOverlaysDudesUserSettings
}
