package overlays_dudes

import (
	"context"

	"github.com/google/uuid"
	"github.com/twirapp/twir/libs/repositories/overlays_dudes/model"
)

type Repository interface {
	GetByID(ctx context.Context, id uuid.UUID) (model.OverlaysDudes, error)
	GetManyByChannelID(ctx context.Context, channelID string) ([]model.OverlaysDudes, error)
	Create(ctx context.Context, input CreateInput) (model.OverlaysDudes, error)
	Update(ctx context.Context, id uuid.UUID, input UpdateInput) (model.OverlaysDudes, error)
	Delete(ctx context.Context, id uuid.UUID) error
}

type CreateInput struct {
	ChannelID string

	// Dude settings
	DudeColor          string
	DudeEyesColor      string
	DudeCosmeticsColor string
	DudeMaxLifeTime    int
	DudeGravity        int
	DudeScale          float32
	DudeSoundsEnabled  bool
	DudeSoundsVolume   float32
	DudeVisibleName    bool
	DudeGrowTime       int
	DudeGrowMaxScale   int
	DudeMaxOnScreen    int
	DudeDefaultSprite  string

	// Message box settings
	MessageBoxEnabled      bool
	MessageBoxBorderRadius int
	MessageBoxBoxColor     string
	MessageBoxFontFamily   string
	MessageBoxFontSize     int
	MessageBoxPadding      int
	MessageBoxShowTime     int
	MessageBoxFill         string

	// Name box settings
	NameBoxFontFamily         string
	NameBoxFontSize           int
	NameBoxFill               []string
	NameBoxLineJoin           string
	NameBoxStrokeThickness    int
	NameBoxStroke             string
	NameBoxFillGradientStops  []float32
	NameBoxFillGradientType   int
	NameBoxFontStyle          string
	NameBoxFontVariant        string
	NameBoxFontWeight         int
	NameBoxDropShadow         bool
	NameBoxDropShadowAlpha    float32
	NameBoxDropShadowAngle    float32
	NameBoxDropShadowBlur     float32
	NameBoxDropShadowDistance float32
	NameBoxDropShadowColor    string

	// Ignore settings
	IgnoreCommands bool
	IgnoreUsers    bool
	IgnoredUsers   []string

	// Spitter emote settings
	SpitterEmoteEnabled bool
}

type UpdateInput struct {
	// Dude settings
	DudeColor          *string
	DudeEyesColor      *string
	DudeCosmeticsColor *string
	DudeMaxLifeTime    *int
	DudeGravity        *int
	DudeScale          *float32
	DudeSoundsEnabled  *bool
	DudeSoundsVolume   *float32
	DudeVisibleName    *bool
	DudeGrowTime       *int
	DudeGrowMaxScale   *int
	DudeMaxOnScreen    *int
	DudeDefaultSprite  *string

	// Message box settings
	MessageBoxEnabled      *bool
	MessageBoxBorderRadius *int
	MessageBoxBoxColor     *string
	MessageBoxFontFamily   *string
	MessageBoxFontSize     *int
	MessageBoxPadding      *int
	MessageBoxShowTime     *int
	MessageBoxFill         *string

	// Name box settings
	NameBoxFontFamily         *string
	NameBoxFontSize           *int
	NameBoxFill               *[]string
	NameBoxLineJoin           *string
	NameBoxStrokeThickness    *int
	NameBoxStroke             *string
	NameBoxFillGradientStops  *[]float32
	NameBoxFillGradientType   *int
	NameBoxFontStyle          *string
	NameBoxFontVariant        *string
	NameBoxFontWeight         *int
	NameBoxDropShadow         *bool
	NameBoxDropShadowAlpha    *float32
	NameBoxDropShadowAngle    *float32
	NameBoxDropShadowBlur     *float32
	NameBoxDropShadowDistance *float32
	NameBoxDropShadowColor    *string

	// Ignore settings
	IgnoreCommands *bool
	IgnoreUsers    *bool
	IgnoredUsers   *[]string

	// Spitter emote settings
	SpitterEmoteEnabled *bool
}
