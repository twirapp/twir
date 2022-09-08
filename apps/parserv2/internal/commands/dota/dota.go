package dota

import (
	"context"
	"fmt"
	"time"
	model "tsuwari/parser/internal/models"

	"github.com/go-redis/redis/v9"
	req "github.com/imroc/req/v3"
	"github.com/samber/lo"
	"gorm.io/gorm"
)

const (
	GAME_NOT_FOUND = "Game not found"
	NO_ACCOUNTS    = "You do not have any accounts added."
)

var ApiInstance = req.C().
	SetCommonRetryCount(2)

var Colors = [...]string{
	"Blue",
	"Teal",
	"Purple",
	"Yellow",
	"Orange",
	"Ping",
	"Gray",
	"Light Blue",
	"Green",
	"Brown",
}

func GetPlayerHero(heroId int, index *int) string {
	const Unknown = "Unknown"

	if index != nil && *index > len(Colors)-1 {
		return Unknown
	}

	if heroId == 0 && index != nil {
		color := Colors[*index]

		if color != "" {
			return color
		}

		return Unknown
	}

	if heroId == 0 && index == nil {
		return Unknown
	}

	hero, ok := lo.Find(DotaHeroes, func(h Hero) bool {
		return h.ID == heroId
	})

	if !ok {
		return Unknown
	}

	if hero.ShortName != nil {
		return *hero.ShortName
	}

	return hero.LocalizedName
}

type GetGamesOpts struct {
	Db       *gorm.DB
	Accounts []string
	Take     *int
	Redis    *redis.Conn
}

func GetGames(opts GetGamesOpts) *[]Game {
	ctx := context.TODO()
	rpsCount := 0

	for _, acc := range opts.Accounts {
		rps, err := opts.Redis.MGet(ctx, fmt.Sprintf("dotaRps:%v", acc)).Result()
		if err != nil {
			continue
		}
		rps = lo.Filter(rps, func(r interface{}, _ int) bool {
			return r != nil
		})

		rpsCount = rpsCount + len(rps)
	}

	if rpsCount == 0 {
		return nil
	}

	cachedGamesCount := 0

	for _, acc := range opts.Accounts {
		games, err := opts.Redis.MGet(ctx, fmt.Sprintf("dotaMatches:%v", acc)).Result()
		if err != nil {
			continue
		}
		games = lo.Filter(games, func(r interface{}, _ int) bool {
			return r != nil
		})

		cachedGamesCount = cachedGamesCount + len(games)
	}

	if cachedGamesCount == 0 {
		return nil
	}

	dbGames := []model.DotaMatchWithRelation{}

	if opts.Take == nil {
		opts.Take = lo.ToPtr(1)
	}

	err := opts.Db.
		Where("players && ?", opts.Accounts).
		Order(`"startedAt" DESC`).
		Joins("GameMode").
		Joins("PlayersCards").
		Find(&dbGames).
		Take(*opts.Take).
		Error

	if err != nil {
		fmt.Println(err)
		return nil
	}

	if len(dbGames) == 0 {
		return nil
	}

	mappedGames := lo.Map(dbGames, func(game model.DotaMatchWithRelation, _ int) Game {
		g := Game{
			ActivateTime:              game.StartedAt,
			LobbyType:                 game.LobbyType,
			GameMode:                  game.GameModeID,
			AvarageMmr:                game.AvarageMmr,
			WeekedTourneyBracketRound: &game.WeekendTourneyBracketRound.String,
			WeekedTourneySkillLevel:   &game.WeekendTourneySkillLevel.String,
			MatchId:                   &game.MatchID,
			LobbyId:                   game.LobbyID,
		}

		g.Players = lo.Map(game.Players, func(p int64, i int) Player {
			return Player{
				AccountId: p,
				HeroId:    game.PlayersHeroes[i],
			}
		})

		return g
	})

	return &mappedGames
}

type Player struct {
	AccountId int64 `json:"account_id"`
	HeroId    int64 `json:"hero_id"`
}

type Game struct {
	ActivateTime              time.Time `json:"activate_time"`
	LobbyType                 int32     `json:"lobby_type"`
	GameMode                  int32     `json:"game_mode"`
	AvarageMmr                int32     `json:"avarage_mmr"`
	Players                   []Player  `json:"players"`
	WeekedTourneyBracketRound *string   `json:"weeked_tourney_bracket_round"`
	WeekedTourneySkillLevel   *string   `json:"weeked_tourney_skill_level"`
	MatchId                   *string   `json:"match_id"`
	LobbyId                   string    `json:"lobby_id"`
}

func GetAccountsByChannelId(db *gorm.DB, channelId string) *[]string {
	accounts := []model.ChannelsDotaAccounts{}
	err := db.Where(`"channelId" = ?`, channelId).Select(&accounts).Error

	if err != nil {
		return nil
	}

	mappedAccounts := lo.Map(accounts, func(a model.ChannelsDotaAccounts, _ int) string {
		return a.ID
	})

	return &mappedAccounts
}
