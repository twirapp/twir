package dota

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/satont/tsuwari/apps/parser/internal/types"

	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/grpc/generated/dota"

	variables_cache "github.com/satont/tsuwari/apps/parser/internal/variablescache"

	"github.com/samber/lo"
	uuid "github.com/satori/go.uuid"
)

var GmCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "gm",
		Description: lo.ToPtr("Game medals from current game."),
		Permission:  "BROADCASTER",
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
			Take:     lo.ToPtr(1),
		})

		if games == nil || len(*games) == 0 {
			result.Result = append(result.Result, GAME_NOT_FOUND)
			return result
		}

		game := lo.FromPtr(games)[0]

		cards := []model.DotaMatchesCards{}
		if game.PlayersCards != nil {
			cards = *game.PlayersCards
		}

		playersForGet := lo.Filter(game.Players, func(p Player, _ int) bool {
			return !lo.SomeBy(cards, func(card model.DotaMatchesCards) bool {
				id, _ := strconv.Atoi(card.AccountID)

				return id == p.AccountId
			})
		})

		wg := sync.WaitGroup{}
		lock := sync.Mutex{}
		for _, player := range playersForGet {
			wg.Add(1)
			go func(player Player) {
				defer wg.Done()

				req, err := ctx.Services.DotaGrpc.GetPlayerCard(context.Background(), &dota.GetPlayerCardRequest{
					AccountId: int64(player.AccountId),
				})
				if err != nil {
					return
				}

				lock.Lock()

				card := model.DotaMatchesCards{
					ID:        uuid.NewV4().String(),
					MatchID:   game.ID,
					AccountID: req.AccountId,
					// RankTier: lo.If(data.RankTier != nil, sql.NullInt64{
					// 	Int64: *data.RankTier,
					// 	Valid: true,
					// }).Else(sql.NullInt64{}),
					// LeaderboardRank: lo.If(data.LeaderboardRank != nil, sql.NullInt64{
					// 	Int64: *data.LeaderboardRank,
					// 	Valid: true,
					// }).Else(sql.NullInt64{}),
				}
				if req.RankTier != nil {
					card.RankTier = sql.NullInt64{
						Int64: *req.RankTier,
						Valid: true,
					}
				}
				if req.LeaderboardRank != nil {
					card.LeaderboardRank = sql.NullInt64{
						Int64: *req.LeaderboardRank,
						Valid: true,
					}
				}
				cards = append(cards, card)
				lock.Unlock()

				err = ctx.Services.Db.Create(&card).Error
				if err != nil {
					fmt.Println(err)
				}
			}(player)
		}

		wg.Wait()

		resultArray := [10]string{}

		for _, card := range cards {
			player, idx, ok := lo.FindIndexOf(game.Players, func(p Player) bool {
				id, _ := strconv.Atoi(card.AccountID)
				return p.AccountId == id
			})

			if !ok {
				continue
			}

			hero := GetPlayerHero(player.HeroId, lo.ToPtr(idx))
			medal, ok := lo.Find(DotaMedals, func(m Medal) bool {
				return m.Tier == int(card.RankTier.Int64)
			})

			if !ok {
				medal = Medal{
					Tier: 0,
					Name: "Unknown",
				}
			}

			rank := lo.
				If(
					card.LeaderboardRank.Valid,
					fmt.Sprintf("#%v", card.LeaderboardRank.Int64),
				).
				Else(fmt.Sprintf("%s", medal.Name))

			resultArray[idx] = fmt.Sprintf("%s %s", hero, rank)
		}

		result.Result = append(result.Result, strings.Join(resultArray[:], ", "))
		return result
	},
}
