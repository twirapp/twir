package dota2

import (
	"context"
	"errors"
	"fmt"

	"github.com/lib/pq"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
)

const mmrArgName = "mmr"

var Mmr = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:    "mmr",
		Module:  "DOTA",
		IsReply: true,
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		if _, err := requireDotaSettings(
			ctx,
			parseCtx,
			func(settings model.ChannelsDotaSettingsCommands) bool { return settings.Mmr },
		); err != nil {
			return nil, err
		}

		data, err := getDotaData(ctx, parseCtx)
		if err != nil {
			return nil, err
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Dota.Outputs.Mmr.SetVars(
						locales.KeysCommandsDotaOutputsMmrVars{
							Mmr:   data.Mmr,
							Medal: medalName(ctx, medalForMMR(data.Mmr)),
						},
					),
				),
			},
		}, nil
	},
}

var MmrSet = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:     "mmr set",
		RolesIDS: pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:   "DOTA",
		IsReply:  true,
	},
	Args: []command_arguments.Arg{
		command_arguments.Int{
			Name: mmrArgName,
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Hints.Mmr)
			},
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		if _, err := requireDotaSettings(ctx, parseCtx, nil); err != nil {
			return nil, err
		}

		mmrArg := parseCtx.ArgsParser.Get(mmrArgName)
		if mmrArg == nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Errors.MmrRequired),
			}
		}

		mmr := mmrArg.Int()
		if mmr < 0 {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Errors.MmrNegative),
			}
		}

		updateResult := parseCtx.Services.Gorm.WithContext(ctx).
			Model(&model.ChannelsDotaSettings{}).
			Where("channel_id = ?", parseCtx.Channel.DBChannelID).
			UpdateColumn("mmr", mmr)
		if updateResult.Error != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Errors.UpdateMmr),
				Err:     fmt.Errorf("update Dota MMR: %w", updateResult.Error),
			}
		}
		if updateResult.RowsAffected == 0 {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Dota.Errors.SettingsNotFound),
				Err:     errors.New("Dota MMR update affected no rows"),
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Dota.Success.MmrUpdated.SetVars(
						locales.KeysCommandsDotaSuccessMmrUpdatedVars{Mmr: mmr},
					),
				),
			},
		}, nil
	},
}
