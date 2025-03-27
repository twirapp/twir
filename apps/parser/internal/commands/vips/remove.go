package vips

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/nicklaw5/helix/v2"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/twitch"
)

var Remove = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "vips remove",
		Description: null.StringFrom("Remove vip form user."),
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
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		twitchClient, err := twitch.NewUserClient(
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

		if len(parseCtx.Mentions) == 0 {
			return nil, &types.CommandHandlerError{
				Message: "you should tag user with @",
			}
		}

		user := parseCtx.Mentions[0]

		vipResp, err := twitchClient.RemoveChannelVip(
			&helix.RemoveChannelVipParams{
				BroadcasterID: parseCtx.Channel.ID,
				UserID:        user.UserId,
			},
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("cannot remove vip: %s", err),
				Err:     err,
			}
		}
		if vipResp.ErrorMessage != "" {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("cannot remove vip: %s", vipResp.ErrorMessage),
			}
		}

		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf("✅ removed vip from user %s", user.UserLogin),
			},
		}
		
		return result, nil
	},
}
