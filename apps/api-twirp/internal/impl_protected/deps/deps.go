package deps

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Deps struct {
	Redis *redis.Client
	Db    *gorm.DB
}
