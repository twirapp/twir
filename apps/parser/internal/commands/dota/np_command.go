package dota

import (
	"fmt"
	"tsuwari/parser/internal/types"

	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var NpAccCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "np",
		Description: lo.ToPtr("Notable players from current dota game"),
		Permission:  "VIEWER",
		Visible:     true,
		Module:      lo.ToPtr("DOTA"),
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		accounts := GetAccountsByChannelId(ctx.Services.Db, ctx.ChannelId)

		if accounts == nil || len(*accounts) == 0 {
			result.Result = append(result.Result, NO_ACCOUNTS)
			return result
		}

		games := GetGames(GetGamesOpts{
			Db:       ctx.Services.Db,
			Accounts: *accounts,
			Take:     lo.ToPtr(1),
			Redis:    ctx.Services.Redis,
		})

		if games == nil || len(*games) == 0 {
			result.Result = append(result.Result, GAME_NOT_FOUND)
			return result
		}

		game := lo.FromPtr(games)[0]
		avgMmr := lo.
			If(
				game.GameMode == 22 && game.LobbyType == 7,
				fmt.Sprintf(" (%v)mmr", game.AvarageMmr)).
			Else("")

		gameMode := GetGameModeById(game.GameMode)
		modeName := lo.If(gameMode == nil, "Unknown").Else(gameMode.Name)

		result.Result = append(result.Result, fmt.Sprintf("%s%s", modeName, avgMmr))
		return result
	},
}
