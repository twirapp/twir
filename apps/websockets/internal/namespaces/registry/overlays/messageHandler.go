package overlays

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/url"
	"strings"

	"github.com/goccy/go-json"
	"github.com/olahol/melody"
	"github.com/samber/lo"
	"github.com/satont/twir/apps/websockets/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/parser"
)

type parseLayerVariablesMessage struct {
	LayerID string `json:"layerId"`
}

type overlayGetLayersMessage struct {
	OverlayID string `json:"overlayId"`
}

type overlayGetLayersResponse struct {
	EventName string                      `json:"eventName"`
	Layers    []model.ChannelOverlayLayer `json:"layers"`
}

func (c *Registry) handleMessage(session *melody.Session, msg []byte) {
	var message types.WebSocketMessage
	if err := json.Unmarshal(msg, &message); err != nil {
		c.services.Logger.Error(err)
		return
	}

	switch message.EventName {
	case "parseLayerVariables":
		var data parseLayerVariablesMessage
		bytes, _ := json.Marshal(message.Data)
		if err := json.Unmarshal(bytes, &data); err != nil {
			c.services.Logger.Error(err)
			return
		}

		var layer model.ChannelOverlayLayer
		if err := c.services.Gorm.
			Preload("Overlay").
			Find(&layer, "id = ?", data.LayerID).
			Error; err != nil {
			c.services.Logger.Error(err)
			return
		}

		if layer.ID.String() == "" || layer.Overlay == nil {
			c.services.Logger.Error("overlay not found")
			return
		}

		text, err := url.QueryUnescape(layer.Settings.HtmlOverlayHTML)
		if err != nil {
			c.services.Logger.Error(err)
			return
		}

		res, err := c.services.Grpc.Parser.ParseTextResponse(
			context.TODO(),
			&parser.ParseTextRequestData{
				Sender: &parser.Sender{},
				Channel: &parser.Channel{
					Id: layer.Overlay.ChannelID,
				},
				Message: &parser.Message{
					Text: text,
				},
				ParseVariables: lo.ToPtr(true),
			},
		)
		if err != nil {
			c.services.Logger.Error(err)
			return
		}

		responses := strings.Join(res.Responses, " ")
		hash := base64.StdEncoding.EncodeToString([]byte(responses))

		if err := session.Write(
			[]byte(fmt.Sprintf(
				`{"eventName":"parsedLayerVariables", "data": "%s", "layerId": "%s"}`,
				hash,
				layer.ID.String(),
			)),
		); err != nil {
			c.services.Logger.Error(err)
		}
	case "getLayers":
		var data overlayGetLayersMessage
		bytes, _ := json.Marshal(message.Data)
		if err := json.Unmarshal(bytes, &data); err != nil {
			c.services.Logger.Error(err)
			return
		}

		socketUserId, ok := session.Get("userId")
		if !ok {
			c.services.Logger.Error("userId not found")
			return
		}

		var overlay model.ChannelOverlay
		if err := c.services.Gorm.
			Preload("Layers").
			Find(&overlay, "channel_id = ? AND id = ?", socketUserId, data.OverlayID).
			Error; err != nil {
			c.services.Logger.Error(err)
			return
		}

		if overlay.ChannelID == "" {
			c.services.Logger.Error("overlay not found")
			return
		}

		responseBytes, err := json.Marshal(
			&overlayGetLayersResponse{
				EventName: "layers",
				Layers:    overlay.Layers,
			},
		)
		if err != nil {
			c.services.Logger.Error(err)
			return
		}

		if err := session.Write(responseBytes); err != nil {
			c.services.Logger.Error(err)
		}
	}
}
