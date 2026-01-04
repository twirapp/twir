package channels_overlays

import (
	"fmt"

	"github.com/google/uuid"
)

const customOverlayWsRouterKey = "api.customOverlaySettings"

func CreateCustomOverlayWsRouterKey(channelID string, id uuid.UUID) string {
	return fmt.Sprintf("%s.%s.%s", customOverlayWsRouterKey, channelID, id)
}
