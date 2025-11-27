package overlays

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/twirapp/twir/apps/api/internal/helpers"
	"github.com/twirapp/twir/apps/api/internal/impl_deps"
	"github.com/twirapp/twir/libs/api/messages/overlays"
	"github.com/twirapp/twir/libs/bus-core/parser"
	model "github.com/twirapp/twir/libs/gomodels"
	"github.com/twirapp/twir/libs/grpc/websockets"
	"github.com/twirapp/twir/libs/logger"
	"github.com/twitchtv/twirp"
	"google.golang.org/protobuf/types/known/emptypb"
	"gorm.io/gorm"
)

type Overlays struct {
	*impl_deps.Deps
}

func (c *Overlays) convertToRpcType(t model.ChannelOverlayType) overlays.OverlayLayerType {
	switch t {
	case model.ChannelOverlayTypeHTML:
		return overlays.OverlayLayerType_HTML
	default:
		return overlays.OverlayLayerType_HTML
	}
}

func (c *Overlays) convertToDbType(t overlays.OverlayLayerType) model.ChannelOverlayType {
	switch t {
	case overlays.OverlayLayerType_HTML:
		return model.ChannelOverlayTypeHTML
	default:
		return model.ChannelOverlayTypeHTML
	}
}

func textToBase64(text string) string {
	return base64.StdEncoding.EncodeToString([]byte(text))
}

func (c *Overlays) validateReq(
	name string,
	width, height int32,
	layers []*overlays.CreateLayer,
) error {
	if utf8.RuneCountInString(name) > 30 {
		return twirp.NewError(twirp.InvalidArgument, "name is too long")
	}

	if width < 0 {
		return twirp.NewError(twirp.InvalidArgument, "width must be positive")
	}

	if height < 0 {
		return twirp.NewError(twirp.InvalidArgument, "height must be positive")
	}

	if len(layers) > 15 {
		return twirp.NewError(twirp.InvalidArgument, "too many layers")
	}

	return nil
}

func base64ToText(text string) string {
	bytes, err := base64.StdEncoding.DecodeString(text)
	if err != nil {
		return ""
	}
	return string(bytes)
}

func (c *Overlays) convertEntity(entity model.ChannelOverlay) *overlays.Overlay {
	id := entity.ID.String()

	layers := make([]*overlays.OverlayLayer, len(entity.Layers))
	for i, l := range entity.Layers {
		layers[i] = &overlays.OverlayLayer{
			Id:   l.ID.String(),
			Type: c.convertToRpcType(l.Type),
			Settings: &overlays.OverlayLayerSettings{
				HtmlOverlayHtml:                        base64ToText(l.Settings.HtmlOverlayHTML),
				HtmlOverlayCss:                         base64ToText(l.Settings.HtmlOverlayCSS),
				HtmlOverlayJs:                          base64ToText(l.Settings.HtmlOverlayJS),
				HtmlOverlayHtmlDataPollSecondsInterval: int32(l.Settings.HtmlOverlayDataPollSecondsInterval),
			},
			OverlayId:               id,
			PosX:                    int32(l.PosX),
			PosY:                    int32(l.PosY),
			Width:                   int32(l.Width),
			Height:                  int32(l.Height),
			CreatedAt:               fmt.Sprint(l.CreatedAt.UnixMilli()),
			UpdatedAt:               fmt.Sprint(l.UpdatedAt.UnixMilli()),
			PeriodicallyRefetchData: l.PeriodicallyRefetchData,
		}
	}

	return &overlays.Overlay{
		Id:        entity.ID.String(),
		ChannelId: entity.ChannelID,
		Name:      entity.Name,
		CreatedAt: fmt.Sprint(entity.CreatedAt.UnixMilli()),
		UpdatedAt: fmt.Sprint(entity.UpdatedAt.UnixMilli()),
		Layers:    layers,
		Width:     int32(entity.Width),
		Height:    int32(entity.Height),
	}
}

func (c *Overlays) OverlaysGetAll(ctx context.Context, _ *emptypb.Empty) (
	*overlays.GetAllResponse,
	error,
) {
	dashboardId := ctx.Value("dashboardId").(string)

	var entities []model.ChannelOverlay
	if err := c.Db.
		WithContext(ctx).
		Preload("Layers").
		Order("created_at DESC").
		Find(&entities, "channel_id = ?", dashboardId).
		Error; err != nil {
		return nil, err
	}

	res := make([]*overlays.Overlay, len(entities))
	for i, e := range entities {
		res[i] = c.convertEntity(e)
	}

	return &overlays.GetAllResponse{
		Overlays: res,
	}, nil
}

func (c *Overlays) OverlaysGetOne(ctx context.Context, req *overlays.GetByIdRequest) (
	*overlays.Overlay, error,
) {
	dashboardId := ctx.Value("dashboardId").(string)

	entity := model.ChannelOverlay{}
	if err := c.Db.
		WithContext(ctx).
		Preload("Layers").
		Find(&entity, "channel_id = ? AND id = ?", dashboardId, req.GetId()).
		Error; err != nil {
		return nil, err
	}

	if entity.ID.String() == "" {
		return nil, twirp.NewError(twirp.NotFound, "not found")
	}

	return c.convertEntity(entity), nil
}

const tenMinutes = 10 * 60

