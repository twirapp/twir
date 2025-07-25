package postgres

import (
	"context"
	"database/sql"
	"log/slog"
	"os"

	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/pterm/pterm"
	cfg "github.com/twirapp/twir/libs/config"
	_ "github.com/twirapp/twir/libs/migrations/postgres"
	"github.com/twirapp/twir/libs/migrations/seeds"
)

func Migrate(ctx context.Context, config *cfg.Config, migrationsPath string) error {
	opts, err := pq.ParseURL(config.DatabaseUrl)
	if err != nil {
		return err
	}

	db, err := sql.Open(string(goose.DialectPostgres), opts)
	if err != nil {
		return err
	}

	if err := goose.SetDialect(string(goose.DialectPostgres)); err != nil {
		return err
	}

	provider, err := goose.NewProvider(
		goose.DialectPostgres,
		db,
		os.DirFS(migrationsPath),
		goose.WithAllowOutofOrder(true),
	)
	if err != nil {
		return err
	}

	if _, err := provider.Up(ctx); err != nil {
		pterm.Error.Println(err)
		return err
	}

	slog.SetDefault(slog.New(pterm.NewSlogHandler(&pterm.DefaultLogger)))

	if err := seeds.CreateDefaultBot(db, config); err != nil {
		return err
	}

	if err := seeds.CreateIntegrations(db, config); err != nil {
		return err
	}

	return nil
}
