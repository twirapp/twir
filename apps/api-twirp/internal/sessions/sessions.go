package sessions

import (
	"encoding/gob"
	"github.com/alexedwards/scs/goredisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/nicklaw5/helix/v2"
	"github.com/redis/go-redis/v9"
	model "github.com/satont/twir/libs/gomodels"
	"time"
)

func New(redisClient *redis.Client) *scs.SessionManager {
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour * 31
	sessionManager.Store = goredisstore.New(redisClient)

	gob.Register(model.Users{})
	gob.Register(helix.User{})

	return sessionManager
}
