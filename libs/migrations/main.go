package main

import (
	"database/sql"
	"embed"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/pressly/goose/v3"
	"github.com/satont/twir/libs/config"
	"github.com/satont/twir/libs/crypto"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
	"time"
)

// TODO: add .go when there will be some go files in migrations
//
//go:embed migrations/*.sql
var embedMigrations embed.FS

const driver = "postgres"

func main() {
	config, err := cfg.New()
	if err != nil {
		panic(err)
	}

	opts, err := pq.ParseURL(config.DatabaseUrl)
	if err != nil {
		panic(err)
	}

	db, err := sql.Open(driver, opts)

	goose.SetBaseFS(embedMigrations)

	if err := goose.SetDialect(driver); err != nil {
		panic(err)
	}

	if err := goose.Up(db, "migrations"); err != nil {
		panic(err)
	}

	seedBot(db, config)
}

type TwitchResponse struct {
	UserID    string   `json:"user_id"`
	ExpiresIn int      `json:"expires_in"`
	Scopes    []string `json:"scopes"`
}

func seedBot(db *sql.DB, config *cfg.Config) {
	row := db.QueryRow(`SELECT "id" FROM bots where type = $1`, "DEFAULT")

	// check is row exists
	var id string
	err := row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		panic(err)
	}

	if id != "" {
		slog.Info("✅ Bot already exists, skipping...")
		return
	}

	if config.BotAccessToken == "" || config.BotRefreshToken == "" {
		panic("🚨 Missed bot access token or bot refresh token")
	}

	accessToken, err := crypto.Encrypt(config.BotAccessToken, config.TokensCipherKey)
	if err != nil {
		panic(err)
	}
	refreshToken, err := crypto.Encrypt(config.BotRefreshToken, config.TokensCipherKey)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("GET", "https://id.twitch.tv/oauth2/validate", nil)
	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "OAuth "+config.BotAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	if resp.StatusCode != 200 {
		panic("🚨 Invalid bot access token " + string(body))
	}

	token := TwitchResponse{}
	if err = json.Unmarshal(body, &token); err != nil {
		panic(err)
	}

	scopes := pq.StringArray{}

	for _, scope := range token.Scopes {
		scopes = append(scopes, scope)
	}

	rows, err := db.Query(
		`INSERT INTO "tokens" (
			"id",
			"accessToken",
			"refreshToken",
			"expiresIn",
			"obtainmentTimestamp",
			"scopes"
			) VALUES ($1, $2, $3, $4, $5, $6) RETURNING "id"`,
		uuid.New().String(),
		accessToken,
		refreshToken,
		token.ExpiresIn,
		time.Now().UTC(),
		scopes,
	)
	if err != nil {
		panic(err)
	}

	var tokenId string
	for rows.Next() {
		if err := rows.Scan(&tokenId); err != nil {
			panic(err)
		}
	}

	if tokenId == "" {
		panic("🚨 Failed to create bot access token")
	}

	_, err = db.Exec(`INSERT INTO "bots" ("id", "type", "tokenId") VALUES ($1, $2, $3)`, token.UserID, "DEFAULT", tokenId)
	if err != nil {
		panic(err)
	}

	slog.Info("✅ Bot access token is created")
}
