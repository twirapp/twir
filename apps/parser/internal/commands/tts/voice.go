package tts

import (
	"context"

	"github.com/guregu/null"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"

	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
)

const (
	ttsVoiceArgName = "voice"
)

var VoiceCommand = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "tts voice",
		Description: null.StringFrom("Change tts voice"),
		Module:      "TTS",
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name:     ttsVoiceArgName,
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
				Message: i18n.GetCtx(ctx, locales.Translations.Errors.Generic.GettingChannelSettings),
				Err:     err,
			}
		}

		if channelSettings == nil {
			result.Result = []string{i18n.GetCtx(ctx, locales.Translations.Commands.Tts.Errors.NotConfigured)}
			return result, nil
		}

		userSettings, _, err := parseCtx.Services.TTSService.GetUserSettings(
			ctx,
			parseCtx.Channel.ID,
			parseCtx.Sender.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Errors.Generic.GettingUserSettings),
				Err:     err,
			}
		}

		textArg := parseCtx.ArgsParser.Get(ttsVoiceArgName)

		if textArg == nil {
			result.Result = append(
				result.Result,
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Tts.Info.Voice.
						SetVars(locales.KeysCommandsTtsInfoVoiceVars{GlobalVoice: channelSettings.Voice, UserVoice: lo.IfF(
							userSettings != nil, func() string {
								return userSettings.Voice
							},
						).Else("not setted")}),
				),
			)
			return result, nil
		}

		text := textArg.String()

		wantedVoice, err := parseCtx.Services.TTSService.ValidateVoice(ctx, parseCtx.Channel.ID, text)
		if err != nil {
			result.Result = append(result.Result, err.Error())
			return result, nil
		}

		if parseCtx.Channel.ID == parseCtx.Sender.ID {
			// Update channel settings
			channelSettings.Voice = wantedVoice.Name
			err := parseCtx.Services.TTSService.UpdateChannelSettings(
				ctx,
				parseCtx.Channel.ID,
				channelSettings,
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(ctx, locales.Translations.Errors.Generic.UpdatingSettings),
					Err:     err,
				}
			}
		} else {
			// Update user settings
			if userSettings == nil {
				_, err := parseCtx.Services.TTSService.CreateUserSettings(
					ctx,
					parseCtx.Channel.ID,
					parseCtx.Sender.ID,
					50,
					50,
					wantedVoice.Name,
				)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: i18n.GetCtx(ctx, locales.Translations.Errors.Generic.CreateSettings),
						Err:     err,
					}
				}
			} else {
				userSettings.Voice = wantedVoice.Name
				err := parseCtx.Services.TTSService.UpdateUserSettings(
					ctx,
					parseCtx.Channel.ID,
					parseCtx.Sender.ID,
					userSettings,
				)
				if err != nil {
					return nil, &types.CommandHandlerError{
						Message: i18n.GetCtx(ctx, locales.Translations.Errors.Generic.UpdatingSettings),
						Err:     err,
					}
				}
			}
		}

		result.Result = append(result.Result, i18n.GetCtx(ctx, locales.Translations.Commands.Tts.Info.ChangeVoice.SetVars(locales.KeysCommandsTtsInfoChangeVoiceVars{NewVoice: wantedVoice.Name})))

		parseCtx.Services.TTSCache.Invalidate(ctx, parseCtx.Channel.ID)

		return result, nil
	},
}
