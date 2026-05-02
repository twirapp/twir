package auth

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/danielgtaylor/huma/v2"
	"github.com/google/uuid"
	httpdelivery "github.com/twirapp/twir/apps/api-gql/internal/delivery/http"
	"github.com/twirapp/twir/libs/crypto"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	kickbotsrepo "github.com/twirapp/twir/libs/repositories/kick_bots"
	usersrepo "github.com/twirapp/twir/libs/repositories/users"
	usersmodel "github.com/twirapp/twir/libs/repositories/users/model"
)

const kickBotSetupKvPrefix = "kick_bot_setup"

type kickBotCallbackInput struct {
	Code  string
	State string
}

type kickBotSetupState struct {
	CodeVerifier string `json:"code_verifier"`
	AdminUserID  string `json:"admin_user_id"`
}

func (a *Auth) handleKickBotCallback(
	ctx context.Context,
	input kickBotCallbackInput,
) (*httpdelivery.BaseOutputJson[authResponseDto], error) {
	if a.kickProvider == nil {
		return nil, huma.Error500InternalServerError("Kick provider is not configured", fmt.Errorf("kick provider is nil"))
	}

	stateBytes, err := a.kv.Get(ctx, kickBotSetupKvPrefix+":"+input.State).Bytes()
	if err != nil {
		return nil, huma.Error400BadRequest("Invalid or expired state", fmt.Errorf("state not found in kv"))
	}

	var setupState kickBotSetupState
	if err := json.Unmarshal(stateBytes, &setupState); err != nil {
		return nil, huma.Error400BadRequest("Invalid state", err)
	}

	a.logger.InfoContext(ctx, "kick bot callback: processing", slog.String("admin_user_id", setupState.AdminUserID))

	tokens, err := a.kickProvider.ExchangeBotSetupCode(ctx, input.Code, setupState.CodeVerifier)
	if err != nil {
		a.logger.ErrorContext(ctx, "kick bot callback: failed to exchange code", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot exchange code", err)
	}

	platformUser, err := a.kickProvider.GetUser(ctx, tokens.AccessToken)
	if err != nil {
		a.logger.ErrorContext(ctx, "kick bot callback: failed to get user from kick", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot get user data from kick", err)
	}

	encryptedAccessToken, err := crypto.Encrypt(tokens.AccessToken, a.config.TokensCipherKey)
	if err != nil {
		a.logger.ErrorContext(ctx, "kick bot callback: failed to encrypt access token", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot encrypt access token", err)
	}

	encryptedRefreshToken, err := crypto.Encrypt(tokens.RefreshToken, a.config.TokensCipherKey)
	if err != nil {
		a.logger.ErrorContext(ctx, "kick bot callback: failed to encrypt refresh token", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot encrypt refresh token", err)
	}

	existingBot, err := a.kickBotsRepo.GetDefault(ctx)
	if err != nil && !errors.Is(err, kickbotsrepo.ErrNotFound) {
		a.logger.ErrorContext(ctx, "kick bot callback: failed to get default bot", logger.Error(err))
		return nil, huma.Error500InternalServerError("Cannot get default bot", err)
	}

	if errors.Is(err, kickbotsrepo.ErrNotFound) {
		internalUser, userErr := a.usersRepo.GetByPlatformID(ctx, platform.PlatformKick, platformUser.ID)
		if userErr != nil && !errors.Is(userErr, usersmodel.ErrNotFound) {
			a.logger.ErrorContext(ctx, "kick bot callback: failed to get internal user", logger.Error(userErr))
			return nil, huma.Error500InternalServerError("Cannot get internal user", userErr)
		}

		if errors.Is(userErr, usersmodel.ErrNotFound) {
			internalUser, userErr = a.usersRepo.Create(ctx, usersrepo.CreateInput{
				Platform:   platform.PlatformKick,
				PlatformID: platformUser.ID,
			})
			if userErr != nil {
				a.logger.ErrorContext(ctx, "kick bot callback: failed to create internal user", logger.Error(userErr))
				return nil, huma.Error500InternalServerError("Cannot create internal user", userErr)
			}
		}

		_, err = a.kickBotsRepo.Create(ctx, kickbotsrepo.CreateInput{
			Type:                "DEFAULT",
			AccessToken:         encryptedAccessToken,
			RefreshToken:        encryptedRefreshToken,
			Scopes:              tokens.Scopes,
			ExpiresIn:           tokens.ExpiresIn,
			ObtainmentTimestamp: time.Now().UTC(),
			KickUserID:          internalUser.ID,
			KickUserLogin:       platformUser.Login,
		})
		if err != nil {
			a.logger.ErrorContext(ctx, "kick bot callback: failed to create bot", logger.Error(err))
			return nil, huma.Error500InternalServerError("Cannot create kick bot", err)
		}
	} else {
		botID, err := uuid.Parse(existingBot.ID)
		if err != nil {
			a.logger.ErrorContext(ctx, "kick bot callback: failed to parse bot id", logger.Error(err))
			return nil, huma.Error500InternalServerError("Invalid bot id", err)
		}
		_, err = a.kickBotsRepo.UpdateToken(ctx, botID, kickbotsrepo.UpdateTokenInput{
			AccessToken:         encryptedAccessToken,
			RefreshToken:        encryptedRefreshToken,
			Scopes:              tokens.Scopes,
			ExpiresIn:           tokens.ExpiresIn,
			ObtainmentTimestamp: time.Now().UTC(),
		})
		if err != nil {
			a.logger.ErrorContext(ctx, "kick bot callback: failed to update bot token", logger.Error(err))
			return nil, huma.Error500InternalServerError("Cannot update kick bot token", err)
		}
	}

	_ = a.kv.Delete(ctx, kickBotSetupKvPrefix+":"+input.State)

	return httpdelivery.CreateBaseOutputJson(authResponseDto{RedirectTo: "/dashboard/admin-panel"}), nil
}
