package seventv

import (
	"context"
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	command_arguments "github.com/satont/twir/apps/parser/internal/command-arguments"
	"github.com/satont/twir/apps/parser/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	seventvintegration "github.com/twirapp/twir/libs/integrations/seventv"
)

const emoteOldNameArgName = "oldName"
const emoteNewNameArgName = "newName"

var EmoteRename = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "7tv rename",
		Description: null.StringFrom("Rename 7tv emote in current set"),
		RolesIDS:    pq.StringArray{model.ChannelRoleTypeModerator.String()},
		Module:      "7tv",
		Visible:     true,
		IsReply:     true,
		Aliases:     []string{},
		Enabled:     false,
	},
	SkipToxicityCheck: true,
	Args: []command_arguments.Arg{
		command_arguments.String{
			Name: emoteOldNameArgName,
		},
		command_arguments.String{
			Name: emoteNewNameArgName,
		},
	},
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		client := seventvintegration.NewClient(parseCtx.Services.Config.SevenTvToken)

		sevenTvUser, err := client.GetProfileByTwitchId(ctx, parseCtx.Channel.ID)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("Failed to get 7tv profile: %v", err),
				Err:     err,
			}
		}
		if sevenTvUser.Users.UserByConnection.Style.ActiveEmoteSetId == nil {
			return &types.CommandsHandlerResult{
				Result: []string{
					`❌ No active emote set`,
				},
			}, nil
		}

		oldName := parseCtx.ArgsParser.Get(emoteOldNameArgName).String()
		newName := parseCtx.ArgsParser.Get(emoteNewNameArgName).String()

		var foundEmoteId string
		for _, emote := range sevenTvUser.Users.UserByConnection.Style.ActiveEmoteSet.Emotes.Items {
			if emote.Alias == oldName {
				foundEmoteId = emote.Emote.Id
				break
			}

			if emote.Emote.DefaultName == oldName {
				foundEmoteId = emote.Emote.Id
				break
			}
		}

		if foundEmoteId == "" {
			return &types.CommandsHandlerResult{
				Result: []string{
					fmt.Sprintf(
						`Emote "%s" not found in set %s`,
						oldName,
						sevenTvUser.Users.UserByConnection.Style.ActiveEmoteSet.Name,
					),
				},
			}, nil
		}

		err = client.RenameEmote(
			ctx,
			*sevenTvUser.Users.UserByConnection.Style.ActiveEmoteSetId,
			foundEmoteId,
			newName,
		)
		if err != nil {
			return nil, &types.CommandHandlerError{
				Message: fmt.Sprintf("Failed to rename 7tv emote: %v", err),
				Err:     err,
			}
		}

		return &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(`✅ Emote "%s" renamed to "%s"`, oldName, newName),
			},
		}, nil
	},
}
