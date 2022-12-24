package redis_storage

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/gofiber/storage/redis"
)

type RedisStorage struct {
	*redis.Storage
}

func NewCache(url string) *RedisStorage {
	store := redis.New(redis.Config{
		URL:       url,
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * runtime.GOMAXPROCS(0),
	})

	storage := &RedisStorage{
		Storage: store,
	}

	return storage
}

func (c *RedisStorage) DeleteByMethod(key string, method string) {
	baseKey := fmt.Sprintf("%s_%s", key, strings.ToUpper(method))
	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		c.Delete(baseKey)
	}()

	go func() {
		defer wg.Done()
		c.Delete(fmt.Sprintf("%s_body", baseKey))
	}()

	wg.Wait()
}
