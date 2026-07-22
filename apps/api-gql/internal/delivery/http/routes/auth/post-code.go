package auth

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/danielgtaylor/huma/v2"
	"github.com/nicklaw5/helix/v2"
	httpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	platformentity "github.com/twirapp/twir/libs/entities/platform"
)

type authBody struct {
	Code  string `json:"code" minLength:"20" required:"true"`
	State string `json:"state" required:"true"`
}

type authResponseDto struct {
	RedirectTo string `json:"redirect_to"`
}

func (a *Auth) handleAuthPostCode(
	ctx context.Context,
	input authBody,
) (*httpdelivery.BaseOutputJson[authResponseDto], error) {
	result, err := a.completePlatformCode(ctx, platformCodeInput{
		Platform: platformentity.PlatformTwitch,
		Code:     input.Code,
		State:    input.State,
	})
	if err != nil {
		return nil, a.platformAuthHTTPError(err)
	}

	return a.completeTwitchAuthResponse(ctx, result)
}

func (a *Auth) completeTwitchAuthResponse(
	ctx context.Context,
	result platformCodeResult,
) (*httpdelivery.BaseOutputJson[authResponseDto], error) {
	if result.PlatformUser == nil {
		return nil, huma.Error500InternalServerError("Cannot get user data from twitch", fmt.Errorf("twitch user not found"))
	}

	if err := a.sessions.SetSessionTwitchUser(ctx, helix.User{
		ID:              result.PlatformUser.ID,
		Login:           result.PlatformUser.Login,
		DisplayName:     result.PlatformUser.DisplayName,
		ProfileImageURL: result.PlatformUser.Avatar,
	}); err != nil {
		return nil, huma.Error500InternalServerError("Cannot set twitch user", err)
	}

	a.logger.InfoContext(ctx, "twitch auth: completed successfully", slog.String("channel_id", result.AuthResult.Channel.ID.String()), slog.String("user_id", result.AuthResult.SessionUserID.String()))

	return httpdelivery.CreateBaseOutputJson(authResponseDto{RedirectTo: result.RedirectTo}), nil
}
