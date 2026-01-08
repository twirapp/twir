package overlays

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/olahol/melody"
	"github.com/twirapp/twir/apps/websockets/types"
	"github.com/twirapp/twir/libs/bus-core/parser"
	customoverlayentity "github.com/twirapp/twir/libs/entities/custom_overlay"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/repositories/channels_overlays"
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

type instantSaveLayerMessage struct {
	OverlayID string                 `json:"overlayId"`
	Layers    []instantSaveLayerData `json:"layers"`
}

type instantSaveLayerData struct {
	ID       string  `json:"id"`
	PosX     int     `json:"posX"`
	PosY     int     `json:"posY"`
	Rotation int     `json:"rotation"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
	Visible  bool    `json:"visible"`
	Opacity  float64 `json:"opacity"`
}

func textToBase64(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

func base64ToText(text string) (string, bool) {
	bytes, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return "", false
	}
	return string(bytes), true
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

		// Handle both old base64-encoded data and new plain text data
		// Try to decode as base64. If successful, use decoded text.
		// Otherwise, use the original text (new GraphQL format stores plain text).
		text := layer.Settings.HtmlOverlayHTML
		if decodedText, ok := base64ToText(text); ok {
			text = decodedText
		}

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
	case "instantSaveLayerPositions":
		var data instantSaveLayerMessage
		bytes, _ := json.Marshal(message.Data)
		if err := json.Unmarshal(bytes, &data); err != nil {
			c.logger.Error(err.Error())
			return
		}

		socketUserId, ok := session.Get("userId")
		if !ok {
			return
		}

		overlayIdParsed, err := uuid.Parse(data.OverlayID)
		if err != nil {
			c.logger.Error("invalid overlay ID", "error", err)
			return
		}

		// Verify overlay belongs to user
		overlay, err := c.channelsOverlaysRepository.GetByID(
			context.Background(),
			overlayIdParsed,
		)
		if err != nil {
			c.logger.Error("failed to get overlay", "error", err)
			return
		}

		if overlay.ChannelID != socketUserId {
			c.logger.Error("overlay does not belong to user", "userId", socketUserId, "overlayId", data.OverlayID)
			return
		}

		// Publish update via wsRouter to notify GraphQL subscribers
		wsRouterKey := fmt.Sprintf("api.customOverlaySettings.%s.%s", overlay.ChannelID, overlay.ID)

		e := customoverlayentity.ChannelOverlay{
			ID:        overlay.ID,
			ChannelID: overlay.ChannelID,
			Name:      overlay.Name,
			CreatedAt: overlay.CreatedAt,
			UpdatedAt: overlay.UpdatedAt,
			Width:     overlay.Width,
			Height:    overlay.Height,
			InstaSave: overlay.InstaSave,
			Layers:    []customoverlayentity.ChannelOverlayLayer{},
		}

		for _, layer := range overlay.Layers {
			var foundInputLayer *instantSaveLayerData
			for _, inputLayer := range data.Layers {
				if inputLayer.ID == layer.ID.String() {
					foundInputLayer = &inputLayer
					break
				}
			}

			if foundInputLayer == nil {
				continue
			}

			e.Layers = append(e.Layers, customoverlayentity.ChannelOverlayLayer{
				ID:        layer.ID,
				OverlayID: layer.OverlayID,
				Type:      customoverlayentity.ChannelOverlayType(layer.Type),
				PosX:      foundInputLayer.PosX,
				PosY:      foundInputLayer.PosY,
				Rotation:  foundInputLayer.Rotation,
				Settings: customoverlayentity.ChannelOverlayLayerSettings{
					HtmlOverlayHTML:                    layer.Settings.HtmlOverlayHTML,
					HtmlOverlayCSS:                     layer.Settings.HtmlOverlayCSS,
					HtmlOverlayJS:                      layer.Settings.HtmlOverlayJS,
					HtmlOverlayDataPollSecondsInterval: layer.Settings.HtmlOverlayDataPollSecondsInterval,
					ImageUrl:                           layer.Settings.ImageUrl,
				},
				CreatedAt:               layer.CreatedAt,
				UpdatedAt:               layer.UpdatedAt,
				Width:                   foundInputLayer.Width,
				Height:                  foundInputLayer.Height,
				PeriodicallyRefetchData: layer.PeriodicallyRefetchData,
				Locked:                  layer.Locked,
				Visible:                 foundInputLayer.Visible,
				Opacity:                 foundInputLayer.Opacity,
			})
		}

		if err := c.wsRouter.Publish(wsRouterKey, e); err != nil {
			c.logger.Error("failed to publish overlay update", "error", err)
		}

		for _, layerData := range data.Layers {
			layerID, err := uuid.Parse(layerData.ID)
			if err != nil {
				c.logger.Error("invalid layer ID", "error", err)
				continue
			}

			go func() {
				_, e := c.channelsOverlaysRepository.UpdateLayer(context.TODO(), layerID, channels_overlays.LayerUpdateInput{
					PosX:     &layerData.PosX,
					PosY:     &layerData.PosY,
					Rotation: &layerData.Rotation,
					Width:    &layerData.Width,
					Height:   &layerData.Height,
					Visible:  &layerData.Visible,
					Opacity:  &layerData.Opacity,
				})
				if e != nil {
					c.logger.Error("failed to update layer", "error", e)
				}
			}()
		}

		// Send acknowledgment
		if err := session.Write(
			[]byte(`{"eventName":"instantSaveAck"}`),
		); err != nil {
			c.logger.Error(err.Error())
		}
	}
}
