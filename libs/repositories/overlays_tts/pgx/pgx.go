package pgx

import (
	"context"
	"errors"
	"fmt"
	"strings"

	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/overlays_tts"
	ttsmodel "github.com/twirapp/twir/libs/repositories/overlays_tts/model"
)

type Opts struct {
	PgxPool *pgxpool.Pool
}

func New(opts Opts) *Pgx {
	return &Pgx{
		pool:   opts.PgxPool,
		getter: trmpgx.DefaultCtxGetter,
	}
}

func NewFx(pool *pgxpool.Pool) *Pgx {
	return New(Opts{PgxPool: pool})
}

var _ overlays_tts.Repository = (*Pgx)(nil)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (p *Pgx) GetByChannelID(ctx context.Context, channelID string) (
	ttsmodel.TTSOverlay,
	error,
) {
	query := `
SELECT
	id,
	channel_id,
	created_at,
	updated_at,
	enabled,
	voice,
	disallowed_voices,
	pitch,
	rate,
	volume,
	do_not_read_twitch_emotes,
	do_not_read_emoji,
	do_not_read_links,
	allow_users_choose_voice_in_main_command,
	max_symbols,
	read_chat_messages,
	read_chat_messages_nicknames
FROM channels_overlays_tts
WHERE channel_id = $1
LIMIT 1;
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	row := conn.QueryRow(ctx, query, channelID)
	overlay := ttsmodel.TTSOverlay{
		Settings: &ttsmodel.TTSOverlaySettings{},
	}

	err := row.Scan(
		&overlay.ID,
		&overlay.ChannelID,
		&overlay.CreatedAt,
		&overlay.UpdatedAt,
		&overlay.Settings.Enabled,
		&overlay.Settings.Voice,
		&overlay.Settings.DisallowedVoices,
		&overlay.Settings.Pitch,
		&overlay.Settings.Rate,
		&overlay.Settings.Volume,
		&overlay.Settings.DoNotReadTwitchEmotes,
		&overlay.Settings.DoNotReadEmoji,
		&overlay.Settings.DoNotReadLinks,
		&overlay.Settings.AllowUsersChooseVoiceInMainCommand,
		&overlay.Settings.MaxSymbols,
		&overlay.Settings.ReadChatMessages,
		&overlay.Settings.ReadChatMessagesNicknames,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ttsmodel.TTSOverlay{}, overlays_tts.ErrNotFound
		}
		return ttsmodel.TTSOverlay{}, fmt.Errorf(
			"tts overlay get by channel ID: %w",
			err,
		)
	}

	return overlay, nil
}

func (p *Pgx) Create(
	ctx context.Context,
	input overlays_tts.CreateInput,
) (ttsmodel.TTSOverlay, error) {
	query := `
INSERT INTO channels_overlays_tts (
	channel_id,
	enabled,
	voice,
	disallowed_voices,
	pitch,
	rate,
	volume,
	do_not_read_twitch_emotes,
	do_not_read_emoji,
	do_not_read_links,
	allow_users_choose_voice_in_main_command,
	max_symbols,
	read_chat_messages,
	read_chat_messages_nicknames
)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14);
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(
		ctx,
		query,
		input.ChannelID,
		input.Settings.Enabled,
		input.Settings.Voice,
		input.Settings.DisallowedVoices,
		input.Settings.Pitch,
		input.Settings.Rate,
		input.Settings.Volume,
		input.Settings.DoNotReadTwitchEmotes,
		input.Settings.DoNotReadEmoji,
		input.Settings.DoNotReadLinks,
		input.Settings.AllowUsersChooseVoiceInMainCommand,
		input.Settings.MaxSymbols,
		input.Settings.ReadChatMessages,
		input.Settings.ReadChatMessagesNicknames,
	)
	if err != nil {
		return ttsmodel.TTSOverlay{}, fmt.Errorf("tts overlay create: %w", err)
	}

	return p.GetByChannelID(ctx, input.ChannelID)
}

