package vips

import (
	"context"
	"fmt"
	"time"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	command_arguments "github.com/twirapp/twir/apps/parser/internal/command-arguments"
	"github.com/twirapp/twir/apps/parser/internal/types"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/twitch"
	scheduledvipsrepository "github.com/twirapp/twir/libs/repositories/scheduled_vips"
	"github.com/xhit/go-str2duration/v2"
)

var Add = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "vips add",
		Description: null.StringFrom("Add vip to user, can be scheduled for unvip."),
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
			Name:     "unvip_in",
			Hint:     "time in, example: 1w5d1m5s",
			Optional: true,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		var unvipAt *time.Time
		unvipArg := parseCtx.ArgsParser.Get("unvip_in")
		if unvipArg != nil {
			duration, err := str2duration.ParseDuration(unvipArg.String())
			if err != nil {
				return nil, &types.CommandHandlerError{
					Message: "invalid duration",
					Err:     err,
				}
			}

			newUnvipAt := time.Now().Add(duration)
			unvipAt = &newUnvipAt
		}

		twitchClient, err := twitch.NewUserClient(
			parseCtx.Channel.ID,
			*parseCtx.Services.Config,
			parseCtx.Services.Bus,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot create broadcaster twitch client",
				Err:     err,
			}
		}

		if len(parseCtx.Mentions) == 0 {
			return nil, &types.CommandHandlerError{
				Message: "you should tag user with @",
			}
		}

		user := parseCtx.Mentions[0]

		var dbUser model.Users
		if err := parseCtx.Services.Gorm.
			WithContext(ctx).
			Where(`"id" = ?`, user.UserId).
			Preload("Stats", `"channelId" = ? AND "userId" = ?`, parseCtx.Channel.ID, user.UserId).
			First(&dbUser).Error; err != nil {
			return nil, &types.CommandHandlerError{
				Message: "cannot find user in database",
				Err:     err,
			}
		}

		if dbUser.Stats != nil && (dbUser.Stats.IsMod || dbUser.Stats.IsVip) {
			return nil, &types.CommandHandlerError{
				Message: "user already vip or moderator",
			}
		}

		trErr := parseCtx.Services.TrmManager.Do(
			ctx,
			func(trCtx context.Context) error {
				vipResp, err := twitchClient.AddChannelVip(
					&helix.AddChannelVipParams{
						BroadcasterID: parseCtx.Channel.ID,
						UserID:        user.UserId,
					},
				)
				if err != nil {
					return &types.CommandHandlerError{
						Message: fmt.Sprintf("cannot add vip: %s", err),
						Err:     err,
					}
				}
				if vipResp.ErrorMessage != "" {
					return &types.CommandHandlerError{
						Message: fmt.Sprintf("cannot add vip: %s", vipResp.ErrorMessage),
					}
				}

				if unvipAt == nil {
					return nil
				}

				err = parseCtx.Services.ScheduledVipsRepo.Create(
					trCtx,
					scheduledvipsrepository.CreateInput{
						ChannelID: parseCtx.Channel.ID,
						UserID:    user.UserId,
						RemoveAt:  unvipAt,
					},
				)
				if err != nil {
					return &types.CommandHandlerError{
						Message: "cannot create scheduled vip in database",
						Err:     err,
					}
				}

				return nil
			},
		)
		if trErr != nil {
			return nil, &types.CommandHandlerError{
				Message: trErr.Error(),
				Err:     trErr,
			}
		}

		result := &types.CommandsHandlerResult{
			Result: []string{},
		}

		if unvipAt != nil {
			result.Result = append(
				result.Result,
				fmt.Sprintf(
					"✅ added vip to %s, will be removed at %s",
					user.UserName,
					unvipAt.Format("2006-01-02 15:04:05"),
				),
			)
		} else {
			result.Result = append(result.Result, fmt.Sprintf("✅ added vip to %s", user.UserName))
		}

		return result, nil
	},
}
