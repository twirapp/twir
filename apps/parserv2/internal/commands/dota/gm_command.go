package dota

import (
	"database/sql"
	"fmt"
	"strconv"
	"sync"
	"time"
	"tsuwari/parser/internal/types"

	model "tsuwari/parser/internal/models"

	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/golang/protobuf/proto"
	"github.com/samber/lo"
	dotanats "github.com/satont/tsuwari/nats/dota"
	uuid "github.com/satori/go.uuid"
)

var GmCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "gm",
		Description: lo.ToPtr("Game medals from current game."),
		Permission:  "BROADCASTER",
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

		if games == nil || len(*games) == 0 {
			return []string{GAME_NOT_FOUND}
		}

		game := lo.FromPtr(games)[0]

		fmt.Println(game.PlayersCards)
		cards := *game.PlayersCards
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

				bytes, _ := proto.Marshal(&dotanats.GetPlayerCardRequest{
					AccountId: int64(player.AccountId),
				})

				response, err := ctx.Services.Nats.Request(
					"dota.getProfileCard",
					bytes,
					5*time.Second,
				)
				if err != nil {
					return
				}

				data := dotanats.GetPlayerCardResponse{}
				if err = proto.Unmarshal(response.Data, &data); err != nil {
					return
				}
				lock.Lock()
				card := model.DotaMatchesCards{
					ID:        uuid.NewV4().String(),
					MatchID:   game.ID,
					AccountID: data.AccountId,
					RankTier: sql.NullInt64{
						Int64: *data.RankTier,
						Valid: data.RankTier != nil,
					},
					LeaderboardRank: sql.NullInt64{
						Int64: *data.LeaderboardRank,
						Valid: data.LeaderboardRank != nil,
					},
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

		return []string{}
	},
}
