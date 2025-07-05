package overlays_dudes

import (
	"fmt"

	"github.com/google/uuid"
)

const dudesWsRouterKey = "api.newDudesOverlaySettings"

func CreateDudesWsRouterKey(channelID string, id uuid.UUID) string {
	return fmt.Sprintf("%s.%s.%s", dudesWsRouterKey, channelID, id)
}
