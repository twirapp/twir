package tts

import (
	"context"
	"fmt"
	"strings"

	"github.com/guregu/null"
	model "github.com/satont/tsuwari/libs/gomodels"

	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	"go.uber.org/zap"
)

var VoiceCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts voice",
		Description: null.StringFrom("Change tts voice"),
		Module:      "TTS",
		IsReply:     true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}
		channelSettings, channelModel := getSettings(ctx, parseCtx.Services.Gorm, parseCtx.Channel.ID, "")

		if channelSettings == nil {
			result.Result = append(result.Result, "TTS is not configured.")
			return result
		}

		userSettings, currentUserModel := getSettings(
			ctx,
			parseCtx.Services.Gorm,
			parseCtx.Channel.ID,
			parseCtx.Sender.ID,
		)

		if parseCtx.Text == nil {
			result.Result = append(
				result.Result,
				fmt.Sprintf(
					"Global voice: %s | Your voice: %s",
					channelSettings.Voice,
					lo.IfF(userSettings != nil, func() string {
						return userSettings.Voice
					}).Else("not setted"),
				))
			return result
		}

		voices := getVoices(ctx, parseCtx.Services.Config)
		if len(voices) == 0 {
			result.Result = append(result.Result, "No voices found")
			return result
		}

		wantedVoice, ok := lo.Find(voices, func(item Voice) bool {
			return item.Name == strings.ToLower(*parseCtx.Text)
		})

		if !ok {
			result.Result = append(result.Result, fmt.Sprintf("Voice %s not found", *parseCtx.Text))
			return result
		}

		_, isDisallowed := lo.Find(channelSettings.DisallowedVoices, func(item string) bool {
			return item == wantedVoice.Name
		})

		if isDisallowed {
			result.Result = append(
				result.Result,
				fmt.Sprintf("Voice %s is disallowed for usage", wantedVoice.Name),
			)
			return result
		}

		if parseCtx.Channel.ID == parseCtx.Sender.ID {
			channelSettings.Voice = wantedVoice.Name
			err := updateSettings(ctx, parseCtx.Services.Gorm, channelModel, channelSettings)
			if err != nil {
				zap.S().Error(err)
				result.Result = append(result.Result, "Error while updating settings")
				return result
			}
		} else {
			if userSettings == nil {
				_, _, err := createUserSettings(
					ctx,
					parseCtx.Services.Gorm,
					50,
					50,
					wantedVoice.Name,
					parseCtx.Channel.ID,
					parseCtx.Sender.ID,
				)
				if err != nil {
					zap.S().Error(err)
					result.Result = append(result.Result, "Error while creating settings")
					return result
				}
			} else {

				userSettings.Voice = wantedVoice.Name
				err := updateSettings(ctx, parseCtx.Services.Gorm, currentUserModel, userSettings)
				if err != nil {
					zap.S().Error(err)
					result.Result = append(result.Result, "Error while updating settings")
					return result
				}
			}
		}

		result.Result = append(result.Result, fmt.Sprintf("Voice changed to %s", wantedVoice.Name))

		return result
	},
}
