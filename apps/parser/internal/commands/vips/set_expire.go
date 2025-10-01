package vips

import (
	"context"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	scheduledvipmodel "github.com/twirapp/twir/libs/repositories/scheduled_vips/model"
	"github.com/twirapp/twir/libs/twitch"
	"github.com/xhit/go-str2duration/v2"
)

var SetExpire = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "vips setexpire",
		Description: null.StringFrom("Set new expiration time for vip."),
		RolesIDS: pq.StringArray{
			model.ChannelRoleTypeModerator.String(),
		},
		Module:  "VIPS",
		Aliases: pq.StringArray{},
		Visible: true,
		IsReply: true,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: "user",
			Hint: i18n.Get(locales.Translations.Commands.Vips.Hints.User),
		},
		command_arguments.VariadicString{
			Name: "unvip_in",
			Hint: i18n.Get(locales.Translations.Commands.Vips.Hints.UnvipIn),
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		if len(parseCtx.Mentions) == 0 {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Errors.Generic.ShouldMentionWithAt),
			}
		}

		channelTwitchClient, err := twitch.NewUserClient(
			parseCtx.Channel.ID,
			*parseCtx.Services.Config,
			parseCtx.Services.Bus,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Errors.Generic.BroadcasterClient),
				Err:     err,
			}
		}

		unvipArg := parseCtx.ArgsParser.Get("unvip_in").String()
		duration, err := str2duration.ParseDuration(unvipArg)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Vips.Errors.InvalidDuration),
				Err:     err,
			}
		}

		newUnvipAt := time.Now().Add(duration)

		user := parseCtx.Mentions[0]

		vip, err := parseCtx.Services.ScheduledVipsRepo.GetByUserAndChannelID(
			ctx,
			user.UserId,
			parseCtx.Channel.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(ctx, locales.Translations.Commands.Vips.Errors.CannotGetListFromDb),
				Err:     err,
			}
		}
		if vip == scheduledvipmodel.Nil {
			err := parseCtx.Services.ScheduledVipsRepo.Create(
				ctx,
				scheduledvipsrepository.CreateInput{
					ChannelID: parseCtx.Channel.ID,
					UserID:    user.UserId,
				},
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(ctx, locales.Translations.Commands.Vips.Errors.CannotCreateScheduledInDb),
					Err:     err,
				}
			}
		} else {
			err = parseCtx.Services.ScheduledVipsRepo.Update(
				ctx,
				vip.ID,
				scheduledvipsrepository.UpdateInput{
					RemoveAt: &newUnvipAt,
				},
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(ctx, locales.Translations.Commands.Vips.Errors.CannotUpdate),
					Err:     err,
				}
			}
		}

		// ignore error
		channelTwitchClient.AddChannelVip(
			&helix.AddChannelVipParams{
				BroadcasterID: parseCtx.Channel.ID,
				UserID:        user.UserId,
			},
		)

		result := &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Vips.Errors.Updated.SetVars(
						locales.KeysCommandsVipsErrorsUpdatedVars{
							UserName: user.UserName,
							EndTime:  newUnvipAt.Format("2006-01-02 15:04:05"),
						},
					),
				),
			},
		}

		return result, nil
	},
}
