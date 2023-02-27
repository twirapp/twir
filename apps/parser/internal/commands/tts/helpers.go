package tts

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/guregu/null"
	"github.com/imroc/req/v3"
	"github.com/samber/do"
	"github.com/samber/lo"
	"github.com/satont/tsuwari/apps/parser/internal/di"
	config "github.com/satont/tsuwari/libs/config"
	model "github.com/satont/tsuwari/libs/gomodels"
	"github.com/satont/tsuwari/libs/types/types/api/modules"
	"github.com/tidwall/gjson"
	"gorm.io/gorm"
)

func getSettings(channelId, userId string) (*modules.TTSSettings, *model.ChannelModulesSettings) {
	db := do.MustInvoke[gorm.DB](di.Provider)

	settings := &model.ChannelModulesSettings{}
	query := db.
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

type voice struct {
	Name    string
	Country string
}

func getVoices() []voice {
	cfg := do.MustInvoke[config.Config](di.Provider)
	data := map[string]any{}
	_, err := req.R().SetSuccessResult(&data).Get(fmt.Sprintf("http://%s/info", cfg.TTSServiceUrl))
	if err != nil {
		fmt.Println(err)
		return nil
	}

	bytes, err := json.Marshal(&data)
	if err != nil {
		fmt.Println(err)
		return nil
	}

	parsedJson := gjson.ParseBytes(bytes)
	voices := []voice{}
	parsedJson.Get("rhvoice_wrapper_voices_info").ForEach(func(key, value gjson.Result) bool {
		voices = append(voices, voice{
			Name:    key.String(),
			Country: value.Get("country").String(),
		})

		return true
	})

	return voices
}

func updateSettings(entity *model.ChannelModulesSettings, settings *modules.TTSSettings) error {
	db := do.MustInvoke[gorm.DB](di.Provider)

	bytes, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	return db.Model(entity).Updates(map[string]interface{}{"settings": bytes}).Error
}

func createUserSettings(rate, pitch int, voice, channelId, userId string) (
	*model.ChannelModulesSettings,
	*modules.TTSSettings,
	error,
) {
	db := do.MustInvoke[gorm.DB](di.Provider)
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

	err = db.Create(userModel).Error
	if err != nil {
		return nil, nil, err
	}

	return userModel, userSettings, nil
}

func switchEnableState(channelId string, newState bool) error {
	channelSettings, channelModele := getSettings(channelId, "")

	if channelSettings == nil {
		return errors.New("Tts not configured")
	}

	channelSettings.Enabled = &newState
	err := updateSettings(channelModele, channelSettings)
	if err != nil {
		return err
	}

	return nil
}
