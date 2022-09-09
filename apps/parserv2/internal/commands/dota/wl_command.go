package dota

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"
	"tsuwari/parser/internal/types"
	"tsuwari/parser/internal/variables/stream"

	model "tsuwari/parser/internal/models"

	variables_cache "tsuwari/parser/internal/variablescache"

	"github.com/samber/lo"
)

type MatchResponse struct {
	Result struct {
		MatchId    int                       `json:"match_id"`
		Players    []model.DotaMatchesPlayer `json:"players"`
		RadiantWin bool                      `json:"radiant_win"`
		GameMode   int32                     `json:"game_mode"`
	}
}

var WlCommand = types.DefaultCommand{
	Command: types.Command{
		Name:        "wl",
		Description: lo.ToPtr("Score for played games on stream"),
		Permission:  "VIEWER",
		Visible:     true,
		Module:      lo.ToPtr("DOTA"),
	},
	Handler: func(ctx variables_cache.ExecutionContext) []string {
		streamString, err := ctx.Services.Redis.Get(context.TODO(), "streams:"+ctx.ChannelId).
			Result()
		streamData := stream.HelixStream{}

		if err != nil || streamString == "" {
			return []string{"Stream not found"}
		}

		err = json.Unmarshal([]byte(streamString), &streamData)

		if err != nil {
			fmt.Println(err)
			return []string{"Something wen't wrong on getting stream."}
		}

		accounts := GetAccountsByChannelId(ctx.Services.Db, ctx.ChannelId)

		if accounts == nil || len(*accounts) == 0 {
			return []string{NO_ACCOUNTS}
		}

		dbGames := []model.DotaMatchWithRelation{}

		err = ctx.Services.Db.
			Table("dota_matches").
			Where(
				`"dota_matches"."players" && ARRAY[?]::int[] 
				AND "startedAt" >= ? 
				AND "lobby_type" IN ?
				`,
				*accounts,
				streamData.StartedAt.Add(-time.Minute*10),
				[2]int{0, 7},
			).
			Order(`"startedAt" DESC`).
			Joins("GameMode").
			Joins("Result").
			Find(&dbGames).Error

		if err != nil {
			fmt.Println(err)
			return []string{"Something wen't wrong on fetching games."}
		}

		dbGames = lo.Filter(dbGames, func(g model.DotaMatchWithRelation, _ int) bool {
			return !lo.SomeBy(g.PlayersHeroes, func(id int64) bool {
				return id == 0
			})
		})

		if len(dbGames) == 0 {
			return []string{"Games not played on the stream."}
		}

		matchesData := []model.DotaMatchesResults{}

		for _, match := range dbGames {
			if match.Result != nil {
				matchesData = append(matchesData, *match.Result)
			}
		}

		gamesForRequest := lo.Filter(dbGames, func(g model.DotaMatchWithRelation, _ int) bool {
			return g.Result == nil
		})

		gamesForRequestWg := sync.WaitGroup{}
		for _, game := range gamesForRequest {
			gamesForRequestWg.Add(1)

			go func(game model.DotaMatchWithRelation) {
				defer gamesForRequestWg.Done()
				data := MatchResponse{}
				r, err := ApiInstance.
					R().
					SetQueryParam("match_id", game.MatchID).
					SetQueryParam("key", "2B5C2069282D28E79B60B494489E31C5").
					SetResult(&data).
					Get("http://api.steampowered.com/IDOTA2Match_570/GetMatchDetails/v1")

				if err != nil || r.StatusCode != 200 {
					fmt.Println(r.StatusCode, err)
					return
				}
				dataForPush := model.DotaMatchesResults{
					MatchID:    strconv.Itoa(data.Result.MatchId),
					Players:    data.Result.Players,
					RadiantWin: data.Result.RadiantWin,
					GameMode:   data.Result.GameMode,
				}

				matchesData = append(matchesData, dataForPush)

				err = ctx.Services.Db.Create(&model.DotaMatchesResults{
					MatchID:    strconv.Itoa(data.Result.MatchId),
					Players:    []byte(data.Result.Players),
					RadiantWin: data.Result.RadiantWin,
					GameMode:   data.Result.GameMode,
				}).Error
				if err != nil {
					fmt.Println(err)
				}
			}(game)
		}

		gamesForRequestWg.Wait()

		return []string{""}
	},
}
