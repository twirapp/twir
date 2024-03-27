package seeds

import (
	"database/sql"
	"log/slog"

	cfg "github.com/satont/twir/libs/config"
)

func CreateIntegrations(db *sql.DB, config *cfg.Config) error {
	_, err := db.Query(
		`INSERT INTO integrations (service) VALUES ($1) ON CONFLICT DO NOTHING`,
		"DONATEPAY",
	)
	if err != nil {
		return err
	}

	_, err = db.Query(
		`INSERT INTO integrations (service) VALUES ($1) ON CONFLICT DO NOTHING`,
		"VALORANT",
	)
	if err != nil {
		return err
	}

	_, err = db.Query(
		`INSERT INTO integrations (service) VALUES ($1) ON CONFLICT DO NOTHING`,
		"DONATE_STREAM",
	)
	if err != nil {
		return err
	}

	_, err = db.Query(
		`INSERT INTO integrations (service) VALUES ($1) ON CONFLICT DO NOTHING`,
		"DONATELLO",
	)
	if err != nil {
		return err
	}

	_, err = db.Query(
		`INSERT INTO integrations (service) VALUES ($1) ON CONFLICT DO NOTHING`,
		"NIGHTBOT",
	)
	if err != nil {
		return err
	}

	slog.Info("Integrations created")

	return nil
}