func (p *Pgx) Update(
	ctx context.Context,
	channelID string,
	input overlays_tts.UpdateInput,
) (ttsmodel.TTSOverlay, error) {
	query := `
UPDATE channels_overlays_tts
SET
	enabled = $1,
	voice = $2,
	disallowed_voices = $3,
	pitch = $4,
	rate = $5,
	volume = $6,
	do_not_read_twitch_emotes = $7,
	do_not_read_emoji = $8,
	do_not_read_links = $9,
	allow_users_choose_voice_in_main_command = $10,
	max_symbols = $11,
	read_chat_messages = $12,
	read_chat_messages_nicknames = $13,
	updated_at = now()
WHERE channel_id = $14
RETURNING channel_id
`

	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(
		ctx,
		query,
		input.Settings.Enabled,
		input.Settings.Voice,
		input.Settings.DisallowedVoices,
		input.Settings.Pitch,
		input.Settings.Rate,
		input.Settings.Volume,
		input.Settings.DoNotReadTwitchEmotes,
		input.Settings.DoNotReadEmoji,
		input.Settings.DoNotReadLinks,
		input.Settings.AllowUsersChooseVoiceInMainCommand,
		input.Settings.MaxSymbols,
		input.Settings.ReadChatMessages,
		input.Settings.ReadChatMessagesNicknames,
		channelID,
	)
	if err != nil {
		return ttsmodel.TTSOverlay{}, fmt.Errorf("tts overlay update: %w", err)
	}

	return p.GetByChannelID(ctx, channelID)
}

func (p *Pgx) GetOrCreate(
	ctx context.Context,
	channelID string,
) (ttsmodel.TTSOverlay, error) {
	overlay, err := p.GetByChannelID(ctx, channelID)
	if err == nil {
		return overlay, nil
	}
	if !errors.Is(err, overlays_tts.ErrNotFound) {
		return ttsmodel.TTSOverlay{}, err
	}
	// Create with default settings
	defaultSettings := ttsmodel.TTSOverlaySettings{
		Enabled:                            true,
		Voice:                              "aleksandr",
		DisallowedVoices:                   []string{},
		Pitch:                              0,
		Rate:                               0,
		Volume:                             70,
		DoNotReadTwitchEmotes:              true,
		DoNotReadEmoji:                     true,
		DoNotReadLinks:                     true,
		AllowUsersChooseVoiceInMainCommand: false,
		MaxSymbols:                         500,
		ReadChatMessages:                   false,
		ReadChatMessagesNicknames:          false,
	}
	return p.Create(
		ctx,
		overlays_tts.CreateInput{
			ChannelID: channelID,
			Settings:  defaultSettings,
		},
	)
}

