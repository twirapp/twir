package tts

import (
	"context"

	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"github.com/twirapp/twir/libs/types/types/api/modules"
)

type Opts struct {
	TTSSettingsCacher *generic_cacher.GenericCacher[modules.TTSSettings]
}

func New(opts Opts) *Service {
	return &Service{
		ttsSettingsCacher: opts.TTSSettingsCacher,
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
