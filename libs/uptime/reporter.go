package uptime

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/twirapp/twir/libs/logger"
	rediskeys "github.com/twirapp/twir/libs/redis_keys"
)

const (
	defaultInterval = 15 * time.Second
	defaultTTL      = 5 * time.Minute
)

type ReporterOpts struct {
	ServiceName string
	Instance    string
	Interval    time.Duration
	TTL         time.Duration
	Meta        map[string]string
}

type Reporter struct {
	client    *redis.Client
	service   string
	instance  string
	interval  time.Duration
	ttl       time.Duration
	startedAt time.Time

	metaMu sync.RWMutex
	meta   map[string]string

	startOnce sync.Once
}

func NewReporter(client *redis.Client, opts ReporterOpts) (*Reporter, error) {
	if opts.ServiceName == "" {
		return nil, fmt.Errorf("service name is required")
	}
	if opts.Interval < 0 {
		return nil, fmt.Errorf("interval must not be negative")
	}
	if opts.TTL < 0 {
		return nil, fmt.Errorf("TTL must not be negative")
	}

	instance, err := resolveInstance(opts.Instance)
	if err != nil {
		return nil, fmt.Errorf("resolve instance: %w", err)
	}

	interval := opts.Interval
	if interval == 0 {
		interval = defaultInterval
	}

	ttl := opts.TTL
	if ttl == 0 {
		ttl = defaultTTL
	}

	meta := make(map[string]string, len(opts.Meta))
	for key, value := range opts.Meta {
		meta[key] = value
	}

	return &Reporter{
		client:    client,
		service:   opts.ServiceName,
		instance:  instance,
		interval:  interval,
		ttl:       ttl,
		startedAt: time.Now().UTC(),
		meta:      meta,
	}, nil
}

func (r *Reporter) Instance() string {
	return r.instance
}

func (r *Reporter) SetMeta(key string, value string) {
	r.metaMu.Lock()
	defer r.metaMu.Unlock()

	r.meta[key] = value
}

func (r *Reporter) DeleteMeta(key string) {
	r.metaMu.Lock()
	defer r.metaMu.Unlock()

	delete(r.meta, key)
}

func (r *Reporter) Start(ctx context.Context) {
	r.startOnce.Do(func() {
		go func() {
			r.report(ctx)

			ticker := time.NewTicker(r.interval)
			defer ticker.Stop()

			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					r.report(ctx)
				}
			}
		}()
	})
}

func (r *Reporter) report(ctx context.Context) {
	if err := r.beat(ctx); err != nil {
		slog.ErrorContext(ctx, "uptime reporter beat failed", logger.Error(err))
	}
}

func (r *Reporter) beat(ctx context.Context) error {
	meta, err := r.metaJSON()
	if err != nil {
		return fmt.Errorf("marshal metadata: %w", err)
	}

	key := rediskeys.UptimeInstanceKey(r.service, r.instance)
	now := time.Now().UTC().Unix()
	pipeline := r.client.TxPipeline()
	pipeline.HSet(
		ctx,
		key,
		"service", r.service,
		"instance", r.instance,
		"startedAt", r.startedAt.Unix(),
		"updatedAt", now,
	)
	if meta == "" {
		pipeline.HDel(ctx, key, "meta")
	} else {
		pipeline.HSet(ctx, key, "meta", meta)
	}
	pipeline.PExpire(ctx, key, r.ttl)

	if _, err := pipeline.Exec(ctx); err != nil {
		return fmt.Errorf("write uptime heartbeat: %w", err)
	}

	return nil
}

func (r *Reporter) metaJSON() (string, error) {
	r.metaMu.RLock()
	defer r.metaMu.RUnlock()

	if len(r.meta) == 0 {
		return "", nil
	}

	meta, err := json.Marshal(r.meta)
	if err != nil {
		return "", fmt.Errorf("encode metadata: %w", err)
	}

	return string(meta), nil
}

func resolveInstance(instance string) (string, error) {
	if instance != "" {
		return instance, nil
	}

	if replica := os.Getenv("REPLICA"); replica != "" {
		return "slot-" + replica, nil
	}

	if hostname, err := os.Hostname(); err == nil && hostname != "" {
		if len(hostname) > 12 {
			hostname = hostname[:12]
		}
		return hostname, nil
	}

	suffix := make([]byte, 4)
	if _, err := rand.Read(suffix); err != nil {
		return "", fmt.Errorf("generate random suffix: %w", err)
	}

	return "rand-" + hex.EncodeToString(suffix), nil
}
