package tts

import (
	"context"

	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/types/types/api/modules"
)

func New(ttsSettingsCacher *generic_cacher.GenericCacher[modules.TTSSettings]) *Service {
	return &Service{
		ttsSettingsCacher: ttsSettingsCacher,
	}
}

type Service struct {
	ttsSettingsCacher *generic_cacher.GenericCacher[modules.TTSSettings]
}

func (s *Service) GetChannelTTSSettings(ctx context.Context, channelID string) (
	modules.TTSSettings,
	error,
) {
	ttsSettings, err := s.ttsSettingsCacher.Get(ctx, channelID)
	if err != nil {
		return modules.TTSSettings{}, err
	}

	return ttsSettings, nil
}
