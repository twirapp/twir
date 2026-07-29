package app

import (
	goredis "github.com/redis/go-redis/v9"
	cfg "github.com/twirapp/twir/libs/config"
	oauthvkvideo "github.com/twirapp/twir/libs/oauth/vkvideo"
	tokensrepository "github.com/twirapp/twir/libs/repositories/tokens"
	tokensrepositorypgx "github.com/twirapp/twir/libs/repositories/tokens/datasources/postgres"
	"go.uber.org/fx"
)

var vkVideoOAuthModule = fx.Options(
	fx.Provide(
		fx.Annotate(tokensrepositorypgx.NewFx, fx.As(new(tokensrepository.Repository))),
		newVKVideoUserTokenSource,
	),
)

func newVKVideoUserTokenSource(config cfg.Config, redis *goredis.Client, repository tokensrepository.Repository) (oauthvkvideo.UserTokenSource, error) {
	return oauthvkvideo.NewUserTokenSource(oauthvkvideo.SourceOptions{
		ClientID: config.VKVideoClientID, ClientSecret: config.VKVideoClientSecret,
		RedirectURL: config.GetVkCallbackUrl(), APIBaseURL: config.VKVideoAPIBaseURL,
		AuthBaseURL: config.VKVideoAuthBaseURL, DevAPIBaseURL: config.VKVideoDevAPIBaseURL,
		Redis: redis, CipherKey: config.TokensCipherKey,
	}, repository)
}
