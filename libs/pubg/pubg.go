package pubg

import (
	"context"
	"errors"
	"fmt"
	"strings"
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
	maxRetries int
}

var ErrOverloaded = errors.New("overloaded, try again later")
var ErrPubg = errors.New("pubg error")

func NewClient(
	cacheStore store.StoreInterface,
	maxRetries int,
	apiKeys ...string,
) (*Client, error) {
	if len(apiKeys) == 0 {
		return nil, errors.New("You must specify at least one API key. ")
	}

	c := cache.New[any](cacheStore)
	return &Client{
		sender:     pubg.NewClient(apiKeys[0], nil),
		apiKeys:    apiKeys,
		cache:      c,
		marshal:    marshaler.New(c),
		maxRetries: maxRetries,
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

const maxRetries = 3

func (c *Client) GetLifetimeStats(
	ctx context.Context,
	accountID string,
) (*pubg.LifetimeStatsPlayer, error) {
	c.senderMu.Lock()
	defer c.senderMu.Unlock()

	for retry := 0; retry < c.maxRetries; retry++ {
		lifetimeStats := &pubg.LifetimeStatsPlayer{}
		_, err := c.marshal.Get(ctx, fmt.Sprintf("lifetimeStats:%s", accountID), lifetimeStats)
		switch err {
		case nil:
			return lifetimeStats, nil
		case redis.Nil:
		default:
		}

		lifetimeStats, err = c.sender.LifetimeStatsPlayer(pubg.SteamPlatform, accountID)
		if err != nil {
			if strings.Contains(err.Error(), "429") {
				c.currentKey = (c.currentKey + 1) % len(c.apiKeys)
				c.sender = pubg.NewClient(c.apiKeys[c.currentKey], nil)
				if retry == maxRetries-1 {
					return nil, errors.Join(err, ErrOverloaded)
				}
				continue
			}
			return nil, errors.Join(err, ErrPubg)
		}

		err = c.marshal.Set(
			ctx,
			fmt.Sprintf("lifetimeStats:%s", accountID),
			lifetimeStats,
			store.WithExpiration(3*time.Hour),
		)
		return lifetimeStats, err
	}

	return nil, ErrOverloaded
}
