package oauth

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisFixture struct {
	admin     *redis.Client
	source    *redis.Client
	container string
}

func newRedisFixture(t *testing.T) redisFixture {
	t.Helper()
	address := os.Getenv("TWIR_OAUTH_TEST_REDIS_ADDR")
	container := os.Getenv("TWIR_OAUTH_TEST_REDIS_CONTAINER")
	if address == "" || container == "" {
		t.Fatal("TWIR_OAUTH_TEST_REDIS_ADDR and TWIR_OAUTH_TEST_REDIS_CONTAINER must be configured by the real-Redis test task")
	}
	admin := redis.NewClient(&redis.Options{
		Addr: address, ReadTimeout: time.Second, WriteTimeout: time.Second,
		ContextTimeoutEnabled: true,
	})
	source := redis.NewClient(&redis.Options{
		Addr: address, ReadTimeout: -1, WriteTimeout: -1,
	})
	t.Cleanup(func() {
		_ = source.Close()
		_ = admin.Close()
	})
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	if err := admin.FlushDB(ctx).Err(); err != nil {
		t.Fatalf("flush Redis: %v", err)
	}
	return redisFixture{admin: admin, source: source, container: container}
}

func newIsolatedRedisFixture(t *testing.T) redisFixture {
	t.Helper()
	container := fmt.Sprintf("twir-oauth-kill-%d", time.Now().UnixNano())
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	output, err := exec.CommandContext(
		ctx,
		"docker", "run", "--detach", "--name", container,
		"--publish", "127.0.0.1::6379", "redis:7.4-alpine",
	).CombinedOutput()
	cancel()
	if err != nil {
		t.Fatalf("start isolated Redis: %v: %s", err, output)
	}
	t.Cleanup(func() {
		cleanupContext, cleanupCancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cleanupCancel()
		_ = exec.CommandContext(cleanupContext, "docker", "rm", "--force", container).Run()
	})
	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	output, err = exec.CommandContext(ctx, "docker", "port", container, "6379/tcp").CombinedOutput()
	cancel()
	if err != nil {
		t.Fatalf("resolve isolated Redis port: %v: %s", err, output)
	}
	mapping := strings.TrimSpace(string(output))
	separator := strings.LastIndexByte(mapping, ':')
	if separator < 0 {
		t.Fatalf("invalid isolated Redis port mapping %q", mapping)
	}
	fixture := redisFixture{
		admin: redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:" + mapping[separator+1:], ReadTimeout: time.Second,
			WriteTimeout: time.Second, ContextTimeoutEnabled: true,
		}),
		source: redis.NewClient(&redis.Options{
			Addr: "127.0.0.1:" + mapping[separator+1:], ReadTimeout: -1, WriteTimeout: -1,
		}),
		container: container,
	}
	t.Cleanup(func() {
		_ = fixture.source.Close()
		_ = fixture.admin.Close()
	})
	fixture.waitReady(t)
	return fixture
}

func (fixture redisFixture) locker(t *testing.T) *RedisLocker {
	return fixture.lockerFor(t, fixture.source)
}

func (fixture redisFixture) lockerFor(t *testing.T, client *redis.Client) *RedisLocker {
	t.Helper()
	locker, err := NewRedisLocker(client, RedisLockerOptions{
		Prefix: "oauth-acceptance", TTL: 600 * time.Millisecond,
		RenewEvery: 150 * time.Millisecond, Timeout: 100 * time.Millisecond,
	})
	if err != nil {
		t.Fatal(err)
	}
	return locker
}

func (fixture redisFixture) keys(t *testing.T) []string {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	keys, err := fixture.admin.Keys(ctx, "oauth-acceptance:*").Result()
	if err != nil {
		t.Fatalf("list Redis keys: %v", err)
	}
	return keys
}

func (fixture redisFixture) docker(t *testing.T, action string) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if output, err := exec.CommandContext(ctx, "docker", action, fixture.container).CombinedOutput(); err != nil {
		t.Fatalf("docker %s: %v: %s", action, err, output)
	}
}

func (fixture redisFixture) waitReady(t *testing.T) {
	t.Helper()
	deadline := time.NewTimer(5 * time.Second)
	defer deadline.Stop()
	ticker := time.NewTicker(25 * time.Millisecond)
	defer ticker.Stop()
	for {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second)
		err := exec.CommandContext(ctx, "docker", "exec", fixture.container, "redis-cli", "ping").Run()
		cancel()
		if err == nil {
			return
		}
		select {
		case <-deadline.C:
			t.Fatal("Redis did not become ready")
		case <-ticker.C:
		}
	}
}
