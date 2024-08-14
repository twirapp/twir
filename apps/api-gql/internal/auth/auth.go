package auth

import (
	"context"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/alexedwards/scs/goredisstore"
	"github.com/alexedwards/scs/v2"
	"github.com/gin-gonic/gin"
	"github.com/nicklaw5/helix/v2"
	"github.com/redis/go-redis/v9"
	model "github.com/satont/twir/libs/gomodels"
	"go.uber.org/fx"
	"gorm.io/gorm"
)

type Opts struct {
	fx.In

	Redis *redis.Client
	Gorm  *gorm.DB
}

type Auth struct {
	sessionManager *scs.SessionManager
	gorm           *gorm.DB
}

func NewSessions(opts Opts) *Auth {
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour * 31
	sessionManager.Store = goredisstore.New(opts.Redis)

	gob.Register(model.Users{})
	gob.Register(helix.User{})

	return &Auth{
		sessionManager: sessionManager,
		gorm:           opts.Gorm,
	}
}

const SESSION_KEY = "__session__"

func (s *Auth) Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println(r)
				c.String(500, "Internal Server Error")
			}
		}()

		cookie, err := c.Cookie(s.sessionManager.Cookie.Name)
		if err != nil {
			cookie = ""
		}

		session, err := s.sessionManager.Load(c.Request.Context(), cookie)
		if err != nil {
			s.sessionManager.ErrorFunc(c.Writer, c.Request, err)
			return
		}

		c.Set(SESSION_KEY, session)

		// ctx := context.WithValue(c.Request.Context(), SESSION_KEY, session)
		c.Request = c.Request.WithContext(session)

		sessionToken, expiryTime, err := s.sessionManager.Commit(session)
		if err != nil {
			panic(err)
		}

		s.sessionManager.WriteSessionCookie(session, c.Writer, sessionToken, expiryTime)

		c.Next()
	}
}

func (s *Auth) Put(ctx context.Context, key string, val interface{}) {
	s.sessionManager.Put(ctx, key, val)
	s.sessionManager.Commit(ctx)
}
