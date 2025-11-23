package tts

import (
	"context"
	"errors"
	"time"

	"github.com/samber/lo"
	"github.com/twirapp/kv"
	"github.com/twirapp/twir/libs/repositories/overlays_tts"
	"github.com/twirapp/twir/libs/types/types/api/modules"

	generic_cacher "github.com/twirapp/twir/libs/cache/generic-cacher"
)

func NewTTSSettings(
	repository overlays_tts.Repository,
	kv kv.KV,
) *generic_cacher.GenericCacher[modules.TTSSettings] {
	return generic_cacher.New[modules.TTSSettings](
		generic_cacher.Opts[modules.TTSSettings]{
			KV:        kv,
			KeyPrefix: "cache:twir:tts-settings:channel:",
			LoadFn: func(ctx context.Context, key string) (modules.TTSSettings, error) {
				overlay, err := repository.GetByChannelID(ctx, key)
				if err != nil {
					if errors.Is(err, overlays_tts.ErrNotFound) {
						return modules.TTSSettings{}, err
					}
					return modules.TTSSettings{}, err
				}

				if overlay.Settings == nil {
					return modules.TTSSettings{}, overlays_tts.ErrNotFound
				}

				data := modules.TTSSettings{
					Enabled:                            lo.ToPtr(overlay.Settings.Enabled),
					Voice:                              overlay.Settings.Voice,
					DisallowedVoices:                   overlay.Settings.DisallowedVoices,
					Pitch:                              int(overlay.Settings.Pitch),
					Rate:                               int(overlay.Settings.Rate),
					Volume:                             int(overlay.Settings.Volume),
					DoNotReadTwitchEmotes:              overlay.Settings.DoNotReadTwitchEmotes,
					DoNotReadEmoji:                     overlay.Settings.DoNotReadEmoji,
					DoNotReadLinks:                     overlay.Settings.DoNotReadLinks,
					AllowUsersChooseVoiceInMainCommand: overlay.Settings.AllowUsersChooseVoiceInMainCommand,
					MaxSymbols:                         int(overlay.Settings.MaxSymbols),
					ReadChatMessages:                   overlay.Settings.ReadChatMessages,
					ReadChatMessagesNicknames:          overlay.Settings.ReadChatMessagesNicknames,
				}

				return data, nil
			},
			Ttl: 24 * time.Hour,
		},
	)
}
