package modules

import (
	"context"
	"encoding/json"
	"fmt"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/modules_obs_websocket"
	"github.com/satont/twir/libs/types/types/api/modules"
	"google.golang.org/protobuf/types/known/emptypb"
)

var keysForGet = []string{"obs:sources:%s", "obs:audio-sources:%s", "obs:scenes:%s"}

func (c *Modules) ModulesOBSWebsocketGet(
	ctx context.Context,
	_ *emptypb.Empty,
) (*modules_obs_websocket.GetResponse, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.ChannelModulesSettings{}
	if err := c.Db.WithContext(ctx).Where(`"channelId" = ?`, dashboardId).First(entity).Error; err != nil {
		return nil, err
	}

	settings := &modules.OBSWebSocketSettings{}
	if err := json.Unmarshal(entity.Settings, settings); err != nil {
		return nil, err
	}

	results := make([][]string, 3)
	for i, key := range keysForGet {
		val := c.Redis.Get(ctx, fmt.Sprintf(key, dashboardId)).Val()
		if err := json.Unmarshal([]byte(val), &results[i]); err != nil {
			return nil, err
		}
	}

	return &modules_obs_websocket.GetResponse{
		ServerPort:     uint32(settings.ServerPort),
		ServerAddress:  settings.ServerAddress,
		ServerPassword: settings.ServerPassword,
		Sources:        results[0],
		AudioSources:   results[1],
		Scenes:         results[2],
	}, nil
}

func (c *Modules) ModulesOBSWebsocketUpdate(ctx context.Context, request *modules_obs_websocket.PostRequest) (*emptypb.Empty, error) {
	dashboardId := ctx.Value("dashboardId").(string)
	entity := &model.ChannelModulesSettings{}
	if err := c.Db.WithContext(ctx).Where(`"channelId" = ?`, dashboardId).First(entity).Error; err != nil {
		return nil, err
	}

	settings := &modules.OBSWebSocketSettings{
		ServerPort:     int(request.ServerPort),
		ServerAddress:  request.ServerAddress,
		ServerPassword: request.ServerPassword,
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