func (c *Overlays) OverlaysUpdate(ctx context.Context, req *overlays.UpdateRequest) (
	*overlays.Overlay, error,
) {
	dashboardId := ctx.Value("dashboardId").(string)

	err := c.validateReq(req.Name, req.Width, req.Height, req.Layers)
	if err != nil {
		return nil, err
	}

	var entity model.ChannelOverlay
	if err := c.Db.
		WithContext(ctx).
		Preload("Layers").
		Find(&entity, "channel_id = ? AND id = ?", dashboardId, req.Id).
		Error; err != nil {
		return nil, err
	}

	if entity.ID.String() == "" {
		return nil, twirp.NewError(twirp.NotFound, "not found")
	}

	err = c.Db.Transaction(
		func(tx *gorm.DB) error {
			entity.Name = req.Name
			entity.UpdatedAt = time.Now()

			if err := tx.Where(
				"overlay_id = ?",
				entity.ID,
			).Delete(&model.ChannelOverlayLayer{}).Error; err != nil {
				return err
			}

			entity.Layers = nil

			for _, l := range req.Layers {
				if l.Settings.HtmlOverlayHtmlDataPollSecondsInterval < 5 ||
					l.Settings.HtmlOverlayHtmlDataPollSecondsInterval > tenMinutes {
					return twirp.NewError(twirp.InvalidArgument, "invalid poll interval")
				}

				layer := model.ChannelOverlayLayer{
					ID:   uuid.New(),
					Type: c.convertToDbType(l.Type),
					Settings: model.ChannelOverlayLayerSettings{
						HtmlOverlayHTML:                    textToBase64(l.Settings.HtmlOverlayHtml),
						HtmlOverlayCSS:                     textToBase64(l.Settings.HtmlOverlayCss),
						HtmlOverlayJS:                      textToBase64(l.Settings.HtmlOverlayJs),
						HtmlOverlayDataPollSecondsInterval: int(l.Settings.HtmlOverlayHtmlDataPollSecondsInterval),
					},
					OverlayID:               entity.ID,
					PosX:                    int(l.PosX),
					PosY:                    int(l.PosY),
					Width:                   int(l.Width),
					Height:                  int(l.Height),
					PeriodicallyRefetchData: l.PeriodicallyRefetchData,
				}

				if err := tx.Save(&layer).Error; err != nil {
					return err
				}
			}

			return tx.Save(&entity).Error
		},
	)

	if err != nil {
		return nil, err
	}

	_, err = c.Grpc.Websockets.RefreshOverlaySettings(
		ctx,
		&websockets.RefreshOverlaysRequest{ChannelId: dashboardId},
	)
	if err != nil {
		c.Logger.Error("failed to refresh overlays", logger.Error(err))
	}

	return c.OverlaysGetOne(ctx, &overlays.GetByIdRequest{Id: req.Id})
}

func (c *Overlays) OverlaysDelete(ctx context.Context, req *overlays.DeleteRequest) (
	*emptypb.Empty, error,
) {
	dashboardId := ctx.Value("dashboardId").(string)

	if err := c.Db.
		WithContext(ctx).
		Delete(&model.ChannelOverlay{}, "channel_id = ? and id = ?", dashboardId, req.Id).
		Error; err != nil {
		return nil, err
	}

	return &emptypb.Empty{}, nil
}

func (c *Overlays) OverlaysCreate(ctx context.Context, req *overlays.CreateRequest) (
	*overlays.Overlay, error,
) {
	dashboardId := ctx.Value("dashboardId").(string)

	err := c.validateReq(req.Name, req.Width, req.Height, req.Layers)
	if err != nil {
		return nil, err
	}

	entity := model.ChannelOverlay{
		ID:        uuid.New(),
		ChannelID: dashboardId,
		Name:      req.Name,
	}

	err = c.Db.Transaction(
		func(tx *gorm.DB) error {
			if err := tx.Save(&entity).Error; err != nil {
				return err
			}

			for _, l := range req.Layers {
				if l.Settings.HtmlOverlayHtmlDataPollSecondsInterval < 5 ||
					l.Settings.HtmlOverlayHtmlDataPollSecondsInterval > tenMinutes {
					return twirp.NewError(twirp.InvalidArgument, "invalid poll interval")
				}

				layer := model.ChannelOverlayLayer{
					ID:   uuid.New(),
					Type: c.convertToDbType(l.Type),
					Settings: model.ChannelOverlayLayerSettings{
						HtmlOverlayHTML:                    textToBase64(l.Settings.HtmlOverlayHtml),
						HtmlOverlayCSS:                     textToBase64(l.Settings.HtmlOverlayCss),
						HtmlOverlayJS:                      textToBase64(l.Settings.HtmlOverlayJs),
						HtmlOverlayDataPollSecondsInterval: int(l.Settings.HtmlOverlayHtmlDataPollSecondsInterval),
					},
					OverlayID:               entity.ID,
					PosX:                    int(l.PosX),
					PosY:                    int(l.PosY),
					Width:                   int(l.Width),
					Height:                  int(l.Height),
					PeriodicallyRefetchData: l.PeriodicallyRefetchData,
				}

				if err := tx.Save(&layer).Error; err != nil {
					return err
				}
			}

			return nil
		},
	)

	if err != nil {
		return nil, err
	}

	return c.OverlaysGetOne(ctx, &overlays.GetByIdRequest{Id: entity.ID.String()})
}

func (c *Overlays) OverlaysParseHtml(ctx context.Context, req *overlays.ParseHtmlOverlayRequest) (
	*overlays.ParseHtmlOverlayResponse, error,
) {
	dashboardId, err := helpers.GetSelectedDashboardIDFromContext(ctx)
	if err != nil {
		return nil, err
	}

	res, err := c.Bus.Parser.ParseVariablesInText.Request(
		ctx, parser.ParseVariablesInTextRequest{
			ChannelID: dashboardId,
			Text:      base64ToText(req.GetHtml()),
		},
	)
	if err != nil {
		return nil, err
	}

	return &overlays.ParseHtmlOverlayResponse{
		Html: textToBase64(res.Data.Text),
	}, nil
}
