package fiber

import (
	"fmt"
	"github.com/gofiber/storage/redis"
	config "github.com/satont/tsuwari/libs/config"
	"runtime"
	"sync"
)

type RedisCacheStorage struct {
	*redis.Storage
}

func NewCache(cfg *config.Config) *RedisCacheStorage {
	store := redis.New(redis.Config{
		URL:       cfg.RedisUrl,
		Database:  0,
		Reset:     false,
		TLSConfig: nil,
		PoolSize:  10 * (runtime.NumCPU() + 1),
	})

	storage := &RedisCacheStorage{
		Storage: store,
	}

	return storage
}

func (c *RedisCacheStorage) BuildKey(path string) string {
	return fmt.Sprintf("fiber:cache:%s", path)
}

func (c *RedisCacheStorage) DeleteGet(path string) {
	baseKey := fmt.Sprintf("%s_%s", c.BuildKey(path), "GET")
	fmt.Println("baseKey", baseKey)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err := c.Delete(baseKey)
		if err != nil {
			fmt.Println(err)
		}
	}()

	go func() {
		defer wg.Done()
		err := c.Delete(fmt.Sprintf("%s_body", baseKey))
		if err != nil {
			fmt.Println(err)
		}
	}()

	wg.Wait()
}
