package twitch

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strings"
	"time"

	"github.com/goccy/go-json"
	"github.com/redis/go-redis/v9"
	twlib "github.com/satont/twir/libs/twitch"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

const searchCategoryKey = "cache:twir:twitch:categories:search:"

func (c *CachedTwitchClient) SearchCategory(
	ctx context.Context,
	query string,
) (*twlib.FoundCategory, error) {
	if query == "" {
		return nil, nil
	}
	span := trace.SpanFromContext(ctx)
	defer span.End()
	span.SetAttributes(
		attribute.String("query", query),
	)

	query = strings.ToLower(query)

	querySum := sha256.Sum256([]byte(query))
	queryHash := hex.EncodeToString(querySum[:])

	cachedBytes, err := c.redis.Get(
		ctx,
		searchCategoryKey+queryHash,
	).Bytes()
	if err != nil && !errors.Is(err, redis.Nil) {
		return nil, err
	}
	if len(cachedBytes) > 0 {
		var cachedCategory twlib.FoundCategory
		if err := json.Unmarshal(cachedBytes, &cachedCategory); err != nil {
			return nil, err
		}

		return &cachedCategory, nil
	}

	category, err := twlib.SearchCategory(ctx, query)
	if err != nil {
		return nil, err
	}

	categoryBytes, err := json.Marshal(category)
	if err != nil {
		return nil, err
	}

	if err := c.redis.Set(
		ctx,
		searchCategoryKey+queryHash,
		categoryBytes,
		31*24*time.Hour,
	).Err(); err != nil {
		return nil, err
	}

	return category, nil
}
