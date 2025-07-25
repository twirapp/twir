package seeds

import (
	"database/sql"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"github.com/lib/pq"
	cfg "github.com/twirapp/twir/libs/config"
	"github.com/twirapp/twir/libs/crypto"
	"io"
	"log/slog"
	"net/http"
	"time"
)

type TwitchResponse struct {
	UserID    string   `json:"user_id"`
	ExpiresIn int      `json:"expires_in"`
	Scopes    []string `json:"scopes"`
}

func CreateDefaultBot(db *sql.DB, config *cfg.Config) error {
	row := db.QueryRow(`SELECT "id" FROM bots where type = $1`, "DEFAULT")

	// check is row exists
	var id string
	err := row.Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return err
	}

	if id != "" {
		slog.Info("✅ Bot already exists, skipping...")
		return nil
	}

	if config.BotAccessToken == "" || config.BotRefreshToken == "" {
		return errors.New("🚨 Missed bot access token or bot refresh token")
	}

	accessToken, err := crypto.Encrypt(config.BotAccessToken, config.TokensCipherKey)
	if err != nil {
		return err
	}
	refreshToken, err := crypto.Encrypt(config.BotRefreshToken, config.TokensCipherKey)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", "https://id.twitch.tv/oauth2/validate", nil)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "OAuth "+config.BotAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		panic("🚨 Invalid bot access token " + string(body))
	}

	token := TwitchResponse{}
	if err = json.Unmarshal(body, &token); err != nil {
		return err
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
		return err
	}

	var tokenId string
	for rows.Next() {
		if err := rows.Scan(&tokenId); err != nil {
			return err
		}
	}

	if tokenId == "" {
		panic("🚨 Failed to create bot access token")
	}

	_, err = db.Exec(`INSERT INTO "bots" ("id", "type", "tokenId") VALUES ($1, $2, $3)`, token.UserID, "DEFAULT", tokenId)
	if err != nil {
		return err
	}

	slog.Info("✅ Bot access token is created")

	return nil
}
