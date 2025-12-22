package vips

import (
	"context"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	scheduledvipsentity "github.com/twirapp/twir/libs/entities/scheduled_vips"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
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
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(ctx, locales.Translations.Commands.Vips.Hints.User)
			},
		},
		command_arguments.VariadicString{
			Name: "unvip_in",
			HintFunc: func(ctx context.Context) string {
				return i18n.GetCtx(ctx, locales.Translations.Commands.Vips.Hints.UnvipIn)
			},
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

		var (
			unvipAt   *time.Time
			unvipType *scheduledvipsentity.RemoveType
		)
		unvipArg := parseCtx.ArgsParser.Get("unvip_in").String()

		if unvipArg == "stream_end" {
			unvipType = lo.ToPtr(scheduledvipsentity.RemoveTypeStreamEnd)
		} else {
			unvipType = lo.ToPtr(scheduledvipsentity.RemoveTypeTime)
			duration, err := str2duration.ParseDuration(unvipArg)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(ctx, locales.Translations.Commands.Vips.Errors.InvalidDuration),
					Err:     err,
				}
			}

			newUnvipAt := time.Now().Add(duration)
			unvipAt = &newUnvipAt
		}

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
		if vip.IsNil() {
			err := parseCtx.Services.ScheduledVipsRepo.Create(
				ctx,
				scheduledvipsrepository.CreateInput{
					ChannelID:  parseCtx.Channel.ID,
					UserID:     user.UserId,
					RemoveType: unvipType,
					RemoveAt:   unvipAt,
				},
			)
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: i18n.GetCtx(
						ctx,
						locales.Translations.Commands.Vips.Errors.CannotCreateScheduledInDb,
					),
					Err: err,
				}
			}
		} else {
			err = parseCtx.Services.ScheduledVipsRepo.Update(
				ctx,
				vip.ID,
				scheduledvipsrepository.UpdateInput{
					RemoveAt:   unvipAt,
					RemoveType: unvipType,
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

		localeVars := locales.KeysCommandsVipsErrorsUpdatedVars{
			UserName: user.UserName,
		}
		if unvipType != nil && *unvipType == scheduledvipsentity.RemoveTypeTime && unvipAt != nil {
			localeVars.EndTime = unvipAt.Format("2006-01-02 15:04:05")
		} else if unvipType != nil && *unvipType == scheduledvipsentity.RemoveTypeStreamEnd {
			localeVars.EndTime = "stream end"
		} else {
			localeVars.EndTime = "never"
		}

		result := &types.CommandsHandlerResult{
			Result: []string{
				i18n.GetCtx(
					ctx,
					locales.Translations.Commands.Vips.Errors.Updated.SetVars(localeVars),
				),
			},
		}

		return result, nil
	},
}
