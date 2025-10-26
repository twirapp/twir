package stats

import (
	"context"
	"fmt"
	"strings"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/twirapp/twir/apps/parser/internal/types"
	"github.com/twirapp/twir/apps/parser/internal/variables/user"

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
		var vars []string

		vars = append(vars, fmt.Sprintf("$(%s)", user.Watched.Name))
		vars = append(vars, fmt.Sprintf("$(%s)", user.Messages.Name))
		vars = append(vars, fmt.Sprintf("$(%s)", user.Emotes.Name))
		vars = append(vars, fmt.Sprintf("$(%s)", user.UsedChannelPoints.Name))
		vars = append(vars, fmt.Sprintf("$(%s)", user.SongsRequested.Name))

		result := &types.CommandsHandlerResult{
			Result: []string{strings.Join(vars, " · ")},
		}

		return result, nil
	},
}
