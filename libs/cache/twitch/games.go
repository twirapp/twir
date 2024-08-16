package twitch

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/nicklaw5/helix/v2"
)

const gamesGetKey = "cache:twir:twitch:games:by-id:"

// 2 weeks
const gamesGetTTL = 14 * 24 * time.Hour

func buildGamesCacheKeyForId(gameID string) string {
	return gamesGetKey + gameID
}

func (c *CachedTwitchClient) GetGame(
	ctx context.Context,
	gameID string,
) (
	*helix.Game,
	error,
) {
	if gameID == "" {
		return nil, nil
	}

	if bytes, _ := c.redis.Get(
		ctx,
		buildGamesCacheKeyForId(gameID),
	).Bytes(); len(bytes) > 0 {
		var game helix.Game
		if err := json.Unmarshal(bytes, &game); err != nil {
			return nil, err
		}

		return &game, nil
	}

	twitchGetGamesReq, err := c.client.GetGames(
		&helix.GamesParams{
			IDs: []string{gameID},
		},
	)
	if err != nil {
		return nil, err
	}
	if twitchGetGamesReq.ErrorMessage != "" {
		return nil, fmt.Errorf(twitchGetGamesReq.ErrorMessage)
	}

	if len(twitchGetGamesReq.Data.Games) == 0 {
		return nil, nil
	}

	game := twitchGetGamesReq.Data.Games[0]

	gameBytes, err := json.Marshal(game)
	if err != nil {
		return nil, err
	}

	if err := c.redis.Set(
		ctx,
		buildGamesCacheKeyForId(gameID),
		gameBytes,
		gamesGetTTL,
	).Err(); err != nil {
		return nil, err
	}

	return &game, nil
}
