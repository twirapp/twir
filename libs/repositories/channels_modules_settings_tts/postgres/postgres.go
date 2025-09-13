package postgres

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/Masterminds/squirrel"
	trmpgx "github.com/avito-tech/go-transaction-manager/drivers/pgxv5/v2"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/twirapp/twir/libs/repositories/channels_modules_settings_tts"
	"github.com/twirapp/twir/libs/repositories/channels_modules_settings_tts/model"
	ttstypes "github.com/twirapp/twir/libs/types/types/api/modules"
)

const TTSModuleType = "tts"

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

var _ channels_modules_settings_tts.Repository = (*Pgx)(nil)
var sq = squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)

type Pgx struct {
	pool   *pgxpool.Pool
	getter *trmpgx.CtxGetter
}

func (c *Pgx) scanToModel(dbRow dbResult) (model.ChannelModulesSettingsTTS, error) {
	var ttsSettings ttstypes.TTSSettings
	if err := json.Unmarshal(dbRow.Settings, &ttsSettings); err != nil {
		return model.Nil, err
	}

	return model.ChannelModulesSettingsTTS{
		ID:                                 dbRow.ID,
		ChannelID:                          dbRow.ChannelID,
		Type:                               dbRow.Type,
		UserID:                             dbRow.UserID,
		Enabled:                            ttsSettings.Enabled,
		Rate:                               ttsSettings.Rate,
		Volume:                             ttsSettings.Volume,
		Pitch:                              ttsSettings.Pitch,
		Voice:                              ttsSettings.Voice,
		AllowUsersChooseVoiceInMainCommand: ttsSettings.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         ttsSettings.MaxSymbols,
		DisallowedVoices:                   ttsSettings.DisallowedVoices,
		DoNotReadEmoji:                     ttsSettings.DoNotReadEmoji,
		DoNotReadTwitchEmotes:              ttsSettings.DoNotReadTwitchEmotes,
		DoNotReadLinks:                     ttsSettings.DoNotReadLinks,
		ReadChatMessages:                   ttsSettings.ReadChatMessages,
		ReadChatMessagesNicknames:          ttsSettings.ReadChatMessagesNicknames,
	}, nil
}

type dbResult struct {
	ID        string          `db:"id"`
	Settings  json.RawMessage `db:"settings"`
	Type      string          `db:"type"`
	ChannelID string          `db:"channelId"`
	UserID    *string         `db:"userId"`
}

// Channel-level operations (userId is null)
func (c *Pgx) GetByChannelID(
	ctx context.Context,
	channelID string,
) (model.ChannelModulesSettingsTTS, error) {
	query, args, err := sq.Select("id", "type", "settings", `"channelId"`, `"userId"`).
		From("channels_modules_settings").
		Where(squirrel.Eq{`"channelId"`: channelID, "type": TTSModuleType, `"userId"`: nil}).
		ToSql()
	if err != nil {
		return model.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	dbRow, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[dbResult])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, channels_modules_settings_tts.ErrNotFound
		}
		return model.Nil, err
	}

	return c.scanToModel(dbRow)
}

