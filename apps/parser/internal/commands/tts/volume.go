package tts

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
)

const (
	ttsVolumeArgName = "volume"
)

var VolumeCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts volume",
		Description: null.StringFrom("Change tts volume. This is not per user, it's global for the channel."),
		Module:      "TTS",
		IsReply:     true,
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeBroadcaster.String()},
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.Int{
			Name:     ttsVolumeArgName,
			Min:      lo.ToPtr(1),
			Max:      lo.ToPtr(100),
			Optional: true,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		result := &types.CommandsHandlerResult{}
		channelSettings, channelEntity := getSettings(
			ctx,
			parseCtx.Services.Gorm,
			parseCtx.Channel.ID,
			"",
		)

		if channelSettings == nil {
			return result, nil
		}

		volumeArg := parseCtx.ArgsParser.Get(ttsVolumeArgName)

		if volumeArg == nil {
			result.Result = append(
				result.Result,
				fmt.Sprintf("Global volume: %v", channelSettings.Volume),
			)
			return result, nil
		}

		volume := volumeArg.Int()

		channelSettings.Volume = volume
		err := updateSettings(ctx, parseCtx.Services.Gorm, channelEntity, channelSettings)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "error while updating settings",
				Err:     err,
			}
		}

		result.Result = append(result.Result, fmt.Sprintf("Volume changed to %v", volume))

		return result, nil
	},
}
