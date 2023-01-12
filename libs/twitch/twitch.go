package twitch

import (
	"context"
	cfg "github.com/satont/tsuwari/libs/config"
	"github.com/satont/tsuwari/libs/grpc/generated/tokens"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
	"time"

	helix "github.com/satont/go-helix/v2"
)

func rateLimitCallback(lastResponse *helix.Response) error {
	if lastResponse.GetRateLimitRemaining() > 0 {
		return nil
	}

	var reset64 int64
	reset64 = int64(lastResponse.GetRateLimitReset())

	currentTime := time.Now().UTC().Unix()

	if currentTime < reset64 {
		timeDiff := time.Duration(reset64 - currentTime)
		if timeDiff > 0 {
			time.Sleep(timeDiff * time.Second)
		}
	}

	return nil
}

func NewAppClient(config cfg.Config, tokensGrpc tokens.TokensClient) (*helix.Client, error) {
	appToken, err := tokensGrpc.RequestAppToken(
		context.Background(),
		&emptypb.Empty{},
		grpc.WaitForReady(true),
	)
	if err != nil {
		return nil, err
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:       config.TwitchClientId,
		ClientSecret:   config.TwitchClientSecret,
		RedirectURI:    config.TwitchCallbackUrl,
		RateLimitFunc:  rateLimitCallback,
		AppAccessToken: appToken.AccessToken,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewUserClient(userID string, config cfg.Config, tokensGrpc tokens.TokensClient) (*helix.Client, error) {
	userToken, err := tokensGrpc.RequestUserToken(
		context.Background(),
		&tokens.GetUserTokenRequest{UserId: userID},
		grpc.WaitForReady(true),
	)
	if err != nil {
		return nil, err
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:        config.TwitchClientId,
		ClientSecret:    config.TwitchClientSecret,
		RedirectURI:     config.TwitchCallbackUrl,
		RateLimitFunc:   rateLimitCallback,
		UserAccessToken: userToken.AccessToken,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewBotClient(botID string, config cfg.Config, tokensGrpc tokens.TokensClient) (*helix.Client, error) {
	botToken, err := tokensGrpc.RequestBotToken(
		context.Background(),
		&tokens.GetBotTokenRequest{BotId: botID},
		grpc.WaitForReady(true),
	)
	if err != nil {
		return nil, err
	}

	client, err := helix.NewClient(&helix.Options{
		ClientID:        config.TwitchClientId,
		ClientSecret:    config.TwitchClientSecret,
		RedirectURI:     config.TwitchCallbackUrl,
		RateLimitFunc:   rateLimitCallback,
		UserAccessToken: botToken.AccessToken,
	})
	if err != nil {
		return nil, err
	}

	return client, nil
}
