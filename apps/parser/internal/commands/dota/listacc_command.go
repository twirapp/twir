package dota

// import (
// 	"github.com/guregu/null"
// 	"github.com/lib/pq"
// 	model "github.com/twirapp/twir/libs/gomodels"
// 	"strings"

// 	"github.com/twirapp/twir/apps/parser/internal/types"
// 	variables_cache "github.com/twirapp/twir/apps/parser/internal/variablescache"
// )

// var ListAccCommand = &types.DefaultCommand{
// 	ChannelsCommands: &model.ChannelsCommands{
// 		Name:        "dota listacc",
// 		Description: null.StringFrom("List of added dota accounts"),
// 		RolesIDS:    pq.StringArray{model.ChannelRoleTypeBroadcaster.String()},
// 		Module:      "DOTA",
// 		IsReply:     true,
// 	},
// 	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
// 		result := &types.CommandsHandlerResult{
// 			Result: make([]string, 0),
// 		}

// 		accounts := GetAccountsByChannelId(ctx.ChannelId)

// 		if accounts == nil || len(*accounts) == 0 {
// 			result.Result = append(result.Result, NO_ACCOUNTS)
// 			return result
// 		}

// 		result.Result = append(result.Result, strings.Join(*accounts, ", "))
// 		return result
// 	},
// }
