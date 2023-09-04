package timers

import (
	"context"
	"time"

	"github.com/samber/lo"
	"github.com/satont/twir/apps/scheduler/internal/types"
	model "github.com/satont/twir/libs/gomodels"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

func NewWatched(ctx context.Context, services *types.Services) {
	timeTick := lo.If(services.Config.AppEnv != "production", 15*time.Second).Else(5 * time.Minute)
	ticker := time.NewTicker(timeTick)

	go func() {
		for {
			select {
			case <-ctx.Done():
				ticker.Stop()
				return
			case <-ticker.C:
				incrementWatchedExpr := services.Gorm.
					Model(&model.UsersStats{}).
					Where(`"userId" IN (?)`, services.Gorm.Table("users_online").Select(`"userId"`))

				err := incrementWatchedExpr.
					WithContext(ctx).
					Update("watched", gorm.Expr("watched + ?", timeTick.Milliseconds())).Error
				if err != nil {
					zap.S().Error(err)
				}
			}
		}
	}()
}
