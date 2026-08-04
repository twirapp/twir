package pastebins

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	"github.com/twirapp/twir/apps/api-gql/internal/services/clientinfo"
	"github.com/twirapp/twir/apps/api-gql/internal/services/pastebins"
	config "github.com/twirapp/twir/libs/config"
)

type Dependencies struct {
	API               huma.API
	Config            config.Config
	Service           *pastebins.Service
	Sessions          *auth.Auth
	ClientInfoService *clientinfo.Service
}

type Registration struct{}

func RegisterRoutes(deps Dependencies) Registration {
	newProfile(ProfileOpts{Service: deps.Service, Sessions: deps.Sessions}).Register(deps.API)
	newGetById(GetByIdOpts{Service: deps.Service}).Register(deps.API)
	newCreate(CreateOpts{
		Service: deps.Service, Sessions: deps.Sessions, ClientInfoService: deps.ClientInfoService,
	}).Register(deps.API)
	newDelete(CreateOpts{Service: deps.Service, Sessions: deps.Sessions, ClientInfoService: deps.ClientInfoService}).Register(deps.API)
	return Registration{}
}

type pasteBinOutputDto struct {
	ID          string     `json:"id" example:"KKMEa"`
	CreatedAt   time.Time  `json:"created_at" example:"2025-04-30T22:14:07.788043Z" format:"date-time"`
	Content     string     `json:"content" example:"Hello world"`
	ExpireAt    *time.Time `json:"expire_at" example:"2025-04-30T22:14:07.788043Z" format:"date-time" nullable:"true"`
	OwnerUserID *string    `json:"owner_user_id" example:"1234567890" nullable:"true"`
}
