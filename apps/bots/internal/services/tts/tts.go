package tts

import (
	"context"

	"github.com/satont/twir/libs/types/types/api/modules"
	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

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
