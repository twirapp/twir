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

	"github.com/lib/pq"
	"github.com/samber/lo"
)

type MatchPlayer struct {
	AccountId  int `json:"account_id"`
	TeamNumber int `json:"team_number"`
	HeroId     int `json:"hero_id"`
}

type MatchResult struct {
	Error      *string       `json:"error"`
	MatchID    int           `json:"match_id"`
	Players    []MatchPlayer `json:"players"`
	RadiantWin bool          `json:"radiant_win"`
	GameMode   int32         `json:"game_mode"`
}

type MatchResponse struct {
	Result MatchResult `json:"result"`
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
			return []string{"Something went wrong on getting stream."}
		}

		accounts := GetAccountsByChannelId(ctx.Services.Db, ctx.ChannelId)

		if accounts == nil || len(*accounts) == 0 {
			return []string{NO_ACCOUNTS}
		}

		dbGames := []model.DotaMatchWithRelation{}
		intAccounts := lo.Map(*accounts, func(a string, _ int) int {
			acc, _ := strconv.Atoi(a)
			return acc
		})

		err = ctx.Services.Db.
			Table("dota_matches").
			Where(
				`"dota_matches"."players" && ?
				AND "startedAt" >= ? 
				AND "lobby_type" IN ?
				`,
				pq.Array(intAccounts),
				streamData.StartedAt.Add(-time.Minute*10),
				[2]int{0, 7},
			).
			Order(`"startedAt" DESC`).
			Joins("GameMode").
			Joins("Result").
			Find(&dbGames).Error

		if err != nil {
			fmt.Println(err)
			return []string{"Something went wrong on fetching games."}
		}

		dbGames = lo.Filter(dbGames, func(g model.DotaMatchWithRelation, _ int) bool {
			return !lo.SomeBy(g.PlayersHeroes, func(id int64) bool {
				return id == 0
			})
		})

		if len(dbGames) == 0 {
			return []string{"Games not played on the stream."}
		}

		matchesData := []MatchResult{}

		for _, match := range dbGames {
			if match.Result != nil {
				matchId, _ := strconv.Atoi(match.MatchID)
				players := []MatchPlayer{}
				_ = json.Unmarshal([]byte(match.Result.Players), &players)
				result := MatchResult{
					Error:      nil,
					MatchID:    matchId,
					RadiantWin: match.Result.RadiantWin,
					GameMode:   match.Result.GameMode,
					Players:    players,
				}
				matchesData = append(matchesData, result)
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
					SetQueryParams(map[string]string{
						"match_id": game.MatchID,
						"key":      "2B5C2069282D28E79B60B494489E31C5",
					}).
					SetResult(&data).
					Get("http://api.steampowered.com/IDOTA2Match_570/GetMatchDetails/v1")

				if err != nil || r.StatusCode != 200 || data.Result.Error != nil {
					fmt.Println(r.StatusCode, err, *data.Result.Error)
					return
				}

				dataForPush := MatchResult{
					MatchID:    data.Result.MatchID,
					Players:    data.Result.Players,
					RadiantWin: data.Result.RadiantWin,
					GameMode:   data.Result.GameMode,
				}

				matchesData = append(matchesData, dataForPush)
				players, _ := json.Marshal(data.Result.Players)
				err = ctx.Services.Db.
					Table("dota_matches_results").
					Create(map[string]interface{}{
						"match_id":    strconv.Itoa(data.Result.MatchID),
						"players":     string(players),
						"radiant_win": data.Result.RadiantWin,
						"game_mode":   data.Result.GameMode,
					}).
					Error
				if err != nil {
					fmt.Println(game.MatchID, err)
				}
			}(game)
		}

		gamesForRequestWg.Wait()

		return []string{""}
	},
}
