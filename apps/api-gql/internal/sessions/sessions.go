package sessions

import (
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
)

type Opts struct {
	fx.In

	Redis *redis.Client
}

type Sessions struct {
	sessionManager *scs.SessionManager
}

func New(opts Opts) *Sessions {
	sessionManager := scs.New()
	sessionManager.Lifetime = 24 * time.Hour * 31
	sessionManager.Store = goredisstore.New(opts.Redis)

	gob.Register(model.Users{})
	gob.Register(helix.User{})

	return &Sessions{
		sessionManager: sessionManager,
	}
}

const SESSION_KEY = "__session__"

func (s *Sessions) Middleware() gin.HandlerFunc {
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

		/*
			headerKey := "X-Session"
			headerKeyExpiry := "X-Session-Expiry"

			ctx, err := mgr.Load(r.Context(), r.Header.Get(headerKey))
			if err != nil {
				log.Output(2, err.Error())
				http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				return
			}

			// // replace "inner" ctx with session (wrap around)
			bw := &bufferedResponseWriter{ResponseWriter: w}
			sr := r.WithContext(ctx)
			next.ServeHTTP(bw, sr)

			if s.Status(ctx) == scs.Modified {
				token, expiry, err := s.Commit(ctx)
				if err != nil {
					log.Output(2, err.Error())
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
					return
				}

				w.Header().Set(headerKey, token)
				w.Header().Set(headerKeyExpiry, expiry.Format(http.TimeFormat))
			}

			if bw.code != 0 {
				w.WriteHeader(bw.code)
			}
			w.Write(bw.buf.Bytes())
		*/
	}
}