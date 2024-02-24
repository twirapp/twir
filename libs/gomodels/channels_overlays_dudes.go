package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
)

type ChannelsOverlaysDudes struct {
	ID                        uuid.UUID       `gorm:"type:uuid;primary_key;"  json:"id"`
	ChannelID                 string          `gorm:"type:text;not null;" json:"channelId"`
	DudeColor                 string          `gorm:"type:text;not null;" json:"dudeColor"`
	DudeEyesColor             string          `gorm:"type:text;not null;" json:"dudeEyesColor"`
	DudeCosmeticsColor        string          `gorm:"type:text;not null;" json:"dudeCosmeticsColor"`
	DudeMaxLifeTime           int32           `gorm:"type:integer;not null;" json:"dudeMaxLifeTime"`
	DudeGravity               int32           `gorm:"type:integer;not null;" json:"dudeGravity"`
	DudeScale                 float32         `gorm:"type:integer;not null;" json:"dudeScale"`
	DudeSoundsEnabled         bool            `gorm:"type:boolean;not null;" json:"dudeSoundsEnabled"`
	DudeSoundsVolume          float32         `gorm:"type:real;not null;"  json:"dudeSoundsVolume"`
	DudeVisibleName           bool            `gorm:"type:boolean;not null;"  json:"dudeVisibleName"`
	DudeGrowTime              int32           `gorm:"type:integer;not null;"  json:"dudeGrowTime"`
	DudeGrowMaxScale          int32           `gorm:"type:integer;not null;"  json:"dudeGrowMaxScale"`
	MessageBoxEnabled         bool            `gorm:"type:boolean;not null;"  json:"messageBoxEnabled"`
	MessageBoxBorderRadius    int32           `gorm:"type:integer;not null;"  json:"messageBoxBorderRadius"`
	MessageBoxBoxColor        string          `gorm:"type:text;not null;"  json:"messageBoxBoxColor"`
	MessageBoxFontFamily      string          `gorm:"type:text;not null;"  json:"messageBoxFontFamily"`
	MessageBoxFontSize        int32           `gorm:"type:integer;not null;"  json:"messageBoxFontSize"`
	MessageBoxPadding         int32           `gorm:"type:integer;not null;"  json:"messageBoxPadding"`
	MessageBoxShowTime        int32           `gorm:"type:integer;not null;"  json:"messageBoxShowTime"`
	MessageBoxFill            string          `gorm:"type:text;not null;"  json:"messageBoxFill"`
	NameBoxFontFamily         string          `gorm:"type:text;not null;" json:"nameBoxFontFamily"`
	NameBoxFontSize           int32           `gorm:"type:integer;not null;" json:"nameBoxFontSize"`
	NameBoxFill               pq.StringArray  `gorm:"type:text[];not null;" json:"nameBoxFill"`
	NameBoxLineJoin           string          `gorm:"type:text;not null;" json:"nameBoxLineJoin"`
	NameBoxStrokeThickness    int32           `gorm:"type:integer;not null;" json:"nameBoxStrokeThickness"`
	NameBoxStroke             string          `gorm:"type:text;not null;" json:"nameBoxStroke"`
	NameBoxFillGradientStops  pq.Float32Array `gorm:"type:text;not null;" json:"nameBoxFillGradientStops"`
	NameBoxFillGradientType   int32           `gorm:"type:integer;not null;" json:"nameBoxFillGradientType"`
	NameBoxFontStyle          string          `gorm:"type:text;not null;" json:"nameBoxFontStyle"`
	NameBoxFontVariant        string          `gorm:"type:text;not null;" json:"nameBoxFontVariant"`
	NameBoxFontWeight         int32           `gorm:"type:integer;not null;" json:"nameBoxFontWeight"`
	NameBoxDropShadow         bool            `gorm:"type:boolean;not null;" json:"nameBoxDropShadow"`
	NameBoxDropShadowAlpha    float32         `gorm:"type:real;not null;" json:"nameBoxDropShadowAlpha"`
	NameBoxDropShadowAngle    float32         `gorm:"type:real;not null;" json:"nameBoxDropShadowAngle"`
	NameBoxDropShadowBlur     float32         `gorm:"type:real;not null;" json:"nameBoxDropShadowBlur"`
	NameBoxDropShadowDistance float32         `gorm:"type:real;not null;" json:"nameBoxDropShadowDistance"`
	NameBoxDropShadowColor    string          `gorm:"type:text;not null;" json:"nameBoxDropShadowColor"`

	IgnoreCommands bool           `gorm:"type:boolean;not null;"  json:"ignoreCommands"`
	IgnoreUsers    bool           `gorm:"type:boolean;not null;"  json:"ignoreUsers"`
	IgnoredUsers   pq.StringArray `gorm:"type:text[];not null;" json:"ignoredUsers"`

	SpitterEmoteEnabled bool `gorm:"type:boolean;not null;" json:"spitterEmoteEnabled"`

	CreatedAt time.Time `gorm:"type:datetime;not null;" json:"createdAt"`
}

func (ChannelsOverlaysDudes) TableName() string {
	return "channels_overlays_dudes"
}
