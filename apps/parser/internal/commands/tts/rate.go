package tts

import (
	"fmt"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	"strconv"
)

var RateCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "tts rate",
		Description: lo.ToPtr("Change tts rate"),
		Permission:  "VIEWER",
		Visible:     false,
		Module:      lo.ToPtr("TTS"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{}
		channelSettings, channelModele := getSettings(ctx.ChannelId, "")

		if channelSettings == nil {
			result.Result = append(result.Result, "TTS is not configured.")
			return result
		}

		userSettings, currentUserModel := getSettings(ctx.ChannelId, ctx.SenderId)

		if ctx.Text == nil {
			result.Result = append(
				result.Result,
				fmt.Sprintf(
					"Global rate: %v | Your rate: %v",
					channelSettings.Rate,
					lo.IfF(userSettings != nil, func() int {
						return userSettings.Rate
					}).Else(channelSettings.Rate),
				))
			return result
		}

		rate, err := strconv.Atoi(*ctx.Text)
		if err != nil {
			result.Result = append(result.Result, "Rate must be a number")
			return result
		}

		if rate < 0 || rate > 100 {
			result.Result = append(result.Result, "Rate must be between 0 and 100")
			return result
		}

		if ctx.ChannelId == ctx.SenderId {
			channelSettings.Rate = rate
			err := updateSettings(channelModele, channelSettings)
			if err != nil {
				fmt.Println(err)
				result.Result = append(result.Result, "Error while updating settings")
				return result
			}
		} else {
			if userSettings == nil {
				_, _, err := createUserSettings(rate, 50, channelSettings.Voice, ctx.ChannelId, ctx.SenderId)
				if err != nil {
					fmt.Println(err)
					result.Result = append(result.Result, "Error while creating settings")
					return result
				}
			} else {
				userSettings.Rate = rate
				err := updateSettings(currentUserModel, userSettings)
				if err != nil {
					fmt.Println(err)
					result.Result = append(result.Result, "Error while updating settings")
					return result
				}
			}
		}

		result.Result = append(result.Result, fmt.Sprintf("Rate changed to %v", rate))

		return result
	},
}
