package overlays_dudes

import (
	"context"

	"github.com/google/uuid"
	"github.com/samber/lo"
	"github.com/twirapp/twir/apps/api-gql/internal/entity"
	"github.com/twirapp/twir/apps/api-gql/internal/wsrouter"
	"github.com/twirapp/twir/libs/repositories/overlays_dudes"
	"github.com/twirapp/twir/libs/repositories/overlays_dudes/model"
	"go.uber.org/fx"
)

type Opts struct {
	fx.In

	OverlaysDudesRepository overlays_dudes.Repository
	WsRouter                wsrouter.WsRouter
}

func New(opts Opts) *Service {
	return &Service{
		overlaysDudesRepository: opts.OverlaysDudesRepository,
		wsRouter:                opts.WsRouter,
	}
}

type Service struct {
	overlaysDudesRepository overlays_dudes.Repository
	wsRouter                wsrouter.WsRouter
}

func (s *Service) modelToEntity(m model.OverlaysDudes) entity.DudesOverlaySettings {
	return entity.DudesOverlaySettings{
		ID: m.ID,
		DudeSettings: entity.DudesDudeSettings{
			Color:          m.DudeColor,
			EyesColor:      m.DudeEyesColor,
			CosmeticsColor: m.DudeCosmeticsColor,
			MaxLifeTime:    m.DudeMaxLifeTime,
			Gravity:        m.DudeGravity,
			Scale:          float64(m.DudeScale),
			SoundsEnabled:  m.DudeSoundsEnabled,
			SoundsVolume:   float64(m.DudeSoundsVolume),
			VisibleName:    m.DudeVisibleName,
			GrowTime:       m.DudeGrowTime,
			GrowMaxScale:   m.DudeGrowMaxScale,
			MaxOnScreen:    m.DudeMaxOnScreen,
			DefaultSprite:  m.DudeDefaultSprite,
		},
		MessageBoxSettings: entity.DudesMessageBoxSettings{
			Enabled:      m.MessageBoxEnabled,
			BorderRadius: m.MessageBoxBorderRadius,
			BoxColor:     m.MessageBoxBoxColor,
			FontFamily:   m.MessageBoxFontFamily,
			FontSize:     m.MessageBoxFontSize,
			Padding:      m.MessageBoxPadding,
			ShowTime:     m.MessageBoxShowTime,
			Fill:         m.MessageBoxFill,
		},
		NameBoxSettings: entity.DudesNameBoxSettings{
			FontFamily:      m.NameBoxFontFamily,
			FontSize:        m.NameBoxFontSize,
			Fill:            m.NameBoxFill,
			LineJoin:        m.NameBoxLineJoin,
			StrokeThickness: m.NameBoxStrokeThickness,
			Stroke:          m.NameBoxStroke,
			FillGradientStops: lo.Map(
				m.NameBoxFillGradientStops,
				func(f float32, _ int) float64 { return float64(f) },
			),
			FillGradientType:   m.NameBoxFillGradientType,
			FontStyle:          m.NameBoxFontStyle,
			FontVariant:        m.NameBoxFontVariant,
			FontWeight:         m.NameBoxFontWeight,
			DropShadow:         m.NameBoxDropShadow,
			DropShadowAlpha:    float64(m.NameBoxDropShadowAlpha),
			DropShadowAngle:    float64(m.NameBoxDropShadowAngle),
			DropShadowBlur:     float64(m.NameBoxDropShadowBlur),
			DropShadowDistance: float64(m.NameBoxDropShadowDistance),
			DropShadowColor:    m.NameBoxDropShadowColor,
		},
		IgnoreSettings: entity.DudesIgnoreSettings{
			IgnoreCommands: m.IgnoreCommands,
			IgnoreUsers:    m.IgnoreUsers,
			Users:          m.IgnoredUsers,
		},
		SpitterEmoteSettings: entity.DudesSpitterEmoteSettings{
			Enabled: m.SpitterEmoteEnabled,
		},
		CreatedAt: m.CreatedAt,
	}
}

type CreateInput struct {
	ChannelID            string
	DudeSettings         CreateDudeSettingsInput
	MessageBoxSettings   CreateMessageBoxSettingsInput
	NameBoxSettings      CreateNameBoxSettingsInput
	IgnoreSettings       CreateIgnoreSettingsInput
	SpitterEmoteSettings CreateSpitterEmoteSettingsInput
}

