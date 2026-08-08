package scheduled_vips

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/services/scheduledvips"
)

type Registration struct{}

func RegisterRoutes(api huma.API, service *scheduledvips.Service, sessions *auth.Auth) Registration {
	newCreate(CreateOpts{Service: service, Sessions: sessions}).Register(api)
	newList(ListOpts{Service: service, Sessions: sessions}).Register(api)
	newDelete(DeleteOpts{Service: service, Sessions: sessions}).Register(api)
	return Registration{}
}

type scheduledVipOutputDto struct {
	ID         string     `json:"id" example:"550e8400-e29b-41d4-a716-446655440000"`
	UserID     string     `json:"user_id" example:"123456789"`
	ChannelID  string     `json:"channel_id" example:"987654321"`
	CreatedAt  time.Time  `json:"created_at" format:"date-time"`
	RemoveAt   *time.Time `json:"remove_at,omitempty" format:"date-time" nullable:"true"`
	RemoveType *string    `json:"remove_type,omitempty" nullable:"true" example:"time" enum:"time,stream_end"`
}
