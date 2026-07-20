package stats

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"strconv"
	"time"

	"github.com/twirapp/kv"
	kvoptions "github.com/twirapp/kv/options"
	"github.com/twirapp/twir/libs/integrations/opendota"
	"github.com/twirapp/twir/libs/integrations/stratz"
)

const (
	winProbabilityCachePrefix = "cache:twir:dota:wp:"
	notablePlayersCachePrefix = "cache:twir:dota:np:"
	heroesCacheKey            = "cache:twir:dota:heroes"

	winProbabilityCacheTTL = 30 * time.Second
	notablePlayersCacheTTL = 5 * time.Minute
	heroesCacheTTL         = 24 * time.Hour
)

type LastGame struct {
	HeroName  string
	Kills     int
	Deaths    int
	Assists   int
	Win       bool
	DurationS int
}

type Stats struct {
	stratz   *stratz.Client
	opendota *opendota.Client
	kv       kv.KV
	logger   *slog.Logger
}

func New(
	stratzClient *stratz.Client,
	opendotaClient *opendota.Client,
	kvStore kv.KV,
	logger *slog.Logger,
) *Stats {
	return &Stats{
		stratz:   stratzClient,
		opendota: opendotaClient,
		kv:       kvStore,
		logger:   logger,
	}
}

func getCached[T any](ctx context.Context, kvStore kv.KV, key string) (T, bool, error) {
	var value T

	cacheBytes, err := kvStore.Get(ctx, key).Bytes()
	if err != nil {
		if errors.Is(err, kv.ErrKeyNil) {
			return value, false, nil
		}
		return value, false, fmt.Errorf("failed to get %q from cache: %w", key, err)
	}

	if len(cacheBytes) == 0 {
		return value, false, nil
	}

	if err := json.Unmarshal(cacheBytes, &value); err != nil {
		return value, false, fmt.Errorf("failed to unmarshal cached %q: %w", key, err)
	}

	return value, true, nil
}

func setCached(ctx context.Context, kvStore kv.KV, key string, ttl time.Duration, value any) error {
	cacheBytes, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal %q for cache: %w", key, err)
	}

	if err := kvStore.Set(ctx, key, cacheBytes, kvoptions.WithExpire(ttl)); err != nil {
		return fmt.Errorf("failed to set %q to cache: %w", key, err)
	}

	return nil
}

func (s *Stats) WinProbability(ctx context.Context, matchID int64) (float64, error) {
	key := winProbabilityCachePrefix + strconv.FormatInt(matchID, 10)

	cached, hit, err := getCached[float64](ctx, s.kv, key)
	if err != nil {
		s.logger.Warn("failed to read win probability cache", slog.Any("err", err))
	} else if hit {
		return cached, nil
	}

	if s.stratz == nil || !s.stratz.Enabled() {
		return 0, nil
	}

	probability, err := s.stratz.WinProbability(ctx, matchID)
	if err != nil {
		s.logger.Error(
			"failed to fetch win probability",
			slog.Any("err", err),
			slog.Int64("match_id", matchID),
		)
		return 0, err
	}

	if err := setCached(ctx, s.kv, key, winProbabilityCacheTTL, probability); err != nil {
		s.logger.Warn("failed to cache win probability", slog.Any("err", err))
	}

	return probability, nil
}

func (s *Stats) NotablePlayers(
	ctx context.Context,
	matchID int64,
	streamerAccountID string,
) ([]string, error) {
	key := notablePlayersCachePrefix + strconv.FormatInt(matchID, 10)

	cached, hit, err := getCached[[]string](ctx, s.kv, key)
	if err != nil {
		s.logger.Warn("failed to read notable players cache", slog.Any("err", err))
	} else if hit {
		return cached, nil
	}

	if s.stratz == nil || !s.stratz.Enabled() {
		return nil, nil
	}

	notables, err := s.stratz.NotablePlayers(ctx, matchID)
	if err != nil {
		s.logger.Error(
			"failed to fetch notable players",
			slog.Any("err", err),
			slog.Int64("match_id", matchID),
		)
		return nil, err
	}

	proPlayers, err := s.opendota.ProPlayers(ctx)
	if err != nil {
		s.logger.Error("failed to fetch pro players", slog.Any("err", err))
		return nil, err
	}

	prosByID := make(map[int64]opendota.ProPlayer, len(proPlayers))
	for _, pro := range proPlayers {
		prosByID[pro.AccountID] = pro
	}

	names := make([]string, 0, len(notables))
	for _, notable := range notables {
		if streamerAccountID != "" &&
			strconv.FormatInt(notable.AccountID, 10) == streamerAccountID {
			continue
		}

		name := notable.Name
		if pro, ok := prosByID[notable.AccountID]; ok && pro.Name != "" {
			name = pro.Name
		}
		if name == "" {
			continue
		}

		names = append(names, name)
	}

	if err := setCached(ctx, s.kv, key, notablePlayersCacheTTL, names); err != nil {
		s.logger.Warn("failed to cache notable players", slog.Any("err", err))
	}

	return names, nil
}

func (s *Stats) heroes(ctx context.Context) (map[int]string, error) {
	cached, hit, err := getCached[map[int]string](ctx, s.kv, heroesCacheKey)
	if err != nil {
		s.logger.Warn("failed to read heroes cache", slog.Any("err", err))
	} else if hit {
		return cached, nil
	}

	heroes, err := s.opendota.Heroes(ctx)
	if err != nil {
		return nil, err
	}

	if err := setCached(ctx, s.kv, heroesCacheKey, heroesCacheTTL, heroes); err != nil {
		s.logger.Warn("failed to cache heroes", slog.Any("err", err))
	}

	return heroes, nil
}

func (s *Stats) LastGame(ctx context.Context, accountID int64) (*LastGame, error) {
	matches, err := s.opendota.RecentMatches(ctx, accountID)
	if err != nil {
		s.logger.Error(
			"failed to fetch recent matches",
			slog.Any("err", err),
			slog.Int64("account_id", accountID),
		)
		return nil, err
	}

	if len(matches) == 0 {
		return nil, nil
	}

	heroes, err := s.heroes(ctx)
	if err != nil {
		s.logger.Warn(
			"failed to fetch heroes, hero name will be unknown",
			slog.Any("err", err),
		)
		heroes = map[int]string{}
	}

	last := matches[0]
	heroName, ok := heroes[last.HeroID]
	if !ok {
		heroName = "Unknown Hero"
	}

	return &LastGame{
		HeroName:  heroName,
		Kills:     last.Kills,
		Deaths:    last.Deaths,
		Assists:   last.Assists,
		Win:       last.Won(),
		DurationS: last.Duration,
	}, nil
}
