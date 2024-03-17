package pubg

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/NovikovRoman/pubg"
	"github.com/eko/gocache/lib/v4/cache"
	"github.com/eko/gocache/lib/v4/marshaler"
	"github.com/eko/gocache/lib/v4/store"
	"github.com/redis/go-redis/v9"
	"github.com/samber/lo"
)

type Client struct {
	senderMu   sync.Mutex
	sender     *pubg.Client
	apiKeys    []string
	currentKey int
	cache      cache.CacheInterface[any]
	marshal    *marshaler.Marshaler
}

func NewClient(cacheStore store.StoreInterface, apiKeys ...string) (*Client, error) {
	if len(apiKeys) == 0 {
		return nil, errors.New("You must specify at least one API key. ")
	}

	c := cache.New[any](cacheStore)
	return &Client{
		sender:  pubg.NewClient(apiKeys[0], nil),
		apiKeys: apiKeys,
		cache:   c,
		marshal: marshaler.New(c),
	}, nil
}

func (c *Client) GetPlayerByNickname(ctx context.Context, nickname string) (*pubg.Players, error) {
	c.senderMu.Lock()
	defer c.senderMu.Unlock()

	players := &pubg.Players{}
	_, err := c.marshal.Get(ctx, fmt.Sprintf("player:%s", nickname), players)
	switch err {
	case nil:
		return players, nil
	case redis.Nil:
	default:
		return nil, err
	}

	players, err = c.sender.PlayersByNames(pubg.SteamPlatform, nickname)
	if err != nil {
		if _, ok := err.(pubg.ErrTooManyRequest); ok {
			c.currentKey = (c.currentKey + 1) % len(c.apiKeys)
			c.sender = pubg.NewClient(c.apiKeys[c.currentKey], nil)
			return c.GetPlayerByNickname(ctx, nickname)
		}

		return nil, err
	}

	err = c.marshal.Set(ctx, fmt.Sprintf("player:%s", nickname), players)
	if err != nil {
		return nil, err
	}

	return players, nil
}

func (c *Client) GetCurrentSeason(ctx context.Context) (*string, error) {
	c.senderMu.Lock()
	defer c.senderMu.Unlock()

	var seasonId string
	_, err := c.marshal.Get(ctx, "seasonId", &seasonId)
	switch err {
	case nil:
		return &seasonId, nil
	case redis.Nil:
	default:
		// return nil, err
	}

	seasons, err := c.sender.Seasons(pubg.SteamPlatform)
	if err != nil {
		if _, ok := err.(pubg.ErrTooManyRequest); ok {
			c.currentKey = (c.currentKey + 1) % len(c.apiKeys)
			c.sender = pubg.NewClient(c.apiKeys[c.currentKey], nil)
			return c.GetCurrentSeason(ctx)
		}

		return nil, err
	}

	for _, season := range seasons.Data {
		if season.Attributes.IsCurrentSeason {
			err = c.marshal.Set(
				ctx,
				"seasonId",
				string(season.ID),
				store.WithExpiration(30*24*time.Hour),
			)
			if err != nil {
				return nil, err
			}

			return lo.ToPtr(string(season.ID)), nil
		}
	}

	return nil, errors.New("error in pubg api")
}

func (c *Client) GetLifetimeStats(
	ctx context.Context,
	accountID string,
) (*pubg.LifetimeStatsPlayer, error) {
	lifetimeStats := &pubg.LifetimeStatsPlayer{}
	_, err := c.marshal.Get(ctx, fmt.Sprintf("lifetimeStats:%s", accountID), lifetimeStats)
	switch err {
	case nil:
		return lifetimeStats, nil
	case redis.Nil:
	default:
		// return nil, err
	}

	lifetimeStats, err = c.sender.LifetimeStatsPlayer(pubg.SteamPlatform, accountID)
	if err != nil {
		if _, ok := err.(pubg.ErrTooManyRequest); ok {
			c.currentKey = (c.currentKey + 1) % len(c.apiKeys)
			c.sender = pubg.NewClient(c.apiKeys[c.currentKey], nil)
			return c.GetLifetimeStats(ctx, accountID)
		}
	}

	err = c.marshal.Set(
		ctx,
		fmt.Sprintf("lifetimeStats:%s", accountID),
		lifetimeStats,
		store.WithExpiration(3*time.Hour),
	)

	return lifetimeStats, err
}