type CreateDudeSettingsInput struct {
	Color          string
	EyesColor      string
	CosmeticsColor string
	MaxLifeTime    int
	Gravity        int
	Scale          float64
	SoundsEnabled  bool
	SoundsVolume   float64
	VisibleName    bool
	GrowTime       int
	GrowMaxScale   int
	MaxOnScreen    int
	DefaultSprite  string
}

type CreateMessageBoxSettingsInput struct {
	Enabled      bool
	BorderRadius int
	BoxColor     string
	FontFamily   string
	FontSize     int
	Padding      int
	ShowTime     int
	Fill         string
}

type CreateNameBoxSettingsInput struct {
	FontFamily         string
	FontSize           int
	Fill               []string
	LineJoin           string
	StrokeThickness    int
	Stroke             string
	FillGradientStops  []float64
	FillGradientType   int
	FontStyle          string
	FontVariant        string
	FontWeight         int
	DropShadow         bool
	DropShadowAlpha    float64
	DropShadowAngle    float64
	DropShadowBlur     float64
	DropShadowDistance float64
	DropShadowColor    string
}

type CreateIgnoreSettingsInput struct {
	IgnoreCommands bool
	IgnoreUsers    bool
	Users          []string
}

type CreateSpitterEmoteSettingsInput struct {
	Enabled bool
}

func (s *Service) GetByID(ctx context.Context, id uuid.UUID) (entity.DudesOverlaySettings, error) {
	foundModel, err := s.overlaysDudesRepository.GetByID(ctx, id)
	if err != nil {
		return entity.DudesOverlaySettingsNil, err
	}

	return s.modelToEntity(foundModel), nil
}

func (s *Service) GetManyByChannelID(
	ctx context.Context,
	channelID string,
) ([]entity.DudesOverlaySettings, error) {
	models, err := s.overlaysDudesRepository.GetManyByChannelID(ctx, channelID)
	if err != nil {
		return nil, err
	}

	entities := make([]entity.DudesOverlaySettings, len(models))
	for i, o := range models {
		entities[i] = s.modelToEntity(o)
	}

	return entities, nil
}

func (s *Service) Create(ctx context.Context, input CreateInput) (
	entity.DudesOverlaySettings,
	error,
) {
	foundModel, err := s.overlaysDudesRepository.Create(
		ctx, overlays_dudes.CreateInput{
			ChannelID:          input.ChannelID,
			DudeColor:          input.DudeSettings.Color,
			DudeEyesColor:      input.DudeSettings.EyesColor,
			DudeCosmeticsColor: input.DudeSettings.CosmeticsColor,
			DudeMaxLifeTime:    input.DudeSettings.MaxLifeTime,
			DudeGravity:        input.DudeSettings.Gravity,
			DudeScale:          float32(input.DudeSettings.Scale),
			DudeSoundsEnabled:  input.DudeSettings.SoundsEnabled,
			DudeSoundsVolume:   float32(input.DudeSettings.SoundsVolume),
			DudeVisibleName:    input.DudeSettings.VisibleName,
			DudeGrowTime:       input.DudeSettings.GrowTime,
			DudeGrowMaxScale:   input.DudeSettings.GrowMaxScale,
			DudeMaxOnScreen:    input.DudeSettings.MaxOnScreen,
			DudeDefaultSprite:  input.DudeSettings.DefaultSprite,

			MessageBoxEnabled:      input.MessageBoxSettings.Enabled,
			MessageBoxBorderRadius: input.MessageBoxSettings.BorderRadius,
			MessageBoxBoxColor:     input.MessageBoxSettings.BoxColor,
			MessageBoxFontFamily:   input.MessageBoxSettings.FontFamily,
			MessageBoxFontSize:     input.MessageBoxSettings.FontSize,
			MessageBoxPadding:      input.MessageBoxSettings.Padding,
			MessageBoxShowTime:     input.MessageBoxSettings.ShowTime,
			MessageBoxFill:         input.MessageBoxSettings.Fill,

			NameBoxFontFamily:      input.NameBoxSettings.FontFamily,
			NameBoxFontSize:        input.NameBoxSettings.FontSize,
			NameBoxFill:            input.NameBoxSettings.Fill,
			NameBoxLineJoin:        input.NameBoxSettings.LineJoin,
			NameBoxStrokeThickness: input.NameBoxSettings.StrokeThickness,
			NameBoxStroke:          input.NameBoxSettings.Stroke,
			NameBoxFillGradientStops: lo.Map(
				input.NameBoxSettings.FillGradientStops,
				func(f float64, _ int) float32 { return float32(f) },
			),
			NameBoxFillGradientType:   input.NameBoxSettings.FillGradientType,
			NameBoxFontStyle:          input.NameBoxSettings.FontStyle,
			NameBoxFontVariant:        input.NameBoxSettings.FontVariant,
			NameBoxFontWeight:         input.NameBoxSettings.FontWeight,
			NameBoxDropShadow:         input.NameBoxSettings.DropShadow,
			NameBoxDropShadowAlpha:    float32(input.NameBoxSettings.DropShadowAlpha),
			NameBoxDropShadowAngle:    float32(input.NameBoxSettings.DropShadowAngle),
			NameBoxDropShadowBlur:     float32(input.NameBoxSettings.DropShadowBlur),
			NameBoxDropShadowDistance: float32(input.NameBoxSettings.DropShadowDistance),
			NameBoxDropShadowColor:    input.NameBoxSettings.DropShadowColor,

			IgnoreCommands: input.IgnoreSettings.IgnoreCommands,
			IgnoreUsers:    input.IgnoreSettings.IgnoreUsers,
			IgnoredUsers:   input.IgnoreSettings.Users,

			SpitterEmoteEnabled: input.SpitterEmoteSettings.Enabled,
		},
	)
	if err != nil {
		return entity.DudesOverlaySettingsNil, err
	}

	return s.modelToEntity(foundModel), nil
}

