package tts

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"strings"
)

var VoiceCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "tts voice",
		Description: lo.ToPtr("Change tts voice"),
		Permission:  "VIEWER",
		Visible:     true,
		Module:      lo.ToPtr("TTS"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		//webSocketsGrpc := do.MustInvoke[websockets.WebsocketClient](di.Provider)
		result := &types.CommandsHandlerResult{}
		channelSettings, channelModel := getSettings(ctx.ChannelId, "")

		if channelSettings == nil {
			result.Result = append(result.Result, "TTS is not configured.")
			return result
		}

		userSettings, currentUserModel := getSettings(ctx.ChannelId, ctx.SenderId)

		if ctx.Text == nil {
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

		voices := getVoices()
		if len(voices) == 0 {
			result.Result = append(result.Result, "No voices found")
			return result
		}

		wantedVoice, ok := lo.Find(voices, func(item voice) bool {
			return item.Name == strings.ToLower(*ctx.Text)
		})

		if !ok {
			result.Result = append(result.Result, fmt.Sprintf("Voice %s not found", *ctx.Text))
			return result
		}

		if ctx.ChannelId == ctx.SenderId {
			channelSettings.Voice = wantedVoice.Name
			err := updateSettings(channelModel, channelSettings)
			if err != nil {
				fmt.Println(err)
				result.Result = append(result.Result, "Error while updating settings")
				return result
			}
		} else {
			if userSettings == nil {
				_, _, err := createUserSettings(50, 50, wantedVoice.Name, ctx.ChannelId, ctx.SenderId)
				if err != nil {
					fmt.Println(err)
					result.Result = append(result.Result, "Error while creating settings")
					return result
				}
			} else {

				userSettings.Voice = wantedVoice.Name
				err := updateSettings(currentUserModel, userSettings)
				if err != nil {
					fmt.Println(err)
					result.Result = append(result.Result, "Error while updating settings")
					return result
				}
			}
		}

		result.Result = append(result.Result, fmt.Sprintf("Voice changed to %s", wantedVoice.Name))

		return result
	},
}