func (c *Pgx) CreateForChannel(
	ctx context.Context,
	input channels_modules_settings_tts.CreateOrUpdateInput,
) (model.ChannelModulesSettingsTTS, error) {
	id := uuid.New().String()

	ttsSettings := ttstypes.TTSSettings{
		Enabled:                            input.Enabled,
		Rate:                               input.Rate,
		Volume:                             input.Volume,
		Pitch:                              input.Pitch,
		Voice:                              input.Voice,
		AllowUsersChooseVoiceInMainCommand: input.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         input.MaxSymbols,
		DisallowedVoices:                   input.DisallowedVoices,
		DoNotReadEmoji:                     input.DoNotReadEmoji,
		DoNotReadTwitchEmotes:              input.DoNotReadTwitchEmotes,
		DoNotReadLinks:                     input.DoNotReadLinks,
		ReadChatMessages:                   input.ReadChatMessages,
		ReadChatMessagesNicknames:          input.ReadChatMessagesNicknames,
	}

	settingsJSON, err := json.Marshal(ttsSettings)
	if err != nil {
		return model.Nil, err
	}

	query, args, err := sq.Insert("channels_modules_settings").
		Columns("id", "type", "settings", `"channelId"`, `"userId"`).
		Values(id, TTSModuleType, string(settingsJSON), input.ChannelID, nil).
		ToSql()
	if err != nil {
		return model.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	return model.ChannelModulesSettingsTTS{
		ID:                                 id,
		ChannelID:                          input.ChannelID,
		UserID:                             nil,
		Enabled:                            input.Enabled,
		Rate:                               input.Rate,
		Volume:                             input.Volume,
		Pitch:                              input.Pitch,
		Voice:                              input.Voice,
		AllowUsersChooseVoiceInMainCommand: input.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         input.MaxSymbols,
		DisallowedVoices:                   input.DisallowedVoices,
		DoNotReadEmoji:                     input.DoNotReadEmoji,
		DoNotReadTwitchEmotes:              input.DoNotReadTwitchEmotes,
		DoNotReadLinks:                     input.DoNotReadLinks,
		ReadChatMessages:                   input.ReadChatMessages,
		ReadChatMessagesNicknames:          input.ReadChatMessagesNicknames,
	}, nil
}

func (c *Pgx) UpdateForChannel(
	ctx context.Context,
	channelID string,
	input channels_modules_settings_tts.CreateOrUpdateInput,
) (model.ChannelModulesSettingsTTS, error) {
	ttsSettings := ttstypes.TTSSettings{
		Enabled:                            input.Enabled,
		Rate:                               input.Rate,
		Volume:                             input.Volume,
		Pitch:                              input.Pitch,
		Voice:                              input.Voice,
		AllowUsersChooseVoiceInMainCommand: input.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         input.MaxSymbols,
		DisallowedVoices:                   input.DisallowedVoices,
		DoNotReadEmoji:                     input.DoNotReadEmoji,
		DoNotReadTwitchEmotes:              input.DoNotReadTwitchEmotes,
		DoNotReadLinks:                     input.DoNotReadLinks,
		ReadChatMessages:                   input.ReadChatMessages,
		ReadChatMessagesNicknames:          input.ReadChatMessagesNicknames,
	}

	settingsJSON, err := json.Marshal(ttsSettings)
	if err != nil {
		return model.Nil, err
	}

	query, args, err := sq.Update("channels_modules_settings").
		Set("settings", string(settingsJSON)).
		Where(squirrel.Eq{`"channelId"`: channelID, "type": TTSModuleType, `"userId"`: nil}).
		ToSql()
	if err != nil {
		return model.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	result, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	if result.RowsAffected() == 0 {
		return model.Nil, channels_modules_settings_tts.ErrNotFound
	}

	return c.GetByChannelID(ctx, channelID)
}

func (c *Pgx) DeleteForChannel(ctx context.Context, channelID string) error {
	query, args, err := sq.Delete("channels_modules_settings").
		Where(squirrel.Eq{`"channelId"`: channelID, "type": TTSModuleType, `"userId"`: nil}).
		ToSql()
	if err != nil {
		return err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	result, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return channels_modules_settings_tts.ErrNotFound
	}

	return nil
}

// User-specific operations (userId is set)
func (c *Pgx) GetByChannelIDAndUserID(
	ctx context.Context,
	channelID, userID string,
) (model.ChannelModulesSettingsTTS, error) {
	query, args, err := sq.Select("id", "type", "settings", `"channelId"`, `"userId"`).
		From("channels_modules_settings").
		Where(squirrel.Eq{`"channelId"`: channelID, `"userId"`: userID, "type": TTSModuleType}).
		ToSql()
	if err != nil {
		return model.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	dbRow, err := pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[dbResult])
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return model.Nil, channels_modules_settings_tts.ErrNotFound
		}
		return model.Nil, err
	}

	return c.scanToModel(dbRow)
}

