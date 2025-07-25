package clickhouse

import (
	"context"
	"os"

	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/pressly/goose/v3"
	"github.com/pterm/pterm"
	cfg "github.com/twirapp/twir/libs/config"
	_ "github.com/twirapp/twir/libs/migrations/clickhouse"
)

func Migrate(ctx context.Context, config *cfg.Config, migrationsPath string) error {
	dbOptions, err := clickhouse.ParseDSN(config.ClickhouseUrl)
	if err != nil {
		return err
	}

	goose.ResetGlobalMigrations()

	db := clickhouse.OpenDB(dbOptions)
	if err := goose.SetDialect(string(goose.DialectClickHouse)); err != nil {
		return err
	}

	provider, err := goose.NewProvider(
		goose.DialectClickHouse,
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

	return nil
}
