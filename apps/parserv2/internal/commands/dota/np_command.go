package dota

import (
	"fmt"
	"strings"
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

		games := GetGames(GetGamesOpts{
			Db:       ctx.Services.Db,
			Accounts: *accounts,
			Take:     lo.ToPtr(1),
			Redis:    ctx.Services.Redis,
		})

		if games == nil || len(*games) == 0 {
			return []string{"Game not found."}
		}

		result := lo.Map(*games, func(game Game, _ int) string {
			avgMmr := lo.
				If(
					game.GameMode == 22 && game.LobbyType == 7,
					fmt.Sprintf(" (%v)mmr", game.AvarageMmr)).
				Else("")

			var modeName string
			gameMode := GetGameModeById(game.GameMode)

			if gameMode == nil {
				modeName = "Unknown"
			} else {
				modeName = gameMode.Name
			}

			return fmt.Sprintf("%s%s", modeName, avgMmr)
		})

		return []string{strings.Join(result, " | ")}
	},
}
