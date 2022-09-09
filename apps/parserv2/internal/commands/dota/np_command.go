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
	Handler: func(ctx variables_cache.ExecutionContext) []string {
		accounts := GetAccountsByChannelId(ctx.Services.Db, ctx.ChannelId)

		if accounts == nil || len(*accounts) == 0 {
			return []string{NO_ACCOUNTS}
		}

		games := *GetGames(GetGamesOpts{
			Db:       ctx.Services.Db,
			Accounts: *accounts,
			Take:     lo.ToPtr(1),
			Redis:    ctx.Services.Redis,
		})

		if games == nil || len(games) == 0 {
			return []string{"Game not found."}
		}

		game := games[0]
		avgMmr := lo.
			If(
				game.GameMode == 22 && game.LobbyType == 7,
				fmt.Sprintf(" (%v)mmr", game.AvarageMmr)).
			Else("")

		gameMode := GetGameModeById(game.GameMode)
		modeName := lo.If(gameMode == nil, "Unknown").Else(gameMode.Name)

		return []string{fmt.Sprintf("%s%s", modeName, avgMmr)}
	},
}
