package dota

// import (
// 	"context"
// 	"fmt"
// 	"github.com/go-redis/redis/v9"
// 	"github.com/samber/do"
// 	"github.com/twirapp/twir/apps/parser/internal/di"
// 	"go.uber.org/zap"
// 	"strings"
// 	"time"

// 	model "github.com/twirapp/twir/libs/gomodels"

// 	req "github.com/imroc/req/v3"
// 	"github.com/samber/lo"
// 	"gorm.io/gorm"
// )

// const (
// 	GAME_NOT_FOUND   = "Game not found"
// 	NO_ACCOUNTS      = "No accounts added."
// 	WRONG_ACCOUNT_ID = "Wrong account id."
// )

// var ApiInstance = req.C().
// 	SetCommonRetryCount(2)

// var Colors = [...]string{
// 	"Blue",
// 	"Teal",
// 	"Purple",
// 	"Yellow",
// 	"Orange",
// 	"Pink",
// 	"Gray",
// 	"Light Blue",
// 	"Green",
// 	"Brown",
// }

// func GetPlayerHero(heroId int, index *int) string {
// 	const Unknown = "Unknown"

// 	if index != nil && *index > len(Colors)-1 {
// 		return Unknown
// 	}

// 	if heroId == 0 && index != nil {
// 		color := Colors[*index]

// 		if color != "" {
// 			return color
// 		}

// 		return Unknown
// 	}

// 	if heroId == 0 && index == nil {
// 		return Unknown
// 	}

// 	hero, ok := lo.Find(DotaHeroes, func(h Hero) bool {
// 		return h.ID == heroId
// 	})

// 	if !ok {
// 		return Unknown
// 	}

// 	if hero.ShortName != nil {
// 		return *hero.ShortName
// 	}

// 	return hero.LocalizedName
// }

// type GetGamesOpts struct {
// 	Accounts []string
// 	Take     *int
// }

// func GetGames(opts GetGamesOpts) *[]Game {
// 	db := do.MustInvoke[gorm.DB](di.Provider)
// 	redisClient := do.MustInvoke[redis.Client](di.Provider)

// 	ctx := context.TODO()
// 	rpsCount := 0

// 	for _, acc := range opts.Accounts {
// 		rps, err := redisClient.MGet(ctx, fmt.Sprintf("dotaRps:%v", acc)).Result()
// 		if err != nil {
// 			continue
// 		}
// 		rps = lo.Filter(rps, func(r interface{}, _ int) bool {
// 			return r != nil
// 		})

// 		rpsCount = rpsCount + len(rps)
// 	}

// 	if rpsCount == 0 {
// 		return nil
// 	}

// 	cachedGamesCount := 0

// 	for _, acc := range opts.Accounts {
// 		games, err := redisClient.MGet(ctx, fmt.Sprintf("dotaMatches:%v", acc)).Result()
// 		if err != nil {
// 			continue
// 		}
// 		games = lo.Filter(games, func(r interface{}, _ int) bool {
// 			return r != nil
// 		})

// 		cachedGamesCount = cachedGamesCount + len(games)
// 	}

// 	if cachedGamesCount == 0 {
// 		return nil
// 	}

// 	dbGames := []model.DotaMatchWithRelation{}

// 	if opts.Take == nil {
// 		opts.Take = lo.ToPtr(1)
// 	}

// 	err := db.
// 		Table("dota_matches").
// 		Where(`players && ?`, fmt.Sprintf("{%s}", strings.Join(opts.Accounts, ","))).
// 		Order(`"startedAt" DESC`).
// 		Joins("GameMode").
// 		Find(&dbGames).Error
// 	if err != nil {
// 		fmt.Println("GetGames:", err)
// 		return nil
// 	}

// 	for i, v := range dbGames {
// 		players := []model.DotaMatchesCards{}
// 		err := db.
// 			Where("match_id = ?", v.ID).
// 			Find(&players).Error
// 		if err == nil {
// 			dbGames[i].PlayersCards = &players
// 		}
// 	}

// 	if len(dbGames) == 0 {
// 		return nil
// 	}

// 	mappedGames := lo.Map(dbGames, func(game model.DotaMatchWithRelation, _ int) Game {
// 		g := Game{
// 			ID:                        game.ID,
// 			ActivateTime:              game.StartedAt,
// 			LobbyType:                 game.LobbyType,
// 			GameMode:                  game.GameModeID,
// 			AvarageMmr:                game.AvarageMmr,
// 			WeekedTourneyBracketRound: &game.WeekendTourneyBracketRound.String,
// 			WeekedTourneySkillLevel:   &game.WeekendTourneySkillLevel.String,
// 			MatchId:                   &game.MatchID,
// 			LobbyId:                   game.LobbyID,
// 			PlayersCards:              game.PlayersCards,
// 		}

// 		g.Players = lo.Map(game.Players, func(p int64, i int) Player {
// 			return Player{
// 				AccountId: int(p),
// 				HeroId:    int(game.PlayersHeroes[i]),
// 			}
// 		})

// 		return g
// 	})

// 	return &mappedGames
// }

// type Player struct {
// 	AccountId int `json:"account_id"`
// 	HeroId    int `json:"hero_id"`
// }

// type Game struct {
// 	ID                        string                    `json:"id"`
// 	ActivateTime              time.Time                 `json:"activate_time"`
// 	LobbyType                 int32                     `json:"lobby_type"`
// 	GameMode                  int32                     `json:"game_mode"`
// 	AvarageMmr                int32                     `json:"avarage_mmr"`
// 	Players                   []Player                  `json:"players"`
// 	WeekedTourneyBracketRound *string                   `json:"weeked_tourney_bracket_round"`
// 	WeekedTourneySkillLevel   *string                   `json:"weeked_tourney_skill_level"`
// 	MatchId                   *string                   `json:"match_id"`
// 	LobbyId                   string                    `json:"lobby_id"`
// 	PlayersCards              *[]model.DotaMatchesCards `json:"playersCards"`
// }

// func GetAccountsByChannelId(channelId string) *[]string {
// 	db := do.MustInvoke[gorm.DB](di.Provider)

// 	accounts := []model.ChannelsDotaAccounts{}
// 	err := db.
// 		Table("channels_dota_accounts").
// 		Where(`"channelId" = ?`, channelId).
// 		Find(&accounts).
// 		Error
// 	if err != nil {
// 		zap.S().Error(err)
// 		return nil
// 	}

// 	mappedAccounts := lo.Map(accounts, func(a model.ChannelsDotaAccounts, _ int) string {
// 		return a.ID
// 	})

// 	return &mappedAccounts
// }

// func GetGameModeById(id int32) *GameMode {
// 	mode, ok := lo.Find(DotaGameModes, func(mode GameMode) bool {
// 		return mode.ID == int(id)
// 	})

// 	if !ok {
// 		return nil
// 	}

// 	return &mode
// }
