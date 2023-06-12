package tts

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/satont/tsuwari/apps/bots/pkg/tlds"
	"net/url"
	"strings"

	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	"github.com/samber/lo"
	config "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/types/types/api/modules"
	"github.com/tidwall/gjson"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func getSettings(
	ctx context.Context,
	db *gorm.DB,
	channelId, userId string,
) (*modules.TTSSettings, *model.ChannelModulesSettings) {
	settings := &model.ChannelModulesSettings{}
	query := db.
		WithContext(ctx).
		Where(`"channelId" = ?`, channelId).
		Where(`"type" = ?`, "tts")

	if userId != "" {
		query = query.Where(`"userId" = ?`, userId)
	} else {
		query = query.Where(`"userId" IS NULL`)
	}

	err := query.First(&settings).Error
	if err != nil {
		return nil, nil
	}

	data := modules.TTSSettings{}
	err = json.Unmarshal(settings.Settings, &data)
	if err != nil {
		return nil, nil
	}

	return &data, settings
}

type Voice struct {
	Name    string
	Country string
}

func getVoices(ctx context.Context, cfg *config.Config) []Voice {
	data := map[string]any{}
	_, err := req.R().
		SetContext(ctx).
		SetSuccessResult(&data).
		Get(fmt.Sprintf("http://%s/info", cfg.TTSServiceUrl))
	if err != nil {
		zap.S().Error(err)
		return nil
	}

	bytes, err := json.Marshal(&data)
	if err != nil {
		zap.S().Error(err)
		return nil
	}

	parsedJson := gjson.ParseBytes(bytes)
	voices := []Voice{}
	parsedJson.Get("rhvoice_wrapper_voices_info").ForEach(func(key, value gjson.Result) bool {
		voices = append(voices, Voice{
			Name:    key.String(),
			Country: value.Get("country").String(),
		})

		return true
	})

	return voices
}

func updateSettings(ctx context.Context, db *gorm.DB, entity *model.ChannelModulesSettings, settings *modules.TTSSettings) error {
	bytes, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	return db.
		Model(entity).
		WithContext(ctx).
		Updates(map[string]interface{}{"settings": bytes}).
		Error
}

func createUserSettings(
	ctx context.Context,
	db *gorm.DB,
	rate, pitch int,
	voice, channelId, userId string,
) (
	*model.ChannelModulesSettings,
	*modules.TTSSettings,
	error,
) {
	userModel := &model.ChannelModulesSettings{
		ID:        uuid.New().String(),
		Type:      "tts",
		Settings:  nil,
		ChannelId: channelId,
		UserId:    null.StringFrom(userId),
	}

	userSettings := &modules.TTSSettings{
		Enabled: lo.ToPtr(true),
		Rate:    rate,
		Volume:  70,
		Pitch:   pitch,
		Voice:   voice,
	}

	bytes, err := json.Marshal(userSettings)
	if err != nil {
		return nil, nil, err
	}

	userModel.Settings = bytes

	err = db.WithContext(ctx).Create(userModel).Error
	if err != nil {
		return nil, nil, err
	}

	return userModel, userSettings, nil
}

func switchEnableState(ctx context.Context, db *gorm.DB, channelId string, newState bool) error {
	channelSettings, channelModele := getSettings(ctx, db, channelId, "")

	if channelSettings == nil {
		return errors.New("Tts not configured")
	}

	channelSettings.Enabled = &newState
	err := updateSettings(ctx, db, channelModele, channelSettings)
	if err != nil {
		return err
	}

	return nil
}

func IsUrl(str string) bool {
	u, err := url.Parse(str)
	if err != nil {
		return false
	}

	return u.Scheme != "" && u.Host != "" && lo.SomeBy(tlds.TLDS, func(item string) bool {
		return strings.HasSuffix(u.Host, item)
	})
}
