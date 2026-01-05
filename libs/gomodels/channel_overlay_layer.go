package model

import (
	"database/sql/driver"
	"errors"
	"time"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type ChannelOverlayLayer struct {
	ID                      uuid.UUID                   `gorm:"primary_key;column:id;type:UUID;"  json:"id"`
	Type                    ChannelOverlayType          `gorm:"column:type;type:TEXT;"  json:"type"`
	Settings                ChannelOverlayLayerSettings `gorm:"column:settings;type:JSONB;"  json:"settings"`
	OverlayID               uuid.UUID                   `gorm:"column:overlay_id;type:UUID;"  json:"overlayId"`
	PosX                    int                         `gorm:"column:pos_x;type:INTEGER;"  json:"pos_x"`
	PosY                    int                         `gorm:"column:pos_y;type:INTEGER;"  json:"pos_y"`
	Width                   int                         `gorm:"column:width;type:INTEGER;"  json:"width"`
	Rotation                int                         `gorm:"column:rotation;type:INTEGER;default:0;"  json:"rotation"`
	Height                  int                         `gorm:"column:height;type:INTEGER;"  json:"height"`
	CreatedAt               time.Time                   `gorm:"column:created_at;data:timestamp;"  json:"createdAt"`
	UpdatedAt               time.Time                   `gorm:"column:updated_at;data:timestamp;"  json:"updatedAt"`
	PeriodicallyRefetchData bool                        `gorm:"column:periodically_refetch_data;type:BOOLEAN"  json:"periodically_refetch_data"`

	Overlay *ChannelOverlay `gorm:"foreignKey:OverlayID" json:"overlay"`
}

func (c ChannelOverlayLayer) TableName() string {
	return "channels_overlays_layers"
}

// ChannelOverlayType types
type ChannelOverlayType string

func (e ChannelOverlayType) String() string {
	return string(e)
}

const (
	ChannelOverlayTypeHTML ChannelOverlayType = "HTML"
)

// ChannelOverlayLayerSettings settings
type ChannelOverlayLayerSettings struct {
	HtmlOverlayHTML                    string `json:"htmlOverlayHtml,omitempty"`
	HtmlOverlayCSS                     string `json:"htmlOverlayCss,omitempty"`
	HtmlOverlayJS                      string `json:"htmlOverlayJs,omitempty"`
	HtmlOverlayDataPollSecondsInterval int    `json:"htmlOverlayDataPollSecondsInterval,omitempty"`
}

func (a ChannelOverlayLayerSettings) Value() (driver.Value, error) {
	b, err := json.Marshal(a)
	if err != nil {
		return nil, err
	}

	return string(b), nil
}

func (a *ChannelOverlayLayerSettings) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &a)
}
