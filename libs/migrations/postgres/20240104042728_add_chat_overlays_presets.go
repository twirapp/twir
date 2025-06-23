package postgres

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/pressly/goose/v3"
)

func init() {
	goose.AddMigrationContext(upAddChatOverlaysPresets, downAddChatOverlaysPresets)
}

func upAddChatOverlaysPresets(ctx context.Context, tx *sql.Tx) error {
	_, err := tx.Query(
		`
CREATE TABLE IF NOT EXISTS channels_overlays_chat (
	id uuid PRIMARY KEY DEFAULT uuid_generate_v4(),
	message_hide_timeout integer NOT NULL,
	message_show_delay integer NOT NULL,
	preset text NOT NULL,
	font_family text NOT NULL,
	font_size integer NOT NULL,
	font_weight integer NOT NULL,
	font_style text NOT NULL,
	hide_commands boolean NOT NULL,
	hide_bots boolean NOT NULL,
	show_badges boolean NOT NULL,
	show_announce_badge boolean NOT NULL,
	text_shadow_color text NOT NULL,
	text_shadow_size integer NOT NULL,
	chat_background_color text NOT NULL,
	direction text NOT NULL,
	channel_id text NOT NULL,
	created_at timestamp NOT NULL default now(),
	FOREIGN KEY (channel_id) REFERENCES channels(id) ON DELETE CASCADE
);
`,
	)

	if err != nil {
		return err
	}

	existedSettingsQuery := `
SELECT id, settings, "channelId" FROM channels_modules_settings WHERE type = 'chat_overlay'
`

	rows, err := tx.QueryContext(ctx, existedSettingsQuery)
	if err != nil {
		return err
	}
	defer rows.Close()

	var existedSettings []struct {
		id            string
		channelId     string
		settingsBytes []byte
	}

	for rows.Next() {
		var id string
		var settingsBytes []byte
		var channelId string
		if err := rows.Scan(&id, &settingsBytes, &channelId); err != nil {
			return err
		}

		existedSettings = append(
			existedSettings, struct {
				id            string
				channelId     string
				settingsBytes []byte
			}{
				id:            id,
				settingsBytes: settingsBytes,
				channelId:     channelId,
			},
		)
	}

	for _, entity := range existedSettings {
		var oldSettings AddChatOverlaysPresets20240104042728
		if err := json.Unmarshal(entity.settingsBytes, &oldSettings); err != nil {
			return err
		}

		_, err = tx.Exec(
			`
INSERT INTO channels_overlays_chat(
	message_hide_timeout,
	message_show_delay,
	preset,
	font_family,
	font_size,
	font_weight,
	font_style,
	hide_commands,
	hide_bots,
	show_badges,
	show_announce_badge,
	text_shadow_color,
	text_shadow_size,
	chat_background_color,
	direction,
	channel_id
) VALUES (
	$1,
	$2,
	$3,
	$4,
	$5,
	$6,
	$7,
	$8,
	$9,
	$10,
	$11,
	$12,
	$13,
	$14,
	$15,
	$16
);
`,
			oldSettings.MessageHideTimeout,
			oldSettings.MessageShowDelay,
			oldSettings.Preset,
			oldSettings.FontFamily,
			oldSettings.FontSize,
			oldSettings.FontWeight,
			oldSettings.FontStyle,
			oldSettings.HideCommands,
			oldSettings.HideBots,
			oldSettings.ShowBadges,
			oldSettings.ShowAnnounceBadge,
			oldSettings.TextShadowColor,
			oldSettings.TextShadowSize,
			oldSettings.ChatBackgroundColor,
			oldSettings.Direction,
			entity.channelId,
		)

		if err != nil {
			return err
		}

		_, err := tx.Exec(`DELETE FROM channels_modules_settings WHERE id = $1`, entity.id)
		if err != nil {
			return err
		}
	}

	return nil
}

func downAddChatOverlaysPresets(ctx context.Context, tx *sql.Tx) error {
	// This code is executed when the migration is rolled back.
	return nil
}

type AddChatOverlaysPresets20240104042728 struct {
	Preset              string `json:"preset"`
	FontFamily          string `json:"fontFamily"`
	FontStyle           string `json:"fontStyle"`
	TextShadowColor     string `json:"textShadowColor"`
	ChatBackgroundColor string `json:"chatBackgroundColor"`
	Direction           string `json:"direction"`
	MessageHideTimeout  uint32 `json:"messageHideTimeout"`
	MessageShowDelay    uint32 `json:"messageShowDelay"`
	FontSize            uint32 `json:"fontSize"`
	FontWeight          uint32 `json:"fontWeight"`
	TextShadowSize      uint32 `json:"textShadowSize"`
	HideCommands        bool   `json:"hideCommands"`
	HideBots            bool   `json:"hideBots"`
	ShowBadges          bool   `json:"showBadges"`
	ShowAnnounceBadge   bool   `json:"showAnnounceBadge"`
}
