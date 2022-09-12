package dota

import (
	"fmt"
	"strconv"
	"strings"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/pkg/helpers"

	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

var LgCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "lg",
		Description: lo.ToPtr("Players from the latest game"),
		Permission:  "VIEWER",
		Visible:     true,
		Module:      lo.ToPtr("DOTA"),
	},
	Handler: func(ctx variables_cache.ExecutionContext) []string {
		accounts := GetAccountsByChannelId(ctx.Services.Db, ctx.ChannelId)

		if len(*accounts) == 0 {
			return []string{NO_ACCOUNTS}
		}

		games := GetGames(GetGamesOpts{
			Db:       ctx.Services.Db,
			Redis:    ctx.Services.Redis,
			Accounts: *accounts,
			Take:     lo.ToPtr(2),
		})

		if games == nil || len(*games) < 2 {
			return []string{GAME_NOT_FOUND}
		}

		currGame := lo.FromPtr(games)[0]
		prevGame := lo.FromPtr(games)[1]

		result := []string{}

		for idx, player := range currGame.Players {
			owner := helpers.Contains(*accounts, strconv.Itoa(player.AccountId))
			if owner {
				continue
			}

			prevPlayer, _, ok := lo.FindIndexOf(prevGame.Players, func(p Player) bool {
				return player.AccountId == p.AccountId
			})

			if !ok {
				continue
			}

			prevHero := GetPlayerHero(prevPlayer.HeroId, nil)
			currHero := GetPlayerHero(player.HeroId, lo.ToPtr(idx))
			result = append(result, fmt.Sprintf("%s played as %s", currHero, prevHero))
		}

		if len(result) == 0 {
			return []string{"Not played with anyone from last game."}
		}

		return []string{strings.Join(result, ", ")}
	},
}
