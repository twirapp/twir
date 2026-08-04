package app

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"time"

	"github.com/twirapp/twir/apps/twitch-mock/internal/admin"
	"github.com/twirapp/twir/apps/twitch-mock/internal/config"
	"github.com/twirapp/twir/apps/twitch-mock/internal/handlers"
	"github.com/twirapp/twir/apps/twitch-mock/internal/state"
	"github.com/twirapp/twir/apps/twitch-mock/internal/websocket"
	"github.com/goforj/wire"
	"github.com/twirapp/twir/libs/baseapp/lifecycle"
)

var ProviderSet = wire.NewSet(
	lifecycle.New,
	NewLogger,
	config.New,
	state.New,
	handlers.New,
	websocket.New,
	admin.New,
	NewApplication,
)

type Application struct {
	lifecycle *lifecycle.Lifecycle
}

func (a *Application) Run() error {
	return a.lifecycle.Run()
}

func NewLogger() *slog.Logger {
	return slog.Default()
}

func NewApplication(
	lc *lifecycle.Lifecycle,
	logger *slog.Logger,
	cfg *config.Config,
	httpServer *handlers.Server,
	wsServer *websocket.Server,
	adminServer *admin.Server,
) *Application {
	apiSrv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           httpServer.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	wsSrv := &http.Server{
		Addr:              cfg.WSAddr,
		Handler:           wsServer.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	adminSrv := &http.Server{
		Addr:              cfg.AdminAddr,
		Handler:           adminServer.Handler(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	lc.Append(lifecycle.Hook{
		OnStart: func(context.Context) error {
			start := func(name string, srv *http.Server) {
				go func() {
					logger.Info("server starting", slog.String("server", name), slog.String("addr", srv.Addr))
					if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
						logger.Error("server stopped", slog.String("server", name), slog.Any("error", err))
					}
				}()
			}

			start("http", apiSrv)
			start("websocket", wsSrv)
			start("admin", adminSrv)

			return nil
		},
		OnStop: func(ctx context.Context) error {
			return errors.Join(
				apiSrv.Shutdown(ctx),
				wsSrv.Shutdown(ctx),
				adminSrv.Shutdown(ctx),
			)
		},
	})

	return &Application{lifecycle: lc}
}