type UpdateInput struct {
	DudeSettings         *UpdateDudeSettingsInput
	MessageBoxSettings   *UpdateMessageBoxSettingsInput
	NameBoxSettings      *UpdateNameBoxSettingsInput
	IgnoreSettings       *UpdateIgnoreSettingsInput
	SpitterEmoteSettings *UpdateSpitterEmoteSettingsInput
}

type UpdateDudeSettingsInput struct {
	Color          *string
	EyesColor      *string
	CosmeticsColor *string
	MaxLifeTime    *int
	Gravity        *int
	Scale          *float64
	SoundsEnabled  *bool
	SoundsVolume   *float64
	VisibleName    *bool
	GrowTime       *int
	GrowMaxScale   *int
	MaxOnScreen    *int
	DefaultSprite  *string
}

type UpdateMessageBoxSettingsInput struct {
	Enabled      *bool
	BorderRadius *int
	BoxColor     *string
	FontFamily   *string
	FontSize     *int
	Padding      *int
	ShowTime     *int
	Fill         *string
}

type UpdateNameBoxSettingsInput struct {
	FontFamily         *string
	FontSize           *int
	Fill               *[]string
	LineJoin           *string
	StrokeThickness    *int
	Stroke             *string
	FillGradientStops  *[]float64
	FillGradientType   *int
	FontStyle          *string
	FontVariant        *string
	FontWeight         *int
	DropShadow         *bool
	DropShadowAlpha    *float64
	DropShadowAngle    *float64
	DropShadowBlur     *float64
	DropShadowDistance *float64
	DropShadowColor    *string
}

type UpdateIgnoreSettingsInput struct {
	IgnoreCommands *bool
	IgnoreUsers    *bool
	Users          *[]string
}

type UpdateSpitterEmoteSettingsInput struct {
	Enabled *bool
}

