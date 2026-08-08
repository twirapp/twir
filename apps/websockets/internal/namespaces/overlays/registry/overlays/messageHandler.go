package overlays

import (
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"log/slog"
	"math"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/olahol/melody"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/websockets/types"
	"github.com/twirapp/twir/libs/bus-core/parser"
	customoverlayentity "github.com/twirapp/twir/libs/entities/custom_overlay"
	"github.com/twirapp/twir/libs/entities/platform"
	model "github.com/twirapp/twir/libs/gomodels"
	twirlogger "github.com/twirapp/twir/libs/logger"
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

// overlayEditorEventMeta models only the routing fields of editor collaboration
// events; the rest of the payload is forwarded to peers opaquely.
type overlayEditorEventMeta struct {
	OverlayID string `json:"overlayId"`
	ClientID  string `json:"clientId"`
}

// Geometry fields are float64 because browsers can send fractional pixels
// (center snapping, align, distribute); they are rounded before persisting.
type instantSaveLayerData struct {
	ID       string   `json:"id"`
	PosX     float64  `json:"posX"`
	PosY     float64  `json:"posY"`
	Rotation float64  `json:"rotation"`
	Width    float64  `json:"width"`
	Height   float64  `json:"height"`
	Visible  bool     `json:"visible"`
	Opacity  float64  `json:"opacity"`
	ZIndex   *float64 `json:"zIndex"`
}

func roundToInt(v float64) int {
	return int(math.Round(v))
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

		channelID, err := uuid.Parse(layer.Overlay.ChannelID)
		if err != nil {
			c.logger.Error(err.Error())
			return
		}
		platformSource := platform.PlatformTwitch
		res, err := c.bus.Parser.ParseVariablesInText.Request(
			context.Background(),
			parser.ParseVariablesInTextRequest{
				ChannelID:      channelID,
				Text:           text,
				PlatformSource: &platformSource,
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

			zIndex := layer.ZIndex
			if foundInputLayer.ZIndex != nil {
				zIndex = roundToInt(*foundInputLayer.ZIndex)
			}

			e.Layers = append(e.Layers, customoverlayentity.ChannelOverlayLayer{
				ID:        layer.ID,
				OverlayID: layer.OverlayID,
				Type:      customoverlayentity.ChannelOverlayType(layer.Type),
				PosX:      roundToInt(foundInputLayer.PosX),
				PosY:      roundToInt(foundInputLayer.PosY),
				Rotation:  roundToInt(foundInputLayer.Rotation),
				Settings: customoverlayentity.ChannelOverlayLayerSettings{
					HtmlOverlayHTML:                    layer.Settings.HtmlOverlayHTML,
					HtmlOverlayCSS:                     layer.Settings.HtmlOverlayCSS,
					HtmlOverlayJS:                      layer.Settings.HtmlOverlayJS,
					HtmlOverlayDataPollSecondsInterval: layer.Settings.HtmlOverlayDataPollSecondsInterval,
					ImageUrl:                           layer.Settings.ImageUrl,
					TextContent:                        layer.Settings.TextContent,
					TextFontFamily:                     layer.Settings.TextFontFamily,
					TextFontSize:                       layer.Settings.TextFontSize,
					TextFontWeight:                     layer.Settings.TextFontWeight,
					TextFontStyle:                      layer.Settings.TextFontStyle,
					TextColor:                          layer.Settings.TextColor,
					TextAlign:                          layer.Settings.TextAlign,
					TextAlignVertical:                  layer.Settings.TextAlignVertical,
					TextStrokeWidth:                    layer.Settings.TextStrokeWidth,
					TextStrokeColor:                    layer.Settings.TextStrokeColor,
					TextShadowColor:                    layer.Settings.TextShadowColor,
					TextShadowBlur:                     layer.Settings.TextShadowBlur,
					TextShadowOffsetX:                  layer.Settings.TextShadowOffsetX,
					TextShadowOffsetY:                  layer.Settings.TextShadowOffsetY,
					TextLineHeight:                     layer.Settings.TextLineHeight,
					TextLetterSpacing:                  layer.Settings.TextLetterSpacing,
					TextTransform:                      layer.Settings.TextTransform,
					VideoUrl:                           layer.Settings.VideoUrl,
					VideoLoop:                          layer.Settings.VideoLoop,
					VideoMuted:                         layer.Settings.VideoMuted,
					IframeUrl:                          layer.Settings.IframeUrl,
					IframeScale:                        layer.Settings.IframeScale,
					WidgetKey:                          layer.Settings.WidgetKey,
					YoutubeVideoID:                     layer.Settings.YoutubeVideoID,
					YoutubeAutoplay:                    layer.Settings.YoutubeAutoplay,
					YoutubeLoop:                        layer.Settings.YoutubeLoop,
					YoutubeMuted:                       layer.Settings.YoutubeMuted,
					EmoteUrl:                           layer.Settings.EmoteUrl,
					EmoteName:                          layer.Settings.EmoteName,
					EmoteProvider:                      layer.Settings.EmoteProvider,
				},
				CreatedAt:               layer.CreatedAt,
				UpdatedAt:               layer.UpdatedAt,
				Width:                   roundToInt(foundInputLayer.Width),
				Height:                  roundToInt(foundInputLayer.Height),
				PeriodicallyRefetchData: layer.PeriodicallyRefetchData,
				Locked:                  layer.Locked,
				Visible:                 foundInputLayer.Visible,
				Opacity:                 foundInputLayer.Opacity,
				ZIndex:                  zIndex,
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

			posX := roundToInt(layerData.PosX)
			posY := roundToInt(layerData.PosY)
			rotation := roundToInt(layerData.Rotation)
			width := roundToInt(layerData.Width)
			height := roundToInt(layerData.Height)
			var zIndex *int
			if layerData.ZIndex != nil {
				zIndex = lo.ToPtr(roundToInt(*layerData.ZIndex))
			}

			go func() {
				_, e := c.channelsOverlaysRepository.UpdateLayer(context.TODO(), layerID, channels_overlays.LayerUpdateInput{
					PosX:     &posX,
					PosY:     &posY,
					Rotation: &rotation,
					Width:    &width,
					Height:   &height,
					Visible:  &layerData.Visible,
					Opacity:  &layerData.Opacity,
					ZIndex:   zIndex,
				})
				if e != nil {
					// Layers created on the client exist only in the builder state until
					// the overlay is fully saved, so a missing row here is expected.
					if errors.Is(e, channels_overlays.ErrNotFound) {
						c.logger.Debug(
							"skipping instant save for not yet persisted layer",
							twirlogger.Error(e),
							slog.String("layer_id", layerID.String()),
						)
					} else {
						c.logger.Error("failed to update layer", twirlogger.Error(e))
					}
				}
			}()
		}

		// Send acknowledgment
		if err := session.Write(
			[]byte(`{"eventName":"instantSaveAck"}`),
		); err != nil {
			c.logger.Error(err.Error())
		}
	case
		"overlayEditorLayerAdd",
		"overlayEditorLayerRemove",
		"overlayEditorLayerUpdate",
		"overlayEditorLayerPositions",
		"overlayEditorLayersReorder",
		"overlayEditorSettingsUpdate",
		"overlayEditorProjectReplace",
		"overlayEditorSyncRequest",
		"overlayEditorSyncState":
		c.handleOverlayEditorEvent(session, message)
	}
}

// handleOverlayEditorEvent rebroadcasts overlay editor collaboration events to
// every session of the same channel (including the sender, which filters its own
// echo by clientId), so multiple browser tabs editing one overlay stay in sync.
func (c *Registry) handleOverlayEditorEvent(session *melody.Session, message types.WebSocketMessage) {
	var meta overlayEditorEventMeta
	bytes, err := json.Marshal(message.Data)
	if err != nil {
		c.logger.Error("failed to marshal overlay editor event data", twirlogger.Error(err))
		return
	}
	if err := json.Unmarshal(bytes, &meta); err != nil {
		c.logger.Error("failed to unmarshal overlay editor event meta", twirlogger.Error(err))
		return
	}

	socketUserId, ok := session.Get("userId")
	if !ok {
		return
	}
	channelID, ok := socketUserId.(string)
	if !ok {
		return
	}

	cacheKey := fmt.Sprintf("overlayEditorAccess:%s", meta.OverlayID)
	if _, ok := session.Get(cacheKey); !ok {
		overlayID, err := uuid.Parse(meta.OverlayID)
		if err != nil {
			c.logger.Error("invalid overlay ID in editor event", twirlogger.Error(err))
			return
		}

		overlay, err := c.channelsOverlaysRepository.GetByID(context.Background(), overlayID)
		if err != nil {
			c.logger.Error("failed to get overlay for editor event", twirlogger.Error(err))
			return
		}

		if overlay.ChannelID != channelID {
			c.logger.Error(
				"overlay editor event rejected: overlay does not belong to user",
				slog.String("user_id", channelID),
				slog.String("overlay_id", meta.OverlayID),
			)
			return
		}

		session.Set(cacheKey, struct{}{})
	}

	if err := c.SendEvent(channelID, message.EventName, message.Data); err != nil {
		c.logger.Error("failed to broadcast overlay editor event", twirlogger.Error(err))
	}
}
