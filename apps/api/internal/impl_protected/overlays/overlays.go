package overlays

import (
	"context"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/google/uuid"
	"github.com/satont/twir/apps/api/internal/impl_deps"
	model "github.com/satont/twir/libs/gomodels"
	"github.com/satont/twir/libs/grpc/generated/api/overlays"
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

func (c *Overlays) convertEntity(entity model.ChannelOverlay) *overlays.Overlay {
	id := entity.ID.String()

	layers := make([]*overlays.OverlayLayer, len(entity.Layers))
	for i, l := range entity.Layers {
		layers[i] = &overlays.OverlayLayer{
			Id:   l.ID.String(),
			Type: c.convertToRpcType(l.Type),
			Settings: &overlays.OverlayLayerSettings{
				HtmlOverlayHtml:                        l.Settings.HtmlOverlayHTML,
				HtmlOverlayCss:                         l.Settings.HtmlOverlayCSS,
				HtmlOverlayJs:                          l.Settings.HtmlOverlayJS,
				HtmlOverlayHtmlDataPollSecondsInterval: int32(l.Settings.HtmlOverlayDataPollSecondsInterval),
			},
			OverlayId: id,
			PosX:      int32(l.PosX),
			PosY:      int32(l.PosY),
			Width:     int32(l.Width),
			Height:    int32(l.Height),
			CreatedAt: fmt.Sprint(l.CreatedAt.UnixMilli()),
			UpdatedAt: fmt.Sprint(l.UpdatedAt.UnixMilli()),
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

	if utf8.RuneCountInString(req.Name) > 50 {
		return nil, twirp.NewError(twirp.InvalidArgument, "name is too long")
	}

	if req.Width < 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "width must be positive")
	}

	if req.Height < 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "height must be positive")
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

	err := c.Db.Transaction(
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
						HtmlOverlayHTML:                    l.Settings.HtmlOverlayHtml,
						HtmlOverlayCSS:                     l.Settings.HtmlOverlayCss,
						HtmlOverlayJS:                      l.Settings.HtmlOverlayJs,
						HtmlOverlayDataPollSecondsInterval: int(l.Settings.HtmlOverlayHtmlDataPollSecondsInterval),
					},
					OverlayID: entity.ID,
					PosX:      int(l.PosX),
					PosY:      int(l.PosY),
					Width:     int(l.Width),
					Height:    int(l.Height),
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

	if utf8.RuneCountInString(req.Name) > 50 {
		return nil, twirp.NewError(twirp.InvalidArgument, "name is too long")
	}

	if req.Width < 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "width must be positive")
	}

	if req.Height < 0 {
		return nil, twirp.NewError(twirp.InvalidArgument, "height must be positive")
	}

	entity := model.ChannelOverlay{
		ID:        uuid.New(),
		ChannelID: dashboardId,
		Name:      req.Name,
	}

	err := c.Db.Transaction(
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
						HtmlOverlayHTML:                    l.Settings.HtmlOverlayHtml,
						HtmlOverlayCSS:                     l.Settings.HtmlOverlayCss,
						HtmlOverlayJS:                      l.Settings.HtmlOverlayJs,
						HtmlOverlayDataPollSecondsInterval: int(l.Settings.HtmlOverlayHtmlDataPollSecondsInterval),
					},
					OverlayID: entity.ID,
					PosX:      int(l.PosX),
					PosY:      int(l.PosY),
					Width:     int(l.Width),
					Height:    int(l.Height),
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
