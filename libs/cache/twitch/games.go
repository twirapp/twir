package twitch

import (
	"context"
	"fmt"
	"time"

	"github.com/goccy/go-json"
	"github.com/nicklaw5/helix/v2"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
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

	span := trace.SpanFromContext(ctx)
	defer span.End()

	span.SetAttributes(
		attribute.String("gameID", gameID),
	)

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

// GetGames fetches games and returns them in the same order as the input IDs.
func (c *CachedTwitchClient) GetGames(
	ctx context.Context,
	gameIDs []string,
) (
	[]*helix.Game,
	error,
) {
	if len(gameIDs) == 0 {
		return nil, nil
	}

	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(
		attribute.String("gameIDs", fmt.Sprintf("%v", gameIDs)),
	)

	fetchedGames := map[string]helix.Game{}
	gamesForRequest := make([]string, 0, len(gameIDs))

	for _, gameID := range gameIDs {
		if bytes, _ := c.redis.Get(
			ctx,
			buildGamesCacheKeyForId(gameID),
		).Bytes(); len(bytes) > 0 {
			var game helix.Game
			if err := json.Unmarshal(bytes, &game); err != nil {
				return nil, err
			}

			fetchedGames[gameID] = game
			continue
		} else {
			gamesForRequest = append(gamesForRequest, gameID)
		}
	}

	if len(gamesForRequest) > 0 {
		twitchGetGamesReq, err := c.client.GetGames(
			&helix.GamesParams{
				IDs: gamesForRequest,
			},
		)
		if err != nil {
			return nil, err
		}
		if twitchGetGamesReq.ErrorMessage != "" {
			return nil, fmt.Errorf(twitchGetGamesReq.ErrorMessage)
		}

		for _, game := range twitchGetGamesReq.Data.Games {
			fetchedGames[game.ID] = game

			gameBytes, err := json.Marshal(game)
			if err != nil {
				return nil, err
			}

			if err := c.redis.Set(
				ctx,
				buildGamesCacheKeyForId(game.ID),
				gameBytes,
				gamesGetTTL,
			).Err(); err != nil {
				return nil, err
			}
		}
	}

	games := make([]*helix.Game, len(gameIDs))
	for i, gameID := range gameIDs {
		if _, ok := fetchedGames[gameID]; !ok {
			continue
		}
		game := fetchedGames[gameID]
		games[i] = &game
	}

	return games, nil
}
