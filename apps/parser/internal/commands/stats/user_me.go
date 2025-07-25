package stats

import (
	"context"
	"fmt"
	"strings"

	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/variables/user"

	"github.com/guregu/null"
	"github.com/lib/pq"

	model "github.com/twirapp/twir/libs/gomodels"
)

var UserMe = &types.DefaultCommand{
	ChannelsCommands: &model.ChannelsCommands{
		Name:        "me",
		Description: null.StringFrom("Prints user statistic."),
		RolesIDS:    pq.StringArray{},
		Module:      "STATS",
		Aliases:     pq.StringArray{"stats"},
		Visible:     true,
		IsReply:     true,
	},
	SkipToxicityCheck: true,
	Handler: func(ctx context.Context, parseCtx *types.ParseContext) (
		*types.CommandsHandlerResult,
		error,
	) {
		var slice []string

		slice = append(slice, fmt.Sprintf("$(%s) watched", user.Watched.Name))
		slice = append(slice, fmt.Sprintf("$(%s) messages", user.Messages.Name))
		slice = append(slice, fmt.Sprintf("$(%s) used emotes", user.Emotes.Name))
		slice = append(slice, fmt.Sprintf("$(%s) used points", user.UsedChannelPoints.Name))
		slice = append(slice, fmt.Sprintf("$(%s) songs requestes", user.SongsRequested.Name))

		result := &types.CommandsHandlerResult{
			Result: []string{strings.Join(slice, " · ")},
		}

		return result, nil
	},
}
