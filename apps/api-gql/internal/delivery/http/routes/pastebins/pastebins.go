package pastebins

import (
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/twirapp/twir/apps/api-gql/internal/auth"
	httpbase "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/apps/api-gql/internal/services/pastebins"
	config "github.com/twirapp/twir/libs/config"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	Api      huma.API
	Config   config.Config
	Service  *pastebins.Service
	Sessions *auth.Auth
}

var FxModule = fx.Provide(
	httpbase.AsFxRoute(newProfile),
	httpbase.AsFxRoute(newGetById),
	httpbase.AsFxRoute(newCreate),
	httpbase.AsFxRoute(newDelete),
)

type pasteBinOutputDto struct {
	ID          string     `json:"id" example:"KKMEa"`
	CreatedAt   time.Time  `json:"created_at" example:"2025-04-30T22:14:07.788043Z" format:"date-time"`
	Content     string     `json:"content" example:"Hello world"`
	ExpireAt    *time.Time `json:"expire_at" example:"2025-04-30T22:14:07.788043Z" format:"date-time" nullable:"true"`
	OwnerUserID *string    `json:"owner_user_id" example:"1234567890" nullable:"true"`
}
