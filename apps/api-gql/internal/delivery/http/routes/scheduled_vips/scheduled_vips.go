package scheduled_vips

import (
	"time"

	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"go.uber.org/fx"
)

var FxModule = fx.Provide(
	httpbase.AsFxRoute(newCreate),
	httpbase.AsFxRoute(newList),
	httpbase.AsFxRoute(newDelete),
)

type scheduledVipOutputDto struct {
	ID         string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID     string     `json:"user_id" example:"123456789"`
	ChannelID  string     `json:"channel_id" example:"987654321"`
	CreatedAt  time.Time  `json:"created_at" format:"date-time"`
	RemoveAt   *time.Time `json:"remove_at,omitempty" format:"date-time" nullable:"true"`
	RemoveType *string    `json:"remove_type,omitempty" nullable:"true" example:"time" enum:"time,stream_end"`
}
