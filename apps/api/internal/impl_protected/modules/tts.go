package modules

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/imroc/req/v3"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/modules_tts"
	"github.com/satont/twir/libs/types/types/api/modules"
	"google.golang.org/protobuf/types/known/emptypb"
)

const TTSType = "tts"

func (c *Modules) ModulesTTSGet(ctx context.Context, empty *emptypb.Empty) (*modules_tts.GetResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? AND "type" = ? AND "userId" = ?`, dashboardId, TTSType, nil).
		First(entity).Error; err != nil {
		return nil, err
	}

	settings := &modules.TTSSettings{}
	if err := json.Unmarshal(entity.Settings, settings); err != nil {
		return nil, err
	}

	return &modules_tts.GetResponse{
		Data: &modules_tts.Settings{
			Enabled:                            *settings.Enabled,
			Rate:                               uint32(settings.Rate),
			Volume:                             uint32(settings.Volume),
			Pitch:                              uint32(settings.Pitch),
			Voice:                              settings.Voice,
			AllowUsersChooseVoiceInMainCommand: settings.AllowUsersChooseVoiceInMainCommand,
			MaxSymbols:                         uint32(settings.MaxSymbols),
			DisallowedVoices:                   settings.DisallowedVoices,
			DoNotReadEmoji:                     settings.DoNotReadEmoji,
			DoNotReadTwitchEmotes:              settings.DoNotReadTwitchEmotes,
			DoNotReadLinks:                     settings.DoNotReadLinks,
			ReadChatMessages:                   settings.ReadChatMessages,
			ReadChatMessagesNicknames:          settings.ReadChatMessagesNicknames,
		},
	}, nil
}

var ttsParseError = fmt.Errorf("internal error: can't parse tts microservice response")

func (c *Modules) ModulesTTSUpdate(
	ctx context.Context,
	request *modules_tts.PostRequest,
) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.ChannelModulesSettings{}
	if err := c.Db.
		WithContext(ctx).
		Where(`"channelId" = ? AND "type" = ? AND "userId" = ?`, dashboardId, TTSType, nil).
		First(entity).Error; err != nil {
		return nil, err
	}

	settings := &modules.TTSSettings{
		Enabled:                            &request.Data.Enabled,
		Rate:                               int(request.Data.Rate),
		Volume:                             int(request.Data.Volume),
		Pitch:                              int(request.Data.Pitch),
		Voice:                              request.Data.Voice,
		AllowUsersChooseVoiceInMainCommand: request.Data.AllowUsersChooseVoiceInMainCommand,
		MaxSymbols:                         int(request.Data.MaxSymbols),
		DisallowedVoices:                   request.Data.DisallowedVoices,
		DoNotReadEmoji:                     request.Data.DoNotReadEmoji,
		DoNotReadTwitchEmotes:              request.Data.DoNotReadTwitchEmotes,
		DoNotReadLinks:                     request.Data.DoNotReadLinks,
		ReadChatMessages:                   request.Data.ReadChatMessages,
		ReadChatMessagesNicknames:          request.Data.ReadChatMessagesNicknames,
	}
	bytes, err := json.Marshal(settings)
	if err != nil {
		return nil, err
	}

	entity.Settings = bytes
	if err := c.Db.WithContext(ctx).Save(entity).Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Modules) ModulesTTSGetInfo(
	ctx context.Context,
	_ *emptypb.Empty,
) (*modules_tts.GetInfoResponse, error) {
	result := map[string]any{}
	resp, err := req.
		R().
		SetContext(ctx).
		SetSuccessResult(&result).
		Get(fmt.Sprintf("http://%s/info", c.Config.TTSServiceUrl))
	if err != nil {
		return nil, err
	}
	if !resp.IsSuccessState() {
		return nil, fmt.Errorf("tts service is not available")
	}

	defaultFormat, ok := result["DEFAULT_FORMAT"].(string)
	if !ok {
		return nil, ttsParseError
	}

	defaultVoice, ok := result["DEFAULT_VOICE"].(string)
	if !ok {
		return nil, ttsParseError
	}

	formats := make(map[string]string)
	respFormats, ok := result["FORMATS"].(map[string]interface{})
	if !ok {
		return nil, ttsParseError
	}
	for k, v := range respFormats {
		formats[k] = v.(string)
	}

	supportedVoices := make([]string, 0)
	respSupportedVoices, ok := result["SUPPORTED_VOICES"].([]interface{})
	if !ok {
		return nil, ttsParseError
	}
	for _, v := range respSupportedVoices {
		supportedVoices = append(supportedVoices, v.(string))
	}

	return &modules_tts.GetInfoResponse{
		DefaultFormat: defaultFormat,
		DefaultVoice:  defaultVoice,
		Formats: &modules_tts.GetInfoResponse_Formats{
			Flac: formats["flac"],
			Mp3:  formats["mp3"],
			Opus: formats["opus"],
			Wav:  formats["wav"],
		},
		SupportedVoices: supportedVoices,
		VoicesInfo:      nil,
	}, nil
}
