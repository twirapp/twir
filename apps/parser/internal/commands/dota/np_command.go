package dota

// import (
// 	"fmt"
// 	"github.com/guregu/null"
// 	"github.com/lib/pq"
// 	model "github.com/twirapp/twir/libs/gomodels"

// 	"github.com/twirapp/twir/apps/parser/internal/types"

// 	variables_cache "github.com/twirapp/twir/apps/parser/internal/variablescache"

// 	"github.com/samber/lo"
// )

// var NpAccCommand = &types.DefaultCommand{
// 	ChannelsCommands: &model.ChannelsCommands{
// 		Name:        "np",
// 		Description: null.StringFrom("Notable players from current dota game"),
// 		RolesIDS:    pq.StringArray{},
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

// 		games := GetGames(GetGamesOpts{
// 			Accounts: *accounts,
// 			Take:     lo.ToPtr(1),
// 		})

// 		if games == nil || len(*games) == 0 {
// 			result.Result = append(result.Result, GAME_NOT_FOUND)
// 			return result
// 		}

// 		game := lo.FromPtr(games)[0]
// 		avgMmr := lo.
// 			If(
// 				game.GameMode == 22 && game.LobbyType == 7,
// 				fmt.Sprintf(" (%v)mmr", game.AvarageMmr)).
// 			Else("")

// 		gameMode := GetGameModeById(game.GameMode)
// 		modeName := lo.If(gameMode == nil, "Unknown").Else(gameMode.Name)

// 		result.Result = append(result.Result, fmt.Sprintf("%s%s", modeName, avgMmr))
// 		return result
// 	},
// }
