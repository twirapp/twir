package tts

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/samber/lo"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
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

		channelSettings, _, err := parseCtx.Services.TTSService.GetChannelSettings(
			ctx,
			parseCtx.Channel.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "error while getting channel settings",
				Err:     err,
			}
		}

		if channelSettings == nil {
			result.Result = []string{"TTS is not configured for this channel"}
			return result, nil
		}

		volumeArg := parseCtx.ArgsParser.Get(ttsVolumeArgName)
		if volumeArg == nil {
			result.Result = []string{fmt.Sprintf("Current volume: %v", channelSettings.Volume)}
			return result, nil
		}

		volume := volumeArg.Int()
		channelSettings.Volume = volume

		err = parseCtx.Services.TTSService.UpdateChannelSettings(
			ctx,
			parseCtx.Channel.ID,
			channelSettings,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "error while updating settings",
				Err:     err,
			}
		}

		result.Result = []string{fmt.Sprintf("TTS volume changed to %v", volume)}

		parseCtx.Services.TTSCache.Invalidate(ctx, parseCtx.Channel.ID)

		return result, nil
	},
}
