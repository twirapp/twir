package vips

import (
	"context"
	"fmt"
	"slices"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/locales"
	"github.com/twirapp/twir/libs/cache/twitch"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/i18n"
	scheduledvipmodel "github.com/twirapp/twir/libs/repositories/scheduled_vips/model"
)

var List = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "vips list",
		Description: null.StringFrom("List of scheduled vips."),
		RolesIDS: pq.StringArray{
			model.ChannelRoleTypeModerator.String(),
		},
		Module:  "VIPS",
		Aliases: pq.StringArray{},
		Visible: true,
		IsReply: true,
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		scheduledVips, err := parseCtx.Services.ScheduledVipsRepo.GetManyByChannelID(
			ctx,
			parseCtx.Channel.ID,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					parseCtx.Services.I18n,
					locales.Translations.Commands.Vips.CannotGetListFromDb,
				),
				Err: err,
			}
		}

		result := &types.CommandsHandlerResult{
			Result: []string{},
		}
		if len(scheduledVips) == 0 {
			result.Result = []string{
				i18n.GetCtx(
					ctx,
					parseCtx.Services.I18n,
					locales.Translations.Commands.Vips.NoScheduledVips,
				),
			}
			return result, nil
		}

		usersIds := make([]string, 0, len(scheduledVips))
		for _, vip := range scheduledVips {
			usersIds = append(usersIds, vip.UserID)
		}

		twitchUsers, err := parseCtx.Services.CacheTwitchClient.GetUsersByIds(ctx, usersIds)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: i18n.GetCtx(
					ctx,
					parseCtx.Services.I18n,
					locales.Translations.Errors.Generic.CannotFindUsersTwitch,
				),
				Err: err,
			}
		}

		scheduledVipsWithUsers := make([]scheduledVip, 0, len(scheduledVips))
		for _, vip := range scheduledVips {
			twitchUser, twitchUserOk := lo.Find(
				twitchUsers, func(user twitch.TwitchUser) bool {
					return user.ID == vip.UserID
				},
			)

			v := scheduledVip{
				UserID: vip.UserID,
				Vip:    vip,
			}
			if twitchUserOk {
				v.User = &twitchUser.User
			}

			scheduledVipsWithUsers = append(
				scheduledVipsWithUsers,
				v,
			)
		}

		slices.SortFunc(
			scheduledVipsWithUsers, func(a, b scheduledVip) int {
				return a.Vip.RemoveAt.Compare(*b.Vip.RemoveAt)
			},
		)

		var resultedString strings.Builder

		for _, vip := range scheduledVipsWithUsers {
			if resultedString.Len() > 0 {
				resultedString.WriteString(" · ")
			}

			resultedString.WriteString(vip.User.DisplayName)
			resultedString.WriteString(" ")
			resultedString.WriteString(
				fmt.Sprintf(
					"(%s)",
					vip.Vip.RemoveAt.Format("2006-01-02 15:04:05"),
				),
			)
		}

		result.Result = []string{resultedString.String()}

		return result, nil
	},
}

type scheduledVip struct {
	UserID string
	User   *helix.User
	Vip    scheduledvipmodel.ScheduledVip
}