func (s *Service) Update(
	ctx context.Context,
	id uuid.UUID,
	input UpdateInput,
) (entity.DudesOverlaySettings, error) {
	updateInput := overlays_dudes.UpdateInput{}

	if input.DudeSettings != nil {
		updateInput.DudeColor = input.DudeSettings.Color
		updateInput.DudeEyesColor = input.DudeSettings.EyesColor
		updateInput.DudeCosmeticsColor = input.DudeSettings.CosmeticsColor
		updateInput.DudeMaxLifeTime = input.DudeSettings.MaxLifeTime
		updateInput.DudeGravity = input.DudeSettings.Gravity
		if input.DudeSettings.Scale != nil {
			updateInput.DudeScale = lo.ToPtr(float32(*input.DudeSettings.Scale))
		}
		updateInput.DudeSoundsEnabled = input.DudeSettings.SoundsEnabled
		if input.DudeSettings.SoundsVolume != nil {
			updateInput.DudeSoundsVolume = lo.ToPtr(float32(*input.DudeSettings.SoundsVolume))
		}
		updateInput.DudeVisibleName = input.DudeSettings.VisibleName
		updateInput.DudeGrowTime = input.DudeSettings.GrowTime
		updateInput.DudeGrowMaxScale = input.DudeSettings.GrowMaxScale
		updateInput.DudeMaxOnScreen = input.DudeSettings.MaxOnScreen
		updateInput.DudeDefaultSprite = input.DudeSettings.DefaultSprite
	}

	if input.MessageBoxSettings != nil {
		updateInput.MessageBoxEnabled = input.MessageBoxSettings.Enabled
		updateInput.MessageBoxBorderRadius = input.MessageBoxSettings.BorderRadius
		updateInput.MessageBoxBoxColor = input.MessageBoxSettings.BoxColor
		updateInput.MessageBoxFontFamily = input.MessageBoxSettings.FontFamily
		updateInput.MessageBoxFontSize = input.MessageBoxSettings.FontSize
		updateInput.MessageBoxPadding = input.MessageBoxSettings.Padding
		updateInput.MessageBoxShowTime = input.MessageBoxSettings.ShowTime
		updateInput.MessageBoxFill = input.MessageBoxSettings.Fill
	}

	if input.NameBoxSettings != nil {
		updateInput.NameBoxFontFamily = input.NameBoxSettings.FontFamily
		updateInput.NameBoxFontSize = input.NameBoxSettings.FontSize
		updateInput.NameBoxFill = input.NameBoxSettings.Fill
		updateInput.NameBoxLineJoin = input.NameBoxSettings.LineJoin
		updateInput.NameBoxStrokeThickness = input.NameBoxSettings.StrokeThickness
		updateInput.NameBoxStroke = input.NameBoxSettings.Stroke
		if input.NameBoxSettings.FillGradientStops != nil {
			updateInput.NameBoxFillGradientStops = lo.ToPtr(
				lo.Map(
					*input.NameBoxSettings.FillGradientStops,
					func(f float64, _ int) float32 { return float32(f) },
				),
			)
		}
		updateInput.NameBoxFillGradientType = input.NameBoxSettings.FillGradientType
		updateInput.NameBoxFontStyle = input.NameBoxSettings.FontStyle
		updateInput.NameBoxFontVariant = input.NameBoxSettings.FontVariant
		updateInput.NameBoxFontWeight = input.NameBoxSettings.FontWeight
		updateInput.NameBoxDropShadow = input.NameBoxSettings.DropShadow
		if input.NameBoxSettings.DropShadowAlpha != nil {
			updateInput.NameBoxDropShadowAlpha = lo.ToPtr(float32(*input.NameBoxSettings.DropShadowAlpha))
		}
		if input.NameBoxSettings.DropShadowAngle != nil {
			updateInput.NameBoxDropShadowAngle = lo.ToPtr(float32(*input.NameBoxSettings.DropShadowAngle))
		}
		if input.NameBoxSettings.DropShadowBlur != nil {
			updateInput.NameBoxDropShadowBlur = lo.ToPtr(float32(*input.NameBoxSettings.DropShadowBlur))
		}
		if input.NameBoxSettings.DropShadowDistance != nil {
			updateInput.NameBoxDropShadowDistance = lo.ToPtr(float32(*input.NameBoxSettings.DropShadowDistance))
		}
		updateInput.NameBoxDropShadowColor = input.NameBoxSettings.DropShadowColor
	}

	if input.IgnoreSettings != nil {
		updateInput.IgnoreCommands = input.IgnoreSettings.IgnoreCommands
		updateInput.IgnoreUsers = input.IgnoreSettings.IgnoreUsers
		updateInput.IgnoredUsers = input.IgnoreSettings.Users
	}

	if input.SpitterEmoteSettings != nil {
		updateInput.SpitterEmoteEnabled = input.SpitterEmoteSettings.Enabled
	}

	updatedModel, err := s.overlaysDudesRepository.Update(ctx, id, updateInput)
	if err != nil {
		return entity.DudesOverlaySettingsNil, err
	}

	if err := s.wsRouter.Publish(
		CreateDudesWsRouterKey(updatedModel.ChannelID, id),
		s.modelToEntity(updatedModel),
	); err != nil {
		return entity.DudesOverlaySettings{}, err
	}

	return s.modelToEntity(updatedModel), nil
}

func (s *Service) Delete(ctx context.Context, id uuid.UUID) error {
	return s.overlaysDudesRepository.Delete(ctx, id)
}
