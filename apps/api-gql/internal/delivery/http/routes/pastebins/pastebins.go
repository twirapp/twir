package pastebins

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/services/clientinfo"
	"github.com/twirapp/twir/apps/api-gql/internal/services/pastebins"
	config "github.com/twirapp/twir/libs/config"
)

type Registration struct{}

func RegisterRoutes(api huma.API, config config.Config, service *pastebins.Service, sessions *auth.Auth, clientInfoService *clientinfo.Service) Registration {
	newProfile(ProfileOpts{Service: service, Sessions: sessions}).Register(api)
	newGetById(GetByIdOpts{Service: service}).Register(api)
	newCreate(CreateOpts{
		Service: service, Sessions: sessions, ClientInfoService: clientInfoService,
	}).Register(api)
	newDelete(CreateOpts{Service: service, Sessions: sessions, ClientInfoService: clientInfoService}).Register(api)
	return Registration{}
}

type pasteBinOutputDto struct {
	ID          string     `json:"id" example:"KKMEa"`
	CreatedAt   time.Time  `json:"created_at" example:"2025-04-30T22:14:07.788043Z" format:"date-time"`
	Content     string     `json:"content" example:"Hello world"`
	ExpireAt    *time.Time `json:"expire_at" example:"2025-04-30T22:14:07.788043Z" format:"date-time" nullable:"true"`
	OwnerUserID *string    `json:"owner_user_id" example:"1234567890" nullable:"true"`
}
