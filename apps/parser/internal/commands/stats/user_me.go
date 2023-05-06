package stats

import (
	"fmt"

	"github.com/guregu/null"
	"github.com/lib/pq"
	"github.com/satont/tsuwari/apps/parser/internal/types"
	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"
	model "github.com/satont/tsuwari/libs/gomodels"
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
	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: []string{
				fmt.Sprintf(
					"$(%s) used emotes · $(%s) watched · $(%s) messages · $(%s) used points · $(%s) songs requested",
					"user.emotes",
					"user.watched",
					"user.messages",
					"user.usedChannelPoints",
					"user.songs.requested.count",
				),
			},
		}

		return result
	},
}