func (p *Pgx) GetUserSettings(
	ctx context.Context,
	channelID, userID string,
) (ttsmodel.TTSUserSettings, error) {
	query := `
SELECT
	id,
	channel_id,
	user_id,
	voice,
	rate,
	pitch,
	created_at,
	updated_at
FROM channels_overlays_tts_users
WHERE channel_id = $1 AND user_id = $2
LIMIT 1;
`
	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	row := conn.QueryRow(ctx, query, channelID, userID)
	var settings ttsmodel.TTSUserSettings
	err := row.Scan(
		&settings.ID,
		&settings.ChannelID,
		&settings.UserID,
		&settings.Voice,
		&settings.Rate,
		&settings.Pitch,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ttsmodel.TTSUserSettings{}, overlays_tts.ErrNotFound
		}
		return ttsmodel.TTSUserSettings{}, fmt.Errorf(
			"tts user settings get by channel and user ID: %w",
			err,
		)
	}
	return settings, nil
}
func (p *Pgx) CreateUserSettings(
	ctx context.Context,
	input overlays_tts.CreateUserSettingsInput,
) (ttsmodel.TTSUserSettings, error) {
	query := `
INSERT INTO channels_overlays_tts_users (
channel_id,
user_id,
voice,
rate,
pitch
)
VALUES ($1, $2, $3, $4, $5)
RETURNING id, channel_id, user_id, voice, rate, pitch, created_at, updated_at;
`
	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	row := conn.QueryRow(
		ctx,
		query,
		input.ChannelID,
		input.UserID,
		input.Voice,
		input.Rate,
		input.Pitch,
	)
	var settings ttsmodel.TTSUserSettings
	err := row.Scan(
		&settings.ID,
		&settings.ChannelID,
		&settings.UserID,
		&settings.Voice,
		&settings.Rate,
		&settings.Pitch,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	if err != nil {
		return ttsmodel.TTSUserSettings{}, fmt.Errorf("tts user settings create: %w", err)
	}
	return settings, nil
}
func (p *Pgx) UpdateUserSettings(
	ctx context.Context,
	channelID, userID string,
	input overlays_tts.UpdateUserSettingsInput,
) (ttsmodel.TTSUserSettings, error) {
	updates := []string{}
	args := []interface{}{}
	argIndex := 1
	if input.Voice != nil {
		updates = append(updates, fmt.Sprintf("voice = $%d", argIndex))
		args = append(args, *input.Voice)
		argIndex++
	}
	if input.Rate != nil {
		updates = append(updates, fmt.Sprintf("rate = $%d", argIndex))
		args = append(args, *input.Rate)
		argIndex++
	}
	if input.Pitch != nil {
		updates = append(updates, fmt.Sprintf("pitch = $%d", argIndex))
		args = append(args, *input.Pitch)
		argIndex++
	}
	if len(updates) == 0 {
		return p.GetUserSettings(ctx, channelID, userID)
	}
	updates = append(updates, "updated_at = now()")
	query := fmt.Sprintf(
		`
UPDATE channels_overlays_tts_users
SET %s
WHERE channel_id = $%d AND user_id = $%d
RETURNING id, channel_id, user_id, voice, rate, pitch, created_at, updated_at;
`, strings.Join(updates, ", "), argIndex, argIndex+1,
	)
	args = append(args, channelID, userID)
	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	row := conn.QueryRow(ctx, query, args...)
	var settings ttsmodel.TTSUserSettings
	err := row.Scan(
		&settings.ID,
		&settings.ChannelID,
		&settings.UserID,
		&settings.Voice,
		&settings.Rate,
		&settings.Pitch,
		&settings.CreatedAt,
		&settings.UpdatedAt,
	)
	if err != nil {
		return ttsmodel.TTSUserSettings{}, fmt.Errorf("tts user settings update: %w", err)
	}
	return settings, nil
}
func (p *Pgx) GetAllUserSettings(
	ctx context.Context,
	channelID string,
) ([]ttsmodel.TTSUserSettings, error) {
	query := `
SELECT
	id,
	channel_id,
	user_id,
	voice,
	rate,
	pitch,
	created_at,
	updated_at
FROM channels_overlays_tts_users
WHERE channel_id = $1
ORDER BY user_id DESC;
`
	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	rows, err := conn.Query(ctx, query, channelID)
	if err != nil {
		return nil, fmt.Errorf("tts user settings get all: %w", err)
	}
	defer rows.Close()
	var settings []ttsmodel.TTSUserSettings
	for rows.Next() {
		var s ttsmodel.TTSUserSettings
		err := rows.Scan(
			&s.ID,
			&s.ChannelID,
			&s.UserID,
			&s.Voice,
			&s.Rate,
			&s.Pitch,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("tts user settings scan: %w", err)
		}
		settings = append(settings, s)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("tts user settings rows error: %w", err)
	}
	return settings, nil
}
func (p *Pgx) DeleteUserSettings(
	ctx context.Context,
	channelID, userID string,
) error {
	query := `
DELETE FROM channels_overlays_tts_users
WHERE channel_id = $1 AND user_id = $2;
`
	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(ctx, query, channelID, userID)
	if err != nil {
		return fmt.Errorf("tts user settings delete: %w", err)
	}
	return nil
}
func (p *Pgx) DeleteMultipleUserSettings(
	ctx context.Context,
	channelID string,
	userIDs []string,
) error {
	if len(userIDs) == 0 {
		return nil
	}
	query := `
DELETE FROM channels_overlays_tts_users
WHERE channel_id = $1 AND user_id = ANY($2);
`
	conn := p.getter.DefaultTrOrDB(ctx, p.pool)
	_, err := conn.Exec(ctx, query, channelID, userIDs)
	if err != nil {
		return fmt.Errorf("tts user settings delete multiple: %w", err)
	}
	return nil
}
func (p *Pgx) GetOrCreateUserSettings(
	ctx context.Context,
	channelID, userID string,
	defaults overlays_tts.CreateUserSettingsInput,
) (ttsmodel.TTSUserSettings, error) {
	settings, err := p.GetUserSettings(ctx, channelID, userID)
	if err == nil {
		return settings, nil
	}
	if !errors.Is(err, overlays_tts.ErrNotFound) {
		return ttsmodel.TTSUserSettings{}, err
	}
	if defaults.ChannelID == "" {
		defaults.ChannelID = channelID
	}
	if defaults.UserID == "" {
		defaults.UserID = userID
	}
	if defaults.Voice == "" {
		defaults.Voice = ""
	}
	if defaults.Rate == 0 {
		defaults.Rate = 50
	}
	if defaults.Pitch == 0 {
		defaults.Pitch = 50
	}
	return p.CreateUserSettings(ctx, defaults)
}
