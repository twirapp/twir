package overlays

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/olahol/melody"
	"github.com/satont/twir/apps/websockets/types"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/bus-core/parser"
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

func textToBase64(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

func base64ToText(text string) string {
	bytes, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func (c *Registry) handleMessage(session *melody.Session, msg []byte) {
	var message types.WebSocketMessage
	if err := json.Unmarshal(msg, &message); err != nil {
		c.logger.Error(err.Error())
		return
	}

	switch message.EventName {
	case "parseLayerVariables":
		var data parseLayerVariablesMessage
		bytes, _ := json.Marshal(message.Data)
		if err := json.Unmarshal(bytes, &data); err != nil {
			c.logger.Error(err.Error())
			return
		}

		var layer model.ChannelOverlayLayer
		if err := c.gorm.
			Preload("Overlay").
			Find(&layer, "id = ?", data.LayerID).
			Error; err != nil {
			c.logger.Error(err.Error())
			return
		}

		if layer.ID.String() == "" || layer.Overlay == nil {
			return
		}

		text := base64ToText(layer.Settings.HtmlOverlayHTML)

		res, err := c.bus.Parser.ParseVariablesInText.Request(
			context.Background(),
			parser.ParseVariablesInTextRequest{
				ChannelID: layer.Overlay.ChannelID,
				Text:      text,
			},
		)
		if err != nil {
			c.logger.Error(err.Error())
			return
		}

		if err := session.Write(
			[]byte(fmt.Sprintf(
				`{"eventName":"parsedLayerVariables", "data": "%s", "layerId": "%s"}`,
				textToBase64(res.Data.Text),
				layer.ID.String(),
			)),
		); err != nil {
			c.logger.Error(err.Error())
		}
	case "getLayers":
		var data overlayGetLayersMessage
		bytes, _ := json.Marshal(message.Data)
		if err := json.Unmarshal(bytes, &data); err != nil {
			c.logger.Error(err.Error())
			return
		}

		socketUserId, ok := session.Get("userId")
		if !ok {
			return
		}

		var overlay model.ChannelOverlay
		if err := c.gorm.
			Preload("Layers").
			Find(&overlay, "channel_id = ? AND id = ?", socketUserId, data.OverlayID).
			Error; err != nil {
			c.logger.Error(err.Error())
			return
		}

		if overlay.ChannelID == "" {
			return
		}

		responseBytes, err := json.Marshal(
			&overlayGetLayersResponse{
				EventName: "layers",
				Layers:    overlay.Layers,
			},
		)
		if err != nil {
			c.logger.Error(err.Error())
			return
		}

		if err := session.Write(responseBytes); err != nil {
			c.logger.Error(err.Error())
		}
	}
}
