package twitch

import (
	"context"
	"fmt"
	"time"

	"github.com/nicklaw5/helix/v2"
	buscore "github.com/twirapp/twir/libs/bus-core"
	"github.com/twirapp/twir/libs/bus-core/tokens"
	cfg "github.com/twirapp/twir/libs/config"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type userIDCtxKey struct{}

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

func NewAppClient(config cfg.Config, twirBus *buscore.Bus) (*helix.Client, error) {
	return NewAppClientWithContext(context.Background(), config, twirBus)
}

func NewAppClientWithContext(
	ctx context.Context,
	config cfg.Config,
	twirBus *buscore.Bus,
) (
	*helix.Client, error,
) {
	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(
		attribute.String("twitch.client.type", "app"),
		attribute.String("twitch.client.client_id", config.TwitchClientId),
		attribute.String("twitch.client.redirect_uri", config.GetTwitchCallbackUrl()),
	)

	appToken, err := twirBus.Tokens.RequestAppToken.Request(
		ctx,
		struct{}{},
	)
	if err != nil {
		return nil, err
	}

	client, err := helix.NewClientWithContext(
		ctx, &helix.Options{
			ClientID:       config.TwitchClientId,
			ClientSecret:   config.TwitchClientSecret,
			RedirectURI:    config.GetTwitchCallbackUrl(),
			RateLimitFunc:  rateLimitCallback,
			AppAccessToken: appToken.Data.AccessToken,
			HTTPClient:     createHttpClient(),
		},
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}

func NewUserClient(userID string, config cfg.Config, twirBus *buscore.Bus) (
	*helix.Client,
	error,
) {
	return NewUserClientWithContext(context.Background(), userID, config, twirBus)
}

func NewUserClientWithContext(
	ctx context.Context,
	userID string,
	config cfg.Config,
	twirBus *buscore.Bus,
) (*helix.Client, error) {
	ctx = context.WithValue(ctx, userIDCtxKey{}, userID)
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("twitch.client.type", "user"),
		attribute.String("twitch.client.user_id", userID),
		attribute.String("twitch.client.client_id", config.TwitchClientId),
		attribute.String("twitch.client.redirect_uri", config.GetTwitchCallbackUrl()),
	)

	userToken, err := twirBus.Tokens.RequestUserToken.Request(
		ctx,
		tokens.GetUserTokenRequest{UserId: userID},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot request user token from tokens service: %w", err)
	}

	client, err := helix.NewClientWithContext(
		ctx,
		&helix.Options{
			ClientID:        config.TwitchClientId,
			ClientSecret:    config.TwitchClientSecret,
			RedirectURI:     config.GetTwitchCallbackUrl(),
			RateLimitFunc:   rateLimitCallback,
			UserAccessToken: userToken.Data.AccessToken,
			HTTPClient:      createHttpClient(),
		},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot create helix client: %w", err)
	}

	return client, nil
}

func NewBotClient(botID string, config cfg.Config, twirBus *buscore.Bus) (
	*helix.Client,
	error,
) {
	return NewBotClientWithContext(context.Background(), botID, config, twirBus)
}

func NewBotClientWithContext(
	ctx context.Context, botID string, config cfg.Config, twirBus *buscore.Bus,
) (*helix.Client, error) {
	ctx = context.WithValue(ctx, userIDCtxKey{}, botID)
	span := trace.SpanFromContext(ctx)
	span.SetAttributes(
		attribute.String("twitch.client.type", "bot"),
		attribute.String("twitch.client.bot_id", botID),
		attribute.String("twitch.client.client_id", config.TwitchClientId),
		attribute.String("twitch.client.redirect_uri", config.GetTwitchCallbackUrl()),
	)

	botToken, err := twirBus.Tokens.RequestBotToken.Request(
		ctx,
		tokens.GetBotTokenRequest{BotId: botID},
	)
	if err != nil {
		return nil, err
	}

	client, err := helix.NewClientWithContext(
		ctx, &helix.Options{
			ClientID:        config.TwitchClientId,
			ClientSecret:    config.TwitchClientSecret,
			RedirectURI:     config.GetTwitchCallbackUrl(),
			RateLimitFunc:   rateLimitCallback,
			UserAccessToken: botToken.Data.AccessToken,
			HTTPClient:      createHttpClient(),
		},
	)
	if err != nil {
		return nil, err
	}

	return client, nil
}
