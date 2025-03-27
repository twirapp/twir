package vips

import (
	"context"
	"fmt"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	scheduledvipmodel "github.com/twirapp/twir/libs/repositories/scheduled_vips/model"
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
			Hint: "@username",
		},
		command_arguments.VariadicString{
			Name: "unvip_in",
			Hint: "time in, example: 1w5d1m5s",
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		if len(parseCtx.Mentions) == 0 {
			return nil, &types.CommandHandlerError{
				Message: "you should tag user with @",
			}
		}

		channelTwitchClient, err := twitch.NewUserClient(
			parseCtx.Channel.ID,
			*parseCtx.Services.Config,
			parseCtx.Services.GrpcClients.Tokens,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create broadcaster twitch client",
				Err:     err,
			}
		}

		unvipArg := parseCtx.ArgsParser.Get("unvip_in").String()
		duration, err := str2duration.ParseDuration(unvipArg)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "invalid duration",
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
				Message: "cannot get scheduled vip",
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
					Message: "cannot create scheduled vip",
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
					Message: "cannot update scheduled vip",
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
				fmt.Sprintf(
					"✅ updated vip for user %s new expriation time %s",
					user.UserName,
					newUnvipAt.Format("2006-01-02 15:04:05"),
				),
			},
		}

		return result, nil
	},
}