func (c *Pgx) CreateForUser(
	ctx context.Context,
	input channels_modules_settings_tts.CreateOrUpdateInput,
) (model.ChannelModulesSettingsTTS, error) {
	if input.UserID == nil || *input.UserID == "" {
		return model.Nil, errors.New("userID is required for user-specific settings")
	}

	id := uuid.New().String()

	ttsSettings := ttstypes.TTSSettings{
		Enabled:                            input.Enabled,
		Rate:                               input.Rate,
		Volume:                             input.Volume,
		Pitch:                              input.Pitch,
		Voice:                              input.Voice,
		AllowUsersChooseVoiceInMainCommand: input.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         input.MaxSymbols,
		DisallowedVoices:                   input.DisallowedVoices,
		DoNotReadEmoji:                     input.DoNotReadEmoji,
		DoNotReadTwitchEmotes:              input.DoNotReadTwitchEmotes,
		DoNotReadLinks:                     input.DoNotReadLinks,
		ReadChatMessages:                   input.ReadChatMessages,
		ReadChatMessagesNicknames:          input.ReadChatMessagesNicknames,
	}

	settingsJSON, err := json.Marshal(ttsSettings)
	if err != nil {
		return model.Nil, err
	}

	query, args, err := sq.Insert("channels_modules_settings").
		Columns("id", "type", "settings", `"channelId"`, `"userId"`).
		Values(id, TTSModuleType, string(settingsJSON), input.ChannelID, input.UserID).
		ToSql()
	if err != nil {
		return model.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	_, err = conn.Exec(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	return model.ChannelModulesSettingsTTS{
		ID:                                 id,
		ChannelID:                          input.ChannelID,
		UserID:                             input.UserID,
		Enabled:                            input.Enabled,
		Rate:                               input.Rate,
		Volume:                             input.Volume,
		Pitch:                              input.Pitch,
		Voice:                              input.Voice,
		AllowUsersChooseVoiceInMainCommand: input.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         input.MaxSymbols,
		DisallowedVoices:                   input.DisallowedVoices,
		DoNotReadEmoji:                     input.DoNotReadEmoji,
		DoNotReadTwitchEmotes:              input.DoNotReadTwitchEmotes,
		DoNotReadLinks:                     input.DoNotReadLinks,
		ReadChatMessages:                   input.ReadChatMessages,
		ReadChatMessagesNicknames:          input.ReadChatMessagesNicknames,
	}, nil
}

func (c *Pgx) UpdateForUser(
	ctx context.Context,
	channelID, userID string,
	input channels_modules_settings_tts.CreateOrUpdateInput,
) (model.ChannelModulesSettingsTTS, error) {
	ttsSettings := ttstypes.TTSSettings{
		Enabled:                            input.Enabled,
		Rate:                               input.Rate,
		Volume:                             input.Volume,
		Pitch:                              input.Pitch,
		Voice:                              input.Voice,
		AllowUsersChooseVoiceInMainCommand: input.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         input.MaxSymbols,
		DisallowedVoices:                   input.DisallowedVoices,
		DoNotReadEmoji:                     input.DoNotReadEmoji,
		DoNotReadTwitchEmotes:              input.DoNotReadTwitchEmotes,
		DoNotReadLinks:                     input.DoNotReadLinks,
		ReadChatMessages:                   input.ReadChatMessages,
		ReadChatMessagesNicknames:          input.ReadChatMessagesNicknames,
	}

	settingsJSON, err := json.Marshal(ttsSettings)
	if err != nil {
		return model.Nil, err
	}

	query, args, err := sq.Update("channels_modules_settings").
		Set("settings", string(settingsJSON)).
		Where(squirrel.Eq{`"channelId"`: channelID, `"userId"`: userID, "type": TTSModuleType}).
		ToSql()
	if err != nil {
		return model.Nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	result, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return model.Nil, err
	}

	if result.RowsAffected() == 0 {
		return model.Nil, channels_modules_settings_tts.ErrNotFound
	}

	return c.GetByChannelIDAndUserID(ctx, channelID, userID)
}

func (c *Pgx) DeleteForUser(ctx context.Context, channelID, userID string) error {
	query, args, err := sq.Delete("channels_modules_settings").
		Where(squirrel.Eq{`"channelId"`: channelID, `"userId"`: userID, "type": TTSModuleType}).
		ToSql()
	if err != nil {
		return err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	result, err := conn.Exec(ctx, query, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return channels_modules_settings_tts.ErrNotFound
	}

	return nil
}

// Get all settings for a channel (both channel-level and user-specific)
func (c *Pgx) GetAllByChannelID(
	ctx context.Context,
	channelID string,
) ([]model.ChannelModulesSettingsTTS, error) {
	query, args, err := sq.Select("id", "type", "settings", `"channelId"`, `"userId"`).
		From("channels_modules_settings").
		Where(squirrel.Eq{`"channelId"`: channelID, "type": TTSModuleType}).
		ToSql()
	if err != nil {
		return nil, err
	}

	conn := c.getter.DefaultTrOrDB(ctx, c.pool)
	rows, err := conn.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}

	dbRows, err := pgx.CollectRows(rows, pgx.RowToStructByName[dbResult])
	if err != nil {
		return nil, err
	}

	result := make([]model.ChannelModulesSettingsTTS, 0, len(dbRows))
	for _, dbRow := range dbRows {
		settings, err := c.scanToModel(dbRow)
		if err != nil {
			return nil, err
		}
		result = append(result, settings)
	}

	return result, nil
}
