package app

import (
	"github.com/avito-tech/go-transaction-manager/trm/v2"
	goredis "github.com/redis/go-redis/v9"
	vkchat "github.com/twirapp/twir/apps/bots/internal/vk"
	cfg "github.com/twirapp/twir/libs/config"
	vkintegrations "github.com/twirapp/twir/libs/integrations/vk"
	oauthvkvideo "github.com/twirapp/twir/libs/oauth/vkvideo"
	tokensrepository "github.com/twirapp/twir/libs/repositories/tokens"
	tokensrepositorypgx "github.com/twirapp/twir/libs/repositories/tokens/datasources/postgres"
	vkvideobotsrepository "github.com/twirapp/twir/libs/repositories/vk_video_bots"
	"go.uber.org/fx"
)

var vkVideoOAuthModule = fx.Options(
	fx.Provide(
		fx.Annotate(tokensrepositorypgx.NewFx, fx.As(new(tokensrepository.Repository))),
		newVKVideoOAuthChatClient,
	),
)

func newVKVideoOAuthChatClient(config cfg.Config, redis *goredis.Client, userRepository tokensrepository.Repository, botRepository vkvideobotsrepository.Repository, transactionRunner trm.Manager, videoChat *vkintegrations.VideoChatClient) (*vkchat.ChatClient, error) {
	userSource, err := oauthvkvideo.NewUserTokenSource(oauthvkvideo.SourceOptions{
		ClientID: config.VKVideoClientID, ClientSecret: config.VKVideoClientSecret,
		RedirectURL: config.GetVkCallbackUrl(), APIBaseURL: config.VKVideoAPIBaseURL,
		AuthBaseURL: config.VKVideoAuthBaseURL, DevAPIBaseURL: config.VKVideoDevAPIBaseURL,
		Redis: redis, CipherKey: config.TokensCipherKey,
	}, userRepository)
	if err != nil {
		return nil, err
	}
	botSource, err := oauthvkvideo.NewSingletonBotTokenSource(oauthvkvideo.SourceOptions{
		ClientID: config.VKVideoClientID, ClientSecret: config.VKVideoClientSecret,
		RedirectURL: config.GetVkVideoBotCallbackUrl(), APIBaseURL: config.VKVideoAPIBaseURL,
		AuthBaseURL: config.VKVideoAuthBaseURL, DevAPIBaseURL: config.VKVideoDevAPIBaseURL,
		Redis: redis, CipherKey: config.TokensCipherKey, TransactionRunner: transactionRunner,
	}, botRepository)
	if err != nil {
		return nil, err
	}
	return vkchat.NewChatClient(userSource, botSource, videoChat), nil
}
