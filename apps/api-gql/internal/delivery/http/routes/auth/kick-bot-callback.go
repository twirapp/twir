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
	buscoreeventsub "github.com/twirapp/twir/libs/bus-core/eventsub"
	"github.com/twirapp/twir/libs/crypto"
	kickbotentity "github.com/twirapp/twir/libs/entities/kick_bot"
	"github.com/twirapp/twir/libs/entities/platform"
	"github.com/twirapp/twir/libs/logger"
	channelplatformsrepo "github.com/twirapp/twir/libs/repositories/channel_platforms"
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

		createdBot, createErr := a.kickBotsRepo.Create(ctx, kickbotsrepo.CreateInput{
			Type:                "DEFAULT",
			AccessToken:         encryptedAccessToken,
			RefreshToken:        encryptedRefreshToken,
			Scopes:              tokens.Scopes,
			ExpiresIn:           tokens.ExpiresIn,
			ObtainmentTimestamp: time.Now().UTC(),
			KickUserID:          internalUser.ID,
			KickUserLogin:       platformUser.Login,
		})
		if createErr != nil {
			a.logger.ErrorContext(ctx, "kick bot callback: failed to create bot", logger.Error(createErr))
			return nil, huma.Error500InternalServerError("Cannot create kick bot", createErr)
		}

		affectedChannelIDs, repairErr := a.assignDefaultKickBotToChannels(ctx, createdBot)
		if repairErr != nil {
			a.logger.ErrorContext(ctx, "kick bot callback: failed to backfill channels with default kick bot", logger.Error(repairErr))
		} else {
			a.publishKickChannelResubscribe(ctx, affectedChannelIDs)
		}
	} else {
		_, err = a.kickBotsRepo.UpdateToken(ctx, existingBot.ID, kickbotsrepo.UpdateTokenInput{
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

		affectedChannelIDs, repairErr := a.assignDefaultKickBotToChannels(ctx, existingBot)
		if repairErr != nil {
			a.logger.ErrorContext(ctx, "kick bot callback: failed to backfill channels with default kick bot", logger.Error(repairErr))
		} else {
			a.publishKickChannelResubscribe(ctx, affectedChannelIDs)
		}
	}

	_ = a.kv.Delete(ctx, kickBotSetupKvPrefix+":"+input.State)

	return httpdelivery.CreateBaseOutputJson(authResponseDto{RedirectTo: "/dashboard/admin-panel"}), nil
}

func (a *Auth) assignDefaultKickBotToChannels(ctx context.Context, kickBot kickbotentity.KickBot) ([]string, error) {
	affectedChannelIDs := make([]string, 0)
	if a.transactionRunner == nil {
		return nil, fmt.Errorf("auth transaction runner is not configured")
	}

	err := a.transactionRunner.Do(ctx, func(txCtx context.Context) error {
		channels, listErr := a.channelsRepo.GetAllByBindingPlatform(txCtx, platform.PlatformKick)
		if listErr != nil {
			return fmt.Errorf("list Kick channel bindings: %w", listErr)
		}

		for _, channel := range channels {
			for _, binding := range channel.Bindings {
				if binding.Platform != platform.PlatformKick {
					continue
				}

				botConfig, needsDefaultBot, configErr := kickBotConfigWithDefault(binding.BotConfig, kickBot.ID)
				if configErr != nil {
					return fmt.Errorf("prepare Kick binding %s bot configuration: %w", binding.ID, configErr)
				}
				if !needsDefaultBot {
					continue
				}

				botUserID := kickBot.KickUserID
				_, updateErr := a.channelPlatformsRepo.Update(txCtx, binding.ID, channelplatformsrepo.UpdateInput{
					UserID:            binding.UserID,
					PlatformChannelID: binding.PlatformChannelID,
					Enabled:           binding.Enabled,
					BotUserID:         &botUserID,
					BotConfig:         botConfig,
				})
				if updateErr != nil {
					return fmt.Errorf("update Kick binding %s: %w", binding.ID, updateErr)
				}

				affectedChannelIDs = append(affectedChannelIDs, channel.ID.String())
			}
		}

		return nil
	})
	if err != nil {
		return nil, err
	}

	return affectedChannelIDs, nil
}

func kickBotConfigWithDefault(botConfig json.RawMessage, kickBotID uuid.UUID) (json.RawMessage, bool, error) {
	config := make(map[string]json.RawMessage)
	if len(botConfig) > 0 {
		if err := json.Unmarshal(botConfig, &config); err != nil {
			return nil, false, err
		}
	}

	if rawKickBotID, ok := config["kick_bot_id"]; ok {
		var existingKickBotID string
		if err := json.Unmarshal(rawKickBotID, &existingKickBotID); err == nil && existingKickBotID != "" {
			return botConfig, false, nil
		}
	}

	encodedKickBotID, err := json.Marshal(kickBotID.String())
	if err != nil {
		return nil, false, fmt.Errorf("encode kick bot ID: %w", err)
	}
	config["kick_bot_id"] = encodedKickBotID

	updatedConfig, err := json.Marshal(config)
	if err != nil {
		return nil, false, fmt.Errorf("encode Kick binding configuration: %w", err)
	}

	return updatedConfig, true, nil
}

func (a *Auth) publishKickChannelResubscribe(ctx context.Context, channelIDs []string) {
	for _, channelID := range channelIDs {
		if a.eventSubPublisher == nil {
			return
		}
		if err := a.eventSubPublisher.Publish(
			ctx,
			buscoreeventsub.EventsubSubscribeToAllEventsRequest{ChannelID: channelID, Platform: platform.PlatformKick},
		); err != nil {
			a.logger.ErrorContext(ctx, "cannot publish kick eventsub subscribe", logger.Error(err), slog.String("channel_id", channelID))
		}
	}
}
