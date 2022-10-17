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
		Visible:     false,
		Module:      lo.ToPtr("DOTA"),
		IsReply:     true,
	},
	Handler: func(ctx variables_cache.ExecutionContext) *types.CommandsHandlerResult {
		result := &types.CommandsHandlerResult{
			Result: make([]string, 0),
		}

		accounts := GetAccountsByChannelId(ctx.Services.Db, ctx.ChannelId)

		if len(*accounts) == 0 {
			result.Result = append(result.Result, NO_ACCOUNTS)
			return result
		}

		games := GetGames(GetGamesOpts{
			Db:       ctx.Services.Db,
			Redis:    ctx.Services.Redis,
			Accounts: *accounts,
			Take:     lo.ToPtr(2),
		})

		if games == nil || len(*games) < 2 {
			result.Result = append(result.Result, GAME_NOT_FOUND)
			return result
		}

		currGame := lo.FromPtr(games)[0]
		prevGame := lo.FromPtr(games)[1]

		resultArray := []string{}

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
			resultArray = append(resultArray, fmt.Sprintf("%s played as %s", currHero, prevHero))
		}

		if len(resultArray) == 0 {
			result.Result = append(result.Result, "Not played with anyone from last game.")
			return result
		}

		result.Result = append(result.Result, strings.Join(resultArray, ", "))
		return result
	},
}
