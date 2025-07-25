package dota

// import (
// 	"fmt"
// 	"github.com/guregu/null"
// 	"github.com/lib/pq"
// 	model "github.com/twirapp/twir/libs/gomodels"
// 	"strconv"
// 	"strings"

// 	"github.com/twirapp/twir/apps/parser/internal/types"
// 	"github.com/twirapp/twir/apps/parser/pkg/helpers"

// 	variables_cache "github.com/twirapp/twir/apps/parser/internal/variablescache"

// 	"github.com/samber/lo"
// )

// var LgCommand = &types.DefaultCommand{
// 	ChannelsCommands: &model.ChannelsCommands{
// 		Name:        "lg",
// 		Description: null.StringFrom("Players from the latest game"),
// 		RolesIDS:    pq.StringArray{},
// 		Module:      "DOTA",
// 		IsReply:     true,
// 	},
// 	Handler: func(ctx *variables_cache.ExecutionContext) *types.CommandsHandlerResult {
// 		result := &types.CommandsHandlerResult{
// 			Result: make([]string, 0),
// 		}

// 		accounts := GetAccountsByChannelId(ctx.ChannelId)

// 		if len(*accounts) == 0 {
// 			result.Result = append(result.Result, NO_ACCOUNTS)
// 			return result
// 		}

// 		games := GetGames(GetGamesOpts{
// 			Accounts: *accounts,
// 			Take:     lo.ToPtr(2),
// 		})

// 		if games == nil || len(*games) < 2 {
// 			result.Result = append(result.Result, GAME_NOT_FOUND)
// 			return result
// 		}

// 		currGame := lo.FromPtr(games)[0]
// 		prevGame := lo.FromPtr(games)[1]

// 		resultArray := []string{}

// 		for idx, player := range currGame.Players {
// 			owner := helpers.Contains(*accounts, strconv.Itoa(player.AccountId))
// 			if owner {
// 				continue
// 			}

// 			prevPlayer, _, ok := lo.FindIndexOf(prevGame.Players, func(p Player) bool {
// 				return player.AccountId == p.AccountId
// 			})

// 			if !ok {
// 				continue
// 			}

// 			prevHero := GetPlayerHero(prevPlayer.HeroId, nil)
// 			currHero := GetPlayerHero(player.HeroId, lo.ToPtr(idx))
// 			resultArray = append(resultArray, fmt.Sprintf("%s played as %s", currHero, prevHero))
// 		}

// 		if len(resultArray) == 0 {
// 			result.Result = append(result.Result, "Not played with anyone from last game.")
// 			return result
// 		}

// 		result.Result = append(result.Result, strings.Join(resultArray, ", "))
// 		return result
// 	},
// }
