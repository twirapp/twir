package uptime

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
	rediskeys "github.com/twirapp/twir/libs/redis_keys"
)

type InstanceStatus struct {
	Service   string
	Instance  string
	StartedAt time.Time
	UpdatedAt time.Time
	Meta      map[string]string
}

func (s InstanceStatus) Uptime(now time.Time) time.Duration {
	return now.Sub(s.StartedAt)
}

func (s InstanceStatus) IsStale(now time.Time, threshold time.Duration) bool {
	return now.Sub(s.UpdatedAt) > threshold
}

func ReadAll(ctx context.Context, client *redis.Client) ([]InstanceStatus, error) {
	keys, err := readKeys(ctx, client)
	if err != nil {
		return nil, fmt.Errorf("scan uptime keys: %w", err)
	}
	if len(keys) == 0 {
		return []InstanceStatus{}, nil
	}

	pipeline := client.Pipeline()
	commands := make([]*redis.MapStringStringCmd, 0, len(keys))
	for _, key := range keys {
		commands = append(commands, pipeline.HGetAll(ctx, key))
	}
	if _, err := pipeline.Exec(ctx); err != nil {
		successfulRead := false
		for _, command := range commands {
			if command.Err() == nil {
				successfulRead = true
				break
			}
		}
		if !successfulRead {
			return nil, fmt.Errorf("read uptime hashes: %w", err)
		}
	}

	statuses := make([]InstanceStatus, 0, len(commands))
	for _, command := range commands {
		if command.Err() != nil {
			continue
		}
		status, ok := parseStatus(command.Val())
		if !ok {
			continue
		}
		statuses = append(statuses, status)
	}

	sort.Slice(statuses, func(first int, second int) bool {
		if statuses[first].Service == statuses[second].Service {
			return statuses[first].Instance < statuses[second].Instance
		}
		return statuses[first].Service < statuses[second].Service
	})

	return statuses, nil
}

func readKeys(ctx context.Context, client *redis.Client) ([]string, error) {
	keys := []string{}
	var cursor uint64
	for {
		batch, nextCursor, err := client.Scan(ctx, cursor, rediskeys.UptimeScanPattern, 0).Result()
		if err != nil {
			return nil, fmt.Errorf("scan cursor %d: %w", cursor, err)
		}
		keys = append(keys, batch...)
		if nextCursor == 0 {
			return keys, nil
		}
		cursor = nextCursor
	}
}

func parseStatus(values map[string]string) (InstanceStatus, bool) {
	service := values["service"]
	instance := values["instance"]
	if service == "" || instance == "" {
		return InstanceStatus{}, false
	}

	startedAt, err := parseUnixTime(values["startedAt"])
	if err != nil {
		return InstanceStatus{}, false
	}
	updatedAt, err := parseUnixTime(values["updatedAt"])
	if err != nil {
		return InstanceStatus{}, false
	}

	meta := map[string]string{}
	if value, ok := values["meta"]; ok {
		if err := json.Unmarshal([]byte(value), &meta); err != nil {
			return InstanceStatus{}, false
		}
	}

	return InstanceStatus{
		Service:   service,
		Instance:  instance,
		StartedAt: startedAt,
		UpdatedAt: updatedAt,
		Meta:      meta,
	}, true
}

func parseUnixTime(value string) (time.Time, error) {
	seconds, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		return time.Time{}, fmt.Errorf("parse unix seconds: %w", err)
	}

	return time.Unix(seconds, 0).UTC(), nil
}
